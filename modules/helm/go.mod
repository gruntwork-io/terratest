module github.com/gruntwork-io/terratest/modules/helm

replace github.com/gruntwork-io/terratest/modules/shell => ../shell

replace github.com/gruntwork-io/terratest/modules/files => ../files

replace github.com/gruntwork-io/terratest/modules/http-helper => ../http-helper

replace github.com/gruntwork-io/terratest/modules/random => ../random

replace github.com/gruntwork-io/terratest/modules/k8s => ../k8s

replace github.com/gruntwork-io/terratest/modules/aws => ../aws

replace github.com/gruntwork-io/terratest/modules/environment => ../environment

replace github.com/gruntwork-io/terratest/modules/logger => ../logger

replace github.com/gruntwork-io/terratest/modules/retry => ../retry

replace github.com/gruntwork-io/terratest/modules/collections => ../collections

replace github.com/gruntwork-io/terratest/modules/customerrors => ../customerrors

replace github.com/gruntwork-io/terratest/modules/ssh => ../ssh

go 1.13

require (
	github.com/ghodss/yaml v1.0.0
	github.com/gruntwork-io/gruntwork-cli v0.6.1
	github.com/gruntwork-io/terratest/modules/files v0.0.0-00010101000000-000000000000
	github.com/gruntwork-io/terratest/modules/k8s v0.0.0-00010101000000-000000000000
	github.com/gruntwork-io/terratest/modules/shell v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.4.0
)
