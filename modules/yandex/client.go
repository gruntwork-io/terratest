package yandex

import (
	"context"
	"os"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/yandex-cloud/go-genproto/yandex/cloud/endpoint"
	"github.com/yandex-cloud/go-sdk/iamkey"
	"github.com/yandex-cloud/go-sdk/pkg/requestid"
	"github.com/yandex-cloud/go-sdk/pkg/retry"
	"google.golang.org/grpc/codes"

	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

import ycsdk "github.com/yandex-cloud/go-sdk"

const MaxRetries = 3

// NewYandexClientE creates an Yandex.Cloud client.
func NewYandexClientE(t *testing.T) (*ycsdk.SDK, error) {
	logger.Logf(t, "Initialize Yandex.Cloud client")

	sdkConfig := ycsdk.Config{}

	switch {
	case os.Getenv("YC_TOKEN") == "" && os.Getenv("YC_SERVICE_ACCOUNT_KEY_FILE") == "":
		logger.Logf(t, "Use Instance Service Account for authentication")
		sdkConfig.Credentials = ycsdk.InstanceServiceAccount()

	case os.Getenv("YC_TOKEN") != "":
		logger.Log(t, "Use OAuth token for authentication")
		sdkConfig.Credentials = ycsdk.OAuthToken(os.Getenv("YC_TOKEN"))

	case os.Getenv("YC_SERVICE_ACCOUNT_KEY_FILE") != "":
		logger.Logf(t, " Use Service Account key file %q for authentication", os.Getenv("YC_SERVICE_ACCOUNT_KEY_FILE"))
		key, err := iamkey.ReadFromJSONFile(os.Getenv("YC_SERVICE_ACCOUNT_KEY_FILE"))
		if err != nil {
			return nil, err
		}

		credentials, err := ycsdk.ServiceAccountKey(key)
		if err != nil {
			return nil, err
		}

		sdkConfig.Credentials = credentials
	}

	requestIDInterceptor := requestid.Interceptor()

	retryInterceptor := retry.Interceptor(
		retry.WithMax(MaxRetries),
		retry.WithCodes(codes.Unavailable),
		retry.WithAttemptHeader(true),
		retry.WithBackoff(retry.DefaultBackoff()))

	// Make sure retry interceptor is above id interceptor.
	// Now we will have new request id for every retry attempt.
	interceptorChain := grpc_middleware.ChainUnaryClient(retryInterceptor, requestIDInterceptor)

	userAgentMD := metadata.Pairs("user-agent", "Terratest")

	sdk, err := ycsdk.Build(context.Background(), sdkConfig,
		grpc.WithDefaultCallOptions(grpc.Header(&userAgentMD)),
		grpc.WithUnaryInterceptor(interceptorChain))

	if err != nil {
		return nil, err
	}

	if _, err = sdk.ApiEndpoint().ApiEndpoint().List(context.Background(), &endpoint.ListApiEndpointsRequest{}); err != nil {
		return nil, err
	}

	return sdk, nil
}
