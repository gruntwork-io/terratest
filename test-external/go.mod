module github.com/gruntwork-io/terratest/test-external

go 1.24.0

require (
	github.com/gruntwork-io/terratest/modules/logger v0.0.0
	github.com/gruntwork-io/terratest/modules/terragrunt v0.0.0
	github.com/gruntwork-io/terratest/modules/testing v0.0.0
	github.com/stretchr/testify v1.11.1
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/gruntwork-io/terratest/internal/lib v0.0.0 // indirect
	github.com/gruntwork-io/terratest/modules/retry v0.0.0 // indirect
	github.com/gruntwork-io/terratest/modules/shell v0.0.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/net v0.47.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/gruntwork-io/terratest/internal/lib => ../internal/lib
	github.com/gruntwork-io/terratest/modules/files => ../modules/files
	github.com/gruntwork-io/terratest/modules/logger => ../modules/logger
	github.com/gruntwork-io/terratest/modules/random => ../modules/random
	github.com/gruntwork-io/terratest/modules/retry => ../modules/retry
	github.com/gruntwork-io/terratest/modules/shell => ../modules/shell
	github.com/gruntwork-io/terratest/modules/terragrunt => ../modules/terragrunt
	github.com/gruntwork-io/terratest/modules/testing => ../modules/testing
)
