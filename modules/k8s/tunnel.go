package k8s

// The following code is a fork of the Helm client. The main differences are:
// - Support testing context for better logging
// - Support resources other than pods
// See: https://github.com/helm/helm/blob/master/pkg/kube/tunnel.go

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/testing"
)

// Global lock to synchronize port selections
var globalMutex sync.Mutex

// KubeResourceType is an enum representing known resource types that can support port forwarding
type KubeResourceType int

const (
	// ResourceTypePod is a k8s pod kind identifier
	ResourceTypePod KubeResourceType = iota
	// ResourceTypeDeployment is a k8s deployment kind identifier
	ResourceTypeDeployment
	// ResourceTypeService is a k8s service kind identifier
	ResourceTypeService
)

func (resourceType KubeResourceType) String() string {
	switch resourceType {
	case ResourceTypeDeployment:
		return "deploy"
	case ResourceTypePod:
		return "pod"
	case ResourceTypeService:
		return "svc"
	default:
		// This should not happen
		return "UNKNOWN_RESOURCE_TYPE"
	}
}

// makeLabels is a helper to format a map of label key and value pairs into a single string for use as a selector.
func makeLabels(labels map[string]string) string {
	out := []string{}
	for key, value := range labels {
		out = append(out, fmt.Sprintf("%s=%s", key, value))
	}
	return strings.Join(out, ",")
}

// Tunnel is the main struct that configures and manages port forwading tunnels to Kubernetes resources.
type Tunnel struct {
	out            io.Writer
	localPort      int
	remotePort     int
	kubectlOptions *KubectlOptions
	resourceType   KubeResourceType
	resourceName   string
	logger         logger.TestLogger
	stopChan       chan struct{}
	readyChan      chan struct{}
}

// NewTunnel creates a new tunnel with NewTunnelWithLogger, setting logger.Terratest as the logger.
func NewTunnel(kubectlOptions *KubectlOptions, resourceType KubeResourceType, resourceName string, local int, remote int) *Tunnel {
	return NewTunnelWithLogger(kubectlOptions, resourceType, resourceName, local, remote, logger.Terratest)
}

// NewTunnelWithLogger will create a new Tunnel struct with the provided logger.
// Note that if you use 0 for the local port, an open port on the host system
// will be selected automatically, and the Tunnel struct will be updated with the selected port.
func NewTunnelWithLogger(
	kubectlOptions *KubectlOptions,
	resourceType KubeResourceType,
	resourceName string,
	local int,
	remote int,
	logger logger.TestLogger,
) *Tunnel {
	return &Tunnel{
		out:            ioutil.Discard,
		localPort:      local,
		remotePort:     remote,
		kubectlOptions: kubectlOptions,
		resourceType:   resourceType,
		resourceName:   resourceName,
		logger:         logger,
		stopChan:       make(chan struct{}, 1),
		readyChan:      make(chan struct{}, 1),
	}
}

// Endpoint returns the tunnel endpoint
func (tunnel *Tunnel) Endpoint() string {
	return fmt.Sprintf("localhost:%d", tunnel.localPort)
}

// Close disconnects a tunnel connection by closing the StopChan, thereby stopping the goroutine.
func (tunnel *Tunnel) Close() {
	close(tunnel.stopChan)
}

// getAttachablePodForResource will find a pod that can be port forwarded to given the provided resource type and return
// the name.
func (tunnel *Tunnel) getAttachablePodForResourceE(t testing.TestingT) (string, error) {
	switch tunnel.resourceType {
	case ResourceTypePod:
		return tunnel.resourceName, nil
	case ResourceTypeService:
		return tunnel.getAttachablePodForServiceE(t)
	case ResourceTypeDeployment:
		return tunnel.getAttachablePodForDeploymentE(t)
	default:
		return "", UnknownKubeResourceType{tunnel.resourceType}
	}
}

// getAttachablePodForDeploymentE will find an active pod associated with the Deployment and return the pod name.
func (tunnel *Tunnel) getAttachablePodForDeploymentE(t testing.TestingT) (string, error) {
	deploy, err := GetDeploymentE(t, tunnel.kubectlOptions, tunnel.resourceName)
	if err != nil {
		return "", err
	}
	selectorLabelsOfPods := makeLabels(deploy.Spec.Selector.MatchLabels)
	deploymentPods, err := ListPodsE(t, tunnel.kubectlOptions, metav1.ListOptions{LabelSelector: selectorLabelsOfPods})
	if err != nil {
		return "", err
	}
	for _, pod := range deploymentPods {
		if IsPodAvailable(&pod) {
			return pod.Name, nil
		}
	}
	return "", DeploymentNotAvailable{deploy}
}

// getAttachablePodForServiceE will find an active pod associated with the Service and return the pod name.
func (tunnel *Tunnel) getAttachablePodForServiceE(t testing.TestingT) (string, error) {
	service, err := GetServiceE(t, tunnel.kubectlOptions, tunnel.resourceName)
	if err != nil {
		return "", err
	}
	selectorLabelsOfPods := makeLabels(service.Spec.Selector)
	servicePods, err := ListPodsE(t, tunnel.kubectlOptions, metav1.ListOptions{LabelSelector: selectorLabelsOfPods})
	if err != nil {
		return "", err
	}
	for _, pod := range servicePods {
		if IsPodAvailable(&pod) {
			return pod.Name, nil
		}
	}
	return "", ServiceNotAvailable{service}
}

// ForwardPort opens a tunnel to a kubernetes resource, as specified by the provided tunnel struct. This will fail the
// test if there is an error attempting to open the port.
func (tunnel *Tunnel) ForwardPort(t testing.TestingT) {
	require.NoError(t, tunnel.ForwardPortE(t))
}

// ForwardPortE opens a tunnel to a kubernetes resource, as specified by the provided tunnel struct.
func (tunnel *Tunnel) ForwardPortE(t testing.TestingT) error {
	tunnel.logger.Logf(
		t,
		"Creating a port forwarding tunnel for resource %s/%s routing local port %d to remote port %d",
		tunnel.resourceType.String(),
		tunnel.resourceName,
		tunnel.localPort,
		tunnel.remotePort,
	)

	// Prepare a kubernetes client for the client-go library
	clientset, err := GetKubernetesClientFromOptionsE(t, tunnel.kubectlOptions)
	if err != nil {
		tunnel.logger.Logf(t, "Error creating a new Kubernetes client: %s", err)
		return err
	}
	kubeConfigPath, err := tunnel.kubectlOptions.GetConfigPath(t)
	if err != nil {
		tunnel.logger.Logf(t, "Error getting kube config path: %s", err)
		return err
	}
	config, err := LoadApiClientConfigE(kubeConfigPath, tunnel.kubectlOptions.ContextName)
	if err != nil {
		tunnel.logger.Logf(t, "Error loading Kubernetes config: %s", err)
		return err
	}

	// Find the pod to port forward to
	podName, err := tunnel.getAttachablePodForResourceE(t)
	if err != nil {
		tunnel.logger.Logf(t, "Error finding available pod: %s", err)
		return err
	}
	tunnel.logger.Logf(t, "Selected pod %s to open port forward to", podName)

	// Build a url to the portforward endpoint
	// example: http://localhost:8080/api/v1/namespaces/helm/pods/tiller-deploy-9itlq/portforward
	postEndpoint := clientset.CoreV1().RESTClient().Post()
	namespace := tunnel.kubectlOptions.Namespace
	portForwardCreateURL := postEndpoint.
		Resource("pods").
		Namespace(namespace).
		Name(podName).
		SubResource("portforward").
		URL()

	tunnel.logger.Logf(t, "Using URL %s to create portforward", portForwardCreateURL)

	// Construct the spdy client required by the client-go portforward library
	transport, upgrader, err := spdy.RoundTripperFor(config)
	if err != nil {
		tunnel.logger.Logf(t, "Error creating http client: %s", err)
		return err
	}
	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: transport}, "POST", portForwardCreateURL)

	// If the localport is 0, get an available port before continuing. We do this here instead of relying on the
	// underlying portforwarder library, because the portforwarder library does not expose the selected local port in a
	// machine readable manner.
	// Synchronize on the global lock to avoid race conditions with concurrently selecting the same available port,
	// since there is a brief moment between `GetAvailablePort` and `portforwader.ForwardPorts` where the selected port
	// is available for selection again.
	if tunnel.localPort == 0 {
		tunnel.logger.Logf(t, "Requested local port is 0. Selecting an open port on host system")
		tunnel.localPort, err = GetAvailablePortE(t)
		if err != nil {
			tunnel.logger.Logf(t, "Error getting available port: %s", err)
			return err
		}
		tunnel.logger.Logf(t, "Selected port %d", tunnel.localPort)
		globalMutex.Lock()
		defer globalMutex.Unlock()
	}

	// Construct a new PortForwarder struct that manages the instructed port forward tunnel
	ports := []string{fmt.Sprintf("%d:%d", tunnel.localPort, tunnel.remotePort)}
	portforwarder, err := portforward.New(dialer, ports, tunnel.stopChan, tunnel.readyChan, tunnel.out, tunnel.out)
	if err != nil {
		tunnel.logger.Logf(t, "Error creating port forwarding tunnel: %s", err)
		return err
	}

	// Open the tunnel in a goroutine so that it is available in the background. Report errors to the main goroutine via
	// a new channel.
	errChan := make(chan error)
	go func() {
		errChan <- portforwarder.ForwardPorts()
	}()

	// Wait for an error or the tunnel to be ready
	select {
	case err = <-errChan:
		tunnel.logger.Logf(t, "Error starting port forwarding tunnel: %s", err)
		return err
	case <-portforwarder.Ready:
		tunnel.logger.Logf(t, "Successfully created port forwarding tunnel")
		return nil
	}
}

// GetAvailablePort retrieves an available port on the host machine. This delegates the port selection to the golang net
// library by starting a server and then checking the port that the server is using. This will fail the test if it could
// not find an available port.
func GetAvailablePort(t testing.TestingT) int {
	port, err := GetAvailablePortE(t)
	require.NoError(t, err)
	return port
}

// GetAvailablePortE retrieves an available port on the host machine. This delegates the port selection to the golang net
// library by starting a server and then checking the port that the server is using.
func GetAvailablePortE(t testing.TestingT) (int, error) {
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0, err
	}
	defer l.Close()

	_, p, err := net.SplitHostPort(l.Addr().String())
	if err != nil {
		return 0, err
	}
	port, err := strconv.Atoi(p)
	if err != nil {
		return 0, err
	}
	return port, err
}
