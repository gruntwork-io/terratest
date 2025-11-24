module github.com/gruntwork-io/terratest/modules/retry

go 1.24.0

require (
	github.com/gruntwork-io/terratest/modules/logger v0.1.0
	github.com/gruntwork-io/terratest/modules/testing v0.1.0
	github.com/stretchr/testify v1.10.0
	golang.org/x/net v0.47.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/gruntwork-io/terratest/modules/logger => ../logger
	github.com/gruntwork-io/terratest/modules/testing => ../testing
)
