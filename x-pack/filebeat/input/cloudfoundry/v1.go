// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

package cloudfoundry

import (
	"github.com/pkg/errors"

	v2 "github.com/elastic/beats/v7/filebeat/input/v2"
	"github.com/elastic/beats/v7/filebeat/input/v2/input-stateless"
	"github.com/elastic/beats/v7/x-pack/libbeat/common/cloudfoundry"
)

// inputV1 defines a udp input to receive event on a specific host:port.
type inputV1 struct {
	config cloudfoundry.Config
}

func configureV1(config cloudfoundry.Config) (*inputV1, error) {
	return &inputV1{config: config}, nil
}

func (i *inputV1) Name() string { return "cloudfoundry-v1" }

func (i *inputV1) Test(_ v2.TestContext) error {
	return nil
}

func (i *inputV1) Run(ctx v2.Context, publisher stateless.Publisher) error {
	log := ctx.Logger
	hub := cloudfoundry.NewHub(&i.config, "filebeat", log)

	log.Info("Starting cloudfoundry input")
	defer log.Info("Stopped cloudfoundry input")

	callbacks := cloudfoundry.DopplerCallbacks{
		Log: func(evt cloudfoundry.Event) {
			publisher.Publish(createEvent(evt))
		},
		Error: func(evt cloudfoundry.EventError) {
			publisher.Publish(createEvent(&evt))
		},
	}

	consumer, err := hub.DopplerConsumer(callbacks)
	if err != nil {
		return errors.Wrapf(err, "initializing doppler consumer")
	}

	consumer.Run()
	<-ctx.Cancelation.Done()
	consumer.Stop()
	return nil
}
