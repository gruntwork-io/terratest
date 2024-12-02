//go:build kubeall || helm
// +build kubeall helm

// NOTE: we have build tags to differentiate kubernetes tests from non-kubernetes tests, and further differentiate helm
// tests. This is done because minikube is heavy and can interfere with docker related tests in terratest. Similarly,
// helm can overload the minikube system and thus interfere with the other kubernetes tests. To avoid overloading the
// system, we run the kubernetes tests and helm tests separately from the others.

package helm

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"

	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/random"
)

// Test that we can render locally a remote chart (e.g bitnami/nginx)
func TestRemoteChartRender(t *testing.T) {
	const (
		remoteChartSource  = "https://charts.bitnami.com/bitnami"
		remoteChartName    = "nginx"
		remoteChartVersion = "13.2.24"
	)

	t.Parallel()

	namespaceName := fmt.Sprintf(
		"%s-%s",
		strings.ToLower(t.Name()),
		strings.ToLower(random.UniqueId()),
	)

	releaseName := remoteChartName

	options := &Options{
		SetValues: map[string]string{
			"image.repository": remoteChartName,
			"image.registry":   "",
			"image.tag":        remoteChartVersion,
		},
		KubectlOptions: k8s.NewKubectlOptions("", "", namespaceName),
		Logger:         logger.Discard,
	}

	// Run RenderTemplate to render the template and capture the output. Note that we use the version without `E`, since
	// we want to assert that the template renders without any errors.
	output := RenderRemoteTemplate(t, options, remoteChartSource, releaseName, []string{"templates/deployment.yaml"})

	// Now we use kubernetes/client-go library to render the template output into the Deployment struct. This will
	// ensure the Deployment resource is rendered correctly.
	var deployment appsv1.Deployment
	UnmarshalK8SYaml(t, output, &deployment)

	// Verify the namespace matches the expected supplied namespace.
	require.Equal(t, namespaceName, deployment.Namespace)

	// Finally, we verify the deployment pod template spec is set to the expected container image value
	expectedContainerImage := remoteChartName + ":" + remoteChartVersion
	deploymentContainers := deployment.Spec.Template.Spec.Containers
	require.Equal(t, len(deploymentContainers), 1)
	require.Equal(t, deploymentContainers[0].Image, expectedContainerImage)
}

// Test that we can dump all the manifest locally a remote chart (e.g bitnami/nginx)
// so that I can use them later to compare between two versions of the same chart for example
func TestRemoteChartRenderDump(t *testing.T) {
	t.Parallel()
	renderChartDump(t, "13.2.20", t.TempDir())
}

// Test that we can diff all the manifest to a local snapshot using a remote chart (e.g bitnami/nginx)
func TestRemoteChartRenderDiff(t *testing.T) {
	t.Parallel()

	initialSnapshot := t.TempDir()
	updatedSnapshot := t.TempDir()
	renderChartDump(t, "5.0.0", initialSnapshot)
	output := renderChartDump(t, "5.1.0", updatedSnapshot)

	options := &Options{
		Logger:       logger.Default,
		SnapshotPath: initialSnapshot,
	}
	// diff in: spec.initContainers.preserve-logs-symlinks.imag, spec.containers.nginx.image, tls certificates
	require.Equal(t, 5, DiffAgainstSnapshot(t, options, output, "nginx"))
}

// render chart dump and return the rendered output
func renderChartDump(t *testing.T, remoteChartVersion, snapshotDir string) string {
	const (
		remoteChartSource = "https://charts.bitnami.com/bitnami"
		remoteChartName   = "nginx"
		// need to set a fix name for the namespace, so it is not flag as a difference
		namespaceName = "dump-ns"
	)

	releaseName := remoteChartName

	options := &Options{
		SetValues: map[string]string{
			"image.repository": remoteChartName,
			"image.registry":   "",
			"image.tag":        remoteChartVersion,
		},
		KubectlOptions: k8s.NewKubectlOptions("", "", namespaceName),
		Logger:         logger.Discard,
	}

	// Run RenderTemplate to render the template and capture the output. Note that we use the version without `E`, since
	// we want to assert that the template renders without any errors.
	output := RenderRemoteTemplate(t, options, remoteChartSource, releaseName, []string{})

	// Now we use kubernetes/client-go library to render the template output into the Deployment struct. This will
	// ensure the Deployment resource is rendered correctly.
	var deployment appsv1.Deployment
	UnmarshalK8SYaml(t, output, &deployment)

	// Verify the namespace matches the expected supplied namespace.
	require.Equal(t, namespaceName, deployment.Namespace)

	// write chart manifest to a local filesystem directory
	options = &Options{
		Logger:       logger.Default,
		SnapshotPath: snapshotDir,
	}
	UpdateSnapshot(t, options, output, releaseName)
	return output
}
