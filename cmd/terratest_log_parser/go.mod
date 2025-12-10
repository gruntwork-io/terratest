module github.com/gruntwork-io/terratest/cmd/terratest_log_parser

go 1.24.0

replace (
	github.com/gruntwork-io/terratest/modules/collections => ../../modules/collections
	github.com/gruntwork-io/terratest/modules/files => ../../modules/files
	github.com/gruntwork-io/terratest/modules/logger => ../../modules/logger
	github.com/gruntwork-io/terratest/modules/random => ../../modules/random
	github.com/gruntwork-io/terratest/modules/shell => ../../modules/shell
	github.com/gruntwork-io/terratest/modules/testing => ../../modules/testing
)

exclude github.com/gruntwork-io/terratest v0.46.16

require (
	github.com/gruntwork-io/go-commons v0.17.2
	github.com/gruntwork-io/terratest/modules/logger v1.0.0
	github.com/sirupsen/logrus v1.9.3
	github.com/urfave/cli/v2 v2.10.3
)

require (
	github.com/cpuguy83/go-md2man/v2 v2.0.7 // indirect
	github.com/go-errors/errors v1.4.2 // indirect
	github.com/jstemmer/go-junit-report v1.0.0 // indirect
	github.com/mattn/go-zglob v0.0.3 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/xrash/smetrics v0.0.0-20201216005158-039620a65673 // indirect
	golang.org/x/sys v0.38.0 // indirect
)
