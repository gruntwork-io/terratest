package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/essentialcontacts/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

// GetEssentialContactAttrs returns the settings Google Cloud holds for the given essential contact
// on a project, so a test can assert on what was actually created rather than only that it exists.
// Google assigns a contact's id when it is created, so the caller passes the id it got back.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetEssentialContactAttrs(t testing.TestingT, ctx context.Context, projectID string, contactID string) *essentialcontacts.GoogleCloudEssentialcontactsV1Contact {
	contact, err := GetEssentialContactAttrsE(t, ctx, projectID, contactID)
	require.NoError(t, err)

	return contact
}

// GetEssentialContactAttrsE returns the settings Google Cloud holds for the given essential contact
// on a project.
// The ctx parameter supports cancellation and timeouts.
func GetEssentialContactAttrsE(t testing.TestingT, ctx context.Context, projectID string, contactID string) (*essentialcontacts.GoogleCloudEssentialcontactsV1Contact, error) {
	logger.Default.Logf(t, "Getting settings for essential contact %s on project %s", contactID, projectID)

	service, err := NewEssentialContactsServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetEssentialContactAttrsWithClient(ctx, service, projectID, contactID)
}

// GetEssentialContactAttrsWithClient returns the settings Google Cloud holds for the given
// essential contact on a project using the supplied *essentialcontacts.Service. Prefer this variant
// in unit tests where the service is backed by an httptest fake server (see
// essentialcontacts_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetEssentialContactAttrsWithClient(ctx context.Context, service *essentialcontacts.Service, projectID string, contactID string) (*essentialcontacts.GoogleCloudEssentialcontactsV1Contact, error) {
	name := fmt.Sprintf("projects/%s/contacts/%s", projectID, contactID)

	contact, err := service.Projects.Contacts.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the essential contact %s does not exist on project %s", contactID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for essential contact %s on project %s: %w", contactID, projectID, err)
	}

	return contact, nil
}

// NewEssentialContactsServiceE creates an Essential Contacts service authenticated the same way
// every other client in this module is.
// The ctx parameter supports cancellation and timeouts.
func NewEssentialContactsServiceE(t testing.TestingT, ctx context.Context) (*essentialcontacts.Service, error) {
	return essentialcontacts.NewService(ctx, append(withOptions(), option.WithScopes(essentialcontacts.CloudPlatformScope))...)
}
