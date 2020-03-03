module github.com/gruntwork-io/terratest/modules/oci

replace github.com/gruntwork-io/terratest/modules/logger => ../logger

replace github.com/gruntwork-io/terratest/modules/random => ../random

go 1.13

require (
	github.com/gruntwork-io/terratest/modules/logger v0.0.0-00010101000000-000000000000
	github.com/gruntwork-io/terratest/modules/random v0.0.0-00010101000000-000000000000
	github.com/oracle/oci-go-sdk v7.1.0+incompatible
)
