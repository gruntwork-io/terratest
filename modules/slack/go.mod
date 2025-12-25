module github.com/gruntwork-io/terratest/modules/slack/v2

go 1.24.0

require (
	github.com/gruntwork-io/terratest/modules/environment/v2 v2.0.0
	github.com/gruntwork-io/terratest/modules/random/v2 v2.0.0
	github.com/gruntwork-io/terratest/modules/retry/v2 v2.0.0
	github.com/gruntwork-io/terratest/modules/testing/v2 v2.0.0
	github.com/slack-go/slack v0.17.3
	github.com/stretchr/testify v1.11.1
)

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/gorilla/websocket v1.5.4-0.20250319132907-e064f32e3674 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	golang.org/x/net v0.47.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/gruntwork-io/terratest/modules/environment/v2 => ../environment
	github.com/gruntwork-io/terratest/modules/random/v2 => ../random
	github.com/gruntwork-io/terratest/modules/retry/v2 => ../retry
	github.com/gruntwork-io/terratest/modules/testing/v2 => ../testing
)
