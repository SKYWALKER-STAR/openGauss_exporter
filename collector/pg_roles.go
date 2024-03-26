// Copyright 2024 The Prometheus Authors
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

const (
	rolesSubsystem = "roles"
)

var (
	pgRolesConnectionLimitsDesc = prometheus.NewDesc(
		prometheus.BuildFQName(
			namespace,
			rolesSubsystem,
			"connection_limit",
		),
		"Connection limit set for the role",
		[]string{"rolname"}, nil,
	)

	pgRolesConnectionLimitsQuery = "SELECT pg_roles.rolname, pg_roles.rolconnlimit FROM pg_roles"
)

type PGRolesCollector struct {}

func (PGRolesCollector) Name() string {
	return "Connection limit"
}

func (PGRolesCollector) Help() string {
	return "PG Role connection limit"
}

func (PGRolesCollector) Version() float64 {
	return 1.0
}

// Update implements Collector and exposes roles connection limits.
// It is called by the Prometheus registry when collecting metrics.
func (PGRolesCollector) Scrape(ctx context.Context,db *sql.DB,ch chan <- prometheus.Metric,logger log.Logger) error {
	// Query the list of databases
	rows, err := db.QueryContext(ctx,
		pgRolesConnectionLimitsQuery,
	)

	if err != nil {
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var rolname sql.NullString
		var connLimit sql.NullInt64
		if err := rows.Scan(&rolname, &connLimit); err != nil {
			return err
		}

		if !rolname.Valid {
			continue
		}
		rolnameLabel := rolname.String

		if !connLimit.Valid {
			continue
		}
		connLimitMetric := float64(connLimit.Int64)

		ch <- prometheus.MustNewConstMetric(
			pgRolesConnectionLimitsDesc,
			prometheus.GaugeValue, connLimitMetric, rolnameLabel,
		)
	}

	return rows.Err()
}
