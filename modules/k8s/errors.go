package k8s

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	networkingv1beta1 "k8s.io/api/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// IngressNotAvailable is returned when a Kubernetes service is not yet available to accept traffic.
type IngressNotAvailable struct {
	ingress *networkingv1.Ingress
}

// Error is a simple function to return a formatted error message as a string
func (err IngressNotAvailable) Error() string {
	return fmt.Sprintf("Ingress %s is not available", err.ingress.Name)
}

// IngressNotAvailableV1Beta1 is returned when a Kubernetes service is not yet available to accept traffic.
type IngressNotAvailableV1Beta1 struct {
	ingress *networkingv1beta1.Ingress
}

// Error is a simple function to return a formatted error message as a string
func (err IngressNotAvailableV1Beta1) Error() string {
	return fmt.Sprintf("Ingress %s is not available", err.ingress.Name)
}

// UnknownKubeResourceType is returned if the given resource type does not match the list of known resource types.
type UnknownKubeResourceType struct {
	ResourceType KubeResourceType
}

func (err UnknownKubeResourceType) Error() string {
	return fmt.Sprintf("ResourceType ID %d is unknown", err.ResourceType)
}

// DesiredNumberOfPodsNotCreated is returned when the number of pods matching a filter condition does not match the
// desired number of Pods.
type DesiredNumberOfPodsNotCreated struct {
	Filter       metav1.ListOptions
	DesiredCount int
}

// Error is a simple function to return a formatted error message as a string
func (err DesiredNumberOfPodsNotCreated) Error() string {
	return fmt.Sprintf("Desired number of pods (%d) matching filter %v not yet created", err.DesiredCount, err.Filter)
}

// ServiceAccountTokenNotAvailable is returned when a Kubernetes ServiceAccount does not have a token provisioned yet.
type ServiceAccountTokenNotAvailable struct {
	Name string
}

// Error is a simple function to return a formatted error message as a string
func (err ServiceAccountTokenNotAvailable) Error() string {
	return fmt.Sprintf("ServiceAccount %s does not have a token yet.", err.Name)
}

// DeploymentNotAvailable is returned when a Kubernetes deployment is not yet available to accept traffic.
type DeploymentNotAvailable struct {
	deploy *appsv1.Deployment
}

// Error is a simple function to return a formatted error message as a string
func (err DeploymentNotAvailable) Error() string {
	dc := getDeploymentCondition(err.deploy, appsv1.DeploymentProgressing)
	if dc == nil {
		return fmt.Sprintf(
			"Deployment %s is not available, missing '%s' condition",
			err.deploy.Name,
			appsv1.DeploymentProgressing,
		)
	}
	return fmt.Sprintf(
		"Deployment %s is not available as '%s' condition indicates that the Deployment is not complete, status: %v, reason: %s, message: %s",
		err.deploy.Name,
		appsv1.DeploymentProgressing,
		dc.Status,
		dc.Reason,
		dc.Message,
	)
}

// NewDeploymentNotAvailableError returnes a DeploymentNotAvailable struct when Kubernetes deems a deployment is not available
func NewDeploymentNotAvailableError(deploy *appsv1.Deployment) DeploymentNotAvailable {
	return DeploymentNotAvailable{deploy}
}

// PodNotAvailable is returned when a Kubernetes service is not yet available to accept traffic.
type PodNotAvailable struct {
	pod *corev1.Pod
}

// Error is a simple function to return a formatted error message as a string
func (err PodNotAvailable) Error() string {
	return fmt.Sprintf("Pod %s is not available, reason: %s, message: %s", err.pod.Name, err.pod.Status.Reason, err.pod.Status.Message)
}

// NewPodNotAvailableError returnes a PodNotAvailable struct when Kubernetes deems a pod is not available
func NewPodNotAvailableError(pod *corev1.Pod) PodNotAvailable {
	return PodNotAvailable{pod}
}

// JobNotSucceeded is returned when a Kubernetes job is not Succeeded
type JobNotSucceeded struct {
	job *batchv1.Job
}

// Error is a simple function to return a formatted error message as a string
func (err JobNotSucceeded) Error() string {
	return fmt.Sprintf("Job %s is not Succeeded", err.job.Name)
}

// NewJobNotSucceeded returnes a JobNotSucceeded when the status of the job is not Succeeded
func NewJobNotSucceeded(job *batchv1.Job) JobNotSucceeded {
	return JobNotSucceeded{job}
}

// ServiceNotAvailable is returned when a Kubernetes service is not yet available to accept traffic.
type ServiceNotAvailable struct {
	service *corev1.Service
}

// Error is a simple function to return a formatted error message as a string
func (err ServiceNotAvailable) Error() string {
	return fmt.Sprintf("Service %s is not available", err.service.Name)
}

// NewServiceNotAvailableError returnes a ServiceNotAvailable struct when Kubernetes deems a service is not available
func NewServiceNotAvailableError(service *corev1.Service) ServiceNotAvailable {
	return ServiceNotAvailable{service}
}

// UnknownServiceType is returned when a Kubernetes service has a type that is not yet handled by the test functions.
type UnknownServiceType struct {
	service *corev1.Service
}

// Error is a simple function to return a formatted error message as a string
func (err UnknownServiceType) Error() string {
	return fmt.Sprintf("Service %s has an unknown service type", err.service.Name)
}

// NewUnknownServiceTypeError returns an UnknownServiceType struct when is it deemed that Kubernetes does not know the service type provided
func NewUnknownServiceTypeError(service *corev1.Service) UnknownServiceType {
	return UnknownServiceType{service}
}

// UnknownServicePort is returned when the given service port is not an exported port of the service.
type UnknownServicePort struct {
	service *corev1.Service
	port    int32
}

// Error is a simple function to return a formatted error message as a string
func (err UnknownServicePort) Error() string {
	return fmt.Sprintf("Port %d is not a part of the service %s", err.port, err.service.Name)
}

// NewUnknownServicePortError returns an UnknownServicePort struct when it is deemed that Kuberenetes does not know of the provided Service Port
func NewUnknownServicePortError(service *corev1.Service, port int32) UnknownServicePort {
	return UnknownServicePort{service, port}
}

// PersistentVolumeNotInStatus is returned when a Kubernetes PersistentVolume is not in the expected status phase
type PersistentVolumeNotInStatus struct {
	pv            *corev1.PersistentVolume
	pvStatusPhase *corev1.PersistentVolumePhase
}

// Error is a simple function to return a formatted error message as a string
func (err PersistentVolumeNotInStatus) Error() string {
	return fmt.Sprintf("Pv %s is not '%s'", err.pv.Name, *err.pvStatusPhase)
}

// NewPersistentVolumeNotInStatusError returns a PersistentVolumeNotInStatus struct when the given Persistent Volume is not in the expected status phase
func NewPersistentVolumeNotInStatusError(pv *corev1.PersistentVolume, pvStatusPhase *corev1.PersistentVolumePhase) PersistentVolumeNotInStatus {
	return PersistentVolumeNotInStatus{pv, pvStatusPhase}
}

// PersistentVolumeClaimNotInStatus is returned when a Kubernetes PersistentVolumeClaim is not in the expected status phase
type PersistentVolumeClaimNotInStatus struct {
	pvc            *corev1.PersistentVolumeClaim
	pvcStatusPhase *corev1.PersistentVolumeClaimPhase
}

// Error is a simple function to return a formatted error message as a string
func (err PersistentVolumeClaimNotInStatus) Error() string {
	return fmt.Sprintf("PVC %s is not '%s'", err.pvc.Name, *err.pvcStatusPhase)
}

// NewPersistentVolumeClaimNotInStatusError returns a PersistentVolumeClaimNotInStatus struct when the given PersistentVolumeClaim is not in the expected status phase
func NewPersistentVolumeClaimNotInStatusError(pvc *corev1.PersistentVolumeClaim, pvcStatusPhase *corev1.PersistentVolumeClaimPhase) PersistentVolumeClaimNotInStatus {
	return PersistentVolumeClaimNotInStatus{pvc, pvcStatusPhase}
}

// NoNodesInKubernetes is returned when the Kubernetes cluster has no nodes registered.
type NoNodesInKubernetes struct{}

// Error is a simple function to return a formatted error message as a string
func (err NoNodesInKubernetes) Error() string {
	return "There are no nodes in the Kubernetes cluster"
}

// NewNoNodesInKubernetesError returns a NoNodesInKubernetes struct when it is deemed that there are no Kubernetes nodes registered
func NewNoNodesInKubernetesError() NoNodesInKubernetes {
	return NoNodesInKubernetes{}
}

// NodeHasNoHostname is returned when a Kubernetes node has no discernible hostname
type NodeHasNoHostname struct {
	node *corev1.Node
}

// Error is a simple function to return a formatted error message as a string
func (err NodeHasNoHostname) Error() string {
	return fmt.Sprintf("Node %s has no hostname", err.node.Name)
}

// NewNodeHasNoHostnameError returns a NodeHasNoHostname struct when it is deemed that the provided node has no hostname
func NewNodeHasNoHostnameError(node *corev1.Node) NodeHasNoHostname {
	return NodeHasNoHostname{node}
}

// MalformedNodeID is returned when a Kubernetes node has a malformed node id scheme
type MalformedNodeID struct {
	node *corev1.Node
}

// Error is a simple function to return a formatted error message as a string
func (err MalformedNodeID) Error() string {
	return fmt.Sprintf("Node %s has malformed ID %s", err.node.Name, err.node.Spec.ProviderID)
}

// NewMalformedNodeIDError returns a MalformedNodeID struct when Kubernetes deems that a NodeID is malformed
func NewMalformedNodeIDError(node *corev1.Node) MalformedNodeID {
	return MalformedNodeID{node}
}

// JSONPathMalformedJSONErr is returned when the jsonpath unmarshal routine fails to parse the given JSON blob.
type JSONPathMalformedJSONErr struct {
	underlyingErr error
}

func (err JSONPathMalformedJSONErr) Error() string {
	return fmt.Sprintf("Error unmarshaling original json blob: %s", err.underlyingErr)
}

// JSONPathMalformedJSONPathErr is returned when the jsonpath unmarshal routine fails to parse the given JSON path
// string.
type JSONPathMalformedJSONPathErr struct {
	underlyingErr error
}

func (err JSONPathMalformedJSONPathErr) Error() string {
	return fmt.Sprintf("Error parsing json path: %s", err.underlyingErr)
}

// JSONPathExtractJSONPathErr is returned when the jsonpath unmarshal routine fails to extract the given JSON path from
// the JSON blob.
type JSONPathExtractJSONPathErr struct {
	underlyingErr error
}

func (err JSONPathExtractJSONPathErr) Error() string {
	return fmt.Sprintf("Error extracting json path from blob: %s", err.underlyingErr)
}

// JSONPathMalformedJSONPathResultErr is returned when the jsonpath unmarshal routine fails to unmarshal the resulting
// data from extraction.
type JSONPathMalformedJSONPathResultErr struct {
	underlyingErr error
}

func (err JSONPathMalformedJSONPathResultErr) Error() string {
	return fmt.Sprintf("Error unmarshaling json path output: %s", err.underlyingErr)
}
