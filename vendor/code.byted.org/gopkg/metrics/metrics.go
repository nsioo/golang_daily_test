// Package metrics provides a goroutine safe metrics client package metrics
// if TCE_HOST_IP is setted, will use this env value as host address
package metrics

import (
	"errors"
	"os"
	"strings"
	"sync"
	"time"

	"code.byted.org/gopkg/env"
)

type metricsType int

const (
	metricsTypeCounter metricsType = iota
	metricsTypeTimer
	metricsTypeStore
	metricsTypeTsStore
	metricsTypeRateCounter
	metricsTypeMeter
)

func (t metricsType) String() string {
	switch t {
	case metricsTypeCounter:
		return "counter"
	case metricsTypeStore:
		return "store"
	case metricsTypeTsStore:
		return "ts_store"
	case metricsTypeTimer:
		return "timer"
	case metricsTypeRateCounter:
		return "rate_counter"
	case metricsTypeMeter:
		return "meter"
	}
	return "unknown"
}

const (
	BlackholeAddr = "blackhole"

	asyncGoroutines = 4                      // the goroutine number of async tasks
	asyncChanBuffer = 1000                   // the buffer size of global chan for async tasks
	maxPendingSize  = 1000                   // the max size of caching metric entity , if the cache is full, will flush it to underlying connection via async chan
	flushInterval   = 200 * time.Millisecond // default interval of cli.Flush

	// DO NOT MODIFY IT IF YOU DONT KNOWN WHAT YOU ARE DOING
	maxBunchBytes = 32 << 10 // 32kb
)

var (
	DefaultMetricsServer = "127.0.0.1:9123"

	errDuplicatedMetrics    = errors.New("duplicated metrics name")
	errDuplicatedTag        = errors.New("duplicated metrics tag")
	errEmitUndefinedMetrics = errors.New("emit undefined metrics")
	errEmitBadMetricsType   = errors.New("emit bad metrics type")
	errEmitBufferFull       = errors.New("emit buffer full")
	errUnKnownValue         = errors.New("Unkown metrics value")
	errMetricsName          = errors.New("metrics name err")
	errTagName              = errors.New("metrics tag name err")
	errTagValue             = errors.New("metrics tag value err")
)

type globalTags struct {
	Tags []T
	Data []byte
	sync.RWMutex
}

var gTags globalTags

func (g *globalTags) AddTag(name, value string) {
	g.Lock()
	defer g.Unlock()
	g.removeTagWithoutLock(name)
	t := Tag(name, value)
	g.Tags = append(g.Tags, t)
	g.Data = appendTags(g.Data, []T{t})
}

func (g *globalTags) RemoveTag(name string) bool {
	g.Lock()
	defer g.Unlock()
	return g.removeTagWithoutLock(name)
}

func (g *globalTags) removeTagWithoutLock(name string) bool {
	tags := g.Tags[:0]
	for _, t := range g.Tags {
		if t.Name != name {
			tags = append(tags, t)
		}
	}
	if len(tags) == len(g.Tags) {
		return false
	}
	g.Data = appendTags(g.Data[:0], tags)
	g.Tags = tags
	return true
}

func (g *globalTags) Bytes() []byte {
	g.RLock()
	defer g.RUnlock()
	return g.Data
}

func (g *globalTags) Reset() {
	g.Lock()
	defer g.Unlock()
	g.Tags = g.Tags[:0]
	g.Data = g.Data[:0]
}

// AddGlobalTag adds a tag to global var, it effects all the metrics
func AddGlobalTag(name, value string) {
	gTags.AddTag(name, value)
}

// RemoveGlobalTag removes a tag from global var
func RemoveGlobalTag(name string) {
	gTags.RemoveTag(name)
}

// ResetGlobalTag clears all global tags
func ResetGlobalTag() {
	gTags.Reset()
}

func init() {
	initValidCharTable()
	if host := strings.TrimSpace(os.Getenv("TCE_HOST_IP")); host != "" {
		DefaultMetricsServer = host + ":9123"
		AddGlobalTag("env_type", "tce")
		AddGlobalTag("pod_name", env.PodName())
	}
	AddGlobalTag("_psm", os.Getenv("TCE_PSM"))
	AddGlobalTag("deploy_stage", os.Getenv("TCE_STAGE"))

	if env.Cluster() != "" {
		AddGlobalTag("cluster", env.Cluster())
	}

	if env.HasIPV6() {
		AddGlobalTag("host_v6", env.HostIPV6())
	}

	go flushClients()
}

var (
	clientsMu sync.RWMutex

	v1Clients = map[string]*MetricsClient{}
	v2Clients = map[string]*MetricsClientV2{}
)

func flushClients() {
	type Flusher interface {
		Flush()
	}
	for range time.Tick(flushInterval) {
		clientsMu.RLock()
		clients := make([]Flusher, 0, len(v1Clients)+len(v2Clients))
		for _, cli := range v1Clients {
			clients = append(clients, cli)
		}
		for _, cli := range v2Clients {
			if !cli.isFlushLoopRunning() { // if the flush loop is running, we leave flushing cache to MetricsClientV2
				clients = append(clients, cli)
			}
		}
		clientsMu.RUnlock()
		for _, cli := range clients { // call Flush without lock
			cli.Flush()
		}
	}
}
