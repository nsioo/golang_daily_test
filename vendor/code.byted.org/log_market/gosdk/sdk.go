package gosdk

import (
	"errors"
	"fmt"
	"net"
	"os"
	"code.byted.org/gopkg/metrics"
	"sync"
)

var (
	debugMode      bool
	raddr          *net.UnixAddr
	metricsClient  *metrics.MetricsClient
	network        = "unixpacket"
	socketPath     = "/opt/tmp/ttlogagent/unixpacket_v2.sock"
	ErrChannelFull = errors.New("[logagent-gosdk] channel full")
	ErrMsgNil      = errors.New("msg cannot be nil")
	ErrStop        = errors.New("gosdk had exited gracefully ")
	mu             sync.RWMutex
	senders        map[string]*Sender
)

func init() {
	var err error
	raddr, err = net.ResolveUnixAddr("unixpacket", socketPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[logagent-gosdk] init error:%v", err)
		return
	}
	mode := os.Getenv("__sdk_mode")
	if mode == "debug" {
		fmt.Println("[logagent-gosdk] use debug mode")
		debugMode = true
	}
	metricsClient = metrics.NewDefaultMetricsClient("toutiao.infra.ttlogagent.gosdk", true)
	senders = make(map[string]*Sender)
}

// preserve for compatibility
func Init() {
}

func Send(taskName string, m *Msg) error {
	mu.RLock()
	sender, ok := senders[taskName]
	mu.RUnlock()
	if !ok {
		mu.Lock()
		sender = NewDefaultSender(taskName)
		senders[taskName] = sender
		mu.Unlock()
	}
	if err := sender.Send(taskName, m); err != nil {
		printToStderrIfDebug("[logagent-gosdk] channel full")
		metricsClient.EmitCounter(taskName+".send", 1, "", asyncWriteErrorTag)
		return ErrChannelFull
	}
	return nil
}

func GracefullyExit() {
	mu.Lock()
	defer mu.Unlock()
	for name, s := range senders {
		s.Stop()
		delete(senders, name)
	}
}

func printToStderrIfDebug(msg string) {
	if !debugMode {
		return
	}
	fmt.Fprintln(os.Stderr, msg)
}
