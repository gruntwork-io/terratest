module github.com/gruntwork-io/terratest/modules/terraform

go 1.25.0

require (
	github.com/gruntwork-io/terratest v0.46.16
	github.com/gruntwork-io/terratest/modules/files v0.0.0
	github.com/gruntwork-io/terratest/modules/formatting v0.0.0
	github.com/gruntwork-io/terratest/modules/logger v0.0.0
	github.com/gruntwork-io/terratest/modules/retry v0.0.0
	github.com/gruntwork-io/terratest/modules/shell v0.0.0
	github.com/gruntwork-io/terratest/modules/testing v0.0.0
	github.com/hashicorp/go-multierror v1.1.1
	github.com/hashicorp/hcl/v2 v2.22.0
	github.com/hashicorp/terraform-json v0.23.0
	github.com/jinzhu/copier v0.0.0-20190924061706-b57f9002281a
	github.com/stretchr/testify v1.11.1
	github.com/tmccombs/hcl2json v0.6.4
	github.com/zclconf/go-cty v1.15.0
)

require (
	github.com/agext/levenshtein v1.2.3 // indirect
	github.com/apparentlymart/go-textseg/v15 v15.0.0 // indirect
	github.com/bgentry/go-netrc v0.0.0-20140422174119-9fd32a8b3d3d // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/hashicorp/errwrap v1.0.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-getter/v2 v2.2.3 // indirect
	github.com/hashicorp/go-safetemp v1.0.0 // indirect
	github.com/hashicorp/go-version v1.7.0 // indirect
	github.com/klauspost/compress v1.16.7 // indirect
	github.com/mattn/go-zglob v0.0.6 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/go-testing-interface v1.14.1 // indirect
	github.com/mitchellh/go-wordwrap v1.0.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/ulikunitz/xz v0.5.10 // indirect
	golang.org/x/crypto v0.43.0 // indirect
	golang.org/x/mod v0.28.0 // indirect
	golang.org/x/net v0.46.0 // indirect
	golang.org/x/sync v0.17.0 // indirect
	golang.org/x/sys v0.37.0 // indirect
	golang.org/x/text v0.30.0 // indirect
	golang.org/x/tools v0.37.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/gruntwork-io/terratest => ../..
	github.com/gruntwork-io/terratest/modules/files => ../files
	github.com/gruntwork-io/terratest/modules/formatting => ../formatting
	github.com/gruntwork-io/terratest/modules/logger => ../logger
	github.com/gruntwork-io/terratest/modules/retry => ../retry
	github.com/gruntwork-io/terratest/modules/shell => ../shell
	github.com/gruntwork-io/terratest/modules/testing => ../testing
)
