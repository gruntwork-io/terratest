package k8s_test

import (
	"context"
	"encoding/json"
	"net/http"
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

// TestKubectlOptionsMarshalsWithRestConfigSet covers the other unserializable field. rest.Config holds func-typed
// fields, and encoding/json rejects a func field whether or not it is set, so before the json:"-" tag any marshal
// of options built by NewKubectlOptionsWithRestConfig failed with "unsupported type: transport.WrapperFunc".
func TestKubectlOptionsMarshalsWithRestConfigSet(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name   string
		config *rest.Config
	}{
		{"plain config", &rest.Config{Host: "https://example.com"}},
		{"config with WrapTransport set", &rest.Config{
			Host: "https://example.com",
			WrapTransport: func(rt http.RoundTripper) http.RoundTripper {
				return rt
			},
		}},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			options := k8s.NewKubectlOptionsWithRestConfig(testCase.config, "default")

			raw, err := json.Marshal(options) //nolint:musttag // KubectlOptions does not have json tags
			require.NoError(t, err, "options carrying a RestConfig must still marshal")
			assert.NotContains(t, string(raw), "RestConfig", "the config must be omitted from JSON")
		})
	}
}
