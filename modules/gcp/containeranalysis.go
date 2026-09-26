package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/containeranalysis/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

// GetContainerAnalysisNoteAttrs returns the settings Google Cloud holds for the given note, so a test can assert on
// what was actually created rather than only that it exists. A note is the description of a kind of finding; an occurrence
// is one instance of it.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetContainerAnalysisNoteAttrs(t testing.TestingT, ctx context.Context, projectID string, noteID string) *containeranalysis.Note {
	note, err := GetContainerAnalysisNoteAttrsE(t, ctx, projectID, noteID)
	require.NoError(t, err)

	return note
}

// GetContainerAnalysisNoteAttrsE returns the settings Google Cloud holds for the given note.
// The ctx parameter supports cancellation and timeouts.
func GetContainerAnalysisNoteAttrsE(t testing.TestingT, ctx context.Context, projectID string, noteID string) (*containeranalysis.Note, error) {
	logger.Default.Logf(t, "Getting settings for note %s in project %s", noteID, projectID)

	service, err := NewContainerAnalysisServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetContainerAnalysisNoteAttrsWithClient(ctx, service, projectID, noteID)
}

// GetContainerAnalysisNoteAttrsWithClient returns the settings Google Cloud holds for the given note using the supplied
// *containeranalysis.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see containeranalysis_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetContainerAnalysisNoteAttrsWithClient(ctx context.Context, service *containeranalysis.Service, projectID string, noteID string) (*containeranalysis.Note, error) {
	name := fmt.Sprintf("projects/%s/notes/%s", projectID, noteID)

	note, err := service.Projects.Notes.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the note %s does not exist in project %s", noteID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for note %s in project %s: %w", noteID, projectID, err)
	}

	return note, nil
}

// GetContainerAnalysisOccurrenceAttrs returns the settings Google Cloud holds for the given occurrence, so a test can assert on
// what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetContainerAnalysisOccurrenceAttrs(t testing.TestingT, ctx context.Context, projectID string, occurrenceID string) *containeranalysis.Occurrence {
	occurrence, err := GetContainerAnalysisOccurrenceAttrsE(t, ctx, projectID, occurrenceID)
	require.NoError(t, err)

	return occurrence
}

// GetContainerAnalysisOccurrenceAttrsE returns the settings Google Cloud holds for the given occurrence.
// The ctx parameter supports cancellation and timeouts.
func GetContainerAnalysisOccurrenceAttrsE(t testing.TestingT, ctx context.Context, projectID string, occurrenceID string) (*containeranalysis.Occurrence, error) {
	logger.Default.Logf(t, "Getting settings for occurrence %s in project %s", occurrenceID, projectID)

	service, err := NewContainerAnalysisServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetContainerAnalysisOccurrenceAttrsWithClient(ctx, service, projectID, occurrenceID)
}

// GetContainerAnalysisOccurrenceAttrsWithClient returns the settings Google Cloud holds for the given occurrence using the supplied
// *containeranalysis.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see containeranalysis_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetContainerAnalysisOccurrenceAttrsWithClient(ctx context.Context, service *containeranalysis.Service, projectID string, occurrenceID string) (*containeranalysis.Occurrence, error) {
	name := fmt.Sprintf("projects/%s/occurrences/%s", projectID, occurrenceID)

	occurrence, err := service.Projects.Occurrences.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the occurrence %s does not exist in project %s", occurrenceID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for occurrence %s in project %s: %w", occurrenceID, projectID, err)
	}

	return occurrence, nil
}

// NewContainerAnalysisServiceE creates a Container Analysis service authenticated the same way every other client in this
// module is.
// The ctx parameter supports cancellation and timeouts.
func NewContainerAnalysisServiceE(t testing.TestingT, ctx context.Context) (*containeranalysis.Service, error) {
	return containeranalysis.NewService(ctx, append(withOptions(), option.WithScopes(containeranalysis.CloudPlatformScope))...)
}
