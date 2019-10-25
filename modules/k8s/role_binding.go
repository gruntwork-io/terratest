package k8s

import (
	"testing"

	"github.com/stretchr/testify/require"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetRoleBinding returns a Kubernetes role binfing resource in the provided namespace with the given name. The namespace used
// is the one provided in the KubectlOptions. This will fail the test if there is an error.
func GetRoleBinding(t *testing.T, options *KubectlOptions, roleBindingName string) *rbacv1.RoleBinding {
	roleBinding, err := GetRoleBindingE(t, options, roleBindingName)
	require.NoError(t, err)
	return roleBinding
}

// GetRoleBindingE returns a Kubernetes role binding resource in the provided namespace with the given name. The namespace used
// is the one provided in the KubectlOptions.
func GetRoleBindingE(t *testing.T, options *KubectlOptions, roleBindingName string) (*rbacv1.RoleBinding, error) {
	clientset, err := GetKubernetesClientFromOptionsE(t, options)
	if err != nil {
		return nil, err
	}
	return clientset.RbacV1().RoleBindings(options.Namespace).Get(roleBindingName, metav1.GetOptions{})
}
