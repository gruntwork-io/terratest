module github.com/gruntwork-io/terratest/modules/gcp

go 1.24.0

require (
	github.com/gruntwork-io/terratest/modules/collections v0.1.0
	github.com/gruntwork-io/terratest/modules/environment v0.1.0
	github.com/gruntwork-io/terratest/modules/logger v0.1.0
	github.com/gruntwork-io/terratest/modules/random v0.1.0
	github.com/gruntwork-io/terratest/modules/retry v0.1.0
	github.com/gruntwork-io/terratest/modules/ssh v0.1.0
	github.com/gruntwork-io/terratest/modules/testing v0.1.0
)

replace (
	github.com/gruntwork-io/terratest/modules/collections => ../collections
	github.com/gruntwork-io/terratest/modules/environment => ../environment
	github.com/gruntwork-io/terratest/modules/logger => ../logger
	github.com/gruntwork-io/terratest/modules/random => ../random
	github.com/gruntwork-io/terratest/modules/retry => ../retry
	github.com/gruntwork-io/terratest/modules/ssh => ../ssh
	github.com/gruntwork-io/terratest/modules/testing => ../testing
)
