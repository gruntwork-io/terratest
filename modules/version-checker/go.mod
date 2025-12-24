module github.com/gruntwork-io/terratest/modules/version-checker/v2

go 1.24.0

require (
	github.com/gruntwork-io/terratest/modules/shell/v2 v2.0.0
	github.com/gruntwork-io/terratest/modules/terraform/v2 v2.0.0
	github.com/gruntwork-io/terratest/modules/testing/v2 v2.0.0
	github.com/hashicorp/go-version v1.7.0
	github.com/stretchr/testify v1.11.1
)

require (
	github.com/agext/levenshtein v1.2.3 // indirect
	github.com/apparentlymart/go-textseg/v15 v15.0.0 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/hashicorp/errwrap v1.0.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/hcl/v2 v2.24.0 // indirect
	github.com/hashicorp/terraform-json v0.27.2 // indirect
	github.com/jinzhu/copier v0.4.0 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/mitchellh/go-wordwrap v1.0.1 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	github.com/tmccombs/hcl2json v0.6.8 // indirect
	github.com/zclconf/go-cty v1.16.4 // indirect
	golang.org/x/mod v0.30.0 // indirect
	golang.org/x/sync v0.18.0 // indirect
	golang.org/x/text v0.31.0 // indirect
	golang.org/x/tools v0.39.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/gruntwork-io/terratest/modules/shell/v2 => ../shell
	github.com/gruntwork-io/terratest/modules/terraform/v2 => ../terraform
	github.com/gruntwork-io/terratest/modules/testing/v2 => ../testing
)
