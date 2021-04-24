package logs

import (
	"fmt"
	"os"
	"sync"
	"time"

	"code.byted.org/gopkg/logs/clients/databus"
	"code.byted.org/gopkg/metrics"
)

const (
	DATABUS_METRICS_PREFIX = "web.databus.collect"
	DATABUS_METRICS_SUCC   = "success.throughput"
	DATABUS_METRICS_ERR    = "error.throughput"

	DATABUS_DEFAULT_CHANNEL = "__LOG__"

	AUTO_BATCH_SIZE     = 30
	AUTO_BATCH_CAPACITY = databus.PACKET_SIZE_LIMIT
	AUTO_BATCH_TIMEOUT  = time.Millisecond * 500
)

type DatabusProvider struct {
	channel  string
	key      string
	level    int
	databusC *databus.DatabusCollector
	metricsC *metrics.MetricsClient

	autoBatch      [AUTO_BATCH_SIZE][]byte
	autoBatchLen   int
	autoBatchCount int
	lock           sync.Mutex
	last           time.Time
}

// Deprecated: the Databus provider is not supported.
func NewDatabusProvider(key string) *DatabusProvider {
	return &DatabusProvider{
		key:     key,
		channel: DATABUS_DEFAULT_CHANNEL,
	}
}

// Deprecated: the Databus provider is not supported
func NewDatabusProviderWithChannel(key, channel string) *DatabusProvider {
	return &DatabusProvider{
		key:     key,
		channel: channel,
	}
}

// Deprecated: the Databus provider is not supported
func (dp *DatabusProvider) Init() error {
	dp.databusC = databus.NewDefaultCollector()
	dp.metricsC = metrics.NewDefaultMetricsClient(DATABUS_METRICS_PREFIX, true)
	return nil
}

func (dp *DatabusProvider) SetLevel(l int) {
	dp.lock.Lock()
	defer dp.lock.Unlock()

	dp.level = l
}

func (dp *DatabusProvider) WriteMsg(msg string, level int) error {
	dp.lock.Lock()
	defer dp.lock.Unlock()
	if level < dp.level {
		return nil
	}

	var err error
	if dp.autoBatchCount == AUTO_BATCH_SIZE || // count limit
		dp.autoBatchLen+len(msg) > AUTO_BATCH_CAPACITY || // capacity limit
		dp.last.Add(AUTO_BATCH_TIMEOUT).Before(time.Now()) { // timeout limit
		err = dp.sendAndReset()
	}

	dp.autoBatch[dp.autoBatchCount] = []byte(msg)
	dp.autoBatchCount += 1
	dp.autoBatchLen += len(msg)

	return err
}

func (dp *DatabusProvider) Destroy() error {
	dp.lock.Lock()
	defer dp.lock.Unlock()
	dp.sendAndReset() // flush auto batch buffer
	return dp.databusC.Close()
}

func (dp *DatabusProvider) Flush() error {
	dp.lock.Lock()
	defer dp.lock.Unlock()
	return dp.sendAndReset()
}

func (dp *DatabusProvider) sendAndReset() error {
	msgs := make([]*databus.ApplicationMessage, 0, dp.autoBatchCount)
	for i := 0; i < dp.autoBatchCount; i++ {
		app := new(databus.ApplicationMessage)
		app.Key = []byte(dp.key)
		var zero int32 = 0
		app.Codec = &zero
		app.Value = dp.autoBatch[i]
		msgs = append(msgs, app)
	}

	err := dp.databusC.CollectArray(dp.channel, msgs)
	metricsStr := DATABUS_METRICS_SUCC
	if err != nil {
		metricsStr = DATABUS_METRICS_ERR
		fmt.Fprintf(os.Stderr, "DatabusProvider send err: %s\n", err)
	}
	dp.metricsC.EmitCounter(metricsStr, dp.autoBatchCount, "", map[string]string{"psm": dp.key})

	// reset
	dp.autoBatchCount = 0
	dp.autoBatchLen = 0
	dp.last = time.Now()
	return err
}
