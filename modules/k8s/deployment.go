package k8s

import (
	"context"
	"fmt"
	"time"

	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/retry"
	"github.com/gruntwork-io/terratest/modules/testing"
)

// ListDeployments will look for deployments in the given namespace that match the given filters and return them. This will
// fail the test if there is an error.
func ListDeployments(t testing.TestingT, options *KubectlOptions, filters metav1.ListOptions) []appsv1.Deployment {
	deployment, err := ListDeploymentsE(t, options, filters)
	require.NoError(t, err)
	return deployment
}

// ListDeploymentsE will look for deployments in the given namespace that match the given filters and return them.
func ListDeploymentsE(t testing.TestingT, options *KubectlOptions, filters metav1.ListOptions) ([]appsv1.Deployment, error) {
	clientset, err := GetKubernetesClientFromOptionsE(t, options)
	if err != nil {
		return nil, err
	}
	deployments, err := clientset.AppsV1().Deployments(options.Namespace).List(context.Background(), filters)
	if err != nil {
		return nil, err
	}
	return deployments.Items, nil
}

// GetDeployment returns a Kubernetes deployment resource in the provided namespace with the given name. This will
// fail the test if there is an error.
func GetDeployment(t testing.TestingT, options *KubectlOptions, deploymentName string) *appsv1.Deployment {
	deployment, err := GetDeploymentE(t, options, deploymentName)
	require.NoError(t, err)
	return deployment
}

// GetDeploymentE returns a Kubernetes deployment resource in the provided namespace with the given name.
func GetDeploymentE(t testing.TestingT, options *KubectlOptions, deploymentName string) (*appsv1.Deployment, error) {
	clientset, err := GetKubernetesClientFromOptionsE(t, options)
	if err != nil {
		return nil, err
	}
	return clientset.AppsV1().Deployments(options.Namespace).Get(context.Background(), deploymentName, metav1.GetOptions{})
}

// WaitUntilDeploymentAvailableE waits until all pods within the deployment are ready and started,
// retrying the check for the specified amount of times, sleeping
// for the provided duration between each try.
// This will fail the test if there is an error.
func WaitUntilDeploymentAvailable(t testing.TestingT, options *KubectlOptions, deploymentName string, retries int, sleepBetweenRetries time.Duration) {
	require.NoError(t, WaitUntilDeploymentAvailableE(t, options, deploymentName, retries, sleepBetweenRetries))
}

// WaitUntilDeploymentAvailableE waits until all pods within the deployment are ready and started,
// retrying the check for the specified amount of times, sleeping
// for the provided duration between each try.
func WaitUntilDeploymentAvailableE(
	t testing.TestingT,
	options *KubectlOptions,
	deploymentName string,
	retries int,
	sleepBetweenRetries time.Duration,
) error {
	statusMsg := fmt.Sprintf("Wait for deployment %s to be provisioned.", deploymentName)
	message, err := retry.DoWithRetryE(
		t,
		statusMsg,
		retries,
		sleepBetweenRetries,
		func() (string, error) {
			deployment, err := GetDeploymentE(t, options, deploymentName)
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
		logger.Logf(t, "Timedout waiting for Deployment to be provisioned: %s", err)
		return err
	}
	logger.Logf(t, message)
	return nil
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
