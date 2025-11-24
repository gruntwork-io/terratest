module github.com/gruntwork-io/terratest/modules/slack

go 1.24.0

require (
	github.com/gruntwork-io/terratest/modules/environment v0.1.0
	github.com/gruntwork-io/terratest/modules/random v0.1.0
	github.com/gruntwork-io/terratest/modules/retry v0.1.0
	github.com/gruntwork-io/terratest/modules/testing v0.1.0
)

replace (
	github.com/gruntwork-io/terratest/modules/environment => ../environment
	github.com/gruntwork-io/terratest/modules/random => ../random
	github.com/gruntwork-io/terratest/modules/retry => ../retry
	github.com/gruntwork-io/terratest/modules/testing => ../testing
)
