module github.com/gruntwork-io/terratest/modules/ssh

replace github.com/gruntwork-io/terratest/modules/logger => ../logger

replace github.com/gruntwork-io/terratest/modules/collections => ../collections

replace github.com/gruntwork-io/terratest/modules/customerrors => ../customerrors

replace github.com/gruntwork-io/terratest/modules/files => ../files

go 1.13

require (
	github.com/gruntwork-io/terratest/modules/collections v0.0.0-00010101000000-000000000000
	github.com/gruntwork-io/terratest/modules/customerrors v0.0.0-00010101000000-000000000000
	github.com/gruntwork-io/terratest/modules/files v0.0.0-00010101000000-000000000000
	github.com/gruntwork-io/terratest/modules/logger v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.4.0
	golang.org/x/crypto v0.0.0-20200210191831-6ca56c2f2e2b
)
