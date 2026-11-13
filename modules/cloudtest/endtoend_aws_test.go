package cloudtest

import (
	"context"
	"strings"
	"sync"
	"testing"
	"time"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	ec2sdk "github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/smithy-go"
	awsinv "github.com/james00012/cloud-inventory/aws"
)

// ec2World is an account: regions, the instances in each, and the regions that
// refuse. The library's own EC2 type runs against it, so these tests exercise
// its region routing, its fan-out and its terminate wait rather than a
// stand-in for them.
type ec2World struct {
	mu         sync.Mutex
	regions    []string
	refuse     map[string]string
	instances  map[string][]*ec2types.Instance
	next       int
	ran        []string
	terminated []string

	// keepRunning makes terminate a no-op, which is how a stuck delete looks.
	keepRunning bool
}

type ec2Client struct {
	world  *ec2World
	region string
}

func (w *ec2World) clientFor(region string) awsinv.EC2API {
	return &ec2Client{world: w, region: region}
}

func (c *ec2Client) deny() error {
	if code, refused := c.world.refuse[c.region]; refused {
		return &smithy.GenericAPIError{Code: code, Message: "refused by the fake"}
	}
	return nil
}

func (c *ec2Client) DescribeRegions(ctx context.Context, in *ec2sdk.DescribeRegionsInput, opts ...func(*ec2sdk.Options)) (*ec2sdk.DescribeRegionsOutput, error) {
	c.world.mu.Lock()
	defer c.world.mu.Unlock()

	out := &ec2sdk.DescribeRegionsOutput{}
	for _, region := range c.world.regions {
		out.Regions = append(out.Regions, ec2types.Region{RegionName: awssdk.String(region)})
	}
	return out, nil
}

func (c *ec2Client) DescribeInstances(ctx context.Context, in *ec2sdk.DescribeInstancesInput, opts ...func(*ec2sdk.Options)) (*ec2sdk.DescribeInstancesOutput, error) {
	if err := c.deny(); err != nil {
		return nil, err
	}

	c.world.mu.Lock()
	defer c.world.mu.Unlock()

	wanted := map[string]bool{}
	for _, id := range in.InstanceIds {
		wanted[id] = true
	}

	var found []ec2types.Instance
	for _, instance := range c.world.instances[c.region] {
		if len(wanted) > 0 && !wanted[*instance.InstanceId] {
			continue
		}
		// An instance reaches running on the read after its launch, which is
		// what gives a wait something to wait for.
		if instance.State.Name == ec2types.InstanceStateNamePending {
			instance.State.Name = ec2types.InstanceStateNameRunning
			found = append(found, *instance)
			continue
		}
		found = append(found, *instance)
	}
	if len(found) == 0 {
		return &ec2sdk.DescribeInstancesOutput{}, nil
	}
	return &ec2sdk.DescribeInstancesOutput{
		Reservations: []ec2types.Reservation{{Instances: found}},
	}, nil
}

func (c *ec2Client) RunInstances(ctx context.Context, in *ec2sdk.RunInstancesInput, opts ...func(*ec2sdk.Options)) (*ec2sdk.RunInstancesOutput, error) {
	if err := c.deny(); err != nil {
		return nil, err
	}

	c.world.mu.Lock()
	defer c.world.mu.Unlock()

	c.world.next++
	id := "i-000000" + string(rune('0'+c.world.next))
	launched := time.Date(2026, 11, 1, 0, 0, 0, 0, time.UTC)

	var tags []ec2types.Tag
	for _, spec := range in.TagSpecifications {
		tags = append(tags, spec.Tags...)
	}

	instance := &ec2types.Instance{
		InstanceId:   awssdk.String(id),
		InstanceType: in.InstanceType,
		ImageId:      in.ImageId,
		LaunchTime:   &launched,
		State:        &ec2types.InstanceState{Name: ec2types.InstanceStateNamePending},
		Tags:         tags,
	}
	c.world.instances[c.region] = append(c.world.instances[c.region], instance)
	c.world.ran = append(c.world.ran, c.region+"/"+id)

	return &ec2sdk.RunInstancesOutput{Instances: []ec2types.Instance{*instance}}, nil
}

func (c *ec2Client) TerminateInstances(ctx context.Context, in *ec2sdk.TerminateInstancesInput, opts ...func(*ec2sdk.Options)) (*ec2sdk.TerminateInstancesOutput, error) {
	if err := c.deny(); err != nil {
		return nil, err
	}

	c.world.mu.Lock()
	defer c.world.mu.Unlock()

	for _, id := range in.InstanceIds {
		for _, instance := range c.world.instances[c.region] {
			if *instance.InstanceId != id {
				continue
			}
			if c.world.keepRunning {
				continue
			}
			// Straight to terminated, so the library's wait is satisfied by
			// its first poll and the test does not sleep.
			instance.State.Name = ec2types.InstanceStateNameTerminated
			c.world.terminated = append(c.world.terminated, c.region+"/"+id)
		}
	}
	return &ec2sdk.TerminateInstancesOutput{}, nil
}

func (w *ec2World) stateOf(region, id string) ec2types.InstanceStateName {
	w.mu.Lock()
	defer w.mu.Unlock()

	for _, instance := range w.instances[region] {
		if *instance.InstanceId == id {
			return instance.State.Name
		}
	}
	return ""
}

func newEC2World(regions ...string) *ec2World {
	return &ec2World{
		regions:   regions,
		refuse:    map[string]string{},
		instances: map[string][]*ec2types.Instance{},
	}
}

func TestCreateWaitAndCleanupDriveTheAwsLibraryEndToEnd(t *testing.T) {
	world := newEC2World("us-east-1", "eu-west-1")
	provider := awsinv.NewEC2InstancesWithClients(world.clientFor)
	rec := &recorder{}

	made := Create(rec, context.Background(), provider, awsinv.Scope{Regions: []string{"eu-west-1"}},
		awsinv.EC2InstanceSpec{
			Region: "eu-west-1",
			RunInstancesInput: &ec2sdk.RunInstancesInput{
				ImageId:      awssdk.String("ami-12345678"),
				InstanceType: ec2types.InstanceTypeT3Micro,
				TagSpecifications: []ec2types.TagSpecification{{
					ResourceType: ec2types.ResourceTypeInstance,
					Tags:         []ec2types.Tag{{Key: awssdk.String("Name"), Value: awssdk.String("e2e-web")}},
				}},
			},
		})

	if rec.failed {
		t.Fatalf("create failed the test: %v %v", rec.errorf, rec.fatalf)
	}

	// The id carries the region, because an instance id is only unique with it,
	// and the constructor is what a caller uses to rebuild it.
	instanceID := strings.TrimPrefix(string(made.ResourceID()), "eu-west-1/")
	if made.ResourceID() != awsinv.InstanceID("eu-west-1", instanceID) {
		t.Fatalf("got id %q", made.ResourceID())
	}

	// The contract's five methods read the SDK struct underneath, including the
	// name that EC2 keeps as a tag rather than a field.
	if made.DisplayName() != "e2e-web" {
		t.Fatalf("got display name %q", made.DisplayName())
	}
	if made.CreatedAt().IsZero() {
		t.Fatal("CreatedAt did not read LaunchTime")
	}
	if made.InstanceType != ec2types.InstanceTypeT3Micro {
		t.Fatalf("got instance type %q", made.InstanceType)
	}

	// The wait reads through the library's own Get, which routes by the
	// region inside the id, and compares against the SDK's own constant.
	ready := WaitUntil(rec, context.Background(), provider, made.ResourceID(), "running",
		func(i awsinv.EC2Instance) bool {
			return i.State != nil && i.State.Name == ec2types.InstanceStateNameRunning
		}, 5, time.Millisecond)

	if ready.State.Name != ec2types.InstanceStateNameRunning {
		t.Fatalf("waited and got state %q", ready.State.Name)
	}
	if world.stateOf("us-east-1", instanceID) != "" {
		t.Fatal("the instance was made in the wrong region")
	}

	rec.runCleanups()

	if got := world.stateOf("eu-west-1", instanceID); got != ec2types.InstanceStateNameTerminated {
		t.Fatalf("the cleanup left the instance %q", got)
	}
	if len(world.terminated) != 1 || world.terminated[0] != "eu-west-1/"+instanceID {
		t.Fatalf("terminated the wrong thing: %v", world.terminated)
	}
	if rec.failed {
		t.Fatalf("cleanup reported a failure: %v", rec.errorf)
	}
}

func TestAssertNoneReportsALeakThroughTheAwsLibrary(t *testing.T) {
	world := newEC2World("us-east-1", "eu-west-1")
	provider := awsinv.NewEC2InstancesWithClients(world.clientFor)

	// One instance per region, only one of them a leftover from a test run.
	for region, name := range map[string]string{"us-east-1": "prod-db", "eu-west-1": "e2e-stale"} {
		rec := &recorder{}
		Create(rec, context.Background(), provider, awsinv.Scope{},
			awsinv.EC2InstanceSpec{
				Region: region,
				RunInstancesInput: &ec2sdk.RunInstancesInput{
					TagSpecifications: []ec2types.TagSpecification{{
						ResourceType: ec2types.ResourceTypeInstance,
						Tags:         []ec2types.Tag{{Key: awssdk.String("Name"), Value: awssdk.String(name)}},
					}},
				},
			})
		if rec.failed {
			t.Fatalf("seeding %s failed: %v", region, rec.fatalf)
		}
	}

	rec := &recorder{}

	// No regions named, so this fans out over every region the account offers.
	AssertNone(rec, context.Background(), provider, awsinv.Scope{},
		func(i awsinv.EC2Instance) bool { return strings.HasPrefix(i.DisplayName(), "e2e-") })

	if len(rec.errorf) != 1 {
		t.Fatalf("expected exactly the one leak, got %v", rec.errorf)
	}
	if !strings.Contains(rec.errorf[0], "left behind") {
		t.Fatalf("did not report it as a leak: %v", rec.errorf)
	}
}

func TestListFailsTheTestWhenARegionRefuses(t *testing.T) {
	world := newEC2World("us-east-1", "eu-west-1")
	world.refuse["eu-west-1"] = "AuthFailure"
	provider := awsinv.NewEC2InstancesWithClients(world.clientFor)

	rec := &recorder{}

	List(rec, context.Background(), provider, awsinv.Scope{})

	if len(rec.fatalf) != 1 || !strings.Contains(rec.fatalf[0], "proves nothing") {
		t.Fatalf("a refused region did not fail the test: %v %v", rec.fatalf, rec.errorf)
	}
}

func TestListFailsTheTestWhenTheAccountDoesNotOfferANamedRegion(t *testing.T) {
	world := newEC2World("us-east-1")
	provider := awsinv.NewEC2InstancesWithClients(world.clientFor)

	rec := &recorder{}

	// A region the account does not offer is a question nobody answered, so a
	// leak check over it must not read as empty.
	List(rec, context.Background(), provider, awsinv.Scope{Regions: []string{"us-east-1", "ap-south-1"}})

	if len(rec.fatalf) != 1 || !strings.Contains(rec.fatalf[0], "proves nothing") {
		t.Fatalf("an unoffered region did not fail the test: %v %v", rec.fatalf, rec.errorf)
	}
}

func TestCleanupTimeoutBoundsAProviderThatWaitsForever(t *testing.T) {
	// The library waits for an instance to reach terminated, for minutes. A
	// fake that never gets there stands in for a stuck delete.
	world := newEC2World("us-east-1")
	world.keepRunning = true
	provider := awsinv.NewEC2InstancesWithClients(world.clientFor)

	was := CleanupTimeout
	CleanupTimeout = 50 * time.Millisecond
	t.Cleanup(func() { CleanupTimeout = was })

	rec := &recorder{}
	Create(rec, context.Background(), provider, awsinv.Scope{},
		awsinv.EC2InstanceSpec{Region: "us-east-1", RunInstancesInput: &ec2sdk.RunInstancesInput{}})
	if rec.failed {
		t.Fatalf("create failed: %v", rec.fatalf)
	}

	started := time.Now()
	rec.runCleanups()
	took := time.Since(started)

	if len(rec.errorf) != 1 || !strings.Contains(rec.errorf[0], "deleting") {
		t.Fatalf("a stuck delete did not report: %v", rec.errorf)
	}
	if took > 2*time.Second {
		t.Fatalf("the cleanup timeout did not bound the provider's wait: took %s", took)
	}
}
