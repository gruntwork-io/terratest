module github.com/gruntwork-io/terratest/modules/aws/v2

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

replace github.com/gruntwork-io/terratest/modules/opa/v2 => ../opa

replace github.com/gruntwork-io/terratest/modules/azure/v2 => ../azure

replace github.com/gruntwork-io/terratest/modules/gcp/v2 => ../gcp

replace github.com/gruntwork-io/terratest/modules/k8s/v2 => ../k8s

replace github.com/gruntwork-io/terratest/modules/helm/v2 => ../helm

replace github.com/gruntwork-io/terratest/modules/terraform/v2 => ../terraform

replace github.com/gruntwork-io/terratest/modules/terragrunt/v2 => ../terragrunt

replace github.com/gruntwork-io/terratest/modules/test-structure/v2 => ../test-structure

require (
	github.com/aws/aws-sdk-go-v2 v1.41.7
	github.com/aws/aws-sdk-go-v2/config v1.32.18
	github.com/aws/aws-sdk-go-v2/credentials v1.19.17
	github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager v0.1.22
	github.com/aws/aws-sdk-go-v2/service/acm v1.39.0
	github.com/aws/aws-sdk-go-v2/service/autoscaling v1.66.2
	github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs v1.74.0
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.57.4
	github.com/aws/aws-sdk-go-v2/service/ec2 v1.304.0
	github.com/aws/aws-sdk-go-v2/service/ecr v1.57.2
	github.com/aws/aws-sdk-go-v2/service/ecs v1.80.0
	github.com/aws/aws-sdk-go-v2/service/iam v1.53.10
	github.com/aws/aws-sdk-go-v2/service/kms v1.52.0
	github.com/aws/aws-sdk-go-v2/service/lambda v1.90.1
	github.com/aws/aws-sdk-go-v2/service/rds v1.118.2
	github.com/aws/aws-sdk-go-v2/service/route53 v1.62.7
	github.com/aws/aws-sdk-go-v2/service/s3 v1.101.0
	github.com/aws/aws-sdk-go-v2/service/secretsmanager v1.41.7
	github.com/aws/aws-sdk-go-v2/service/sns v1.39.17
	github.com/aws/aws-sdk-go-v2/service/sqs v1.42.27
	github.com/aws/aws-sdk-go-v2/service/ssm v1.68.6
	github.com/aws/aws-sdk-go-v2/service/sts v1.42.1
	github.com/aws/smithy-go v1.25.1
	github.com/go-sql-driver/mysql v1.10.0
	github.com/google/uuid v1.6.0
	github.com/gruntwork-io/go-commons v0.17.2
	github.com/gruntwork-io/terratest/modules/core/v2 v2.0.0-00010101000000-000000000000
	github.com/gruntwork-io/terratest/modules/ssh/v2 v2.0.0-00010101000000-000000000000
	github.com/hashicorp/go-multierror v1.1.1
	github.com/jackc/pgx/v5 v5.9.2
	github.com/pquerna/otp v1.5.0
	github.com/stretchr/testify v1.11.1
)

require (
	filippo.io/edwards25519 v1.2.0 // indirect
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.7.10 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.18.23 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.4.23 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.7.23 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.4.24 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.9 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.9.15 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.12.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.13.23 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.19.23 // indirect
	github.com/aws/aws-sdk-go-v2/service/signin v1.0.11 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.30.17 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.36.0 // indirect
	github.com/boombuler/barcode v1.0.1-0.20190219062509-6c824513bacc // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-errors/errors v1.4.2 // indirect
	github.com/hashicorp/errwrap v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/mattn/go-zglob v0.0.6 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/urfave/cli/v2 v2.10.3 // indirect
	github.com/xrash/smetrics v0.0.0-20201216005158-039620a65673 // indirect
	golang.org/x/crypto v0.52.0 // indirect
	golang.org/x/sync v0.20.0 // indirect
	golang.org/x/sys v0.45.0 // indirect
	golang.org/x/text v0.37.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
