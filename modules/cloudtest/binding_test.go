package cloudtest

import (
	"context"
	"testing"
	"time"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	ec2sdk "github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	inv "github.com/james00012/cloud-inventory"
	awsinv "github.com/james00012/cloud-inventory/aws"
	gcpinv "github.com/james00012/cloud-inventory/gcp"
	runv2 "google.golang.org/api/run/v2"
)

// The wrapper is written once. These assertions are the whole adoption cost per
// resource type: the real providers already satisfy it, with no code in this
// package naming EC2 or Cloud Run.
var (
	_ Managed[awsinv.EC2Instance, awsinv.EC2InstanceSpec, awsinv.Scope] = (*awsinv.EC2Instances)(nil)
	_ inv.Getter[awsinv.EC2Instance]                                    = (*awsinv.EC2Instances)(nil)
	_ inv.Lister[awsinv.EC2Instance, awsinv.Scope]                      = (*awsinv.EC2Instances)(nil)

	_ Managed[gcpinv.CloudRunService, gcpinv.CloudRunServiceSpec, gcpinv.Scope] = (*gcpinv.CloudRunServices)(nil)
	_ inv.Getter[gcpinv.CloudRunService]                                        = (*gcpinv.CloudRunServices)(nil)
	_ inv.Lister[gcpinv.CloudRunService, gcpinv.Scope]                          = (*gcpinv.CloudRunServices)(nil)
)

// TestWhatAUserWritesAWS is what a Terratest user's test looks like against the
// wrapper. It does not run without credentials; it is here because it has to
// compile, which is the claim being made.
func TestWhatAUserWritesAWS(t *testing.T) {
	t.Skip("needs AWS credentials; compiled, not run")

	ctx := context.Background()
	instances := awsinv.NewEC2Instances(awssdk.Config{})

	// One create call. Cleanup is registered by Create, and the spec carries
	// the SDK's own RunInstancesInput, so there is no Terratest vocabulary.
	instance := Create(t, ctx, instances, awsinv.Scope{Regions: []string{"us-east-1"}},
		awsinv.EC2InstanceSpec{
			Region: "us-east-1",
			RunInstancesInput: &ec2sdk.RunInstancesInput{
				ImageId:      awssdk.String("ami-12345678"),
				InstanceType: ec2types.InstanceTypeT3Micro,
				MinCount:     awssdk.Int32(1),
				MaxCount:     awssdk.Int32(1),
			},
		})

	// The value handed back is the AWS SDK struct, so the assertion reads the
	// SDK's own fields and compares against its own enum constants.
	ready := WaitUntil(t, ctx, instances, instance.ResourceID(), "running",
		func(i awsinv.EC2Instance) bool {
			return i.State != nil && i.State.Name == ec2types.InstanceStateNameRunning
		}, 30, 10*time.Second)

	if ready.PublicIpAddress == nil {
		t.Fatal("expected a public address")
	}

	// Reading a resource the test did not create, which is what most of
	// Terratest does today. The id comes from the provider's constructor, so a
	// caller never writes the format by hand.
	byID := Get(t, ctx, instances, awsinv.InstanceID("us-east-1", "i-0abc123"))
	if byID.InstanceType != ec2types.InstanceTypeT3Micro {
		t.Fatalf("got %s", byID.InstanceType)
	}
}

// TestWhatAUserWritesGCP is the same test against a different cloud, through the
// same functions.
func TestWhatAUserWritesGCP(t *testing.T) {
	t.Skip("needs GCP credentials; compiled, not run")

	ctx := context.Background()
	services, err := gcpinv.NewCloudRunServices(ctx)
	if err != nil {
		t.Fatal(err)
	}

	service := Create(t, ctx, services, gcpinv.Scope{ProjectID: "my-project"},
		gcpinv.CloudRunServiceSpec{
			ServiceID: "e2e-api",
			Location:  "us-central1",
			GoogleCloudRunV2Service: &runv2.GoogleCloudRunV2Service{
				Template: &runv2.GoogleCloudRunV2RevisionTemplate{
					Containers: []*runv2.GoogleCloudRunV2Container{{Image: "gcr.io/cloudrun/hello"}},
				},
			},
		})

	// Again the provider's own struct, so Uri and Template are the real fields.
	if service.Uri == "" {
		t.Fatal("expected a service URL")
	}

	// This line is the replacement for GetCloudRunServiceAttrs(t, ctx,
	// projectID, region, serviceID): the same three arguments a caller knows,
	// and the provider builds the resource path.
	attrs := Get(t, ctx, services, gcpinv.ServiceID("my-project", "us-central1", "e2e-api"))
	if attrs.Template == nil {
		t.Fatal("expected a revision template")
	}

	// A leak check that fails rather than passing quietly when a location
	// could not be read.
	AssertNone(t, ctx, services, gcpinv.Scope{ProjectID: "my-project"},
		func(s gcpinv.CloudRunService) bool { return s.DisplayName() == "e2e-api-stale" })
}
