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

	_ "github.com/blang/semver/v4"
	"github.com/go-kit/log"
	_ "github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
)

const xlogLocationSubsystem = "xlog_location"

type PGXlogLocationCollector struct {
}

var (
	xlogLocationBytes = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, xlogLocationSubsystem, "bytes"),
		"Postgres LSN (log sequence number) being generated on primary or replayed on replica (truncated to low 52 bits)",
		[]string{},
		prometheus.Labels{},
	)

	xlogLocationQuery = `
	SELECT CASE
		WHEN pg_is_in_recovery() THEN (pg_last_xlog_replay_location() - '0/0') % (2^52)::bigint
		ELSE (pg_current_xlog_location() - '0/0') % (2^52)::bigint
	END AS bytes
	`
)

func (PGXlogLocationCollector) Name() string {
	return "xlog location collector"
}

func (PGXlogLocationCollector) Help() string {
	return "xlog location information"
}

func (PGXlogLocationCollector) Version() float64 {
	return 1.0
}

func (PGXlogLocationCollector) Scrape(ctx context.Context, db *sql.DB, ch chan<- prometheus.Metric,logger log.Logger) error {

	// xlog was renmaed to WAL in PostgreSQL 10
	// https://wiki.postgresql.org/wiki/New_in_postgres_10#Renaming_of_.22xlog.22_to_.22wal.22_Globally_.28and_location.2Flsn.29

	rows, err := db.QueryContext(ctx,
		xlogLocationQuery)

	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var bytes float64

		if err := rows.Scan(&bytes); err != nil {
			return err
		}

		ch <- prometheus.MustNewConstMetric(
			xlogLocationBytes,
			prometheus.GaugeValue,
			bytes,
		)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	return nil
}
