package azure

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v6"
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
		{
			name: "IPInRange",
			ip:   "10.0.0.5",
			cidr: "10.0.0.0/24",
			want: true,
		},
		{
			name: "IPOutOfRange",
			ip:   "10.0.1.5",
			cidr: "10.0.0.0/24",
			want: false,
		},
		{
			name: "IPAtNetworkBoundary",
			ip:   "10.0.0.0",
			cidr: "10.0.0.0/24",
			want: true,
		},
		{
			name: "IPAtBroadcast",
			ip:   "10.0.0.255",
			cidr: "10.0.0.0/24",
			want: true,
		},
		{
			name:    "InvalidIP",
			ip:      "not-an-ip",
			cidr:    "10.0.0.0/24",
			wantErr: true,
		},
		{
			name:    "EmptyIP",
			ip:      "",
			cidr:    "10.0.0.0/24",
			wantErr: true,
		},
		{
			name:    "InvalidCIDR",
			ip:      "10.0.0.5",
			cidr:    "not-a-cidr",
			wantErr: true,
		},
		{
			name: "IPv6InRange",
			ip:   "fd00::5",
			cidr: "fd00::/64",
			want: true,
		},
		{
			name: "IPv6OutOfRange",
			ip:   "fd01::5",
			cidr: "fd00::/64",
			want: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := checkIPInCIDR(tc.ip, tc.cidr)
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

			got := extractDNSServerIPs(tc.vnet)
			assert.Equal(t, tc.want, got)
		})
	}
}
