package azure_test

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v6"
	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAvsContainsVM(t *testing.T) {
	t.Parallel()

	vms := []*armcompute.SubResource{
		{ID: to.Ptr("/subscriptions/sub/resourceGroups/rg/providers/Microsoft.Compute/virtualMachines/VM-ONE")},
		{ID: to.Ptr("/subscriptions/sub/resourceGroups/rg/providers/Microsoft.Compute/virtualMachines/VM-TWO")},
	}

	tests := []struct {
		name    string
		vmName  string
		want    bool
		wantErr bool
	}{
		{"ExactCaseMatch", "VM-ONE", true, false},
		{"CaseInsensitiveMatch", "vm-one", true, false},
		{"MixedCaseMatch", "Vm-Two", true, false},
		{"NotFound", "vm-three", false, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := azure.AvsContainsVM(tc.vmName, "test-avs", vms)
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

	_, err := azure.AvsContainsVM("any-vm", "test-avs", []*armcompute.SubResource{})
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

			got := azure.ExtractAvsVMNames(tc.vms)
			assert.Equal(t, tc.want, got)
		})
	}
}
