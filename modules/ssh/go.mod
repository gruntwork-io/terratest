module github.com/gruntwork-io/terratest/modules/ssh

go 1.24.0

require (
	github.com/gruntwork-io/terratest/modules/collections v0.55.0
	github.com/gruntwork-io/terratest/modules/files v0.55.0
	github.com/gruntwork-io/terratest/modules/logger v0.55.0
	github.com/gruntwork-io/terratest/modules/retry v0.55.0
	github.com/gruntwork-io/terratest/modules/testing v0.55.0
	github.com/hashicorp/go-multierror v1.1.1
	github.com/stretchr/testify v1.11.1
	golang.org/x/crypto v0.44.0
)

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/hashicorp/errwrap v1.0.0 // indirect
	github.com/mattn/go-zglob v0.0.3 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	golang.org/x/net v0.47.0 // indirect
	golang.org/x/sys v0.38.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/gruntwork-io/terratest/modules/collections => ../collections
	github.com/gruntwork-io/terratest/modules/files => ../files
	github.com/gruntwork-io/terratest/modules/logger => ../logger
	github.com/gruntwork-io/terratest/modules/retry => ../retry
	github.com/gruntwork-io/terratest/modules/testing => ../testing
)
