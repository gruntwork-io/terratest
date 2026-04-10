package azure_test

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v6"
	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheckIPInCIDR(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		ip      string
		cidr    string
		want    bool
		wantErr bool
	}{
		{"IPInRange", "10.0.0.5", "10.0.0.0/24", true, false},
		{"IPOutOfRange", "10.0.1.5", "10.0.0.0/24", false, false},
		{"IPAtNetworkBoundary", "10.0.0.0", "10.0.0.0/24", true, false},
		{"IPAtBroadcast", "10.0.0.255", "10.0.0.0/24", true, false},
		{"InvalidIP", "not-an-ip", "10.0.0.0/24", false, true},
		{"EmptyIP", "", "10.0.0.0/24", false, true},
		{"InvalidCIDR", "10.0.0.5", "not-a-cidr", false, true},
		{"IPv6InRange", "fd00::5", "fd00::/64", true, false},
		{"IPv6OutOfRange", "fd01::5", "fd00::/64", false, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := azure.CheckIPInCIDR(tc.ip, tc.cidr)
			if tc.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestExtractDNSServerIPs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		vnet *armnetwork.VirtualNetwork
		want []string
	}{
		{
			name: "TwoDNSServers",
			vnet: &armnetwork.VirtualNetwork{
				Properties: &armnetwork.VirtualNetworkPropertiesFormat{
					DhcpOptions: &armnetwork.DhcpOptions{
						DNSServers: []*string{to.Ptr("8.8.8.8"), to.Ptr("8.8.4.4")},
					},
				},
			},
			want: []string{"8.8.8.8", "8.8.4.4"},
		},
		{
			name: "NilDhcpOptions",
			vnet: &armnetwork.VirtualNetwork{
				Properties: &armnetwork.VirtualNetworkPropertiesFormat{
					DhcpOptions: nil,
				},
			},
			want: nil,
		},
		{
			name: "EmptyDNSServers",
			vnet: &armnetwork.VirtualNetwork{
				Properties: &armnetwork.VirtualNetworkPropertiesFormat{
					DhcpOptions: &armnetwork.DhcpOptions{
						DNSServers: []*string{},
					},
				},
			},
			want: []string{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := azure.ExtractDNSServerIPs(tc.vnet)
			assert.Equal(t, tc.want, got)
		})
	}
}
