package k8s

import (
	"context"
	"os"
	"path/filepath"
	"sort"

	gwErrors "github.com/gruntwork-io/go-commons/errors"
	"github.com/stretchr/testify/require"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"

	"github.com/gruntwork-io/terratest/modules/environment"
	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/testing"
)

// LoadConfigFromPath will load a ClientConfig object from a file path that points to a location on disk containing a
// kubectl config.
func LoadConfigFromPath(path string) clientcmd.ClientConfig {
	config := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: path},
		&clientcmd.ConfigOverrides{})

	return config
}

// LoadAPIClientConfigE loads a ClientConfig object from a file path that points to a location on disk containing a
// kubectl config, with the requested context loaded.
func LoadAPIClientConfigE(configPath string, contextName string) (*restclient.Config, error) {
	overrides := clientcmd.ConfigOverrides{}
	if contextName != "" {
		overrides.CurrentContext = contextName
	}

	config := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: configPath},
		&overrides)

	return config.ClientConfig()
}

// LoadApiClientConfigE loads a ClientConfig object from a file path that points to a location on disk containing a
// kubectl config, with the requested context loaded.
//
// Deprecated: Use LoadAPIClientConfigE instead.
func LoadApiClientConfigE(configPath string, contextName string) (*restclient.Config, error) { //nolint:staticcheck,revive // deprecated
	return LoadAPIClientConfigE(configPath, contextName)
}

// DeleteConfigContextE will remove the context specified at the provided name, and remove any clusters and authinfos
// that are orphaned as a result of it. The config path is either specified in the environment variable KUBECONFIG or at
// the user's home directory under `.kube/config`.
func DeleteConfigContextE(t testing.TestingT, contextName string) error {
	kubeConfigPath, err := GetKubeConfigPathE(t)
	if err != nil {
		return err
	}

	return DeleteConfigContextWithPathE(t, kubeConfigPath, contextName)
}

// DeleteConfigContextWithPathE will remove the context specified at the provided name, and remove any clusters and
// authinfos that are orphaned as a result of it.
func DeleteConfigContextWithPathE(t testing.TestingT, kubeConfigPath string, contextName string) error {
	logger.Default.Logf(t, "Removing kubectl config context %s from config at path %s", contextName, kubeConfigPath)

	// Load config and get data structure representing config info
	config := LoadConfigFromPath(kubeConfigPath)

	rawConfig, err := config.RawConfig()
	if err != nil {
		return err
	}

	// Check if the context we want to delete actually exists, and if so, delete it.
	_, ok := rawConfig.Contexts[contextName]
	if !ok {
		logger.Default.Logf(t, "WARNING: Could not find context %s from config at path %s", contextName, kubeConfigPath)
		return nil
	}

	delete(rawConfig.Contexts, contextName)

	// If the removing context is the current context, be sure to set a new one
	if contextName == rawConfig.CurrentContext {
		if err := setNewContext(&rawConfig); err != nil {
			return err
		}
	}

	// Finally, clean up orphaned clusters and authinfos and then save config
	RemoveOrphanedClusterAndAuthInfoConfig(&rawConfig)

	if err := clientcmd.ModifyConfig(config.ConfigAccess(), rawConfig, false); err != nil {
		return err
	}

	logger.Default.Logf(
		t,
		"Removed context %s from config at path %s and any orphaned clusters and authinfos",
		contextName,
		kubeConfigPath)

	return nil
}

// setNewContext will pick the alphebetically first available context from the list of contexts in the config to use as
// the new current context
func setNewContext(config *api.Config) error {
	// Sort contextNames and pick the first one
	contextNames := make([]string, 0, len(config.Contexts))
	for name := range config.Contexts {
		contextNames = append(contextNames, name)
	}

	sort.Strings(contextNames)

	if len(contextNames) > 0 {
		config.CurrentContext = contextNames[0]
	} else {
		return ErrNoAvailableContexts
	}

	return nil
}

// RemoveOrphanedClusterAndAuthInfoConfig will remove all configurations related to clusters and users that have no
// contexts associated with it
func RemoveOrphanedClusterAndAuthInfoConfig(config *api.Config) {
	newAuthInfos := map[string]*api.AuthInfo{}

	newClusters := map[string]*api.Cluster{}
	for _, context := range config.Contexts {
		newClusters[context.Cluster] = config.Clusters[context.Cluster]
		newAuthInfos[context.AuthInfo] = config.AuthInfos[context.AuthInfo]
	}

	config.AuthInfos = newAuthInfos
	config.Clusters = newClusters
}

// GetKubeConfigPathContextE determines which file path to use as the kubectl config path.
// The ctx parameter is accepted for API consistency.
func GetKubeConfigPathContextE(t testing.TestingT, ctx context.Context) (string, error) {
	kubeConfigPath := environment.GetFirstNonEmptyEnvVarOrEmptyString(t, []string{"KUBECONFIG"})
	if kubeConfigPath == "" {
		configPath, err := KubeConfigPathFromHomeDirE()
		if err != nil {
			return "", err
		}

		kubeConfigPath = configPath
	}

	return kubeConfigPath, nil
}

// GetKubeConfigPathContext determines which file path to use as the kubectl config path.
// The ctx parameter is accepted for API consistency.
// This will fail the test if there is an error.
func GetKubeConfigPathContext(t testing.TestingT, ctx context.Context) string {
	t.Helper()
	path, err := GetKubeConfigPathContextE(t, ctx)
	require.NoError(t, err)

	return path
}

// GetKubeConfigPathE determines which file path to use as the kubectl config path
//
// Deprecated: Use [GetKubeConfigPathContextE] instead.
func GetKubeConfigPathE(t testing.TestingT) (string, error) {
	return GetKubeConfigPathContextE(t, context.Background())
}

// KubeConfigPathFromHomeDirE returns a string to the default Kubernetes config path in the home directory. This will
// error if the home directory can not be determined.
func KubeConfigPathFromHomeDirE() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configPath := filepath.Join(home, ".kube", "config")

	return configPath, err
}

// CopyHomeKubeConfigToTemp will copy the kubeconfig in the home directory to a temp file. This will fail the test if
// there are any errors.
func CopyHomeKubeConfigToTemp(t testing.TestingT) string {
	path, err := CopyHomeKubeConfigToTempE(t)
	if err != nil {
		if path != "" {
			_ = os.Remove(path)
		}

		t.Fatal(err)
	}

	return path
}

// CopyHomeKubeConfigToTempE will copy the kubeconfig in the home directory to a temp file.
func CopyHomeKubeConfigToTempE(t testing.TestingT) (string, error) {
	configPath, err := KubeConfigPathFromHomeDirE()
	if err != nil {
		return "", err
	}

	tmpConfig, err := os.CreateTemp("", "")
	if err != nil {
		return "", gwErrors.WithStackTrace(err)
	}
	defer tmpConfig.Close()

	err = files.CopyFile(configPath, tmpConfig.Name())

	return tmpConfig.Name(), err
}

// UpsertConfigContext will update or insert a new context to the provided config, binding the provided cluster to the
// provided user.
func UpsertConfigContext(config *api.Config, contextName string, clusterName string, userName string) {
	config.Contexts[contextName] = &api.Context{Cluster: clusterName, AuthInfo: userName}
}
