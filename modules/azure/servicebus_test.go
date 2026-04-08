//go:build azure
// +build azure

// NOTE: We use build tags to differentiate azure testing because we currently do not have azure access setup for
// CircleCI.

package azure

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

/*
The below tests are currently stubbed out, with the expectation that they will throw errors. These tests can be extended.
*/

func TestListServiceBusNamespaceNamesE(t *testing.T) {
	t.Parallel()

	_, err := ListServiceBusNamespaceNamesContextE(context.Background(), "")

	require.Error(t, err)
}

func TestListServiceBusNamespaceIDsByResourceGroupE(t *testing.T) {
	t.Parallel()

	_, err := ListServiceBusNamespaceIDsByResourceGroupContextE(context.Background(), "", "")

	require.Error(t, err)
}

func TestListNamespaceAuthRulesE(t *testing.T) {
	t.Parallel()

	_, err := ListNamespaceAuthRulesContextE(context.Background(), "", "", "")

	require.Error(t, err)
}

func TestListNamespaceTopicsE(t *testing.T) {
	t.Parallel()

	_, err := ListNamespaceTopicsContextE(context.Background(), "", "", "")

	require.Error(t, err)
}

func TestListTopicAuthRulesE(t *testing.T) {
	t.Parallel()

	_, err := ListTopicAuthRulesContextE(context.Background(), "", "", "", "")

	require.Error(t, err)
}

func TestListTopicSubscriptionsNameE(t *testing.T) {
	t.Parallel()

	_, err := ListTopicSubscriptionsNameContextE(context.Background(), "", "", "", "")

	require.Error(t, err)
}
