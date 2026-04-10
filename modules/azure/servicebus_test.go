//go:build azure
// +build azure

// NOTE: We use build tags to differentiate azure testing because we currently do not have azure access setup for
// CircleCI.

package azure_test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/stretchr/testify/require"
)

/*
The below tests are currently stubbed out, with the expectation that they will throw errors. These tests can be extended.
*/

func TestListServiceBusNamespaceNamesE(t *testing.T) {
	t.Parallel()

	_, err := azure.ListServiceBusNamespaceNamesContextE(t.Context(), "")

	require.Error(t, err)
}

func TestListServiceBusNamespaceIDsByResourceGroupE(t *testing.T) {
	t.Parallel()

	_, err := azure.ListServiceBusNamespaceIDsByResourceGroupContextE(t.Context(), "", "")

	require.Error(t, err)
}

func TestListNamespaceAuthRulesE(t *testing.T) {
	t.Parallel()

	_, err := azure.ListNamespaceAuthRulesContextE(t.Context(), "", "", "")

	require.Error(t, err)
}

func TestListNamespaceTopicsE(t *testing.T) {
	t.Parallel()

	_, err := azure.ListNamespaceTopicsContextE(t.Context(), "", "", "")

	require.Error(t, err)
}

func TestListTopicAuthRulesE(t *testing.T) {
	t.Parallel()

	_, err := azure.ListTopicAuthRulesContextE(t.Context(), "", "", "", "")

	require.Error(t, err)
}

func TestListTopicSubscriptionsNameE(t *testing.T) {
	t.Parallel()

	_, err := azure.ListTopicSubscriptionsNameContextE(t.Context(), "", "", "", "")

	require.Error(t, err)
}
