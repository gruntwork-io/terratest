package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/dataplex/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

// GetDataplexLakeAttrs returns the settings Google Cloud holds for the given Dataplex lake, so a
// test can assert on what was actually created rather than only that it exists. A lake lives in a
// location, which has to be given.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDataplexLakeAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, lakeID string) *dataplex.GoogleCloudDataplexV1Lake {
	lake, err := GetDataplexLakeAttrsE(t, ctx, projectID, location, lakeID)
	require.NoError(t, err)

	return lake
}

// GetDataplexLakeAttrsE returns the settings Google Cloud holds for the given Dataplex lake.
// The ctx parameter supports cancellation and timeouts.
func GetDataplexLakeAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, lakeID string) (*dataplex.GoogleCloudDataplexV1Lake, error) {
	logger.Default.Logf(t, "Getting settings for Dataplex lake %s in project %s", lakeID, projectID)

	service, err := NewDataplexServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetDataplexLakeAttrsWithClient(ctx, service, projectID, location, lakeID)
}

// GetDataplexLakeAttrsWithClient returns the settings Google Cloud holds for the given Dataplex
// lake using the supplied *dataplex.Service. Prefer this variant in unit tests where the service is
// backed by an httptest fake server (see dataplex_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetDataplexLakeAttrsWithClient(ctx context.Context, service *dataplex.Service, projectID string, location string, lakeID string) (*dataplex.GoogleCloudDataplexV1Lake, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/lakes/%s", projectID, location, lakeID)

	lake, err := service.Projects.Locations.Lakes.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("Dataplex lake %s does not exist in project %s", lakeID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Dataplex lake %s in project %s: %w", lakeID, projectID, err)
	}

	return lake, nil
}

// NewDataplexServiceE creates a Dataplex service authenticated the same way every other client in
// this module is.
// The ctx parameter supports cancellation and timeouts.
func NewDataplexServiceE(t testing.TestingT, ctx context.Context) (*dataplex.Service, error) {
	return dataplex.NewService(ctx, append(withOptions(), option.WithScopes(dataplex.CloudPlatformScope))...)
}

// GetDataplexZoneAttrs returns the settings Google Cloud holds for the given Dataplex zone, so a test can assert on
// what was actually created rather than only that it exists. A zone belongs to a lake, so it is named by the
// lake's id as well as its own.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDataplexZoneAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, lakeID string, zoneID string) *dataplex.GoogleCloudDataplexV1Zone {
	zone, err := GetDataplexZoneAttrsE(t, ctx, projectID, location, lakeID, zoneID)
	require.NoError(t, err)

	return zone
}

// GetDataplexZoneAttrsE returns the settings Google Cloud holds for the given Dataplex zone.
// The ctx parameter supports cancellation and timeouts.
func GetDataplexZoneAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, lakeID string, zoneID string) (*dataplex.GoogleCloudDataplexV1Zone, error) {
	logger.Default.Logf(t, "Getting settings for Dataplex zone %s in lake %s in %s in project %s", zoneID, lakeID, location, projectID)

	service, err := NewDataplexServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetDataplexZoneAttrsWithClient(ctx, service, projectID, location, lakeID, zoneID)
}

// GetDataplexZoneAttrsWithClient returns the settings Google Cloud holds for the given Dataplex zone using the supplied
// *dataplex.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see dataplex_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetDataplexZoneAttrsWithClient(ctx context.Context, service *dataplex.Service, projectID string, location string, lakeID string, zoneID string) (*dataplex.GoogleCloudDataplexV1Zone, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/lakes/%s/zones/%s", projectID, location, lakeID, zoneID)

	zone, err := service.Projects.Locations.Lakes.Zones.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Dataplex zone %s does not exist in lake %s in %s in project %s", zoneID, lakeID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Dataplex zone %s in lake %s in %s in project %s: %w", zoneID, lakeID, location, projectID, err)
	}

	return zone, nil
}

// GetDataplexAssetAttrs returns the settings Google Cloud holds for the given Dataplex asset, so a test can assert on
// what was actually created rather than only that it exists. An asset is the storage a zone points at, and is named
// by the lake and the zone as well as itself.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDataplexAssetAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, lakeID string, zoneID string, assetID string) *dataplex.GoogleCloudDataplexV1Asset {
	asset, err := GetDataplexAssetAttrsE(t, ctx, projectID, location, lakeID, zoneID, assetID)
	require.NoError(t, err)

	return asset
}

// GetDataplexAssetAttrsE returns the settings Google Cloud holds for the given Dataplex asset.
// The ctx parameter supports cancellation and timeouts.
func GetDataplexAssetAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, lakeID string, zoneID string, assetID string) (*dataplex.GoogleCloudDataplexV1Asset, error) {
	logger.Default.Logf(t, "Getting settings for Dataplex asset %s in zone %s in lake %s in %s in project %s", assetID, zoneID, lakeID, location, projectID)

	service, err := NewDataplexServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetDataplexAssetAttrsWithClient(ctx, service, projectID, location, lakeID, zoneID, assetID)
}

// GetDataplexAssetAttrsWithClient returns the settings Google Cloud holds for the given Dataplex asset using the supplied
// *dataplex.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see dataplex_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetDataplexAssetAttrsWithClient(ctx context.Context, service *dataplex.Service, projectID string, location string, lakeID string, zoneID string, assetID string) (*dataplex.GoogleCloudDataplexV1Asset, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/lakes/%s/zones/%s/assets/%s", projectID, location, lakeID, zoneID, assetID)

	asset, err := service.Projects.Locations.Lakes.Zones.Assets.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Dataplex asset %s does not exist in zone %s in lake %s in %s in project %s", assetID, zoneID, lakeID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Dataplex asset %s in zone %s in lake %s in %s in project %s: %w", assetID, zoneID, lakeID, location, projectID, err)
	}

	return asset, nil
}

// GetDataplexEntryGroupAttrs returns the settings Google Cloud holds for the given Dataplex entry group, so a test can assert on
// what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDataplexEntryGroupAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, entryGroupID string) *dataplex.GoogleCloudDataplexV1EntryGroup {
	entryGroup, err := GetDataplexEntryGroupAttrsE(t, ctx, projectID, location, entryGroupID)
	require.NoError(t, err)

	return entryGroup
}

// GetDataplexEntryGroupAttrsE returns the settings Google Cloud holds for the given Dataplex entry group.
// The ctx parameter supports cancellation and timeouts.
func GetDataplexEntryGroupAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, entryGroupID string) (*dataplex.GoogleCloudDataplexV1EntryGroup, error) {
	logger.Default.Logf(t, "Getting settings for Dataplex entry group %s in %s in project %s", entryGroupID, location, projectID)

	service, err := NewDataplexServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetDataplexEntryGroupAttrsWithClient(ctx, service, projectID, location, entryGroupID)
}

// GetDataplexEntryGroupAttrsWithClient returns the settings Google Cloud holds for the given Dataplex entry group using the supplied
// *dataplex.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see dataplex_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetDataplexEntryGroupAttrsWithClient(ctx context.Context, service *dataplex.Service, projectID string, location string, entryGroupID string) (*dataplex.GoogleCloudDataplexV1EntryGroup, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/entryGroups/%s", projectID, location, entryGroupID)

	entryGroup, err := service.Projects.Locations.EntryGroups.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Dataplex entry group %s does not exist in %s in project %s", entryGroupID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Dataplex entry group %s in %s in project %s: %w", entryGroupID, location, projectID, err)
	}

	return entryGroup, nil
}

// GetDataplexEntryTypeAttrs returns the settings Google Cloud holds for the given Dataplex entry type, so a test can assert on
// what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDataplexEntryTypeAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, entryTypeID string) *dataplex.GoogleCloudDataplexV1EntryType {
	entryType, err := GetDataplexEntryTypeAttrsE(t, ctx, projectID, location, entryTypeID)
	require.NoError(t, err)

	return entryType
}

// GetDataplexEntryTypeAttrsE returns the settings Google Cloud holds for the given Dataplex entry type.
// The ctx parameter supports cancellation and timeouts.
func GetDataplexEntryTypeAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, entryTypeID string) (*dataplex.GoogleCloudDataplexV1EntryType, error) {
	logger.Default.Logf(t, "Getting settings for Dataplex entry type %s in %s in project %s", entryTypeID, location, projectID)

	service, err := NewDataplexServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetDataplexEntryTypeAttrsWithClient(ctx, service, projectID, location, entryTypeID)
}

// GetDataplexEntryTypeAttrsWithClient returns the settings Google Cloud holds for the given Dataplex entry type using the supplied
// *dataplex.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see dataplex_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetDataplexEntryTypeAttrsWithClient(ctx context.Context, service *dataplex.Service, projectID string, location string, entryTypeID string) (*dataplex.GoogleCloudDataplexV1EntryType, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/entryTypes/%s", projectID, location, entryTypeID)

	entryType, err := service.Projects.Locations.EntryTypes.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Dataplex entry type %s does not exist in %s in project %s", entryTypeID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Dataplex entry type %s in %s in project %s: %w", entryTypeID, location, projectID, err)
	}

	return entryType, nil
}

// GetDataplexAspectTypeAttrs returns the settings Google Cloud holds for the given Dataplex aspect type, so a test can assert on
// what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDataplexAspectTypeAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, aspectTypeID string) *dataplex.GoogleCloudDataplexV1AspectType {
	aspectType, err := GetDataplexAspectTypeAttrsE(t, ctx, projectID, location, aspectTypeID)
	require.NoError(t, err)

	return aspectType
}

// GetDataplexAspectTypeAttrsE returns the settings Google Cloud holds for the given Dataplex aspect type.
// The ctx parameter supports cancellation and timeouts.
func GetDataplexAspectTypeAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, aspectTypeID string) (*dataplex.GoogleCloudDataplexV1AspectType, error) {
	logger.Default.Logf(t, "Getting settings for Dataplex aspect type %s in %s in project %s", aspectTypeID, location, projectID)

	service, err := NewDataplexServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetDataplexAspectTypeAttrsWithClient(ctx, service, projectID, location, aspectTypeID)
}

// GetDataplexAspectTypeAttrsWithClient returns the settings Google Cloud holds for the given Dataplex aspect type using the supplied
// *dataplex.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see dataplex_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetDataplexAspectTypeAttrsWithClient(ctx context.Context, service *dataplex.Service, projectID string, location string, aspectTypeID string) (*dataplex.GoogleCloudDataplexV1AspectType, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/aspectTypes/%s", projectID, location, aspectTypeID)

	aspectType, err := service.Projects.Locations.AspectTypes.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Dataplex aspect type %s does not exist in %s in project %s", aspectTypeID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Dataplex aspect type %s in %s in project %s: %w", aspectTypeID, location, projectID, err)
	}

	return aspectType, nil
}

// GetDataplexEntryAttrs returns the settings Google Cloud holds for the given Dataplex entry, so a test can assert on
// what was actually created rather than only that it exists. An entry belongs to an entry group, so it is named by
// the group's id as well as its own.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDataplexEntryAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, entryGroupID string, entryID string) *dataplex.GoogleCloudDataplexV1Entry {
	entry, err := GetDataplexEntryAttrsE(t, ctx, projectID, location, entryGroupID, entryID)
	require.NoError(t, err)

	return entry
}

// GetDataplexEntryAttrsE returns the settings Google Cloud holds for the given Dataplex entry.
// The ctx parameter supports cancellation and timeouts.
func GetDataplexEntryAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, entryGroupID string, entryID string) (*dataplex.GoogleCloudDataplexV1Entry, error) {
	logger.Default.Logf(t, "Getting settings for Dataplex entry %s in entry group %s in %s in project %s", entryID, entryGroupID, location, projectID)

	service, err := NewDataplexServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetDataplexEntryAttrsWithClient(ctx, service, projectID, location, entryGroupID, entryID)
}

// GetDataplexEntryAttrsWithClient returns the settings Google Cloud holds for the given Dataplex entry using the supplied
// *dataplex.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see dataplex_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetDataplexEntryAttrsWithClient(ctx context.Context, service *dataplex.Service, projectID string, location string, entryGroupID string, entryID string) (*dataplex.GoogleCloudDataplexV1Entry, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/entryGroups/%s/entries/%s", projectID, location, entryGroupID, entryID)

	entry, err := service.Projects.Locations.EntryGroups.Entries.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Dataplex entry %s does not exist in entry group %s in %s in project %s", entryID, entryGroupID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Dataplex entry %s in entry group %s in %s in project %s: %w", entryID, entryGroupID, location, projectID, err)
	}

	return entry, nil
}

// GetDataplexGlossaryAttrs returns the settings Google Cloud holds for the given Dataplex glossary, so a test can assert on
// what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDataplexGlossaryAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, glossaryID string) *dataplex.GoogleCloudDataplexV1Glossary {
	glossary, err := GetDataplexGlossaryAttrsE(t, ctx, projectID, location, glossaryID)
	require.NoError(t, err)

	return glossary
}

// GetDataplexGlossaryAttrsE returns the settings Google Cloud holds for the given Dataplex glossary.
// The ctx parameter supports cancellation and timeouts.
func GetDataplexGlossaryAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, glossaryID string) (*dataplex.GoogleCloudDataplexV1Glossary, error) {
	logger.Default.Logf(t, "Getting settings for Dataplex glossary %s in %s in project %s", glossaryID, location, projectID)

	service, err := NewDataplexServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetDataplexGlossaryAttrsWithClient(ctx, service, projectID, location, glossaryID)
}

// GetDataplexGlossaryAttrsWithClient returns the settings Google Cloud holds for the given Dataplex glossary using the supplied
// *dataplex.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see dataplex_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetDataplexGlossaryAttrsWithClient(ctx context.Context, service *dataplex.Service, projectID string, location string, glossaryID string) (*dataplex.GoogleCloudDataplexV1Glossary, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/glossaries/%s", projectID, location, glossaryID)

	glossary, err := service.Projects.Locations.Glossaries.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Dataplex glossary %s does not exist in %s in project %s", glossaryID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Dataplex glossary %s in %s in project %s: %w", glossaryID, location, projectID, err)
	}

	return glossary, nil
}

// GetDataplexGlossaryCategoryAttrs returns the settings Google Cloud holds for the given glossary category, so a test can assert on
// what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDataplexGlossaryCategoryAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, glossaryID string, categoryID string) *dataplex.GoogleCloudDataplexV1GlossaryCategory {
	category, err := GetDataplexGlossaryCategoryAttrsE(t, ctx, projectID, location, glossaryID, categoryID)
	require.NoError(t, err)

	return category
}

// GetDataplexGlossaryCategoryAttrsE returns the settings Google Cloud holds for the given glossary category.
// The ctx parameter supports cancellation and timeouts.
func GetDataplexGlossaryCategoryAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, glossaryID string, categoryID string) (*dataplex.GoogleCloudDataplexV1GlossaryCategory, error) {
	logger.Default.Logf(t, "Getting settings for glossary category %s in glossary %s in %s in project %s", categoryID, glossaryID, location, projectID)

	service, err := NewDataplexServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetDataplexGlossaryCategoryAttrsWithClient(ctx, service, projectID, location, glossaryID, categoryID)
}

// GetDataplexGlossaryCategoryAttrsWithClient returns the settings Google Cloud holds for the given glossary category using the supplied
// *dataplex.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see dataplex_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetDataplexGlossaryCategoryAttrsWithClient(ctx context.Context, service *dataplex.Service, projectID string, location string, glossaryID string, categoryID string) (*dataplex.GoogleCloudDataplexV1GlossaryCategory, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/glossaries/%s/categories/%s", projectID, location, glossaryID, categoryID)

	category, err := service.Projects.Locations.Glossaries.Categories.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the glossary category %s does not exist in glossary %s in %s in project %s", categoryID, glossaryID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for glossary category %s in glossary %s in %s in project %s: %w", categoryID, glossaryID, location, projectID, err)
	}

	return category, nil
}

// GetDataplexGlossaryTermAttrs returns the settings Google Cloud holds for the given glossary term, so a test can assert on
// what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDataplexGlossaryTermAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, glossaryID string, termID string) *dataplex.GoogleCloudDataplexV1GlossaryTerm {
	term, err := GetDataplexGlossaryTermAttrsE(t, ctx, projectID, location, glossaryID, termID)
	require.NoError(t, err)

	return term
}

// GetDataplexGlossaryTermAttrsE returns the settings Google Cloud holds for the given glossary term.
// The ctx parameter supports cancellation and timeouts.
func GetDataplexGlossaryTermAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, glossaryID string, termID string) (*dataplex.GoogleCloudDataplexV1GlossaryTerm, error) {
	logger.Default.Logf(t, "Getting settings for glossary term %s in glossary %s in %s in project %s", termID, glossaryID, location, projectID)

	service, err := NewDataplexServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetDataplexGlossaryTermAttrsWithClient(ctx, service, projectID, location, glossaryID, termID)
}

// GetDataplexGlossaryTermAttrsWithClient returns the settings Google Cloud holds for the given glossary term using the supplied
// *dataplex.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see dataplex_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetDataplexGlossaryTermAttrsWithClient(ctx context.Context, service *dataplex.Service, projectID string, location string, glossaryID string, termID string) (*dataplex.GoogleCloudDataplexV1GlossaryTerm, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/glossaries/%s/terms/%s", projectID, location, glossaryID, termID)

	term, err := service.Projects.Locations.Glossaries.Terms.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the glossary term %s does not exist in glossary %s in %s in project %s", termID, glossaryID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for glossary term %s in glossary %s in %s in project %s: %w", termID, glossaryID, location, projectID, err)
	}

	return term, nil
}

// GetDataplexDataScanAttrs returns the settings Google Cloud holds for the given Dataplex data scan, so a test can assert on
// what was actually created rather than only that it exists. The full view is asked for, because the default answer
// carries no scan specification and a caller asserting on one would read nil.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDataplexDataScanAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, scanID string) *dataplex.GoogleCloudDataplexV1DataScan {
	scan, err := GetDataplexDataScanAttrsE(t, ctx, projectID, location, scanID)
	require.NoError(t, err)

	return scan
}

// GetDataplexDataScanAttrsE returns the settings Google Cloud holds for the given Dataplex data scan.
// The ctx parameter supports cancellation and timeouts.
func GetDataplexDataScanAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, scanID string) (*dataplex.GoogleCloudDataplexV1DataScan, error) {
	logger.Default.Logf(t, "Getting settings for Dataplex data scan %s in %s in project %s", scanID, location, projectID)

	service, err := NewDataplexServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetDataplexDataScanAttrsWithClient(ctx, service, projectID, location, scanID)
}

// GetDataplexDataScanAttrsWithClient returns the settings Google Cloud holds for the given Dataplex data scan using the supplied
// *dataplex.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see dataplex_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetDataplexDataScanAttrsWithClient(ctx context.Context, service *dataplex.Service, projectID string, location string, scanID string) (*dataplex.GoogleCloudDataplexV1DataScan, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/dataScans/%s", projectID, location, scanID)

	scan, err := service.Projects.Locations.DataScans.Get(name).View("FULL").Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Dataplex data scan %s does not exist in %s in project %s", scanID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Dataplex data scan %s in %s in project %s: %w", scanID, location, projectID, err)
	}

	return scan, nil
}
