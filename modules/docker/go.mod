module github.com/gruntwork-io/terratest/modules/docker

go 1.24.0

require (
	github.com/gruntwork-io/terratest/modules/collections v0.55.0
	github.com/gruntwork-io/terratest/modules/git v0.55.0
	github.com/gruntwork-io/terratest/modules/http-helper v0.55.0
	github.com/gruntwork-io/terratest/modules/logger v0.55.0
	github.com/gruntwork-io/terratest/modules/random v0.55.0
	github.com/gruntwork-io/terratest/modules/shell v0.55.0
	github.com/gruntwork-io/terratest/modules/testing v0.55.0
	github.com/hashicorp/go-multierror v1.1.1
	github.com/stretchr/testify v1.11.1
	gotest.tools/v3 v3.0.3
)

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/hashicorp/errwrap v1.0.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/gruntwork-io/terratest/modules/collections => ../collections
	github.com/gruntwork-io/terratest/modules/git => ../git
	github.com/gruntwork-io/terratest/modules/http-helper => ../http-helper
	github.com/gruntwork-io/terratest/modules/logger => ../logger
	github.com/gruntwork-io/terratest/modules/random => ../random
	github.com/gruntwork-io/terratest/modules/shell => ../shell
	github.com/gruntwork-io/terratest/modules/testing => ../testing
)
