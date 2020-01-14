package mage

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/pkg/errors"

	"github.com/elastic/beats/libbeat/tests/terraform"
)

// defaultDir is the default location of the terraform files inside the beats repository
var defaultDir = "./testing/terraform"

func NewTerraform() *terraform.Terraform {
	mg.Deps(checkTerraform)

	esBeatsDir, err := ElasticBeatsDir()
	if err != nil {
		panic(errors.Wrap(err, "failed determine libbeat dir location"))
	}

	tf := &terraform.Terraform{
		Dir:     filepath.Join(esBeatsDir, defaultDir),
		Stderr:  os.Stderr,
		Targets: targetsFromEnv(),
	}
	if mg.Verbose() {
		tf.Stdout = os.Stdout
	}

	return tf
}

func checkTerraform() error {
	err := sh.Run("terraform", "version")
	if err != nil {
		return errors.New("terraform command not available")
	}
	return nil
}

func targetsFromEnv() []string {
	targetsEnv := os.Getenv("TERRAFORM_TARGETS")
	if targetsEnv == "" {
		return nil
	}
	return strings.Split(targetsEnv, ",")
}
