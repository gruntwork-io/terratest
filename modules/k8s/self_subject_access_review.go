package k8s

import (
	"context"

	"github.com/gruntwork-io/go-commons/errors"
	"github.com/stretchr/testify/require"
	authv1 "k8s.io/api/authorization/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/gruntwork-io/terratest/modules/testing"
)

// CanIDoContextE returns whether or not the provided action is allowed by the client configured by the provided kubectl option.
// This will return an error if there are problems accessing the kubernetes API (but not if the action is simply denied).
// The ctx parameter supports cancellation and timeouts.
//
//nolint:gocritic // hugeParam: cannot change public function signature
func CanIDoContextE(t testing.TestingT, ctx context.Context, options *KubectlOptions, action authv1.ResourceAttributes) (bool, error) {
	clientset, err := GetKubernetesClientFromOptionsContextE(t, ctx, options)
	if err != nil {
		return false, err
	}

	check := authv1.SelfSubjectAccessReview{
		Spec: authv1.SelfSubjectAccessReviewSpec{ResourceAttributes: &action},
	}

	resp, err := clientset.AuthorizationV1().SelfSubjectAccessReviews().Create(ctx, &check, metav1.CreateOptions{})
	if err != nil {
		return false, errors.WithStackTrace(err)
	}

	if !resp.Status.Allowed {
		options.Logger.Logf(t, "Denied action %s on resource %s with name '%s' for reason %s", action.Verb, action.Resource, action.Name, resp.Status.Reason)
	}

	return resp.Status.Allowed, nil
}

// CanIDoContext returns whether or not the provided action is allowed by the client configured by the provided kubectl option.
// The ctx parameter supports cancellation and timeouts.
// This will fail if there are any errors accessing the kubernetes API (but not if the action is denied).
//
//nolint:gocritic // hugeParam: cannot change public function signature
func CanIDoContext(t testing.TestingT, ctx context.Context, options *KubectlOptions, action authv1.ResourceAttributes) bool {
	t.Helper()
	allowed, err := CanIDoContextE(t, ctx, options, action)
	require.NoError(t, err)

	return allowed
}

// CanIDo returns whether or not the provided action is allowed by the client configured by the provided kubectl option.
// This will fail if there are any errors accessing the kubernetes API (but not if the action is denied).
//
// Deprecated: Use [CanIDoContext] instead.
//
//nolint:gocritic // hugeParam: cannot change public function signature
func CanIDo(t testing.TestingT, options *KubectlOptions, action authv1.ResourceAttributes) bool {
	t.Helper()

	return CanIDoContext(t, context.Background(), options, action)
}

// CanIDoE returns whether or not the provided action is allowed by the client configured by the provided kubectl option.
// This will an error if there are problems accessing the kubernetes API (but not if the action is simply denied).
//
// Deprecated: Use [CanIDoContextE] instead.
//
//nolint:gocritic // hugeParam: cannot change public function signature
func CanIDoE(t testing.TestingT, options *KubectlOptions, action authv1.ResourceAttributes) (bool, error) {
	return CanIDoContextE(t, context.Background(), options, action)
}
