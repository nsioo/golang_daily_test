package logs

import (
	"code.byted.org/gopkg/env"
	"code.byted.org/gopkg/metrics"
)

var (
	metricsClient *metrics.MetricsClientV2

	metricsTagWarn  = []metrics.T{{Name: "level", Value: "WARNING"}, {Name: "cluster", Value: env.Cluster()}}
	metricsTagError = []metrics.T{{Name: "level", Value: "ERROR"}, {Name: "cluster", Value: env.Cluster()}}
	metricsTagFatal = []metrics.T{{Name: "level", Value: "CRITICAL"}, {Name: "cluster", Value: env.Cluster()}} // 和py统一, 将fatal打成critical
	metricsLim      = 4                                                                                        //  只打Warn及以上的日志,
)

func init() {
	metricsClient = metrics.NewDefaultMetricsClientV2("toutiao.service.log", true)
}

// FIXME: it's strange to do metrics in logs SDK
// TODO: Use LogHook for metrics, refer: https://github.com/sirupsen/logrus/blob/master/hooks.go#L8
func doMetrics(logLevel int, psm string) error {
	if logLevel < metricsLim {
		return nil
	}

	if len(psm) == 0 {
		return nil
	}

	var err error
	if logLevel == 4 { // warning
		err = metricsClient.EmitCounter(psm+".throughput", 1, metricsTagWarn...)
	} else if logLevel == 5 { // error
		err = metricsClient.EmitCounter(psm+".throughput", 1, metricsTagError...)
	} else if logLevel == 6 { // fatal
		err = metricsClient.EmitCounter(psm+".throughput", 1, metricsTagFatal...)
	}
	return err
}
