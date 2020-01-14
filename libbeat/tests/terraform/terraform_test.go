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
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var terraformDir = "./test/local"

func TestMain(m *testing.M) {
	// Cleanup terraform state files before and after running the tests
	cleanup()

	result := m.Run()

	cleanup()
	os.Exit(result)
}

func TestTerraform(t *testing.T) {
	tmpDir, err := ioutil.TempDir(os.TempDir(), "beats-terraform-")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	tf := Terraform{
		Dir: terraformDir,
		Vars: Vars{
			"install_dir": tmpDir,
		},
	}
	verbose(&tf)

	err = tf.Init()
	require.NoError(t, err)

	err = tf.Apply()
	require.NoError(t, err)

	path, err := tf.Output("file_path")
	require.NoError(t, err)

	d, err := ioutil.ReadFile(path)
	require.NoError(t, err)

	assert.Equal(t, "some content\n", string(d))

	err = tf.Destroy()
	require.NoError(t, err)

	_, err = os.Stat(path)
	assert.True(t, os.IsNotExist(err), "file should have been removed")
}

func TestTerraformWithTarget(t *testing.T) {
	tmpDir, err := ioutil.TempDir(os.TempDir(), "beats-terraform-")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	tf := Terraform{
		Dir: terraformDir,
		Vars: Vars{
			"install_dir": tmpDir,
		},
		Targets: Targets{
			"local_file.empty_file",
		},
	}
	verbose(&tf)

	err = tf.Init()
	require.NoError(t, err)

	err = tf.Apply()
	require.NoError(t, err)

	_, err = tf.Output("file_path")
	require.Error(t, err, "this output shouldn't exist")

	path, err := tf.Output("empty_file_path")
	require.NoError(t, err)

	_, err = os.Stat(path)
	assert.NoError(t, err, "empty file should exist")

	err = tf.Destroy()
	require.NoError(t, err)

	_, err = os.Stat(path)
	assert.True(t, os.IsNotExist(err), "file should have been removed")
}

func verbose(tf *Terraform) {
	if testing.Verbose() {
		tf.Stdout = os.Stdout
		tf.Stderr = os.Stderr
	}
}

func cleanup() {
	files := []string{
		// ".terraform",
		"terraform.tfstate",
		"terraform.tfstate.backup",
	}

	for _, f := range files {
		os.RemoveAll(filepath.Join(terraformDir, f))
	}
}
