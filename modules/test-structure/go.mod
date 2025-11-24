module github.com/gruntwork-io/terratest/modules/test-structure

go 1.24.0

require (
	github.com/gruntwork-io/terratest/modules/aws v0.1.0
	github.com/gruntwork-io/terratest/modules/collections v0.1.0
	github.com/gruntwork-io/terratest/modules/files v0.1.0
	github.com/gruntwork-io/terratest/modules/git v0.1.0
	github.com/gruntwork-io/terratest/modules/k8s v0.1.0
	github.com/gruntwork-io/terratest/modules/logger v0.1.0
	github.com/gruntwork-io/terratest/modules/opa v0.1.0
	github.com/gruntwork-io/terratest/modules/packer v0.1.0
	github.com/gruntwork-io/terratest/modules/ssh v0.1.0
	github.com/gruntwork-io/terratest/modules/terraform v0.1.0
	github.com/gruntwork-io/terratest/modules/testing v0.1.0
)

replace (
	github.com/gruntwork-io/terratest/modules/aws => ../aws
	github.com/gruntwork-io/terratest/modules/collections => ../collections
	github.com/gruntwork-io/terratest/modules/files => ../files
	github.com/gruntwork-io/terratest/modules/git => ../git
	github.com/gruntwork-io/terratest/modules/k8s => ../k8s
	github.com/gruntwork-io/terratest/modules/logger => ../logger
	github.com/gruntwork-io/terratest/modules/opa => ../opa
	github.com/gruntwork-io/terratest/modules/packer => ../packer
	github.com/gruntwork-io/terratest/modules/ssh => ../ssh
	github.com/gruntwork-io/terratest/modules/terraform => ../terraform
	github.com/gruntwork-io/terratest/modules/testing => ../testing
)
