package k8s

import (
	"context"
	"fmt"
	"time"

	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/gruntwork-io/terratest/modules/retry"
	"github.com/gruntwork-io/terratest/modules/testing"
)

// ListDeploymentsContextE looks up deployments in the given namespace that match the given filters and return them.
// The ctx parameter supports cancellation and timeouts.
//
//nolint:gocritic // hugeParam: cannot change public function signature
func ListDeploymentsContextE(t testing.TestingT, ctx context.Context, options *KubectlOptions, filters metav1.ListOptions) ([]appsv1.Deployment, error) {
	clientset, err := GetKubernetesClientFromOptionsContextE(t, ctx, options)
	if err != nil {
		return nil, err
	}

	deployments, err := clientset.AppsV1().Deployments(options.Namespace).List(ctx, filters)
	if err != nil {
		return nil, err
	}

	return deployments.Items, nil
}

// ListDeploymentsContext looks up deployments in the given namespace that match the given filters and return them.
// The ctx parameter supports cancellation and timeouts.
// This will fail the test if there is an error.
//
//nolint:gocritic // hugeParam: cannot change public function signature
func ListDeploymentsContext(t testing.TestingT, ctx context.Context, options *KubectlOptions, filters metav1.ListOptions) []appsv1.Deployment {
	t.Helper()
	deployment, err := ListDeploymentsContextE(t, ctx, options, filters)
	require.NoError(t, err)

	return deployment
}

// ListDeployments will look for deployments in the given namespace that match the given filters and return them. This will
// fail the test if there is an error.
//
// Deprecated: Use [ListDeploymentsContext] instead.
//
//nolint:gocritic // hugeParam: cannot change public function signature
func ListDeployments(t testing.TestingT, options *KubectlOptions, filters metav1.ListOptions) []appsv1.Deployment {
	t.Helper()

	return ListDeploymentsContext(t, context.Background(), options, filters)
}

// ListDeploymentsE will look for deployments in the given namespace that match the given filters and return them.
//
// Deprecated: Use [ListDeploymentsContextE] instead.
//
//nolint:gocritic // hugeParam: cannot change public function signature
func ListDeploymentsE(t testing.TestingT, options *KubectlOptions, filters metav1.ListOptions) ([]appsv1.Deployment, error) {
	return ListDeploymentsContextE(t, context.Background(), options, filters)
}

// GetDeploymentContextE returns a Kubernetes deployment resource in the provided namespace with the given name.
// The ctx parameter supports cancellation and timeouts.
func GetDeploymentContextE(t testing.TestingT, ctx context.Context, options *KubectlOptions, deploymentName string) (*appsv1.Deployment, error) {
	clientset, err := GetKubernetesClientFromOptionsContextE(t, ctx, options)
	if err != nil {
		return nil, err
	}

	return clientset.AppsV1().Deployments(options.Namespace).Get(ctx, deploymentName, metav1.GetOptions{})
}

// GetDeploymentContext returns a Kubernetes deployment resource in the provided namespace with the given name.
// The ctx parameter supports cancellation and timeouts.
// This will fail the test if there is an error.
func GetDeploymentContext(t testing.TestingT, ctx context.Context, options *KubectlOptions, deploymentName string) *appsv1.Deployment {
	t.Helper()
	deployment, err := GetDeploymentContextE(t, ctx, options, deploymentName)
	require.NoError(t, err)

	return deployment
}

// GetDeployment returns a Kubernetes deployment resource in the provided namespace with the given name. This will
// fail the test if there is an error.
//
// Deprecated: Use [GetDeploymentContext] instead.
func GetDeployment(t testing.TestingT, options *KubectlOptions, deploymentName string) *appsv1.Deployment {
	t.Helper()

	return GetDeploymentContext(t, context.Background(), options, deploymentName)
}

// GetDeploymentE returns a Kubernetes deployment resource in the provided namespace with the given name.
//
// Deprecated: Use [GetDeploymentContextE] instead.
func GetDeploymentE(t testing.TestingT, options *KubectlOptions, deploymentName string) (*appsv1.Deployment, error) {
	return GetDeploymentContextE(t, context.Background(), options, deploymentName)
}

// WaitUntilDeploymentAvailableContextE waits until all pods within the deployment are ready and started,
// retrying the check for the specified amount of times, sleeping for the provided duration between each try.
// The ctx parameter supports cancellation and timeouts.
func WaitUntilDeploymentAvailableContextE( //nolint:dupl // similar retry pattern across resource types is intentional
	t testing.TestingT,
	ctx context.Context,
	options *KubectlOptions,
	deploymentName string,
	retries int,
	sleepBetweenRetries time.Duration,
) error {
	statusMsg := fmt.Sprintf("Wait for deployment %s to be provisioned.", deploymentName)

	message, err := retry.DoWithRetryContextE(
		t,
		ctx,
		statusMsg,
		retries,
		sleepBetweenRetries,
		func() (string, error) {
			deployment, err := GetDeploymentContextE(t, ctx, options, deploymentName)
			if err != nil {
				return "", err
			}

			if !IsDeploymentAvailable(deployment) {
				return "", NewDeploymentNotAvailableError(deployment)
			}

			return "Deployment is now available", nil
		},
	)
	if err != nil {
		options.Logger.Logf(t, "Timedout waiting for Deployment to be provisioned: %s", err)
		return err
	}

	options.Logger.Logf(t, "%s", message)

	return nil
}

// WaitUntilDeploymentAvailableContext waits until all pods within the deployment are ready and started,
// retrying the check for the specified amount of times, sleeping for the provided duration between each try.
// The ctx parameter supports cancellation and timeouts.
// This will fail the test if there is an error.
func WaitUntilDeploymentAvailableContext(t testing.TestingT, ctx context.Context, options *KubectlOptions, deploymentName string, retries int, sleepBetweenRetries time.Duration) {
	t.Helper()
	require.NoError(t, WaitUntilDeploymentAvailableContextE(t, ctx, options, deploymentName, retries, sleepBetweenRetries))
}

// WaitUntilDeploymentAvailable waits until all pods within the deployment are ready and started,
// retrying the check for the specified amount of times, sleeping
// for the provided duration between each try.
// This will fail the test if there is an error.
//
// Deprecated: Use [WaitUntilDeploymentAvailableContext] instead.
func WaitUntilDeploymentAvailable(t testing.TestingT, options *KubectlOptions, deploymentName string, retries int, sleepBetweenRetries time.Duration) {
	t.Helper()
	WaitUntilDeploymentAvailableContext(t, context.Background(), options, deploymentName, retries, sleepBetweenRetries)
}

// WaitUntilDeploymentAvailableE waits until all pods within the deployment are ready and started,
// retrying the check for the specified amount of times, sleeping
// for the provided duration between each try.
//
// Deprecated: Use [WaitUntilDeploymentAvailableContextE] instead.
func WaitUntilDeploymentAvailableE(
	t testing.TestingT,
	options *KubectlOptions,
	deploymentName string,
	retries int,
	sleepBetweenRetries time.Duration,
) error {
	return WaitUntilDeploymentAvailableContextE(t, context.Background(), options, deploymentName, retries, sleepBetweenRetries)
}

// IsDeploymentAvailable returns true if all pods within the deployment are ready and started
func IsDeploymentAvailable(deploy *appsv1.Deployment) bool {
	dc := getDeploymentCondition(deploy, appsv1.DeploymentProgressing)
	return dc != nil && dc.Status == v1.ConditionTrue && dc.Reason == "NewReplicaSetAvailable"
}

func getDeploymentCondition(deploy *appsv1.Deployment, cType appsv1.DeploymentConditionType) *appsv1.DeploymentCondition {
	for idx := range deploy.Status.Conditions {
		dc := &deploy.Status.Conditions[idx]
		if dc.Type == cType {
			return dc
		}
	}

	return nil
}
