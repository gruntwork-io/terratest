package azure

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type azureMLTestHarness struct {
	resourceGroup string
	workspace     string
	compute       string
}

func getTestHarness() *azureMLTestHarness {
	return &azureMLTestHarness{
		resourceGroup: "test-resource-group",
		workspace:     "test-workspace",
		compute:       "test-compute",
	}
}
func TestMachinelearningWorkspaceExistsE(t *testing.T) {
	t.Parallel()

	subscriptionID := ""
	harness := getTestHarness()

	_, err := MachinelearningWorkspaceExistsE(t, harness.resourceGroup, harness.workspace, subscriptionID)
	require.Error(t, err)
}

func TestGetMachinelearningWorkspaceE(t *testing.T) {
	t.Parallel()

	subscriptionID := ""
	harness := getTestHarness()

	_, err := GetMachinelearningWorkspaceE(t, harness.resourceGroup, harness.workspace, subscriptionID)
	require.Error(t, err)
}

func TestMachinelearningComputeExistsE(t *testing.T) {
	t.Parallel()

	subscriptionID := ""
	harness := getTestHarness()

	_, err := MachinelearningComputeExistsE(t, harness.resourceGroup, harness.workspace, harness.compute, subscriptionID)
	require.Error(t, err)
}

func TestGetMachinelearningComputeE(t *testing.T) {
	t.Parallel()

	subscriptionID := ""
	harness := getTestHarness()

	_, err := GetMachinelearningComputeE(t, harness.resourceGroup, harness.workspace, harness.compute, subscriptionID)
	require.Error(t, err)
}
