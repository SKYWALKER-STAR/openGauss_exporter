package collector

import (
	"context"
	"database/sql"

	_ "dbmanager"

	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
)


const (
	test_metrics = "testForFun"
)

var (
	gausstest = prometheus.NewDesc(
		prometheus.BuildFQName(namespace,test_metrics,"Second"),
		"This is just for fun",
		[]string{},
		nil,
	)
)
type TestMetrics struct {}

func (TestMetrics) Name() string {
	return "testForFun"
}

func (TestMetrics) Help() string {
	return "This is test for fun"
}

func (TestMetrics) Version() float64 {
	return 1.0
}

func (TestMetrics) Scrape(ctx context.Context,db *sql.DB,ch chan <- prometheus.Metric,logger log.Logger) error {
	ch <- prometheus.MustNewConstMetric(
		gausstest,prometheus.GaugeValue,5127,
	)

	return nil
}
