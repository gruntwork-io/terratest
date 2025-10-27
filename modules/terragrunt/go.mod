module github.com/gruntwork-io/terratest/modules/terragrunt

go 1.25.0

require (
	github.com/gruntwork-io/terratest/modules/files v0.0.0
	github.com/gruntwork-io/terratest/modules/formatting v0.0.0-00010101000000-000000000000
	github.com/gruntwork-io/terratest/modules/logger v0.0.0
	github.com/gruntwork-io/terratest/modules/retry v0.0.0-00010101000000-000000000000
	github.com/gruntwork-io/terratest/modules/shell v0.0.0
	github.com/gruntwork-io/terratest/modules/testing v0.0.0
	github.com/stretchr/testify v1.11.1
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/mattn/go-zglob v0.0.6 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/net v0.46.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/gruntwork-io/terratest => ../..
	github.com/gruntwork-io/terratest/modules/files => ../files
	github.com/gruntwork-io/terratest/modules/formatting => ../formatting
	github.com/gruntwork-io/terratest/modules/logger => ../logger
	github.com/gruntwork-io/terratest/modules/retry => ../retry
	github.com/gruntwork-io/terratest/modules/shell => ../shell
	github.com/gruntwork-io/terratest/modules/testing => ../testing
)
