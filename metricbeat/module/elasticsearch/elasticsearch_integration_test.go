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

package elasticsearch_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/tests/compose"
	"github.com/elastic/beats/metricbeat/helper/elastic"
	mbtest "github.com/elastic/beats/metricbeat/mb/testing"
	"github.com/elastic/beats/metricbeat/module/elasticsearch"
	_ "github.com/elastic/beats/metricbeat/module/elasticsearch/ccr"
	_ "github.com/elastic/beats/metricbeat/module/elasticsearch/cluster_stats"
	_ "github.com/elastic/beats/metricbeat/module/elasticsearch/index"
	_ "github.com/elastic/beats/metricbeat/module/elasticsearch/index_recovery"
	_ "github.com/elastic/beats/metricbeat/module/elasticsearch/index_summary"
	_ "github.com/elastic/beats/metricbeat/module/elasticsearch/ml_job"
	_ "github.com/elastic/beats/metricbeat/module/elasticsearch/node"
	_ "github.com/elastic/beats/metricbeat/module/elasticsearch/node_stats"
	_ "github.com/elastic/beats/metricbeat/module/elasticsearch/shard"
)

var metricSets = []string{
	"ccr",
	"cluster_stats",
	"index",
	"index_recovery",
	"index_summary",
	"ml_job",
	"node",
	"node_stats",
	"shard",
}

func TestElasticsearch(t *testing.T) {
	runner := compose.TestRunner{
		Service: "elasticsearch",
		Options: compose.RunnerOptions{
			"ELASTICSEARCH_VERSION": {
				// "7.0.0-alpha2",
				"6.5.4",
				"6.4.3",
				"6.3.2",
				"6.2.4",
				"5.6.14",
			},
		},
		Parallel: true,
	}

	runner.Run(t, compose.Suite{
		"Fetch": testFetch,
		"Data":  testData,
	})
}

func testFetch(t *testing.T, r compose.R) {
	if v := r.Option("ELASTICSEARCH_VERSION"); v == "6.2.4" || v == "5.6.14" {
		t.Skip("This test fails on this version")
	}

	host := r.Host()
	err := createIndex(host)
	assert.NoError(t, err)

	err = enableTrialLicense(host)
	assert.NoError(t, err)

	err = createMLJob(host)
	assert.NoError(t, err)

	err = createCCRStats(host)
	assert.NoError(t, err)

	for _, metricSet := range metricSets {
		t.Run(metricSet, func(t *testing.T) {
			checkSkip(t, metricSet, host)
			f := mbtest.NewReportingMetricSetV2(t, getConfig(metricSet, host))
			events, errs := mbtest.ReportingFetchV2(f)

			assert.Empty(t, errs)
			if !assert.NotEmpty(t, events) {
				t.FailNow()
			}
			t.Logf("%s/%s event: %+v", f.Module().Name(), f.Name(),
				events[0].BeatEvent("elasticsearch", metricSet).Fields.StringToPrint())
		})
	}
}

func testData(t *testing.T, r compose.R) {
	for _, metricSet := range metricSets {
		t.Run(metricSet, func(t *testing.T) {
			checkSkip(t, metricSet, r.Host())
			f := mbtest.NewReportingMetricSetV2(t, getConfig(metricSet, r.Host()))
			err := mbtest.WriteEventsReporterV2(f, t, metricSet)
			if err != nil {
				t.Fatal("write", err)
			}
		})
	}
}

// GetConfig returns config for elasticsearch module
func getConfig(metricset string, host string) map[string]interface{} {
	return map[string]interface{}{
		"module":                     elasticsearch.ModuleName,
		"metricsets":                 []string{metricset},
		"hosts":                      []string{host},
		"index_recovery.active_only": false,
	}
}

// createIndex creates and elasticsearch index in case it does not exit yet
func createIndex(host string) error {
	client := &http.Client{}

	if checkExists("http://" + host + "/testindex") {
		return nil
	}

	req, err := http.NewRequest("PUT", "http://"+host+"/testindex", nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("HTTP error %d: %s", resp.StatusCode, resp.Status)
	}

	return nil
}

// createIndex creates and elasticsearch index in case it does not exit yet
func enableTrialLicense(host string) error {
	client := &http.Client{}

	enableXPackURL := "/_xpack/license/start_trial?acknowledge=true"

	req, err := http.NewRequest("POST", "http://"+host+enableXPackURL, nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("could not enable trial license, response = %v", string(body))
	}

	return nil
}

func createMLJob(host string) error {

	mlJob, err := ioutil.ReadFile("ml_job/_meta/test/test_job.json")
	if err != nil {
		return err
	}

	jobURL := "/_xpack/ml/anomaly_detectors/total-requests"

	if checkExists("http://" + host + jobURL) {
		return nil
	}

	body, resp, err := httpPutJSON(host, jobURL, mlJob)

	if resp.StatusCode != 200 {
		return fmt.Errorf("HTTP error loading ml job %d: %s, %s", resp.StatusCode, resp.Status, string(body))
	}

	return nil
}

func createCCRStats(host string) error {
	err := setupCCRRemote(host)
	if err != nil {
		return err
	}

	err = createCCRLeaderIndex(host)
	if err != nil {
		return err
	}

	err = createCCRFollowerIndex(host)
	if err != nil {
		return err
	}

	return nil
}

func setupCCRRemote(host string) error {
	remoteSettings, err := ioutil.ReadFile("ccr/_meta/test/test_remote_settings.json")
	if err != nil {
		return err
	}

	settingsURL := "/_cluster/settings"
	_, _, err = httpPutJSON(host, settingsURL, remoteSettings)
	return err
}

func createCCRLeaderIndex(host string) error {
	leaderIndex, err := ioutil.ReadFile("ccr/_meta/test/test_leader_index.json")
	if err != nil {
		return err
	}

	indexURL := "/pied_piper"
	_, _, err = httpPutJSON(host, indexURL, leaderIndex)
	return err
}

func createCCRFollowerIndex(host string) error {
	followerIndex, err := ioutil.ReadFile("ccr/_meta/test/test_follower_index.json")
	if err != nil {
		return err
	}

	followURL := "/rats/_ccr/follow"
	_, _, err = httpPutJSON(host, followURL, followerIndex)
	return err
}

func checkExists(url string) bool {
	resp, err := http.Get(url)
	if err != nil {
		return false
	}
	resp.Body.Close()

	// Entry exists
	if resp.StatusCode == 200 {
		return true
	}
	return false
}

func checkSkip(t *testing.T, metricset string, host string) {
	if metricset != "ccr" {
		return
	}

	version, err := getElasticsearchVersion(host)
	if err != nil {
		t.Fatal("getting elasticsearch version", err)
	}

	isCCRStatsAPIAvailable := elastic.IsFeatureAvailable(version, elasticsearch.CCRStatsAPIAvailableVersion)

	if !isCCRStatsAPIAvailable {
		t.Skip("elasticsearch CCR stats API is not available until " + elasticsearch.CCRStatsAPIAvailableVersion.String())
	}
}

func getElasticsearchVersion(elasticsearchHostPort string) (*common.Version, error) {
	resp, err := http.Get("http://" + elasticsearchHostPort + "/")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data common.MapStr
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	version, err := data.GetValue("version.number")
	if err != nil {
		return nil, err
	}

	return common.NewVersion(version.(string))
}

func httpPutJSON(host, path string, body []byte) ([]byte, *http.Response, error) {
	req, err := http.NewRequest("PUT", "http://"+host+path, bytes.NewReader(body))
	if err != nil {
		return nil, nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	return body, resp, nil
}
