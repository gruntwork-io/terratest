package k8s

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"k8s.io/client-go/rest"
)

// awsProviderIDScheme is the scheme used in the provider ID of a node backed by an AWS EC2 instance.
const awsProviderIDScheme = "aws"

// NodePublicIPLookup resolves the public IP addresses of the given cloud provider instance IDs in the given region,
// returning a map of instance ID to public IP. Instances with no public IP may be omitted from the map.
//
// This is an escape hatch, not the usual path. Node addresses recorded by a cloud controller manager already carry
// the instance's public IP as an ExternalIP, and FindNodeHostnameContextE prefers that, so most clusters need no
// lookup at all. Set one only when your cluster does not advertise an ExternalIP.
//
// It exists so that k8s can resolve a node's externally reachable address without depending on any cloud provider
// module. It matches the signature of aws.GetPublicIpsOfEc2InstancesContextE, so on EKS you can wire it up
// directly:
//
//	options := k8s.NewKubectlOptions("", "", "default")
//	options.NodePublicIPLookup = aws.GetPublicIpsOfEc2InstancesContextE
type NodePublicIPLookup func(t testing.TestingT, ctx context.Context, instanceIDs []string, region string) (map[string]string, error)

// KubectlOptions represents common options necessary to specify for all Kubectl calls
type KubectlOptions struct {
	Env map[string]string
	// NodePublicIPLookup, if set, resolves the public IP of a node whose provider ID identifies a cloud instance and
	// whose node object carries no ExternalIP address. It is consulted only by the FindNodeHostnameWithOptions
	// functions, and only after the node's own ExternalIP has been checked, so most callers can leave it nil.
	// It is skipped when serializing options, since a function cannot be represented as JSON.
	NodePublicIPLookup NodePublicIPLookup `json:"-"`
	// RestConfig is not serialized. A rest.Config cannot be rebuilt from JSON, so rather than drop it silently
	// MarshalJSON refuses to encode options that carry one. See ErrRestConfigNotSerializable.
	RestConfig     *rest.Config `json:"-"`
	Logger         *logger.Logger
	ContextName    string
	ConfigPath     string
	Namespace      string
	RequestTimeout time.Duration
	InClusterAuth  bool
}

// NewKubectlOptions will return a pointer to new instance of KubectlOptions with the configured options
func NewKubectlOptions(contextName string, configPath string, namespace string) *KubectlOptions {
	return &KubectlOptions{
		ContextName: contextName,
		ConfigPath:  configPath,
		Namespace:   namespace,
		Env:         map[string]string{},
	}
}

// NewKubectlOptionsWithInClusterAuth will return a pointer to a new instance of KubectlOptions with the InClusterAuth field set to true
func NewKubectlOptionsWithInClusterAuth() *KubectlOptions {
	return &KubectlOptions{
		InClusterAuth: true,
	}
}

// NewKubectlOptionsWithRestConfig will return a pointer to a new instance of KubectlOptions with pre-built config object
func NewKubectlOptionsWithRestConfig(config *rest.Config, namespace string) *KubectlOptions {
	return &KubectlOptions{
		Namespace:  namespace,
		RestConfig: config,
	}
}

// GetConfigPath will return a sensible default if the config path is not set on the options.
func (kubectlOptions *KubectlOptions) GetConfigPath(t testing.TestingT) (string, error) {
	// We predeclare `err` here so that we can update `kubeConfigPath` in the if block below. Otherwise, go complains
	// saying `err` is undefined.
	var err error

	kubeConfigPath := kubectlOptions.ConfigPath
	if kubeConfigPath == "" {
		kubeConfigPath, err = GetKubeConfigPathContextE(t, context.Background())
		if err != nil {
			return "", err
		}
	}

	return kubeConfigPath, nil
}

// MarshalJSON implements json.Marshaler.
//
// It exists to make the loss explicit. RestConfig and NodePublicIPLookup cannot be encoded, and encoding/json
// rejects their func-typed contents outright, so both are tagged json:"-". Dropping NodePublicIPLookup is benign,
// since a missing lookup degrades to a visible node resolution failure. Dropping RestConfig is not, so this
// returns ErrRestConfigNotSerializable instead of writing options that would silently target the wrong cluster.
func (options KubectlOptions) MarshalJSON() ([]byte, error) {
	if options.RestConfig != nil {
		return nil, ErrRestConfigNotSerializable
	}

	// alias drops the MarshalJSON method, so this does not recurse.
	type alias KubectlOptions

	return json.Marshal(alias(options))
}
