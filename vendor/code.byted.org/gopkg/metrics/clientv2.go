package metrics

import (
	"fmt"
	"sync"
	"time"
)

type MetricsClientV2 struct {
	mu sync.RWMutex

	server  string
	prefix  string
	nocheck bool

	t *tagCache
	c *metricCache

	// it's fine with using slice, RegisterCounter is only called once for each counter,
	// and use slice here for improving the performace of iteration
	counters []*counterImpl
	timers   []*timerImpl

	flushLoopRunning bool

	metrictypes map[string]metricsType
}

func NewMetricsClientV2(server, prefix string, nocheck bool) *MetricsClientV2 {
	ckey := fmt.Sprintf("%v|%v|%v", server, prefix, nocheck)
	clientsMu.Lock()
	defer clientsMu.Unlock()

	if cli := v2Clients[ckey]; cli != nil {
		return cli
	}

	cli := &MetricsClientV2{server: server, prefix: prefix, nocheck: nocheck}
	cli.metrictypes = make(map[string]metricsType)
	cli.t = newTagCache()
	cli.c = newMetricCache(newSender(server))

	v2Clients[ckey] = cli
	return cli
}

func NewDefaultMetricsClientV2(prefix string, nocheck bool) *MetricsClientV2 {
	return NewMetricsClientV2(DefaultMetricsServer, prefix, nocheck)
}

// SetBlock indicates metrics entities will be sent to underlying connection directly if cache is full without using async chan
func (m *MetricsClientV2) SetBlock(v bool) {
	m.c.SetBlock(v)
}

// SetFlushInterval sets flush interval of cache.
// by default. clients are flushed by a global routine,
// with SetFlushInterval, it will disable the global one by setting flushLoopRunning=true, and has its own timer.
// NOTE: the method can only be called once
func (m *MetricsClientV2) SetFlushInterval(d time.Duration) {
	if d < time.Millisecond { // 0 or bug?
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.flushLoopRunning {
		return
	}
	m.flushLoopRunning = true
	go m.flushloop(d)
}

func (m *MetricsClientV2) isFlushLoopRunning() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.flushLoopRunning
}

func (m *MetricsClientV2) DefineCounter(name string) error {
	return m.defineMetrics(name, metricsTypeCounter)
}

func (m *MetricsClientV2) DefineRateCounter(name string) error {
	return m.defineMetrics(name, metricsTypeRateCounter)
}

// DefineMeter meter combines counter & rate_couter
// meter(m) = counter(m) + rate_couter(m.rate)
// requires metricserver2 above 1.0.0.65
func (m *MetricsClientV2) DefineMeter(name string) error {
	return m.defineMetrics(name, metricsTypeMeter)
}

func (m *MetricsClientV2) DefineTimer(name string) error {
	return m.defineMetrics(name, metricsTypeTimer)
}

func (m *MetricsClientV2) DefineStore(name string) error {
	return m.defineMetrics(name, metricsTypeStore)
}

func (m *MetricsClientV2) defineMetrics(name string, mt metricsType) error {
	if m.nocheck {
		return nil
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	t, ok := m.metrictypes[name]
	if !ok {
		m.metrictypes[name] = mt
		return nil
	}
	if mt != t {
		return errDuplicatedMetrics
	}
	return nil
}

func (m *MetricsClientV2) EmitCounter(name string, value interface{}, tags ...T) error {
	return m.emit(metricsTypeCounter, name, value, 0, tags...)
}

func (m *MetricsClientV2) EmitRateCounter(name string, value interface{}, tags ...T) error {
	return m.emit(metricsTypeRateCounter, name, value, 0, tags...)
}

func (m *MetricsClientV2) EmitMeter(name string, value interface{}, tags ...T) error {
	return m.emit(metricsTypeMeter, name, value, 0, tags...)
}

func (m *MetricsClientV2) EmitTimer(name string, value interface{}, tags ...T) error {
	return m.emit(metricsTypeTimer, name, value, 0, tags...)
}

func (m *MetricsClientV2) EmitStore(name string, value interface{}, tags ...T) error {
	return m.emit(metricsTypeStore, name, value, 0, tags...)
}

// EmitStoreWithTime is same as EmitStore except it emit store metrics with time
func (m *MetricsClientV2) EmitStoreWithTime(name string, value interface{}, t time.Time, tags ...T) error {
	if t.IsZero() {
		return m.emit(metricsTypeTsStore, name, value, 0, tags...)
	}
	return m.emit(metricsTypeTsStore, name, value, t.Unix(), tags...)
}

// Flush sends any cached data to the metrics server
func (m *MetricsClientV2) Flush() {
	m.flushCountersAndTimers()
	m.c.Flush()
}

func (m *MetricsClientV2) flushCountersAndTimers() {
	m.mu.RLock()
	counters := append(make([]*counterImpl, 0, len(m.counters)), m.counters...)
	timers := append(make([]*timerImpl, 0, len(m.timers)), m.timers...)
	m.mu.RUnlock()
	for _, c := range counters {
		if err := c.Flush(); err != nil {
			logfunc("[WARN] Flush %q err: %s", c.e.name)
		}
	}
	for _, t := range timers {
		if err := t.Flush(); err != nil {
			logfunc("[WARN] Flush %q err: %s", t.e.name)
		}
	}
}

func (m *MetricsClientV2) emit(mt metricsType, name string, value interface{}, ts int64, tags ...T) error {
	if !m.nocheck {
		m.mu.RLock()
		t, ok := m.metrictypes[name]
		m.mu.RUnlock()
		if !ok {
			return errEmitUndefinedMetrics
		}
		if t != mt {
			// we reuse DefineStore by metricsTypeStore for metricsTypeTsStore
			if t != metricsTypeStore || mt != metricsTypeTsStore {
				return errEmitBadMetricsType
			}
		}
	}
	v, err := toFloat64(value)
	if err != nil {
		return err
	}
	if v == 0 && (mt == metricsTypeCounter || mt == metricsTypeRateCounter) { // meaningless
		return nil
	}

	if !isValidString(name) {
		return errMetricsName
	}

	return m.c.Send(m.t.MakeMetricEntity(mt, m.prefix, name, v, ts, tags))
}

func (m *MetricsClientV2) flushloop(interval time.Duration) {
	for range time.Tick(interval) {
		m.Flush()
	}
}

func verifyNameAndTags(name string, tags []T) error {
	if !isValidString(name) {
		return errMetricsName
	}
	m := make(map[string]bool)
	for _, t := range tags {
		if !isValidString(t.Name) {
			return errTagName
		}
		if !isValidString(t.Value) {
			return errTagValue
		}
		if m[t.Name] {
			return errDuplicatedTag
		}
		m[t.Name] = true
	}
	return nil
}

// RegisterCounter registers a new counter according to name and tags. It use a int64 for counting and then sends the diff to tsdb.
// It returns err if name and tags combination is existed or any other errs detected, user have to save and reuse the return value
// the `maxPendingSize`(=1000) will NOT be applied to the counter, but `flushInterval`(=200ms) do, and it works well with cli.Flush.
func (m *MetricsClientV2) RegisterCounter(name string, tags ...T) (Counter, error) {
	return m.registerCounter(metricsTypeCounter, name, tags...)
}

// RegisterRateCounter registers a new 'rate_counter' type  counter. same as RegisterCounter, except the type of metrics.
func (m *MetricsClientV2) RegisterRateCounter(name string, tags ...T) (Counter, error) {
	return m.registerCounter(metricsTypeRateCounter, name, tags...)
}

func (m *MetricsClientV2) registerCounter(tp metricsType, name string, tags ...T) (Counter, error) {
	if err := verifyNameAndTags(name, tags); err != nil {
		return nil, err
	}
	c := &counterImpl{
		e: m.t.MakeMetricEntity(tp, m.prefix, name, 0, 0, tags),
		c: m.c,
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	key := c.e.Key()
	for i := range m.counters {
		if m.counters[i].e.Key() == key {
			return nil, errDuplicatedMetrics
		}
	}
	m.counters = append(m.counters, c)
	return c, nil
}

// RegisterTimers registers a new timer according to name and tags. It use [][]float64 for holding values in memory and then sends it all to tsdb.
// It returns err if name and tags combination is existed or any other errs detected, user have to save and reuse the return value.
// the `maxPendingSize`(=1000) and flushInterval`(=200ms) will be applied to the timer, and it works well with cli.Flush.
func (m *MetricsClientV2) RegisterTimer(name string, tags ...T) (Timer, error) {
	if err := verifyNameAndTags(name, tags); err != nil {
		return nil, err
	}
	t := &timerImpl{
		e: m.t.MakeMetricEntity(metricsTypeTimer, m.prefix, name, 0, 0, tags),
		c: m.c,
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	key := t.e.Key()
	for i := range m.timers {
		if m.timers[i].e.Key() == key {
			return nil, errDuplicatedMetrics
		}
	}
	m.timers = append(m.timers, t)
	return t, nil
}
