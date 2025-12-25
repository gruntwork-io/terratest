module github.com/gruntwork-io/terratest/modules/http-helper/v2

go 1.24.0

require (
	github.com/gruntwork-io/terratest/modules/logger/v2 v2.0.0
	github.com/gruntwork-io/terratest/modules/random/v2 v2.0.0
	github.com/gruntwork-io/terratest/modules/retry/v2 v2.0.0
	github.com/gruntwork-io/terratest/modules/testing/v2 v2.0.0
	github.com/stretchr/testify v1.11.1
)

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	golang.org/x/net v0.47.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/gruntwork-io/terratest/modules/logger/v2 => ../logger
	github.com/gruntwork-io/terratest/modules/random/v2 => ../random
	github.com/gruntwork-io/terratest/modules/retry/v2 => ../retry
	github.com/gruntwork-io/terratest/modules/testing/v2 => ../testing
)
