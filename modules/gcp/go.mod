module github.com/gruntwork-io/terratest/modules/gcp/v2

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

replace github.com/gruntwork-io/terratest/modules/aws/v2 => ../aws

replace github.com/gruntwork-io/terratest/modules/azure/v2 => ../azure

replace github.com/gruntwork-io/terratest/modules/k8s/v2 => ../k8s

replace github.com/gruntwork-io/terratest/modules/helm/v2 => ../helm

replace github.com/gruntwork-io/terratest/modules/terraform/v2 => ../terraform

replace github.com/gruntwork-io/terratest/modules/terragrunt/v2 => ../terragrunt

replace github.com/gruntwork-io/terratest/modules/test-structure/v2 => ../test-structure

require (
	cloud.google.com/go/cloudbuild v1.30.0
	cloud.google.com/go/pubsub/v2 v2.6.0
	cloud.google.com/go/storage v1.62.2
	github.com/google/go-containerregistry v0.21.6
	github.com/gruntwork-io/terratest/modules/core/v2 v2.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.11.1
	golang.org/x/oauth2 v0.36.0
	google.golang.org/api v0.280.0
	google.golang.org/grpc v1.81.1
)

require (
	cel.dev/expr v0.25.1 // indirect
	cloud.google.com/go v0.123.0 // indirect
	cloud.google.com/go/auth v0.20.0 // indirect
	cloud.google.com/go/auth/oauth2adapt v0.2.8 // indirect
	cloud.google.com/go/compute/metadata v0.9.0 // indirect
	cloud.google.com/go/iam v1.7.0 // indirect
	cloud.google.com/go/longrunning v0.9.0 // indirect
	cloud.google.com/go/monitoring v1.24.3 // indirect
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/detectors/gcp v1.31.0 // indirect
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/metric v0.55.0 // indirect
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/internal/resourcemapping v0.55.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/cncf/xds/go v0.0.0-20260202195803-dba9d589def2 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/docker/cli v29.4.3+incompatible // indirect
	github.com/docker/docker-credential-helpers v0.9.3 // indirect
	github.com/envoyproxy/go-control-plane/envoy v1.37.0 // indirect
	github.com/envoyproxy/protoc-gen-validate v1.3.3 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/go-jose/go-jose/v4 v4.1.4 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/google/s2a-go v0.1.9 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.3.15 // indirect
	github.com/googleapis/gax-go/v2 v2.22.0 // indirect
	github.com/klauspost/compress v1.18.6 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.1.1 // indirect
	github.com/planetscale/vtprotobuf v0.6.1-0.20240319094008-0393e58bdf10 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/sirupsen/logrus v1.9.4 // indirect
	github.com/spiffe/go-spiffe/v2 v2.6.0 // indirect
	go.einride.tech/aip v0.83.0 // indirect
	go.opencensus.io v0.24.0 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/contrib/detectors/gcp v1.42.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.67.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.67.0 // indirect
	go.opentelemetry.io/otel v1.43.0 // indirect
	go.opentelemetry.io/otel/metric v1.43.0 // indirect
	go.opentelemetry.io/otel/sdk v1.43.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v1.43.0 // indirect
	go.opentelemetry.io/otel/trace v1.43.0 // indirect
	golang.org/x/crypto v0.51.0 // indirect
	golang.org/x/net v0.54.0 // indirect
	golang.org/x/sync v0.20.0 // indirect
	golang.org/x/sys v0.44.0 // indirect
	golang.org/x/text v0.37.0 // indirect
	golang.org/x/time v0.15.0 // indirect
	google.golang.org/genproto v0.0.0-20260319201613-d00831a3d3e7 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20260401024825-9d38bb4040a9 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260511170946-3700d4141b60 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
