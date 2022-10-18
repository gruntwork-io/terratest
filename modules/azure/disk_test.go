//go:build azure
// +build azure

// NOTE: We use build tags to differentiate azure testing because we currently do not have azure access setup for
// CircleCI.

package azure

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetDiskE(t *testing.T) {
	t.Parallel()

	diskName := ""
	rgName := ""
	subID := ""

	_, err := GetDiskE(diskName, rgName, subID)

	require.Error(t, err)
}
