package external_test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/core/v2/random"
	"github.com/gruntwork-io/terratest/modules/terraform/v2"
	"github.com/stretchr/testify/assert"
)

// Proves a real external consumer can pull a tier-0 module (core) and a top-tier
// module (terraform, which transitively requires httphelper/opa/ssh) purely through
// published go.mods under GOWORK=off.
func TestExternalConsumer(t *testing.T) {
	t.Parallel()

	assert.Len(t, random.UniqueId(), 6)

	opts := terraform.Options{TerraformDir: "."}
	assert.Equal(t, ".", opts.TerraformDir)
}
