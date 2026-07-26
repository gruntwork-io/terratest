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

func TestFindNodeHostnameUsesPublicIPLookup(t *testing.T) {
	t.Parallel()

	var gotInstanceIDs []string

	var gotRegion string

	options := k8s.NewKubectlOptions("", "", "default")
	options.NodePublicIPLookup = func(_ gotesting.TestingT, _ context.Context, instanceIDs []string, region string) (map[string]string, error) {
		gotInstanceIDs, gotRegion = instanceIDs, region

		return map[string]string{testInstanceID: testPublicIP}, nil
	}

	hostname, err := k8s.FindNodeHostnameContextE(t, t.Context(), options, nodeWithProviderID(testProviderID))
	require.NoError(t, err)

	assert.Equal(t, testPublicIP, hostname)
	assert.Equal(t, []string{testInstanceID}, gotInstanceIDs)
	assert.Equal(t, testExpectedReg, gotRegion)
}

func TestFindNodeHostnameFallsBackWhenNoLookupConfigured(t *testing.T) {
	t.Parallel()

	options := k8s.NewKubectlOptions("", "", "default")

	hostname, err := k8s.FindNodeHostnameContextE(t, t.Context(), options, nodeWithProviderID(testProviderID))
	require.NoError(t, err)

	assert.Equal(t, testHostname, hostname, "should fall back to the internal hostname when no lookup is configured")
}

func TestFindNodeHostnameFallsBackWhenInstanceHasNoPublicIP(t *testing.T) {
	t.Parallel()

	options := k8s.NewKubectlOptions("", "", "default")
	options.NodePublicIPLookup = func(_ gotesting.TestingT, _ context.Context, _ []string, _ string) (map[string]string, error) {
		return map[string]string{}, nil
	}

	hostname, err := k8s.FindNodeHostnameContextE(t, t.Context(), options, nodeWithProviderID(testProviderID))
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

	hostname, err := k8s.FindNodeHostnameContextE(t, t.Context(), options, nodeWithProviderID("gce://project/us-central1-a/instance-1"))
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

	_, err := k8s.FindNodeHostnameContextE(t, t.Context(), options, nodeWithProviderID(testProviderID))
	require.ErrorIs(t, err, expectedErr)
}
