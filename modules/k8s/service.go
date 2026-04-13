package k8s

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/retry"
	"github.com/gruntwork-io/terratest/modules/testing"
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

// ListServices will look for services in the given namespace that match the given filters and return them. This will
// fail the test if there is an error.
//
// Deprecated: Use [ListServicesContext] instead.
//
//nolint:gocritic // hugeParam: cannot change public function signature
func ListServices(t testing.TestingT, options *KubectlOptions, filters metav1.ListOptions) []corev1.Service {
	t.Helper()

	return ListServicesContext(t, context.Background(), options, filters)
}

// ListServicesE will look for services in the given namespace that match the given filters and return them.
//
// Deprecated: Use [ListServicesContextE] instead.
//
//nolint:gocritic // hugeParam: cannot change public function signature
func ListServicesE(t testing.TestingT, options *KubectlOptions, filters metav1.ListOptions) ([]corev1.Service, error) {
	return ListServicesContextE(t, context.Background(), options, filters)
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

// GetService returns a Kubernetes service resource in the provided namespace with the given name. This will
// fail the test if there is an error.
//
// Deprecated: Use [GetServiceContext] instead.
func GetService(t testing.TestingT, options *KubectlOptions, serviceName string) *corev1.Service {
	t.Helper()

	return GetServiceContext(t, context.Background(), options, serviceName)
}

// GetServiceE returns a Kubernetes service resource in the provided namespace with the given name.
//
// Deprecated: Use [GetServiceContextE] instead.
func GetServiceE(t testing.TestingT, options *KubectlOptions, serviceName string) (*corev1.Service, error) {
	return GetServiceContextE(t, context.Background(), options, serviceName)
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

			isMinikube, err := IsMinikubeE(t, options) //nolint:contextcheck // IsMinikubeE not yet context-aware
			if err != nil {
				return "", err
			}

			// For minikube, all services will be available immediately so we only do the check if we are not on
			// minikube.
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

// WaitUntilServiceAvailable waits until the service endpoint is ready to accept traffic.
//
// Deprecated: Use [WaitUntilServiceAvailableContext] instead.
func WaitUntilServiceAvailable(t testing.TestingT, options *KubectlOptions, serviceName string, retries int, sleepBetweenRetries time.Duration) {
	t.Helper()
	WaitUntilServiceAvailableContext(t, context.Background(), options, serviceName, retries, sleepBetweenRetries)
}

// IsServiceAvailable returns true if the service endpoint is ready to accept traffic. Note that for Minikube, this
// function is moot as all services, even LoadBalancer, is available immediately.
func IsServiceAvailable(service *corev1.Service) bool {
	// Only the LoadBalancer type has a delay. All other service types are available if the resource exists.
	switch service.Spec.Type {
	case corev1.ServiceTypeLoadBalancer:
		ingress := service.Status.LoadBalancer.Ingress
		// The load balancer is ready if it has at least one ingress point
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

// GetServiceEndpoint will return the service access point. If the service endpoint is not ready, will fail the test
// immediately.
//
// Deprecated: Use [GetServiceEndpointContext] instead.
func GetServiceEndpoint(t testing.TestingT, options *KubectlOptions, service *corev1.Service, servicePort int) string {
	t.Helper()

	return GetServiceEndpointContext(t, context.Background(), options, service, servicePort)
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
		// ClusterIP service type will map directly to service port
		return fmt.Sprintf("%s:%d", service.Spec.ClusterIP, servicePort), nil
	case corev1.ServiceTypeNodePort:
		return findEndpointForNodePortServiceContext(t, ctx, options, service, int32(servicePort))
	case corev1.ServiceTypeExternalName:
		return "", NewUnknownServiceTypeError(service)
	case corev1.ServiceTypeLoadBalancer:
		// For minikube, LoadBalancer service is exactly the same as NodePort service
		isMinikube, err := IsMinikubeE(t, options) //nolint:contextcheck // IsMinikubeE not yet context-aware
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
			return fmt.Sprintf("%s:%d", ingress[0].IP, servicePort), nil
		}
		// Load Balancer service type will map directly to service port
		return fmt.Sprintf("%s:%d", ingress[0].Hostname, servicePort), nil
	default:
		return "", NewUnknownServiceTypeError(service)
	}
}

// GetServiceEndpointE will return the service access point using the following logic:
//   - For ClusterIP service type, return the URL that maps to ClusterIP and Service Port
//   - For NodePort service type, identify the public IP of the node (if it exists, otherwise return the bound hostname),
//     and the assigned node port for the provided service port, and return the URL that maps to node ip and node port.
//   - For LoadBalancer service type, return the publicly accessible hostname of the load balancer.
//     If the hostname is empty, it will return the public IP of the LoadBalancer.
//   - All other service types are not supported.
//
// Deprecated: Use [GetServiceEndpointContextE] instead.
func GetServiceEndpointE(t testing.TestingT, options *KubectlOptions, service *corev1.Service, servicePort int) (string, error) {
	return GetServiceEndpointContextE(t, context.Background(), options, service, servicePort)
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

	node, err := pickRandomNodeE(t, options) //nolint:contextcheck // pickRandomNodeE not yet context-aware
	if err != nil {
		return "", err
	}

	nodeHostname, err := FindNodeHostnameContextE(t, ctx, node)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s:%d", nodeHostname, nodePort), nil
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

// FindNodePortE returns the allocated NodePort for the given servicePort from the service definition.
//
// Deprecated: Use [FindNodePortContextE] instead.
func FindNodePortE(service *corev1.Service, servicePort int32) (int32, error) {
	return FindNodePortContextE(context.Background(), service, servicePort)
}

// pickRandomNode will pick a random node in the kubernetes cluster
func pickRandomNodeE(t testing.TestingT, options *KubectlOptions) (corev1.Node, error) {
	nodes, err := GetNodesE(t, options)
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
//nolint:gocritic // hugeParam: cannot change public function signature
func FindNodeHostnameContextE(t testing.TestingT, ctx context.Context, node corev1.Node) (string, error) {
	nodeIDUri, err := url.Parse(node.Spec.ProviderID)
	if err != nil {
		return "", err
	}

	switch nodeIDUri.Scheme {
	case "aws":
		return findAwsNodeHostnameContextE(t, ctx, &node, nodeIDUri)
	default:
		return findDefaultNodeHostnameE(&node)
	}
}

// FindNodeHostnameE returns the hostname or IP address of the given node, preferring the external IP when available.
//
// Deprecated: Use [FindNodeHostnameContextE] instead.
//
//nolint:gocritic // hugeParam: cannot change public function signature
func FindNodeHostnameE(t testing.TestingT, node corev1.Node) (string, error) {
	t.Helper()

	return FindNodeHostnameContextE(t, context.Background(), node)
}

// findAwsNodeHostname will return the public ip of the node, assuming the node is an AWS EC2 instance.
// If the instance does not have a public IP, will return the internal hostname as recorded on the Kubernetes node
// object.
// expectedAWSIDPathParts is the number of path segments in an AWS provider ID (empty, availability zone, instance ID).
const expectedAWSIDPathParts = 3

func findAwsNodeHostnameContextE(t testing.TestingT, ctx context.Context, node *corev1.Node, awsIDUri *url.URL) (string, error) {
	// Path is /AVAILABILITY_ZONE/INSTANCE_ID
	parts := strings.Split(awsIDUri.Path, "/")
	if len(parts) != expectedAWSIDPathParts {
		return "", NewMalformedNodeIDError(node)
	}

	instanceID := parts[2]
	availabilityZone := parts[1]
	// Availability Zone name is known to be region code + 1 letter
	// https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/using-regions-availability-zones.html
	region := availabilityZone[:len(availabilityZone)-1]

	ipMap, err := aws.GetPublicIpsOfEc2InstancesContextE(t, ctx, []string{instanceID}, region)
	if err != nil {
		return "", err
	}

	publicIP, containsIP := ipMap[instanceID]
	if !containsIP || publicIP == "" {
		// return default hostname
		return findDefaultNodeHostnameE(node)
	}

	return publicIP, nil
}

// findDefaultNodeHostname returns the hostname recorded on the Kubernetes node object.
func findDefaultNodeHostnameE(node *corev1.Node) (string, error) {
	for _, address := range node.Status.Addresses {
		if address.Type == corev1.NodeHostName {
			return address.Address, nil
		}
	}

	return "", NewNodeHasNoHostnameError(node)
}
