module github.com/gruntwork-io/terratest/modules/oci

go 1.24.0

require (
	github.com/gruntwork-io/terratest/modules/logger v0.55.0
	github.com/gruntwork-io/terratest/modules/random v0.55.0
	github.com/gruntwork-io/terratest/modules/testing v0.55.0
	github.com/oracle/oci-go-sdk v24.3.0+incompatible
)

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
)

replace (
	github.com/gruntwork-io/terratest/modules/logger => ../logger
	github.com/gruntwork-io/terratest/modules/random => ../random
	github.com/gruntwork-io/terratest/modules/testing => ../testing
)
