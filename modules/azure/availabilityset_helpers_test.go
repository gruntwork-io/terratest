package azure

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAvsContainsVM(t *testing.T) {
	t.Parallel()

	vms := []*armcompute.SubResource{
		{ID: to.Ptr("/subscriptions/sub/resourceGroups/rg/providers/Microsoft.Compute/virtualMachines/VM-ONE")},
		{ID: to.Ptr("/subscriptions/sub/resourceGroups/rg/providers/Microsoft.Compute/virtualMachines/VM-TWO")},
	}

	tests := []struct { //nolint:govet // fieldalignment not worth optimizing in test structs
		name    string
		vmName  string
		want    bool
		wantErr bool
	}{
		{
			name:   "ExactCaseMatch",
			vmName: "VM-ONE",
			want:   true,
		},
		{
			name:   "CaseInsensitiveMatch",
			vmName: "vm-one",
			want:   true,
		},
		{
			name:   "MixedCaseMatch",
			vmName: "Vm-Two",
			want:   true,
		},
		{
			name:    "NotFound",
			vmName:  "vm-three",
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := avsContainsVM(tc.vmName, "test-avs", vms)
			if tc.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestAvsContainsVMEmptyList(t *testing.T) {
	t.Parallel()

	_, err := avsContainsVM("any-vm", "test-avs", []*armcompute.SubResource{})
	require.Error(t, err)
}

func TestExtractAvsVMNames(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		vms  []*armcompute.SubResource
		want []string
	}{
		{
			name: "MultipleVMs",
			vms: []*armcompute.SubResource{
				{ID: to.Ptr("/subscriptions/sub/resourceGroups/rg/providers/Microsoft.Compute/virtualMachines/VM-ONE")},
				{ID: to.Ptr("/subscriptions/sub/resourceGroups/rg/providers/Microsoft.Compute/virtualMachines/VM-TWO")},
			},
			want: []string{"VM-ONE", "VM-TWO"},
		},
		{
			name: "EmptyList",
			vms:  []*armcompute.SubResource{},
			want: nil,
		},
		{
			name: "SkipsInvalidIDs",
			vms: []*armcompute.SubResource{
				{ID: to.Ptr("/subscriptions/sub/resourceGroups/rg/providers/Microsoft.Compute/virtualMachines/VALID-VM")},
				{ID: to.Ptr("")},
			},
			want: []string{"VALID-VM"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := extractAvsVMNames(tc.vms)
			assert.Equal(t, tc.want, got)
		})
	}
}
