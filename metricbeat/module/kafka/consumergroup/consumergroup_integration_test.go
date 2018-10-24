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

package consumergroup

import (
	"context"
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"github.com/Shopify/sarama"
	"github.com/pkg/errors"

	"github.com/elastic/beats/libbeat/tests/compose"
	mbtest "github.com/elastic/beats/metricbeat/mb/testing"
)

const (
	kafkaDefaultHost = "localhost"
	kafkaDefaultPort = "9092"
)

func TestData(t *testing.T) {
	compose.EnsureUp(t, "kafka")

	ctx := context.Background()
	defer ctx.Done()

	c, err := startConsumer(t, ctx, "metricbeat-test")
	if err != nil {
		t.Fatal(errors.Wrap(err, "starting kafka consumer"))
	}
	defer c.Close()

	ms := mbtest.NewReportingMetricSetV2(t, getConfig())
	for retries := 0; retries < 3; retries++ {
		err = mbtest.WriteEventsReporterV2(ms, t, "")
		if err == nil {
			return
		}
		time.Sleep(500 * time.Millisecond)
	}
	t.Fatal("write", err)
}

type dummyConsumerGroupHandler struct{}

func (dummyConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (dummyConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }
func (h dummyConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		// Just consume it
		sess.MarkMessage(msg, "")
	}
	return nil
}

func startConsumer(t *testing.T, ctx context.Context, topic string) (io.Closer, error) {
	brokers := []string{getTestKafkaHost()}
	topics := []string{topic}

	config := sarama.NewConfig()
	// Minimal version required for consumer groups
	config.Version = sarama.V0_10_2_0
	cg, err := sarama.NewConsumerGroup(brokers, "test-group", config)
	if err != nil {
		return nil, errors.Wrap(err, "creating consumer group")
	}

	go func() {
		err = cg.Consume(ctx, topics, &dummyConsumerGroupHandler{})
		if err != nil {
			t.Logf("Failed to consume: %v", err)
		}
	}()
	return cg, nil
}

func getConfig() map[string]interface{} {
	return map[string]interface{}{
		"module":     "kafka",
		"metricsets": []string{"consumergroup"},
		"hosts":      []string{getTestKafkaHost()},
	}
}

func getTestKafkaHost() string {
	return fmt.Sprintf("%v:%v",
		getenv("KAFKA_HOST", kafkaDefaultHost),
		getenv("KAFKA_PORT", kafkaDefaultPort),
	)
}

func getenv(name, defaultValue string) string {
	return strDefault(os.Getenv(name), defaultValue)
}

func strDefault(a, defaults string) string {
	if len(a) == 0 {
		return defaults
	}
	return a
}
