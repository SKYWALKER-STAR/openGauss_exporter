package collector

import (
	"context"
	"database/sql"
	"github.com/go-kit/log"
	_ "gitee.com/opengauss/openGauss-connector-go-pq"
	"github.com/prometheus/client_golang/prometheus"
)

type Scraper interface {

	// Name of the Scraper. Should be unique
	Name() string

	// Help describes the role of the Scraper
	// Example: "Collect from SHOW ENGINE INNODB STATUS"
	Help() string

	// Version of GaussDB from which scraper is available
	Version() float64

	// Scrape collects data from database connection and sends it over channel as prometheus metric.
	Scrape(ctx context.Context,db *sql.DB,ch chan <- prometheus.Metric,logger log.Logger) error
}
