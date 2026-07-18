package k8s

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
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

		return "UNKNOWN_RESOURCE_TYPE"
	}
}

// makeLabels is a helper to format a map of label key and value pairs into a single string for use as a selector.
func makeLabels(labels map[string]string) string {
	out := make([]string, 0, len(labels))
	for key, value := range labels {
		out = append(out, fmt.Sprintf("%s=%s", key, value))
	}

	return strings.Join(out, ",")
}

// Tunnel is the main struct that configures and manages port forwading tunnels to Kubernetes resources.
type Tunnel struct {
	out            io.Writer
	logger         logger.TestLogger
	kubectlOptions *KubectlOptions
	stopChan       chan struct{}
	readyChan      chan struct{}
	resourceName   string
	localPort      int
	remotePort     int
	resourceType   KubeResourceType
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
		out:            io.Discard,
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
	deploy, err := GetDeploymentContextE(t, context.Background(), tunnel.kubectlOptions, tunnel.resourceName)
	if err != nil {
		return "", err
	}

	selectorLabelsOfPods := makeLabels(deploy.Spec.Selector.MatchLabels)

	deploymentPods, err := ListPodsContextE(t, context.Background(), tunnel.kubectlOptions, metav1.ListOptions{LabelSelector: selectorLabelsOfPods})
	if err != nil {
		return "", err
	}

	for i := range deploymentPods {
		if IsPodAvailable(&deploymentPods[i]) {
			return deploymentPods[i].Name, nil
		}
	}

	return "", DeploymentNotAvailable{deploy}
}

// getAttachablePodForServiceE will find an active pod associated with the Service and return the pod name.
func (tunnel *Tunnel) getAttachablePodForServiceE(t testing.TestingT) (string, error) {
	service, err := GetServiceContextE(t, context.Background(), tunnel.kubectlOptions, tunnel.resourceName)
	if err != nil {
		return "", err
	}

	selectorLabelsOfPods := makeLabels(service.Spec.Selector)

	servicePods, err := ListPodsContextE(t, context.Background(), tunnel.kubectlOptions, metav1.ListOptions{LabelSelector: selectorLabelsOfPods})
	if err != nil {
		return "", err
	}

	for i := range servicePods {
		if IsPodAvailable(&servicePods[i]) {
			return servicePods[i].Name, nil
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

	clientset, err := GetKubernetesClientFromOptionsContextE(t, context.Background(), tunnel.kubectlOptions)
	if err != nil {
		tunnel.logger.Logf(t, "Error creating a new Kubernetes client: %s", err)
		return err
	}

	config := tunnel.kubectlOptions.RestConfig
	if config == nil {
		kubeConfigPath, err := tunnel.kubectlOptions.GetConfigPath(t)
		if err != nil {
			tunnel.logger.Logf(t, "Error getting kube config path: %s", err)
			return err
		}

		config, err = LoadAPIClientConfigE(kubeConfigPath, tunnel.kubectlOptions.ContextName)
		if err != nil {
			tunnel.logger.Logf(t, "Error loading Kubernetes config: %s", err)
			return err
		}
	}

	podName, err := tunnel.getAttachablePodForResourceE(t)
	if err != nil {
		tunnel.logger.Logf(t, "Error finding available pod: %s", err)
		return err
	}

	tunnel.logger.Logf(t, "Selected pod %s to open port forward to", podName)

	var targetPort = tunnel.remotePort

	if tunnel.resourceType == ResourceTypeService {
		service := GetServiceContext(t, context.Background(), tunnel.kubectlOptions, tunnel.resourceName)

		var portFound = false

		for _, portSpec := range service.Spec.Ports {
			if portSpec.Port == int32(targetPort) {
				if portSpec.TargetPort.Type == intstr.String {
					pod, err := GetPodContextE(t, context.Background(), tunnel.kubectlOptions, podName)
					if err != nil {
						return err
					}

					targetPort, err = getPodPortByName(pod, portSpec.TargetPort.String())
					if err != nil {
						tunnel.logger.Logf(t, "Error selecting port by name: %s", err)
						return err
					}

					portFound = true

					break
				}

				targetPort = portSpec.TargetPort.IntValue()
				portFound = true

				break
			}
		}

		if !portFound {
			return TargetPortNotFoundError{TargetPort: targetPort, ServiceName: tunnel.resourceName}
		}
	}

	postEndpoint := clientset.CoreV1().RESTClient().Post()
	namespace := tunnel.kubectlOptions.Namespace
	portForwardCreateURL := postEndpoint.
		Resource("pods").
		Namespace(namespace).
		Name(podName).
		SubResource("portforward").
		URL()

	tunnel.logger.Logf(t, "Using URL %s to create portforward", portForwardCreateURL)

	transport, upgrader, err := spdy.RoundTripperFor(config)
	if err != nil {
		tunnel.logger.Logf(t, "Error creating http client: %s", err)
		return err
	}

	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: transport}, "POST", portForwardCreateURL)

	if tunnel.localPort == 0 {
		tunnel.logger.Logf(t, "Requested local port is 0. Selecting an open port on host system")

		tunnel.localPort, err = GetAvailablePortContextE(t, context.Background())
		if err != nil {
			tunnel.logger.Logf(t, "Error getting available port: %s", err)
			return err
		}

		tunnel.logger.Logf(t, "Selected port %d", tunnel.localPort)
		globalMutex.Lock()
		defer globalMutex.Unlock()
	}

	ports := []string{fmt.Sprintf("%d:%d", tunnel.localPort, targetPort)}

	portforwarder, err := portforward.New(dialer, ports, tunnel.stopChan, tunnel.readyChan, tunnel.out, tunnel.out)
	if err != nil {
		tunnel.logger.Logf(t, "Error creating port forwarding tunnel: %s", err)
		return err
	}

	errChan := make(chan error)

	go func() {
		errChan <- portforwarder.ForwardPorts()
	}()

	select {
	case err = <-errChan:
		tunnel.logger.Logf(t, "Error starting port forwarding tunnel: %s", err)
		return err
	case <-portforwarder.Ready:
		tunnel.logger.Logf(t, "Successfully created port forwarding tunnel")
		return nil
	}
}

// GetAvailablePortContext retrieves an available port on the host machine using the provided context. This delegates the
// port selection to the golang net library by starting a server and then checking the port that the server is using.
// This will fail the test if it could not find an available port.
func GetAvailablePortContext(t testing.TestingT, ctx context.Context) int {
	port, err := GetAvailablePortContextE(t, ctx)
	require.NoError(t, err)

	return port
}

// GetAvailablePortContextE retrieves an available port on the host machine using the provided context. This delegates
// the port selection to the golang net library by starting a server and then checking the port that the server is using.
func GetAvailablePortContextE(t testing.TestingT, ctx context.Context) (int, error) {
	l, err := (&net.ListenConfig{}).Listen(ctx, "tcp", ":0")
	if err != nil {
		return 0, err
	}

	defer func() { _ = l.Close() }()

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

func getPodPortByName(pod *corev1.Pod, portName string) (int, error) {
	if pod == nil {
		return 0, ErrNilPod
	}

	for i := range pod.Spec.Containers {
		for j := range pod.Spec.Containers[i].Ports {
			if pod.Spec.Containers[i].Ports[j].Name == portName {
				return int(pod.Spec.Containers[i].Ports[j].ContainerPort), nil
			}
		}
	}

	return 0, PortNotFoundInPodError{PortName: portName, PodName: pod.Name}
}
