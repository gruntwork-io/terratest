package k8s

import (
	"context"

	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetRole returns a Kubernetes role resource in the provided namespace with the given name. The namespace used
// is the one provided in the KubectlOptions. This will fail the test if there is an error.
func GetRole(t testing.TestingT, options *KubectlOptions, roleName string) *rbacv1.Role {
	role, err := GetRoleE(t, options, roleName)
	require.NoError(t, err)
	return role
}

// GetRoleE returns a Kubernetes role resource in the provided namespace with the given name. The namespace used
// is the one provided in the KubectlOptions.
func GetRoleE(t testing.TestingT, options *KubectlOptions, roleName string) (*rbacv1.Role, error) {
	clientset, err := GetKubernetesClientFromOptionsE(t, options)
	if err != nil {
		return nil, err
	}
	return clientset.RbacV1().Roles(options.Namespace).Get(context.Background(), roleName, metav1.GetOptions{})
}
