module github.com/gruntwork-io/terratest/modules/test-structure/v2

go 1.26.0

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

replace github.com/gruntwork-io/terratest/modules/aws/v2 => ../aws

replace github.com/gruntwork-io/terratest/modules/azure/v2 => ../azure

replace github.com/gruntwork-io/terratest/modules/gcp/v2 => ../gcp

replace github.com/gruntwork-io/terratest/modules/k8s/v2 => ../k8s

replace github.com/gruntwork-io/terratest/modules/helm/v2 => ../helm

replace github.com/gruntwork-io/terratest/modules/terraform/v2 => ../terraform

replace github.com/gruntwork-io/terratest/modules/terragrunt/v2 => ../terragrunt

require (
	github.com/gruntwork-io/go-commons v0.17.2
	github.com/gruntwork-io/terratest/modules/aws/v2 v2.0.0-00010101000000-000000000000
	github.com/gruntwork-io/terratest/modules/core/v2 v2.0.0-00010101000000-000000000000
	github.com/gruntwork-io/terratest/modules/k8s/v2 v2.0.0-00010101000000-000000000000
	github.com/gruntwork-io/terratest/modules/opa/v2 v2.0.0-00010101000000-000000000000
	github.com/gruntwork-io/terratest/modules/packer/v2 v2.0.0-00010101000000-000000000000
	github.com/gruntwork-io/terratest/modules/ssh/v2 v2.0.0-00010101000000-000000000000
	github.com/gruntwork-io/terratest/modules/terraform/v2 v2.0.0-00010101000000-000000000000
	github.com/mattn/go-zglob v0.0.6
	github.com/stretchr/testify v1.11.1
)

require (
	filippo.io/edwards25519 v1.2.0 // indirect
	github.com/agext/levenshtein v1.2.3 // indirect
	github.com/apparentlymart/go-textseg/v15 v15.0.0 // indirect
	github.com/aws/aws-sdk-go-v2 v1.41.7 // indirect
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.7.10 // indirect
	github.com/aws/aws-sdk-go-v2/config v1.32.18 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.19.17 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.18.23 // indirect
	github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager v0.1.22 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.4.23 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.7.23 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.4.24 // indirect
	github.com/aws/aws-sdk-go-v2/service/acm v1.39.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/autoscaling v1.66.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs v1.74.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.57.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/ec2 v1.304.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/ecr v1.57.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/ecs v1.80.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/iam v1.53.10 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.9 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.9.15 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.12.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.13.23 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.19.23 // indirect
	github.com/aws/aws-sdk-go-v2/service/kms v1.52.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/lambda v1.90.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/rds v1.118.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/route53 v1.62.7 // indirect
	github.com/aws/aws-sdk-go-v2/service/s3 v1.101.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/secretsmanager v1.41.7 // indirect
	github.com/aws/aws-sdk-go-v2/service/signin v1.0.11 // indirect
	github.com/aws/aws-sdk-go-v2/service/sns v1.39.17 // indirect
	github.com/aws/aws-sdk-go-v2/service/sqs v1.42.27 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssm v1.68.6 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.30.17 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.36.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.42.1 // indirect
	github.com/aws/smithy-go v1.25.1 // indirect
	github.com/bgentry/go-netrc v0.0.0-20140422174119-9fd32a8b3d3d // indirect
	github.com/boombuler/barcode v1.0.1-0.20190219062509-6c824513bacc // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.5 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/emicklei/go-restful/v3 v3.13.0 // indirect
	github.com/fxamacker/cbor/v2 v2.9.0 // indirect
	github.com/go-errors/errors v1.4.2 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-openapi/jsonpointer v0.21.0 // indirect
	github.com/go-openapi/jsonreference v0.20.2 // indirect
	github.com/go-openapi/swag v0.23.0 // indirect
	github.com/go-sql-driver/mysql v1.10.0 // indirect
	github.com/google/gnostic-models v0.7.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/websocket v1.5.4-0.20250319132907-e064f32e3674 // indirect
	github.com/gruntwork-io/terratest v1.0.0 // indirect
	github.com/gruntwork-io/terratest/modules/shell/v2 v2.0.0-00010101000000-000000000000 // indirect
	github.com/hashicorp/errwrap v1.0.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-getter/v2 v2.2.3 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-safetemp v1.0.0 // indirect
	github.com/hashicorp/go-version v1.9.0 // indirect
	github.com/hashicorp/hcl/v2 v2.24.0 // indirect
	github.com/hashicorp/terraform-json v0.27.2 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgx/v5 v5.9.2 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/jinzhu/copier v0.4.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/compress v1.16.5 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/go-testing-interface v1.14.1 // indirect
	github.com/mitchellh/go-wordwrap v1.0.1 // indirect
	github.com/moby/spdystream v0.5.1 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.3-0.20250322232337-35a7c28c31ee // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/pquerna/otp v1.5.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/spf13/pflag v1.0.9 // indirect
	github.com/tmccombs/hcl2json v0.6.9 // indirect
	github.com/ulikunitz/xz v0.5.10 // indirect
	github.com/urfave/cli/v2 v2.10.3 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	github.com/xrash/smetrics v0.0.0-20201216005158-039620a65673 // indirect
	github.com/zclconf/go-cty v1.18.1 // indirect
	go.yaml.in/yaml/v2 v2.4.3 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	golang.org/x/crypto v0.52.0 // indirect
	golang.org/x/exp v0.0.0-20221106115401-f9659909a136 // indirect
	golang.org/x/mod v0.35.0 // indirect
	golang.org/x/net v0.54.0 // indirect
	golang.org/x/oauth2 v0.36.0 // indirect
	golang.org/x/sync v0.20.0 // indirect
	golang.org/x/sys v0.45.0 // indirect
	golang.org/x/term v0.43.0 // indirect
	golang.org/x/text v0.37.0 // indirect
	golang.org/x/time v0.15.0 // indirect
	golang.org/x/tools v0.44.0 // indirect
	google.golang.org/protobuf v1.36.12-0.20260120151049-f2248ac996af // indirect
	gopkg.in/evanphx/json-patch.v4 v4.13.0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/api v0.36.1 // indirect
	k8s.io/apimachinery v0.36.1 // indirect
	k8s.io/client-go v0.36.1 // indirect
	k8s.io/klog/v2 v2.140.0 // indirect
	k8s.io/kube-openapi v0.0.0-20260317180543-43fb72c5454a // indirect
	k8s.io/streaming v0.36.1 // indirect
	k8s.io/utils v0.0.0-20260210185600-b8788abfbbc2 // indirect
	sigs.k8s.io/json v0.0.0-20250730193827-2d320260d730 // indirect
	sigs.k8s.io/randfill v1.0.0 // indirect
	sigs.k8s.io/structured-merge-diff/v6 v6.3.2 // indirect
	sigs.k8s.io/yaml v1.6.0 // indirect
)
