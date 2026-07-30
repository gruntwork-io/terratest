package k8s_test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/k8s/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
)

// TestGetServiceEndpointBracketsIPv6 pins that endpoints are host/port joined rather than concatenated. An IPv6
// address concatenated with ":port" is ambiguous and unusable in a URL, and this matters more now that node
// resolution prefers the ExternalIP recorded on the node, which is IPv6 on a dual-stack or IPv6 cluster.
func TestGetServiceEndpointBracketsIPv6(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name      string
		clusterIP string
		expected  string
	}{
		{"ipv4 is unchanged", "10.0.0.1", "10.0.0.1:80"},
		{"ipv6 is bracketed", "2600:1f18:abcd::1", "[2600:1f18:abcd::1]:80"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			service := &corev1.Service{
				Spec: corev1.ServiceSpec{Type: corev1.ServiceTypeClusterIP, ClusterIP: testCase.clusterIP},
			}

			endpoint, err := k8s.GetServiceEndpointContextE(t, t.Context(), nil, service, 80)
			require.NoError(t, err)
			assert.Equal(t, testCase.expected, endpoint)
		})
	}
}
