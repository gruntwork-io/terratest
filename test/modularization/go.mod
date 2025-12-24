module github.com/gruntwork-io/terratest/test/modularization/v2

go 1.24.0

require (
	github.com/gruntwork-io/terratest/modules/aws/v2 v2.0.0
	github.com/gruntwork-io/terratest/modules/azure/v2 v2.0.0
	github.com/gruntwork-io/terratest/modules/collections/v2 v2.0.0
	github.com/gruntwork-io/terratest/modules/dns-helper/v2 v2.0.0
	github.com/gruntwork-io/terratest/modules/environment/v2 v2.0.0
	github.com/gruntwork-io/terratest/modules/gcp/v2 v2.0.0
	github.com/gruntwork-io/terratest/modules/git/v2 v2.0.0
	github.com/gruntwork-io/terratest/modules/helm/v2 v2.0.0
	github.com/gruntwork-io/terratest/modules/http-helper/v2 v2.0.0
	github.com/gruntwork-io/terratest/modules/logger/v2 v2.0.0
	github.com/gruntwork-io/terratest/modules/oci/v2 v2.0.0
	github.com/gruntwork-io/terratest/modules/ssh/v2 v2.0.0
	github.com/gruntwork-io/terratest/modules/terraform/v2 v2.0.0
	github.com/gruntwork-io/terratest/modules/terragrunt/v2 v2.0.0
	github.com/gruntwork-io/terratest/modules/testing/v2 v2.0.0
	github.com/stretchr/testify v1.11.1
)

require (
	cel.dev/expr v0.16.1 // indirect
	cloud.google.com/go v0.116.0 // indirect
	cloud.google.com/go/auth v0.10.2 // indirect
	cloud.google.com/go/auth/oauth2adapt v0.2.5 // indirect
	cloud.google.com/go/cloudbuild v1.19.0 // indirect
	cloud.google.com/go/compute/metadata v0.7.0 // indirect
	cloud.google.com/go/iam v1.2.2 // indirect
	cloud.google.com/go/longrunning v0.6.2 // indirect
	cloud.google.com/go/monitoring v1.21.2 // indirect
	cloud.google.com/go/storage v1.47.0 // indirect
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/Azure/azure-sdk-for-go v68.0.0+incompatible // indirect
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.20.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/azidentity v1.13.1 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/internal v1.11.2 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/appcontainers/armappcontainers/v3 v3.1.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/appservice/armappservice/v2 v2.3.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/datafactory/armdatafactory/v9 v9.1.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/keyvault/armkeyvault v1.5.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/mysql/armmysql v1.2.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/postgresql/armpostgresql v1.2.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources v1.2.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/sql/armsql v1.2.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/synapse/armsynapse v0.8.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azcertificates v1.4.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys v1.4.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azsecrets v1.4.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/internal v1.2.0 // indirect
	github.com/Azure/go-autorest v14.2.0+incompatible // indirect
	github.com/Azure/go-autorest/autorest v0.11.30 // indirect
	github.com/Azure/go-autorest/autorest/adal v0.9.22 // indirect
	github.com/Azure/go-autorest/autorest/azure/auth v0.5.13 // indirect
	github.com/Azure/go-autorest/autorest/azure/cli v0.4.6 // indirect
	github.com/Azure/go-autorest/autorest/date v0.3.0 // indirect
	github.com/Azure/go-autorest/autorest/to v0.4.1 // indirect
	github.com/Azure/go-autorest/autorest/validation v0.3.2 // indirect
	github.com/Azure/go-autorest/logger v0.2.1 // indirect
	github.com/Azure/go-autorest/tracing v0.6.0 // indirect
	github.com/AzureAD/microsoft-authentication-library-for-go v1.6.0 // indirect
	github.com/BurntSushi/toml v1.5.0 // indirect
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/detectors/gcp v1.24.1 // indirect
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/metric v0.48.1 // indirect
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/internal/resourcemapping v0.48.1 // indirect
	github.com/agext/levenshtein v1.2.3 // indirect
	github.com/apparentlymart/go-textseg/v15 v15.0.0 // indirect
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
	github.com/bgentry/go-netrc v0.0.0-20140422174119-9fd32a8b3d3d // indirect
	github.com/boombuler/barcode v1.0.1-0.20190219062509-6c824513bacc // indirect
	github.com/census-instrumentation/opencensus-proto v0.4.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/cncf/xds/go v0.0.0-20240905190251-b4127c9b8d78 // indirect
	github.com/containerd/stargz-snapshotter/estargz v0.18.1 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.7 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/dimchansky/utfbom v1.1.1 // indirect
	github.com/docker/cli v29.0.3+incompatible // indirect
	github.com/docker/distribution v2.8.3+incompatible // indirect
	github.com/docker/docker-credential-helpers v0.9.3 // indirect
	github.com/envoyproxy/go-control-plane v0.13.0 // indirect
	github.com/envoyproxy/protoc-gen-validate v1.1.0 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/go-errors/errors v1.4.2 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-sql-driver/mysql v1.9.3 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.0 // indirect
	github.com/golang-jwt/jwt/v5 v5.3.0 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/gonvenience/bunt v1.4.2 // indirect
	github.com/gonvenience/idem v0.0.2 // indirect
	github.com/gonvenience/neat v1.3.16 // indirect
	github.com/gonvenience/term v1.0.4 // indirect
	github.com/gonvenience/text v1.0.9 // indirect
	github.com/gonvenience/ytbx v1.4.7 // indirect
	github.com/google/go-containerregistry v0.20.7 // indirect
	github.com/google/s2a-go v0.1.8 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.3.4 // indirect
	github.com/googleapis/gax-go/v2 v2.14.0 // indirect
	github.com/gruntwork-io/go-commons v0.17.2 // indirect
	github.com/gruntwork-io/terratest/modules/files/v2 v2.0.0 // indirect
	github.com/gruntwork-io/terratest/modules/opa/v2 v2.0.0 // indirect
	github.com/gruntwork-io/terratest/modules/random/v2 v2.0.0 // indirect
	github.com/gruntwork-io/terratest/modules/retry/v2 v2.0.0 // indirect
	github.com/gruntwork-io/terratest/modules/shell/v2 v2.0.0 // indirect
	github.com/hashicorp/errwrap v1.0.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-getter/v2 v2.2.3 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-safetemp v1.0.0 // indirect
	github.com/hashicorp/go-version v1.7.0 // indirect
	github.com/hashicorp/hcl/v2 v2.24.0 // indirect
	github.com/hashicorp/terraform-json v0.27.2 // indirect
	github.com/homeport/dyff v1.10.2 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgx/v5 v5.7.6 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/jinzhu/copier v0.4.0 // indirect
	github.com/klauspost/compress v1.18.1 // indirect
	github.com/kylelemons/godebug v1.1.0 // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/mattn/go-ciede2000 v0.0.0-20170301095244-782e8c62fec3 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-zglob v0.0.3 // indirect
	github.com/miekg/dns v1.1.68 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/go-ps v1.0.0 // indirect
	github.com/mitchellh/go-testing-interface v1.14.1 // indirect
	github.com/mitchellh/go-wordwrap v1.0.1 // indirect
	github.com/mitchellh/hashstructure v1.1.0 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.1.1 // indirect
	github.com/oracle/oci-go-sdk v24.3.0+incompatible // indirect
	github.com/pkg/browser v0.0.0-20240102092130-5ac0b6a4141c // indirect
	github.com/planetscale/vtprotobuf v0.6.1-0.20240319094008-0393e58bdf10 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/pquerna/otp v1.5.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/sergi/go-diff v1.4.0 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/texttheater/golang-levenshtein v1.0.1 // indirect
	github.com/tmccombs/hcl2json v0.6.8 // indirect
	github.com/ulikunitz/xz v0.5.10 // indirect
	github.com/urfave/cli/v2 v2.10.3 // indirect
	github.com/vbatts/tar-split v0.12.2 // indirect
	github.com/virtuald/go-ordered-json v0.0.0-20170621173500-b18e6e673d74 // indirect
	github.com/xrash/smetrics v0.0.0-20201216005158-039620a65673 // indirect
	github.com/zclconf/go-cty v1.16.4 // indirect
	go.opencensus.io v0.24.0 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/contrib/detectors/gcp v1.29.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.54.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.61.0 // indirect
	go.opentelemetry.io/otel v1.36.0 // indirect
	go.opentelemetry.io/otel/metric v1.36.0 // indirect
	go.opentelemetry.io/otel/sdk v1.36.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v1.36.0 // indirect
	go.opentelemetry.io/otel/trace v1.36.0 // indirect
	golang.org/x/crypto v0.45.0 // indirect
	golang.org/x/exp v0.0.0-20231006140011-7918f672742d // indirect
	golang.org/x/mod v0.30.0 // indirect
	golang.org/x/net v0.47.0 // indirect
	golang.org/x/oauth2 v0.33.0 // indirect
	golang.org/x/sync v0.18.0 // indirect
	golang.org/x/sys v0.38.0 // indirect
	golang.org/x/term v0.37.0 // indirect
	golang.org/x/text v0.31.0 // indirect
	golang.org/x/time v0.9.0 // indirect
	golang.org/x/tools v0.39.0 // indirect
	google.golang.org/api v0.206.0 // indirect
	google.golang.org/genproto v0.0.0-20241113202542-65e8d215514f // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20241209162323-e6fa225c2576 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20241209162323-e6fa225c2576 // indirect
	google.golang.org/grpc v1.68.1 // indirect
	google.golang.org/grpc/stats/opentelemetry v0.0.0-20240907200651-3ffb98b2c93a // indirect
	google.golang.org/protobuf v1.36.6 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/gruntwork-io/terratest/internal/lib/v2 => ../../internal/lib
	github.com/gruntwork-io/terratest/modules/aws/v2 => ../../modules/aws
	github.com/gruntwork-io/terratest/modules/azure/v2 => ../../modules/azure
	github.com/gruntwork-io/terratest/modules/collections/v2 => ../../modules/collections
	github.com/gruntwork-io/terratest/modules/dns-helper/v2 => ../../modules/dns-helper
	github.com/gruntwork-io/terratest/modules/environment/v2 => ../../modules/environment
	github.com/gruntwork-io/terratest/modules/files/v2 => ../../modules/files
	github.com/gruntwork-io/terratest/modules/gcp/v2 => ../../modules/gcp
	github.com/gruntwork-io/terratest/modules/git/v2 => ../../modules/git
	github.com/gruntwork-io/terratest/modules/helm/v2 => ../../modules/helm
	github.com/gruntwork-io/terratest/modules/http-helper/v2 => ../../modules/http-helper
	github.com/gruntwork-io/terratest/modules/k8s/v2 => ../../modules/k8s
	github.com/gruntwork-io/terratest/modules/logger/v2 => ../../modules/logger
	github.com/gruntwork-io/terratest/modules/oci/v2 => ../../modules/oci
	github.com/gruntwork-io/terratest/modules/opa/v2 => ../../modules/opa
	github.com/gruntwork-io/terratest/modules/random/v2 => ../../modules/random
	github.com/gruntwork-io/terratest/modules/retry/v2 => ../../modules/retry
	github.com/gruntwork-io/terratest/modules/shell/v2 => ../../modules/shell
	github.com/gruntwork-io/terratest/modules/ssh/v2 => ../../modules/ssh
	github.com/gruntwork-io/terratest/modules/terraform/v2 => ../../modules/terraform
	github.com/gruntwork-io/terratest/modules/terragrunt/v2 => ../../modules/terragrunt
	github.com/gruntwork-io/terratest/modules/testing/v2 => ../../modules/testing
)
