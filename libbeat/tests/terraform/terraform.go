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

package terraform

import (
	"bytes"
	"io"
	"os/exec"
	"strings"
)

type Terraform struct {
	Dir     string
	Vars    Vars
	Targets Targets

	Stdout io.Writer
	Stderr io.Writer
}

func (t *Terraform) Init() error {
	return t.terraformCmd("init").Run()
}

func (t *Terraform) Apply() error {
	args := t.buildArgs("apply", "-auto-approve")
	return t.terraformCmd(args...).Run()
}

func (t *Terraform) Destroy() error {
	args := t.buildArgs("destroy", "-auto-approve")
	return t.terraformCmd(args...).Run()
}

func (t *Terraform) Output(name string) (string, error) {
	var stdout bytes.Buffer
	cmd := t.terraformCmd("output", name)
	cmd.Stdout = &stdout
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(stdout.String()), nil
}

func (t *Terraform) buildArgs(args ...string) []string {
	vars := t.Vars.args()
	args = append(args, vars...)
	targets := t.Targets.args()
	args = append(args, targets...)
	return args
}

func (t *Terraform) terraformCmd(args ...string) *exec.Cmd {
	terraform := exec.Command("terraform", args...)
	terraform.Dir = t.Dir
	terraform.Stdout = t.Stdout
	terraform.Stderr = t.Stderr
	return terraform
}

type Vars map[string]string

func (v Vars) args() (args []string) {
	for name, value := range v {
		args = append(args, "-var="+name+"="+value)
	}
	return
}

type Targets []string

func (t Targets) args() (args []string) {
	args = make([]string, len(t))
	for i, target := range t {
		args[i] = "-target=" + target
	}
	return
}
