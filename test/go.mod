module test

replace github.com/gruntwork-io/terratest/modules/docker => ../modules/docker

replace github.com/gruntwork-io/terratest/modules/helm => ../modules/helm

replace github.com/gruntwork-io/terratest/modules/http-helper => ../modules/http-helper

replace github.com/gruntwork-io/terratest/modules/k8s => ../modules/k8s

replace github.com/gruntwork-io/terratest/modules/random => ../modules/random

replace github.com/gruntwork-io/terratest/modules/logger => ../modules/logger

replace github.com/gruntwork-io/terratest/modules/aws => ../modules/aws

replace github.com/gruntwork-io/terratest/modules/packer => ../modules/packer

replace github.com/gruntwork-io/terratest/modules/gcp => ../modules/gcp

replace github.com/gruntwork-io/terratest/modules/oci => ../modules/oci

replace github.com/gruntwork-io/terratest/modules/terraform => ../modules/terraform

replace github.com/gruntwork-io/terratest/modules/azure => ../modules/azure

replace github.com/gruntwork-io/terratest/modules/test-structure => ../modules/test-structure

replace github.com/gruntwork-io/terratest/modules/retry => ../modules/retry

replace github.com/gruntwork-io/terratest/modules/ssh => ../modules/ssh

replace github.com/gruntwork-io/terratest/modules/collections => ../modules/collections

replace github.com/gruntwork-io/terratest/modules/customerrors => ../modules/customerrors

replace github.com/gruntwork-io/terratest/modules/files => ../modules/files

replace github.com/gruntwork-io/terratest/modules/shell => ../modules/shell

replace github.com/gruntwork-io/terratest/modules/environment => ../modules/environment

go 1.13

require (
	github.com/Azure/azure-sdk-for-go v38.1.0+incompatible
	github.com/aws/aws-sdk-go v1.28.14
	github.com/gruntwork-io/terratest/modules/aws v0.0.0-00010101000000-000000000000
	github.com/gruntwork-io/terratest/modules/azure v0.0.0-00010101000000-000000000000
	github.com/gruntwork-io/terratest/modules/docker v0.0.0-00010101000000-000000000000
	github.com/gruntwork-io/terratest/modules/gcp v0.0.0-00010101000000-000000000000
	github.com/gruntwork-io/terratest/modules/helm v0.0.0-00010101000000-000000000000
	github.com/gruntwork-io/terratest/modules/http-helper v0.0.0-00010101000000-000000000000
	github.com/gruntwork-io/terratest/modules/k8s v0.0.0-00010101000000-000000000000
	github.com/gruntwork-io/terratest/modules/logger v0.0.0-00010101000000-000000000000
	github.com/gruntwork-io/terratest/modules/oci v0.0.0-00010101000000-000000000000
	github.com/gruntwork-io/terratest/modules/packer v0.0.0-00010101000000-000000000000
	github.com/gruntwork-io/terratest/modules/random v0.0.0-00010101000000-000000000000
	github.com/gruntwork-io/terratest/modules/retry v0.0.0-00010101000000-000000000000
	github.com/gruntwork-io/terratest/modules/ssh v0.0.0-00010101000000-000000000000
	github.com/gruntwork-io/terratest/modules/terraform v0.0.0-00010101000000-000000000000
	github.com/gruntwork-io/terratest/modules/test-structure v0.0.0-00010101000000-000000000000
	github.com/magiconair/properties v1.8.1
	github.com/stretchr/testify v1.4.0
	k8s.io/api v0.17.0
)
