module github.com/gruntwork-io/terratest/modules/logger

go 1.25.0

require (
	github.com/gruntwork-io/go-commons v0.8.0
	github.com/gruntwork-io/terratest v0.46.16
	github.com/gruntwork-io/terratest/modules/files v0.0.0
	github.com/gruntwork-io/terratest/modules/shell v0.0.0
	github.com/gruntwork-io/terratest/modules/testing v0.0.0
	github.com/jstemmer/go-junit-report v1.0.0
	github.com/sirupsen/logrus v1.9.3
	github.com/stretchr/testify v1.11.1
)

require (
	github.com/cpuguy83/go-md2man/v2 v2.0.5 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-errors/errors v1.4.2 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/mattn/go-zglob v0.0.6 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/urfave/cli v1.22.16 // indirect
	golang.org/x/sys v0.37.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/gruntwork-io/terratest => ../..
	github.com/gruntwork-io/terratest/modules/files => ../files
	github.com/gruntwork-io/terratest/modules/shell => ../shell
	github.com/gruntwork-io/terratest/modules/testing => ../testing
)
