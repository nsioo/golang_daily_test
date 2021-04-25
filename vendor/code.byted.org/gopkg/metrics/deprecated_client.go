package metrics

import (
	"fmt"
	"sync"
)

// Deprecated: use MetricsClientV2
type MetricsClient struct {
	mu sync.RWMutex

	NamespacePrefix string
	AllMetrics      map[string]map[string]metricsType
	Server          string
	IgnoreCheck     bool

	t *tagCache
	c *metricCache

	warnN int64
}

// Deprecated: use NewMetricsClientV2
func NewMetricsClient(server, namespacePrefix string, ignoreCheck bool) *MetricsClient {
	ckey := fmt.Sprintf("%v|%v|%v", server, namespacePrefix, ignoreCheck)
	clientsMu.Lock()
	defer clientsMu.Unlock()

	if cli := v1Clients[ckey]; cli != nil {
		return cli
	}

	cli := &MetricsClient{
		NamespacePrefix: namespacePrefix,
		AllMetrics:      make(map[string]map[string]metricsType),
		Server:          server,
		IgnoreCheck:     ignoreCheck,
	}
	cli.t = newTagCache()
	cli.c = newMetricCache(newSender(server))

	v1Clients[ckey] = cli
	return cli
}

// Deprecated: use NewDefaultMetricsClientV2
func NewDefaultMetricsClient(namespacePrefix string, ignoreCheck bool) *MetricsClient {
	return NewMetricsClient(DefaultMetricsServer, namespacePrefix, ignoreCheck)
}

func (mc *MetricsClient) DefineCounter(name, prefix string) error {
	return mc.defineMetrics(name, prefix, metricsTypeCounter)
}

func (mc *MetricsClient) DefineTimer(name, prefix string) error {
	return mc.defineMetrics(name, prefix, metricsTypeTimer)
}

func (mc *MetricsClient) DefineStore(name, prefix string) error {
	return mc.defineMetrics(name, prefix, metricsTypeStore)
}

func (mc *MetricsClient) defineMetrics(name, prefix string, mt metricsType) error {
	// mc.IgnoreCheck won't be modified, not need lock.
	if mc.IgnoreCheck {
		return nil
	}
	if len(prefix) == 0 {
		prefix = mc.NamespacePrefix
	}
	mc.mu.Lock()
	defer mc.mu.Unlock()
	m := mc.AllMetrics[prefix]
	if m == nil {
		m = make(map[string]metricsType)
		mc.AllMetrics[prefix] = m
	}
	t, ok := m[name]
	if !ok {
		m[name] = mt
		return nil
	}
	if mt != t {
		return errDuplicatedMetrics
	}
	return nil
}

func (mc *MetricsClient) EmitCounter(name string, value interface{}, prefix string, tagkv map[string]string) error {
	return mc.emit(metricsTypeCounter, name, value, prefix, tagkv)
}

func (mc *MetricsClient) EmitTimer(name string, value interface{}, prefix string, tagkv map[string]string) error {
	return mc.emit(metricsTypeTimer, name, value, prefix, tagkv)
}

func (mc *MetricsClient) EmitStore(name string, value interface{}, prefix string, tagkv map[string]string) error {
	return mc.emit(metricsTypeStore, name, value, prefix, tagkv)
}

// Flush sends any cached data to the metrics server
func (m *MetricsClient) Flush() {
	m.c.Flush()
}

const maxWarnN = 3

func (m *MetricsClient) emit(mt metricsType, name string, value interface{},
	prefix string, tagkv map[string]string) error {

	//if atomic.LoadInt64(&m.warnN) < maxWarnN { // fast check without invalidating cpu cache
	//	if atomic.AddInt64(&m.warnN, 1) <= maxWarnN {
	//		logfunc("[WARN] gopkg/metrics: [type=%s prefix=%q name=%q] MetricsClient is deprecated, use MetricsClientV2, check README.md for help.", mt, prefix, name)
	//	}
	//}

	if len(prefix) == 0 {
		prefix = m.NamespacePrefix
	}
	if !m.IgnoreCheck {
		m.mu.RLock()
		types, ok1 := m.AllMetrics[prefix]
		t, ok2 := types[name] // read from nil is safe
		m.mu.RUnlock()
		if !ok1 || !ok2 {
			return errEmitUndefinedMetrics
		}
		if t != mt {
			return errEmitBadMetricsType
		}
	}
	v, err := toFloat64(value)
	if err != nil {
		return err
	}
	if mt == metricsTypeCounter && v == 0 { // meaningless
		return nil
	}
	tags := Map2Tags(tagkv)
	return m.c.Send(m.t.MakeMetricEntity(mt, prefix, name, v, 0, tags))
}

// If you use the default metricsClient, then the NamespacePrefix is "",
// so you can fill in "prefix" when using DefineCounter, DefineTimer etc.
// and EmitCounter, EmitTimer etc.
// default metrics client won't ignore metrics check.
var metricsClient = NewDefaultMetricsClient("", false)

func DefineCounter(name, prefix string) error {
	return metricsClient.DefineCounter(name, prefix)
}

func DefineTimer(name, prefix string) error {
	return metricsClient.DefineStore(name, prefix)
}

func DefineStore(name, prefix string) error {
	return metricsClient.DefineStore(name, prefix)
}

func EmitCounter(name string, value interface{}, prefix string, tagkv map[string]string) error {
	return metricsClient.EmitCounter(name, value, prefix, tagkv)
}

func EmitTimer(name string, value interface{}, prefix string, tagkv map[string]string) error {
	return metricsClient.EmitTimer(name, value, prefix, tagkv)
}

func EmitStore(name string, value interface{}, prefix string, tagkv map[string]string) error {
	return metricsClient.EmitStore(name, value, prefix, tagkv)
}
