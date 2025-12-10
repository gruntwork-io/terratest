module github.com/gruntwork-io/terratest/modules/logger

go 1.24.0

replace (
	github.com/gruntwork-io/terratest/modules/files => ../files
	github.com/gruntwork-io/terratest/modules/random => ../random
	github.com/gruntwork-io/terratest/modules/shell => ../shell
	github.com/gruntwork-io/terratest/modules/testing => ../testing
)

exclude github.com/gruntwork-io/terratest v0.46.16

require (
	github.com/gruntwork-io/go-commons v0.17.2
	github.com/gruntwork-io/terratest/modules/files v1.0.0
	github.com/gruntwork-io/terratest/modules/random v1.0.0
	github.com/gruntwork-io/terratest/modules/shell v1.0.0
	github.com/gruntwork-io/terratest/modules/testing v1.0.0
	github.com/jstemmer/go-junit-report v1.0.0
	github.com/sirupsen/logrus v1.9.3
	github.com/stretchr/testify v1.11.1
)

require (
	github.com/cpuguy83/go-md2man/v2 v2.0.7 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/go-errors/errors v1.4.2 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/mattn/go-zglob v0.0.3 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/urfave/cli/v2 v2.10.3 // indirect
	github.com/xrash/smetrics v0.0.0-20201216005158-039620a65673 // indirect
	golang.org/x/sys v0.38.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
