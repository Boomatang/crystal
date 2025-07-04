package workflow

import "github.com/prometheus/client_golang/prometheus"

var (
	ActionTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "action_total",
			Help: "Total number of actions taken",
		},
		[]string{"action"},
	)
	ActionDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "action_duration_seconds",
			Help:    "Duration of an action set",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"action"},
	)
	NodeCount = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "Node_Count",
			Help: "Number of nodes in the topology",
		},
	)
)

func MustRegister() {
	prometheus.MustRegister(ActionDuration)
	prometheus.MustRegister(ActionTotal)
	prometheus.MustRegister(NodeCount)
}
