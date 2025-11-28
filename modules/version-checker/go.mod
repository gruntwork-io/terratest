module github.com/gruntwork-io/terratest/modules/version-checker

go 1.24.0

require (
	github.com/gruntwork-io/terratest/modules/shell v0.55.0
	github.com/gruntwork-io/terratest/modules/terraform v0.55.0
	github.com/gruntwork-io/terratest/modules/testing v0.55.0
)

replace (
	github.com/gruntwork-io/terratest/modules/shell => ../shell
	github.com/gruntwork-io/terratest/modules/terraform => ../terraform
	github.com/gruntwork-io/terratest/modules/testing => ../testing
)
