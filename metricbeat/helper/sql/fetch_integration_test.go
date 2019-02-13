// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

// +build integration

package sql

import (
	"database/sql"
	"fmt"
	"math/rand"
	"testing"
	"time"

	//_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/elastic/beats/libbeat/tests/compose"
	mbmysql "github.com/elastic/beats/metricbeat/module/mysql"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func TestFetchWithMySQL(t *testing.T) {
	compose.EnsureUp(t, "mysql")

	// XXX: Remove dependencies on mysql module
	db, err := sql.Open("mysql", mbmysql.GetMySQLEnvDSN())
	require.NoError(t, err)

	testFetch(t, db)
}

func testFetch(t *testing.T, db *sql.DB) {
	dbName := fmt.Sprintf("testdb%d", rand.Uint64())
	tableName := fmt.Sprintf("%s.%s", dbName, "test")
	defer db.Exec(fmt.Sprintf("DROP DATABASE %s", dbName))
	_, err := db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
	require.NoError(t, err)
	_, err = db.Exec(fmt.Sprintf("CREATE TABLE %s (a1 integer, a2 varchar(10))", tableName))
	require.NoError(t, err)
	_, err = db.Exec(fmt.Sprintf("INSERT INTO %s VALUES (0, 'foo'), (1, 'bar')", tableName))
	require.NoError(t, err)

	cases := []struct {
		title    string
		query    string
		args     []interface{}
		expected []map[string]interface{}
	}{
		{
			title: "all values",
			query: fmt.Sprintf("SELECT * FROM %s", tableName),
			expected: []map[string]interface{}{
				{
					"a1": 0,
					"a2": "foo",
				},
				{
					"a1": 1,
					"a2": "bar",
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.title, func(t *testing.T) {
			result, err := Fetch(db, c.query, c.args...)
			require.NoError(t, err)
			assert.Equal(t, result, c.expected)
		})
	}
}
