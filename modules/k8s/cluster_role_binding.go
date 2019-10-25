package k8s

import (
	"testing"

	"github.com/stretchr/testify/require"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetClusterRoleBinding returns a Kubernetes role binfing resource with the given name. This will fail the test if there is an error.
func GetClusterRoleBinding(t *testing.T, options *KubectlOptions, crBindingName string) *rbacv1.ClusterRoleBinding {
	crBinding, err := GetClusterRoleBindingE(t, options, crBindingName)
	require.NoError(t, err)
	return crBinding
}

// GetClusterRoleBindingE returns a Kubernetes role binding resource with the given name.
func GetClusterRoleBindingE(t *testing.T, options *KubectlOptions, crBindingName string) (*rbacv1.ClusterRoleBinding, error) {
	clientset, err := GetKubernetesClientFromOptionsE(t, options)
	if err != nil {
		return nil, err
	}
	return clientset.RbacV1().ClusterRoleBindings().Get(crBindingName, metav1.GetOptions{})
}
