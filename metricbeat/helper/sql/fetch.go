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

package sql

import (
	"database/sql"
	"reflect"

	"github.com/elastic/beats/libbeat/common"
	s "github.com/elastic/beats/libbeat/common/schema"
	"github.com/pkg/errors"
)

// FetchWithSchema fetches the result of a query and applies an schema on each row
func FetchWithSchema(db *sql.DB, schema s.Schema, query string, args ...interface{}) ([]common.MapStr, error) {
	if schema == nil {
		return nil, errors.New("nil schema")
	}

	results, err := Fetch(db, query, args...)
	if err != nil {
		return nil, err
	}

	data := make([]common.MapStr, len(results))
	for i, r := range results {
		d, err := schema.Apply(r)
		if err != nil {
			return nil, err
		}
		data[i] = d
	}
	return data, nil
}

// Fetch fetches the result of a query
func Fetch(db *sql.DB, query string, args ...interface{}) ([]map[string]interface{}, error) {
	if db == nil {
		return nil, errors.New("database not initialized")
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get column names")
	}

	var results []map[string]interface{}
	for rows.Next() {
		values, err := prepareValues(rows)
		if err != nil {
			return nil, errors.Wrap(err, "failed to prepare values")
		}

		err = rows.Scan(values...)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan row")
		}

		row := make(map[string]interface{})
		for i, column := range columns {
			row[column] = values[i]
		}

		results = append(results, row)
	}
	return results, nil
}

// prepareValues gets the column names and an array ready to be used as value for Scan
func prepareValues(rows *sql.Rows) ([]interface{}, error) {
	if rows == nil {
		return nil, errors.New("rows not initialized")
	}

	types, err := rows.ColumnTypes()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get column types")
	}

	values := make([]interface{}, len(types))
	for i, t := range types {
		v := reflect.New(t.ScanType()).Interface()
		switch v.(type) {
		case *sql.NullInt64:
			var n int
			values[i] = &n
		case *sql.RawBytes:
			var s string
			values[i] = &s
		case sql.Scanner:
			values[i] = v
		default:
			values[i] = v
		}
	}

	return values, nil
}
