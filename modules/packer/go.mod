module github.com/gruntwork-io/terratest/modules/packer

go 1.24.0

require (
	github.com/gruntwork-io/terratest/modules/logger v1.0.0
	github.com/gruntwork-io/terratest/modules/retry v1.0.0
	github.com/gruntwork-io/terratest/modules/shell v1.0.0
	github.com/gruntwork-io/terratest/modules/testing v1.0.0
	github.com/hashicorp/go-multierror v1.1.1
	github.com/hashicorp/go-version v1.7.0
	github.com/stretchr/testify v1.11.1
)

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/hashicorp/errwrap v1.0.0 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	golang.org/x/net v0.47.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/gruntwork-io/terratest/modules/logger => ../logger
	github.com/gruntwork-io/terratest/modules/retry => ../retry
	github.com/gruntwork-io/terratest/modules/shell => ../shell
	github.com/gruntwork-io/terratest/modules/testing => ../testing
)
