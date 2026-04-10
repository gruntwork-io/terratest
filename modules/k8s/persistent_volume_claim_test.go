//go:build kubeall || kubernetes
// +build kubeall kubernetes

// NOTE: we have build tags to differentiate kubernetes tests from non-kubernetes tests. This is done because minikube
// is heavy and can interfere with docker related tests in terratest. Specifically, many of the tests start to fail with
// `connection refused` errors from `minikube`. To avoid overloading the system, we run the kubernetes tests and helm
// tests separately from the others. This may not be necessary if you have a sufficiently powerful machine.  We
// recommend at least 4 cores and 16GB of RAM if you want to run all the tests together.

package k8s_test

import (
	"strings"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/k8s"

	"github.com/stretchr/testify/require"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"github.com/gruntwork-io/terratest/modules/random"
)

func TestListPersistentVolumeClaimsReturnsPersistentVolumeClaimsInNamespace(t *testing.T) {
	t.Parallel()

	pvcName := "test-dummy-pvc"
	namespace := strings.ToLower(random.UniqueID())
	options := k8s.NewKubectlOptions("", "", namespace)

	configData := renderFixtureYamlTemplate(namespace, pvcName)
	defer k8s.KubectlDeleteFromString(t, options, configData)

	k8s.KubectlApplyFromString(t, options, configData)

	pvcs := k8s.ListPersistentVolumeClaims(t, options, metav1.ListOptions{})
	require.Len(t, pvcs, 1)
	pvc := pvcs[0]
	require.Equal(t, pvc.Name, pvcName)
	require.Equal(t, pvc.Namespace, namespace)
}

func TestListPersistentVolumeClaimsReturnsZeroPersistentVolumeClaimsIfNoneCreated(t *testing.T) {
	t.Parallel()

	namespace := strings.ToLower(random.UniqueID())
	options := k8s.NewKubectlOptions("", "", namespace)

	k8s.CreateNamespace(t, options, namespace)
	defer k8s.DeleteNamespace(t, options, namespace)

	pvcs := k8s.ListPersistentVolumeClaims(t, options, metav1.ListOptions{})
	require.Empty(t, pvcs)
}

func TestGetPersistentVolumeClaimEReturnsErrorForNonExistantPersistentVolumeClaim(t *testing.T) {
	t.Parallel()

	options := k8s.NewKubectlOptions("", "", "default")
	_, err := k8s.GetPersistentVolumeClaimE(t, options, "non-existent")
	require.Error(t, err)
}

func TestGetPersistentVolumeClaimReturnsCorrectPersistentVolumeClaimInCorrectNamespace(t *testing.T) {
	t.Parallel()

	pvcName := "test-dummy-pvc"
	namespace := strings.ToLower(random.UniqueID())
	options := k8s.NewKubectlOptions("", "", namespace)

	configData := renderFixtureYamlTemplate(namespace, pvcName)
	defer k8s.KubectlDeleteFromString(t, options, configData)

	k8s.KubectlApplyFromString(t, options, configData)

	pvc := k8s.GetPersistentVolumeClaim(t, options, pvcName)
	require.Equal(t, pvc.Name, pvcName)
	require.Equal(t, pvc.Namespace, namespace)
}

func TestWaitUntilPersistentVolumeClaimInGivenStatusPhase(t *testing.T) {
	t.Parallel()

	pvcName := "test-dummy-pvc"
	namespace := strings.ToLower(random.UniqueID())
	pvcBoundStatusPhase := corev1.ClaimBound
	options := k8s.NewKubectlOptions("", "", namespace)

	configData := renderFixtureYamlTemplate(namespace, pvcName)
	defer k8s.KubectlDeleteFromString(t, options, configData)

	k8s.KubectlApplyFromString(t, options, configData)

	k8s.WaitUntilPersistentVolumeClaimInStatus(t, options, pvcName, &pvcBoundStatusPhase, 60, 1*time.Second)
}

func TestWaitUntilPersistentVolumeClaimInStatusEReturnsErrorWhenWaitingForAnUnexistentPvc(t *testing.T) {
	t.Parallel()

	pvcBoundStatusPhase := corev1.ClaimBound
	options := k8s.NewKubectlOptions("", "", "default")
	err := k8s.WaitUntilPersistentVolumeClaimInStatusE(t, options, "non-existent", &pvcBoundStatusPhase, 3, 1*time.Second)
	require.Error(t, err)
}

func TestWaitUntilPersistentVolumeClaimInStatusEReturnsErrorWhenTimesOut(t *testing.T) {
	t.Parallel()

	pvcName := "test-dummy-pvc"
	pvcLostStatusPhase := corev1.ClaimLost
	namespace := strings.ToLower(random.UniqueID())
	options := k8s.NewKubectlOptions("", "", namespace)

	configData := renderFixtureYamlTemplate(namespace, pvcName)
	defer k8s.KubectlDeleteFromString(t, options, configData)

	k8s.KubectlApplyFromString(t, options, configData)

	err := k8s.WaitUntilPersistentVolumeClaimInStatusE(t, options, pvcName, &pvcLostStatusPhase, 5, 1*time.Second)
	require.Error(t, err)
}

func TestIsPersistentVolumeClaimInStatusReturnsFalseIfPvcIsNil(t *testing.T) {
	t.Parallel()

	result := k8s.IsPersistentVolumeClaimInStatus(nil, nil)
	require.False(t, result)
}

const pvcFixtureYamlTemplate = `---
apiVersion: v1
kind: Namespace
metadata:
  name: __namespace__
---
apiVersion: v1
kind: PersistentVolume
metadata:
  name: __namespace__
spec:
  capacity:
    storage: 10Mi
  accessModes:
    - ReadWriteOnce
  hostPath:
    path: "/tmp/__namespace__"
---
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
    namespace: __namespace__
    name: __pvcName__
spec:
    accessModes:
        - ReadWriteOnce
    resources:
        requests:
            storage: 10Mi
---
apiVersion: v1
kind: Pod
metadata:
  name: test-pvc-pod
  namespace: __namespace__
spec:
  volumes:
    - name: test-pvc-volume
      persistentVolumeClaim:
        claimName: __pvcName__
  containers:
    - name: test-pvc-image
      image: nginx
      volumeMounts:
        - mountPath: "/tmp/foo"
          name: test-pvc-volume
`

func renderFixtureYamlTemplate(namespace, pvcName string) string {
	return strings.ReplaceAll(strings.ReplaceAll(pvcFixtureYamlTemplate, "__namespace__", namespace), "__pvcName__", pvcName)
}
