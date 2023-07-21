package balancer

import (
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	RequestDuration   *prometheus.HistogramVec
	ProxyDistribution *prometheus.CounterVec
}

func NewMetrics() *Metrics {
	return &Metrics{
		RequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: "balancer_app",
				Name:      "response_time_seconds",
				Help:      "Duration of requests.",
				Buckets:   []float64{0.1, 0.15, 0.2, 0.25, 0.3, 1, 3, 5},
			}, []string{"method"}),
		ProxyDistribution: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "balancer_app",
				Name:      "proxy_distribution",
				Help:      "Which proxy server proceed the method",
			},
			[]string{"method", "server", "success"},
		),
	}
}

func (m *Metrics) Register() error {
	if err := prometheus.Register(m.ProxyDistribution); err != nil {
		return err
	}

	if err := prometheus.Register(m.RequestDuration); err != nil {
		return err
	}
	return nil
}

func (m *Metrics) Stat(method, server string, startTime time.Time, err error) {
	isSuccess := fmt.Sprintf("%v", err == nil)
	m.ProxyDistribution.WithLabelValues(method, server, isSuccess).Inc()
	m.RequestDuration.WithLabelValues(method).Observe(time.Since(startTime).Seconds())
}
