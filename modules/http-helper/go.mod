module github.com/gruntwork-io/terratest/modules/http-helper

replace github.com/gruntwork-io/terratest/modules/logger => ../logger

replace github.com/gruntwork-io/terratest/modules/random => ../random

replace github.com/gruntwork-io/terratest/modules/retry => ../retry

go 1.13

require (
	github.com/gruntwork-io/terratest/modules/logger v0.0.0-00010101000000-000000000000
	github.com/gruntwork-io/terratest/modules/random v0.0.0-00010101000000-000000000000
	github.com/gruntwork-io/terratest/modules/retry v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.4.0
)
