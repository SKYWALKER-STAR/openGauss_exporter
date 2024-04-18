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

const postmasterSubsystem = "postmaster"

type PGPostmasterCollector struct {
}

var (
	pgPostMasterStartTimeSeconds = prometheus.NewDesc(
		prometheus.BuildFQName(
			namespace,
			postmasterSubsystem,
			"start_time_seconds",
		),
		"Time at which postmaster started",
		[]string{}, nil,
	)

	pgPostmasterQuery = "SELECT extract(epoch from pg_postmaster_start_time) from pg_postmaster_start_time();"
)

func (PGPostmasterCollector) Name() string {
	return "PG Postmaster collector"
}

func (PGPostmasterCollector) Help() string {
	return "PG Postmaster information"
}

func (PGPostmasterCollector) Version() float64 {
	return 1.0
}

func (PGPostmasterCollector) Scrape(ctx context.Context, db *sql.DB, ch chan<- prometheus.Metric,logger log.Logger) error {
	row := db.QueryRowContext(ctx,
		pgPostmasterQuery)

	var startTimeSeconds sql.NullFloat64
	err := row.Scan(&startTimeSeconds)
	if err != nil {
		return err
	}
	startTimeSecondsMetric := 0.0
	if startTimeSeconds.Valid {
		startTimeSecondsMetric = startTimeSeconds.Float64
	}
	ch <- prometheus.MustNewConstMetric(
		pgPostMasterStartTimeSeconds,
		prometheus.GaugeValue, startTimeSecondsMetric,
	)
	return nil
}