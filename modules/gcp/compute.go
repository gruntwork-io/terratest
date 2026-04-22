package gcp

import (
	"context"
	"fmt"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/retry"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/compute/v1"
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
	// GetInstanceIDs gets the IDs of Instances in the given Instance Group.
	// This will fail the test if there is an error.
	GetInstanceIDs(t testing.TestingT) []string

	// GetInstanceIDsE gets the IDs of Instances in the given Instance Group.
	GetInstanceIDsE(t testing.TestingT) ([]string, error)

	// GetInstanceIDsContextE gets the IDs of Instances in the given Instance Group.
	// The ctx parameter supports cancellation and timeouts.
	GetInstanceIDsContextE(t testing.TestingT, ctx context.Context) ([]string, error)

	// Deprecated: Use [InstanceGroup.GetInstanceIDs] instead.
	GetInstanceIds(t testing.TestingT) []string //nolint:staticcheck,revive // preserving deprecated method name

	// Deprecated: Use [InstanceGroup.GetInstanceIDsE] instead.
	GetInstanceIdsE(t testing.TestingT) ([]string, error) //nolint:staticcheck,revive // preserving deprecated method name
}

// FetchInstance queries GCP to return an instance of the Compute Instance type.
// This will fail the test if there is an error.
//
// Deprecated: Use [FetchInstanceContext] instead.
func FetchInstance(t testing.TestingT, projectID string, name string) *Instance {
	return FetchInstanceContext(t, context.Background(), projectID, name)
}

// FetchInstanceContext queries GCP to return an instance of the Compute Instance type.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchInstanceContext(t testing.TestingT, ctx context.Context, projectID string, name string) *Instance {
	instance, err := FetchInstanceContextE(t, ctx, projectID, name)
	require.NoError(t, err)

	return instance
}

// FetchInstanceE queries GCP to return an instance of the Compute Instance type.
//
// Deprecated: Use [FetchInstanceContextE] instead.
func FetchInstanceE(t testing.TestingT, projectID string, name string) (*Instance, error) {
	return FetchInstanceContextE(t, context.Background(), projectID, name)
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
	// If we want to fetch an Instance without knowing its Zone, we have to query GCP for all Instances in the project
	// and match on name.
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

// FetchImage queries GCP to return a new instance of the Compute Image type.
// This will fail the test if there is an error.
//
// Deprecated: Use [FetchImageContext] instead.
func FetchImage(t testing.TestingT, projectID string, name string) *Image {
	return FetchImageContext(t, context.Background(), projectID, name)
}

// FetchImageContext queries GCP to return a new instance of the Compute Image type.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchImageContext(t testing.TestingT, ctx context.Context, projectID string, name string) *Image {
	image, err := FetchImageContextE(t, ctx, projectID, name)
	require.NoError(t, err)

	return image
}

// FetchImageE queries GCP to return a new instance of the Compute Image type.
//
// Deprecated: Use [FetchImageContextE] instead.
func FetchImageE(t testing.TestingT, projectID string, name string) (*Image, error) {
	return FetchImageContextE(t, context.Background(), projectID, name)
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

// FetchRegionalInstanceGroup queries GCP to return a new instance of the Regional Instance Group type.
// This will fail the test if there is an error.
//
// Deprecated: Use [FetchRegionalInstanceGroupContext] instead.
func FetchRegionalInstanceGroup(t testing.TestingT, projectID string, region string, name string) *RegionalInstanceGroup {
	return FetchRegionalInstanceGroupContext(t, context.Background(), projectID, region, name)
}

// FetchRegionalInstanceGroupContext queries GCP to return a new instance of the Regional Instance Group type.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchRegionalInstanceGroupContext(t testing.TestingT, ctx context.Context, projectID string, region string, name string) *RegionalInstanceGroup {
	instanceGroup, err := FetchRegionalInstanceGroupContextE(t, ctx, projectID, region, name)
	require.NoError(t, err)

	return instanceGroup
}

// FetchRegionalInstanceGroupE queries GCP to return a new instance of the Regional Instance Group type.
//
// Deprecated: Use [FetchRegionalInstanceGroupContextE] instead.
func FetchRegionalInstanceGroupE(t testing.TestingT, projectID string, region string, name string) (*RegionalInstanceGroup, error) {
	return FetchRegionalInstanceGroupContextE(t, context.Background(), projectID, region, name)
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

// FetchZonalInstanceGroup queries GCP to return a new instance of the Zonal Instance Group type.
// This will fail the test if there is an error.
//
// Deprecated: Use [FetchZonalInstanceGroupContext] instead.
func FetchZonalInstanceGroup(t testing.TestingT, projectID string, zone string, name string) *ZonalInstanceGroup {
	return FetchZonalInstanceGroupContext(t, context.Background(), projectID, zone, name)
}

// FetchZonalInstanceGroupContext queries GCP to return a new instance of the Zonal Instance Group type.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchZonalInstanceGroupContext(t testing.TestingT, ctx context.Context, projectID string, zone string, name string) *ZonalInstanceGroup {
	instanceGroup, err := FetchZonalInstanceGroupContextE(t, ctx, projectID, zone, name)
	require.NoError(t, err)

	return instanceGroup
}

// FetchZonalInstanceGroupE queries GCP to return a new instance of the Zonal Instance Group type.
//
// Deprecated: Use [FetchZonalInstanceGroupContextE] instead.
func FetchZonalInstanceGroupE(t testing.TestingT, projectID string, zone string, name string) (*ZonalInstanceGroup, error) {
	return FetchZonalInstanceGroupContextE(t, context.Background(), projectID, zone, name)
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

// GetPublicIP gets the public IP address of the given Compute Instance.
// This will fail the test if there is an error.
//
// Deprecated: Use [Instance.GetPublicIPContext] instead.
func (i *Instance) GetPublicIP(t testing.TestingT) string {
	return i.GetPublicIPContext(t, context.Background())
}

// GetPublicIPContext gets the public IP address of the given Compute Instance.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func (i *Instance) GetPublicIPContext(t testing.TestingT, ctx context.Context) string {
	ip, err := i.GetPublicIPContextE(t, ctx)
	require.NoError(t, err)

	return ip
}

// GetPublicIPE gets the public IP address of the given Compute Instance.
//
// Deprecated: Use [Instance.GetPublicIPContextE] instead.
func (i *Instance) GetPublicIPE(t testing.TestingT) (string, error) {
	return i.GetPublicIPContextE(t, context.Background())
}

// GetPublicIp gets the public IP address of the given Compute Instance.
// This will fail the test if there is an error.
//
// Deprecated: Use [Instance.GetPublicIP] instead.
func (i *Instance) GetPublicIp(t testing.TestingT) string { //nolint:staticcheck,revive // preserving deprecated method name
	return i.GetPublicIP(t)
}

// GetPublicIpE gets the public IP address of the given Compute Instance.
//
// Deprecated: Use [Instance.GetPublicIPE] instead.
func (i *Instance) GetPublicIpE(t testing.TestingT) (string, error) { //nolint:staticcheck,revive // preserving deprecated method name
	return i.GetPublicIPE(t)
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

// GetLabels returns all the tags for the given Compute Instance.
//
// Deprecated: Use [Instance.GetLabelsContext] instead.
func (i *Instance) GetLabels(t testing.TestingT) map[string]string {
	return i.GetLabelsContext(t, context.Background())
}

// GetLabelsContext returns all the tags for the given Compute Instance.
// The ctx parameter supports cancellation and timeouts.
func (i *Instance) GetLabelsContext(t testing.TestingT, ctx context.Context) map[string]string {
	labels, err := i.GetLabelsContextE(t, ctx)
	require.NoError(t, err)

	return labels
}

// GetLabelsE returns all the tags for the given Compute Instance.
//
// Deprecated: Use [Instance.GetLabelsContextE] instead.
func (i *Instance) GetLabelsE(t testing.TestingT) (map[string]string, error) {
	return i.GetLabelsContextE(t, context.Background())
}

// GetLabelsContextE returns all the tags for the given Compute Instance.
// The ctx parameter supports cancellation and timeouts.
func (i *Instance) GetLabelsContextE(t testing.TestingT, ctx context.Context) (map[string]string, error) {
	return i.Labels, nil
}

// GetZone returns the Zone in which the Compute Instance is located.
//
// Deprecated: Use [Instance.GetZoneContext] instead.
func (i *Instance) GetZone(t testing.TestingT) string {
	return i.GetZoneContext(t, context.Background())
}

// GetZoneContext returns the Zone in which the Compute Instance is located.
// The ctx parameter supports cancellation and timeouts.
func (i *Instance) GetZoneContext(t testing.TestingT, ctx context.Context) string {
	zone, err := i.GetZoneContextE(t, ctx)
	require.NoError(t, err)

	return zone
}

// GetZoneE returns the Zone in which the Compute Instance is located.
//
// Deprecated: Use [Instance.GetZoneContextE] instead.
func (i *Instance) GetZoneE(t testing.TestingT) (string, error) {
	return i.GetZoneContextE(t, context.Background())
}

// GetZoneContextE returns the Zone in which the Compute Instance is located.
// The ctx parameter supports cancellation and timeouts.
func (i *Instance) GetZoneContextE(t testing.TestingT, ctx context.Context) (string, error) {
	return ZoneURLToZone(i.Zone), nil
}

// SetLabels adds the tags to the given Compute Instance.
// This will fail the test if there is an error.
//
// Deprecated: Use [Instance.SetLabelsContext] instead.
func (i *Instance) SetLabels(t testing.TestingT, labels map[string]string) {
	i.SetLabelsContext(t, context.Background(), labels)
}

// SetLabelsContext adds the tags to the given Compute Instance.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func (i *Instance) SetLabelsContext(t testing.TestingT, ctx context.Context, labels map[string]string) {
	err := i.SetLabelsContextE(t, ctx, labels)
	require.NoError(t, err)
}

// SetLabelsE adds the tags to the given Compute Instance.
//
// Deprecated: Use [Instance.SetLabelsContextE] instead.
func (i *Instance) SetLabelsE(t testing.TestingT, labels map[string]string) error {
	return i.SetLabelsContextE(t, context.Background(), labels)
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

// SetLabelsWithClient adds the tags to the given Compute Instance using the supplied *compute.Service.
// Prefer this variant in unit tests where the service is backed by an httptest fake server
// (see compute_unit_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func (i *Instance) SetLabelsWithClient(ctx context.Context, service *compute.Service, labels map[string]string) error {
	req := compute.InstancesSetLabelsRequest{Labels: labels, LabelFingerprint: i.LabelFingerprint}

	if _, err := service.Instances.SetLabels(i.projectID, ZoneURLToZone(i.Zone), i.Name, &req).Context(ctx).Do(); err != nil {
		return fmt.Errorf("Instances.SetLabels(%s) got error: %w", i.Name, err)
	}

	return nil
}

// GetMetadata gets the given Compute Instance's metadata.
//
// Deprecated: Use [Instance.GetMetadataContext] instead.
func (i *Instance) GetMetadata(t testing.TestingT) []*compute.MetadataItems {
	return i.GetMetadataContext(t, context.Background())
}

// GetMetadataContext gets the given Compute Instance's metadata.
// The ctx parameter supports cancellation and timeouts.
func (i *Instance) GetMetadataContext(t testing.TestingT, ctx context.Context) []*compute.MetadataItems {
	metadata, err := i.GetMetadataContextE(t, ctx)
	require.NoError(t, err)

	return metadata
}

// GetMetadataE gets the given Compute Instance's metadata.
//
// Deprecated: Use [Instance.GetMetadataContextE] instead.
func (i *Instance) GetMetadataE(t testing.TestingT) ([]*compute.MetadataItems, error) {
	return i.GetMetadataContextE(t, context.Background())
}

// GetMetadataContextE gets the given Compute Instance's metadata.
// The ctx parameter supports cancellation and timeouts.
func (i *Instance) GetMetadataContextE(t testing.TestingT, ctx context.Context) ([]*compute.MetadataItems, error) {
	return i.Metadata.Items, nil
}

// SetMetadata sets the given Compute Instance's metadata.
// This will fail the test if there is an error.
//
// Deprecated: Use [Instance.SetMetadataContext] instead.
func (i *Instance) SetMetadata(t testing.TestingT, metadata map[string]string) {
	i.SetMetadataContext(t, context.Background(), metadata)
}

// SetMetadataContext sets the given Compute Instance's metadata.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func (i *Instance) SetMetadataContext(t testing.TestingT, ctx context.Context, metadata map[string]string) {
	err := i.SetMetadataContextE(t, ctx, metadata)
	require.NoError(t, err)
}

// SetMetadataE adds the given metadata map to the existing metadata of the given Compute Instance.
//
// Deprecated: Use [Instance.SetMetadataContextE] instead.
func (i *Instance) SetMetadataE(t testing.TestingT, metadata map[string]string) error {
	return i.SetMetadataContextE(t, context.Background(), metadata)
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

// AddSSHKey adds the given public SSH key to the Compute Instance. Users can SSH in with the given username.
// This will fail the test if there is an error.
//
// Deprecated: Use [Instance.AddSSHKeyContext] instead.
func (i *Instance) AddSSHKey(t testing.TestingT, username string, publicKey string) {
	i.AddSSHKeyContext(t, context.Background(), username, publicKey)
}

// AddSSHKeyContext adds the given public SSH key to the Compute Instance. Users can SSH in with the given username.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func (i *Instance) AddSSHKeyContext(t testing.TestingT, ctx context.Context, username string, publicKey string) {
	err := i.AddSSHKeyContextE(t, ctx, username, publicKey)
	require.NoError(t, err)
}

// AddSSHKeyE adds the given public SSH key to the Compute Instance. Users can SSH in with the given username.
//
// Deprecated: Use [Instance.AddSSHKeyContextE] instead.
func (i *Instance) AddSSHKeyE(t testing.TestingT, username string, publicKey string) error {
	return i.AddSSHKeyContextE(t, context.Background(), username, publicKey)
}

// AddSshKey adds the given public SSH key to the Compute Instance. Users can SSH in with the given username.
// This will fail the test if there is an error.
//
// Deprecated: Use [Instance.AddSSHKey] instead.
func (i *Instance) AddSshKey(t testing.TestingT, username string, publicKey string) { //nolint:staticcheck,revive // preserving deprecated method name
	i.AddSSHKey(t, username, publicKey)
}

// AddSshKeyE adds the given public SSH key to the Compute Instance. Users can SSH in with the given username.
//
// Deprecated: Use [Instance.AddSSHKeyE] instead.
func (i *Instance) AddSshKeyE(t testing.TestingT, username string, publicKey string) error { //nolint:staticcheck,revive // preserving deprecated method name
	return i.AddSSHKeyE(t, username, publicKey)
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
	// We represent the key in the format required per GCP docs (https://cloud.google.com/compute/docs/instances/adding-removing-ssh-keys)
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

// DeleteImage deletes the given Compute Image.
// This will fail the test if there is an error.
//
// Deprecated: Use [Image.DeleteImageContext] instead.
func (i *Image) DeleteImage(t testing.TestingT) {
	i.DeleteImageContext(t, context.Background())
}

// DeleteImageContext deletes the given Compute Image.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func (i *Image) DeleteImageContext(t testing.TestingT, ctx context.Context) {
	err := i.DeleteImageContextE(t, ctx)
	require.NoError(t, err)
}

// DeleteImageE deletes the given Compute Image.
//
// Deprecated: Use [Image.DeleteImageContextE] instead.
func (i *Image) DeleteImageE(t testing.TestingT) error {
	return i.DeleteImageContextE(t, context.Background())
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

// GetInstanceIDs gets the IDs of Instances in the given Zonal Instance Group.
// This will fail the test if there is an error.
//
// Deprecated: Use [ZonalInstanceGroup.GetInstanceIDsContext] instead.
func (ig *ZonalInstanceGroup) GetInstanceIDs(t testing.TestingT) []string {
	return ig.GetInstanceIDsContext(t, context.Background())
}

// GetInstanceIDsContext gets the IDs of Instances in the given Zonal Instance Group.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func (ig *ZonalInstanceGroup) GetInstanceIDsContext(t testing.TestingT, ctx context.Context) []string {
	ids, err := ig.GetInstanceIDsContextE(t, ctx)
	require.NoError(t, err)

	return ids
}

// GetInstanceIDsE gets the IDs of Instances in the given Zonal Instance Group.
//
// Deprecated: Use [ZonalInstanceGroup.GetInstanceIDsContextE] instead.
func (ig *ZonalInstanceGroup) GetInstanceIDsE(t testing.TestingT) ([]string, error) {
	return ig.GetInstanceIDsContextE(t, context.Background())
}

// GetInstanceIds gets the IDs of Instances in the given Zonal Instance Group.
// This will fail the test if there is an error.
//
// Deprecated: Use [ZonalInstanceGroup.GetInstanceIDs] instead.
func (ig *ZonalInstanceGroup) GetInstanceIds(t testing.TestingT) []string { //nolint:staticcheck,revive // preserving deprecated method name
	return ig.GetInstanceIDs(t)
}

// GetInstanceIdsE gets the IDs of Instances in the given Zonal Instance Group.
//
// Deprecated: Use [ZonalInstanceGroup.GetInstanceIDsE] instead.
func (ig *ZonalInstanceGroup) GetInstanceIdsE(t testing.TestingT) ([]string, error) { //nolint:staticcheck,revive // preserving deprecated method name
	return ig.GetInstanceIDsE(t)
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
func (ig *ZonalInstanceGroup) GetInstanceIDsWithClient(ctx context.Context, service *compute.Service) ([]string, error) { //nolint:dupl // zonal and regional implementations differ in API types
	requestBody := &compute.InstanceGroupsListInstancesRequest{
		InstanceState: "ALL",
	}

	instanceIDs := []string{}
	zone := ZoneURLToZone(ig.Zone)
	req := service.InstanceGroups.ListInstances(ig.projectID, zone, ig.Name, requestBody)

	err := req.Pages(ctx, func(page *compute.InstanceGroupsListInstances) error {
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
		return nil, fmt.Errorf("InstanceGroups.ListInstances(%s) got error: %w", ig.Name, err)
	}

	return instanceIDs, nil
}

// GetInstanceIDs gets the IDs of Instances in the given Regional Instance Group.
// This will fail the test if there is an error.
//
// Deprecated: Use [RegionalInstanceGroup.GetInstanceIDsContext] instead.
func (ig *RegionalInstanceGroup) GetInstanceIDs(t testing.TestingT) []string {
	return ig.GetInstanceIDsContext(t, context.Background())
}

// GetInstanceIDsContext gets the IDs of Instances in the given Regional Instance Group.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func (ig *RegionalInstanceGroup) GetInstanceIDsContext(t testing.TestingT, ctx context.Context) []string {
	ids, err := ig.GetInstanceIDsContextE(t, ctx)
	require.NoError(t, err)

	return ids
}

// GetInstanceIDsE gets the IDs of Instances in the given Regional Instance Group.
//
// Deprecated: Use [RegionalInstanceGroup.GetInstanceIDsContextE] instead.
func (ig *RegionalInstanceGroup) GetInstanceIDsE(t testing.TestingT) ([]string, error) {
	return ig.GetInstanceIDsContextE(t, context.Background())
}

// GetInstanceIds gets the IDs of Instances in the given Regional Instance Group.
// This will fail the test if there is an error.
//
// Deprecated: Use [RegionalInstanceGroup.GetInstanceIDs] instead.
func (ig *RegionalInstanceGroup) GetInstanceIds(t testing.TestingT) []string { //nolint:staticcheck,revive // preserving deprecated method name
	return ig.GetInstanceIDs(t)
}

// GetInstanceIdsE gets the IDs of Instances in the given Regional Instance Group.
//
// Deprecated: Use [RegionalInstanceGroup.GetInstanceIDsE] instead.
func (ig *RegionalInstanceGroup) GetInstanceIdsE(t testing.TestingT) ([]string, error) { //nolint:staticcheck,revive // preserving deprecated method name
	return ig.GetInstanceIDsE(t)
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
func (ig *RegionalInstanceGroup) GetInstanceIDsWithClient(ctx context.Context, service *compute.Service) ([]string, error) { //nolint:dupl // zonal and regional implementations differ in API types
	requestBody := &compute.RegionInstanceGroupsListInstancesRequest{
		InstanceState: "ALL",
	}

	instanceIDs := []string{}
	region := RegionURLToRegion(ig.Region)
	req := service.RegionInstanceGroups.ListInstances(ig.projectID, region, ig.Name, requestBody)

	err := req.Pages(ctx, func(page *compute.RegionInstanceGroupsListInstances) error {
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
		return nil, fmt.Errorf("InstanceGroups.ListInstances(%s) got error: %w", ig.Name, err)
	}

	return instanceIDs, nil
}

// GetInstances returns a collection of Instance structs from the given Zonal Instance Group.
// This will fail the test if there is an error.
//
// Deprecated: Use [ZonalInstanceGroup.GetInstancesContext] instead.
func (ig *ZonalInstanceGroup) GetInstances(t testing.TestingT, projectID string) []*Instance {
	return ig.GetInstancesContext(t, context.Background(), projectID)
}

// GetInstancesContext returns a collection of Instance structs from the given Zonal Instance Group.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func (ig *ZonalInstanceGroup) GetInstancesContext(t testing.TestingT, ctx context.Context, projectID string) []*Instance {
	instances, err := ig.GetInstancesContextE(t, ctx, projectID)
	require.NoError(t, err)

	return instances
}

// GetInstancesE returns a collection of Instance structs from the given Zonal Instance Group.
//
// Deprecated: Use [ZonalInstanceGroup.GetInstancesContextE] instead.
func (ig *ZonalInstanceGroup) GetInstancesE(t testing.TestingT, projectID string) ([]*Instance, error) {
	return ig.GetInstancesContextE(t, context.Background(), projectID)
}

// GetInstancesContextE returns a collection of Instance structs from the given Zonal Instance Group.
// The ctx parameter supports cancellation and timeouts.
func (ig *ZonalInstanceGroup) GetInstancesContextE(t testing.TestingT, ctx context.Context, projectID string) ([]*Instance, error) {
	return getInstancesContextE(t, ctx, ig, projectID)
}

// GetInstances returns a collection of Instance structs from the given Regional Instance Group.
// This will fail the test if there is an error.
//
// Deprecated: Use [RegionalInstanceGroup.GetInstancesContext] instead.
func (ig *RegionalInstanceGroup) GetInstances(t testing.TestingT, projectID string) []*Instance {
	return ig.GetInstancesContext(t, context.Background(), projectID)
}

// GetInstancesContext returns a collection of Instance structs from the given Regional Instance Group.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func (ig *RegionalInstanceGroup) GetInstancesContext(t testing.TestingT, ctx context.Context, projectID string) []*Instance {
	instances, err := ig.GetInstancesContextE(t, ctx, projectID)
	require.NoError(t, err)

	return instances
}

// GetInstancesE returns a collection of Instance structs from the given Regional Instance Group.
//
// Deprecated: Use [RegionalInstanceGroup.GetInstancesContextE] instead.
func (ig *RegionalInstanceGroup) GetInstancesE(t testing.TestingT, projectID string) ([]*Instance, error) {
	return ig.GetInstancesContextE(t, context.Background(), projectID)
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

// GetPublicIPs returns a slice of the public IPs from the given Zonal Instance Group.
// This will fail the test if there is an error.
//
// Deprecated: Use [ZonalInstanceGroup.GetPublicIPsContext] instead.
func (ig *ZonalInstanceGroup) GetPublicIPs(t testing.TestingT, projectID string) []string {
	return ig.GetPublicIPsContext(t, context.Background(), projectID)
}

// GetPublicIPsContext returns a slice of the public IPs from the given Zonal Instance Group.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func (ig *ZonalInstanceGroup) GetPublicIPsContext(t testing.TestingT, ctx context.Context, projectID string) []string {
	ips, err := ig.GetPublicIPsContextE(t, ctx, projectID)
	require.NoError(t, err)

	return ips
}

// GetPublicIPsE returns a slice of the public IPs from the given Zonal Instance Group.
//
// Deprecated: Use [ZonalInstanceGroup.GetPublicIPsContextE] instead.
func (ig *ZonalInstanceGroup) GetPublicIPsE(t testing.TestingT, projectID string) ([]string, error) {
	return ig.GetPublicIPsContextE(t, context.Background(), projectID)
}

// GetPublicIps returns a slice of the public IPs from the given Zonal Instance Group.
// This will fail the test if there is an error.
//
// Deprecated: Use [ZonalInstanceGroup.GetPublicIPs] instead.
func (ig *ZonalInstanceGroup) GetPublicIps(t testing.TestingT, projectID string) []string { //nolint:staticcheck,revive // preserving deprecated method name
	return ig.GetPublicIPs(t, projectID)
}

// GetPublicIpsE returns a slice of the public IPs from the given Zonal Instance Group.
//
// Deprecated: Use [ZonalInstanceGroup.GetPublicIPsE] instead.
func (ig *ZonalInstanceGroup) GetPublicIpsE(t testing.TestingT, projectID string) ([]string, error) { //nolint:staticcheck,revive // preserving deprecated method name
	return ig.GetPublicIPsE(t, projectID)
}

// GetPublicIPsContextE returns a slice of the public IPs from the given Zonal Instance Group.
// The ctx parameter supports cancellation and timeouts.
func (ig *ZonalInstanceGroup) GetPublicIPsContextE(t testing.TestingT, ctx context.Context, projectID string) ([]string, error) {
	return getPublicIPsContextE(t, ctx, ig, projectID)
}

// GetPublicIPs returns a slice of the public IPs from the given Regional Instance Group.
// This will fail the test if there is an error.
//
// Deprecated: Use [RegionalInstanceGroup.GetPublicIPsContext] instead.
func (ig *RegionalInstanceGroup) GetPublicIPs(t testing.TestingT, projectID string) []string {
	return ig.GetPublicIPsContext(t, context.Background(), projectID)
}

// GetPublicIPsContext returns a slice of the public IPs from the given Regional Instance Group.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func (ig *RegionalInstanceGroup) GetPublicIPsContext(t testing.TestingT, ctx context.Context, projectID string) []string {
	ips, err := ig.GetPublicIPsContextE(t, ctx, projectID)
	require.NoError(t, err)

	return ips
}

// GetPublicIPsE returns a slice of the public IPs from the given Regional Instance Group.
//
// Deprecated: Use [RegionalInstanceGroup.GetPublicIPsContextE] instead.
func (ig *RegionalInstanceGroup) GetPublicIPsE(t testing.TestingT, projectID string) ([]string, error) {
	return ig.GetPublicIPsContextE(t, context.Background(), projectID)
}

// GetPublicIps returns a slice of the public IPs from the given Regional Instance Group.
// This will fail the test if there is an error.
//
// Deprecated: Use [RegionalInstanceGroup.GetPublicIPs] instead.
func (ig *RegionalInstanceGroup) GetPublicIps(t testing.TestingT, projectID string) []string { //nolint:staticcheck,revive // preserving deprecated method name
	return ig.GetPublicIPs(t, projectID)
}

// GetPublicIpsE returns a slice of the public IPs from the given Regional Instance Group.
//
// Deprecated: Use [RegionalInstanceGroup.GetPublicIPsE] instead.
func (ig *RegionalInstanceGroup) GetPublicIpsE(t testing.TestingT, projectID string) ([]string, error) { //nolint:staticcheck,revive // preserving deprecated method name
	return ig.GetPublicIPsE(t, projectID)
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

// GetRandomInstance returns a randomly selected Instance from the Zonal Instance Group.
// This will fail the test if there is an error.
//
// Deprecated: Use [ZonalInstanceGroup.GetRandomInstanceContext] instead.
func (ig *ZonalInstanceGroup) GetRandomInstance(t testing.TestingT) *Instance {
	return ig.GetRandomInstanceContext(t, context.Background())
}

// GetRandomInstanceContext returns a randomly selected Instance from the Zonal Instance Group.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func (ig *ZonalInstanceGroup) GetRandomInstanceContext(t testing.TestingT, ctx context.Context) *Instance {
	instance, err := ig.GetRandomInstanceContextE(t, ctx)
	require.NoError(t, err)

	return instance
}

// GetRandomInstanceE returns a randomly selected Instance from the Zonal Instance Group.
//
// Deprecated: Use [ZonalInstanceGroup.GetRandomInstanceContextE] instead.
func (ig *ZonalInstanceGroup) GetRandomInstanceE(t testing.TestingT) (*Instance, error) {
	return ig.GetRandomInstanceContextE(t, context.Background())
}

// GetRandomInstanceContextE returns a randomly selected Instance from the Zonal Instance Group.
// The ctx parameter supports cancellation and timeouts.
func (ig *ZonalInstanceGroup) GetRandomInstanceContextE(t testing.TestingT, ctx context.Context) (*Instance, error) {
	return getRandomInstanceContextE(t, ctx, ig, ig.Name, ig.Region, ig.Size, ig.projectID)
}

// GetRandomInstance returns a randomly selected Instance from the Regional Instance Group.
// This will fail the test if there is an error.
//
// Deprecated: Use [RegionalInstanceGroup.GetRandomInstanceContext] instead.
func (ig *RegionalInstanceGroup) GetRandomInstance(t testing.TestingT) *Instance {
	return ig.GetRandomInstanceContext(t, context.Background())
}

// GetRandomInstanceContext returns a randomly selected Instance from the Regional Instance Group.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func (ig *RegionalInstanceGroup) GetRandomInstanceContext(t testing.TestingT, ctx context.Context) *Instance {
	instance, err := ig.GetRandomInstanceContextE(t, ctx)
	require.NoError(t, err)

	return instance
}

// GetRandomInstanceE returns a randomly selected Instance from the Regional Instance Group.
//
// Deprecated: Use [RegionalInstanceGroup.GetRandomInstanceContextE] instead.
func (ig *RegionalInstanceGroup) GetRandomInstanceE(t testing.TestingT) (*Instance, error) {
	return ig.GetRandomInstanceContextE(t, context.Background())
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

// NewComputeService creates a new Compute service, which is used to make GCE API calls.
// This will fail the test if there is an error.
//
// Deprecated: Use [NewComputeServiceContext] instead.
func NewComputeService(t testing.TestingT) *compute.Service {
	return NewComputeServiceContext(t, context.Background())
}

// NewComputeServiceContext creates a new Compute service, which is used to make GCE API calls.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func NewComputeServiceContext(t testing.TestingT, ctx context.Context) *compute.Service {
	client, err := NewComputeServiceContextE(t, ctx)
	require.NoError(t, err)

	return client
}

// NewComputeServiceE creates a new Compute service, which is used to make GCE API calls.
//
// Deprecated: Use [NewComputeServiceContextE] instead.
func NewComputeServiceE(t testing.TestingT) (*compute.Service, error) {
	return NewComputeServiceContextE(t, context.Background())
}

// NewComputeServiceContextE creates a new Compute service, which is used to make GCE API calls.
// The ctx parameter supports cancellation and timeouts.
func NewComputeServiceContextE(t testing.TestingT, ctx context.Context) (*compute.Service, error) {
	if ts, ok := getStaticTokenSource(); ok {
		return compute.NewService(ctx, option.WithTokenSource(ts))
	}

	// Retrieve the Google OAuth token using a retry loop as it can sometimes return an error.
	// e.g: oauth2: cannot fetch token: Post https://oauth2.googleapis.com/token: net/http: TLS handshake timeout
	// This is loosely based on https://github.com/kubernetes/kubernetes/blob/7e8de5422cb5ad76dd0c147cf4336220d282e34b/pkg/cloudprovider/providers/gce/gce.go#L831.

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

// NewInstancesService creates a new InstancesService, which is used to make a subset of GCE API calls.
// This will fail the test if there is an error.
//
// Deprecated: Use [NewInstancesServiceContext] instead.
func NewInstancesService(t testing.TestingT) *compute.InstancesService {
	return NewInstancesServiceContext(t, context.Background())
}

// NewInstancesServiceContext creates a new InstancesService, which is used to make a subset of GCE API calls.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func NewInstancesServiceContext(t testing.TestingT, ctx context.Context) *compute.InstancesService {
	client, err := NewInstancesServiceContextE(t, ctx)
	require.NoError(t, err)

	return client
}

// NewInstancesServiceE creates a new InstancesService, which is used to make a subset of GCE API calls.
//
// Deprecated: Use [NewInstancesServiceContextE] instead.
func NewInstancesServiceE(t testing.TestingT) (*compute.InstancesService, error) {
	return NewInstancesServiceContextE(t, context.Background())
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

// RandomValidGcpName returns a random, valid name for GCP resources. Many resources in GCP require lowercase letters only.
//
// Deprecated: Use [RandomValidGCPName] instead.
func RandomValidGcpName() string { //nolint:staticcheck,revive // preserving deprecated function name
	return RandomValidGCPName()
}
