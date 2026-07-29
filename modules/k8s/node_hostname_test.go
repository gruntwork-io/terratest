package k8s_test

import (
	"context"
	"errors"
	"testing"

	gotesting "github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/gruntwork-io/terratest/modules/k8s/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
)

const (
	testInstanceID  = "i-0123456789abcdef0"
	testProviderID  = "aws:///us-east-1a/" + testInstanceID
	testHostname    = "ip-10-0-0-1.ec2.internal"
	testPublicIP    = "203.0.113.10"
	testExpectedReg = "us-east-1"
)

func nodeWithProviderID(providerID string) corev1.Node {
	return corev1.Node{
		Spec: corev1.NodeSpec{ProviderID: providerID},
		Status: corev1.NodeStatus{
			Addresses: []corev1.NodeAddress{
				{Type: corev1.NodeHostName, Address: testHostname},
			},
		},
	}
}

// nodeWithExternalIP models what a cloud controller manager records for an instance that has a public IP.
func nodeWithExternalIP(providerID string, externalIP string) corev1.Node {
	node := nodeWithProviderID(providerID)
	node.Status.Addresses = append(
		[]corev1.NodeAddress{{Type: corev1.NodeExternalIP, Address: externalIP}},
		node.Status.Addresses...,
	)

	return node
}

func TestFindNodeHostnameUsesPublicIPLookup(t *testing.T) {
	t.Parallel()

	var gotInstanceIDs []string

	var gotRegion string

	options := k8s.NewKubectlOptions("", "", "default")
	options.NodePublicIPLookup = func(_ gotesting.TestingT, _ context.Context, instanceIDs []string, region string) (map[string]string, error) {
		gotInstanceIDs, gotRegion = instanceIDs, region

		return map[string]string{testInstanceID: testPublicIP}, nil
	}

	hostname, err := k8s.FindNodeHostnameWithOptionsContextE(t, t.Context(), options, nodeWithProviderID(testProviderID))
	require.NoError(t, err)

	assert.Equal(t, testPublicIP, hostname)
	assert.Equal(t, []string{testInstanceID}, gotInstanceIDs)
	assert.Equal(t, testExpectedReg, gotRegion)
}

func TestFindNodeHostnameFallsBackWhenNoLookupConfigured(t *testing.T) {
	t.Parallel()

	options := k8s.NewKubectlOptions("", "", "default")

	hostname, err := k8s.FindNodeHostnameWithOptionsContextE(t, t.Context(), options, nodeWithProviderID(testProviderID))
	require.NoError(t, err)

	assert.Equal(t, testHostname, hostname, "should fall back to the internal hostname when no lookup is configured")
}

func TestFindNodeHostnameFallsBackWhenInstanceHasNoPublicIP(t *testing.T) {
	t.Parallel()

	options := k8s.NewKubectlOptions("", "", "default")
	options.NodePublicIPLookup = func(_ gotesting.TestingT, _ context.Context, _ []string, _ string) (map[string]string, error) {
		return map[string]string{}, nil
	}

	hostname, err := k8s.FindNodeHostnameWithOptionsContextE(t, t.Context(), options, nodeWithProviderID(testProviderID))
	require.NoError(t, err)

	assert.Equal(t, testHostname, hostname)
}

func TestFindNodeHostnameSkipsLookupForNonAWSProvider(t *testing.T) {
	t.Parallel()

	options := k8s.NewKubectlOptions("", "", "default")
	options.NodePublicIPLookup = func(_ gotesting.TestingT, _ context.Context, _ []string, _ string) (map[string]string, error) {
		t.Error("lookup should not be called for a non-AWS provider ID")

		return nil, nil
	}

	hostname, err := k8s.FindNodeHostnameWithOptionsContextE(t, t.Context(), options, nodeWithProviderID("gce://project/us-central1-a/instance-1"))
	require.NoError(t, err)

	assert.Equal(t, testHostname, hostname)
}

func TestFindNodeHostnamePropagatesLookupError(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("describe instances failed")

	options := k8s.NewKubectlOptions("", "", "default")
	options.NodePublicIPLookup = func(_ gotesting.TestingT, _ context.Context, _ []string, _ string) (map[string]string, error) {
		return nil, expectedErr
	}

	_, err := k8s.FindNodeHostnameWithOptionsContextE(t, t.Context(), options, nodeWithProviderID(testProviderID))
	require.ErrorIs(t, err, expectedErr)
}

// TestFindNodeHostnameRejectsMalformedAwsProviderIDs covers provider IDs that satisfy the path-segment count but
// carry an empty availability zone or instance ID. Before the explicit guard, an empty AZ reached
// availabilityZone[:len(availabilityZone)-1] and panicked with a slice-bounds error.
func TestFindNodeHostnameRejectsMalformedAwsProviderIDs(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name       string
		providerID string
	}{
		{"empty availability zone", "aws:////i-0123456789"},
		{"empty instance ID", "aws:///us-east-1d/"},
		{"both empty", "aws:////"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			node := corev1.Node{
				Spec: corev1.NodeSpec{ProviderID: testCase.providerID},
				Status: corev1.NodeStatus{
					Addresses: []corev1.NodeAddress{{Type: corev1.NodeHostName, Address: testHostname}},
				},
			}

			options := k8s.NewKubectlOptions("", "", "default")
			options.NodePublicIPLookup = func(_ gotesting.TestingT, _ context.Context, _ []string, _ string) (map[string]string, error) {
				t.Fatal("lookup must not be reached for a malformed provider ID")

				return nil, nil
			}

			// The guard must fire before the region slice, with and without a lookup configured.
			for _, opts := range []*k8s.KubectlOptions{nil, options} {
				require.NotPanics(t, func() {
					_, err := k8s.FindNodeHostnameWithOptionsContextE(t, t.Context(), opts, node)
					require.Error(t, err)
				})
			}
		})
	}
}

// TestFindNodeHostnamePrefersExternalIPFromNode is the primary path. Cloud controller managers record an
// instance's public IP as an ExternalIP address on the node object, so no cloud API call is needed and no lookup
// has to be configured.
func TestFindNodeHostnamePrefersExternalIPFromNode(t *testing.T) {
	t.Parallel()

	node := nodeWithExternalIP(testProviderID, testPublicIP)

	hostname, err := k8s.FindNodeHostnameContextE(t, t.Context(), node)
	require.NoError(t, err)
	assert.Equal(t, testPublicIP, hostname, "the node's ExternalIP must win without any options")
}

// TestFindNodeHostnameExternalIPSkipsTheLookup pins that the cloud API call is not made when the node already
// advertises an external IP, even if a lookup is configured.
func TestFindNodeHostnameExternalIPSkipsTheLookup(t *testing.T) {
	t.Parallel()

	options := k8s.NewKubectlOptions("", "", "default")
	options.NodePublicIPLookup = func(_ gotesting.TestingT, _ context.Context, _ []string, _ string) (map[string]string, error) {
		t.Fatal("lookup must not be called when the node advertises an ExternalIP")

		return nil, nil
	}

	hostname, err := k8s.FindNodeHostnameWithOptionsContextE(t, t.Context(), options, nodeWithExternalIP(testProviderID, testPublicIP))
	require.NoError(t, err)
	assert.Equal(t, testPublicIP, hostname)
}

// TestFindNodeHostnameExternalIPWorksForNonAwsProviders confirms the preference is provider agnostic, so GKE and
// other clusters that advertise an ExternalIP get it too.
func TestFindNodeHostnameExternalIPWorksForNonAwsProviders(t *testing.T) {
	t.Parallel()

	hostname, err := k8s.FindNodeHostnameContextE(t, t.Context(), nodeWithExternalIP("gce://project/zone/instance", testPublicIP))
	require.NoError(t, err)
	assert.Equal(t, testPublicIP, hostname)
}

// TestFindNodeHostnameWithoutOptionsFallsBackToHostname covers an AWS node with no ExternalIP reached through the
// options-free entry point: there is nothing to query with, so it degrades to the internal hostname.
func TestFindNodeHostnameWithoutOptionsFallsBackToHostname(t *testing.T) {
	t.Parallel()

	hostname, err := k8s.FindNodeHostnameContextE(t, t.Context(), nodeWithProviderID(testProviderID))
	require.NoError(t, err)
	assert.Equal(t, testHostname, hostname)
}

// TestFindNodeHostnameIgnoresEmptyExternalIP guards against an ExternalIP entry with a blank address shadowing the
// rest of the resolution chain.
func TestFindNodeHostnameIgnoresEmptyExternalIP(t *testing.T) {
	t.Parallel()

	hostname, err := k8s.FindNodeHostnameContextE(t, t.Context(), nodeWithExternalIP(testProviderID, ""))
	require.NoError(t, err)
	assert.Equal(t, testHostname, hostname, "a blank ExternalIP must not shadow the hostname")
}
