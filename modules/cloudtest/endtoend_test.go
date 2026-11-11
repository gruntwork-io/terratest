package cloudtest

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	gcpinv "github.com/james00012/cloud-inventory/gcp"
	"google.golang.org/api/option"
	runv1 "google.golang.org/api/run/v1"
	runv2 "google.golang.org/api/run/v2"
)

// realProvider wires the library's own Cloud Run type to a fake server, so these
// tests drive its real create, get, list and delete rather than a stand-in for
// them. It is the closest thing to a live run that needs no credentials.
func realProvider(t *testing.T, handler http.HandlerFunc) *gcpinv.CloudRunServices {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	v1, err := runv1.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	if err != nil {
		t.Fatalf("v1 client: %v", err)
	}
	v2, err := runv2.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	if err != nil {
		t.Fatalf("v2 client: %v", err)
	}

	return gcpinv.NewCloudRunServicesWithClients(v1, v2)
}

func TestCreateAndCleanupDriveTheLibraryEndToEnd(t *testing.T) {
	// The fake holds the service's existence, so the assertions below are about
	// what is left on the server rather than which requests were sent.
	exists := false

	provider := realProvider(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		done := &runv2.GoogleLongrunningOperation{Name: "operations/one", Done: true}

		switch {
		case r.Method == http.MethodDelete:
			exists = false
			_ = json.NewEncoder(w).Encode(done)
		case r.Method == http.MethodPost && !strings.HasSuffix(r.URL.Path, ":wait"):
			exists = true
			_ = json.NewEncoder(w).Encode(done)
		case r.Method == http.MethodPost:
			_ = json.NewEncoder(w).Encode(done)
		case !exists:
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]any{
				"error": map[string]any{"code": 404, "message": "service not found"},
			})
		default:
			_ = json.NewEncoder(w).Encode(&runv2.GoogleCloudRunV2Service{
				Name:       "projects/p/locations/us-central1/services/e2e-api",
				CreateTime: "2026-11-01T00:00:00Z",
				Uri:        "https://e2e-api.example.run.app",
				Template: &runv2.GoogleCloudRunV2RevisionTemplate{
					Containers: []*runv2.GoogleCloudRunV2Container{{Image: "gcr.io/p/api:1"}},
				},
			})
		}
	})

	rec := &recorder{}

	made := Create(rec, context.Background(), provider, gcpinv.Scope{ProjectID: "p"},
		gcpinv.CloudRunServiceSpec{
			ServiceID: "e2e-api",
			Location:  "us-central1",
			GoogleCloudRunV2Service: &runv2.GoogleCloudRunV2Service{
				Template: &runv2.GoogleCloudRunV2RevisionTemplate{
					Containers: []*runv2.GoogleCloudRunV2Container{{Image: "gcr.io/p/api:1"}},
				},
			},
		})

	if rec.failed {
		t.Fatalf("create failed the test: %v %v", rec.errorf, rec.fatalf)
	}

	// The value handed back is the API's own struct, reached through the
	// library's type rather than modelled by it.
	if made.Uri != "https://e2e-api.example.run.app" {
		t.Fatalf("got uri %q", made.Uri)
	}
	if made.Template == nil || made.Template.Containers[0].Image != "gcr.io/p/api:1" {
		t.Fatal("the revision template did not survive the round trip")
	}
	if made.ResourceID() != gcpinv.ServiceID("p", "us-central1", "e2e-api") {
		t.Fatalf("got id %q", made.ResourceID())
	}
	if !exists {
		t.Fatal("the service was gone before the test finished")
	}

	rec.runCleanups()

	if exists {
		t.Fatal("the cleanup left the service behind")
	}
	if rec.failed {
		t.Fatalf("cleanup reported a failure: %v", rec.errorf)
	}

	// The library reads an already deleted service as gone rather than as an
	// error, so a second cleanup is safe.
	rec.runCleanups()
	if rec.failed {
		t.Fatalf("a repeated cleanup failed: %v", rec.errorf)
	}
}

func TestListFailsTheTestWhenTheLibraryCouldNotReadALocation(t *testing.T) {
	provider := realProvider(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if strings.HasSuffix(r.URL.Path, "/locations") {
			_ = json.NewEncoder(w).Encode(&runv1.ListLocationsResponse{
				Locations: []*runv1.Location{
					{LocationId: "us-central1"},
					{LocationId: "europe-west1"},
				},
			})
			return
		}

		// One location answers, the other refuses, which is the partial read
		// a leak check must never read as an empty project.
		if strings.Contains(r.URL.Path, "europe-west1") {
			w.WriteHeader(http.StatusForbidden)
			_ = json.NewEncoder(w).Encode(map[string]any{
				"error": map[string]any{"code": 403, "message": "caller lacks permission"},
			})
			return
		}

		_ = json.NewEncoder(w).Encode(&runv2.GoogleCloudRunV2ListServicesResponse{})
	})

	rec := &recorder{}

	List(rec, context.Background(), provider, gcpinv.Scope{ProjectID: "p"})

	if len(rec.fatalf) != 1 || !strings.Contains(rec.fatalf[0], "proves nothing") {
		t.Fatalf("a refused location did not fail the test: %v %v", rec.fatalf, rec.errorf)
	}
}
