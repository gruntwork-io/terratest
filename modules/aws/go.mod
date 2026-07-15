module github.com/gruntwork-io/terratest/modules/aws/v2

go 1.26.0

require (
	github.com/aws/aws-sdk-go-v2 v1.42.0
	github.com/aws/aws-sdk-go-v2/config v1.32.16
	github.com/aws/aws-sdk-go-v2/credentials v1.19.15
	github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager v0.1.17
	github.com/aws/aws-sdk-go-v2/service/acm v1.38.2
	github.com/aws/aws-sdk-go-v2/service/autoscaling v1.66.1
	github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs v1.69.1
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.57.2
	github.com/aws/aws-sdk-go-v2/service/ec2 v1.297.1
	github.com/aws/aws-sdk-go-v2/service/ecr v1.57.1
	github.com/aws/aws-sdk-go-v2/service/ecrpublic v1.39.6
	github.com/aws/aws-sdk-go-v2/service/ecs v1.78.1
	github.com/aws/aws-sdk-go-v2/service/iam v1.53.8
	github.com/aws/aws-sdk-go-v2/service/kms v1.50.5
	github.com/aws/aws-sdk-go-v2/service/lambda v1.89.1
	github.com/aws/aws-sdk-go-v2/service/rds v1.118.1
	github.com/aws/aws-sdk-go-v2/service/route53 v1.62.6
	github.com/aws/aws-sdk-go-v2/service/s3 v1.99.1
	github.com/aws/aws-sdk-go-v2/service/secretsmanager v1.41.6
	github.com/aws/aws-sdk-go-v2/service/sns v1.39.16
	github.com/aws/aws-sdk-go-v2/service/sqs v1.42.26
	github.com/aws/aws-sdk-go-v2/service/ssm v1.68.5
	github.com/aws/aws-sdk-go-v2/service/sts v1.42.0
	github.com/aws/smithy-go v1.27.1
	github.com/go-sql-driver/mysql v1.8.1
	github.com/google/uuid v1.6.0
	github.com/gruntwork-io/go-commons v0.8.0
	github.com/gruntwork-io/terratest/modules/core/v2 v2.0.0-00010101000000-000000000000
	github.com/gruntwork-io/terratest/modules/ssh/v2 v2.0.0-00010101000000-000000000000
	github.com/hashicorp/go-multierror v1.1.1
	github.com/jackc/pgx/v5 v5.9.0
	github.com/pquerna/otp v1.4.0
	github.com/stretchr/testify v1.11.1
)

require (
	filippo.io/edwards25519 v1.1.1 // indirect
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.7.9 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.18.22 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.4.29 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.7.29 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.4.23 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.8 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.9.14 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.11.22 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.13.22 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.19.22 // indirect
	github.com/aws/aws-sdk-go-v2/service/signin v1.0.10 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.30.16 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.35.20 // indirect
	github.com/boombuler/barcode v1.0.1-0.20190219062509-6c824513bacc // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.5 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/go-errors/errors v1.0.2-0.20180813162953-d98b870cc4e0 // indirect
	github.com/hashicorp/errwrap v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/mattn/go-zglob v0.0.2-0.20190814121620-e3c945676326 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/urfave/cli v1.22.16 // indirect
	golang.org/x/crypto v0.52.0 // indirect
	golang.org/x/sync v0.20.0 // indirect
	golang.org/x/sys v0.45.0 // indirect
	golang.org/x/text v0.37.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
