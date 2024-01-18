package exporter

import (
	_ "fmt"
	_ "log"
	_ "dbmanager"
	"github.com/prometheus/client_golang/prometheus"
	_ "github.com/prometheus/client_golang/prometheus/promhttp"
)

type metrics_operation interface {
	Register()
}

type TestMetrics struct {
	cpuTemp		prometheus.Gauge
	registry	*prometheus.Registry
}

func (t *TestMetrics) Register() {
	t.registry = prometheus.NewRegistry()
	t.registry.MustRegister(t.cpuTemp)
}

func (t *TestMetrics) SetCpuTemp(number float64) {
	t.cpuTemp.Set(number)
}

func (t *TestMetrics) GetRegistry() *prometheus.Registry {
	return t.registry
}

func CreateTestMetrics(name string,help string) *TestMetrics {
	return &TestMetrics {
		cpuTemp: prometheus.NewGauge(prometheus.GaugeOpts {
			Name: name,
			Help: help,
		}),
	}
	
}
