module github.com/gruntwork-io/terratest/modules/terraform

go 1.24.0

require (
	github.com/gruntwork-io/terratest/internal/lib v1.0.0
	github.com/gruntwork-io/terratest/modules/collections v1.0.0
	github.com/gruntwork-io/terratest/modules/files v1.0.0
	github.com/gruntwork-io/terratest/modules/http-helper v1.0.0
	github.com/gruntwork-io/terratest/modules/logger v1.0.0
	github.com/gruntwork-io/terratest/modules/opa v1.0.0
	github.com/gruntwork-io/terratest/modules/random v1.0.0
	github.com/gruntwork-io/terratest/modules/retry v1.0.0
	github.com/gruntwork-io/terratest/modules/shell v1.0.0
	github.com/gruntwork-io/terratest/modules/ssh v1.0.0
	github.com/gruntwork-io/terratest/modules/testing v1.0.0
	github.com/hashicorp/go-multierror v1.1.1
	github.com/hashicorp/hcl/v2 v2.24.0
	github.com/hashicorp/terraform-json v0.27.2
	github.com/jinzhu/copier v0.4.0
	github.com/stretchr/testify v1.11.1
	github.com/tmccombs/hcl2json v0.6.8
	github.com/zclconf/go-cty v1.16.4
)

require (
	github.com/agext/levenshtein v1.2.3 // indirect
	github.com/apparentlymart/go-textseg/v15 v15.0.0 // indirect
	github.com/bgentry/go-netrc v0.0.0-20140422174119-9fd32a8b3d3d // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/go-test/deep v1.1.1 // indirect
	github.com/hashicorp/errwrap v1.0.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-getter/v2 v2.2.3 // indirect
	github.com/hashicorp/go-safetemp v1.0.0 // indirect
	github.com/hashicorp/go-version v1.7.0 // indirect
	github.com/klauspost/compress v1.18.1 // indirect
	github.com/mattn/go-zglob v0.0.3 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/go-testing-interface v1.14.1 // indirect
	github.com/mitchellh/go-wordwrap v1.0.1 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/ulikunitz/xz v0.5.10 // indirect
	golang.org/x/crypto v0.44.0 // indirect
	golang.org/x/mod v0.30.0 // indirect
	golang.org/x/net v0.47.0 // indirect
	golang.org/x/sync v0.18.0 // indirect
	golang.org/x/sys v0.38.0 // indirect
	golang.org/x/text v0.31.0 // indirect
	golang.org/x/tools v0.39.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/gruntwork-io/terratest/internal/lib => ../../internal/lib
	github.com/gruntwork-io/terratest/modules/collections => ../collections
	github.com/gruntwork-io/terratest/modules/files => ../files
	github.com/gruntwork-io/terratest/modules/http-helper => ../http-helper
	github.com/gruntwork-io/terratest/modules/logger => ../logger
	github.com/gruntwork-io/terratest/modules/opa => ../opa
	github.com/gruntwork-io/terratest/modules/random => ../random
	github.com/gruntwork-io/terratest/modules/retry => ../retry
	github.com/gruntwork-io/terratest/modules/shell => ../shell
	github.com/gruntwork-io/terratest/modules/ssh => ../ssh
	github.com/gruntwork-io/terratest/modules/testing => ../testing
)
