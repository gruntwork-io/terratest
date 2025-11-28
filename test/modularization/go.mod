module github.com/gruntwork-io/terratest/test/modularization

go 1.24.0

require (
	github.com/gruntwork-io/terratest/modules/collections v0.55.0
	github.com/gruntwork-io/terratest/modules/dns-helper v0.55.0
	github.com/gruntwork-io/terratest/modules/environment v0.55.0
	github.com/gruntwork-io/terratest/modules/git v0.55.0
	github.com/gruntwork-io/terratest/modules/http-helper v0.55.0
	github.com/gruntwork-io/terratest/modules/logger v0.55.0
	github.com/gruntwork-io/terratest/modules/oci v0.55.0
	github.com/gruntwork-io/terratest/modules/ssh v0.55.0
	github.com/gruntwork-io/terratest/modules/terragrunt v0.55.0
	github.com/gruntwork-io/terratest/modules/testing v0.55.0
	github.com/stretchr/testify v1.11.1
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/gruntwork-io/terratest/internal/lib v0.55.0 // indirect
	github.com/gruntwork-io/terratest/modules/files v0.55.0 // indirect
	github.com/gruntwork-io/terratest/modules/random v0.55.0 // indirect
	github.com/gruntwork-io/terratest/modules/retry v0.55.0 // indirect
	github.com/gruntwork-io/terratest/modules/shell v0.55.0 // indirect
	github.com/hashicorp/errwrap v1.0.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/mattn/go-zglob v0.0.2-0.20190814121620-e3c945676326 // indirect
	github.com/miekg/dns v1.1.68 // indirect
	github.com/oracle/oci-go-sdk v24.3.0+incompatible // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/crypto v0.44.0 // indirect
	golang.org/x/mod v0.24.0 // indirect
	golang.org/x/net v0.47.0 // indirect
	golang.org/x/sync v0.14.0 // indirect
	golang.org/x/sys v0.38.0 // indirect
	golang.org/x/tools v0.33.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/gruntwork-io/terratest/internal/lib => ../../internal/lib
	github.com/gruntwork-io/terratest/modules/collections => ../../modules/collections
	github.com/gruntwork-io/terratest/modules/dns-helper => ../../modules/dns-helper
	github.com/gruntwork-io/terratest/modules/environment => ../../modules/environment
	github.com/gruntwork-io/terratest/modules/files => ../../modules/files
	github.com/gruntwork-io/terratest/modules/git => ../../modules/git
	github.com/gruntwork-io/terratest/modules/http-helper => ../../modules/http-helper
	github.com/gruntwork-io/terratest/modules/logger => ../../modules/logger
	github.com/gruntwork-io/terratest/modules/oci => ../../modules/oci
	github.com/gruntwork-io/terratest/modules/random => ../../modules/random
	github.com/gruntwork-io/terratest/modules/retry => ../../modules/retry
	github.com/gruntwork-io/terratest/modules/shell => ../../modules/shell
	github.com/gruntwork-io/terratest/modules/ssh => ../../modules/ssh
	github.com/gruntwork-io/terratest/modules/terragrunt => ../../modules/terragrunt
	github.com/gruntwork-io/terratest/modules/testing => ../../modules/testing
)
