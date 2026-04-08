package azure //nolint:testpackage // tests access unexported functions

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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newFakeVMClient creates a fake VirtualMachinesClient backed by the given fake server.
func newFakeVMClient(t *testing.T, srv *computefake.VirtualMachinesServer) *armcompute.VirtualMachinesClient {
	t.Helper()

	client, err := armcompute.NewVirtualMachinesClient("fake-sub", &azfake.TokenCredential{},
		&arm.ClientOptions{ClientOptions: policy.ClientOptions{
			Transport: computefake.NewVirtualMachinesServerTransport(srv),
		}})
	require.NoError(t, err)

	return client
}

// ---------------------------------------------------------------------------
// fetchVirtualMachine tests
// ---------------------------------------------------------------------------

func TestFetchVirtualMachine(t *testing.T) {
	t.Parallel()

	tests := []struct { //nolint:govet // fieldalignment not worth optimizing in test structs
		name      string
		wantName  string
		errSubstr string
		server    computefake.VirtualMachinesServer
		wantErr   bool
	}{
		{
			name: "Success",
			server: computefake.VirtualMachinesServer{
				Get: func(_ context.Context, _ string, _ string, _ *armcompute.VirtualMachinesClientGetOptions) (resp azfake.Responder[armcompute.VirtualMachinesClientGetResponse], errResp azfake.ErrorResponder) {
					result := armcompute.VirtualMachinesClientGetResponse{
						VirtualMachine: armcompute.VirtualMachine{
							Name: to.Ptr("test-vm"),
							Properties: &armcompute.VirtualMachineProperties{
								HardwareProfile: &armcompute.HardwareProfile{
									VMSize: to.Ptr(armcompute.VirtualMachineSizeTypesStandardDS1V2),
								},
							},
						},
					}
					resp.SetResponse(http.StatusOK, result, nil)

					return
				},
			},
			wantName: "test-vm",
		},
		{
			name: "NotFound",
			server: computefake.VirtualMachinesServer{
				Get: func(_ context.Context, _ string, _ string, _ *armcompute.VirtualMachinesClientGetOptions) (resp azfake.Responder[armcompute.VirtualMachinesClientGetResponse], errResp azfake.ErrorResponder) {
					errResp.SetResponseError(http.StatusNotFound, "ResourceNotFound")

					return
				},
			},
			wantErr:   true,
			errSubstr: "ResourceNotFound",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			client := newFakeVMClient(t, &tc.server)

			vm, err := fetchVirtualMachine(context.Background(), client, "rg", "vm")
			if tc.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errSubstr)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.wantName, *vm.Name)
		})
	}
}

// ---------------------------------------------------------------------------
// extractVMNics tests
// ---------------------------------------------------------------------------

func TestExtractVMNics(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		vm      *armcompute.VirtualMachine
		want    []string
		wantErr bool
	}{
		{
			name: "TwoValidNICs",
			vm: &armcompute.VirtualMachine{
				Properties: &armcompute.VirtualMachineProperties{
					NetworkProfile: &armcompute.NetworkProfile{
						NetworkInterfaces: []*armcompute.NetworkInterfaceReference{
							{ID: to.Ptr("/subscriptions/sub/resourceGroups/rg/providers/Microsoft.Network/networkInterfaces/nic1")},
							{ID: to.Ptr("/subscriptions/sub/resourceGroups/rg/providers/Microsoft.Network/networkInterfaces/nic2")},
						},
					},
				},
			},
			want: []string{"nic1", "nic2"},
		},
		{
			name: "NilNetworkProfile",
			vm: &armcompute.VirtualMachine{
				Properties: &armcompute.VirtualMachineProperties{
					NetworkProfile: nil,
				},
			},
			want: nil,
		},
		{
			name: "InvalidNICResourceID",
			vm: &armcompute.VirtualMachine{
				Properties: &armcompute.VirtualMachineProperties{
					NetworkProfile: &armcompute.NetworkProfile{
						NetworkInterfaces: []*armcompute.NetworkInterfaceReference{
							{ID: to.Ptr("")},
						},
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := extractVMNics(tc.vm)
			if tc.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

// ---------------------------------------------------------------------------
// extractVMAvailabilitySetID tests
// ---------------------------------------------------------------------------

func TestExtractVMAvailabilitySetID(t *testing.T) {
	t.Parallel()

	tests := []struct { //nolint:govet // fieldalignment not worth optimizing in test structs
		name string
		want string
		vm   *armcompute.VirtualMachine
	}{
		{
			name: "AvailabilitySetPresent",
			vm: &armcompute.VirtualMachine{
				Properties: &armcompute.VirtualMachineProperties{
					AvailabilitySet: &armcompute.SubResource{
						ID: to.Ptr("/subscriptions/sub/resourceGroups/rg/providers/Microsoft.Compute/availabilitySets/myAvSet"),
					},
				},
			},
			want: "myAvSet",
		},
		{
			name: "NilAvailabilitySet",
			vm: &armcompute.VirtualMachine{
				Properties: &armcompute.VirtualMachineProperties{
					AvailabilitySet: nil,
				},
			},
			want: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := extractVMAvailabilitySetID(tc.vm)
			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

// ---------------------------------------------------------------------------
// extractVMImage tests
// ---------------------------------------------------------------------------

func TestExtractVMImage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		vm   *armcompute.VirtualMachine
		want VMImage
	}{
		{
			name: "MarketplaceImage",
			vm: &armcompute.VirtualMachine{
				Properties: &armcompute.VirtualMachineProperties{
					StorageProfile: &armcompute.StorageProfile{
						ImageReference: &armcompute.ImageReference{
							Publisher: to.Ptr("Canonical"),
							Offer:     to.Ptr("UbuntuServer"),
							SKU:       to.Ptr("18.04-LTS"),
							Version:   to.Ptr("latest"),
						},
					},
				},
			},
			want: VMImage{
				Publisher: "Canonical",
				Offer:     "UbuntuServer",
				SKU:       "18.04-LTS",
				Version:   "latest",
			},
		},
		{
			name: "CustomImage",
			vm: &armcompute.VirtualMachine{
				Properties: &armcompute.VirtualMachineProperties{
					StorageProfile: &armcompute.StorageProfile{
						ImageReference: &armcompute.ImageReference{
							ID: to.Ptr("/subscriptions/sub/resourceGroups/rg/providers/Microsoft.Compute/images/myImage"),
						},
					},
				},
			},
			want: VMImage{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := extractVMImage(tc.vm)
			assert.Equal(t, tc.want, got)
		})
	}
}

// ---------------------------------------------------------------------------
// extractVMTags tests
// ---------------------------------------------------------------------------

func TestExtractVMTags(t *testing.T) {
	t.Parallel()

	tests := []struct { //nolint:govet // fieldalignment not worth optimizing in test structs
		name string
		vm   *armcompute.VirtualMachine
		want map[string]string
	}{
		{
			name: "TagsPresent",
			vm: &armcompute.VirtualMachine{
				Tags: map[string]*string{
					"env":   to.Ptr("dev"),
					"owner": to.Ptr("team-a"),
				},
			},
			want: map[string]string{"env": "dev", "owner": "team-a"},
		},
		{
			name: "NilTags",
			vm:   &armcompute.VirtualMachine{},
			want: map[string]string{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := extractVMTags(tc.vm)
			assert.Equal(t, tc.want, got)
		})
	}
}

// ---------------------------------------------------------------------------
// listVirtualMachineNames tests
// ---------------------------------------------------------------------------

func TestListVirtualMachineNames(t *testing.T) {
	t.Parallel()

	srv := &computefake.VirtualMachinesServer{
		NewListPager: func(_ string, _ *armcompute.VirtualMachinesClientListOptions) (resp azfake.PagerResponder[armcompute.VirtualMachinesClientListResponse]) {
			resp.AddPage(http.StatusOK, armcompute.VirtualMachinesClientListResponse{
				VirtualMachineListResult: armcompute.VirtualMachineListResult{
					Value: []*armcompute.VirtualMachine{
						{Name: to.Ptr("vm1")},
						{Name: to.Ptr("vm2")},
					},
				},
			}, nil)

			return
		},
	}

	client := newFakeVMClient(t, srv)
	names, err := listVirtualMachineNames(context.Background(), client, "rg")

	require.NoError(t, err)
	assert.Equal(t, []string{"vm1", "vm2"}, names)
}
