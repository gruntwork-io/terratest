module github.com/gruntwork-io/terratest/modules/dns-helper

go 1.24.0

require (
	github.com/gruntwork-io/terratest/modules/logger v0.1.0
	github.com/gruntwork-io/terratest/modules/retry v0.1.0
	github.com/gruntwork-io/terratest/modules/testing v0.1.0
	github.com/miekg/dns v1.1.68
	github.com/stretchr/testify v1.10.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/mod v0.24.0 // indirect
	golang.org/x/net v0.47.0 // indirect
	golang.org/x/sync v0.14.0 // indirect
	golang.org/x/sys v0.38.0 // indirect
	golang.org/x/tools v0.33.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/gruntwork-io/terratest/modules/logger => ../logger
	github.com/gruntwork-io/terratest/modules/retry => ../retry
	github.com/gruntwork-io/terratest/modules/testing => ../testing
)
