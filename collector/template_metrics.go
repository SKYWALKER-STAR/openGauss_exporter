package collector

import (
	"context"
	"database/sql"

	_ "dbmanager"

	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	/* 定义metrics名称 */
)

type TemplateMetrics struct {}

func (TemplateMetrics) Name() string{
	/* 实现Name函数 */

	return "OK"
}

func (TemplateMetrics) Help() string{
	/* 实现Help函数*/

	return "OK"
}

func (TemplateMetrics) Version() float64 {
	return 1.0
}

func (TemplateMetrics) Scrape(ctx context.Context,db *sql.DB,ch chan <- prometheus.Metric,logger log.Logger) error {
	/* 这里实现抓取逻辑 */
	return nil
}
