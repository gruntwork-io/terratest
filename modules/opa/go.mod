module github.com/gruntwork-io/terratest/modules/opa/v2

go 1.26

replace github.com/gruntwork-io/terratest/modules/core/v2 => ../core

replace github.com/gruntwork-io/terratest/modules/shell/v2 => ../shell

replace github.com/gruntwork-io/terratest/modules/ssh/v2 => ../ssh

replace github.com/gruntwork-io/terratest/modules/http-helper/v2 => ../http-helper

replace github.com/gruntwork-io/terratest/modules/dns-helper/v2 => ../dns-helper

replace github.com/gruntwork-io/terratest/modules/version-checker/v2 => ../version-checker

replace github.com/gruntwork-io/terratest/modules/docker/v2 => ../docker

replace github.com/gruntwork-io/terratest/modules/packer/v2 => ../packer

replace github.com/gruntwork-io/terratest/modules/database/v2 => ../database

replace github.com/gruntwork-io/terratest/modules/slack/v2 => ../slack

replace github.com/gruntwork-io/terratest/modules/oci/v2 => ../oci

replace github.com/gruntwork-io/terratest/modules/aws/v2 => ../aws

replace github.com/gruntwork-io/terratest/modules/azure/v2 => ../azure

replace github.com/gruntwork-io/terratest/modules/gcp/v2 => ../gcp

replace github.com/gruntwork-io/terratest/modules/k8s/v2 => ../k8s

replace github.com/gruntwork-io/terratest/modules/helm/v2 => ../helm

replace github.com/gruntwork-io/terratest/modules/terraform/v2 => ../terraform

replace github.com/gruntwork-io/terratest/modules/terragrunt/v2 => ../terragrunt

replace github.com/gruntwork-io/terratest/modules/test-structure/v2 => ../test-structure

require (
	github.com/gruntwork-io/terratest/modules/core/v2 v2.0.0-00010101000000-000000000000
	github.com/gruntwork-io/terratest/modules/shell/v2 v2.0.0-00010101000000-000000000000
	github.com/hashicorp/go-getter/v2 v2.2.3
	github.com/hashicorp/go-multierror v1.1.1
	github.com/stretchr/testify v1.11.1
	golang.org/x/sync v0.20.0
)

require (
	github.com/bgentry/go-netrc v0.0.0-20140422174119-9fd32a8b3d3d // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/hashicorp/errwrap v1.0.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.0 // indirect
	github.com/hashicorp/go-safetemp v1.0.0 // indirect
	github.com/hashicorp/go-version v1.1.0 // indirect
	github.com/klauspost/compress v1.11.2 // indirect
	github.com/mattn/go-zglob v0.0.6 // indirect
	github.com/mitchellh/go-homedir v1.0.0 // indirect
	github.com/mitchellh/go-testing-interface v1.0.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/ulikunitz/xz v0.5.8 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
