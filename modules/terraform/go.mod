module github.com/gruntwork-io/terratest/modules/terraform

replace github.com/gruntwork-io/terratest/modules/files => ../files

replace github.com/gruntwork-io/terratest/modules/collections => ../collections

replace github.com/gruntwork-io/terratest/modules/logger => ../logger

replace github.com/gruntwork-io/terratest/modules/retry => ../retry

replace github.com/gruntwork-io/terratest/modules/shell => ../shell

replace github.com/gruntwork-io/terratest/modules/ssh => ../ssh

replace github.com/gruntwork-io/terratest/modules/random => ../random

replace github.com/gruntwork-io/terratest/modules/customerrors => ../customerrors

go 1.13

require (
	github.com/gruntwork-io/terratest/modules/collections v0.0.0-00010101000000-000000000000
	github.com/gruntwork-io/terratest/modules/files v0.0.0-00010101000000-000000000000
	github.com/gruntwork-io/terratest/modules/logger v0.0.0-00010101000000-000000000000
	github.com/gruntwork-io/terratest/modules/retry v0.0.0-00010101000000-000000000000
	github.com/gruntwork-io/terratest/modules/shell v0.0.0-00010101000000-000000000000
	github.com/gruntwork-io/terratest/modules/ssh v0.0.0-00010101000000-000000000000
	github.com/magiconair/properties v1.8.1
	github.com/stretchr/testify v1.4.0
)
