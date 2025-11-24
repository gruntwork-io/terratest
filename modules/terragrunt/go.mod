module github.com/gruntwork-io/terratest/modules/terragrunt

go 1.24.0

require (
	github.com/gruntwork-io/terratest/internal/lib v0.1.0
	github.com/gruntwork-io/terratest/modules/files v0.1.0
	github.com/gruntwork-io/terratest/modules/logger v0.1.0
	github.com/gruntwork-io/terratest/modules/retry v0.1.0
	github.com/gruntwork-io/terratest/modules/shell v0.1.0
	github.com/gruntwork-io/terratest/modules/testing v0.1.0
	github.com/stretchr/testify v1.11.1
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/mattn/go-zglob v0.0.2-0.20190814121620-e3c945676326 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/net v0.47.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/gruntwork-io/terratest/internal/lib => ../../internal/lib
	github.com/gruntwork-io/terratest/modules/files => ../files
	github.com/gruntwork-io/terratest/modules/logger => ../logger
	github.com/gruntwork-io/terratest/modules/random => ../random
	github.com/gruntwork-io/terratest/modules/retry => ../retry
	github.com/gruntwork-io/terratest/modules/shell => ../shell
	github.com/gruntwork-io/terratest/modules/testing => ../testing
)
