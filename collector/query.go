// Copyright 2018 The Prometheus Authors
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
	"fmt"
	_ "context"
	_ "strings"
	"net/http"
	"database/sql"

	_ "gitee.com/opengauss/openGauss-connector-go-pq"
)

/*
type QueryTest struct {}
func (q *QueryTest) ServeHTTP(w http.ResponseWriter,r *http.Request) {
	str := "opengauss://exporter:7erMtAlaN3mXmfKwhUcY&@127.0.0.1:5432/postgres?sslmode=disable"
	db, err := sql.Open("opengauss", str)

	if err != nil {
		fmt.Println("Error from Query Test")
	}
	defer db.Close()


}
*/

func QueryTest(w http.ResponseWriter,r *http.Request) {
	str := "opengauss://exporter:7erMtAlaN3mXmfKwhUcY&@127.0.0.1:5432/exportertest?sslmode=disable"
	query := "select * from testtable"

	db, err := sql.Open("opengauss", str)
	if err != nil {
		fmt.Println("Error from Query Test")
	}
	defer db.Close()

	rows,errQuery := db.Query(query)
	if errQuery != nil {
		fmt.Println(errQuery)
	}

	defer rows.Close()

	for rows.Next() {
		var (
			id int64
			name string
		)
		if err := rows.Scan(&id,&name); err != nil {
			fmt.Println(err)
		}

		fmt.Printf("id %d's name is %s\n",id,name)
	}
}
