module github.com/gruntwork-io/terratest/modules/helm/v2

go 1.24.0

require (
	github.com/gonvenience/ytbx v1.4.7
	github.com/gruntwork-io/go-commons v0.17.2
	github.com/gruntwork-io/terratest/modules/files/v2 v2.0.0
	github.com/gruntwork-io/terratest/modules/http-helper/v2 v2.0.0
	github.com/gruntwork-io/terratest/modules/logger/v2 v2.0.0
	github.com/gruntwork-io/terratest/modules/random/v2 v2.0.0
	github.com/gruntwork-io/terratest/modules/shell/v2 v2.0.0
	github.com/gruntwork-io/terratest/modules/testing/v2 v2.0.0
	github.com/homeport/dyff v1.10.2
	github.com/stretchr/testify v1.11.1
	gopkg.in/yaml.v3 v3.0.1
	k8s.io/api v0.34.0
	k8s.io/apimachinery v0.34.0
)

require (
	github.com/BurntSushi/toml v1.5.0 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.7 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/fxamacker/cbor/v2 v2.9.0 // indirect
	github.com/go-errors/errors v1.4.2 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/gonvenience/bunt v1.4.2 // indirect
	github.com/gonvenience/idem v0.0.2 // indirect
	github.com/gonvenience/neat v1.3.16 // indirect
	github.com/gonvenience/term v1.0.4 // indirect
	github.com/gonvenience/text v1.0.9 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/mattn/go-ciede2000 v0.0.0-20170301095244-782e8c62fec3 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-zglob v0.0.3 // indirect
	github.com/mitchellh/go-ps v1.0.0 // indirect
	github.com/mitchellh/hashstructure v1.1.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.3-0.20250322232337-35a7c28c31ee // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/sergi/go-diff v1.4.0 // indirect
	github.com/spf13/pflag v1.0.9 // indirect
	github.com/texttheater/golang-levenshtein v1.0.1 // indirect
	github.com/urfave/cli/v2 v2.10.3 // indirect
	github.com/virtuald/go-ordered-json v0.0.0-20170621173500-b18e6e673d74 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	github.com/xrash/smetrics v0.0.0-20201216005158-039620a65673 // indirect
	go.yaml.in/yaml/v2 v2.4.2 // indirect
	golang.org/x/exp v0.0.0-20231006140011-7918f672742d // indirect
	golang.org/x/net v0.47.0 // indirect
	golang.org/x/sync v0.18.0 // indirect
	golang.org/x/sys v0.38.0 // indirect
	golang.org/x/term v0.37.0 // indirect
	golang.org/x/text v0.31.0 // indirect
	golang.org/x/tools v0.39.0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	k8s.io/klog/v2 v2.130.1 // indirect
	k8s.io/utils v0.0.0-20250604170112-4c0f3b243397 // indirect
	sigs.k8s.io/json v0.0.0-20241014173422-cfa47c3a1cc8 // indirect
	sigs.k8s.io/randfill v1.0.0 // indirect
	sigs.k8s.io/structured-merge-diff/v6 v6.3.0 // indirect
)

replace (
	github.com/gruntwork-io/terratest/modules/files/v2 => ../files
	github.com/gruntwork-io/terratest/modules/http-helper/v2 => ../http-helper
	github.com/gruntwork-io/terratest/modules/k8s/v2 => ../k8s
	github.com/gruntwork-io/terratest/modules/logger/v2 => ../logger
	github.com/gruntwork-io/terratest/modules/random/v2 => ../random
	github.com/gruntwork-io/terratest/modules/shell/v2 => ../shell
	github.com/gruntwork-io/terratest/modules/testing/v2 => ../testing
)
