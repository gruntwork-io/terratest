module github.com/gruntwork-io/terratest/modules/opa/v2

go 1.24.0

require (
	github.com/gruntwork-io/terratest/modules/files/v2 v2.0.0
	github.com/gruntwork-io/terratest/modules/git/v2 v2.0.0
	github.com/gruntwork-io/terratest/modules/logger/v2 v2.0.0
	github.com/gruntwork-io/terratest/modules/shell/v2 v2.0.0
	github.com/gruntwork-io/terratest/modules/testing/v2 v2.0.0
	github.com/hashicorp/go-getter/v2 v2.2.3
	github.com/hashicorp/go-multierror v1.1.1
	github.com/stretchr/testify v1.11.1
)

require (
	github.com/bgentry/go-netrc v0.0.0-20140422174119-9fd32a8b3d3d // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/hashicorp/errwrap v1.0.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-safetemp v1.0.0 // indirect
	github.com/hashicorp/go-version v1.7.0 // indirect
	github.com/klauspost/compress v1.18.1 // indirect
	github.com/mattn/go-zglob v0.0.3 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/go-testing-interface v1.14.1 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/ulikunitz/xz v0.5.10 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/gruntwork-io/terratest/modules/files/v2 => ../files
	github.com/gruntwork-io/terratest/modules/git/v2 => ../git
	github.com/gruntwork-io/terratest/modules/logger/v2 => ../logger
	github.com/gruntwork-io/terratest/modules/shell/v2 => ../shell
	github.com/gruntwork-io/terratest/modules/testing/v2 => ../testing
)
