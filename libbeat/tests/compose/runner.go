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

package compose

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/pkg/errors"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type TestRunner struct {
	Service  string
	Options  RunnerOptions
	Parallel bool
	Timeout  int
}

type Suite map[string]func(t *testing.T, r R)

type RunnerOptions map[string][]string

func (r *TestRunner) scenarios() []map[string]string {
	n := 1
	options := make(map[string][]string)
	for env, values := range r.Options {
		// Allow to override options from environment variables
		value := os.Getenv(env)
		if value != "" {
			values = []string{value}
		}
		options[env] = values
		n *= len(values)
	}

	scenarios := make([]map[string]string, n)
	for variable, values := range options {
		v := 0
		for i, s := range scenarios {
			if s == nil {
				s = make(map[string]string)
				scenarios[i] = s
			}
			s[variable] = values[v]
			v = (v + 1) % len(values)
		}
	}

	return scenarios
}

func (r *TestRunner) runSuite(t *testing.T, tests Suite, ctl R) {
	for name, test := range tests {
		t.Run(name, func(t *testing.T) { test(t, ctl) })
	}
}

func (r *TestRunner) runHostOverride(t *testing.T, tests Suite) bool {
	env := strings.ToUpper(r.Service) + "_HOST"
	host := os.Getenv(env)
	if host == "" {
		return false
	}

	t.Logf("Test host overriden by %s=%s", env, host)

	ctl := &runnerControl{
		host: host,
		t:    t,
	}
	r.runSuite(t, tests, ctl)
	return true
}

func (r *TestRunner) Run(t *testing.T, tests Suite) {
	t.Helper()

	if r.runHostOverride(t, tests) {
		return
	}

	timeout := r.Timeout
	if timeout == 0 {
		timeout = 300
	}

	scenarios := r.scenarios()
	if len(scenarios) == 0 {
		t.Fatal("Test runner configuration doesn't produce scenarios")
	}
	for _, s := range scenarios {
		var vars []string
		for k, v := range s {
			os.Setenv(k, v)
			vars = append(vars, fmt.Sprintf("%s=%s", k, v))
		}
		sort.Strings(vars)
		desc := strings.Join(vars, ",")

		seq := make([]byte, 16)
		rand.Read(seq)
		name := fmt.Sprintf("%s_%x", r.Service, seq)

		project, err := getComposeProject(name)
		if err != nil {
			t.Fatal(err)
		}
		project.SetParameters(s)

		t.Run(desc, func(t *testing.T) {
			if r.Parallel {
				t.Parallel()
			}

			err := project.Start(r.Service)
			// Down() is "idempotent", Start() has several points where it can fail,
			// so run Down() even if Start() fails.
			defer project.Down()
			if err != nil {
				t.Fatal(err)
			}

			err = project.Wait(timeout, r.Service)
			if err != nil {
				t.Fatal(errors.Wrapf(err, "waiting for %s/%s", r.Service, desc))
			}

			host, err := project.Host(r.Service)
			if err != nil {
				t.Fatal(errors.Wrapf(err, "getting host for %s/%s", r.Service, desc))
			}

			ctl := &runnerControl{
				host:     host,
				t:        t,
				scenario: s,
			}
			r.runSuite(t, tests, ctl)
		})

	}
}

type R interface {
	Host() string
}

type runnerControl struct {
	t        *testing.T
	host     string
	scenario map[string]string
}

func (r *runnerControl) Host() string {
	return r.host
}

func (r *runnerControl) Hostname() string {
	hostname, _, _ := net.SplitHostPort(r.host)
	return hostname
}

func (r *runnerControl) Port() string {
	_, port, _ := net.SplitHostPort(r.host)
	return port
}
