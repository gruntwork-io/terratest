package k8s_test

import (
	"context"
	"encoding/json"
	"testing"

	gotesting "github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/gruntwork-io/terratest/modules/k8s/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestKubectlOptionsMarshalsWithLookupSet guards the json:"-" tag on NodePublicIPLookup. KubectlOptions is
// serialized as test data, and encoding/json fails on a func field, so dropping the tag would break every caller
// that saves options with a lookup configured.
func TestKubectlOptionsMarshalsWithLookupSet(t *testing.T) {
	t.Parallel()

	options := k8s.NewKubectlOptions("terratest-context", "~/.kube/config", "default")
	options.NodePublicIPLookup = func(_ gotesting.TestingT, _ context.Context, _ []string, _ string) (map[string]string, error) {
		return nil, nil
	}

	// nolint below: musttag because KubectlOptions has no json tags, and staticcheck SA1026 because it sees the
	// func inside RestConfig. RestConfig is nil here, so the only func in play is NodePublicIPLookup, which is
	// exactly what this test is pinning. Note that marshalling options with RestConfig set fails independently of
	// this field; that is pre-existing behaviour on main and out of scope for this change.
	raw, err := json.Marshal(options) //nolint:musttag,staticcheck // see comment above
	require.NoError(t, err, "options carrying a lookup func must still marshal")
	assert.NotContains(t, string(raw), "NodePublicIPLookup", "the func field must be omitted from JSON")

	var round k8s.KubectlOptions

	require.NoError(t, json.Unmarshal(raw, &round)) //nolint:musttag // KubectlOptions does not have json tags
	assert.Equal(t, "terratest-context", round.ContextName)
	assert.Equal(t, "default", round.Namespace)
	assert.Nil(t, round.NodePublicIPLookup, "func field is intentionally not persisted")
}
