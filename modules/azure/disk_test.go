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

func TestGetDiskE(t *testing.T) {
	t.Parallel()

	diskName := ""
	rgName := ""
	subID := ""

	_, err := azure.GetDiskContextE(t.Context(), diskName, rgName, subID)

	require.Error(t, err)
}
