package gcp

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/bigquerydatatransfer/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

// GetBigQueryTransferConfigAttrs returns the settings Google Cloud holds for the given BigQuery data transfer configuration, so a test can assert on
// what was actually created rather than only that it exists. Google assigns a transfer configuration its id, so the caller
// passes the one it got back rather than a name it chose. The whole resource name is accepted too, because that is what the
// Terraform resource's id is, and Google names the configuration by project number while the caller is holding an id.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetBigQueryTransferConfigAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, configID string) *bigquerydatatransfer.TransferConfig {
	config, err := GetBigQueryTransferConfigAttrsE(t, ctx, projectID, location, configID)
	require.NoError(t, err)

	return config
}

// GetBigQueryTransferConfigAttrsE returns the settings Google Cloud holds for the given BigQuery data transfer configuration.
// The ctx parameter supports cancellation and timeouts.
func GetBigQueryTransferConfigAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, configID string) (*bigquerydatatransfer.TransferConfig, error) {
	logger.Default.Logf(t, "Getting settings for BigQuery transfer configuration %s in %s in project %s", configID, location, projectID)

	service, err := NewBigQueryDataTransferServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetBigQueryTransferConfigAttrsWithClient(ctx, service, projectID, location, configID)
}

// GetBigQueryTransferConfigAttrsWithClient returns the settings Google Cloud holds for the given BigQuery data transfer configuration using the supplied
// *bigquerydatatransfer.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see bigquerydatatransfer_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetBigQueryTransferConfigAttrsWithClient(ctx context.Context, service *bigquerydatatransfer.Service, projectID string, location string, configID string) (*bigquerydatatransfer.TransferConfig, error) {
	// Either the id on its own or the whole resource name: Google answers with a name built from the
	// project number, and the Terraform resource's id is that name, so prefixing it again would ask for
	// a path that does not exist.
	name := configID
	if !strings.Contains(configID, "/") {
		name = fmt.Sprintf("projects/%s/locations/%s/transferConfigs/%s", projectID, location, configID)
	}

	config, err := service.Projects.Locations.TransferConfigs.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the BigQuery transfer configuration %s does not exist in %s in project %s", configID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for BigQuery transfer configuration %s in %s in project %s: %w", configID, location, projectID, err)
	}

	return config, nil
}

// NewBigQueryDataTransferServiceE creates a BigQuery Data Transfer service authenticated the same way every other client in this
// module is.
// The ctx parameter supports cancellation and timeouts.
func NewBigQueryDataTransferServiceE(t testing.TestingT, ctx context.Context) (*bigquerydatatransfer.Service, error) {
	return bigquerydatatransfer.NewService(ctx, append(withOptions(), option.WithScopes(bigquerydatatransfer.CloudPlatformScope))...)
}
