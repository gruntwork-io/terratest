//go:build kubeall || kubernetes
// +build kubeall kubernetes

// NOTE: we have build tags to differentiate kubernetes tests from non-kubernetes tests. This is done because minikube
// is heavy and can interfere with docker related tests in terratest. Specifically, many of the tests start to fail with
// `connection refused` errors from `minikube`. To avoid overloading the system, we run the kubernetes tests and helm
// tests separately from the others. This may not be necessary if you have a sufficiently powerful machine.  We
// recommend at least 4 cores and 16GB of RAM if you want to run all the tests together.

package k8s

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"github.com/gruntwork-io/terratest/modules/random"
)

func TestListPersistentVolumesReturnsAllPersistentVolumes(t *testing.T) {
	t.Parallel()

	numPvFound := 0
	pvNames := map[string]struct{}{
		strings.ToLower(random.UniqueId()): {},
		strings.ToLower(random.UniqueId()): {},
		strings.ToLower(random.UniqueId()): {},
	}

	options := NewKubectlOptions("", "", "")
	for pvName := range pvNames {
		pv := fmt.Sprintf(PvFixtureYamlTemplate, pvName, pvName)
		defer KubectlDeleteFromString(t, options, pv)
		KubectlApplyFromString(t, options, pv)
	}

	pvs := ListPersistentVolumes(t, options, metav1.ListOptions{})
	for _, pv := range pvs {
		if _, ok := pvNames[pv.Name]; ok {
			numPvFound++
		}
	}

	require.Equal(t, numPvFound, len(pvNames))
}

func TestListPersistentVolumesReturnsZeroPersistentVolumesIfNoneCreated(t *testing.T) {
	t.Parallel()

	options := NewKubectlOptions("", "", "")
	pvs := ListPersistentVolumes(t, options, metav1.ListOptions{})
	require.Equal(t, 0, len(pvs))
}

func TestGetPersistentVolumeEReturnsErrorForNonExistentPersistentVolumes(t *testing.T) {
	t.Parallel()

	options := NewKubectlOptions("", "", "")
	_, err := GetPersistentVolumeE(t, options, "non-existent")
	require.Error(t, err)
}

func TestGetPersistentVolumeReturnsCorrectPersistentVolume(t *testing.T) {
	t.Parallel()

	pvName := strings.ToLower(random.UniqueId())
	options := NewKubectlOptions("", "", "")
	configData := fmt.Sprintf(PvFixtureYamlTemplate, pvName, pvName)
	defer KubectlDeleteFromString(t, options, configData)
	KubectlApplyFromString(t, options, configData)

	pv := GetPersistentVolume(t, options, pvName)
	require.Equal(t, pv.Name, pvName)
}

func TestWaitUntilPersistentVolumeInTheGivenStatusPhase(t *testing.T) {
	t.Parallel()

	pvName := strings.ToLower(random.UniqueId())
	pvAvailableStatusPhase := corev1.VolumeAvailable

	options := NewKubectlOptions("", "", pvName)
	configData := fmt.Sprintf(PvFixtureYamlTemplate, pvName, pvName)
	KubectlApplyFromString(t, options, configData)
	defer KubectlDeleteFromString(t, options, configData)

	WaitUntilPersistentVolumeInStatus(t, options, pvName, &pvAvailableStatusPhase, 60, 1*time.Second)
}

const PvFixtureYamlTemplate = `---
apiVersion: v1
kind: PersistentVolume
metadata:
  name: %s
spec:
  capacity:
    storage: 10Mi
  accessModes:
    - ReadWriteOnce
  hostPath:
    path: "/tmp/%s"
`
