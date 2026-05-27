module github.com/gruntwork-io/terratest/cmd/terratest_log_parser/v2

go 1.26

replace github.com/gruntwork-io/terratest/modules/core/v2 => ../../modules/core

replace github.com/gruntwork-io/terratest/modules/shell/v2 => ../../modules/shell

replace github.com/gruntwork-io/terratest/modules/ssh/v2 => ../../modules/ssh

replace github.com/gruntwork-io/terratest/modules/http-helper/v2 => ../../modules/http-helper

replace github.com/gruntwork-io/terratest/modules/dns-helper/v2 => ../../modules/dns-helper

replace github.com/gruntwork-io/terratest/modules/version-checker/v2 => ../../modules/version-checker

replace github.com/gruntwork-io/terratest/modules/docker/v2 => ../../modules/docker

replace github.com/gruntwork-io/terratest/modules/packer/v2 => ../../modules/packer

replace github.com/gruntwork-io/terratest/modules/database/v2 => ../../modules/database

replace github.com/gruntwork-io/terratest/modules/slack/v2 => ../../modules/slack

replace github.com/gruntwork-io/terratest/modules/oci/v2 => ../../modules/oci

replace github.com/gruntwork-io/terratest/modules/opa/v2 => ../../modules/opa

replace github.com/gruntwork-io/terratest/modules/aws/v2 => ../../modules/aws

replace github.com/gruntwork-io/terratest/modules/azure/v2 => ../../modules/azure

replace github.com/gruntwork-io/terratest/modules/gcp/v2 => ../../modules/gcp

replace github.com/gruntwork-io/terratest/modules/k8s/v2 => ../../modules/k8s

replace github.com/gruntwork-io/terratest/modules/helm/v2 => ../../modules/helm

replace github.com/gruntwork-io/terratest/modules/terraform/v2 => ../../modules/terraform

replace github.com/gruntwork-io/terratest/modules/terragrunt/v2 => ../../modules/terragrunt

replace github.com/gruntwork-io/terratest/modules/test-structure/v2 => ../../modules/test-structure

require (
	github.com/gruntwork-io/go-commons v0.17.2
	github.com/gruntwork-io/terratest/modules/core/v2 v2.0.0-00010101000000-000000000000
	github.com/sirupsen/logrus v1.9.4
	github.com/urfave/cli v1.22.17
)

require (
	github.com/cpuguy83/go-md2man/v2 v2.0.7 // indirect
	github.com/go-errors/errors v1.4.2 // indirect
	github.com/jstemmer/go-junit-report v1.0.0 // indirect
	github.com/mattn/go-zglob v0.0.6 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	golang.org/x/sys v0.18.0 // indirect
)

replace github.com/gruntwork-io/go-commons => github.com/gruntwork-io/go-commons v0.8.0
