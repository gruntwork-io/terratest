module github.com/gruntwork-io/terratest/modules/oci

go 1.24.0

require (
	github.com/gruntwork-io/terratest/modules/logger v1.0.0
	github.com/gruntwork-io/terratest/modules/random v1.0.0
	github.com/gruntwork-io/terratest/modules/testing v1.0.0
	github.com/oracle/oci-go-sdk v24.3.0+incompatible
)

replace (
	github.com/gruntwork-io/terratest/modules/logger => ../logger
	github.com/gruntwork-io/terratest/modules/random => ../random
	github.com/gruntwork-io/terratest/modules/testing => ../testing
)
