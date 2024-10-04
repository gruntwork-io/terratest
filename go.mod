module github.com/gruntwork-io/terratest

go 1.21

require (
	cloud.google.com/go v0.110.0 // indirect
	cloud.google.com/go/storage v1.28.1
	github.com/Azure/azure-sdk-for-go v51.0.0+incompatible
	github.com/Azure/go-autorest/autorest v0.11.20
	github.com/Azure/go-autorest/autorest/azure/auth v0.5.8
	github.com/Azure/go-autorest/autorest/to v0.4.0 // indirect
	github.com/Azure/go-autorest/autorest/validation v0.3.1 // indirect
	github.com/aws/aws-lambda-go v1.13.3
	github.com/aws/aws-sdk-go v1.44.122
	github.com/ghodss/yaml v1.0.0
	github.com/go-errors/errors v1.0.2-0.20180813162953-d98b870cc4e0 // indirect
	github.com/go-sql-driver/mysql v1.4.1
	github.com/google/go-containerregistry v0.6.0
	github.com/google/uuid v1.3.0
	github.com/gruntwork-io/go-commons v0.8.0
	github.com/hashicorp/go-getter v1.7.6
	github.com/hashicorp/go-multierror v1.1.0
	github.com/hashicorp/go-version v1.6.0
	github.com/hashicorp/hcl/v2 v2.9.1
	github.com/hashicorp/terraform-json v0.13.0
	github.com/jinzhu/copier v0.0.0-20190924061706-b57f9002281a
	github.com/jstemmer/go-junit-report v0.9.1
	github.com/magiconair/properties v1.8.5
	github.com/mattn/go-zglob v0.0.2-0.20190814121620-e3c945676326
	github.com/miekg/dns v1.1.31
	github.com/mitchellh/go-homedir v1.1.0
	github.com/oracle/oci-go-sdk v7.1.0+incompatible
	github.com/pquerna/otp v1.2.0
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.8.4
	github.com/tmccombs/hcl2json v0.3.3
	github.com/urfave/cli v1.22.2
	github.com/zclconf/go-cty v1.9.1
	golang.org/x/crypto v0.21.0
	golang.org/x/net v0.23.0
	golang.org/x/oauth2 v0.8.0
	google.golang.org/api v0.114.0
	google.golang.org/genproto v0.0.0-20230410155749-daa745c078e1
	k8s.io/api v0.28.4
	k8s.io/apimachinery v0.28.4
	k8s.io/client-go v0.28.4
)

require (
	cloud.google.com/go/cloudbuild v1.9.0
	github.com/aws/aws-sdk-go-v2 v1.32.0
	github.com/aws/aws-sdk-go-v2/config v1.27.41
	github.com/aws/aws-sdk-go-v2/credentials v1.17.39
	github.com/aws/aws-sdk-go-v2/feature/s3/manager v1.17.28
	github.com/aws/aws-sdk-go-v2/service/acm v1.30.0
	github.com/aws/aws-sdk-go-v2/service/autoscaling v1.45.0
	github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs v1.41.0
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.36.0
	github.com/aws/aws-sdk-go-v2/service/ec2 v1.181.0
	github.com/aws/aws-sdk-go-v2/service/ecr v1.36.0
	github.com/aws/aws-sdk-go-v2/service/ecs v1.47.0
	github.com/aws/aws-sdk-go-v2/service/iam v1.37.0
	github.com/aws/aws-sdk-go-v2/service/kms v1.37.0
	github.com/aws/aws-sdk-go-v2/service/lambda v1.63.0
	github.com/aws/aws-sdk-go-v2/service/rds v1.87.0
	github.com/aws/aws-sdk-go-v2/service/route53 v1.45.0
	github.com/aws/aws-sdk-go-v2/service/s3 v1.65.0
	github.com/aws/aws-sdk-go-v2/service/secretsmanager v1.34.0
	github.com/aws/aws-sdk-go-v2/service/sns v1.33.0
	github.com/aws/aws-sdk-go-v2/service/sqs v1.36.0
	github.com/aws/aws-sdk-go-v2/service/ssm v1.55.0
	github.com/aws/aws-sdk-go-v2/service/sts v1.32.0
	github.com/gonvenience/ytbx v1.4.4
	github.com/homeport/dyff v1.6.0
	github.com/slack-go/slack v0.10.3
	gotest.tools/v3 v3.0.3
)

require (
	cloud.google.com/go/compute v1.19.1 // indirect
	cloud.google.com/go/compute/metadata v0.2.3 // indirect
	cloud.google.com/go/iam v0.13.0 // indirect
	cloud.google.com/go/longrunning v0.4.1 // indirect
	github.com/Azure/go-autorest v14.2.0+incompatible // indirect
	github.com/Azure/go-autorest/autorest/adal v0.9.13 // indirect
	github.com/Azure/go-autorest/autorest/azure/cli v0.4.2 // indirect
	github.com/Azure/go-autorest/autorest/date v0.3.0 // indirect
	github.com/Azure/go-autorest/logger v0.2.1 // indirect
	github.com/Azure/go-autorest/tracing v0.6.0 // indirect
	github.com/BurntSushi/toml v1.3.2 // indirect
	github.com/agext/levenshtein v1.2.3 // indirect
	github.com/apparentlymart/go-textseg/v13 v13.0.0 // indirect
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.6.6 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.16.15 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.19 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.19 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.1 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.3.19 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.12.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.4.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.10.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.12.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.18.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.24.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.28.0 // indirect
	github.com/aws/smithy-go v1.22.0 // indirect
	github.com/bgentry/go-netrc v0.0.0-20140422174119-9fd32a8b3d3d // indirect
	github.com/boombuler/barcode v1.0.1-0.20190219062509-6c824513bacc // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dimchansky/utfbom v1.1.1 // indirect
	github.com/docker/cli v20.10.9+incompatible // indirect
	github.com/docker/distribution v2.8.2+incompatible // indirect
	github.com/docker/docker v25.0.6+incompatible // indirect
	github.com/docker/docker-credential-helpers v0.6.3 // indirect
	github.com/emicklei/go-restful/v3 v3.9.0 // indirect
	github.com/form3tech-oss/jwt-go v3.2.2+incompatible // indirect
	github.com/go-logr/logr v1.2.4 // indirect
	github.com/go-openapi/jsonpointer v0.19.6 // indirect
	github.com/go-openapi/jsonreference v0.20.2 // indirect
	github.com/go-openapi/swag v0.22.3 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/gonvenience/bunt v1.3.5 // indirect
	github.com/gonvenience/neat v1.3.12 // indirect
	github.com/gonvenience/term v1.0.2 // indirect
	github.com/gonvenience/text v1.0.7 // indirect
	github.com/gonvenience/wrap v1.1.2 // indirect
	github.com/google/gnostic-models v0.6.8 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.2.3 // indirect
	github.com/googleapis/gax-go/v2 v2.7.1 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/hashicorp/errwrap v1.0.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-safetemp v1.0.0 // indirect
	github.com/imdario/mergo v0.3.11 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/compress v1.15.11 // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-ciede2000 v0.0.0-20170301095244-782e8c62fec3 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/mitchellh/go-ps v1.0.0 // indirect
	github.com/mitchellh/go-testing-interface v1.14.1 // indirect
	github.com/mitchellh/go-wordwrap v1.0.1 // indirect
	github.com/mitchellh/hashstructure v1.1.0 // indirect
	github.com/moby/spdystream v0.2.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.0.2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/satori/go.uuid v1.2.0 // indirect
	github.com/sergi/go-diff v1.3.1 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/texttheater/golang-levenshtein v1.0.1 // indirect
	github.com/ulikunitz/xz v0.5.10 // indirect
	github.com/virtuald/go-ordered-json v0.0.0-20170621173500-b18e6e673d74 // indirect
	go.opencensus.io v0.24.0 // indirect
	golang.org/x/sync v0.4.0 // indirect
	golang.org/x/sys v0.18.0 // indirect
	golang.org/x/term v0.18.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	golang.org/x/time v0.3.0 // indirect
	golang.org/x/xerrors v0.0.0-20220907171357-04be3eba64a2 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/grpc v1.56.3 // indirect
	google.golang.org/protobuf v1.33.0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/klog/v2 v2.100.1 // indirect
	k8s.io/kube-openapi v0.0.0-20230717233707-2695361300d9 // indirect
	k8s.io/utils v0.0.0-20230406110748-d93618cff8a2 // indirect
	sigs.k8s.io/json v0.0.0-20221116044647-bc3834ca7abd // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.2.3 // indirect
	sigs.k8s.io/yaml v1.3.0 // indirect
)
