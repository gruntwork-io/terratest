module github.com/gruntwork-io/terratest/modules/aws

go 1.24.0

require (
	github.com/gruntwork-io/terratest/modules/collections v0.55.0
	github.com/gruntwork-io/terratest/modules/files v0.55.0
	github.com/gruntwork-io/terratest/modules/logger v0.55.0
	github.com/gruntwork-io/terratest/modules/random v0.55.0
	github.com/gruntwork-io/terratest/modules/retry v0.55.0
	github.com/gruntwork-io/terratest/modules/ssh v0.55.0
	github.com/gruntwork-io/terratest/modules/testing v0.55.0
)

replace (
	github.com/gruntwork-io/terratest/modules/collections => ../collections
	github.com/gruntwork-io/terratest/modules/files => ../files
	github.com/gruntwork-io/terratest/modules/logger => ../logger
	github.com/gruntwork-io/terratest/modules/random => ../random
	github.com/gruntwork-io/terratest/modules/retry => ../retry
	github.com/gruntwork-io/terratest/modules/ssh => ../ssh
	github.com/gruntwork-io/terratest/modules/testing => ../testing
)
