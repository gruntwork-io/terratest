package k8s

import (
	"context"

	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetRoleContextE returns a Kubernetes role resource in the provided namespace with the given name.
// The namespace used is the one provided in the KubectlOptions.
// The ctx parameter supports cancellation and timeouts.
func GetRoleContextE(t testing.TestingT, ctx context.Context, options *KubectlOptions, roleName string) (*rbacv1.Role, error) {
	clientset, err := GetKubernetesClientFromOptionsContextE(t, ctx, options)
	if err != nil {
		return nil, err
	}

	return clientset.RbacV1().Roles(options.Namespace).Get(ctx, roleName, metav1.GetOptions{})
}

// GetRoleContext returns a Kubernetes role resource in the provided namespace with the given name.
// The namespace used is the one provided in the KubectlOptions.
// The ctx parameter supports cancellation and timeouts.
// This will fail the test if there is an error.
func GetRoleContext(t testing.TestingT, ctx context.Context, options *KubectlOptions, roleName string) *rbacv1.Role {
	t.Helper()
	role, err := GetRoleContextE(t, ctx, options, roleName)
	require.NoError(t, err)

	return role
}

// GetRole returns a Kubernetes role resource in the provided namespace with the given name. The namespace used
// is the one provided in the KubectlOptions. This will fail the test if there is an error.
//
// Deprecated: Use [GetRoleContext] instead.
func GetRole(t testing.TestingT, options *KubectlOptions, roleName string) *rbacv1.Role {
	t.Helper()

	return GetRoleContext(t, context.Background(), options, roleName)
}

// GetRoleE returns a Kubernetes role resource in the provided namespace with the given name. The namespace used
// is the one provided in the KubectlOptions.
//
// Deprecated: Use [GetRoleContextE] instead.
func GetRoleE(t testing.TestingT, options *KubectlOptions, roleName string) (*rbacv1.Role, error) {
	return GetRoleContextE(t, context.Background(), options, roleName)
}
