module github.com/gruntwork-io/terratest/modules/test-structure

replace github.com/gruntwork-io/terratest/modules/files => ../files

replace github.com/gruntwork-io/terratest/modules/k8s => ../k8s

replace github.com/gruntwork-io/terratest/modules/terraform => ../terraform

replace github.com/gruntwork-io/terratest/modules/aws => ../aws

replace github.com/gruntwork-io/terratest/modules/logger => ../logger

replace github.com/gruntwork-io/terratest/modules/packer => ../packer

replace github.com/gruntwork-io/terratest/modules/collections => ../collections

replace github.com/gruntwork-io/terratest/modules/customerrors => ../customerrors

replace github.com/gruntwork-io/terratest/modules/random => ../random

replace github.com/gruntwork-io/terratest/modules/retry => ../retry

replace github.com/gruntwork-io/terratest/modules/ssh => ../ssh

replace github.com/gruntwork-io/terratest/modules/environment => ../environment

replace github.com/gruntwork-io/terratest/modules/http-helper => ../http-helper

replace github.com/gruntwork-io/terratest/modules/shell => ../shell

go 1.13

require (
	github.com/gruntwork-io/terratest/modules/aws v0.0.0-00010101000000-000000000000
	github.com/gruntwork-io/terratest/modules/files v0.0.0-00010101000000-000000000000
	github.com/gruntwork-io/terratest/modules/k8s v0.0.0-00010101000000-000000000000
	github.com/gruntwork-io/terratest/modules/logger v0.0.0-00010101000000-000000000000
	github.com/gruntwork-io/terratest/modules/packer v0.0.0-00010101000000-000000000000
	github.com/gruntwork-io/terratest/modules/terraform v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.4.0
)
