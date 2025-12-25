module github.com/gruntwork-io/terratest/cmd/pick-instance-type/v2

go 1.24.0

replace (
	github.com/gruntwork-io/terratest/modules/aws/v2 => ../../modules/aws
	github.com/gruntwork-io/terratest/modules/collections/v2 => ../../modules/collections
	github.com/gruntwork-io/terratest/modules/files/v2 => ../../modules/files
	github.com/gruntwork-io/terratest/modules/logger/v2 => ../../modules/logger
	github.com/gruntwork-io/terratest/modules/random/v2 => ../../modules/random
	github.com/gruntwork-io/terratest/modules/retry/v2 => ../../modules/retry
	github.com/gruntwork-io/terratest/modules/shell/v2 => ../../modules/shell
	github.com/gruntwork-io/terratest/modules/ssh/v2 => ../../modules/ssh
	github.com/gruntwork-io/terratest/modules/testing/v2 => ../../modules/testing
)

exclude github.com/gruntwork-io/terratest v0.46.16

require (
	github.com/gruntwork-io/terratest/modules/aws/v2 v2.0.0
	github.com/urfave/cli/v2 v2.10.3
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/aws/aws-sdk-go-v2 v1.40.0 // indirect
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.7.3 // indirect
	github.com/aws/aws-sdk-go-v2/config v1.32.2 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.19.2 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.18.14 // indirect
	github.com/aws/aws-sdk-go-v2/feature/s3/manager v1.20.12 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.4.14 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.7.14 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.4 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.4.14 // indirect
	github.com/aws/aws-sdk-go-v2/service/acm v1.37.15 // indirect
	github.com/aws/aws-sdk-go-v2/service/autoscaling v1.62.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs v1.61.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.53.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/ec2 v1.275.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/ecr v1.54.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/ecs v1.69.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/iam v1.52.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.9.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.11.14 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.13.14 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.19.14 // indirect
	github.com/aws/aws-sdk-go-v2/service/kms v1.49.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/lambda v1.84.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/rds v1.111.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/route53 v1.61.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/s3 v1.92.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/secretsmanager v1.40.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/signin v1.0.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/sns v1.39.7 // indirect
	github.com/aws/aws-sdk-go-v2/service/sqs v1.42.17 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssm v1.67.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.30.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.35.10 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.41.2 // indirect
	github.com/aws/smithy-go v1.23.2 // indirect
	github.com/boombuler/barcode v1.0.1-0.20190219062509-6c824513bacc // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.7 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/go-errors/errors v1.4.2 // indirect
	github.com/go-sql-driver/mysql v1.9.3 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gruntwork-io/go-commons v0.17.2 // indirect
	github.com/gruntwork-io/terratest/modules/collections/v2 v2.0.0 // indirect
	github.com/gruntwork-io/terratest/modules/files/v2 v2.0.0 // indirect
	github.com/gruntwork-io/terratest/modules/logger/v2 v2.0.0 // indirect
	github.com/gruntwork-io/terratest/modules/random/v2 v2.0.0 // indirect
	github.com/gruntwork-io/terratest/modules/retry/v2 v2.0.0 // indirect
	github.com/gruntwork-io/terratest/modules/ssh/v2 v2.0.0 // indirect
	github.com/gruntwork-io/terratest/modules/testing/v2 v2.0.0 // indirect
	github.com/hashicorp/errwrap v1.0.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgx/v5 v5.7.6 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/mattn/go-zglob v0.0.3 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/pquerna/otp v1.5.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/stretchr/testify v1.11.1 // indirect
	github.com/xrash/smetrics v0.0.0-20201216005158-039620a65673 // indirect
	golang.org/x/crypto v0.45.0 // indirect
	golang.org/x/net v0.47.0 // indirect
	golang.org/x/sync v0.18.0 // indirect
	golang.org/x/sys v0.38.0 // indirect
	golang.org/x/text v0.31.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
