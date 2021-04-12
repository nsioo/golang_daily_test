package gosdk

import (
	"fmt"
	"net"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gogo/protobuf/proto"

	"code.byted.org/log_market/gosdk/internal/ratelimit"
)

const (
	statusStop = iota
	statusRunning
)

var (
	sendSuccessTag     = map[string]string{"status": "success"}
	connRateLimitTag   = map[string]string{"status": "error", "reason": "connRateLimit"}
	connWriteErrorTag  = map[string]string{"status": "error", "reason": "connWriteError"}
	asyncWriteErrorTag = map[string]string{"status": "error", "reason": "asyncWrite"}
	connErrorTag       = map[string]string{"status": "error", "reason": "connError"}
	marshalErrorTag    = map[string]string{"status": "error", "reason": "marshalError"}
)

type SenderConfig struct {
	//unix socket buffer size is :  160K on 32-bit linux, 208K on 64-bit linux. PacketSizeLimit must be less than above.
	PacketSizeLimit uint
	//Sender is asynchronous, BatchFlushLatencyMs determinate the max wait time of message which in channel.
	BatchFlushLatencyMs uint
	//when write data to socket timeout
	SocketWriteTimeoutMs uint
	//Sender is asynchronous, ChannelSize
	ChannelSize uint
	//In general, TaskName is P.S.M
	TaskName string
	//the number of worker , one worker mean that one goroutine.
	WorkerNum uint
}

type Sender struct {
	SenderConfig
	currentBatchSizeByte uint
	connRetryLimit       *ratelimit.Bucket
	conn                 net.Conn
	ch                   chan *Msg
	lastFlushTime        time.Time
	batch                *MsgBatch
	exitCh               chan struct{}
	wg                   *sync.WaitGroup
	status               int32
}

type SenderOption func(*Sender)

func NewDefaultSenderConfig(taskName string) SenderConfig {
	return SenderConfig{
		//unix socket buffer size is :  160K on 32-bit linux, 208K on 64-bit linux. PacketSizeLimit must be less than above.
		//and < 32K packet will be more friendly for memory allocate.
		PacketSizeLimit:      30 * 1024,
		BatchFlushLatencyMs:  200,
		ChannelSize:          4096,
		SocketWriteTimeoutMs: 300,
		TaskName:             taskName,
		WorkerNum:            1,
	}
}
func NewDefaultSender(taskName string) *Sender {
	return newSender(NewDefaultSenderConfig(taskName))
}

func NewSender(config SenderConfig, options ...SenderOption) *Sender {
	s := newSender(config)
	for _, op := range options {
		op(s)
	}
	return s
}

func newSender(config SenderConfig) *Sender {
	s := &Sender{
		SenderConfig:   config,
		connRetryLimit: ratelimit.NewBucket(time.Second, 1),
		lastFlushTime:  time.Now(),
		batch: &MsgBatch{
			Msgs: make([]*Msg, 0, 256),
		},
		ch:     make(chan *Msg, config.ChannelSize),
		exitCh: make(chan struct{}),
		wg:     new(sync.WaitGroup),
		status: statusRunning,
	}
	s.wg.Add(int(s.WorkerNum))
	for i := 0; i < int(s.WorkerNum); i++ {
		go s.run()
	}
	return s
}

func (s *Sender) Send(taskName string, msg *Msg) error {
	if atomic.LoadInt32(&s.status) != statusRunning {
		return ErrStop
	}
	if msg == nil {
		return ErrMsgNil
	}
	if msg.Tags == nil {
		tags := make(map[string]string, 1)
		msg.Tags = tags
	}

	if msg.Tags["_taskName"] != taskName {
		msg.Tags["_taskName"] = taskName
	}

	select {
	case s.ch <- msg:
	default:
		return ErrChannelFull
	}
	return nil
}

func (s *Sender) Stop() {
	if !atomic.CompareAndSwapInt32(&s.status, statusRunning, statusStop) {
		return
	}
	close(s.exitCh)
	s.wg.Wait()
	if s.conn != nil {
		s.conn.Close()
	}
}

func (s *Sender) run() {
	if c, err := net.DialUnix(network, nil, raddr); err == nil {
		s.conn = c
	}
	ticker := time.NewTicker(time.Millisecond * time.Duration(s.BatchFlushLatencyMs))
	defer s.wg.Done()
	defer ticker.Stop()

	for {
		select {
		case msg := <-s.ch:
			s.appendMsg(msg)
		case <-ticker.C:
			if uint(time.Now().Sub(s.lastFlushTime))/uint(time.Millisecond) > s.BatchFlushLatencyMs {
				s.flush()
			}
		case <-s.exitCh:
			chanLen := len(s.ch)
			for i := 0; i < chanLen; i++ {
				select {
				case msg1 := <-s.ch:
					s.appendMsg(msg1)
				default:
				}
			}
			s.flush()
			return
		}
	}
}

func (s *Sender) appendMsg(msg *Msg) {
	size := uint(msg.Size())
	if len(msg.Msg) > int(s.PacketSizeLimit) {
		body := msg.Msg[:int(s.PacketSizeLimit)]
		msg = &Msg{
			Tags: msg.Tags,
			Msg:  body,
		}
		s.flush()
	}

	s.batch.Msgs = append(s.batch.Msgs, msg)
	s.currentBatchSizeByte += size
	if s.currentBatchSizeByte >= s.PacketSizeLimit {
		s.flush()
	}
}

func (s *Sender) sendToSocket(buf []byte) (error, map[string]string) {
	if s.conn == nil {
		if s.connRetryLimit.TakeAvailable(1) < 1 {
			return fmt.Errorf("[logagent-gosdk] build connection error=rate limit"), connRateLimitTag
		}
		c, err := net.DialUnix(network, nil, raddr)
		if err != nil {
			return fmt.Errorf("[logagent-gosdk] build connection error=%v", err), connErrorTag
		}
		s.conn = c
	}
	s.conn.SetWriteDeadline(time.Now().Add(time.Millisecond * time.Duration(s.SocketWriteTimeoutMs)))
	if _, err := s.conn.Write(buf); err != nil {
		s.conn.Close()
		s.conn = nil
		return fmt.Errorf("[logagent-gosdk] send  to socket failed, error=%s", err.Error()), connWriteErrorTag
	}
	return nil, nil
}

func (s *Sender) flush() {
	defer func() {
		s.currentBatchSizeByte = 0
		s.batch.Msgs = s.batch.Msgs[:0]
		s.lastFlushTime = time.Now()
	}()
	lenBatchMsg := len(s.batch.Msgs)
	if lenBatchMsg == 0 {
		return
	}

	buf, err := proto.Marshal(s.batch)
	if err != nil {
		printToStderrIfDebug(fmt.Sprintf("[logagent-gosdk] proto marashal batch error=%v", err))
		metricsClient.EmitCounter(s.TaskName+".send", lenBatchMsg, "", marshalErrorTag)
		return
	}
	if err, errorTag := s.sendToSocket(buf); err != nil {
		printToStderrIfDebug("first send to socket  failed :" + err.Error())
		if err, errorTag = s.sendToSocket(buf); err != nil {
			metricsClient.EmitCounter(s.TaskName+".send", lenBatchMsg, "", errorTag)
			if errorTag["reason"] == connRateLimitTag["reason"] {
				printToStderrIfDebug(err.Error())
			} else {
				fmt.Fprintf(os.Stderr, err.Error())
			}
			return
		}
	}
	metricsClient.EmitCounter(s.TaskName+".send", lenBatchMsg, "", sendSuccessTag)
}
