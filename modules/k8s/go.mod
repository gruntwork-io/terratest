module github.com/gruntwork-io/terratest/modules/k8s

go 1.24.0

require (
	github.com/gruntwork-io/terratest/modules/aws v0.55.0
	github.com/gruntwork-io/terratest/modules/environment v0.55.0
	github.com/gruntwork-io/terratest/modules/files v0.55.0
	github.com/gruntwork-io/terratest/modules/http-helper v0.55.0
	github.com/gruntwork-io/terratest/modules/logger v0.55.0
	github.com/gruntwork-io/terratest/modules/random v0.55.0
	github.com/gruntwork-io/terratest/modules/retry v0.55.0
	github.com/gruntwork-io/terratest/modules/shell v0.55.0
	github.com/gruntwork-io/terratest/modules/testing v0.55.0
)

replace (
	github.com/gruntwork-io/terratest/modules/aws => ../aws
	github.com/gruntwork-io/terratest/modules/environment => ../environment
	github.com/gruntwork-io/terratest/modules/files => ../files
	github.com/gruntwork-io/terratest/modules/http-helper => ../http-helper
	github.com/gruntwork-io/terratest/modules/logger => ../logger
	github.com/gruntwork-io/terratest/modules/random => ../random
	github.com/gruntwork-io/terratest/modules/retry => ../retry
	github.com/gruntwork-io/terratest/modules/shell => ../shell
	github.com/gruntwork-io/terratest/modules/testing => ../testing
)
