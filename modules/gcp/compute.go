package gcp

import (
	"context"
	"fmt"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/gruntwork-io/terratest/modules/retry"
	"google.golang.org/api/compute/v1"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/testing"
	"golang.org/x/oauth2/google"
)

// Corresponds to a GCP Compute Instance (https://cloud.google.com/compute/docs/instances/)
type Instance struct {
	projectID string
	*compute.Instance
}

// Corresponds to a GCP Image (https://cloud.google.com/compute/docs/images)
type Image struct {
	projectID string
	*compute.Image
}

// Corresponds to a GCP Zonal Instance Group (https://cloud.google.com/compute/docs/instance-groups/)
type ZonalInstanceGroup struct {
	projectID string
	*compute.InstanceGroup
}

// Corresponds to a GCP Regional Instance Group (https://cloud.google.com/compute/docs/instance-groups/)
type RegionalInstanceGroup struct {
	projectID string
	*compute.InstanceGroup
}

type InstanceGroup interface {
	GetInstanceIds(t testing.TestingT) []string
	GetInstanceIdsE(t testing.TestingT) ([]string, error)
}

// FetchInstance queries GCP to return an instance of the (GCP Compute) Instance type
func FetchInstance(t testing.TestingT, projectID string, name string) *Instance {
	instance, err := FetchInstanceE(t, projectID, name)
	if err != nil {
		t.Fatal(err)
	}

	return instance
}

// FetchInstance queries GCP to return an instance of the (GCP Compute) Instance type
func FetchInstanceE(t testing.TestingT, projectID string, name string) (*Instance, error) {
	logger.Logf(t, "Getting Compute Instance %s", name)

	ctx := context.Background()
	service, err := NewComputeServiceE(t)
	if err != nil {
		t.Fatal(err)
	}

	// If we want to fetch an Instance without knowing its Zone, we have to query GCP for all Instances in the project
	// and match on name.
	instanceAggregatedList, err := service.Instances.AggregatedList(projectID).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("Instances.AggregatedList(%s) got error: %v", projectID, err)
	}

	for _, instanceList := range instanceAggregatedList.Items {
		for _, instance := range instanceList.Instances {
			if name == instance.Name {
				return &Instance{projectID, instance}, nil
			}
		}
	}

	return nil, fmt.Errorf("Compute Instance %s could not be found in project %s", name, projectID)
}

// FetchImage queries GCP to return a new instance of the (GCP Compute) Image type
func FetchImage(t testing.TestingT, projectID string, name string) *Image {
	image, err := FetchImageE(t, projectID, name)
	if err != nil {
		t.Fatal(err)
	}

	return image
}

// FetchImage queries GCP to return a new instance of the (GCP Compute) Image type
func FetchImageE(t testing.TestingT, projectID string, name string) (*Image, error) {
	logger.Logf(t, "Getting Image %s", name)

	ctx := context.Background()
	service, err := NewComputeServiceE(t)
	if err != nil {
		return nil, err
	}

	req := service.Images.Get(projectID, name)
	image, err := req.Context(ctx).Do()
	if err != nil {
		return nil, err
	}

	return &Image{projectID, image}, nil
}

// FetchRegionalInstanceGroup queries GCP to return a new instance of the Regional Instance Group type
func FetchRegionalInstanceGroup(t testing.TestingT, projectID string, region string, name string) *RegionalInstanceGroup {
	instanceGroup, err := FetchRegionalInstanceGroupE(t, projectID, region, name)
	if err != nil {
		t.Fatal(err)
	}

	return instanceGroup
}

// FetchRegionalInstanceGroup queries GCP to return a new instance of the Regional Instance Group type
func FetchRegionalInstanceGroupE(t testing.TestingT, projectID string, region string, name string) (*RegionalInstanceGroup, error) {
	logger.Logf(t, "Getting Regional Instance Group %s", name)

	ctx := context.Background()
	service, err := NewComputeServiceE(t)
	if err != nil {
		return nil, err
	}

	req := service.RegionInstanceGroups.Get(projectID, region, name)
	instanceGroup, err := req.Context(ctx).Do()
	if err != nil {
		return nil, err
	}

	return &RegionalInstanceGroup{projectID, instanceGroup}, nil
}

// FetchZonalInstanceGroup queries GCP to return a new instance of the Regional Instance Group type
func FetchZonalInstanceGroup(t testing.TestingT, projectID string, zone string, name string) *ZonalInstanceGroup {
	instanceGroup, err := FetchZonalInstanceGroupE(t, projectID, zone, name)
	if err != nil {
		t.Fatal(err)
	}

	return instanceGroup
}

// FetchZonalInstanceGroup queries GCP to return a new instance of the Regional Instance Group type
func FetchZonalInstanceGroupE(t testing.TestingT, projectID string, zone string, name string) (*ZonalInstanceGroup, error) {
	logger.Logf(t, "Getting Zonal Instance Group %s", name)

	ctx := context.Background()
	service, err := NewComputeServiceE(t)
	if err != nil {
		return nil, err
	}

	req := service.InstanceGroups.Get(projectID, zone, name)
	instanceGroup, err := req.Context(ctx).Do()
	if err != nil {
		return nil, err
	}

	return &ZonalInstanceGroup{projectID, instanceGroup}, nil
}

// GetPublicIP gets the public IP address of the given Compute Instance.
func (i *Instance) GetPublicIp(t testing.TestingT) string {
	ip, err := i.GetPublicIpE(t)
	if err != nil {
		t.Fatal(err)
	}
	return ip
}

// GetPublicIpE gets the public IP address of the given Compute Instance.
func (i *Instance) GetPublicIpE(t testing.TestingT) (string, error) {
	// If there are no accessConfigs specified, then this instance will have no external internet access:
	// https://cloud.google.com/compute/docs/reference/rest/v1/instances.
	if len(i.NetworkInterfaces[0].AccessConfigs) == 0 {
		return "", fmt.Errorf("Attempted to get public IP of Compute Instance %s, but that Compute Instance does not have a public IP address", i.Name)
	}

	ip := i.NetworkInterfaces[0].AccessConfigs[0].NatIP

	return ip, nil
}

// GetLabels returns all the tags for the given Compute Instance.
func (i *Instance) GetLabels(t testing.TestingT) map[string]string {
	return i.Labels
}

// GetZone returns the Zone in which the Compute Instance is located.
func (i *Instance) GetZone(t testing.TestingT) string {
	return ZoneUrlToZone(i.Zone)
}

// SetLabels adds the tags to the given Compute Instance.
func (i *Instance) SetLabels(t testing.TestingT, labels map[string]string) {
	err := i.SetLabelsE(t, labels)
	if err != nil {
		t.Fatal(err)
	}
}

// SetLabelsE adds the tags to the given Compute Instance.
func (i *Instance) SetLabelsE(t testing.TestingT, labels map[string]string) error {
	logger.Logf(t, "Adding labels to instance %s in zone %s", i.Name, i.Zone)

	ctx := context.Background()
	service, err := NewComputeServiceE(t)
	if err != nil {
		return err
	}

	req := compute.InstancesSetLabelsRequest{Labels: labels, LabelFingerprint: i.LabelFingerprint}
	if _, err := service.Instances.SetLabels(i.projectID, i.GetZone(t), i.Name, &req).Context(ctx).Do(); err != nil {
		return fmt.Errorf("Instances.SetLabels(%s) got error: %v", i.Name, err)
	}

	return nil
}

// GetMetadata gets the given Compute Instance's metadata
func (i *Instance) GetMetadata(t testing.TestingT) []*compute.MetadataItems {
	return i.Metadata.Items
}

// SetMetadata sets the given Compute Instance's metadata
func (i *Instance) SetMetadata(t testing.TestingT, metadata map[string]string) {
	err := i.SetMetadataE(t, metadata)
	if err != nil {
		t.Fatal(err)
	}
}

// SetLabelsE adds the given metadata map to the existing metadata of the given Compute Instance.
func (i *Instance) SetMetadataE(t testing.TestingT, metadata map[string]string) error {
	logger.Logf(t, "Adding metadata to instance %s in zone %s", i.Name, i.Zone)

	ctx := context.Background()
	service, err := NewInstancesServiceE(t)
	if err != nil {
		return err
	}

	metadataItems := newMetadata(t, i.Metadata, metadata)
	req := service.SetMetadata(i.projectID, i.GetZone(t), i.Name, metadataItems)
	if _, err := req.Context(ctx).Do(); err != nil {
		return fmt.Errorf("Instances.SetMetadata(%s) got error: %v", i.Name, err)
	}

	return nil
}

// newMetadata takes in a Compute Instance's existing metadata plus a new set of key-value pairs and returns an updated
// metadata object.
func newMetadata(t testing.TestingT, oldMetadata *compute.Metadata, kvs map[string]string) *compute.Metadata {
	items := []*compute.MetadataItems{}

	for key, val := range kvs {
		item := &compute.MetadataItems{
			Key:   key,
			Value: &val,
		}

		items = append(oldMetadata.Items, item)
	}

	newMetadata := &compute.Metadata{
		Fingerprint: oldMetadata.Fingerprint,
		Items:       items,
	}

	return newMetadata
}

// Add the given public SSH key to the Compute Instance. Users can SSH in with the given username.
func (i *Instance) AddSshKey(t testing.TestingT, username string, publicKey string) {
	err := i.AddSshKeyE(t, username, publicKey)
	if err != nil {
		t.Fatal(err)
	}
}

// Add the given public SSH key to the Compute Instance. Users can SSH in with the given username.
func (i *Instance) AddSshKeyE(t testing.TestingT, username string, publicKey string) error {
	logger.Logf(t, "Adding SSH Key to Compute Instance %s for username %s\n", i.Name, username)

	// We represent the key in the format required per GCP docs (https://cloud.google.com/compute/docs/instances/adding-removing-ssh-keys)
	publicKeyFormatted := strings.TrimSpace(publicKey)
	sshKeyFormatted := fmt.Sprintf("%s:%s %s", username, publicKeyFormatted, username)

	metadata := map[string]string{
		"ssh-keys": sshKeyFormatted,
	}

	err := i.SetMetadataE(t, metadata)
	if err != nil {
		return fmt.Errorf("Failed to add SSH key to Compute Instance: %s", err)
	}

	return nil
}

// DeleteImage deletes the given Compute Image.
func (i *Image) DeleteImage(t testing.TestingT) {
	err := i.DeleteImageE(t)
	if err != nil {
		t.Fatal(err)
	}
}

// DeleteImageE deletes the given Compute Image.
func (i *Image) DeleteImageE(t testing.TestingT) error {
	logger.Logf(t, "Destroying Image %s", i.Name)

	ctx := context.Background()
	service, err := NewComputeServiceE(t)
	if err != nil {
		return err
	}

	if _, err := service.Images.Delete(i.projectID, i.Name).Context(ctx).Do(); err != nil {
		return fmt.Errorf("Images.Delete(%s) got error: %v", i.Name, err)
	}

	return nil
}

// GetInstanceIds gets the IDs of Instances in the given Instance Group.
func (ig *ZonalInstanceGroup) GetInstanceIds(t testing.TestingT) []string {
	ids, err := ig.GetInstanceIdsE(t)
	if err != nil {
		t.Fatal(err)
	}
	return ids
}

// GetInstanceIdsE gets the IDs of Instances in the given Zonal Instance Group.
func (ig *ZonalInstanceGroup) GetInstanceIdsE(t testing.TestingT) ([]string, error) {
	logger.Logf(t, "Get instances for Zonal Instance Group %s", ig.Name)

	ctx := context.Background()
	service, err := NewComputeServiceE(t)
	if err != nil {
		return nil, err
	}

	requestBody := &compute.InstanceGroupsListInstancesRequest{
		InstanceState: "ALL",
	}

	instanceIDs := []string{}
	zone := ZoneUrlToZone(ig.Zone)

	req := service.InstanceGroups.ListInstances(ig.projectID, zone, ig.Name, requestBody)

	err = req.Pages(ctx, func(page *compute.InstanceGroupsListInstances) error {
		for _, instance := range page.Items {
			// For some reason service.InstanceGroups.ListInstances returns us a collection
			// with Instance URLs and we need only the Instance ID for the next call. Use
			// the path functions to chop the Instance ID off the end of the URL.
			instanceID := path.Base(instance.Instance)
			instanceIDs = append(instanceIDs, instanceID)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("InstanceGroups.ListInstances(%s) got error: %v", ig.Name, err)
	}

	return instanceIDs, nil
}

// GetInstanceIds gets the IDs of Instances in the given Regional Instance Group.
func (ig *RegionalInstanceGroup) GetInstanceIds(t testing.TestingT) []string {
	ids, err := ig.GetInstanceIdsE(t)
	if err != nil {
		t.Fatal(err)
	}
	return ids
}

// GetInstanceIdsE gets the IDs of Instances in the given Regional Instance Group.
func (ig *RegionalInstanceGroup) GetInstanceIdsE(t testing.TestingT) ([]string, error) {
	logger.Logf(t, "Get instances for Regional Instance Group %s", ig.Name)

	ctx := context.Background()

	service, err := NewComputeServiceE(t)
	if err != nil {
		return nil, err
	}

	requestBody := &compute.RegionInstanceGroupsListInstancesRequest{
		InstanceState: "ALL",
	}

	instanceIDs := []string{}
	region := RegionUrlToRegion(ig.Region)

	req := service.RegionInstanceGroups.ListInstances(ig.projectID, region, ig.Name, requestBody)

	err = req.Pages(ctx, func(page *compute.RegionInstanceGroupsListInstances) error {
		for _, instance := range page.Items {
			// For some reason service.InstanceGroups.ListInstances returns us a collection
			// with Instance URLs and we need only the Instance ID for the next call. Use
			// the path functions to chop the Instance ID off the end of the URL.
			instanceID := path.Base(instance.Instance)
			instanceIDs = append(instanceIDs, instanceID)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("InstanceGroups.ListInstances(%s) got error: %v", ig.Name, err)
	}

	return instanceIDs, nil
}

// Return a collection of Instance structs from the given Instance Group
func (ig *ZonalInstanceGroup) GetInstances(t testing.TestingT, projectId string) []*Instance {
	return getInstances(t, ig, projectId)
}

// Return a collection of Instance structs from the given Instance Group
func (ig *ZonalInstanceGroup) GetInstancesE(t testing.TestingT, projectId string) ([]*Instance, error) {
	return getInstancesE(t, ig, projectId)
}

// Return a collection of Instance structs from the given Instance Group
func (ig *RegionalInstanceGroup) GetInstances(t testing.TestingT, projectId string) []*Instance {
	return getInstances(t, ig, projectId)
}

// Return a collection of Instance structs from the given Instance Group
func (ig *RegionalInstanceGroup) GetInstancesE(t testing.TestingT, projectId string) ([]*Instance, error) {
	return getInstancesE(t, ig, projectId)
}

// getInstancesE returns a collection of Instance structs from the given Instance Group
func getInstances(t testing.TestingT, ig InstanceGroup, projectId string) []*Instance {
	instances, err := getInstancesE(t, ig, projectId)
	if err != nil {
		t.Fatal(err)
	}

	return instances
}

// getInstancesE returns a collection of Instance structs from the given Instance Group
func getInstancesE(t testing.TestingT, ig InstanceGroup, projectId string) ([]*Instance, error) {
	instanceIds, err := ig.GetInstanceIdsE(t)
	if err != nil {
		return nil, fmt.Errorf("Failed to get Instance Group IDs: %s", err)
	}

	var instances []*Instance

	for _, instanceId := range instanceIds {
		instance, err := FetchInstanceE(t, projectId, instanceId)
		if err != nil {
			return nil, fmt.Errorf("Failed to get Instance: %s", err)
		}

		instances = append(instances, instance)
	}

	return instances, nil
}

// GetPublicIps returns a slice of the public IPs from the given Instance Group
func (ig *ZonalInstanceGroup) GetPublicIps(t testing.TestingT, projectId string) []string {
	return getPublicIps(t, ig, projectId)
}

// GetPublicIpsE returns a slice of the public IPs from the given Instance Group
func (ig *ZonalInstanceGroup) GetPublicIpsE(t testing.TestingT, projectId string) ([]string, error) {
	return getPublicIpsE(t, ig, projectId)
}

// GetPublicIps returns a slice of the public IPs from the given Instance Group
func (ig *RegionalInstanceGroup) GetPublicIps(t testing.TestingT, projectId string) []string {
	return getPublicIps(t, ig, projectId)
}

// GetPublicIpsE returns a slice of the public IPs from the given Instance Group
func (ig *RegionalInstanceGroup) GetPublicIpsE(t testing.TestingT, projectId string) ([]string, error) {
	return getPublicIpsE(t, ig, projectId)
}

// getPublicIps a slice of the public IPs from the given Instance Group
func getPublicIps(t testing.TestingT, ig InstanceGroup, projectId string) []string {
	ips, err := getPublicIpsE(t, ig, projectId)
	if err != nil {
		t.Fatal(err)
	}

	return ips
}

// getPublicIpsE a slice of the public IPs from the given Instance Group
func getPublicIpsE(t testing.TestingT, ig InstanceGroup, projectId string) ([]string, error) {
	instances, err := getInstancesE(t, ig, projectId)
	if err != nil {
		return nil, fmt.Errorf("Failed to get Compute Instances from Instance Group: %s", err)
	}

	var ips []string

	for _, instance := range instances {
		ip := instance.GetPublicIp(t)
		ips = append(ips, ip)
	}

	return ips, nil
}

// getRandomInstance returns a randomly selected Instance from the Regional Instance Group
func (ig *ZonalInstanceGroup) GetRandomInstance(t testing.TestingT) *Instance {
	return getRandomInstance(t, ig, ig.Name, ig.Region, ig.Size, ig.projectID)
}

// getRandomInstanceE returns a randomly selected Instance from the Regional Instance Group
func (ig *ZonalInstanceGroup) GetRandomInstanceE(t testing.TestingT) (*Instance, error) {
	return getRandomInstanceE(t, ig, ig.Name, ig.Region, ig.Size, ig.projectID)
}

// getRandomInstance returns a randomly selected Instance from the Regional Instance Group
func (ig *RegionalInstanceGroup) GetRandomInstance(t testing.TestingT) *Instance {
	return getRandomInstance(t, ig, ig.Name, ig.Region, ig.Size, ig.projectID)
}

// getRandomInstanceE returns a randomly selected Instance from the Regional Instance Group
func (ig *RegionalInstanceGroup) GetRandomInstanceE(t testing.TestingT) (*Instance, error) {
	return getRandomInstanceE(t, ig, ig.Name, ig.Region, ig.Size, ig.projectID)
}

func getRandomInstance(t testing.TestingT, ig InstanceGroup, name string, region string, size int64, projectID string) *Instance {
	instance, err := getRandomInstanceE(t, ig, name, region, size, projectID)
	if err != nil {
		t.Fatal(err)
	}

	return instance
}

func getRandomInstanceE(t testing.TestingT, ig InstanceGroup, name string, region string, size int64, projectID string) (*Instance, error) {
	instanceIDs := ig.GetInstanceIds(t)
	if len(instanceIDs) == 0 {
		return nil, fmt.Errorf("Could not find any instances in Regional Instance Group or Zonal Instance Group %s in Region %s", name, region)
	}

	clusterSize := int(size)
	if len(instanceIDs) != clusterSize {
		return nil, fmt.Errorf("Expected Regional Instance Group or Zonal Instance Group %s in Region %s to have %d instances, but found %d", name, region, clusterSize, len(instanceIDs))
	}

	randIndex := random.Random(0, clusterSize-1)
	instanceID := instanceIDs[randIndex]
	instance := FetchInstance(t, projectID, instanceID)

	return instance, nil
}

// NewComputeService creates a new Compute service, which is used to make GCE API calls.
func NewComputeService(t testing.TestingT) *compute.Service {
	client, err := NewComputeServiceE(t)
	if err != nil {
		t.Fatal(err)
	}
	return client
}

// NewComputeServiceE creates a new Compute service, which is used to make GCE API calls.
func NewComputeServiceE(t testing.TestingT) (*compute.Service, error) {
	ctx := context.Background()

	// Retrieve the Google OAuth token using a retry loop as it can sometimes return an error.
	// e.g: oauth2: cannot fetch token: Post https://oauth2.googleapis.com/token: net/http: TLS handshake timeout
	// This is loosely based on https://github.com/kubernetes/kubernetes/blob/7e8de5422cb5ad76dd0c147cf4336220d282e34b/pkg/cloudprovider/providers/gce/gce.go#L831.

	description := "Attempting to request a Google OAuth2 token"
	maxRetries := 6
	timeBetweenRetries := 10 * time.Second

	var client *http.Client

	msg, retryErr := retry.DoWithRetryE(t, description, maxRetries, timeBetweenRetries, func() (string, error) {
		rawClient, err := google.DefaultClient(ctx, compute.CloudPlatformScope)
		if err != nil {
			return "Error retrieving default GCP client", err
		}
		client = rawClient
		return "Successfully retrieved default GCP client", nil
	})
	logger.Logf(t, msg)

	if retryErr != nil {
		return nil, retryErr
	}

	return compute.New(client)
}

// NewInstancesService creates a new InstancesService service, which is used to make a subset of GCE API calls.
func NewInstancesService(t testing.TestingT) *compute.InstancesService {
	client, err := NewInstancesServiceE(t)
	if err != nil {
		t.Fatal(err)
	}
	return client
}

// NewInstancesServiceE creates a new InstancesService service, which is used to make a subset of GCE API calls.
func NewInstancesServiceE(t testing.TestingT) (*compute.InstancesService, error) {
	service, err := NewComputeServiceE(t)
	if err != nil {
		return nil, fmt.Errorf("Failed to get new Instances Service\n")
	}

	return service.Instances, nil
}

// Return a random, valid name for GCP resources. Many resources in GCP requires lowercase letters only.
func RandomValidGcpName() string {
	id := strings.ToLower(random.UniqueId())
	instanceName := fmt.Sprintf("terratest-%s", id)

	return instanceName
}
