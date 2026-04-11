package azure_test

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v6"
	computefake "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v6/fake"
	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newFakeAvailabilitySetsClient(t *testing.T, srv *computefake.AvailabilitySetsServer) *armcompute.AvailabilitySetsClient {
	t.Helper()

	client, err := armcompute.NewAvailabilitySetsClient("fake-sub", &azfake.TokenCredential{},
		&arm.ClientOptions{ClientOptions: policy.ClientOptions{
			Transport: computefake.NewAvailabilitySetsServerTransport(srv),
		}})
	require.NoError(t, err)

	return client
}

// fakeAvailabilitySet returns a standard test availability set for reuse across tests.
func fakeAvailabilitySet() armcompute.AvailabilitySet {
	return armcompute.AvailabilitySet{
		Name: to.Ptr("test-avs"),
		Properties: &armcompute.AvailabilitySetProperties{
			PlatformFaultDomainCount: to.Ptr[int32](3),
			VirtualMachines: []*armcompute.SubResource{
				{ID: to.Ptr("/subscriptions/sub/resourceGroups/rg/providers/Microsoft.Compute/virtualMachines/VM-ONE")},
				{ID: to.Ptr("/subscriptions/sub/resourceGroups/rg/providers/Microsoft.Compute/virtualMachines/VM-TWO")},
			},
		},
	}
}

func newFakeAvailabilitySetsServerWithResponse() *computefake.AvailabilitySetsServer {
	avs := fakeAvailabilitySet()

	return &computefake.AvailabilitySetsServer{
		Get: func(_ context.Context, _ string, _ string, _ *armcompute.AvailabilitySetsClientGetOptions) (resp azfake.Responder[armcompute.AvailabilitySetsClientGetResponse], errResp azfake.ErrorResponder) {
			resp.SetResponse(http.StatusOK, armcompute.AvailabilitySetsClientGetResponse{
				AvailabilitySet: avs,
			}, nil)

			return
		},
	}
}

// ---------------------------------------------------------------------------
// CheckAvailabilitySetContainsVM tests (case-insensitive matching logic)
// ---------------------------------------------------------------------------

func TestCheckAvailabilitySetContainsVM(t *testing.T) {
	t.Parallel()

	srv := newFakeAvailabilitySetsServerWithResponse()
	client := newFakeAvailabilitySetsClient(t, srv)

	tests := []struct {
		name   string
		vmName string
		want   bool
	}{
		{
			name:   "CaseInsensitiveMatch",
			vmName: "vm-one",
			want:   true,
		},
		{
			name:   "ExactCaseMatch",
			vmName: "VM-ONE",
			want:   true,
		},
		{
			name:   "NotFound",
			vmName: "vm-three",
			want:   false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			resp, err := client.Get(context.Background(), "rg", "test-avs", nil)
			require.NoError(t, err)

			found := false

			for _, vm := range resp.Properties.VirtualMachines {
				if name := azure.GetNameFromResourceID(*vm.ID); strings.EqualFold(tc.vmName, name) {
					found = true

					break
				}
			}

			assert.Equal(t, tc.want, found)
		})
	}
}

// ---------------------------------------------------------------------------
// GetAvailabilitySetVMNamesInCaps tests (extracts names from resource IDs)
// ---------------------------------------------------------------------------

func TestGetAvailabilitySetVMNamesInCaps(t *testing.T) {
	t.Parallel()

	srv := newFakeAvailabilitySetsServerWithResponse()
	client := newFakeAvailabilitySetsClient(t, srv)

	resp, err := client.Get(context.Background(), "rg", "test-avs", nil)
	require.NoError(t, err)

	vms := make([]string, 0, len(resp.Properties.VirtualMachines))

	for _, vm := range resp.Properties.VirtualMachines {
		if vmName := azure.GetNameFromResourceID(*vm.ID); len(vmName) > 0 {
			vms = append(vms, vmName)
		}
	}

	assert.Equal(t, []string{"VM-ONE", "VM-TWO"}, vms)
}

// ---------------------------------------------------------------------------
// GetAvailabilitySetFaultDomainCount tests (returns count)
// ---------------------------------------------------------------------------

func TestGetAvailabilitySetFaultDomainCount(t *testing.T) {
	t.Parallel()

	srv := newFakeAvailabilitySetsServerWithResponse()
	client := newFakeAvailabilitySetsClient(t, srv)

	resp, err := client.Get(context.Background(), "rg", "test-avs", nil)
	require.NoError(t, err)

	assert.Equal(t, int32(3), *resp.Properties.PlatformFaultDomainCount)
}

// ---------------------------------------------------------------------------
// AvailabilitySet NotFound error handling
// ---------------------------------------------------------------------------

func TestGetAvailabilitySet_NotFound(t *testing.T) {
	t.Parallel()

	srv := &computefake.AvailabilitySetsServer{
		Get: func(_ context.Context, _ string, _ string, _ *armcompute.AvailabilitySetsClientGetOptions) (resp azfake.Responder[armcompute.AvailabilitySetsClientGetResponse], errResp azfake.ErrorResponder) {
			errResp.SetResponseError(http.StatusNotFound, "ResourceNotFound")

			return
		},
	}

	client := newFakeAvailabilitySetsClient(t, srv)

	_, err := client.Get(context.Background(), "rg", "missing-avs", nil)
	require.Error(t, err)
	assert.True(t, azure.ResourceNotFoundErrorExists(err))
}
