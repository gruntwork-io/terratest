// In-repo consumer-simulation project. Imports a representative cross-section of
// terratest v2 submodules under GOWORK=off, so CI can verify the external-consumer
// experience: each submodule resolves through its own go.mod, not the workspace.
//
// This module is NOT published. It exists only to be tested by CI.

module github.com/gruntwork-io/terratest/test-external

go 1.26

require (
	github.com/gruntwork-io/terratest/modules/core/v2 v2.0.0
	github.com/gruntwork-io/terratest/modules/shell/v2 v2.0.0
	github.com/stretchr/testify v1.11.1
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/gruntwork-io/terratest/modules/core/v2 => ../modules/core

replace github.com/gruntwork-io/terratest/modules/shell/v2 => ../modules/shell
