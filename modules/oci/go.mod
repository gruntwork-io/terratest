module github.com/gruntwork-io/terratest/modules/oci/v2

go 1.24.0

require (
	github.com/gruntwork-io/terratest/modules/logger/v2 v2.0.0
	github.com/gruntwork-io/terratest/modules/random/v2 v2.0.0
	github.com/gruntwork-io/terratest/modules/testing/v2 v2.0.0
	github.com/oracle/oci-go-sdk v24.3.0+incompatible
)

replace (
	github.com/gruntwork-io/terratest/modules/logger/v2 => ../logger
	github.com/gruntwork-io/terratest/modules/random/v2 => ../random
	github.com/gruntwork-io/terratest/modules/testing/v2 => ../testing
)
