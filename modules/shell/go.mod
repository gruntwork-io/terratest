module github.com/gruntwork-io/terratest/modules/shell

go 1.25.0

require (
	github.com/gruntwork-io/terratest v0.46.16
	github.com/gruntwork-io/terratest/modules/logger v0.0.0
	github.com/gruntwork-io/terratest/modules/testing v0.0.0
	github.com/stretchr/testify v1.11.1
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rogpeppe/go-internal v1.13.1 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/gruntwork-io/terratest => ../..
	github.com/gruntwork-io/terratest/modules/logger => ../logger
	github.com/gruntwork-io/terratest/modules/testing => ../testing
)
