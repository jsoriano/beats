package terraform

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/magefile/mage/mg"

	devtools "github.com/elastic/beats/dev-tools/mage"
)

// Terraform is a set of targets to manage cloud resources with terraform
type Terraform mg.Namespace

// Apply applies a terraform configuration
func (t Terraform) Apply() error {
	fmt.Println(">> terraform: Applying configuration")
	return t.apply()
}

// AWS is the equivalent of apply with TARGETS=module.aws
func (t Terraform) AWS() error {
	fmt.Println(">> terraform: Applying configuration for AWS resorces")
	return t.apply("module.aws")
}

// Destroy destroys a terraform configuration
func (t Terraform) Destroy() error {
	fmt.Println(">> terraform: Destroying configuration")
	tf := devtools.NewTerraform()
	err := tf.Destroy()
	return errors.Wrap(err, "destroying terraform configuration")
}

func (t Terraform) apply(targets ...string) error {
	tf := devtools.NewTerraform()
	tf.Targets = append(tf.Targets, targets...)

	if len(tf.Targets) == 0 {
		return errors.New("some target needs to be specified with TERRAFORM_TARGETS environment variable")
	}

	if err := tf.Init(); err != nil {
		return errors.Wrap(err, "initializing terraform")
	}

	if err := tf.Apply(); err != nil {
		return errors.Wrap(err, "applying terraform configuration")
	}

	return nil
}
