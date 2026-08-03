package k8s_test

import (
	"context"
	"encoding/json"
	"testing"

	gotesting "github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/gruntwork-io/terratest/modules/k8s/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/client-go/rest"
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

	raw, err := json.Marshal(options) //nolint:musttag // KubectlOptions does not have json tags
	require.NoError(t, err, "options carrying a lookup func must still marshal")
	assert.NotContains(t, string(raw), "NodePublicIPLookup", "the func field must be omitted from JSON")

	var round k8s.KubectlOptions

	require.NoError(t, json.Unmarshal(raw, &round)) //nolint:musttag // KubectlOptions does not have json tags
	assert.Equal(t, "terratest-context", round.ContextName)
	assert.Equal(t, "default", round.Namespace)
	assert.Nil(t, round.NodePublicIPLookup, "func field is intentionally not persisted")
}

// TestMarshalJSONRejectsRestConfigOnEveryPath is the point of putting the check in MarshalJSON rather than in
// SaveKubectlOptions. KubectlOptions reaches disk through several writers: the generic teststructure.SaveTestData,
// helm.Options which embeds it, and plain json.Marshal in user code. A check in one saver would leave the others
// silently dropping the config, which is worse than the crash this PR fixes.
func TestMarshalJSONRejectsRestConfigOnEveryPath(t *testing.T) {
	t.Parallel()

	options := k8s.NewKubectlOptionsWithRestConfig(&rest.Config{Host: "https://example.com"}, "default")

	t.Run("direct marshal of the value", func(t *testing.T) {
		t.Parallel()

		_, err := json.Marshal(options) //nolint:musttag // KubectlOptions does not have json tags
		require.ErrorIs(t, err, k8s.ErrRestConfigNotSerializable)
	})

	t.Run("marshal when embedded in another struct", func(t *testing.T) {
		t.Parallel()

		// helm.Options embeds *k8s.KubectlOptions, so this is the shape that would otherwise slip through.
		wrapper := struct {
			KubectlOptions *k8s.KubectlOptions
		}{KubectlOptions: options}

		_, err := json.Marshal(wrapper) //nolint:musttag // KubectlOptions does not have json tags
		require.ErrorIs(t, err, k8s.ErrRestConfigNotSerializable)
	})

	t.Run("options without a RestConfig still marshal", func(t *testing.T) {
		t.Parallel()

		plain := k8s.NewKubectlOptions("ctx", "~/.kube/config", "default")

		raw, err := json.Marshal(plain) //nolint:musttag // KubectlOptions does not have json tags
		require.NoError(t, err)
		assert.Contains(t, string(raw), "ctx")
		assert.NotContains(t, string(raw), "example.com", "no cluster host should leak into the JSON")
	})

	t.Run("in-cluster options round trip", func(t *testing.T) {
		t.Parallel()

		raw, err := json.Marshal(k8s.NewKubectlOptionsWithInClusterAuth()) //nolint:musttag // no json tags
		require.NoError(t, err)

		var round k8s.KubectlOptions

		require.NoError(t, json.Unmarshal(raw, &round)) //nolint:musttag // no json tags
		assert.True(t, round.InClusterAuth, "the documented alternative must survive a round trip")
	})
}
