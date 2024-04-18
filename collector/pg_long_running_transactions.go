// Copyright 2023 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package collector

import (
	"context"
	"database/sql"

	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
)

const longRunningTransactionsSubsystem = "long_running_transactions"

type PGLongRunningTransactionsCollector struct {
}

var (
	longRunningTransactionsCount = prometheus.NewDesc(
		"pg_long_running_transactions",
		"Current number of long running transactions",
		[]string{},
		prometheus.Labels{},
	)

	longRunningTransactionsAgeInSeconds = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, longRunningTransactionsSubsystem, "oldest_timestamp_seconds"),
		"The current maximum transaction age in seconds",
		[]string{},
		prometheus.Labels{},
	)

	longRunningTransactionsQuery = `
	SELECT
		COUNT(*) as transactions,
   		MAX(EXTRACT(EPOCH FROM clock_timestamp())) AS oldest_timestamp_seconds
    FROM pg_catalog.pg_stat_activity
    WHERE state is distinct from 'idle' AND query not like 'autovacuum:%'
	`
)
func (PGLongRunningTransactionsCollector) Name() string {
	return "longRunningTransactions"
}

func (PGLongRunningTransactionsCollector) Help() string {
	return "longRunningTransactions information"
}

func (PGLongRunningTransactionsCollector) Version() float64 {
	return 1.0
}

func (PGLongRunningTransactionsCollector) Scrape(ctx context.Context, db *sql.DB, ch chan<- prometheus.Metric,logger log.Logger) error {
	rows, err := db.QueryContext(ctx,
		longRunningTransactionsQuery)

	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var transactions, ageInSeconds float64

		if err := rows.Scan(&transactions, &ageInSeconds); err != nil {
			return err
		}

		ch <- prometheus.MustNewConstMetric(
			longRunningTransactionsCount,
			prometheus.GaugeValue,
			transactions,
		)
		ch <- prometheus.MustNewConstMetric(
			longRunningTransactionsAgeInSeconds,
			prometheus.GaugeValue,
			ageInSeconds,
		)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	return nil
}