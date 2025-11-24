module github.com/gruntwork-io/terratest/modules/packer

go 1.24.0

require (
	github.com/gruntwork-io/terratest/modules/logger v0.1.0
	github.com/gruntwork-io/terratest/modules/retry v0.1.0
	github.com/gruntwork-io/terratest/modules/shell v0.1.0
	github.com/gruntwork-io/terratest/modules/testing v0.1.0
)

replace (
	github.com/gruntwork-io/terratest/modules/logger => ../logger
	github.com/gruntwork-io/terratest/modules/retry => ../retry
	github.com/gruntwork-io/terratest/modules/shell => ../shell
	github.com/gruntwork-io/terratest/modules/testing => ../testing
)
