package gcp_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gruntwork-io/terratest/modules/gcp/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/essentialcontacts/v1"
	"google.golang.org/api/option"
)

// newFakeEssentialContactsService points a real Essential Contacts client at a local test server, so the Google transport is
// exercised rather than a hand-written stand-in for a type we do not own.
func newFakeEssentialContactsService(t *testing.T, handler http.Handler) *essentialcontacts.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := essentialcontacts.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetEssentialContactAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a contact the terraform-google-
	// management module created, not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Contains(t, r.URL.Path, "/projects/gw-library-test-project/contacts/7", "unexpected path")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/123/contacts/7","email":"gw-lib-test@example.com","languageTag":"en-GB","notificationCategorySubscriptions":["TECHNICAL"],"validationState":"VALID"}`))
	})

	contact, err := gcp.GetEssentialContactAttrsWithClient(context.Background(), newFakeEssentialContactsService(t, handler), "gw-library-test-project", "7")
	require.NoError(t, err)

	assert.Equal(t, "gw-lib-test@example.com", contact.Email)
	assert.Equal(t, "en-GB", contact.LanguageTag)
	assert.Equal(t, []string{"TECHNICAL"}, contact.NotificationCategorySubscriptions)
}

func TestGetEssentialContactAttrsWithClientMissingContact(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	// The error names the contact and everything that identifies it, and each of those is a value no
	// other part of the message contains, or its check could not fail.
	_, err := gcp.GetEssentialContactAttrsWithClient(context.Background(), newFakeEssentialContactsService(t, handler), "gw-library-test-project", "gone")
	require.ErrorContains(t, err, "does not exist")
	require.ErrorContains(t, err, "gone")
	require.ErrorContains(t, err, "gw-library-test-project")
}
