package exporter

import (
	_ "log"
	_ "dbmanager"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	_ "github.com/prometheus/client_golang/prometheus/promhttp"
)

type TemplateMetrics struct {
	cpuTemp		prometheus.Gauge
	registry	*prometheus.Registry
}

func (t *TemplateMetrics) Register() {
	t.registry = prometheus.NewRegistry()
	t.registry.MustRegister(t.cpuTemp)
}

func (t *TemplateMetrics) SetCpuTemp(number float64) {
	t.cpuTemp.Set(number)
}

func (t *TemplateMetrics) GetRegistry() *prometheus.Registry {
	return t.registry
}

func (t *TemplateMetrics) GetDataType() {
	fmt.Printf("%T\n",t.cpuTemp)
}

func CreateTemplateMetrics(name string,help string) *TemplateMetrics {
	return &TemplateMetrics {
		cpuTemp: prometheus.NewGauge(prometheus.GaugeOpts {
			Name: name,
			Help: help,
		}),
	}
}
