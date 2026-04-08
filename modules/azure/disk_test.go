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

func TestGetDiskE(t *testing.T) {
	t.Parallel()

	diskName := ""
	rgName := ""
	subID := ""

	_, err := GetDiskContextE(context.Background(), diskName, rgName, subID)

	require.Error(t, err)
}
