module github.com/gruntwork-io/terratest/modules/k8s

replace github.com/gruntwork-io/terratest/modules/logger => ../logger

replace github.com/gruntwork-io/terratest/modules/random => ../random

replace github.com/gruntwork-io/terratest/modules/retry => ../retry

replace github.com/gruntwork-io/terratest/modules/customerrors => ../customerrors

replace github.com/gruntwork-io/terratest/modules/files => ../files

replace github.com/gruntwork-io/terratest/modules/ssh => ../ssh

replace github.com/gruntwork-io/terratest/modules/collections => ../collections

go 1.13

require (
	github.com/aws/aws-sdk-go v1.28.14
	github.com/go-sql-driver/mysql v1.5.0
	github.com/google/uuid v1.1.1
	github.com/gruntwork-io/terratest/modules/collections v0.0.0-00010101000000-000000000000
	github.com/gruntwork-io/terratest/modules/customerrors v0.0.0-00010101000000-000000000000
	github.com/gruntwork-io/terratest/modules/files v0.0.0-00010101000000-000000000000
	github.com/gruntwork-io/terratest/modules/logger v0.0.0-00010101000000-000000000000
	github.com/gruntwork-io/terratest/modules/random v0.0.0-00010101000000-000000000000
	github.com/gruntwork-io/terratest/modules/retry v0.0.0-00010101000000-000000000000
	github.com/gruntwork-io/terratest/modules/ssh v0.0.0-00010101000000-000000000000
	github.com/pquerna/otp v1.2.0
	github.com/stretchr/testify v1.4.0
)
