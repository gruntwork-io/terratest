//go:build azure
// +build azure

// NOTE: We use build tags to differentiate azure testing because we currently do not have azure access setup for
// CircleCI.

package azure

import (
	"testing"

	"github.com/stretchr/testify/require"
)

/*
The below tests are currently stubbed out, with the expectation that they will throw errors. These tests can be extended.
*/

func TestListServiceBusNamespaceNamesE(t *testing.T) {
	t.Parallel()

	_, err := ListServiceBusNamespaceNamesContextE(t.Context(), "")

	require.Error(t, err)
}

func TestListServiceBusNamespaceIDsByResourceGroupE(t *testing.T) {
	t.Parallel()

	_, err := ListServiceBusNamespaceIDsByResourceGroupContextE(t.Context(), "", "")

	require.Error(t, err)
}

func TestListNamespaceAuthRulesE(t *testing.T) {
	t.Parallel()

	_, err := ListNamespaceAuthRulesContextE(t.Context(), "", "", "")

	require.Error(t, err)
}

func TestListNamespaceTopicsE(t *testing.T) {
	t.Parallel()

	_, err := ListNamespaceTopicsContextE(t.Context(), "", "", "")

	require.Error(t, err)
}

func TestListTopicAuthRulesE(t *testing.T) {
	t.Parallel()

	_, err := ListTopicAuthRulesContextE(t.Context(), "", "", "", "")

	require.Error(t, err)
}

func TestListTopicSubscriptionsNameE(t *testing.T) {
	t.Parallel()

	_, err := ListTopicSubscriptionsNameContextE(t.Context(), "", "", "", "")

	require.Error(t, err)
}
