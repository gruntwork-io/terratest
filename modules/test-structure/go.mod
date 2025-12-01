module github.com/gruntwork-io/terratest/modules/test-structure

go 1.24.0

require (
	github.com/gruntwork-io/terratest/modules/aws v0.55.0
	github.com/gruntwork-io/terratest/modules/collections v0.55.0
	github.com/gruntwork-io/terratest/modules/files v0.55.0
	github.com/gruntwork-io/terratest/modules/git v0.55.0
	github.com/gruntwork-io/terratest/modules/k8s v0.55.0
	github.com/gruntwork-io/terratest/modules/logger v0.55.0
	github.com/gruntwork-io/terratest/modules/opa v0.55.0
	github.com/gruntwork-io/terratest/modules/packer v0.55.0
	github.com/gruntwork-io/terratest/modules/ssh v0.55.0
	github.com/gruntwork-io/terratest/modules/terraform v0.55.0
	github.com/gruntwork-io/terratest/modules/testing v0.55.0
)

require (
	github.com/gruntwork-io/go-commons v0.17.2 // indirect
	golang.org/x/exp v0.0.0-20221106115401-f9659909a136 // indirect
	k8s.io/api v0.28.4 // indirect
	k8s.io/apimachinery v0.28.4 // indirect
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
