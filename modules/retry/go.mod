module github.com/gruntwork-io/terratest/modules/retry

replace github.com/gruntwork-io/terratest/modules/logger => ../logger

go 1.13

require (
	github.com/gruntwork-io/terratest/modules/logger v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.4.0
	golang.org/x/net v0.0.0-20200202094626-16171245cfb2
)
