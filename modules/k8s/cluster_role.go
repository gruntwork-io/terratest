package k8s

import (
	"context"

	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetClusterRoleContextE returns a Kubernetes ClusterRole resource with the given name.
// The ctx parameter supports cancellation and timeouts.
func GetClusterRoleContextE(t testing.TestingT, ctx context.Context, options *KubectlOptions, roleName string) (*rbacv1.ClusterRole, error) {
	clientset, err := GetKubernetesClientFromOptionsContextE(t, ctx, options)
	if err != nil {
		return nil, err
	}

	return clientset.RbacV1().ClusterRoles().Get(ctx, roleName, metav1.GetOptions{})
}

// GetClusterRoleContext returns a Kubernetes ClusterRole resource with the given name.
// The ctx parameter supports cancellation and timeouts.
// This will fail the test if there is an error.
func GetClusterRoleContext(t testing.TestingT, ctx context.Context, options *KubectlOptions, roleName string) *rbacv1.ClusterRole {
	t.Helper()
	role, err := GetClusterRoleContextE(t, ctx, options, roleName)
	require.NoError(t, err)

	return role
}

// GetClusterRole returns a Kubernetes ClusterRole resource with the given name. This will fail the test if there is an error.
//
// Deprecated: Use [GetClusterRoleContext] instead.
func GetClusterRole(t testing.TestingT, options *KubectlOptions, roleName string) *rbacv1.ClusterRole {
	t.Helper()

	return GetClusterRoleContext(t, context.Background(), options, roleName)
}

// GetClusterRoleE returns a Kubernetes ClusterRole resource with the given name.
//
// Deprecated: Use [GetClusterRoleContextE] instead.
func GetClusterRoleE(t testing.TestingT, options *KubectlOptions, roleName string) (*rbacv1.ClusterRole, error) {
	return GetClusterRoleContextE(t, context.Background(), options, roleName)
}
