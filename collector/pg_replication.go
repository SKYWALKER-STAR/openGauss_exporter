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

const replicationSubsystem = "replication"


type PGReplicationCollector struct {
}

var (
	pgReplicationLag = prometheus.NewDesc(
		prometheus.BuildFQName(
			namespace,
			replicationSubsystem,
			"lag_seconds",
		),
		"Replication lag behind master in seconds",
		[]string{}, nil,
	)
	pgReplicationIsReplica = prometheus.NewDesc(
		prometheus.BuildFQName(
			namespace,
			replicationSubsystem,
			"is_replica",
		),
		"Indicates if the server is a replica",
		[]string{}, nil,
	)

	pgReplicationQuery = `SELECT
	CASE
		WHEN NOT pg_is_in_recovery() THEN 0
                WHEN pg_last_wal_receive_lsn () = pg_last_wal_replay_lsn () THEN 0
		ELSE GREATEST (0, EXTRACT(EPOCH FROM (now() - pg_last_xact_replay_timestamp())))
	END AS lag,
	CASE
		WHEN pg_is_in_recovery() THEN 1
		ELSE 0
"	END as is_replica`
)
func (PGReplicationCollector) Name() string {
	return "Replication collector"
}

func (PGReplicationCollector) Help() string {
	return "Replication information"
}

func (PGReplicationCollector) Version() float64 {
	return 1.0
}

func (PGReplicationCollector) Scrape(ctx context.Context, db *sql.DB, ch chan<- prometheus.Metric, logger log.Logger) error {
	row := db.QueryRowContext(ctx,
		pgReplicationQuery,
	)

	var lag float64
	var isReplica int64
	err := row.Scan(&lag, &isReplica)
	if err != nil {
		return err
	}
	ch <- prometheus.MustNewConstMetric(
		pgReplicationLag,
		prometheus.GaugeValue, lag,
	)
	ch <- prometheus.MustNewConstMetric(
		pgReplicationIsReplica,
		prometheus.GaugeValue, float64(isReplica),
	)
	return nil
}
