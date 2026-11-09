package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/datacatalog/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

// GetDataCatalogTaxonomyAttrs returns the settings Google Cloud holds for the given taxonomy, so a test can assert on
// what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDataCatalogTaxonomyAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, taxonomyID string) *datacatalog.GoogleCloudDatacatalogV1Taxonomy {
	taxonomy, err := GetDataCatalogTaxonomyAttrsE(t, ctx, projectID, location, taxonomyID)
	require.NoError(t, err)

	return taxonomy
}

// GetDataCatalogTaxonomyAttrsE returns the settings Google Cloud holds for the given taxonomy.
// The ctx parameter supports cancellation and timeouts.
func GetDataCatalogTaxonomyAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, taxonomyID string) (*datacatalog.GoogleCloudDatacatalogV1Taxonomy, error) {
	logger.Default.Logf(t, "Getting settings for taxonomy %s in %s in project %s", taxonomyID, location, projectID)

	service, err := NewDataCatalogServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetDataCatalogTaxonomyAttrsWithClient(ctx, service, projectID, location, taxonomyID)
}

// GetDataCatalogTaxonomyAttrsWithClient returns the settings Google Cloud holds for the given taxonomy using the supplied
// *datacatalog.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see datacatalog_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetDataCatalogTaxonomyAttrsWithClient(ctx context.Context, service *datacatalog.Service, projectID string, location string, taxonomyID string) (*datacatalog.GoogleCloudDatacatalogV1Taxonomy, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/taxonomies/%s", projectID, location, taxonomyID)

	taxonomy, err := service.Projects.Locations.Taxonomies.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the taxonomy %s does not exist in %s in project %s", taxonomyID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for taxonomy %s in %s in project %s: %w", taxonomyID, location, projectID, err)
	}

	return taxonomy, nil
}

// GetDataCatalogPolicyTagAttrs returns the settings Google Cloud holds for the given policy tag, so a test can assert on
// what was actually created rather than only that it exists. A policy tag belongs to a taxonomy, so it is named by
// the taxonomy's id as well as its own.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDataCatalogPolicyTagAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, taxonomyID string, policyTagID string) *datacatalog.GoogleCloudDatacatalogV1PolicyTag {
	policyTag, err := GetDataCatalogPolicyTagAttrsE(t, ctx, projectID, location, taxonomyID, policyTagID)
	require.NoError(t, err)

	return policyTag
}

// GetDataCatalogPolicyTagAttrsE returns the settings Google Cloud holds for the given policy tag.
// The ctx parameter supports cancellation and timeouts.
func GetDataCatalogPolicyTagAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, taxonomyID string, policyTagID string) (*datacatalog.GoogleCloudDatacatalogV1PolicyTag, error) {
	logger.Default.Logf(t, "Getting settings for policy tag %s in taxonomy %s in %s in project %s", policyTagID, taxonomyID, location, projectID)

	service, err := NewDataCatalogServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetDataCatalogPolicyTagAttrsWithClient(ctx, service, projectID, location, taxonomyID, policyTagID)
}

// GetDataCatalogPolicyTagAttrsWithClient returns the settings Google Cloud holds for the given policy tag using the supplied
// *datacatalog.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see datacatalog_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetDataCatalogPolicyTagAttrsWithClient(ctx context.Context, service *datacatalog.Service, projectID string, location string, taxonomyID string, policyTagID string) (*datacatalog.GoogleCloudDatacatalogV1PolicyTag, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/taxonomies/%s/policyTags/%s", projectID, location, taxonomyID, policyTagID)

	policyTag, err := service.Projects.Locations.Taxonomies.PolicyTags.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the policy tag %s does not exist in taxonomy %s in %s in project %s", policyTagID, taxonomyID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for policy tag %s in taxonomy %s in %s in project %s: %w", policyTagID, taxonomyID, location, projectID, err)
	}

	return policyTag, nil
}

// GetDataCatalogEntryAttrs returns the settings Google Cloud holds for the given Data Catalog entry, so a test can assert on
// what was actually created rather than only that it exists.
// The schema and the integrated system Google reports for an entry are what a test asserts on, and
// neither is echoed back from the create request.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDataCatalogEntryAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, entryGroupID string, entryID string) *datacatalog.GoogleCloudDatacatalogV1Entry {
	result, err := GetDataCatalogEntryAttrsE(t, ctx, projectID, location, entryGroupID, entryID)
	require.NoError(t, err)

	return result
}

// GetDataCatalogEntryAttrsE returns the settings Google Cloud holds for the given Data Catalog entry.
// The ctx parameter supports cancellation and timeouts.
func GetDataCatalogEntryAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, entryGroupID string, entryID string) (*datacatalog.GoogleCloudDatacatalogV1Entry, error) {
	logger.Default.Logf(t, "Getting settings for Data Catalog entry %s in entry group %s in %s of project %s", entryID, entryGroupID, location, projectID)

	service, err := NewDataCatalogServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetDataCatalogEntryAttrsWithClient(ctx, service, projectID, location, entryGroupID, entryID)
}

// GetDataCatalogEntryAttrsWithClient returns the settings Google Cloud holds for the given Data Catalog entry using the supplied
// *datacatalog.Service. Prefer this variant in unit tests where the service is backed by an httptest fake server
// (see datacatalog_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetDataCatalogEntryAttrsWithClient(ctx context.Context, service *datacatalog.Service, projectID string, location string, entryGroupID string, entryID string) (*datacatalog.GoogleCloudDatacatalogV1Entry, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/entryGroups/%s/entries/%s", projectID, location, entryGroupID, entryID)

	result, err := service.Projects.Locations.EntryGroups.Entries.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Data Catalog entry %s does not exist in entry group %s in %s of project %s", entryID, entryGroupID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Data Catalog entry %s in entry group %s in %s of project %s: %w", entryID, entryGroupID, location, projectID, err)
	}

	return result, nil
}

// GetDataCatalogEntryTagsAttrs returns the tags Google Cloud holds on the given Data Catalog entry,
// so a test can assert on the fields a tag was created with rather than only that it exists.
// A tag has no name a caller can predict, so the entry's tags are listed and the caller picks the one
// it wants out of the slice.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDataCatalogEntryTagsAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, entryGroupID string, entryID string) []*datacatalog.GoogleCloudDatacatalogV1Tag {
	tags, err := GetDataCatalogEntryTagsAttrsE(t, ctx, projectID, location, entryGroupID, entryID)
	require.NoError(t, err)

	return tags
}

// GetDataCatalogEntryTagsAttrsE returns the tags Google Cloud holds on the given Data Catalog entry.
// The ctx parameter supports cancellation and timeouts.
func GetDataCatalogEntryTagsAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, entryGroupID string, entryID string) ([]*datacatalog.GoogleCloudDatacatalogV1Tag, error) {
	logger.Default.Logf(t, "Getting the tags on Data Catalog entry %s in entry group %s in %s of project %s", entryID, entryGroupID, location, projectID)

	service, err := NewDataCatalogServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetDataCatalogEntryTagsAttrsWithClient(ctx, service, projectID, location, entryGroupID, entryID)
}

// GetDataCatalogEntryTagsAttrsWithClient returns the tags Google Cloud holds on the given Data
// Catalog entry using the supplied *datacatalog.Service. Prefer this variant in unit tests where the
// service is backed by an httptest fake server (see datacatalog_test.go for the pattern).
// Every page is followed, so a test never has to care how many tags an entry carries.
// The ctx parameter supports cancellation and timeouts.
func GetDataCatalogEntryTagsAttrsWithClient(ctx context.Context, service *datacatalog.Service, projectID string, location string, entryGroupID string, entryID string) ([]*datacatalog.GoogleCloudDatacatalogV1Tag, error) {
	parent := fmt.Sprintf("projects/%s/locations/%s/entryGroups/%s/entries/%s", projectID, location, entryGroupID, entryID)

	var tags []*datacatalog.GoogleCloudDatacatalogV1Tag

	call := service.Projects.Locations.EntryGroups.Entries.Tags.List(parent)
	if err := call.Pages(ctx, func(page *datacatalog.GoogleCloudDatacatalogV1ListTagsResponse) error {
		tags = append(tags, page.Tags...)

		return nil
	}); err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Data Catalog entry %s does not exist in entry group %s in %s of project %s", entryID, entryGroupID, location, projectID)
		}

		return nil, fmt.Errorf("failed to list the tags on Data Catalog entry %s in entry group %s in %s of project %s: %w", entryID, entryGroupID, location, projectID, err)
	}

	return tags, nil
}

// NewDataCatalogServiceE creates a Data Catalog service authenticated the same way every other
// client in this module is.
// The ctx parameter supports cancellation and timeouts.
func NewDataCatalogServiceE(t testing.TestingT, ctx context.Context) (*datacatalog.Service, error) {
	return datacatalog.NewService(ctx, append(withOptions(), option.WithScopes(datacatalog.CloudPlatformScope))...)
}
