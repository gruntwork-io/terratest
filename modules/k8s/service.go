package k8s

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/random"
	"github.com/gruntwork-io/terratest/modules/core/v2/retry"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
)

// ListServicesContextE looks up services in the given namespace that match the given filters and return them.
// The ctx parameter supports cancellation and timeouts.
//
//nolint:gocritic // hugeParam: cannot change public function signature
func ListServicesContextE(t testing.TestingT, ctx context.Context, options *KubectlOptions, filters metav1.ListOptions) ([]corev1.Service, error) {
	clientset, err := GetKubernetesClientFromOptionsContextE(t, ctx, options)
	if err != nil {
		return nil, err
	}

	resp, err := clientset.CoreV1().Services(options.Namespace).List(ctx, filters)
	if err != nil {
		return nil, err
	}

	return resp.Items, nil
}

// ListServicesContext looks up services in the given namespace that match the given filters and return them.
// The ctx parameter supports cancellation and timeouts.
// This will fail the test if there is an error.
//
//nolint:gocritic // hugeParam: cannot change public function signature
func ListServicesContext(t testing.TestingT, ctx context.Context, options *KubectlOptions, filters metav1.ListOptions) []corev1.Service {
	t.Helper()
	services, err := ListServicesContextE(t, ctx, options, filters)
	require.NoError(t, err)

	return services
}

// GetServiceContextE returns a Kubernetes service resource in the provided namespace with the given name.
// The ctx parameter supports cancellation and timeouts.
func GetServiceContextE(t testing.TestingT, ctx context.Context, options *KubectlOptions, serviceName string) (*corev1.Service, error) {
	clientset, err := GetKubernetesClientFromOptionsContextE(t, ctx, options)
	if err != nil {
		return nil, err
	}

	return clientset.CoreV1().Services(options.Namespace).Get(ctx, serviceName, metav1.GetOptions{})
}

// GetServiceContext returns a Kubernetes service resource in the provided namespace with the given name.
// The ctx parameter supports cancellation and timeouts.
// This will fail the test if there is an error.
func GetServiceContext(t testing.TestingT, ctx context.Context, options *KubectlOptions, serviceName string) *corev1.Service {
	t.Helper()
	service, err := GetServiceContextE(t, ctx, options, serviceName)
	require.NoError(t, err)

	return service
}

// WaitUntilServiceAvailableContextE waits until the service endpoint is ready to accept traffic.
// The ctx parameter supports cancellation and timeouts.
func WaitUntilServiceAvailableContextE(t testing.TestingT, ctx context.Context, options *KubectlOptions, serviceName string, retries int, sleepBetweenRetries time.Duration) error {
	statusMsg := fmt.Sprintf("Wait for service %s to be provisioned.", serviceName)

	message, err := retry.DoWithRetryContextE(
		t,
		ctx,
		statusMsg,
		retries,
		sleepBetweenRetries,
		func() (string, error) {
			service, err := GetServiceContextE(t, ctx, options, serviceName)
			if err != nil {
				return "", err
			}

			isMinikube, err := IsMinikubeE(t, options)
			if err != nil {
				return "", err
			}

			if !isMinikube && !IsServiceAvailable(service) {
				return "", NewServiceNotAvailableError(service)
			}

			return "Service is now available", nil
		},
	)
	if err != nil {
		return err
	}

	options.Logger.Logf(t, "%s", message)

	return nil
}

// WaitUntilServiceAvailableContext waits until the service endpoint is ready to accept traffic.
// The ctx parameter supports cancellation and timeouts.
// This will fail the test if there is an error.
func WaitUntilServiceAvailableContext(t testing.TestingT, ctx context.Context, options *KubectlOptions, serviceName string, retries int, sleepBetweenRetries time.Duration) {
	t.Helper()
	err := WaitUntilServiceAvailableContextE(t, ctx, options, serviceName, retries, sleepBetweenRetries)
	require.NoError(t, err)
}

// IsServiceAvailable returns true if the service endpoint is ready to accept traffic. Note that for Minikube, this
// function is moot as all services, even LoadBalancer, is available immediately.
func IsServiceAvailable(service *corev1.Service) bool {

	switch service.Spec.Type {
	case corev1.ServiceTypeLoadBalancer:
		ingress := service.Status.LoadBalancer.Ingress

		return len(ingress) > 0
	case corev1.ServiceTypeClusterIP, corev1.ServiceTypeNodePort, corev1.ServiceTypeExternalName:
		return true
	default:
		return true
	}
}

// GetServiceEndpointContext will return the service access point using the provided context. If the service endpoint is
// not ready, will fail the test immediately.
func GetServiceEndpointContext(t testing.TestingT, ctx context.Context, options *KubectlOptions, service *corev1.Service, servicePort int) string {
	t.Helper()
	endpoint, err := GetServiceEndpointContextE(t, ctx, options, service, servicePort)
	require.NoError(t, err)

	return endpoint
}

// GetServiceEndpointContextE will return the service access point using the provided context and the following logic:
//   - For ClusterIP service type, return the URL that maps to ClusterIP and Service Port
//   - For NodePort service type, identify the public IP of the node (if it exists, otherwise return the bound hostname),
//     and the assigned node port for the provided service port, and return the URL that maps to node ip and node port.
//   - For LoadBalancer service type, return the publicly accessible hostname of the load balancer.
//     If the hostname is empty, it will return the public IP of the LoadBalancer.
//   - All other service types are not supported.
func GetServiceEndpointContextE(t testing.TestingT, ctx context.Context, options *KubectlOptions, service *corev1.Service, servicePort int) (string, error) {
	switch service.Spec.Type {
	case corev1.ServiceTypeClusterIP:

		return net.JoinHostPort(service.Spec.ClusterIP, strconv.Itoa(servicePort)), nil
	case corev1.ServiceTypeNodePort:
		return findEndpointForNodePortServiceContext(t, ctx, options, service, int32(servicePort))
	case corev1.ServiceTypeExternalName:
		return "", NewUnknownServiceTypeError(service)
	case corev1.ServiceTypeLoadBalancer:

		isMinikube, err := IsMinikubeE(t, options)
		if err != nil {
			return "", err
		}

		if isMinikube {
			return findEndpointForNodePortServiceContext(t, ctx, options, service, int32(servicePort))
		}

		ingress := service.Status.LoadBalancer.Ingress
		if len(ingress) == 0 {
			return "", NewServiceNotAvailableError(service)
		}

		if ingress[0].Hostname == "" {
			return net.JoinHostPort(ingress[0].IP, strconv.Itoa(servicePort)), nil
		}

		return net.JoinHostPort(ingress[0].Hostname, strconv.Itoa(servicePort)), nil
	default:
		return "", NewUnknownServiceTypeError(service)
	}
}

// findEndpointForNodePortServiceContext extracts an endpoint that can be reached outside the kubernetes cluster using the
// provided context. NodePort type needs to find the right allocated node port mapped to the service port, as well as
// find out the externally reachable ip (if available).
func findEndpointForNodePortServiceContext(
	t testing.TestingT,
	ctx context.Context,
	options *KubectlOptions,
	service *corev1.Service,
	servicePort int32,
) (string, error) {
	nodePort, err := FindNodePortContextE(ctx, service, servicePort)
	if err != nil {
		return "", err
	}

	node, err := pickRandomNodeE(t, options)
	if err != nil {
		return "", err
	}

	nodeHostname, err := FindNodeHostnameWithOptionsContextE(t, ctx, options, node)
	if err != nil {
		return "", err
	}

	return net.JoinHostPort(nodeHostname, strconv.FormatInt(int64(nodePort), 10)), nil
}

// FindNodePortContextE returns the allocated NodePort for the given servicePort from the service definition.
// The ctx parameter is accepted for API consistency but is not used since this is a local struct lookup.
func FindNodePortContextE(_ context.Context, service *corev1.Service, servicePort int32) (int32, error) {
	for _, port := range service.Spec.Ports {
		if port.Port == servicePort {
			return port.NodePort, nil
		}
	}

	return -1, NewUnknownServicePortError(service, servicePort)
}

// FindNodePortContext returns the allocated NodePort for the given servicePort from the service definition.
// The ctx parameter is accepted for API consistency but is not used since this is a local struct lookup.
// This will fail the test if there is an error.
func FindNodePortContext(t testing.TestingT, ctx context.Context, service *corev1.Service, servicePort int32) int32 {
	t.Helper()

	nodePort, err := FindNodePortContextE(ctx, service, servicePort)
	require.NoError(t, err)

	return nodePort
}

// pickRandomNode will pick a random node in the kubernetes cluster
func pickRandomNodeE(t testing.TestingT, options *KubectlOptions) (corev1.Node, error) {
	nodes, err := GetNodesContextE(t, context.Background(), options)
	if err != nil {
		return corev1.Node{}, err
	}

	if len(nodes) == 0 {
		return corev1.Node{}, NewNoNodesInKubernetesError()
	}

	index := random.Random(0, len(nodes)-1)

	return nodes[index], nil
}

// FindNodeHostnameContext returns the hostname or IP address of the given node using the provided context, preferring
// the external IP when available. This will fail the test if there is an error.
//
//nolint:gocritic // hugeParam: cannot change public function signature
func FindNodeHostnameContext(t testing.TestingT, ctx context.Context, node corev1.Node) string {
	t.Helper()
	hostname, err := FindNodeHostnameContextE(t, ctx, node)
	require.NoError(t, err)

	return hostname
}

// FindNodeHostnameContextE returns the hostname or IP address of the given node using the provided context, preferring
// the external IP when available.
//
// The address is read from the node object itself: cloud controller managers record an instance's public IP as an
// ExternalIP address on the node. For the rare cluster that does not advertise one, see
// FindNodeHostnameWithOptionsContextE, which can fall back to querying the cloud provider directly.
//
//nolint:gocritic // hugeParam: cannot change public function signature
func FindNodeHostnameContextE(t testing.TestingT, ctx context.Context, node corev1.Node) (string, error) {
	return FindNodeHostnameWithOptionsContextE(t, ctx, nil, node)
}

// FindNodeHostnameWithOptionsContext behaves like FindNodeHostnameContext, but consults
// options.NodePublicIPLookup when the node itself does not advertise an external IP. This will fail the test if
// there is an error.
//
//nolint:gocritic // hugeParam: cannot change public function signature
func FindNodeHostnameWithOptionsContext(
	t testing.TestingT,
	ctx context.Context,
	options *KubectlOptions,
	node corev1.Node,
) string {
	t.Helper()
	hostname, err := FindNodeHostnameWithOptionsContextE(t, ctx, options, node)
	require.NoError(t, err)

	return hostname
}

// FindNodeHostnameWithOptionsContextE behaves like FindNodeHostnameContextE, but consults
// options.NodePublicIPLookup when the node itself does not advertise an external IP.
//
// Almost all callers want FindNodeHostnameContextE instead. A lookup is only needed on clusters whose cloud
// controller manager does not record the instance's public IP as an ExternalIP on the node object. See
// NodePublicIPLookup for the wiring.
//
//nolint:gocritic // hugeParam: cannot change public function signature
func FindNodeHostnameWithOptionsContextE(
	t testing.TestingT,
	ctx context.Context,
	options *KubectlOptions,
	node corev1.Node,
) (string, error) {
	// An external IP recorded on the node is authoritative and costs no API call, so prefer it for every provider.
	if externalIP, ok := findNodeAddress(&node, corev1.NodeExternalIP); ok {
		return externalIP, nil
	}

	nodeIDUri, err := url.Parse(node.Spec.ProviderID)
	if err != nil {
		return "", err
	}

	switch nodeIDUri.Scheme {
	case awsProviderIDScheme:
		return findAwsNodeHostnameContextE(t, ctx, options, &node, nodeIDUri)
	default:
		return findDefaultNodeHostnameE(&node)
	}
}

// findNodeAddress returns the first address of the given type recorded on the node, and whether one was found.
func findNodeAddress(node *corev1.Node, addressType corev1.NodeAddressType) (string, bool) {
	for _, address := range node.Status.Addresses {
		if address.Type == addressType && address.Address != "" {
			return address.Address, true
		}
	}

	return "", false
}

// findAwsNodeHostname will return the public ip of the node, assuming the node is an AWS EC2 instance.
// If the instance does not have a public IP, will return the internal hostname as recorded on the Kubernetes node
// object.
// expectedAWSIDPathParts is the number of path segments in an AWS provider ID (empty, availability zone, instance ID).
const expectedAWSIDPathParts = 3

func findAwsNodeHostnameContextE(
	t testing.TestingT,
	ctx context.Context,
	options *KubectlOptions,
	node *corev1.Node,
	awsIDUri *url.URL,
) (string, error) {
	parts := strings.Split(awsIDUri.Path, "/")
	if len(parts) != expectedAWSIDPathParts {
		return "", NewMalformedNodeIDError(node)
	}

	instanceID := parts[2]

	availabilityZone := parts[1]
	if availabilityZone == "" || instanceID == "" {
		return "", NewMalformedNodeIDError(node)
	}

	// An AZ is a region plus a single trailing letter, so dropping that letter yields the region.
	region := availabilityZone[:len(availabilityZone)-1]

	if options == nil || options.NodePublicIPLookup == nil {
		logger.Default.Logf(t, "[WARNING] Node %s is backed by an AWS EC2 instance, but KubectlOptions.NodePublicIPLookup "+
			"is not set, so its public IP cannot be resolved. Falling back to the internal hostname recorded on the node. "+
			"Set options.NodePublicIPLookup = aws.GetPublicIpsOfEc2InstancesContextE to resolve public IPs.", node.Name)

		return findDefaultNodeHostnameE(node)
	}

	ipMap, err := options.NodePublicIPLookup(t, ctx, []string{instanceID}, region)
	if err != nil {
		return "", err
	}

	publicIP, containsIP := ipMap[instanceID]
	if !containsIP || publicIP == "" {

		return findDefaultNodeHostnameE(node)
	}

	return publicIP, nil
}

// findDefaultNodeHostname returns the hostname recorded on the Kubernetes node object.
func findDefaultNodeHostnameE(node *corev1.Node) (string, error) {
	if hostname, ok := findNodeAddress(node, corev1.NodeHostName); ok {
		return hostname, nil
	}

	return "", NewNodeHasNoHostnameError(node)
}
