package gcp

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/documentai/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

// GetDocumentAIProcessorAttrs returns the settings Google Cloud holds for the given Document AI
// processor, so a test can assert on what was actually created rather than only that it exists. A
// processor lives in a multi-region such as `us` or `eu`, which has to be given, and the id is the
// one Google assigns.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDocumentAIProcessorAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, processorID string) *documentai.GoogleCloudDocumentaiV1Processor {
	processor, err := GetDocumentAIProcessorAttrsE(t, ctx, projectID, location, processorID)
	require.NoError(t, err)

	return processor
}

// GetDocumentAIProcessorAttrsE returns the settings Google Cloud holds for the given Document AI
// processor.
// The ctx parameter supports cancellation and timeouts.
func GetDocumentAIProcessorAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, processorID string) (*documentai.GoogleCloudDocumentaiV1Processor, error) {
	logger.Default.Logf(t, "Getting settings for Document AI processor %s in project %s", processorID, projectID)

	service, err := NewDocumentAIServiceE(t, ctx, location)
	if err != nil {
		return nil, err
	}

	return GetDocumentAIProcessorAttrsWithClient(ctx, service, projectID, location, processorID)
}

// GetDocumentAIProcessorAttrsWithClient returns the settings Google Cloud holds for the given
// Document AI processor using the supplied *documentai.Service. Prefer this variant in unit tests
// where the service is backed by an httptest fake server (see documentai_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetDocumentAIProcessorAttrsWithClient(ctx context.Context, service *documentai.Service, projectID string, location string, processorID string) (*documentai.GoogleCloudDocumentaiV1Processor, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/processors/%s", projectID, location, processorID)

	processor, err := service.Projects.Locations.Processors.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("Document AI processor %s does not exist in project %s", processorID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Document AI processor %s in project %s: %w", processorID, projectID, err)
	}

	return processor, nil
}

// locationPattern is what a Google Cloud location identifier may contain. It is checked before a
// location reaches an endpoint, because a value carrying a slash, a colon or an at sign would build
// a URL pointing at a host of the caller's choosing rather than at Google.
var locationPattern = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

// NewDocumentAIServiceE creates a Document AI service authenticated the same way every other client
// in this module is. Document AI answers only on a per-location host, so the location decides which
// endpoint the service talks to, and a location that is not a plain identifier is refused.
// The ctx parameter supports cancellation and timeouts.
func NewDocumentAIServiceE(t testing.TestingT, ctx context.Context, location string) (*documentai.Service, error) {
	if !locationPattern.MatchString(location) {
		return nil, fmt.Errorf("%q is not a valid location: a location may hold only lowercase letters, digits and hyphens", location)
	}

	opts := append(withOptions(), option.WithScopes(documentai.CloudPlatformScope))
	opts = append(opts, option.WithEndpoint(fmt.Sprintf("https://%s-documentai.googleapis.com/", location)))

	return documentai.NewService(ctx, opts...)
}
