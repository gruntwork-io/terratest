module github.com/gruntwork-io/terratest

go 1.13

replace github.com/gruntwork-io/terratest/modules/logger => ./modules/logger

require (
	github.com/aws/aws-lambda-go v1.14.1
	github.com/gruntwork-io/gruntwork-cli v0.6.1
	github.com/gruntwork-io/terratest/modules/logger v0.0.0-00010101000000-000000000000
	github.com/jstemmer/go-junit-report v0.9.1 // indirect
	github.com/sirupsen/logrus v1.4.2
	github.com/urfave/cli v1.22.2
)
