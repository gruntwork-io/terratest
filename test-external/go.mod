// In-repo consumer-simulation project. Imports a cross-section of terratest v2
// submodules under GOWORK=off so CI can verify the external-consumer experience:
// each submodule resolves through its own go.mod and the proxy, not the workspace.
// NOT published; exists only to be tested by CI (red until the v2 tags are cut).
module github.com/gruntwork-io/terratest/test-external

go 1.26

require (
	github.com/gruntwork-io/terratest/modules/core/v2 v2.0.0
	github.com/gruntwork-io/terratest/modules/terraform/v2 v2.0.0
	github.com/stretchr/testify v1.11.1
)
