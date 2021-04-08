package k8s

import (
	"context"

	"github.com/gruntwork-io/go-commons/errors"
	"github.com/stretchr/testify/require"
	authv1 "k8s.io/api/authorization/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/testing"
)

// CanIDo returns whether or not the provided action is allowed by the client configured by the provided kubectl option.
// This will fail if there are any errors accessing the kubernetes API (but not if the action is denied).
func CanIDo(t testing.TestingT, options *KubectlOptions, action authv1.ResourceAttributes) bool {
	allowed, err := CanIDoE(t, options, action)
	require.NoError(t, err)
	return allowed
}

// CanIDoE returns whether or not the provided action is allowed by the client configured by the provided kubectl option.
// This will an error if there are problems accessing the kubernetes API (but not if the action is simply denied).
func CanIDoE(t testing.TestingT, options *KubectlOptions, action authv1.ResourceAttributes) (bool, error) {
	clientset, err := GetKubernetesClientFromOptionsE(t, options)
	if err != nil {
		return false, err
	}
	check := authv1.SelfSubjectAccessReview{
		Spec: authv1.SelfSubjectAccessReviewSpec{ResourceAttributes: &action},
	}
	resp, err := clientset.AuthorizationV1().SelfSubjectAccessReviews().Create(context.Background(), &check, metav1.CreateOptions{})
	if err != nil {
		return false, errors.WithStackTrace(err)
	}
	if !resp.Status.Allowed {
		logger.Logf(t, "Denied action %s on resource %s with name '%s' for reason %s", action.Verb, action.Resource, action.Name, resp.Status.Reason)
	}
	return resp.Status.Allowed, nil
}
