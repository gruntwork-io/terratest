module github.com/gruntwork-io/terratest/modules/terraform

go 1.24.0

require (
	github.com/gruntwork-io/terratest/internal/lib v0.1.0
	github.com/gruntwork-io/terratest/modules/collections v0.1.0
	github.com/gruntwork-io/terratest/modules/files v0.1.0
	github.com/gruntwork-io/terratest/modules/http-helper v0.1.0
	github.com/gruntwork-io/terratest/modules/logger v0.1.0
	github.com/gruntwork-io/terratest/modules/opa v0.1.0
	github.com/gruntwork-io/terratest/modules/random v0.1.0
	github.com/gruntwork-io/terratest/modules/retry v0.1.0
	github.com/gruntwork-io/terratest/modules/shell v0.1.0
	github.com/gruntwork-io/terratest/modules/ssh v0.1.0
	github.com/gruntwork-io/terratest/modules/testing v0.1.0
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
