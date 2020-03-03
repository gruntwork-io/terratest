module github.com/gruntwork-io/terratest/modules/k8s

replace github.com/gruntwork-io/terratest/modules/logger => ../logger

replace github.com/gruntwork-io/terratest/modules/environment => ../environment

replace github.com/gruntwork-io/terratest/modules/files => ../files

replace github.com/gruntwork-io/terratest/modules/random => ../random

replace github.com/gruntwork-io/terratest/modules/shell => ../shell

replace github.com/gruntwork-io/terratest/modules/retry => ../retry

replace github.com/gruntwork-io/terratest/modules/http-helper => ../http-helper

replace github.com/gruntwork-io/terratest/modules/aws => ../aws

replace github.com/gruntwork-io/terratest/modules/collections => ../collections

replace github.com/gruntwork-io/terratest/modules/customerrors => ../customerrors

replace github.com/gruntwork-io/terratest/modules/ssh => ../ssh

go 1.13

require (
	github.com/gruntwork-io/gruntwork-cli v0.6.1
	github.com/gruntwork-io/terratest/modules/aws v0.0.0-00010101000000-000000000000
	github.com/gruntwork-io/terratest/modules/environment v0.0.0-00010101000000-000000000000
	github.com/gruntwork-io/terratest/modules/files v0.0.0-00010101000000-000000000000
	github.com/gruntwork-io/terratest/modules/http-helper v0.0.0-00010101000000-000000000000
	github.com/gruntwork-io/terratest/modules/logger v0.0.0-00010101000000-000000000000
	github.com/gruntwork-io/terratest/modules/random v0.0.0-00010101000000-000000000000
	github.com/gruntwork-io/terratest/modules/retry v0.0.0-00010101000000-000000000000
	github.com/gruntwork-io/terratest/modules/shell v0.0.0-00010101000000-000000000000
	github.com/imdario/mergo v0.3.7 // indirect
	github.com/mitchellh/go-homedir v1.1.0
	github.com/stretchr/testify v1.4.0
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d // indirect
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0 // indirect
	k8s.io/api v0.17.0
	k8s.io/apimachinery v0.17.0
	k8s.io/client-go v0.17.0
)
