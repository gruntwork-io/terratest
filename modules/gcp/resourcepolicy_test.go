package gcp_test

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/gcp/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFetchResourcePolicyWithClient(t *testing.T) {
	t.Parallel()

	// The values are the ones the terraform-google-compute resource policy module sets, because the point of reading settings back
	// is asserting a module configured the resource policy it was asked for.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/regions/us-central1/resourcePolicies/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"gw-library-test","description":"created by terratest","status":"READY","snapshotSchedulePolicy":{"schedule":{"dailySchedule":{"daysInCycle":1,"startTime":"04:00"}},"retentionPolicy":{"maxRetentionDays":2}}}`))
	})

	policy, err := gcp.FetchResourcePolicyWithClient(context.Background(), newFakeComputeService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "gw-library-test", policy.Name)
	assert.Equal(t, "created by terratest", policy.Description)
	require.NotNil(t, policy.SnapshotSchedulePolicy)
	require.NotNil(t, policy.SnapshotSchedulePolicy.RetentionPolicy)
	assert.Equal(t, int64(2), policy.SnapshotSchedulePolicy.RetentionPolicy.MaxRetentionDays)
	require.NotNil(t, policy.SnapshotSchedulePolicy.Schedule.DailySchedule)
	assert.Equal(t, "04:00", policy.SnapshotSchedulePolicy.Schedule.DailySchedule.StartTime)
}
