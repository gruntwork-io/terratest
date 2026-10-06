package gcp

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/random"
	"github.com/gruntwork-io/terratest/modules/core/v2/retry"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

const defaultRetryInterval = 10 * time.Second

// Instance represents a GCP Compute Instance (https://cloud.google.com/compute/docs/instances/).
type Instance struct {
	*compute.Instance
	projectID string
}

// Image represents a GCP Image (https://cloud.google.com/compute/docs/images).
type Image struct {
	*compute.Image
	projectID string
}

// ZonalInstanceGroup represents a GCP Zonal Instance Group (https://cloud.google.com/compute/docs/instance-groups/).
type ZonalInstanceGroup struct {
	*compute.InstanceGroup
	projectID string
}

// RegionalInstanceGroup represents a GCP Regional Instance Group (https://cloud.google.com/compute/docs/instance-groups/).
type RegionalInstanceGroup struct {
	*compute.InstanceGroup
	projectID string
}

// InstanceGroup is an interface for instance groups that can return their instance IDs.
type InstanceGroup interface {
	// GetInstanceIDsContextE gets the IDs of Instances in the given Instance Group.
	// The ctx parameter supports cancellation and timeouts.
	GetInstanceIDsContextE(t testing.TestingT, ctx context.Context) ([]string, error)
}

// FetchInstanceContext queries GCP to return an instance of the Compute Instance type.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchInstanceContext(t testing.TestingT, ctx context.Context, projectID string, name string) *Instance {
	instance, err := FetchInstanceContextE(t, ctx, projectID, name)
	require.NoError(t, err)

	return instance
}

// FetchInstanceContextE queries GCP to return an instance of the Compute Instance type.
// The ctx parameter supports cancellation and timeouts.
func FetchInstanceContextE(t testing.TestingT, ctx context.Context, projectID string, name string) (*Instance, error) {
	logger.Default.Logf(t, "Getting Compute Instance %s", name)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return FetchInstanceWithClient(ctx, service, projectID, name)
}

// FetchInstanceWithClient queries GCP to return an instance of the Compute Instance type
// using the supplied *compute.Service. Prefer this variant in unit tests where the service is
// backed by an httptest fake server (see compute_unit_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func FetchInstanceWithClient(ctx context.Context, service *compute.Service, projectID string, name string) (*Instance, error) {

	instanceAggregatedList, err := service.Instances.AggregatedList(projectID).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("Instances.AggregatedList(%s) got error: %w", projectID, err)
	}

	for _, instanceList := range instanceAggregatedList.Items {
		for _, instance := range instanceList.Instances {
			if name == instance.Name {
				return &Instance{Instance: instance, projectID: projectID}, nil
			}
		}
	}

	return nil, fmt.Errorf("compute Instance %s could not be found in project %s", name, projectID)
}

// FetchNetworkContext queries GCP to return the settings it holds for the given VPC network, so a
// test can assert on what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchNetworkContext(t testing.TestingT, ctx context.Context, projectID string, name string) *compute.Network {
	network, err := FetchNetworkContextE(t, ctx, projectID, name)
	require.NoError(t, err)

	return network
}

// FetchNetworkContextE queries GCP to return the settings it holds for the given VPC network.
// The ctx parameter supports cancellation and timeouts.
func FetchNetworkContextE(t testing.TestingT, ctx context.Context, projectID string, name string) (*compute.Network, error) {
	logger.Default.Logf(t, "Getting network %s", name)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return FetchNetworkWithClient(ctx, service, projectID, name)
}

// FetchNetworkWithClient queries GCP to return the settings it holds for the given VPC network
// using the supplied *compute.Service. Prefer this variant in unit tests where the service is
// backed by an httptest fake server (see compute_unit_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func FetchNetworkWithClient(ctx context.Context, service *compute.Service, projectID string, name string) (*compute.Network, error) {
	network, err := service.Networks.Get(projectID, name).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("Networks.Get(%s, %s) got error: %w", projectID, name, err)
	}

	return network, nil
}

// FetchFirewallContext queries GCP to return the settings it holds for the given firewall rule, so
// a test can assert on what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchFirewallContext(t testing.TestingT, ctx context.Context, projectID string, name string) *compute.Firewall {
	firewall, err := FetchFirewallContextE(t, ctx, projectID, name)
	require.NoError(t, err)

	return firewall
}

// FetchFirewallContextE queries GCP to return the settings it holds for the given firewall rule.
// The ctx parameter supports cancellation and timeouts.
func FetchFirewallContextE(t testing.TestingT, ctx context.Context, projectID string, name string) (*compute.Firewall, error) {
	logger.Default.Logf(t, "Getting firewall rule %s", name)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return FetchFirewallWithClient(ctx, service, projectID, name)
}

// FetchFirewallWithClient queries GCP to return the settings it holds for the given firewall rule
// using the supplied *compute.Service. Prefer this variant in unit tests where the service is
// backed by an httptest fake server (see compute_unit_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func FetchFirewallWithClient(ctx context.Context, service *compute.Service, projectID string, name string) (*compute.Firewall, error) {
	firewall, err := service.Firewalls.Get(projectID, name).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("Firewalls.Get(%s, %s) got error: %w", projectID, name, err)
	}

	return firewall, nil
}

// FetchImageContext queries GCP to return a new instance of the Compute Image type.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchImageContext(t testing.TestingT, ctx context.Context, projectID string, name string) *Image {
	image, err := FetchImageContextE(t, ctx, projectID, name)
	require.NoError(t, err)

	return image
}

// FetchImageContextE queries GCP to return a new instance of the Compute Image type.
// The ctx parameter supports cancellation and timeouts.
func FetchImageContextE(t testing.TestingT, ctx context.Context, projectID string, name string) (*Image, error) {
	logger.Default.Logf(t, "Getting Image %s", name)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return FetchImageWithClient(ctx, service, projectID, name)
}

// FetchImageWithClient queries GCP to return a new instance of the Compute Image type
// using the supplied *compute.Service. Prefer this variant in unit tests where the service is
// backed by an httptest fake server (see compute_unit_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func FetchImageWithClient(ctx context.Context, service *compute.Service, projectID string, name string) (*Image, error) {
	image, err := service.Images.Get(projectID, name).Context(ctx).Do()
	if err != nil {
		return nil, err
	}

	return &Image{Image: image, projectID: projectID}, nil
}

// FetchRegionalInstanceGroupContext queries GCP to return a new instance of the Regional Instance Group type.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchRegionalInstanceGroupContext(t testing.TestingT, ctx context.Context, projectID string, region string, name string) *RegionalInstanceGroup {
	instanceGroup, err := FetchRegionalInstanceGroupContextE(t, ctx, projectID, region, name)
	require.NoError(t, err)

	return instanceGroup
}

// FetchRegionalInstanceGroupContextE queries GCP to return a new instance of the Regional Instance Group type.
// The ctx parameter supports cancellation and timeouts.
func FetchRegionalInstanceGroupContextE(t testing.TestingT, ctx context.Context, projectID string, region string, name string) (*RegionalInstanceGroup, error) {
	logger.Default.Logf(t, "Getting Regional Instance Group %s", name)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return FetchRegionalInstanceGroupWithClient(ctx, service, projectID, region, name)
}

// FetchRegionalInstanceGroupWithClient queries GCP to return a new instance of the Regional Instance
// Group type using the supplied *compute.Service. Prefer this variant in unit tests where the service
// is backed by an httptest fake server (see compute_unit_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func FetchRegionalInstanceGroupWithClient(ctx context.Context, service *compute.Service, projectID string, region string, name string) (*RegionalInstanceGroup, error) {
	instanceGroup, err := service.RegionInstanceGroups.Get(projectID, region, name).Context(ctx).Do()
	if err != nil {
		return nil, err
	}

	return &RegionalInstanceGroup{InstanceGroup: instanceGroup, projectID: projectID}, nil
}

// FetchZonalInstanceGroupContext queries GCP to return a new instance of the Zonal Instance Group type.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchZonalInstanceGroupContext(t testing.TestingT, ctx context.Context, projectID string, zone string, name string) *ZonalInstanceGroup {
	instanceGroup, err := FetchZonalInstanceGroupContextE(t, ctx, projectID, zone, name)
	require.NoError(t, err)

	return instanceGroup
}

// FetchZonalInstanceGroupContextE queries GCP to return a new instance of the Zonal Instance Group type.
// The ctx parameter supports cancellation and timeouts.
func FetchZonalInstanceGroupContextE(t testing.TestingT, ctx context.Context, projectID string, zone string, name string) (*ZonalInstanceGroup, error) {
	logger.Default.Logf(t, "Getting Zonal Instance Group %s", name)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return FetchZonalInstanceGroupWithClient(ctx, service, projectID, zone, name)
}

// FetchZonalInstanceGroupWithClient queries GCP to return a new instance of the Zonal Instance Group
// type using the supplied *compute.Service. Prefer this variant in unit tests where the service is
// backed by an httptest fake server (see compute_unit_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func FetchZonalInstanceGroupWithClient(ctx context.Context, service *compute.Service, projectID string, zone string, name string) (*ZonalInstanceGroup, error) {
	instanceGroup, err := service.InstanceGroups.Get(projectID, zone, name).Context(ctx).Do()
	if err != nil {
		return nil, err
	}

	return &ZonalInstanceGroup{InstanceGroup: instanceGroup, projectID: projectID}, nil
}

// GetPublicIPContext gets the public IP address of the given Compute Instance.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func (i *Instance) GetPublicIPContext(t testing.TestingT, ctx context.Context) string {
	ip, err := i.GetPublicIPContextE(t, ctx)
	require.NoError(t, err)

	return ip
}

// GetPublicIPContextE gets the public IP address of the given Compute Instance.
// The ctx parameter supports cancellation and timeouts.
// Returns an error (rather than panicking) when the instance has no network interfaces
// or the first interface has no accessConfigs — both indicate the instance has no
// external internet access (https://cloud.google.com/compute/docs/reference/rest/v1/instances).
func (i *Instance) GetPublicIPContextE(t testing.TestingT, ctx context.Context) (string, error) {
	if len(i.NetworkInterfaces) == 0 || len(i.NetworkInterfaces[0].AccessConfigs) == 0 {
		return "", fmt.Errorf("attempted to get public IP of Compute Instance %s, but that Compute Instance does not have a public IP address", i.Name)
	}

	return i.NetworkInterfaces[0].AccessConfigs[0].NatIP, nil
}

// GetLabelsContext returns all the tags for the given Compute Instance.
// The ctx parameter supports cancellation and timeouts.
func (i *Instance) GetLabelsContext(t testing.TestingT, ctx context.Context) map[string]string {
	labels, err := i.GetLabelsContextE(t, ctx)
	require.NoError(t, err)

	return labels
}

// GetLabelsContextE returns all the tags for the given Compute Instance.
// The ctx parameter supports cancellation and timeouts.
func (i *Instance) GetLabelsContextE(t testing.TestingT, ctx context.Context) (map[string]string, error) {
	return i.Labels, nil
}

// GetZoneContext returns the Zone in which the Compute Instance is located.
// The ctx parameter supports cancellation and timeouts.
func (i *Instance) GetZoneContext(t testing.TestingT, ctx context.Context) string {
	zone, err := i.GetZoneContextE(t, ctx)
	require.NoError(t, err)

	return zone
}

// GetZoneContextE returns the Zone in which the Compute Instance is located.
// The ctx parameter supports cancellation and timeouts.
func (i *Instance) GetZoneContextE(t testing.TestingT, ctx context.Context) (string, error) {
	return ZoneURLToZone(i.Zone), nil
}

// SetLabelsContext adds the tags to the given Compute Instance.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func (i *Instance) SetLabelsContext(t testing.TestingT, ctx context.Context, labels map[string]string) {
	err := i.SetLabelsContextE(t, ctx, labels)
	require.NoError(t, err)
}

// SetLabelsContextE adds the tags to the given Compute Instance.
// The ctx parameter supports cancellation and timeouts.
func (i *Instance) SetLabelsContextE(t testing.TestingT, ctx context.Context, labels map[string]string) error {
	logger.Default.Logf(t, "Adding labels to instance %s in zone %s", i.Name, i.Zone)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return err
	}

	return i.SetLabelsWithClient(ctx, service, labels)
}

// SetLabelsWithClient merges the given labels into the instance's existing labels using the
// supplied *compute.Service. Keys present in labels overwrite existing values; other labels
// are preserved. Prefer this variant in unit tests where the service is backed by an httptest
// fake server (see compute_unit_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func (i *Instance) SetLabelsWithClient(ctx context.Context, service *compute.Service, labels map[string]string) error {
	merged := make(map[string]string, len(i.Labels)+len(labels))
	for k, v := range i.Labels {
		merged[k] = v
	}

	for k, v := range labels {
		merged[k] = v
	}

	req := compute.InstancesSetLabelsRequest{Labels: merged, LabelFingerprint: i.LabelFingerprint}

	if _, err := service.Instances.SetLabels(i.projectID, ZoneURLToZone(i.Zone), i.Name, &req).Context(ctx).Do(); err != nil {
		return fmt.Errorf("Instances.SetLabels(%s) got error: %w", i.Name, err)
	}

	return nil
}

// GetMetadataContext gets the given Compute Instance's metadata.
// The ctx parameter supports cancellation and timeouts.
func (i *Instance) GetMetadataContext(t testing.TestingT, ctx context.Context) []*compute.MetadataItems {
	metadata, err := i.GetMetadataContextE(t, ctx)
	require.NoError(t, err)

	return metadata
}

// GetMetadataContextE gets the given Compute Instance's metadata.
// The ctx parameter supports cancellation and timeouts.
func (i *Instance) GetMetadataContextE(t testing.TestingT, ctx context.Context) ([]*compute.MetadataItems, error) {
	return i.Metadata.Items, nil
}

// SetMetadataContext sets the given Compute Instance's metadata.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func (i *Instance) SetMetadataContext(t testing.TestingT, ctx context.Context, metadata map[string]string) {
	err := i.SetMetadataContextE(t, ctx, metadata)
	require.NoError(t, err)
}

// SetMetadataContextE adds the given metadata map to the existing metadata of the given Compute Instance.
// The ctx parameter supports cancellation and timeouts.
func (i *Instance) SetMetadataContextE(t testing.TestingT, ctx context.Context, metadata map[string]string) error {
	logger.Default.Logf(t, "Adding metadata to instance %s in zone %s", i.Name, i.Zone)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return err
	}

	return i.SetMetadataWithClient(ctx, service, metadata)
}

// SetMetadataWithClient adds the given metadata map to the existing metadata of the given Compute
// Instance using the supplied *compute.Service. Prefer this variant in unit tests where the service
// is backed by an httptest fake server (see compute_unit_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func (i *Instance) SetMetadataWithClient(ctx context.Context, service *compute.Service, metadata map[string]string) error {
	metadataItems := NewMetadata(i.Metadata, metadata)

	req := service.Instances.SetMetadata(i.projectID, ZoneURLToZone(i.Zone), i.Name, metadataItems)

	if _, err := req.Context(ctx).Do(); err != nil {
		return fmt.Errorf("Instances.SetMetadata(%s) got error: %w", i.Name, err)
	}

	return nil
}

// NewMetadata merges new key-value pairs into existing metadata, preserving unmodified items.
func NewMetadata(oldMetadata *compute.Metadata, kvs map[string]string) *compute.Metadata {
	itemsMap := make(map[string]*string)

	if oldMetadata != nil {
		for _, item := range oldMetadata.Items {
			itemsMap[item.Key] = item.Value
		}
	}

	for key, val := range kvs {
		v := val
		itemsMap[key] = &v
	}

	items := make([]*compute.MetadataItems, 0, len(itemsMap))

	for key, val := range itemsMap {
		items = append(items, &compute.MetadataItems{Key: key, Value: val})
	}

	fingerprint := ""

	if oldMetadata != nil {
		fingerprint = oldMetadata.Fingerprint
	}

	return &compute.Metadata{
		Fingerprint: fingerprint,
		Items:       items,
	}
}

// AddSSHKeyContext adds the given public SSH key to the Compute Instance. Users can SSH in with the given username.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func (i *Instance) AddSSHKeyContext(t testing.TestingT, ctx context.Context, username string, publicKey string) {
	err := i.AddSSHKeyContextE(t, ctx, username, publicKey)
	require.NoError(t, err)
}

// AddSSHKeyContextE adds the given public SSH key to the Compute Instance. Users can SSH in with the given username.
// The ctx parameter supports cancellation and timeouts.
func (i *Instance) AddSSHKeyContextE(t testing.TestingT, ctx context.Context, username string, publicKey string) error {
	logger.Default.Logf(t, "Adding SSH Key to Compute Instance %s for username %s\n", i.Name, username)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return err
	}

	return i.AddSSHKeyWithClient(ctx, service, username, publicKey)
}

// AddSSHKeyWithClient adds the given public SSH key to the Compute Instance using the supplied
// *compute.Service. Users can SSH in with the given username. Prefer this variant in unit tests
// where the service is backed by an httptest fake server (see compute_unit_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func (i *Instance) AddSSHKeyWithClient(ctx context.Context, service *compute.Service, username string, publicKey string) error {

	publicKeyFormatted := strings.TrimSpace(publicKey)
	sshKeyFormatted := fmt.Sprintf("%s:%s %s", username, publicKeyFormatted, username)

	metadata := map[string]string{
		"ssh-keys": sshKeyFormatted,
	}

	if err := i.SetMetadataWithClient(ctx, service, metadata); err != nil {
		return fmt.Errorf("failed to add SSH key to Compute Instance: %w", err)
	}

	return nil
}

// DeleteImageContext deletes the given Compute Image.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func (i *Image) DeleteImageContext(t testing.TestingT, ctx context.Context) {
	err := i.DeleteImageContextE(t, ctx)
	require.NoError(t, err)
}

// DeleteImageContextE deletes the given Compute Image.
// The ctx parameter supports cancellation and timeouts.
func (i *Image) DeleteImageContextE(t testing.TestingT, ctx context.Context) error {
	logger.Default.Logf(t, "Destroying Image %s", i.Name)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return err
	}

	return i.DeleteImageWithClient(ctx, service)
}

// DeleteImageWithClient deletes the given Compute Image using the supplied *compute.Service.
// Prefer this variant in unit tests where the service is backed by an httptest fake server
// (see compute_unit_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func (i *Image) DeleteImageWithClient(ctx context.Context, service *compute.Service) error {
	if _, err := service.Images.Delete(i.projectID, i.Name).Context(ctx).Do(); err != nil {
		return fmt.Errorf("Images.Delete(%s) got error: %w", i.Name, err)
	}

	return nil
}

// GetInstanceIDsContext gets the IDs of Instances in the given Zonal Instance Group.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func (ig *ZonalInstanceGroup) GetInstanceIDsContext(t testing.TestingT, ctx context.Context) []string {
	ids, err := ig.GetInstanceIDsContextE(t, ctx)
	require.NoError(t, err)

	return ids
}

// GetInstanceIDsContextE gets the IDs of Instances in the given Zonal Instance Group.
// The ctx parameter supports cancellation and timeouts.
func (ig *ZonalInstanceGroup) GetInstanceIDsContextE(t testing.TestingT, ctx context.Context) ([]string, error) {
	logger.Default.Logf(t, "Get instances for Zonal Instance Group %s", ig.Name)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return ig.GetInstanceIDsWithClient(ctx, service)
}

// GetInstanceIDsWithClient gets the IDs of Instances in the given Zonal Instance Group using the
// supplied *compute.Service. Prefer this variant in unit tests where the service is backed by an
// httptest fake server (see compute_unit_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func (ig *ZonalInstanceGroup) GetInstanceIDsWithClient(ctx context.Context, service *compute.Service) ([]string, error) {
	requestBody := &compute.InstanceGroupsListInstancesRequest{
		InstanceState: "ALL",
	}

	instanceIDs := []string{}
	zone := ZoneURLToZone(ig.Zone)
	req := service.InstanceGroups.ListInstances(ig.projectID, zone, ig.Name, requestBody)

	err := req.Pages(ctx, func(page *compute.InstanceGroupsListInstances) error {
		for _, instance := range page.Items {

			instanceID := path.Base(instance.Instance)
			instanceIDs = append(instanceIDs, instanceID)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("InstanceGroups.ListInstances(%s) got error: %w", ig.Name, err)
	}

	return instanceIDs, nil
}

// GetInstanceIDsContext gets the IDs of Instances in the given Regional Instance Group.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func (ig *RegionalInstanceGroup) GetInstanceIDsContext(t testing.TestingT, ctx context.Context) []string {
	ids, err := ig.GetInstanceIDsContextE(t, ctx)
	require.NoError(t, err)

	return ids
}

// GetInstanceIDsContextE gets the IDs of Instances in the given Regional Instance Group.
// The ctx parameter supports cancellation and timeouts.
func (ig *RegionalInstanceGroup) GetInstanceIDsContextE(t testing.TestingT, ctx context.Context) ([]string, error) {
	logger.Default.Logf(t, "Get instances for Regional Instance Group %s", ig.Name)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return ig.GetInstanceIDsWithClient(ctx, service)
}

// GetInstanceIDsWithClient gets the IDs of Instances in the given Regional Instance Group using the
// supplied *compute.Service. Prefer this variant in unit tests where the service is backed by an
// httptest fake server (see compute_unit_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func (ig *RegionalInstanceGroup) GetInstanceIDsWithClient(ctx context.Context, service *compute.Service) ([]string, error) {
	requestBody := &compute.RegionInstanceGroupsListInstancesRequest{
		InstanceState: "ALL",
	}

	instanceIDs := []string{}
	region := RegionURLToRegion(ig.Region)
	req := service.RegionInstanceGroups.ListInstances(ig.projectID, region, ig.Name, requestBody)

	err := req.Pages(ctx, func(page *compute.RegionInstanceGroupsListInstances) error {
		for _, instance := range page.Items {

			instanceID := path.Base(instance.Instance)
			instanceIDs = append(instanceIDs, instanceID)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("InstanceGroups.ListInstances(%s) got error: %w", ig.Name, err)
	}

	return instanceIDs, nil
}

// GetInstancesContext returns a collection of Instance structs from the given Zonal Instance Group.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func (ig *ZonalInstanceGroup) GetInstancesContext(t testing.TestingT, ctx context.Context, projectID string) []*Instance {
	instances, err := ig.GetInstancesContextE(t, ctx, projectID)
	require.NoError(t, err)

	return instances
}

// GetInstancesContextE returns a collection of Instance structs from the given Zonal Instance Group.
// The ctx parameter supports cancellation and timeouts.
func (ig *ZonalInstanceGroup) GetInstancesContextE(t testing.TestingT, ctx context.Context, projectID string) ([]*Instance, error) {
	return getInstancesContextE(t, ctx, ig, projectID)
}

// GetInstancesContext returns a collection of Instance structs from the given Regional Instance Group.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func (ig *RegionalInstanceGroup) GetInstancesContext(t testing.TestingT, ctx context.Context, projectID string) []*Instance {
	instances, err := ig.GetInstancesContextE(t, ctx, projectID)
	require.NoError(t, err)

	return instances
}

// GetInstancesContextE returns a collection of Instance structs from the given Regional Instance Group.
// The ctx parameter supports cancellation and timeouts.
func (ig *RegionalInstanceGroup) GetInstancesContextE(t testing.TestingT, ctx context.Context, projectID string) ([]*Instance, error) {
	return getInstancesContextE(t, ctx, ig, projectID)
}

// getInstancesContextE returns a collection of Instance structs from the given Instance Group.
func getInstancesContextE(t testing.TestingT, ctx context.Context, ig InstanceGroup, projectID string) ([]*Instance, error) {
	instanceIDs, err := ig.GetInstanceIDsContextE(t, ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get Instance Group IDs: %w", err)
	}

	var instances []*Instance

	for _, instanceID := range instanceIDs {
		instance, err := FetchInstanceContextE(t, ctx, projectID, instanceID)
		if err != nil {
			return nil, fmt.Errorf("failed to get Instance: %w", err)
		}

		instances = append(instances, instance)
	}

	return instances, nil
}

// GetPublicIPsContext returns a slice of the public IPs from the given Zonal Instance Group.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func (ig *ZonalInstanceGroup) GetPublicIPsContext(t testing.TestingT, ctx context.Context, projectID string) []string {
	ips, err := ig.GetPublicIPsContextE(t, ctx, projectID)
	require.NoError(t, err)

	return ips
}

// GetPublicIPsContextE returns a slice of the public IPs from the given Zonal Instance Group.
// The ctx parameter supports cancellation and timeouts.
func (ig *ZonalInstanceGroup) GetPublicIPsContextE(t testing.TestingT, ctx context.Context, projectID string) ([]string, error) {
	return getPublicIPsContextE(t, ctx, ig, projectID)
}

// GetPublicIPsContext returns a slice of the public IPs from the given Regional Instance Group.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func (ig *RegionalInstanceGroup) GetPublicIPsContext(t testing.TestingT, ctx context.Context, projectID string) []string {
	ips, err := ig.GetPublicIPsContextE(t, ctx, projectID)
	require.NoError(t, err)

	return ips
}

// GetPublicIPsContextE returns a slice of the public IPs from the given Regional Instance Group.
// The ctx parameter supports cancellation and timeouts.
func (ig *RegionalInstanceGroup) GetPublicIPsContextE(t testing.TestingT, ctx context.Context, projectID string) ([]string, error) {
	return getPublicIPsContextE(t, ctx, ig, projectID)
}

// getPublicIPsContextE returns a slice of the public IPs from the given Instance Group.
func getPublicIPsContextE(t testing.TestingT, ctx context.Context, ig InstanceGroup, projectID string) ([]string, error) {
	instances, err := getInstancesContextE(t, ctx, ig, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get Compute Instances from Instance Group: %w", err)
	}

	var ips []string

	for _, instance := range instances {
		ip, err := instance.GetPublicIPContextE(t, ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get public IP for instance: %w", err)
		}

		ips = append(ips, ip)
	}

	return ips, nil
}

// GetRandomInstanceContext returns a randomly selected Instance from the Zonal Instance Group.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func (ig *ZonalInstanceGroup) GetRandomInstanceContext(t testing.TestingT, ctx context.Context) *Instance {
	instance, err := ig.GetRandomInstanceContextE(t, ctx)
	require.NoError(t, err)

	return instance
}

// GetRandomInstanceContextE returns a randomly selected Instance from the Zonal Instance Group.
// The ctx parameter supports cancellation and timeouts.
func (ig *ZonalInstanceGroup) GetRandomInstanceContextE(t testing.TestingT, ctx context.Context) (*Instance, error) {
	return getRandomInstanceContextE(t, ctx, ig, ig.Name, ig.Region, ig.Size, ig.projectID)
}

// GetRandomInstanceContext returns a randomly selected Instance from the Regional Instance Group.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func (ig *RegionalInstanceGroup) GetRandomInstanceContext(t testing.TestingT, ctx context.Context) *Instance {
	instance, err := ig.GetRandomInstanceContextE(t, ctx)
	require.NoError(t, err)

	return instance
}

// GetRandomInstanceContextE returns a randomly selected Instance from the Regional Instance Group.
// The ctx parameter supports cancellation and timeouts.
func (ig *RegionalInstanceGroup) GetRandomInstanceContextE(t testing.TestingT, ctx context.Context) (*Instance, error) {
	return getRandomInstanceContextE(t, ctx, ig, ig.Name, ig.Region, ig.Size, ig.projectID)
}

func getRandomInstanceContextE(t testing.TestingT, ctx context.Context, ig InstanceGroup, name string, region string, size int64, projectID string) (*Instance, error) {
	instanceIDs, err := ig.GetInstanceIDsContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	if len(instanceIDs) == 0 {
		return nil, fmt.Errorf("could not find any instances in Instance Group %s in Region %s", name, region)
	}

	clusterSize := int(size)
	if len(instanceIDs) != clusterSize {
		return nil, fmt.Errorf("expected Instance Group %s in Region %s to have %d instances, but found %d", name, region, clusterSize, len(instanceIDs))
	}

	randIndex := random.Random(0, clusterSize-1)
	instanceID := instanceIDs[randIndex]

	instance, err := FetchInstanceContextE(t, ctx, projectID, instanceID)
	if err != nil {
		return nil, err
	}

	return instance, nil
}

// NewComputeServiceContext creates a new Compute service, which is used to make GCE API calls.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func NewComputeServiceContext(t testing.TestingT, ctx context.Context) *compute.Service {
	client, err := NewComputeServiceContextE(t, ctx)
	require.NoError(t, err)

	return client
}

// FetchNodeGroup returns the settings Google Cloud holds for the given sole-tenant node group, so a test can assert on what was
// actually created rather than only that it exists.
// A node group is the set of physical servers a sole-tenant VM is placed on, so its size and its maintenance policy decide where those VMs can run.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchNodeGroup(t testing.TestingT, ctx context.Context, projectID string, zone string, name string) *compute.NodeGroup {
	group, err := FetchNodeGroupE(t, ctx, projectID, zone, name)
	require.NoError(t, err)

	return group
}

// FetchNodeGroupE returns the settings Google Cloud holds for the given sole-tenant node group.
// The ctx parameter supports cancellation and timeouts.
func FetchNodeGroupE(t testing.TestingT, ctx context.Context, projectID string, zone string, name string) (*compute.NodeGroup, error) {
	logger.Default.Logf(t, "Getting settings for sole-tenant node group %s in %s in project %s", name, zone, projectID)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return FetchNodeGroupWithClient(ctx, service, projectID, zone, name)
}

// FetchNodeGroupWithClient returns the settings Google Cloud holds for the given sole-tenant node group using the supplied
// *compute.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see compute_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func FetchNodeGroupWithClient(ctx context.Context, service *compute.Service, projectID string, zone string, name string) (*compute.NodeGroup, error) {
	group, err := service.NodeGroups.Get(projectID, zone, name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the sole-tenant node group %s in %s in project %s does not exist", name, zone, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for sole-tenant node group %s in %s in project %s: %w", name, zone, projectID, err)
	}

	return group, nil
}

// FetchReservation returns the settings Google Cloud holds for the given reservation, so a test can assert on what was
// actually created rather than only that it exists.
// A reservation holds capacity for VMs that do not exist yet, so what it holds and whether it is shared is the whole point of it.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchReservation(t testing.TestingT, ctx context.Context, projectID string, zone string, name string) *compute.Reservation {
	reservation, err := FetchReservationE(t, ctx, projectID, zone, name)
	require.NoError(t, err)

	return reservation
}

// FetchReservationE returns the settings Google Cloud holds for the given reservation.
// The ctx parameter supports cancellation and timeouts.
func FetchReservationE(t testing.TestingT, ctx context.Context, projectID string, zone string, name string) (*compute.Reservation, error) {
	logger.Default.Logf(t, "Getting settings for reservation %s in %s in project %s", name, zone, projectID)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return FetchReservationWithClient(ctx, service, projectID, zone, name)
}

// FetchReservationWithClient returns the settings Google Cloud holds for the given reservation using the supplied
// *compute.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see compute_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func FetchReservationWithClient(ctx context.Context, service *compute.Service, projectID string, zone string, name string) (*compute.Reservation, error) {
	reservation, err := service.Reservations.Get(projectID, zone, name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the reservation %s in %s in project %s does not exist", name, zone, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for reservation %s in %s in project %s: %w", name, zone, projectID, err)
	}

	return reservation, nil
}

// FetchRegionCommitment returns the settings Google Cloud holds for the given commitment, so a test can assert on what was
// actually created rather than only that it exists.
// A commitment buys a discount for a fixed term, and neither the term nor the resources it covers can be changed once it is made.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchRegionCommitment(t testing.TestingT, ctx context.Context, projectID string, region string, name string) *compute.Commitment {
	commitment, err := FetchRegionCommitmentE(t, ctx, projectID, region, name)
	require.NoError(t, err)

	return commitment
}

// FetchRegionCommitmentE returns the settings Google Cloud holds for the given commitment.
// The ctx parameter supports cancellation and timeouts.
func FetchRegionCommitmentE(t testing.TestingT, ctx context.Context, projectID string, region string, name string) (*compute.Commitment, error) {
	logger.Default.Logf(t, "Getting settings for commitment %s in %s in project %s", name, region, projectID)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return FetchRegionCommitmentWithClient(ctx, service, projectID, region, name)
}

// FetchRegionCommitmentWithClient returns the settings Google Cloud holds for the given commitment using the supplied
// *compute.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see compute_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func FetchRegionCommitmentWithClient(ctx context.Context, service *compute.Service, projectID string, region string, name string) (*compute.Commitment, error) {
	commitment, err := service.RegionCommitments.Get(projectID, region, name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the commitment %s in %s in project %s does not exist", name, region, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for commitment %s in %s in project %s: %w", name, region, projectID, err)
	}

	return commitment, nil
}

// FetchStoragePool returns the settings Google Cloud holds for the given storage pool, so a test can assert on what was
// actually created rather than only that it exists.
// A pool buys capacity and throughput that the disks inside it draw on, so those two numbers are what a disk in the pool actually gets.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchStoragePool(t testing.TestingT, ctx context.Context, projectID string, zone string, name string) *compute.StoragePool {
	pool, err := FetchStoragePoolE(t, ctx, projectID, zone, name)
	require.NoError(t, err)

	return pool
}

// FetchStoragePoolE returns the settings Google Cloud holds for the given storage pool.
// The ctx parameter supports cancellation and timeouts.
func FetchStoragePoolE(t testing.TestingT, ctx context.Context, projectID string, zone string, name string) (*compute.StoragePool, error) {
	logger.Default.Logf(t, "Getting settings for storage pool %s in %s in project %s", name, zone, projectID)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return FetchStoragePoolWithClient(ctx, service, projectID, zone, name)
}

// FetchStoragePoolWithClient returns the settings Google Cloud holds for the given storage pool using the supplied
// *compute.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see compute_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func FetchStoragePoolWithClient(ctx context.Context, service *compute.Service, projectID string, zone string, name string) (*compute.StoragePool, error) {
	pool, err := service.StoragePools.Get(projectID, zone, name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the storage pool %s in %s in project %s does not exist", name, zone, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for storage pool %s in %s in project %s: %w", name, zone, projectID, err)
	}

	return pool, nil
}

// FetchPublicAdvertisedPrefix returns the settings Google Cloud holds for the given public advertised prefix, so a test can assert on what was
// actually created rather than only that it exists.
// This is the range an organisation brings to Google and proves it owns, so the range and its verification status are the point.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchPublicAdvertisedPrefix(t testing.TestingT, ctx context.Context, projectID string, name string) *compute.PublicAdvertisedPrefix {
	prefix, err := FetchPublicAdvertisedPrefixE(t, ctx, projectID, name)
	require.NoError(t, err)

	return prefix
}

// FetchPublicAdvertisedPrefixE returns the settings Google Cloud holds for the given public advertised prefix.
// The ctx parameter supports cancellation and timeouts.
func FetchPublicAdvertisedPrefixE(t testing.TestingT, ctx context.Context, projectID string, name string) (*compute.PublicAdvertisedPrefix, error) {
	logger.Default.Logf(t, "Getting settings for public advertised prefix %s in project %s", name, projectID)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return FetchPublicAdvertisedPrefixWithClient(ctx, service, projectID, name)
}

// FetchPublicAdvertisedPrefixWithClient returns the settings Google Cloud holds for the given public advertised prefix using the supplied
// *compute.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see compute_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func FetchPublicAdvertisedPrefixWithClient(ctx context.Context, service *compute.Service, projectID string, name string) (*compute.PublicAdvertisedPrefix, error) {
	prefix, err := service.PublicAdvertisedPrefixes.Get(projectID, name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the public advertised prefix %s in project %s does not exist", name, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for public advertised prefix %s in project %s: %w", name, projectID, err)
	}

	return prefix, nil
}

// FetchPublicDelegatedPrefix returns the settings Google Cloud holds for the given public delegated prefix, so a test can assert on what was
// actually created rather than only that it exists.
// A delegated prefix carves a region's share out of an advertised prefix, so which parent it came from and whether it is live both matter.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchPublicDelegatedPrefix(t testing.TestingT, ctx context.Context, projectID string, region string, name string) *compute.PublicDelegatedPrefix {
	prefix, err := FetchPublicDelegatedPrefixE(t, ctx, projectID, region, name)
	require.NoError(t, err)

	return prefix
}

// FetchPublicDelegatedPrefixE returns the settings Google Cloud holds for the given public delegated prefix.
// The ctx parameter supports cancellation and timeouts.
func FetchPublicDelegatedPrefixE(t testing.TestingT, ctx context.Context, projectID string, region string, name string) (*compute.PublicDelegatedPrefix, error) {
	logger.Default.Logf(t, "Getting settings for public delegated prefix %s in %s in project %s", name, region, projectID)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return FetchPublicDelegatedPrefixWithClient(ctx, service, projectID, region, name)
}

// FetchPublicDelegatedPrefixWithClient returns the settings Google Cloud holds for the given public delegated prefix using the supplied
// *compute.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see compute_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func FetchPublicDelegatedPrefixWithClient(ctx context.Context, service *compute.Service, projectID string, region string, name string) (*compute.PublicDelegatedPrefix, error) {
	prefix, err := service.PublicDelegatedPrefixes.Get(projectID, region, name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the public delegated prefix %s in %s in project %s does not exist", name, region, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for public delegated prefix %s in %s in project %s: %w", name, region, projectID, err)
	}

	return prefix, nil
}

// FetchInterconnect returns the settings Google Cloud holds for the given interconnect, so a test can assert on what was
// actually created rather than only that it exists.
// An interconnect is a physical cross-connect into Google's network, so its link type and location describe hardware rather than configuration.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchInterconnect(t testing.TestingT, ctx context.Context, projectID string, name string) *compute.Interconnect {
	interconnect, err := FetchInterconnectE(t, ctx, projectID, name)
	require.NoError(t, err)

	return interconnect
}

// FetchInterconnectE returns the settings Google Cloud holds for the given interconnect.
// The ctx parameter supports cancellation and timeouts.
func FetchInterconnectE(t testing.TestingT, ctx context.Context, projectID string, name string) (*compute.Interconnect, error) {
	logger.Default.Logf(t, "Getting settings for interconnect %s in project %s", name, projectID)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return FetchInterconnectWithClient(ctx, service, projectID, name)
}

// FetchInterconnectWithClient returns the settings Google Cloud holds for the given interconnect using the supplied
// *compute.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see compute_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func FetchInterconnectWithClient(ctx context.Context, service *compute.Service, projectID string, name string) (*compute.Interconnect, error) {
	interconnect, err := service.Interconnects.Get(projectID, name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the interconnect %s in project %s does not exist", name, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for interconnect %s in project %s: %w", name, projectID, err)
	}

	return interconnect, nil
}

// FetchInterconnectGroup returns the settings Google Cloud holds for the given interconnect group, so a test can assert on what was
// actually created rather than only that it exists.
// A group is how several interconnects are treated as one unit for redundancy, so which members it holds is the point.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchInterconnectGroup(t testing.TestingT, ctx context.Context, projectID string, name string) *compute.InterconnectGroup {
	group, err := FetchInterconnectGroupE(t, ctx, projectID, name)
	require.NoError(t, err)

	return group
}

// FetchInterconnectGroupE returns the settings Google Cloud holds for the given interconnect group.
// The ctx parameter supports cancellation and timeouts.
func FetchInterconnectGroupE(t testing.TestingT, ctx context.Context, projectID string, name string) (*compute.InterconnectGroup, error) {
	logger.Default.Logf(t, "Getting settings for interconnect group %s in project %s", name, projectID)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return FetchInterconnectGroupWithClient(ctx, service, projectID, name)
}

// FetchInterconnectGroupWithClient returns the settings Google Cloud holds for the given interconnect group using the supplied
// *compute.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see compute_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func FetchInterconnectGroupWithClient(ctx context.Context, service *compute.Service, projectID string, name string) (*compute.InterconnectGroup, error) {
	group, err := service.InterconnectGroups.Get(projectID, name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the interconnect group %s in project %s does not exist", name, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for interconnect group %s in project %s: %w", name, projectID, err)
	}

	return group, nil
}

// FetchInterconnectAttachment returns the settings Google Cloud holds for the given interconnect attachment, so a test can assert on what was
// actually created rather than only that it exists.
// An attachment joins a VLAN on an interconnect to a cloud router, so its bandwidth and the router it reaches decide what traffic can cross.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchInterconnectAttachment(t testing.TestingT, ctx context.Context, projectID string, region string, name string) *compute.InterconnectAttachment {
	attachment, err := FetchInterconnectAttachmentE(t, ctx, projectID, region, name)
	require.NoError(t, err)

	return attachment
}

// FetchInterconnectAttachmentE returns the settings Google Cloud holds for the given interconnect attachment.
// The ctx parameter supports cancellation and timeouts.
func FetchInterconnectAttachmentE(t testing.TestingT, ctx context.Context, projectID string, region string, name string) (*compute.InterconnectAttachment, error) {
	logger.Default.Logf(t, "Getting settings for interconnect attachment %s in %s in project %s", name, region, projectID)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return FetchInterconnectAttachmentWithClient(ctx, service, projectID, region, name)
}

// FetchInterconnectAttachmentWithClient returns the settings Google Cloud holds for the given interconnect attachment using the supplied
// *compute.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see compute_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func FetchInterconnectAttachmentWithClient(ctx context.Context, service *compute.Service, projectID string, region string, name string) (*compute.InterconnectAttachment, error) {
	attachment, err := service.InterconnectAttachments.Get(projectID, region, name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the interconnect attachment %s in %s in project %s does not exist", name, region, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for interconnect attachment %s in %s in project %s: %w", name, region, projectID, err)
	}

	return attachment, nil
}

// FetchCrossSiteNetwork returns the settings Google Cloud holds for the given cross-site network, so a test can assert on what was
// actually created rather than only that it exists.
// A cross-site network is the container that wire groups between two sites belong to.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchCrossSiteNetwork(t testing.TestingT, ctx context.Context, projectID string, name string) *compute.CrossSiteNetwork {
	network, err := FetchCrossSiteNetworkE(t, ctx, projectID, name)
	require.NoError(t, err)

	return network
}

// FetchCrossSiteNetworkE returns the settings Google Cloud holds for the given cross-site network.
// The ctx parameter supports cancellation and timeouts.
func FetchCrossSiteNetworkE(t testing.TestingT, ctx context.Context, projectID string, name string) (*compute.CrossSiteNetwork, error) {
	logger.Default.Logf(t, "Getting settings for cross-site network %s in project %s", name, projectID)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return FetchCrossSiteNetworkWithClient(ctx, service, projectID, name)
}

// FetchCrossSiteNetworkWithClient returns the settings Google Cloud holds for the given cross-site network using the supplied
// *compute.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see compute_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func FetchCrossSiteNetworkWithClient(ctx context.Context, service *compute.Service, projectID string, name string) (*compute.CrossSiteNetwork, error) {
	network, err := service.CrossSiteNetworks.Get(projectID, name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the cross-site network %s in project %s does not exist", name, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for cross-site network %s in project %s: %w", name, projectID, err)
	}

	return network, nil
}

// FetchWireGroup returns the settings Google Cloud holds for the given wire group, so a test can assert on what was
// actually created rather than only that it exists.
// A wire group is the pair of wires carrying traffic between two sites, so its topology and bandwidth are what it is for.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchWireGroup(t testing.TestingT, ctx context.Context, projectID string, crossSiteNetwork string, name string) *compute.WireGroup {
	group, err := FetchWireGroupE(t, ctx, projectID, crossSiteNetwork, name)
	require.NoError(t, err)

	return group
}

// FetchWireGroupE returns the settings Google Cloud holds for the given wire group.
// The ctx parameter supports cancellation and timeouts.
func FetchWireGroupE(t testing.TestingT, ctx context.Context, projectID string, crossSiteNetwork string, name string) (*compute.WireGroup, error) {
	logger.Default.Logf(t, "Getting settings for wire group %s in cross-site network %s in project %s", name, crossSiteNetwork, projectID)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return FetchWireGroupWithClient(ctx, service, projectID, crossSiteNetwork, name)
}

// FetchWireGroupWithClient returns the settings Google Cloud holds for the given wire group using the supplied
// *compute.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see compute_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func FetchWireGroupWithClient(ctx context.Context, service *compute.Service, projectID string, crossSiteNetwork string, name string) (*compute.WireGroup, error) {
	group, err := service.WireGroups.Get(projectID, crossSiteNetwork, name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the wire group %s in cross-site network %s in project %s does not exist", name, crossSiteNetwork, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for wire group %s in cross-site network %s in project %s: %w", name, crossSiteNetwork, projectID, err)
	}

	return group, nil
}

// FetchPreviewFeature returns the settings Google Cloud holds for the given preview feature, so a test can assert on what was
// actually created rather than only that it exists.
// A preview feature is a project-wide opt-in, so whether it is enabled is the only thing it says.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchPreviewFeature(t testing.TestingT, ctx context.Context, projectID string, name string) *compute.PreviewFeature {
	feature, err := FetchPreviewFeatureE(t, ctx, projectID, name)
	require.NoError(t, err)

	return feature
}

// FetchPreviewFeatureE returns the settings Google Cloud holds for the given preview feature.
// The ctx parameter supports cancellation and timeouts.
func FetchPreviewFeatureE(t testing.TestingT, ctx context.Context, projectID string, name string) (*compute.PreviewFeature, error) {
	logger.Default.Logf(t, "Getting settings for preview feature %s in project %s", name, projectID)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return FetchPreviewFeatureWithClient(ctx, service, projectID, name)
}

// FetchPreviewFeatureWithClient returns the settings Google Cloud holds for the given preview feature using the supplied
// *compute.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see compute_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func FetchPreviewFeatureWithClient(ctx context.Context, service *compute.Service, projectID string, name string) (*compute.PreviewFeature, error) {
	feature, err := service.PreviewFeatures.Get(projectID, name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the preview feature %s in project %s does not exist", name, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for preview feature %s in project %s: %w", name, projectID, err)
	}

	return feature, nil
}

// NewInstancesServiceContext creates a new InstancesService, which is used to make a subset of GCE API calls.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func NewInstancesServiceContext(t testing.TestingT, ctx context.Context) *compute.InstancesService {
	client, err := NewInstancesServiceContextE(t, ctx)
	require.NoError(t, err)

	return client
}

// NewInstancesServiceContextE creates a new InstancesService, which is used to make a subset of GCE API calls.
// The ctx parameter supports cancellation and timeouts.
func NewInstancesServiceContextE(t testing.TestingT, ctx context.Context) (*compute.InstancesService, error) {
	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get new Instances Service: %w", err)
	}

	return service.Instances, nil
}

// RandomValidGCPName returns a random, valid name for GCP resources. Many resources in GCP require lowercase letters only.
func RandomValidGCPName() string {
	id := strings.ToLower(random.UniqueID())

	return "terratest-" + id
}

// FetchProjectMetadata queries GCP to return the common instance metadata the given project holds,
// so a test can assert on what was actually set rather than only that a key exists. Metadata is a
// field of the project rather than a resource of its own, so it is read by project alone.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchProjectMetadata(t testing.TestingT, ctx context.Context, projectID string) *compute.Metadata {
	metadata, err := FetchProjectMetadataE(t, ctx, projectID)
	require.NoError(t, err)

	return metadata
}

// FetchProjectMetadataE queries GCP to return the common instance metadata the given project holds.
// The ctx parameter supports cancellation and timeouts.
func FetchProjectMetadataE(t testing.TestingT, ctx context.Context, projectID string) (*compute.Metadata, error) {
	logger.Default.Logf(t, "Getting common instance metadata for project %s", projectID)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return FetchProjectMetadataWithClient(ctx, service, projectID)
}

// FetchProjectMetadataWithClient queries GCP to return the common instance metadata the given
// project holds using the supplied *compute.Service. Prefer this variant in unit tests where the
// service is backed by an httptest fake server (see compute_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func FetchProjectMetadataWithClient(ctx context.Context, service *compute.Service, projectID string) (*compute.Metadata, error) {
	project, err := service.Projects.Get(projectID).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("Projects.Get(%s) got error: %w", projectID, err)
	}

	if project.CommonInstanceMetadata == nil {
		return nil, fmt.Errorf("project %s holds no common instance metadata", projectID)
	}

	return project.CommonInstanceMetadata, nil
}

// FetchNetworkAttachment returns the settings Google Cloud holds for the given network attachment, so
// a test can assert on what was actually created rather than only that it exists. An attachment is how
// a producer's traffic is allowed into a consumer's subnetwork, so who may connect is the whole point
// of it.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchNetworkAttachment(t testing.TestingT, ctx context.Context, projectID string, region string, name string) *compute.NetworkAttachment {
	attachment, err := FetchNetworkAttachmentE(t, ctx, projectID, region, name)
	require.NoError(t, err)

	return attachment
}

// FetchNetworkAttachmentE returns the settings Google Cloud holds for the given network attachment.
// The ctx parameter supports cancellation and timeouts.
func FetchNetworkAttachmentE(t testing.TestingT, ctx context.Context, projectID string, region string, name string) (*compute.NetworkAttachment, error) {
	logger.Default.Logf(t, "Getting settings for network attachment %s in %s in project %s", name, region, projectID)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return FetchNetworkAttachmentWithClient(ctx, service, projectID, region, name)
}

// FetchNetworkAttachmentWithClient returns the settings Google Cloud holds for the given network
// attachment using the supplied *compute.Service. Prefer this variant in unit tests where the service
// is backed by an httptest fake server (see compute_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func FetchNetworkAttachmentWithClient(ctx context.Context, service *compute.Service, projectID string, region string, name string) (*compute.NetworkAttachment, error) {
	attachment, err := service.NetworkAttachments.Get(projectID, region, name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the network attachment %s does not exist in %s in project %s", name, region, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for network attachment %s in %s in project %s: %w", name, region, projectID, err)
	}

	return attachment, nil
}

// FetchDiskIamPolicy returns the IAM policy Google Cloud holds for the given disk, so a test can
// assert on who may use it rather than only that a policy was applied.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchDiskIamPolicy(t testing.TestingT, ctx context.Context, projectID string, zone string, name string) *compute.Policy {
	policy, err := FetchDiskIamPolicyE(t, ctx, projectID, zone, name)
	require.NoError(t, err)

	return policy
}

// FetchDiskIamPolicyE returns the IAM policy Google Cloud holds for the given disk.
// The ctx parameter supports cancellation and timeouts.
func FetchDiskIamPolicyE(t testing.TestingT, ctx context.Context, projectID string, zone string, name string) (*compute.Policy, error) {
	logger.Default.Logf(t, "Getting the IAM policy for disk %s in project %s", name, projectID)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return FetchDiskIamPolicyWithClient(ctx, service, projectID, zone, name)
}

// FetchDiskIamPolicyWithClient returns the IAM policy Google Cloud holds for the given disk using
// the supplied *compute.Service. Prefer this variant in unit tests where the service is backed by an
// httptest fake server (see compute_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func FetchDiskIamPolicyWithClient(ctx context.Context, service *compute.Service, projectID string, zone string, name string) (*compute.Policy, error) {
	// A policy carrying a conditional binding is only returned in full at version 3, so that is what
	// is asked for.
	policy, err := service.Disks.GetIamPolicy(projectID, zone, name).
		OptionsRequestedPolicyVersion(iamPolicyVersionWithConditions).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the disk %s in project %s does not exist", name, projectID)
		}

		return nil, fmt.Errorf("failed to get the IAM policy for disk %s in project %s: %w", name, projectID, err)
	}

	return policy, nil
}

// FetchRegionDiskIamPolicy returns the IAM policy Google Cloud holds for the given regional disk, so a test can
// assert on who may use it rather than only that a policy was applied.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchRegionDiskIamPolicy(t testing.TestingT, ctx context.Context, projectID string, region string, name string) *compute.Policy {
	policy, err := FetchRegionDiskIamPolicyE(t, ctx, projectID, region, name)
	require.NoError(t, err)

	return policy
}

// FetchRegionDiskIamPolicyE returns the IAM policy Google Cloud holds for the given regional disk.
// The ctx parameter supports cancellation and timeouts.
func FetchRegionDiskIamPolicyE(t testing.TestingT, ctx context.Context, projectID string, region string, name string) (*compute.Policy, error) {
	logger.Default.Logf(t, "Getting the IAM policy for regional disk %s in project %s", name, projectID)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return FetchRegionDiskIamPolicyWithClient(ctx, service, projectID, region, name)
}

// FetchRegionDiskIamPolicyWithClient returns the IAM policy Google Cloud holds for the given regional disk using
// the supplied *compute.Service. Prefer this variant in unit tests where the service is backed by an
// httptest fake server (see compute_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func FetchRegionDiskIamPolicyWithClient(ctx context.Context, service *compute.Service, projectID string, region string, name string) (*compute.Policy, error) {
	// A policy carrying a conditional binding is only returned in full at version 3, so that is what
	// is asked for.
	policy, err := service.RegionDisks.GetIamPolicy(projectID, region, name).
		OptionsRequestedPolicyVersion(iamPolicyVersionWithConditions).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the regional disk %s in project %s does not exist", name, projectID)
		}

		return nil, fmt.Errorf("failed to get the IAM policy for regional disk %s in project %s: %w", name, projectID, err)
	}

	return policy, nil
}

// FetchSnapshotIamPolicy returns the IAM policy Google Cloud holds for the given snapshot, so a test can
// assert on who may use it rather than only that a policy was applied.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchSnapshotIamPolicy(t testing.TestingT, ctx context.Context, projectID string, name string) *compute.Policy {
	policy, err := FetchSnapshotIamPolicyE(t, ctx, projectID, name)
	require.NoError(t, err)

	return policy
}

// FetchSnapshotIamPolicyE returns the IAM policy Google Cloud holds for the given snapshot.
// The ctx parameter supports cancellation and timeouts.
func FetchSnapshotIamPolicyE(t testing.TestingT, ctx context.Context, projectID string, name string) (*compute.Policy, error) {
	logger.Default.Logf(t, "Getting the IAM policy for snapshot %s in project %s", name, projectID)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return FetchSnapshotIamPolicyWithClient(ctx, service, projectID, name)
}

// FetchSnapshotIamPolicyWithClient returns the IAM policy Google Cloud holds for the given snapshot using
// the supplied *compute.Service. Prefer this variant in unit tests where the service is backed by an
// httptest fake server (see compute_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func FetchSnapshotIamPolicyWithClient(ctx context.Context, service *compute.Service, projectID string, name string) (*compute.Policy, error) {
	// A policy carrying a conditional binding is only returned in full at version 3, so that is what
	// is asked for.
	policy, err := service.Snapshots.GetIamPolicy(projectID, name).
		OptionsRequestedPolicyVersion(iamPolicyVersionWithConditions).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the snapshot %s in project %s does not exist", name, projectID)
		}

		return nil, fmt.Errorf("failed to get the IAM policy for snapshot %s in project %s: %w", name, projectID, err)
	}

	return policy, nil
}

// FetchInstantSnapshotIamPolicy returns the IAM policy Google Cloud holds for the given instant snapshot, so a test can
// assert on who may use it rather than only that a policy was applied.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchInstantSnapshotIamPolicy(t testing.TestingT, ctx context.Context, projectID string, zone string, name string) *compute.Policy {
	policy, err := FetchInstantSnapshotIamPolicyE(t, ctx, projectID, zone, name)
	require.NoError(t, err)

	return policy
}

// FetchInstantSnapshotIamPolicyE returns the IAM policy Google Cloud holds for the given instant snapshot.
// The ctx parameter supports cancellation and timeouts.
func FetchInstantSnapshotIamPolicyE(t testing.TestingT, ctx context.Context, projectID string, zone string, name string) (*compute.Policy, error) {
	logger.Default.Logf(t, "Getting the IAM policy for instant snapshot %s in project %s", name, projectID)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return FetchInstantSnapshotIamPolicyWithClient(ctx, service, projectID, zone, name)
}

// FetchInstantSnapshotIamPolicyWithClient returns the IAM policy Google Cloud holds for the given instant snapshot using
// the supplied *compute.Service. Prefer this variant in unit tests where the service is backed by an
// httptest fake server (see compute_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func FetchInstantSnapshotIamPolicyWithClient(ctx context.Context, service *compute.Service, projectID string, zone string, name string) (*compute.Policy, error) {
	// A policy carrying a conditional binding is only returned in full at version 3, so that is what
	// is asked for.
	policy, err := service.InstantSnapshots.GetIamPolicy(projectID, zone, name).
		OptionsRequestedPolicyVersion(iamPolicyVersionWithConditions).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the instant snapshot %s in project %s does not exist", name, projectID)
		}

		return nil, fmt.Errorf("failed to get the IAM policy for instant snapshot %s in project %s: %w", name, projectID, err)
	}

	return policy, nil
}

// FetchRegionInstantSnapshotIamPolicy returns the IAM policy Google Cloud holds for the given regional instant snapshot, so a test can
// assert on who may use it rather than only that a policy was applied.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchRegionInstantSnapshotIamPolicy(t testing.TestingT, ctx context.Context, projectID string, region string, name string) *compute.Policy {
	policy, err := FetchRegionInstantSnapshotIamPolicyE(t, ctx, projectID, region, name)
	require.NoError(t, err)

	return policy
}

// FetchRegionInstantSnapshotIamPolicyE returns the IAM policy Google Cloud holds for the given regional instant snapshot.
// The ctx parameter supports cancellation and timeouts.
func FetchRegionInstantSnapshotIamPolicyE(t testing.TestingT, ctx context.Context, projectID string, region string, name string) (*compute.Policy, error) {
	logger.Default.Logf(t, "Getting the IAM policy for regional instant snapshot %s in project %s", name, projectID)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return FetchRegionInstantSnapshotIamPolicyWithClient(ctx, service, projectID, region, name)
}

// FetchRegionInstantSnapshotIamPolicyWithClient returns the IAM policy Google Cloud holds for the given regional instant snapshot using
// the supplied *compute.Service. Prefer this variant in unit tests where the service is backed by an
// httptest fake server (see compute_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func FetchRegionInstantSnapshotIamPolicyWithClient(ctx context.Context, service *compute.Service, projectID string, region string, name string) (*compute.Policy, error) {
	// A policy carrying a conditional binding is only returned in full at version 3, so that is what
	// is asked for.
	policy, err := service.RegionInstantSnapshots.GetIamPolicy(projectID, region, name).
		OptionsRequestedPolicyVersion(iamPolicyVersionWithConditions).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the regional instant snapshot %s in project %s does not exist", name, projectID)
		}

		return nil, fmt.Errorf("failed to get the IAM policy for regional instant snapshot %s in project %s: %w", name, projectID, err)
	}

	return policy, nil
}

// FetchInstanceTemplateIamPolicy returns the IAM policy Google Cloud holds for the given instance template, so a test can
// assert on who may use it rather than only that a policy was applied.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchInstanceTemplateIamPolicy(t testing.TestingT, ctx context.Context, projectID string, name string) *compute.Policy {
	policy, err := FetchInstanceTemplateIamPolicyE(t, ctx, projectID, name)
	require.NoError(t, err)

	return policy
}

// FetchInstanceTemplateIamPolicyE returns the IAM policy Google Cloud holds for the given instance template.
// The ctx parameter supports cancellation and timeouts.
func FetchInstanceTemplateIamPolicyE(t testing.TestingT, ctx context.Context, projectID string, name string) (*compute.Policy, error) {
	logger.Default.Logf(t, "Getting the IAM policy for instance template %s in project %s", name, projectID)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return FetchInstanceTemplateIamPolicyWithClient(ctx, service, projectID, name)
}

// FetchInstanceTemplateIamPolicyWithClient returns the IAM policy Google Cloud holds for the given instance template using
// the supplied *compute.Service. Prefer this variant in unit tests where the service is backed by an
// httptest fake server (see compute_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func FetchInstanceTemplateIamPolicyWithClient(ctx context.Context, service *compute.Service, projectID string, name string) (*compute.Policy, error) {
	// A policy carrying a conditional binding is only returned in full at version 3, so that is what
	// is asked for.
	policy, err := service.InstanceTemplates.GetIamPolicy(projectID, name).
		OptionsRequestedPolicyVersion(iamPolicyVersionWithConditions).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the instance template %s in project %s does not exist", name, projectID)
		}

		return nil, fmt.Errorf("failed to get the IAM policy for instance template %s in project %s: %w", name, projectID, err)
	}

	return policy, nil
}

// FetchInstanceIamPolicy returns the IAM policy Google Cloud holds for the given instance, so a test can
// assert on who may use it rather than only that a policy was applied.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchInstanceIamPolicy(t testing.TestingT, ctx context.Context, projectID string, zone string, name string) *compute.Policy {
	policy, err := FetchInstanceIamPolicyE(t, ctx, projectID, zone, name)
	require.NoError(t, err)

	return policy
}

// FetchInstanceIamPolicyE returns the IAM policy Google Cloud holds for the given instance.
// The ctx parameter supports cancellation and timeouts.
func FetchInstanceIamPolicyE(t testing.TestingT, ctx context.Context, projectID string, zone string, name string) (*compute.Policy, error) {
	logger.Default.Logf(t, "Getting the IAM policy for instance %s in project %s", name, projectID)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return FetchInstanceIamPolicyWithClient(ctx, service, projectID, zone, name)
}

// FetchInstanceIamPolicyWithClient returns the IAM policy Google Cloud holds for the given instance using
// the supplied *compute.Service. Prefer this variant in unit tests where the service is backed by an
// httptest fake server (see compute_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func FetchInstanceIamPolicyWithClient(ctx context.Context, service *compute.Service, projectID string, zone string, name string) (*compute.Policy, error) {
	// A policy carrying a conditional binding is only returned in full at version 3, so that is what
	// is asked for.
	policy, err := service.Instances.GetIamPolicy(projectID, zone, name).
		OptionsRequestedPolicyVersion(iamPolicyVersionWithConditions).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the instance %s in project %s does not exist", name, projectID)
		}

		return nil, fmt.Errorf("failed to get the IAM policy for instance %s in project %s: %w", name, projectID, err)
	}

	return policy, nil
}

// FetchInstantSnapshot returns the settings Google Cloud holds for the given instant snapshot, so a test can assert on
// what was actually created rather than only that it exists. An instant snapshot keeps only the blocks that changed, so it is cheap to take and lives beside its disk rather than in a bucket.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchInstantSnapshot(t testing.TestingT, ctx context.Context, projectID string, zone string, name string) *compute.InstantSnapshot {
	result, err := FetchInstantSnapshotE(t, ctx, projectID, zone, name)
	require.NoError(t, err)

	return result
}

// FetchInstantSnapshotE returns the settings Google Cloud holds for the given instant snapshot.
// The ctx parameter supports cancellation and timeouts.
func FetchInstantSnapshotE(t testing.TestingT, ctx context.Context, projectID string, zone string, name string) (*compute.InstantSnapshot, error) {
	logger.Default.Logf(t, "Getting settings for instant snapshot %s in project %s", name, projectID)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return FetchInstantSnapshotWithClient(ctx, service, projectID, zone, name)
}

// FetchInstantSnapshotWithClient returns the settings Google Cloud holds for the given instant snapshot using the supplied
// *compute.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see compute_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func FetchInstantSnapshotWithClient(ctx context.Context, service *compute.Service, projectID string, zone string, name string) (*compute.InstantSnapshot, error) {
	result, err := service.InstantSnapshots.Get(projectID, zone, name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the instant snapshot %s in project %s does not exist", name, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for instant snapshot %s in project %s: %w", name, projectID, err)
	}

	return result, nil
}

// FetchRegionInstantSnapshot returns the settings Google Cloud holds for the given regional instant snapshot, so a test can assert on
// what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchRegionInstantSnapshot(t testing.TestingT, ctx context.Context, projectID string, region string, name string) *compute.InstantSnapshot {
	result, err := FetchRegionInstantSnapshotE(t, ctx, projectID, region, name)
	require.NoError(t, err)

	return result
}

// FetchRegionInstantSnapshotE returns the settings Google Cloud holds for the given regional instant snapshot.
// The ctx parameter supports cancellation and timeouts.
func FetchRegionInstantSnapshotE(t testing.TestingT, ctx context.Context, projectID string, region string, name string) (*compute.InstantSnapshot, error) {
	logger.Default.Logf(t, "Getting settings for regional instant snapshot %s in project %s", name, projectID)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return FetchRegionInstantSnapshotWithClient(ctx, service, projectID, region, name)
}

// FetchRegionInstantSnapshotWithClient returns the settings Google Cloud holds for the given regional instant snapshot using the supplied
// *compute.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see compute_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func FetchRegionInstantSnapshotWithClient(ctx context.Context, service *compute.Service, projectID string, region string, name string) (*compute.InstantSnapshot, error) {
	result, err := service.RegionInstantSnapshots.Get(projectID, region, name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the regional instant snapshot %s in project %s does not exist", name, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for regional instant snapshot %s in project %s: %w", name, projectID, err)
	}

	return result, nil
}

// FetchNodeTemplate returns the settings Google Cloud holds for the given sole tenant node template, so a test can assert on
// what was actually created rather than only that it exists. A template describes the nodes a group would run; on its own it runs none.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchNodeTemplate(t testing.TestingT, ctx context.Context, projectID string, region string, name string) *compute.NodeTemplate {
	result, err := FetchNodeTemplateE(t, ctx, projectID, region, name)
	require.NoError(t, err)

	return result
}

// FetchNodeTemplateE returns the settings Google Cloud holds for the given sole tenant node template.
// The ctx parameter supports cancellation and timeouts.
func FetchNodeTemplateE(t testing.TestingT, ctx context.Context, projectID string, region string, name string) (*compute.NodeTemplate, error) {
	logger.Default.Logf(t, "Getting settings for sole tenant node template %s in project %s", name, projectID)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return FetchNodeTemplateWithClient(ctx, service, projectID, region, name)
}

// FetchNodeTemplateWithClient returns the settings Google Cloud holds for the given sole tenant node template using the supplied
// *compute.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see compute_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func FetchNodeTemplateWithClient(ctx context.Context, service *compute.Service, projectID string, region string, name string) (*compute.NodeTemplate, error) {
	result, err := service.NodeTemplates.Get(projectID, region, name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the sole tenant node template %s in project %s does not exist", name, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for sole tenant node template %s in project %s: %w", name, projectID, err)
	}

	return result, nil
}

// FetchRegionAutoscaler returns the settings Google Cloud holds for the given regional autoscaler, so a test can assert on
// what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchRegionAutoscaler(t testing.TestingT, ctx context.Context, projectID string, region string, name string) *compute.Autoscaler {
	result, err := FetchRegionAutoscalerE(t, ctx, projectID, region, name)
	require.NoError(t, err)

	return result
}

// FetchRegionAutoscalerE returns the settings Google Cloud holds for the given regional autoscaler.
// The ctx parameter supports cancellation and timeouts.
func FetchRegionAutoscalerE(t testing.TestingT, ctx context.Context, projectID string, region string, name string) (*compute.Autoscaler, error) {
	logger.Default.Logf(t, "Getting settings for regional autoscaler %s in project %s", name, projectID)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return FetchRegionAutoscalerWithClient(ctx, service, projectID, region, name)
}

// FetchRegionAutoscalerWithClient returns the settings Google Cloud holds for the given regional autoscaler using the supplied
// *compute.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see compute_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func FetchRegionAutoscalerWithClient(ctx context.Context, service *compute.Service, projectID string, region string, name string) (*compute.Autoscaler, error) {
	result, err := service.RegionAutoscalers.Get(projectID, region, name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the regional autoscaler %s in project %s does not exist", name, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for regional autoscaler %s in project %s: %w", name, projectID, err)
	}

	return result, nil
}

// FetchInstanceResourcePolicies returns the resource policies Google Cloud holds against the given
// instance, so a test can assert that an attachment landed. An attachment is not a resource of its
// own: Google answers for it as a list on the instance, so that is what this reads.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchInstanceResourcePolicies(t testing.TestingT, ctx context.Context, projectID string, zone string, name string) []string {
	policies, err := FetchInstanceResourcePoliciesE(t, ctx, projectID, zone, name)
	require.NoError(t, err)

	return policies
}

// FetchInstanceResourcePoliciesE returns the resource policies Google Cloud holds against the given
// instance.
// The ctx parameter supports cancellation and timeouts.
func FetchInstanceResourcePoliciesE(t testing.TestingT, ctx context.Context, projectID string, zone string, name string) ([]string, error) {
	logger.Default.Logf(t, "Getting the resource policies attached to instance %s in project %s", name, projectID)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return FetchInstanceResourcePoliciesWithClient(ctx, service, projectID, zone, name)
}

// FetchInstanceResourcePoliciesWithClient returns the resource policies Google Cloud holds against
// the given instance using the supplied *compute.Service. Prefer this variant in unit tests where the
// service is backed by an httptest fake server (see compute_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func FetchInstanceResourcePoliciesWithClient(ctx context.Context, service *compute.Service, projectID string, zone string, name string) ([]string, error) {
	instance, err := service.Instances.Get(projectID, zone, name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the instance %s in project %s does not exist", name, projectID)
		}

		return nil, fmt.Errorf("failed to get the resource policies attached to instance %s in project %s: %w", name, projectID, err)
	}

	return instance.ResourcePolicies, nil
}

// NewComputeServiceContextE creates a new Compute service, which is used to make GCE API calls.
// The ctx parameter supports cancellation and timeouts.
func NewComputeServiceContextE(t testing.TestingT, ctx context.Context) (*compute.Service, error) {
	if ts, ok := getStaticTokenSource(); ok {
		return compute.NewService(ctx, option.WithTokenSource(ts))
	}

	description := "Attempting to request a Google OAuth2 token"
	maxRetries := 6

	var client *http.Client

	msg, retryErr := retry.DoWithRetryContextE(t, ctx, description, maxRetries, defaultRetryInterval, func() (string, error) {
		rawClient, err := google.DefaultClient(ctx, compute.CloudPlatformScope)
		if err != nil {
			return "Error retrieving default GCP client", err
		}

		client = rawClient

		return "Successfully retrieved default GCP client", nil
	})

	logger.Default.Logf(t, "%s", msg)

	if retryErr != nil {
		return nil, retryErr
	}

	return compute.NewService(ctx, option.WithHTTPClient(client))
}
