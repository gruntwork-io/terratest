package azure_test

import (
	"context"
	"net/http"
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

	transport := computefake.NewAvailabilitySetsServerTransport(srv)
	client, err := armcompute.NewAvailabilitySetsClient("fake-sub", &azfake.TokenCredential{}, &arm.ClientOptions{
		ClientOptions: policy.ClientOptions{Transport: transport},
	})
	require.NoError(t, err)

	return client
}

func fakeAvsGetHandler(avsName string, vmIDs []string, faultDomainCount int32) func(context.Context, string, string, *armcompute.AvailabilitySetsClientGetOptions) (azfake.Responder[armcompute.AvailabilitySetsClientGetResponse], azfake.ErrorResponder) {
	return func(_ context.Context, _ string, _ string, _ *armcompute.AvailabilitySetsClientGetOptions) (resp azfake.Responder[armcompute.AvailabilitySetsClientGetResponse], errResp azfake.ErrorResponder) {
		vms := make([]*armcompute.SubResource, len(vmIDs))
		for i, id := range vmIDs {
			vms[i] = &armcompute.SubResource{ID: to.Ptr(id)}
		}

		resp.SetResponse(http.StatusOK, armcompute.AvailabilitySetsClientGetResponse{
			AvailabilitySet: armcompute.AvailabilitySet{
				Name: to.Ptr(avsName),
				Properties: &armcompute.AvailabilitySetProperties{
					VirtualMachines:          vms,
					PlatformFaultDomainCount: to.Ptr(faultDomainCount),
				},
			},
		}, nil)

		return
	}
}

func TestGetAvailabilitySetWithClient(t *testing.T) {
	t.Parallel()

	srv := &computefake.AvailabilitySetsServer{
		Get: fakeAvsGetHandler("my-avs", nil, 2),
	}
	client := newFakeAvailabilitySetsClient(t, srv)

	avs, err := azure.GetAvailabilitySetWithClient(t.Context(), client, "rg", "my-avs")
	require.NoError(t, err)
	assert.Equal(t, "my-avs", *avs.Name)
}

func TestCheckAvailabilitySetContainsVMWithClient(t *testing.T) {
	t.Parallel()

	vmIDs := []string{
		"/subscriptions/sub/resourceGroups/RG/providers/Microsoft.Compute/virtualMachines/VM-ONE",
		"/subscriptions/sub/resourceGroups/RG/providers/Microsoft.Compute/virtualMachines/VM-TWO",
	}

	tests := []struct {
		name    string
		vmName  string
		found   bool
		wantErr bool
	}{
		{name: "exact case match", vmName: "VM-ONE", found: true},
		{name: "case insensitive match", vmName: "vm-one", found: true},
		{name: "not found", vmName: "vm-three", found: false, wantErr: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			srv := &computefake.AvailabilitySetsServer{
				Get: fakeAvsGetHandler("avs", vmIDs, 2),
			}
			client := newFakeAvailabilitySetsClient(t, srv)

			found, err := azure.CheckAvailabilitySetContainsVMWithClient(t.Context(), client, "rg", "avs", tc.vmName)

			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tc.found, found)
		})
	}
}

func TestGetAvailabilitySetVMNamesInCapsWithClient(t *testing.T) {
	t.Parallel()

	vmIDs := []string{
		"/subscriptions/sub/resourceGroups/RG/providers/Microsoft.Compute/virtualMachines/VM-ALPHA",
		"/subscriptions/sub/resourceGroups/RG/providers/Microsoft.Compute/virtualMachines/VM-BETA",
	}

	srv := &computefake.AvailabilitySetsServer{
		Get: fakeAvsGetHandler("avs", vmIDs, 3),
	}
	client := newFakeAvailabilitySetsClient(t, srv)

	names, err := azure.GetAvailabilitySetVMNamesInCapsWithClient(t.Context(), client, "rg", "avs")
	require.NoError(t, err)
	assert.Equal(t, []string{"VM-ALPHA", "VM-BETA"}, names)
}

func TestExtractAvailabilitySetFaultDomainCount(t *testing.T) {
	t.Parallel()

	t.Run("valid", func(t *testing.T) {
		t.Parallel()

		avs := &armcompute.AvailabilitySet{
			Properties: &armcompute.AvailabilitySetProperties{
				PlatformFaultDomainCount: to.Ptr[int32](3),
			},
		}

		count, err := azure.ExtractAvailabilitySetFaultDomainCount(avs)
		require.NoError(t, err)
		assert.Equal(t, int32(3), count)
	})

	t.Run("nil properties", func(t *testing.T) {
		t.Parallel()

		avs := &armcompute.AvailabilitySet{}

		_, err := azure.ExtractAvailabilitySetFaultDomainCount(avs)
		require.Error(t, err)
	})
}
