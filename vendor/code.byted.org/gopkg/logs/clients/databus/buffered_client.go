package databus

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang/protobuf/proto"
)

const (
	DEFAULT_QUEUE_SIZE   = 1000
	DEFAULT_CONSUMER_NUM = 1
)

var ErrBufferBull = errors.New("buffer is full")

type BufferedDatabusCollector struct {
	queued       chan []byte
	socketPath   string
	writeTimeout time.Duration
	maxConn      int
}

func NewDefaultBufferedDatabusCollector() *BufferedDatabusCollector {
	return NewBufferedDatabusCollector(DEFAULT_SOCKET_PATH, DEFAULT_TIMEOUT, DEFAULT_MAX_CONN_NUM, DEFAULT_QUEUE_SIZE, DEFAULT_CONSUMER_NUM)
}

func NewBufferedDatabusCollector(socketPath string, timeout time.Duration, maxConn int, queueSize int, consumerNum int) *BufferedDatabusCollector {
	bc := &BufferedDatabusCollector{
		queued:       make(chan []byte, queueSize),
		socketPath:   socketPath,
		writeTimeout: timeout,
		maxConn:      maxConn,
	}
	for i := 0; i < consumerNum; i++ {
		go bc.consumer()
	}
	return bc
}

func (bc *BufferedDatabusCollector) consumer() {
	collector := NewCollector(bc.socketPath, bc.writeTimeout, 2)
	for {
		select {
		case msg, ok := <-bc.queued:
			if !ok {
				break
			}
			collector.SendProto(msg)
		}
	}
}

func (bc *BufferedDatabusCollector) CollectArray(channel string, messages []*ApplicationMessage) error {
	var payload RequestPayload
	payload.Channel = &channel
	payload.Messages = messages
	buf, err := proto.Marshal(&payload)
	if err != nil {
		return err
	}
	if len(buf) > PACKET_SIZE_LIMIT {
		return fmt.Errorf("CollectArray Err: Packet too large.")
	}
	select {
	case bc.queued <- buf:
		return nil
	default:
		return ErrBufferBull
	}
}
func (bc *BufferedDatabusCollector) Collect(channel string, value []byte, key []byte, codec int32) error {
	var message ApplicationMessage
	message.Key = key
	message.Value = value
	message.Codec = &codec
	var payload RequestPayload
	payload.Channel = &channel
	payload.Messages = []*ApplicationMessage{&message}
	buf, err := proto.Marshal(&payload)
	if err != nil {
		return err
	}

	select {
	case bc.queued <- buf:
		return nil
	default:
		return ErrBufferBull
	}
}
func (bc *BufferedDatabusCollector) Close() error {
	//close(bc.queued)
	return nil
}
