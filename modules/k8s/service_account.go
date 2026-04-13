package k8s

import (
	"context"
	stderrors "errors"
	"time"

	"github.com/gruntwork-io/go-commons/errors"
	"github.com/stretchr/testify/require"
	authenticationv1 "k8s.io/api/authentication/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"

	"github.com/gruntwork-io/terratest/modules/retry"
	"github.com/gruntwork-io/terratest/modules/testing"
)

const (
	// tokenExpirationSeconds is the expiration time for service account tokens requested via the TokenRequest API.
	tokenExpirationSeconds int64 = 3600
	// tokenSecretRetryCount is the number of retries when waiting for a service account secret to be provisioned.
	tokenSecretRetryCount = 30
	// tokenSecretSleepSeconds is the number of seconds to sleep between retries when waiting for a service account secret.
	tokenSecretSleepSeconds = 10
)

// GetServiceAccountContextE returns a Kubernetes service account resource in the provided namespace with the given name.
// The namespace used is the one provided in the KubectlOptions.
// The ctx parameter supports cancellation and timeouts.
func GetServiceAccountContextE(t testing.TestingT, ctx context.Context, options *KubectlOptions, serviceAccountName string) (*corev1.ServiceAccount, error) {
	clientset, err := GetKubernetesClientFromOptionsContextE(t, ctx, options)
	if err != nil {
		return nil, err
	}

	return clientset.CoreV1().ServiceAccounts(options.Namespace).Get(ctx, serviceAccountName, metav1.GetOptions{})
}

// GetServiceAccountContext returns a Kubernetes service account resource in the provided namespace with the given name.
// The namespace used is the one provided in the KubectlOptions.
// The ctx parameter supports cancellation and timeouts.
// This will fail the test if there is an error.
func GetServiceAccountContext(t testing.TestingT, ctx context.Context, options *KubectlOptions, serviceAccountName string) *corev1.ServiceAccount {
	t.Helper()
	serviceAccount, err := GetServiceAccountContextE(t, ctx, options, serviceAccountName)
	require.NoError(t, err)

	return serviceAccount
}

// GetServiceAccount returns a Kubernetes service account resource in the provided namespace with the given name. The
// namespace used is the one provided in the KubectlOptions. This will fail the test if there is an error.
//
// Deprecated: Use [GetServiceAccountContext] instead.
func GetServiceAccount(t testing.TestingT, options *KubectlOptions, serviceAccountName string) *corev1.ServiceAccount {
	t.Helper()

	return GetServiceAccountContext(t, context.Background(), options, serviceAccountName)
}

// GetServiceAccountE returns a Kubernetes service account resource in the provided namespace with the given name. The
// namespace used is the one provided in the KubectlOptions.
//
// Deprecated: Use [GetServiceAccountContextE] instead.
func GetServiceAccountE(t testing.TestingT, options *KubectlOptions, serviceAccountName string) (*corev1.ServiceAccount, error) {
	return GetServiceAccountContextE(t, context.Background(), options, serviceAccountName)
}

// CreateServiceAccountContextE will create a new service account resource in the provided namespace with the given name.
// The namespace used is the one provided in the KubectlOptions.
// The ctx parameter supports cancellation and timeouts.
func CreateServiceAccountContextE(t testing.TestingT, ctx context.Context, options *KubectlOptions, serviceAccountName string) error {
	clientset, err := GetKubernetesClientFromOptionsContextE(t, ctx, options)
	if err != nil {
		return err
	}

	serviceAccount := corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceAccountName,
			Namespace: options.Namespace,
		},
	}
	_, err = clientset.CoreV1().ServiceAccounts(options.Namespace).Create(ctx, &serviceAccount, metav1.CreateOptions{})

	return err
}

// CreateServiceAccountContext will create a new service account resource in the provided namespace with the given name.
// The namespace used is the one provided in the KubectlOptions.
// The ctx parameter supports cancellation and timeouts.
// This will fail the test if there is an error.
func CreateServiceAccountContext(t testing.TestingT, ctx context.Context, options *KubectlOptions, serviceAccountName string) {
	t.Helper()
	require.NoError(t, CreateServiceAccountContextE(t, ctx, options, serviceAccountName))
}

// CreateServiceAccount will create a new service account resource in the provided namespace with the given name. The
// namespace used is the one provided in the KubectlOptions. This will fail the test if there is an error.
//
// Deprecated: Use [CreateServiceAccountContext] instead.
func CreateServiceAccount(t testing.TestingT, options *KubectlOptions, serviceAccountName string) {
	t.Helper()
	CreateServiceAccountContext(t, context.Background(), options, serviceAccountName)
}

// CreateServiceAccountE will create a new service account resource in the provided namespace with the given name. The
// namespace used is the one provided in the KubectlOptions.
//
// Deprecated: Use [CreateServiceAccountContextE] instead.
func CreateServiceAccountE(t testing.TestingT, options *KubectlOptions, serviceAccountName string) error {
	return CreateServiceAccountContextE(t, context.Background(), options, serviceAccountName)
}

// GetServiceAccountAuthTokenContextE will retrieve the ServiceAccount token from the cluster so it can be used to
// authenticate requests as that ServiceAccount.
// On K8s 1.24+, service account tokens are no longer auto-created as secrets, so this uses the TokenRequest API.
// The ctx parameter supports cancellation and timeouts.
func GetServiceAccountAuthTokenContextE(t testing.TestingT, ctx context.Context, kubectlOptions *KubectlOptions, serviceAccountName string) (string, error) {
	clientset, err := GetKubernetesClientFromOptionsContextE(t, ctx, kubectlOptions)
	if err != nil {
		return "", err
	}

	// First try the TokenRequest API (K8s 1.24+)
	expSeconds := tokenExpirationSeconds
	tokenRequest := &authenticationv1.TokenRequest{
		Spec: authenticationv1.TokenRequestSpec{
			ExpirationSeconds: &expSeconds,
		},
	}

	tokenResponse, err := clientset.CoreV1().ServiceAccounts(kubectlOptions.Namespace).CreateToken(
		ctx,
		serviceAccountName,
		tokenRequest,
		metav1.CreateOptions{},
	)
	if err == nil {
		return tokenResponse.Status.Token, nil
	}

	// Fall back to legacy secret-based tokens for older K8s versions
	kubectlOptions.Logger.Logf(t, "TokenRequest API failed (%s), falling back to secret-based tokens", err)

	msg, retryErr := retry.DoWithRetryContextE(
		t,
		ctx,
		"Waiting for ServiceAccount Token to be provisioned",
		tokenSecretRetryCount,
		tokenSecretSleepSeconds*time.Second,
		func() (string, error) {
			kubectlOptions.Logger.Logf(t, "Checking if service account has secret")

			serviceAccount, err := GetServiceAccountContextE(t, ctx, kubectlOptions, serviceAccountName)
			if err != nil {
				return "", err
			}

			if len(serviceAccount.Secrets) == 0 {
				msg := "no secrets on the service account yet"
				kubectlOptions.Logger.Logf(t, "%s", msg)

				return "", stderrors.New(msg)
			}

			return "Service Account has secret", nil
		},
	)
	if retryErr != nil {
		return "", retryErr
	}

	kubectlOptions.Logger.Logf(t, "%s", msg)

	serviceAccount, err := GetServiceAccountContextE(t, ctx, kubectlOptions, serviceAccountName)
	if err != nil {
		return "", err
	}

	if len(serviceAccount.Secrets) != 1 {
		return "", errors.WithStackTrace(ServiceAccountTokenNotAvailable{serviceAccountName})
	}

	secret, err := GetSecretContextE(t, ctx, kubectlOptions, serviceAccount.Secrets[0].Name)
	if err != nil {
		return "", err
	}

	return string(secret.Data["token"]), nil
}

// GetServiceAccountAuthTokenContext will retrieve the ServiceAccount token from the cluster so it can be used to
// authenticate requests as that ServiceAccount.
// The ctx parameter supports cancellation and timeouts.
// This will fail the test if there is an error.
func GetServiceAccountAuthTokenContext(t testing.TestingT, ctx context.Context, kubectlOptions *KubectlOptions, serviceAccountName string) string {
	t.Helper()
	token, err := GetServiceAccountAuthTokenContextE(t, ctx, kubectlOptions, serviceAccountName)
	require.NoError(t, err)

	return token
}

// GetServiceAccountAuthToken will retrieve the ServiceAccount token from the cluster so it can be used to
// authenticate requests as that ServiceAccount. This will fail the test if there is an error.
//
// Deprecated: Use [GetServiceAccountAuthTokenContext] instead.
func GetServiceAccountAuthToken(t testing.TestingT, kubectlOptions *KubectlOptions, serviceAccountName string) string {
	t.Helper()

	return GetServiceAccountAuthTokenContext(t, context.Background(), kubectlOptions, serviceAccountName)
}

// GetServiceAccountAuthTokenE will retrieve the ServiceAccount token from the cluster so it can be used to
// authenticate requests as that ServiceAccount.
// On K8s 1.24+, service account tokens are no longer auto-created as secrets, so this uses the TokenRequest API.
//
// Deprecated: Use [GetServiceAccountAuthTokenContextE] instead.
func GetServiceAccountAuthTokenE(t testing.TestingT, kubectlOptions *KubectlOptions, serviceAccountName string) (string, error) {
	return GetServiceAccountAuthTokenContextE(t, context.Background(), kubectlOptions, serviceAccountName)
}

// AddConfigContextForServiceAccountE will add a new config context that binds the ServiceAccount auth token to the
// Kubernetes cluster of the current config context.
func AddConfigContextForServiceAccountE(
	t testing.TestingT,
	kubectlOptions *KubectlOptions,
	contextName string,
	serviceAccountName string,
	token string,
) error {
	// First load the config context
	config := LoadConfigFromPath(kubectlOptions.ConfigPath)

	rawConfig, err := config.RawConfig()
	if err != nil {
		return errors.WithStackTrace(err)
	}

	// Next get the current cluster
	currentContext := rawConfig.Contexts[rawConfig.CurrentContext]
	currentCluster := currentContext.Cluster

	// Now insert the auth info for the service account
	rawConfig.AuthInfos[serviceAccountName] = &api.AuthInfo{Token: token}

	// We now have enough info to add the new context
	UpsertConfigContext(&rawConfig, contextName, currentCluster, serviceAccountName)

	// Finally, overwrite the config
	if err := clientcmd.ModifyConfig(config.ConfigAccess(), rawConfig, false); err != nil {
		return errors.WithStackTrace(err)
	}

	return nil
}
