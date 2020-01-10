---
layout: collection-browser-doc
title: Terraform basic example
category: tests-for-examples
excerpt: >-
  The basic test written in GoLang to test Terraform.
tags: ["example"]
image: /assets/img/logos/terraform-logo.png
order: 300
nav_title: Examples
nav_title_link: /examples/
---

Full soure code can be found here: [terraform_basic_example_test.go](https://github.com/gruntwork-io/terratest/blob/master/test/terraform_basic_example_test.go).

This test fits to the [Terraform Basic Example]({{site.baseurl}}/examples/infrastructure-as-code-examples/basic-terraform/).

## Source code

```
package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// An example of how to test the simple Terraform module in examples/terraform-basic-example using Terratest.
func TestTerraformBasicExample(t *testing.T) {
	t.Parallel()

	expectedText := "test"
	expectedList := []string{expectedText}
	expectedMap := map[string]string{"expected": expectedText}

	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../examples/terraform-basic-example",

		// Variables to pass to our Terraform code using -var options
		Vars: map[string]interface{}{
			"example": expectedText,

			// We also can see how lists and maps translate between terratest and terraform.
			"example_list": expectedList,
			"example_map":  expectedMap,
		},

		// Variables to pass to our Terraform code using -var-file options
		VarFiles: []string{"varfile.tfvars"},

		// Disable colors in Terraform commands so its easier to parse stdout/stderr
		NoColor: true,
	}

	// At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer terraform.Destroy(t, terraformOptions)

	// This will run `terraform init` and `terraform apply` and fail the test if there are any errors
	terraform.InitAndApply(t, terraformOptions)

	// Run `terraform output` to get the values of output variables
	actualTextExample := terraform.Output(t, terraformOptions, "example")
	actualTextExample2 := terraform.Output(t, terraformOptions, "example2")
	actualExampleList := terraform.OutputList(t, terraformOptions, "example_list")
	actualExampleMap := terraform.OutputMap(t, terraformOptions, "example_map")

	// Verify we're getting back the outputs we expect
	assert.Equal(t, expectedText, actualTextExample)
	assert.Equal(t, expectedText, actualTextExample2)
	assert.Equal(t, expectedList, actualExampleList)
	assert.Equal(t, expectedMap, actualExampleMap)
}
```

## Step-by-step explanation

### 1. Reference variables

First, define variables that will be used to set up Terraform code and then to validate output:

```
expectedText := "test"
expectedList := []string{expectedText}
expectedMap := map[string]string{"expected": expectedText}
```

### 2. Terraform setup

Set up Terraform code with `terraformOptions`:

```
terraformOptions := &terraform.Options{
	// The path to where our Terraform code is located
	TerraformDir: "../examples/terraform-basic-example",

	// Variables to pass to our Terraform code using -var options
	Vars: map[string]interface{}{
		"example": expectedText,

		// We also can see how lists and maps translate between terratest and terraform.
		"example_list": expectedList,
		"example_map":  expectedMap,
	},

	// Variables to pass to our Terraform code using -var-file options
	VarFiles: []string{"varfile.tfvars"},

	// Disable colors in Terraform commands so its easier to parse stdout/stderr
	NoColor: true,
}
```

### 3. Infrastructure initialization

Next, initialize and deploy Terraform code. Following code runs `terraform init` and `terraform deploy`:

```
terraform.InitAndApply(t, terraformOptions)
```

### 4. Validation

Validate the Terraform output:

```
actualTextExample := terraform.Output(t, terraformOptions, "example")
actualTextExample2 := terraform.Output(t, terraformOptions, "example2")
actualExampleList := terraform.OutputList(t, terraformOptions, "example_list")
actualExampleMap := terraform.OutputMap(t, terraformOptions, "example_map")

// Verify we're getting back the outputs we expect
assert.Equal(t, expectedText, actualTextExample)
assert.Equal(t, expectedText, actualTextExample2)
assert.Equal(t, expectedList, actualExampleList)
assert.Equal(t, expectedMap, actualExampleMap)
```

### 5. Clean

Finally, clean up any resources that were created:

```
defer terraform.Destroy(t, terraformOptions)
```
