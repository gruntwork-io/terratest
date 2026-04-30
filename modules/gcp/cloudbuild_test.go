package gcp_test

import (
	"context"
	"net"
	"testing"

	cloudbuild "cloud.google.com/go/cloudbuild/apiv1/v2"
	cloudbuildpb "cloud.google.com/go/cloudbuild/apiv1/v2/cloudbuildpb"
	"github.com/gruntwork-io/terratest/modules/gcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

// fakeCloudBuildServer only implements the methods the terratest *WithClient helpers actually
// call. CreateBuild (long-running op) is left to the build-tagged integration test.
type fakeCloudBuildServer struct {
	cloudbuildpb.UnimplementedCloudBuildServer

	getBuild   func(*cloudbuildpb.GetBuildRequest) *cloudbuildpb.Build
	listBuilds func(*cloudbuildpb.ListBuildsRequest) *cloudbuildpb.ListBuildsResponse
}

func (f *fakeCloudBuildServer) GetBuild(_ context.Context, req *cloudbuildpb.GetBuildRequest) (*cloudbuildpb.Build, error) {
	return f.getBuild(req), nil
}

func (f *fakeCloudBuildServer) ListBuilds(_ context.Context, req *cloudbuildpb.ListBuildsRequest) (*cloudbuildpb.ListBuildsResponse, error) {
	return f.listBuilds(req), nil
}

// newFakeCloudBuildClient runs `srv` on an in-memory bufconn gRPC server and returns a client
// pointed at it. No credentials, no network.
func newFakeCloudBuildClient(t *testing.T, srv *fakeCloudBuildServer) *cloudbuild.Client {
	t.Helper()

	lis := bufconn.Listen(1024 * 1024)
	grpcServer := grpc.NewServer()
	cloudbuildpb.RegisterCloudBuildServer(grpcServer, srv)

	go func() { _ = grpcServer.Serve(lis) }()

	t.Cleanup(func() { grpcServer.Stop(); _ = lis.Close() })

	conn, err := grpc.NewClient("passthrough:///bufnet",
		grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) { return lis.DialContext(ctx) }),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)
	t.Cleanup(func() { _ = conn.Close() })

	client, err := cloudbuild.NewClient(context.Background(), option.WithGRPCConn(conn))
	require.NoError(t, err)
	t.Cleanup(func() { _ = client.Close() })

	return client
}

func TestGetBuildWithClient(t *testing.T) {
	t.Parallel()

	client := newFakeCloudBuildClient(t, &fakeCloudBuildServer{
		getBuild: func(req *cloudbuildpb.GetBuildRequest) *cloudbuildpb.Build {
			return &cloudbuildpb.Build{Id: req.GetId(), ProjectId: req.GetProjectId(), Status: cloudbuildpb.Build_SUCCESS}
		},
	})

	got, err := gcp.GetBuildWithClient(context.Background(), client, "p", "b1")
	require.NoError(t, err)
	assert.Equal(t, "b1", got.GetId())
	assert.Equal(t, cloudbuildpb.Build_SUCCESS, got.GetStatus())
}

func TestGetBuildsWithClient(t *testing.T) {
	t.Parallel()

	client := newFakeCloudBuildClient(t, &fakeCloudBuildServer{
		listBuilds: func(req *cloudbuildpb.ListBuildsRequest) *cloudbuildpb.ListBuildsResponse {
			assert.Equal(t, "p", req.GetProjectId())

			return &cloudbuildpb.ListBuildsResponse{Builds: []*cloudbuildpb.Build{{Id: "a"}, {Id: "b"}}}
		},
	})

	got, err := gcp.GetBuildsWithClient(context.Background(), client, "p")
	require.NoError(t, err)
	require.Len(t, got, 2)
}

func TestGetBuildsForTriggerWithClient(t *testing.T) {
	t.Parallel()

	client := newFakeCloudBuildClient(t, &fakeCloudBuildServer{
		listBuilds: func(_ *cloudbuildpb.ListBuildsRequest) *cloudbuildpb.ListBuildsResponse {
			return &cloudbuildpb.ListBuildsResponse{Builds: []*cloudbuildpb.Build{
				{Id: "a", BuildTriggerId: "match"},
				{Id: "b", BuildTriggerId: "other"},
				{Id: "c", BuildTriggerId: "match"},
			}}
		},
	})

	got, err := gcp.GetBuildsForTriggerWithClient(context.Background(), client, "p", "match")
	require.NoError(t, err)
	require.Len(t, got, 2)
	assert.ElementsMatch(t, []string{"a", "c"}, []string{got[0].GetId(), got[1].GetId()})
}
