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

package health

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/elastic/beats/libbeat/tests/compose"
	mbtest "github.com/elastic/beats/metricbeat/mb/testing"
	"github.com/elastic/beats/metricbeat/module/traefik/mtest"
)

func makeBadRequest(config map[string]interface{}) error {
	host := config["hosts"].([]string)[0]

	resp, err := http.Get("http://" + host + "/foobar")
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

func TestHealth(t *testing.T) {
	mtest.Runner.Run(t, compose.Suite{
		"Fetch": func(t *testing.T, r compose.R) {
			config := mtest.GetConfig("health", r.Host())

			makeBadRequest(config)

			ms := mbtest.NewReportingMetricSetV2(t, config)
			reporter := &mbtest.CapturingReporterV2{}

			ms.Fetch(reporter)
			assert.Nil(t, reporter.GetErrors(), "Errors while fetching metrics")

			event := reporter.GetEvents()[0]
			assert.NotNil(t, event)
			t.Logf("%s/%s event: %+v", ms.Module().Name(), ms.Name(), event)

			responseCount, _ := event.MetricSetFields.GetValue("response.count")
			assert.True(t, responseCount.(int64) >= 1)

			badResponseCount, _ := event.MetricSetFields.GetValue("response.status_codes.404")
			assert.True(t, badResponseCount.(float64) >= 1)
		},
		"Data": func(t *testing.T, r compose.R) {
			ms := mbtest.NewReportingMetricSetV2(t, mtest.GetConfig("health", r.Host()))
			err := mbtest.WriteEventsReporterV2(ms, t, "")
			if err != nil {
				t.Fatal("write", err)
			}
		},
	})
}
