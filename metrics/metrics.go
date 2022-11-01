package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	_counterMsg = promauto.NewCounterVec( //nolint:gochecknoglobals
		prometheus.CounterOpts{ //nolint:promlinter,exhaustruct
			Namespace: "tg",
			Subsystem: "routertext",
			Name:      "msg_count_ops_total",
		},
		[]string{"type"},
	)
	_summaryExecuteTime = promauto.NewSummaryVec( //nolint:gochecknoglobals
		prometheus.SummaryOpts{ //nolint:promlinter,exhaustruct
			Namespace: "tg",
			Subsystem: "routertext",
			Name:      "summary_execute_time_seconds",
			Objectives: map[float64]float64{
				0.5:  0.1,   //nolint:gomnd
				0.9:  0.01,  //nolint:gomnd
				0.99: 0.001, //nolint:gomnd
			},
		},
		[]string{"type"},
	)
)

func CounterMsgInc(name string) {
	_counterMsg.WithLabelValues(name).Inc()
}

func SummaryExecuteTimeObserve(name string, value float64) {
	_summaryExecuteTime.WithLabelValues(name).Observe(value)
}
