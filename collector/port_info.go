package collector

import (
	"context"
	"database/sql"

	_ "dbmanager"

	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	port = "port"
	portInfoQuery = `SELECT setting FROM pg_settings WHERE name = 'port'`;
)

var (
	/* 定义metrics名称 */
	portInfoDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace,port,"info"),
		"gaussDB port information",
		[]string{},
		nil,
	)

)

type PortInfoMetrics struct {}

func (PortInfoMetrics) Name() string{
	/* 实现Name函数 */
	return "PortInfoMetrics"
}

func (PortInfoMetrics) Help() string{
	/* 实现Help函数*/
	return "Metrtics Example"
}

func (PortInfoMetrics) Version() float64 {
	return 1.0
}

func (PortInfoMetrics) Scrape(ctx context.Context,db *sql.DB,ch chan <- prometheus.Metric,logger log.Logger) error {
	/* 这里实现抓取逻辑 */
	var portInfo int
	err := db.QueryRowContext(ctx,portInfoQuery).Scan(&portInfo)
	if err != nil {
		return err
	}

	if portInfo == 0 {
		return nil
	}

	masterLogRows,err := db.QueryContext(ctx,portInfoQuery)
	if err != nil {
		return err
	}
	defer masterLogRows.Close()

	columns,err := masterLogRows.Columns()
	if err != nil {
		return err
	}

	columnCount := len(columns)

	ch <- prometheus.MustNewConstMetric(
		portInfoDesc,prometheus.GaugeValue,float64(columnCount),
	)
	return nil
}
