/*
 * *************************************************************************************
 * Author: Ming
 * Create Date: 2024/1/16
 * Usage: Prometheus四种Metric类型的工厂，用户程序能够通过CreateMetric来创建相应的Metric
 * *************************************************************************************
 */

package metrics

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
)

/* 
 ***************************************************
 * Metrics接口，一组通用的Method
 ***************************************************
*/
type Interface_Metrics interface {
	SetName(name string)
	GetName() string
	SetHelp(help string)
	GetHelp() string
	GetObj() prometheus.Metric

}

/* 
 ***************************************************
 *所有Metric通用的属性
 ***************************************************
*/
type Metrics struct {
	name string
	help string
}

func (m *Metrics) SetName(name string) {
	m.name = name
}

func (m *Metrics) GetName() string {
	return m.name
}

func (m *Metrics) SetHelp(help string) {
	m.help = help
}

func (m *Metrics) GetHelp() string {
	return m.help
}

/* 
 ***************************************************
 *Counter类型的Metrics相关的指标及创建函数
 ***************************************************
*/

type Counter struct {
	Metrics
	obj prometheus.Counter
}

func (c *Counter) GetObj() prometheus.Metric {
	return c.obj
}

func NewCounter(name string,help string) Interface_Metrics {
	return &Counter {
		Metrics: Metrics {
			name: name,
			help: help,
		},
		obj: prometheus.NewCounter(prometheus.CounterOpts {
			Name: name,
			Help: help,
		}),
	}

}
/* 
 ***************************************************
 *Gauge类型的Metrics相关的指标及创建函数
 ***************************************************
*/
type Gauge struct {
	Metrics
	obj prometheus.Gauge
}

func (g *Gauge) GetObj() prometheus.Metric {
	return g.obj
}

func NewGauge(name string,help string) Interface_Metrics {
	return &Gauge {
		Metrics: Metrics {
			name: name,
			help: help,
		},

		obj: prometheus.NewGauge(prometheus.GaugeOpts {
			Name: name,
			Help: help,
		}),

	}
}

/* 
 ***************************************************
 *Histogram类型的Metrics相关的指标及创建函数
 ***************************************************
*/
type Histogram struct {
	Metrics
	obj prometheus.Histogram
}

func (g *Histogram) GetObj() prometheus.Metric {
	return g.obj
}

func NewHistogram(name string,help string) Interface_Metrics {
	return &Histogram {
		Metrics: Metrics {
			name: name,
			help: help,
		},

		obj: prometheus.NewHistogram(prometheus.HistogramOpts {
			Name: name,
			Help: help,
		}),
	}
}

/* 
 ***************************************************
 *Summary类型的Metrics相关的指标及创建函数
 ***************************************************
*/
type Summary struct {
	Metrics
	obj prometheus.Summary
}

func (g *Summary) GetObj() prometheus.Metric {
	return g.obj
}

func NewSummary(name string,help string) Interface_Metrics {
	return &Summary {
		Metrics: Metrics {
			name: name,
			help: help,
		},

		obj: prometheus.NewSummary(prometheus.SummaryOpts {
			Name: name,
			Help: help,
		}),
	}
}

/* 
 ***************************************************
 * Metric工厂
 ***************************************************
*/
func CreateMetric(metricType string,metricName string,metricHelp string) Interface_Metrics {

	switch metricType {
		case "Counter":
			return NewCounter(metricName,metricType)
		case "Gague":
			return NewGauge(metricName,metricType)
		case "Histogram":
			return NewHistogram(metricName,metricType)
		case "Summary":
			return NewSummary(metricName,metricType)
		default:
			fmt.Printf("Error:Unknown Metric Type:[%s]\n",metricType)
			return nil
	}
}
