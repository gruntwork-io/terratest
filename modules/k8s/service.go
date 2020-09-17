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
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/retry"
	"github.com/gruntwork-io/terratest/modules/testing"
)

// ListServices will look for services in the given namespace that match the given filters and return them. This will
// fail the test if there is an error.
func ListServices(t testing.TestingT, options *KubectlOptions, filters metav1.ListOptions) []corev1.Service {
	service, err := ListServicesE(t, options, filters)
	require.NoError(t, err)
	return service
}

// ListServicesE will look for services in the given namespace that match the given filters and return them.
func ListServicesE(t testing.TestingT, options *KubectlOptions, filters metav1.ListOptions) ([]corev1.Service, error) {
	clientset, err := GetKubernetesClientFromOptionsE(t, options)
	if err != nil {
		return nil, err
	}
	resp, err := clientset.CoreV1().Services(options.Namespace).List(context.Background(), filters)
	if err != nil {
		return nil, err
	}
	return resp.Items, nil
}

// GetService returns a Kubernetes service resource in the provided namespace with the given name. This will
// fail the test if there is an error.
func GetService(t testing.TestingT, options *KubectlOptions, serviceName string) *corev1.Service {
	service, err := GetServiceE(t, options, serviceName)
	require.NoError(t, err)
	return service
}

// GetServiceE returns a Kubernetes service resource in the provided namespace with the given name.
func GetServiceE(t testing.TestingT, options *KubectlOptions, serviceName string) (*corev1.Service, error) {
	clientset, err := GetKubernetesClientFromOptionsE(t, options)
	if err != nil {
		return nil, err
	}
	return clientset.CoreV1().Services(options.Namespace).Get(context.Background(), serviceName, metav1.GetOptions{})
}

// WaitUntilServiceAvailable waits until the service endpoint is ready to accept traffic.
func WaitUntilServiceAvailable(t testing.TestingT, options *KubectlOptions, serviceName string, retries int, sleepBetweenRetries time.Duration) {
	statusMsg := fmt.Sprintf("Wait for service %s to be provisioned.", serviceName)
	message := retry.DoWithRetry(
		t,
		statusMsg,
		retries,
		sleepBetweenRetries,
		func() (string, error) {
			service, err := GetServiceE(t, options, serviceName)
			if err != nil {
				return "", err
			}

			isMinikube, err := IsMinikubeE(t, options)
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
	logger.Logf(t, message)
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
	default:
		return true
	}
}

// GetServiceEndpoint will return the service access point. If the service endpoint is not ready, will fail the test
// immediately.
func GetServiceEndpoint(t testing.TestingT, options *KubectlOptions, service *corev1.Service, servicePort int) string {
	endpoint, err := GetServiceEndpointE(t, options, service, servicePort)
	require.NoError(t, err)
	return endpoint
}

// GetServiceEndpointE will return the service access point using the following logic:
// - For ClusterIP service type, return the URL that maps to ClusterIP and Service Port
// - For NodePort service type, identify the public IP of the node (if it exists, otherwise return the bound hostname),
//   and the assigned node port for the provided service port, and return the URL that maps to node ip and node port.
// - For LoadBalancer service type, return the publicly accessible hostname of the load balancer.
//   If the hostname is empty, it will return the public IP of the LoadBalancer.
// - All other service types are not supported.
func GetServiceEndpointE(t testing.TestingT, options *KubectlOptions, service *corev1.Service, servicePort int) (string, error) {
	switch service.Spec.Type {
	case corev1.ServiceTypeClusterIP:
		// ClusterIP service type will map directly to service port
		return fmt.Sprintf("%s:%d", service.Spec.ClusterIP, servicePort), nil
	case corev1.ServiceTypeNodePort:
		return findEndpointForNodePortService(t, options, service, int32(servicePort))
	case corev1.ServiceTypeLoadBalancer:
		// For minikube, LoadBalancer service is exactly the same as NodePort service
		isMinikube, err := IsMinikubeE(t, options)
		if err != nil {
			return "", err
		}
		if isMinikube {
			return findEndpointForNodePortService(t, options, service, int32(servicePort))
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

// Extracts a endpoint that can be reached outside the kubernetes cluster. NodePort type needs to find the right
// allocated node port mapped to the service port, as well as find out the externally reachable ip (if available).
func findEndpointForNodePortService(
	t testing.TestingT,
	options *KubectlOptions,
	service *corev1.Service,
	servicePort int32,
) (string, error) {
	nodePort, err := FindNodePortE(service, int32(servicePort))
	if err != nil {
		return "", err
	}
	node, err := pickRandomNodeE(t, options)
	if err != nil {
		return "", err
	}
	nodeHostname, err := FindNodeHostnameE(t, node)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s:%d", nodeHostname, nodePort), nil
}

// Given the desired servicePort, return the allocated nodeport
func FindNodePortE(service *corev1.Service, servicePort int32) (int32, error) {
	for _, port := range service.Spec.Ports {
		if port.Port == servicePort {
			return port.NodePort, nil
		}
	}
	return -1, NewUnknownServicePortError(service, servicePort)
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

// Given a node, return the ip address, preferring the external IP
func FindNodeHostnameE(t testing.TestingT, node corev1.Node) (string, error) {
	nodeIDUri, err := url.Parse(node.Spec.ProviderID)
	if err != nil {
		return "", err
	}
	switch nodeIDUri.Scheme {
	case "aws":
		return findAwsNodeHostnameE(t, node, nodeIDUri)
	default:
		return findDefaultNodeHostnameE(node)
	}
}

// findAwsNodeHostname will return the public ip of the node, assuming the node is an AWS EC2 instance.
// If the instance does not have a public IP, will return the internal hostname as recorded on the Kubernetes node
// object.
func findAwsNodeHostnameE(t testing.TestingT, node corev1.Node, awsIDUri *url.URL) (string, error) {
	// Path is /AVAILABILITY_ZONE/INSTANCE_ID
	parts := strings.Split(awsIDUri.Path, "/")
	if len(parts) != 3 {
		return "", NewMalformedNodeIDError(&node)
	}
	instanceID := parts[2]
	availabilityZone := parts[1]
	// Availability Zone name is known to be region code + 1 letter
	// https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/using-regions-availability-zones.html
	region := availabilityZone[:len(availabilityZone)-1]

	ipMap, err := aws.GetPublicIpsOfEc2InstancesE(t, []string{instanceID}, region)
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
func findDefaultNodeHostnameE(node corev1.Node) (string, error) {
	for _, address := range node.Status.Addresses {
		if address.Type == corev1.NodeHostName {
			return address.Address, nil
		}
	}
	return "", NewNodeHasNoHostnameError(&node)
}
