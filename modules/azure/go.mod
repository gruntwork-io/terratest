module github.com/gruntwork-io/terratest/modules/azure

go 1.24.0

require (
	github.com/gruntwork-io/terratest/modules/collections v0.1.0
	github.com/gruntwork-io/terratest/modules/random v0.1.0
	github.com/gruntwork-io/terratest/modules/testing v0.1.0
)

replace (
	github.com/gruntwork-io/terratest/modules/collections => ../collections
	github.com/gruntwork-io/terratest/modules/random => ../random
	github.com/gruntwork-io/terratest/modules/testing => ../testing
)
