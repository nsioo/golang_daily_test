package metrics

import (
	"io"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"sync"
	"unsafe"
)

const (
	_emit = "emit"
)

var debugMetrics = os.Getenv("DEBUG_GOPKG_METRICS") != ""

var unixdomainsock = ""

func init() {
	for _, path := range []string{"/opt/tmp/sock/metric.sock", "/tmp/metric.sock"} {
		for i := 0; i < 3; i++ {
			conn, err := net.Dial("unixgram", path)
			if err == nil {
				conn.Close()
				unixdomainsock = path
				return
			}
			if strings.Contains(err.Error(), "no such file") {
				break
			}
		}
	}
}

type metricsWriter struct {
	addr string

	mu    sync.RWMutex
	conns []net.Conn
}

func (w *metricsWriter) connect() (net.Conn, error) {
	if strings.HasPrefix(w.addr, "/") {
		return net.Dial("unixgram", w.addr)
	} else {
		return net.Dial("udp", w.addr)
	}
}

func (w *metricsWriter) getconn() (net.Conn, error) {
	w.mu.Lock()
	if len(w.conns) > 0 {
		conn := w.conns[len(w.conns)-1]
		w.conns = w.conns[:len(w.conns)-1]
		w.mu.Unlock()
		return conn, nil
	}
	w.mu.Unlock()
	return w.connect()
}

func (w *metricsWriter) putconn(conn net.Conn) {
	w.mu.Lock()
	if len(w.conns) < asyncGoroutines {
		w.conns = append(w.conns, conn)
	} else {
		conn.Close()
	}
	w.mu.Unlock()
}

func (w *metricsWriter) Write(b []byte) (int, error) {
	if w.addr == BlackholeAddr {
		return len(b), nil
	}
	conn, err := w.getconn()
	if err != nil {
		if debugMetrics {
			logfunc("[DEBUG] gopkg/metrics: conn err: %s\n", err)
		}
		return 0, err
	}
	n, err := conn.Write(b)
	if err != nil {
		if debugMetrics {
			logfunc("[DEBUG] gopkg/metrics: write err: %s\n", err)
		}
		conn.Close()
	} else {
		w.putconn(conn)
	}
	return n, err
}

var writers struct {
	sync.Mutex
	m map[string]io.Writer // addr => metricsWriter
}

// try to reuse the metricsWriter as well as the connections in the metricsWriter
func getOrCreateWriter(addr string) io.Writer {
	if addr == BlackholeAddr {
		return ioutil.Discard
	}
	writers.Lock()
	defer writers.Unlock()
	if writers.m == nil {
		writers.m = make(map[string]io.Writer)
	}
	w := writers.m[addr]
	if w != nil {
		return w
	}
	w = &metricsWriter{addr: addr}
	writers.m[addr] = w
	return w
}

type sender struct {
	batch bool
	agg   bool

	w io.Writer
}

type senderInterface interface {
	SendCounter(ms []metricEntity)
	Send(ms []metricEntity)
}

func newSender(addr string) senderInterface {
	s := &sender{agg: true}
	if addr == DefaultMetricsServer && unixdomainsock != "" {
		addr = unixdomainsock
		s.batch = true
	}
	s.w = getOrCreateWriter(addr)
	return s
}

type aggregatekey struct {
	prefix string
	name   string
	tt     uintptr // tags pointer
	mt     metricsType
}

type counterAggregator struct {
	keys []aggregatekey
	m    map[aggregatekey]*metricEntity
}

var counterAggregatorPool = sync.Pool{
	New: func() interface{} {
		return &counterAggregator{
			keys: make([]aggregatekey, 0, maxPendingSize),
			m:    make(map[aggregatekey]*metricEntity, maxPendingSize),
		}
	},
}

func (a *counterAggregator) Merge(ms []metricEntity) []metricEntity {
	for i := range ms {
		e := &ms[i]
		k := aggregatekey{
			prefix: e.prefix,
			name:   e.name,
			tt:     uintptr(unsafe.Pointer(e.tt)),
			mt:     e.mt,
		}
		v, ok := a.m[k]
		if ok {
			v.v += e.v
		} else {
			a.m[k] = e
			a.keys = append(a.keys, k)
		}
	}
	p := ms[:0]
	for _, k := range a.keys {
		p = append(p, *a.m[k])
		delete(a.m, k)
	}
	a.keys = a.keys[:0]
	return p
}

func (s *sender) DisableCounterAggregator(v bool) {
	s.agg = !v
}

func (s *sender) SendCounter(ms []metricEntity) {
	if len(ms) == 0 {
		return
	}
	if s.agg {
		a := counterAggregatorPool.Get().(*counterAggregator)
		s.Send(a.Merge(ms))
		counterAggregatorPool.Put(a)
	} else {
		s.Send(ms)
	}
}

func printMs(ms []metricEntity) {
	logfunc("[DEBUG] gopkg/metrics: sending %d metrics to server:\n", len(ms))
	dup := make(map[aggregatekey]bool)
	for _, m := range ms {
		k := aggregatekey{prefix: m.prefix, name: m.name, mt: m.mt}
		if !dup[k] {
			logfunc("[DEBUG] %s '%s.%s'\n", m.mt, m.prefix, m.name)
			dup[k] = true
		}
	}
}

func (s *sender) Send(ms []metricEntity) {
	if debugMetrics {
		printMs(ms)
	}
	if !s.batch {
		p := wbufpool.Get().(*wbuf)
		defer wbufpool.Put(p)
		for _, m := range ms {
			s.w.Write(m.AppendTo(p.b[:0]))
		}
		return
	}
	// send bunch
	for len(ms) > 0 {
		ms = ms[s.sendbunch(ms):]
	}
}

type wbuf struct {
	b []byte

	mem [2 * maxBunchBytes]byte
}

var wbufpool = sync.Pool{
	New: func() interface{} {
		p := new(wbuf)
		p.b = p.mem[:0]
		return p
	},
}

func (s *sender) sendbunch(ms []metricEntity) int {
	if len(ms) == 0 {
		return 0
	}

	// limit to send maxBunchBytes
	k := 0
	n := msgpackArrayHeaderSize
	for i, m := range ms {
		n += m.MarshalSize()
		if n >= maxBunchBytes {
			// if the first metrics is larger than maxBunchBytes, we have nothing to send
			// we can only skip it by returning 1 or the caller will go into an infinite loop
			if i == 0 {
				logfunc("[ERROR] gopkg/metrics: %s '%s.%s' metrics too large\n", m.mt, m.prefix, m.name)
				return 1
			}
			break
		}
		k++
	}

	ms = ms[:k]

	p := wbufpool.Get().(*wbuf)
	defer wbufpool.Put(p)
	p.b = p.b[:0]

	// marshal to p
	p.b = msgpackAppendArrayHeader(p.b, uint16(len(ms)))
	for _, m := range ms {
		p.b = m.AppendTo(p.b)
	}
	s.w.Write(p.b)
	return k
}
