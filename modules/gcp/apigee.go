package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/apigee/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

// Apigee names everything under an organization rather than under a project, and an organization is a
// separate thing a project is attached to, even though its id is usually the project id. Every read here
// therefore takes an organization id where reads for other APIs take a project id.

// GetApigeeOrganizationAttrs returns the settings Google Cloud holds for the given Apigee organization, so a test can assert on what was
// actually created rather than only that it exists.
// An organization is the tenant everything else belongs to, so its runtime type, its analytics region and the network it peers with are set once and shape everything inside it.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetApigeeOrganizationAttrs(t testing.TestingT, ctx context.Context, orgID string) *apigee.GoogleCloudApigeeV1Organization {
	attrs, err := GetApigeeOrganizationAttrsE(t, ctx, orgID)
	require.NoError(t, err)

	return attrs
}

// GetApigeeOrganizationAttrsE returns the settings Google Cloud holds for the given Apigee organization.
// The ctx parameter supports cancellation and timeouts.
func GetApigeeOrganizationAttrsE(t testing.TestingT, ctx context.Context, orgID string) (*apigee.GoogleCloudApigeeV1Organization, error) {
	logger.Default.Logf(t, "Getting settings for Apigee organization %s", orgID)

	service, err := NewApigeeServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetApigeeOrganizationAttrsWithClient(ctx, service, orgID)
}

// GetApigeeOrganizationAttrsWithClient returns the settings Google Cloud holds for the given Apigee organization using the supplied
// *apigee.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see apigee_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetApigeeOrganizationAttrsWithClient(ctx context.Context, service *apigee.Service, orgID string) (*apigee.GoogleCloudApigeeV1Organization, error) {
	name := "organizations/" + orgID

	attrs, err := service.Organizations.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Apigee organization %s does not exist", orgID)
		}

		return nil, fmt.Errorf("failed to get settings for Apigee organization %s: %w", orgID, err)
	}

	return attrs, nil
}

// GetApigeeAPIProxyAttrs returns the settings Google Cloud holds for the given Apigee API proxy, so a test can assert on what was
// actually created rather than only that it exists.
// A proxy is the API itself, so which revisions exist under it is what says whether anything can be deployed.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetApigeeAPIProxyAttrs(t testing.TestingT, ctx context.Context, orgID string, apiID string) *apigee.GoogleCloudApigeeV1ApiProxy {
	attrs, err := GetApigeeAPIProxyAttrsE(t, ctx, orgID, apiID)
	require.NoError(t, err)

	return attrs
}

// GetApigeeAPIProxyAttrsE returns the settings Google Cloud holds for the given Apigee API proxy.
// The ctx parameter supports cancellation and timeouts.
func GetApigeeAPIProxyAttrsE(t testing.TestingT, ctx context.Context, orgID string, apiID string) (*apigee.GoogleCloudApigeeV1ApiProxy, error) {
	logger.Default.Logf(t, "Getting settings for Apigee API proxy %s in organization %s", apiID, orgID)

	service, err := NewApigeeServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetApigeeAPIProxyAttrsWithClient(ctx, service, orgID, apiID)
}

// GetApigeeAPIProxyAttrsWithClient returns the settings Google Cloud holds for the given Apigee API proxy using the supplied
// *apigee.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see apigee_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetApigeeAPIProxyAttrsWithClient(ctx context.Context, service *apigee.Service, orgID string, apiID string) (*apigee.GoogleCloudApigeeV1ApiProxy, error) {
	name := fmt.Sprintf("organizations/%s/apis/%s", orgID, apiID)

	attrs, err := service.Organizations.Apis.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Apigee API proxy %s in organization %s does not exist", apiID, orgID)
		}

		return nil, fmt.Errorf("failed to get settings for Apigee API proxy %s in organization %s: %w", apiID, orgID, err)
	}

	return attrs, nil
}

// GetApigeeAPIProductAttrs returns the settings Google Cloud holds for the given Apigee API product, so a test can assert on what was
// actually created rather than only that it exists.
// A product is the bundle a developer app is given access to, so its approval type and the proxies it names decide what any app key can call.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetApigeeAPIProductAttrs(t testing.TestingT, ctx context.Context, orgID string, productID string) *apigee.GoogleCloudApigeeV1ApiProduct {
	attrs, err := GetApigeeAPIProductAttrsE(t, ctx, orgID, productID)
	require.NoError(t, err)

	return attrs
}

// GetApigeeAPIProductAttrsE returns the settings Google Cloud holds for the given Apigee API product.
// The ctx parameter supports cancellation and timeouts.
func GetApigeeAPIProductAttrsE(t testing.TestingT, ctx context.Context, orgID string, productID string) (*apigee.GoogleCloudApigeeV1ApiProduct, error) {
	logger.Default.Logf(t, "Getting settings for Apigee API product %s in organization %s", productID, orgID)

	service, err := NewApigeeServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetApigeeAPIProductAttrsWithClient(ctx, service, orgID, productID)
}

// GetApigeeAPIProductAttrsWithClient returns the settings Google Cloud holds for the given Apigee API product using the supplied
// *apigee.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see apigee_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetApigeeAPIProductAttrsWithClient(ctx context.Context, service *apigee.Service, orgID string, productID string) (*apigee.GoogleCloudApigeeV1ApiProduct, error) {
	name := fmt.Sprintf("organizations/%s/apiproducts/%s", orgID, productID)

	attrs, err := service.Organizations.Apiproducts.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Apigee API product %s in organization %s does not exist", productID, orgID)
		}

		return nil, fmt.Errorf("failed to get settings for Apigee API product %s in organization %s: %w", productID, orgID, err)
	}

	return attrs, nil
}

// GetApigeeAppGroupAttrs returns the settings Google Cloud holds for the given Apigee app group, so a test can assert on what was
// actually created rather than only that it exists.
// An app group owns apps on behalf of a team rather than a person, so its channel and status decide whose apps these are.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetApigeeAppGroupAttrs(t testing.TestingT, ctx context.Context, orgID string, groupID string) *apigee.GoogleCloudApigeeV1AppGroup {
	attrs, err := GetApigeeAppGroupAttrsE(t, ctx, orgID, groupID)
	require.NoError(t, err)

	return attrs
}

// GetApigeeAppGroupAttrsE returns the settings Google Cloud holds for the given Apigee app group.
// The ctx parameter supports cancellation and timeouts.
func GetApigeeAppGroupAttrsE(t testing.TestingT, ctx context.Context, orgID string, groupID string) (*apigee.GoogleCloudApigeeV1AppGroup, error) {
	logger.Default.Logf(t, "Getting settings for Apigee app group %s in organization %s", groupID, orgID)

	service, err := NewApigeeServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetApigeeAppGroupAttrsWithClient(ctx, service, orgID, groupID)
}

// GetApigeeAppGroupAttrsWithClient returns the settings Google Cloud holds for the given Apigee app group using the supplied
// *apigee.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see apigee_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetApigeeAppGroupAttrsWithClient(ctx context.Context, service *apigee.Service, orgID string, groupID string) (*apigee.GoogleCloudApigeeV1AppGroup, error) {
	name := fmt.Sprintf("organizations/%s/appgroups/%s", orgID, groupID)

	attrs, err := service.Organizations.Appgroups.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Apigee app group %s in organization %s does not exist", groupID, orgID)
		}

		return nil, fmt.Errorf("failed to get settings for Apigee app group %s in organization %s: %w", groupID, orgID, err)
	}

	return attrs, nil
}

// GetApigeeDeveloperAttrs returns the settings Google Cloud holds for the given Apigee developer, so a test can assert on what was
// actually created rather than only that it exists.
// A developer is who owns an app, so the email that identifies them and their status decide whether their keys work.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetApigeeDeveloperAttrs(t testing.TestingT, ctx context.Context, orgID string, email string) *apigee.GoogleCloudApigeeV1Developer {
	attrs, err := GetApigeeDeveloperAttrsE(t, ctx, orgID, email)
	require.NoError(t, err)

	return attrs
}

// GetApigeeDeveloperAttrsE returns the settings Google Cloud holds for the given Apigee developer.
// The ctx parameter supports cancellation and timeouts.
func GetApigeeDeveloperAttrsE(t testing.TestingT, ctx context.Context, orgID string, email string) (*apigee.GoogleCloudApigeeV1Developer, error) {
	logger.Default.Logf(t, "Getting settings for Apigee developer %s in organization %s", email, orgID)

	service, err := NewApigeeServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetApigeeDeveloperAttrsWithClient(ctx, service, orgID, email)
}

// GetApigeeDeveloperAttrsWithClient returns the settings Google Cloud holds for the given Apigee developer using the supplied
// *apigee.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see apigee_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetApigeeDeveloperAttrsWithClient(ctx context.Context, service *apigee.Service, orgID string, email string) (*apigee.GoogleCloudApigeeV1Developer, error) {
	name := fmt.Sprintf("organizations/%s/developers/%s", orgID, email)

	attrs, err := service.Organizations.Developers.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Apigee developer %s in organization %s does not exist", email, orgID)
		}

		return nil, fmt.Errorf("failed to get settings for Apigee developer %s in organization %s: %w", email, orgID, err)
	}

	return attrs, nil
}

// GetApigeeDeveloperAppAttrs returns the settings Google Cloud holds for the given Apigee developer app, so a test can assert on what was
// actually created rather than only that it exists.
// An app is what holds the keys a caller uses, so the products it is bound to decide what those keys may reach.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetApigeeDeveloperAppAttrs(t testing.TestingT, ctx context.Context, orgID string, email string, appID string) *apigee.GoogleCloudApigeeV1DeveloperApp {
	attrs, err := GetApigeeDeveloperAppAttrsE(t, ctx, orgID, email, appID)
	require.NoError(t, err)

	return attrs
}

// GetApigeeDeveloperAppAttrsE returns the settings Google Cloud holds for the given Apigee developer app.
// The ctx parameter supports cancellation and timeouts.
func GetApigeeDeveloperAppAttrsE(t testing.TestingT, ctx context.Context, orgID string, email string, appID string) (*apigee.GoogleCloudApigeeV1DeveloperApp, error) {
	logger.Default.Logf(t, "Getting settings for Apigee developer app %s %s in organization %s", appID, email, orgID)

	service, err := NewApigeeServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetApigeeDeveloperAppAttrsWithClient(ctx, service, orgID, email, appID)
}

// GetApigeeDeveloperAppAttrsWithClient returns the settings Google Cloud holds for the given Apigee developer app using the supplied
// *apigee.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see apigee_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetApigeeDeveloperAppAttrsWithClient(ctx context.Context, service *apigee.Service, orgID string, email string, appID string) (*apigee.GoogleCloudApigeeV1DeveloperApp, error) {
	name := fmt.Sprintf("organizations/%s/developers/%s/apps/%s", orgID, email, appID)

	attrs, err := service.Organizations.Developers.Apps.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Apigee developer app %s %s in organization %s does not exist", appID, email, orgID)
		}

		return nil, fmt.Errorf("failed to get settings for Apigee developer app %s %s in organization %s: %w", appID, email, orgID, err)
	}

	return attrs, nil
}

// GetApigeeEnvironmentAttrs returns the settings Google Cloud holds for the given Apigee environment, so a test can assert on what was
// actually created rather than only that it exists.
// An environment is where a proxy actually runs, so its deployment type and its node config decide how it serves traffic.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetApigeeEnvironmentAttrs(t testing.TestingT, ctx context.Context, orgID string, envID string) *apigee.GoogleCloudApigeeV1Environment {
	attrs, err := GetApigeeEnvironmentAttrsE(t, ctx, orgID, envID)
	require.NoError(t, err)

	return attrs
}

// GetApigeeEnvironmentAttrsE returns the settings Google Cloud holds for the given Apigee environment.
// The ctx parameter supports cancellation and timeouts.
func GetApigeeEnvironmentAttrsE(t testing.TestingT, ctx context.Context, orgID string, envID string) (*apigee.GoogleCloudApigeeV1Environment, error) {
	logger.Default.Logf(t, "Getting settings for Apigee environment %s in organization %s", envID, orgID)

	service, err := NewApigeeServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetApigeeEnvironmentAttrsWithClient(ctx, service, orgID, envID)
}

// GetApigeeEnvironmentAttrsWithClient returns the settings Google Cloud holds for the given Apigee environment using the supplied
// *apigee.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see apigee_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetApigeeEnvironmentAttrsWithClient(ctx context.Context, service *apigee.Service, orgID string, envID string) (*apigee.GoogleCloudApigeeV1Environment, error) {
	name := fmt.Sprintf("organizations/%s/environments/%s", orgID, envID)

	attrs, err := service.Organizations.Environments.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Apigee environment %s in organization %s does not exist", envID, orgID)
		}

		return nil, fmt.Errorf("failed to get settings for Apigee environment %s in organization %s: %w", envID, orgID, err)
	}

	return attrs, nil
}

// GetApigeeEnvironmentGroupAttrs returns the settings Google Cloud holds for the given Apigee environment group, so a test can assert on what was
// actually created rather than only that it exists.
// A group is the set of hostnames traffic arrives on, so the hostnames it holds decide which requests reach the environments attached to it.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetApigeeEnvironmentGroupAttrs(t testing.TestingT, ctx context.Context, orgID string, groupID string) *apigee.GoogleCloudApigeeV1EnvironmentGroup {
	attrs, err := GetApigeeEnvironmentGroupAttrsE(t, ctx, orgID, groupID)
	require.NoError(t, err)

	return attrs
}

// GetApigeeEnvironmentGroupAttrsE returns the settings Google Cloud holds for the given Apigee environment group.
// The ctx parameter supports cancellation and timeouts.
func GetApigeeEnvironmentGroupAttrsE(t testing.TestingT, ctx context.Context, orgID string, groupID string) (*apigee.GoogleCloudApigeeV1EnvironmentGroup, error) {
	logger.Default.Logf(t, "Getting settings for Apigee environment group %s in organization %s", groupID, orgID)

	service, err := NewApigeeServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetApigeeEnvironmentGroupAttrsWithClient(ctx, service, orgID, groupID)
}

// GetApigeeEnvironmentGroupAttrsWithClient returns the settings Google Cloud holds for the given Apigee environment group using the supplied
// *apigee.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see apigee_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetApigeeEnvironmentGroupAttrsWithClient(ctx context.Context, service *apigee.Service, orgID string, groupID string) (*apigee.GoogleCloudApigeeV1EnvironmentGroup, error) {
	name := fmt.Sprintf("organizations/%s/envgroups/%s", orgID, groupID)

	attrs, err := service.Organizations.Envgroups.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Apigee environment group %s in organization %s does not exist", groupID, orgID)
		}

		return nil, fmt.Errorf("failed to get settings for Apigee environment group %s in organization %s: %w", groupID, orgID, err)
	}

	return attrs, nil
}

// GetApigeeEnvironmentGroupAttachmentAttrs returns the settings Google Cloud holds for the given Apigee environment group attachment, so a test can assert on what was
// actually created rather than only that it exists.
// The attachment is what puts one environment behind a group's hostnames, so without it the group routes nowhere.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetApigeeEnvironmentGroupAttachmentAttrs(t testing.TestingT, ctx context.Context, orgID string, groupID string, attachmentID string) *apigee.GoogleCloudApigeeV1EnvironmentGroupAttachment {
	attrs, err := GetApigeeEnvironmentGroupAttachmentAttrsE(t, ctx, orgID, groupID, attachmentID)
	require.NoError(t, err)

	return attrs
}

// GetApigeeEnvironmentGroupAttachmentAttrsE returns the settings Google Cloud holds for the given Apigee environment group attachment.
// The ctx parameter supports cancellation and timeouts.
func GetApigeeEnvironmentGroupAttachmentAttrsE(t testing.TestingT, ctx context.Context, orgID string, groupID string, attachmentID string) (*apigee.GoogleCloudApigeeV1EnvironmentGroupAttachment, error) {
	logger.Default.Logf(t, "Getting settings for Apigee environment group attachment %s %s in organization %s", attachmentID, groupID, orgID)

	service, err := NewApigeeServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetApigeeEnvironmentGroupAttachmentAttrsWithClient(ctx, service, orgID, groupID, attachmentID)
}

// GetApigeeEnvironmentGroupAttachmentAttrsWithClient returns the settings Google Cloud holds for the given Apigee environment group attachment using the supplied
// *apigee.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see apigee_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetApigeeEnvironmentGroupAttachmentAttrsWithClient(ctx context.Context, service *apigee.Service, orgID string, groupID string, attachmentID string) (*apigee.GoogleCloudApigeeV1EnvironmentGroupAttachment, error) {
	name := fmt.Sprintf("organizations/%s/envgroups/%s/attachments/%s", orgID, groupID, attachmentID)

	attrs, err := service.Organizations.Envgroups.Attachments.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Apigee environment group attachment %s %s in organization %s does not exist", attachmentID, groupID, orgID)
		}

		return nil, fmt.Errorf("failed to get settings for Apigee environment group attachment %s %s in organization %s: %w", attachmentID, groupID, orgID, err)
	}

	return attrs, nil
}

// GetApigeeInstanceAttrs returns the settings Google Cloud holds for the given Apigee instance, so a test can assert on what was
// actually created rather than only that it exists.
// An instance is the runtime in one region, so its location, its IP range and its disk encryption key are what it is made of.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetApigeeInstanceAttrs(t testing.TestingT, ctx context.Context, orgID string, instanceID string) *apigee.GoogleCloudApigeeV1Instance {
	attrs, err := GetApigeeInstanceAttrsE(t, ctx, orgID, instanceID)
	require.NoError(t, err)

	return attrs
}

// GetApigeeInstanceAttrsE returns the settings Google Cloud holds for the given Apigee instance.
// The ctx parameter supports cancellation and timeouts.
func GetApigeeInstanceAttrsE(t testing.TestingT, ctx context.Context, orgID string, instanceID string) (*apigee.GoogleCloudApigeeV1Instance, error) {
	logger.Default.Logf(t, "Getting settings for Apigee instance %s in organization %s", instanceID, orgID)

	service, err := NewApigeeServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetApigeeInstanceAttrsWithClient(ctx, service, orgID, instanceID)
}

// GetApigeeInstanceAttrsWithClient returns the settings Google Cloud holds for the given Apigee instance using the supplied
// *apigee.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see apigee_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetApigeeInstanceAttrsWithClient(ctx context.Context, service *apigee.Service, orgID string, instanceID string) (*apigee.GoogleCloudApigeeV1Instance, error) {
	name := fmt.Sprintf("organizations/%s/instances/%s", orgID, instanceID)

	attrs, err := service.Organizations.Instances.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Apigee instance %s in organization %s does not exist", instanceID, orgID)
		}

		return nil, fmt.Errorf("failed to get settings for Apigee instance %s in organization %s: %w", instanceID, orgID, err)
	}

	return attrs, nil
}

// GetApigeeInstanceAttachmentAttrs returns the settings Google Cloud holds for the given Apigee instance attachment, so a test can assert on what was
// actually created rather than only that it exists.
// The attachment is what makes an environment run on one instance, so the environment it names is the whole point.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetApigeeInstanceAttachmentAttrs(t testing.TestingT, ctx context.Context, orgID string, instanceID string, attachmentID string) *apigee.GoogleCloudApigeeV1InstanceAttachment {
	attrs, err := GetApigeeInstanceAttachmentAttrsE(t, ctx, orgID, instanceID, attachmentID)
	require.NoError(t, err)

	return attrs
}

// GetApigeeInstanceAttachmentAttrsE returns the settings Google Cloud holds for the given Apigee instance attachment.
// The ctx parameter supports cancellation and timeouts.
func GetApigeeInstanceAttachmentAttrsE(t testing.TestingT, ctx context.Context, orgID string, instanceID string, attachmentID string) (*apigee.GoogleCloudApigeeV1InstanceAttachment, error) {
	logger.Default.Logf(t, "Getting settings for Apigee instance attachment %s %s in organization %s", attachmentID, instanceID, orgID)

	service, err := NewApigeeServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetApigeeInstanceAttachmentAttrsWithClient(ctx, service, orgID, instanceID, attachmentID)
}

// GetApigeeInstanceAttachmentAttrsWithClient returns the settings Google Cloud holds for the given Apigee instance attachment using the supplied
// *apigee.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see apigee_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetApigeeInstanceAttachmentAttrsWithClient(ctx context.Context, service *apigee.Service, orgID string, instanceID string, attachmentID string) (*apigee.GoogleCloudApigeeV1InstanceAttachment, error) {
	name := fmt.Sprintf("organizations/%s/instances/%s/attachments/%s", orgID, instanceID, attachmentID)

	attrs, err := service.Organizations.Instances.Attachments.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Apigee instance attachment %s %s in organization %s does not exist", attachmentID, instanceID, orgID)
		}

		return nil, fmt.Errorf("failed to get settings for Apigee instance attachment %s %s in organization %s: %w", attachmentID, instanceID, orgID, err)
	}

	return attrs, nil
}

// GetApigeeNatAddressAttrs returns the settings Google Cloud holds for the given Apigee NAT address, so a test can assert on what was
// actually created rather than only that it exists.
// A NAT address is the fixed outbound IP a backend can allow, so the address itself and its state are the point.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetApigeeNatAddressAttrs(t testing.TestingT, ctx context.Context, orgID string, instanceID string, addressID string) *apigee.GoogleCloudApigeeV1NatAddress {
	attrs, err := GetApigeeNatAddressAttrsE(t, ctx, orgID, instanceID, addressID)
	require.NoError(t, err)

	return attrs
}

// GetApigeeNatAddressAttrsE returns the settings Google Cloud holds for the given Apigee NAT address.
// The ctx parameter supports cancellation and timeouts.
func GetApigeeNatAddressAttrsE(t testing.TestingT, ctx context.Context, orgID string, instanceID string, addressID string) (*apigee.GoogleCloudApigeeV1NatAddress, error) {
	logger.Default.Logf(t, "Getting settings for Apigee NAT address %s %s in organization %s", addressID, instanceID, orgID)

	service, err := NewApigeeServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetApigeeNatAddressAttrsWithClient(ctx, service, orgID, instanceID, addressID)
}

// GetApigeeNatAddressAttrsWithClient returns the settings Google Cloud holds for the given Apigee NAT address using the supplied
// *apigee.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see apigee_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetApigeeNatAddressAttrsWithClient(ctx context.Context, service *apigee.Service, orgID string, instanceID string, addressID string) (*apigee.GoogleCloudApigeeV1NatAddress, error) {
	name := fmt.Sprintf("organizations/%s/instances/%s/natAddresses/%s", orgID, instanceID, addressID)

	attrs, err := service.Organizations.Instances.NatAddresses.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Apigee NAT address %s %s in organization %s does not exist", addressID, instanceID, orgID)
		}

		return nil, fmt.Errorf("failed to get settings for Apigee NAT address %s %s in organization %s: %w", addressID, instanceID, orgID, err)
	}

	return attrs, nil
}

// GetApigeeEndpointAttachmentAttrs returns the settings Google Cloud holds for the given Apigee endpoint attachment, so a test can assert on what was
// actually created rather than only that it exists.
// An endpoint attachment is how a proxy reaches a private service, so the service it targets and its state decide whether that call connects.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetApigeeEndpointAttachmentAttrs(t testing.TestingT, ctx context.Context, orgID string, attachmentID string) *apigee.GoogleCloudApigeeV1EndpointAttachment {
	attrs, err := GetApigeeEndpointAttachmentAttrsE(t, ctx, orgID, attachmentID)
	require.NoError(t, err)

	return attrs
}

// GetApigeeEndpointAttachmentAttrsE returns the settings Google Cloud holds for the given Apigee endpoint attachment.
// The ctx parameter supports cancellation and timeouts.
func GetApigeeEndpointAttachmentAttrsE(t testing.TestingT, ctx context.Context, orgID string, attachmentID string) (*apigee.GoogleCloudApigeeV1EndpointAttachment, error) {
	logger.Default.Logf(t, "Getting settings for Apigee endpoint attachment %s in organization %s", attachmentID, orgID)

	service, err := NewApigeeServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetApigeeEndpointAttachmentAttrsWithClient(ctx, service, orgID, attachmentID)
}

// GetApigeeEndpointAttachmentAttrsWithClient returns the settings Google Cloud holds for the given Apigee endpoint attachment using the supplied
// *apigee.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see apigee_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetApigeeEndpointAttachmentAttrsWithClient(ctx context.Context, service *apigee.Service, orgID string, attachmentID string) (*apigee.GoogleCloudApigeeV1EndpointAttachment, error) {
	name := fmt.Sprintf("organizations/%s/endpointAttachments/%s", orgID, attachmentID)

	attrs, err := service.Organizations.EndpointAttachments.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Apigee endpoint attachment %s in organization %s does not exist", attachmentID, orgID)
		}

		return nil, fmt.Errorf("failed to get settings for Apigee endpoint attachment %s in organization %s: %w", attachmentID, orgID, err)
	}

	return attrs, nil
}

// GetApigeeDNSZoneAttrs returns the settings Google Cloud holds for the given Apigee DNS zone, so a test can assert on what was
// actually created rather than only that it exists.
// A zone is how a proxy resolves a private domain, so the domain and the networks that may use it are the point.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetApigeeDNSZoneAttrs(t testing.TestingT, ctx context.Context, orgID string, zoneID string) *apigee.GoogleCloudApigeeV1DnsZone {
	attrs, err := GetApigeeDNSZoneAttrsE(t, ctx, orgID, zoneID)
	require.NoError(t, err)

	return attrs
}

// GetApigeeDNSZoneAttrsE returns the settings Google Cloud holds for the given Apigee DNS zone.
// The ctx parameter supports cancellation and timeouts.
func GetApigeeDNSZoneAttrsE(t testing.TestingT, ctx context.Context, orgID string, zoneID string) (*apigee.GoogleCloudApigeeV1DnsZone, error) {
	logger.Default.Logf(t, "Getting settings for Apigee DNS zone %s in organization %s", zoneID, orgID)

	service, err := NewApigeeServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetApigeeDNSZoneAttrsWithClient(ctx, service, orgID, zoneID)
}

// GetApigeeDNSZoneAttrsWithClient returns the settings Google Cloud holds for the given Apigee DNS zone using the supplied
// *apigee.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see apigee_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetApigeeDNSZoneAttrsWithClient(ctx context.Context, service *apigee.Service, orgID string, zoneID string) (*apigee.GoogleCloudApigeeV1DnsZone, error) {
	name := fmt.Sprintf("organizations/%s/dnsZones/%s", orgID, zoneID)

	attrs, err := service.Organizations.DnsZones.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Apigee DNS zone %s in organization %s does not exist", zoneID, orgID)
		}

		return nil, fmt.Errorf("failed to get settings for Apigee DNS zone %s in organization %s: %w", zoneID, orgID, err)
	}

	return attrs, nil
}

// GetApigeeDataCollectorAttrs returns the settings Google Cloud holds for the given Apigee data collector, so a test can assert on what was
// actually created rather than only that it exists.
// A collector declares one custom analytics field, so its type decides what a proxy may record into it.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetApigeeDataCollectorAttrs(t testing.TestingT, ctx context.Context, orgID string, collectorID string) *apigee.GoogleCloudApigeeV1DataCollector {
	attrs, err := GetApigeeDataCollectorAttrsE(t, ctx, orgID, collectorID)
	require.NoError(t, err)

	return attrs
}

// GetApigeeDataCollectorAttrsE returns the settings Google Cloud holds for the given Apigee data collector.
// The ctx parameter supports cancellation and timeouts.
func GetApigeeDataCollectorAttrsE(t testing.TestingT, ctx context.Context, orgID string, collectorID string) (*apigee.GoogleCloudApigeeV1DataCollector, error) {
	logger.Default.Logf(t, "Getting settings for Apigee data collector %s in organization %s", collectorID, orgID)

	service, err := NewApigeeServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetApigeeDataCollectorAttrsWithClient(ctx, service, orgID, collectorID)
}

// GetApigeeDataCollectorAttrsWithClient returns the settings Google Cloud holds for the given Apigee data collector using the supplied
// *apigee.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see apigee_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetApigeeDataCollectorAttrsWithClient(ctx context.Context, service *apigee.Service, orgID string, collectorID string) (*apigee.GoogleCloudApigeeV1DataCollector, error) {
	name := fmt.Sprintf("organizations/%s/datacollectors/%s", orgID, collectorID)

	attrs, err := service.Organizations.Datacollectors.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Apigee data collector %s in organization %s does not exist", collectorID, orgID)
		}

		return nil, fmt.Errorf("failed to get settings for Apigee data collector %s in organization %s: %w", collectorID, orgID, err)
	}

	return attrs, nil
}

// GetApigeeDatastoreAttrs returns the settings Google Cloud holds for the given Apigee datastore, so a test can assert on what was
// actually created rather than only that it exists.
// A datastore is where exported analytics are written, so its target type and its display name are what it is for.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetApigeeDatastoreAttrs(t testing.TestingT, ctx context.Context, orgID string, datastoreID string) *apigee.GoogleCloudApigeeV1Datastore {
	attrs, err := GetApigeeDatastoreAttrsE(t, ctx, orgID, datastoreID)
	require.NoError(t, err)

	return attrs
}

// GetApigeeDatastoreAttrsE returns the settings Google Cloud holds for the given Apigee datastore.
// The ctx parameter supports cancellation and timeouts.
func GetApigeeDatastoreAttrsE(t testing.TestingT, ctx context.Context, orgID string, datastoreID string) (*apigee.GoogleCloudApigeeV1Datastore, error) {
	logger.Default.Logf(t, "Getting settings for Apigee datastore %s in organization %s", datastoreID, orgID)

	service, err := NewApigeeServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetApigeeDatastoreAttrsWithClient(ctx, service, orgID, datastoreID)
}

// GetApigeeDatastoreAttrsWithClient returns the settings Google Cloud holds for the given Apigee datastore using the supplied
// *apigee.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see apigee_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetApigeeDatastoreAttrsWithClient(ctx context.Context, service *apigee.Service, orgID string, datastoreID string) (*apigee.GoogleCloudApigeeV1Datastore, error) {
	name := fmt.Sprintf("organizations/%s/analytics/datastores/%s", orgID, datastoreID)

	attrs, err := service.Organizations.Analytics.Datastores.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Apigee datastore %s in organization %s does not exist", datastoreID, orgID)
		}

		return nil, fmt.Errorf("failed to get settings for Apigee datastore %s in organization %s: %w", datastoreID, orgID, err)
	}

	return attrs, nil
}

// GetApigeeSharedFlowAttrs returns the settings Google Cloud holds for the given Apigee shared flow, so a test can assert on what was
// actually created rather than only that it exists.
// A shared flow is policy other proxies call into, so which revisions exist under it says whether it can be deployed.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetApigeeSharedFlowAttrs(t testing.TestingT, ctx context.Context, orgID string, flowID string) *apigee.GoogleCloudApigeeV1SharedFlow {
	attrs, err := GetApigeeSharedFlowAttrsE(t, ctx, orgID, flowID)
	require.NoError(t, err)

	return attrs
}

// GetApigeeSharedFlowAttrsE returns the settings Google Cloud holds for the given Apigee shared flow.
// The ctx parameter supports cancellation and timeouts.
func GetApigeeSharedFlowAttrsE(t testing.TestingT, ctx context.Context, orgID string, flowID string) (*apigee.GoogleCloudApigeeV1SharedFlow, error) {
	logger.Default.Logf(t, "Getting settings for Apigee shared flow %s in organization %s", flowID, orgID)

	service, err := NewApigeeServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetApigeeSharedFlowAttrsWithClient(ctx, service, orgID, flowID)
}

// GetApigeeSharedFlowAttrsWithClient returns the settings Google Cloud holds for the given Apigee shared flow using the supplied
// *apigee.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see apigee_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetApigeeSharedFlowAttrsWithClient(ctx context.Context, service *apigee.Service, orgID string, flowID string) (*apigee.GoogleCloudApigeeV1SharedFlow, error) {
	name := fmt.Sprintf("organizations/%s/sharedflows/%s", orgID, flowID)

	attrs, err := service.Organizations.Sharedflows.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Apigee shared flow %s in organization %s does not exist", flowID, orgID)
		}

		return nil, fmt.Errorf("failed to get settings for Apigee shared flow %s in organization %s: %w", flowID, orgID, err)
	}

	return attrs, nil
}

// GetApigeeSpaceAttrs returns the settings Google Cloud holds for the given Apigee space, so a test can assert on what was
// actually created rather than only that it exists.
// A space is how resources in one organization are split between teams, so its display name is how a caller tells spaces apart.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetApigeeSpaceAttrs(t testing.TestingT, ctx context.Context, orgID string, spaceID string) *apigee.GoogleCloudApigeeV1Space {
	attrs, err := GetApigeeSpaceAttrsE(t, ctx, orgID, spaceID)
	require.NoError(t, err)

	return attrs
}

// GetApigeeSpaceAttrsE returns the settings Google Cloud holds for the given Apigee space.
// The ctx parameter supports cancellation and timeouts.
func GetApigeeSpaceAttrsE(t testing.TestingT, ctx context.Context, orgID string, spaceID string) (*apigee.GoogleCloudApigeeV1Space, error) {
	logger.Default.Logf(t, "Getting settings for Apigee space %s in organization %s", spaceID, orgID)

	service, err := NewApigeeServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetApigeeSpaceAttrsWithClient(ctx, service, orgID, spaceID)
}

// GetApigeeSpaceAttrsWithClient returns the settings Google Cloud holds for the given Apigee space using the supplied
// *apigee.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see apigee_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetApigeeSpaceAttrsWithClient(ctx context.Context, service *apigee.Service, orgID string, spaceID string) (*apigee.GoogleCloudApigeeV1Space, error) {
	name := fmt.Sprintf("organizations/%s/spaces/%s", orgID, spaceID)

	attrs, err := service.Organizations.Spaces.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Apigee space %s in organization %s does not exist", spaceID, orgID)
		}

		return nil, fmt.Errorf("failed to get settings for Apigee space %s in organization %s: %w", spaceID, orgID, err)
	}

	return attrs, nil
}

// GetApigeeTargetServerAttrs returns the settings Google Cloud holds for the given Apigee target server, so a test can assert on what was
// actually created rather than only that it exists.
// A target server is the backend a proxy sends requests to, so its host, its port and whether it is enabled decide where traffic lands.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetApigeeTargetServerAttrs(t testing.TestingT, ctx context.Context, orgID string, envID string, serverID string) *apigee.GoogleCloudApigeeV1TargetServer {
	attrs, err := GetApigeeTargetServerAttrsE(t, ctx, orgID, envID, serverID)
	require.NoError(t, err)

	return attrs
}

// GetApigeeTargetServerAttrsE returns the settings Google Cloud holds for the given Apigee target server.
// The ctx parameter supports cancellation and timeouts.
func GetApigeeTargetServerAttrsE(t testing.TestingT, ctx context.Context, orgID string, envID string, serverID string) (*apigee.GoogleCloudApigeeV1TargetServer, error) {
	logger.Default.Logf(t, "Getting settings for Apigee target server %s %s in organization %s", serverID, envID, orgID)

	service, err := NewApigeeServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetApigeeTargetServerAttrsWithClient(ctx, service, orgID, envID, serverID)
}

// GetApigeeTargetServerAttrsWithClient returns the settings Google Cloud holds for the given Apigee target server using the supplied
// *apigee.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see apigee_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetApigeeTargetServerAttrsWithClient(ctx context.Context, service *apigee.Service, orgID string, envID string, serverID string) (*apigee.GoogleCloudApigeeV1TargetServer, error) {
	name := fmt.Sprintf("organizations/%s/environments/%s/targetservers/%s", orgID, envID, serverID)

	attrs, err := service.Organizations.Environments.Targetservers.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Apigee target server %s %s in organization %s does not exist", serverID, envID, orgID)
		}

		return nil, fmt.Errorf("failed to get settings for Apigee target server %s %s in organization %s: %w", serverID, envID, orgID, err)
	}

	return attrs, nil
}

// GetApigeeKeystoreAttrs returns the settings Google Cloud holds for the given Apigee keystore, so a test can assert on what was
// actually created rather than only that it exists.
// A keystore holds the certificates an environment presents or trusts, so which aliases it carries is what it offers.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetApigeeKeystoreAttrs(t testing.TestingT, ctx context.Context, orgID string, envID string, keystoreID string) *apigee.GoogleCloudApigeeV1Keystore {
	attrs, err := GetApigeeKeystoreAttrsE(t, ctx, orgID, envID, keystoreID)
	require.NoError(t, err)

	return attrs
}

// GetApigeeKeystoreAttrsE returns the settings Google Cloud holds for the given Apigee keystore.
// The ctx parameter supports cancellation and timeouts.
func GetApigeeKeystoreAttrsE(t testing.TestingT, ctx context.Context, orgID string, envID string, keystoreID string) (*apigee.GoogleCloudApigeeV1Keystore, error) {
	logger.Default.Logf(t, "Getting settings for Apigee keystore %s %s in organization %s", keystoreID, envID, orgID)

	service, err := NewApigeeServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetApigeeKeystoreAttrsWithClient(ctx, service, orgID, envID, keystoreID)
}

// GetApigeeKeystoreAttrsWithClient returns the settings Google Cloud holds for the given Apigee keystore using the supplied
// *apigee.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see apigee_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetApigeeKeystoreAttrsWithClient(ctx context.Context, service *apigee.Service, orgID string, envID string, keystoreID string) (*apigee.GoogleCloudApigeeV1Keystore, error) {
	name := fmt.Sprintf("organizations/%s/environments/%s/keystores/%s", orgID, envID, keystoreID)

	attrs, err := service.Organizations.Environments.Keystores.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Apigee keystore %s %s in organization %s does not exist", keystoreID, envID, orgID)
		}

		return nil, fmt.Errorf("failed to get settings for Apigee keystore %s %s in organization %s: %w", keystoreID, envID, orgID, err)
	}

	return attrs, nil
}

// GetApigeeKeystoreAliasAttrs returns the settings Google Cloud holds for the given Apigee keystore alias, so a test can assert on what was
// actually created rather than only that it exists.
// An alias is one certificate or key pair inside a keystore, so its type and the certificate it holds are the point. One read serves all three ways an alias can be created.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetApigeeKeystoreAliasAttrs(t testing.TestingT, ctx context.Context, orgID string, envID string, keystoreID string, aliasID string) *apigee.GoogleCloudApigeeV1Alias {
	attrs, err := GetApigeeKeystoreAliasAttrsE(t, ctx, orgID, envID, keystoreID, aliasID)
	require.NoError(t, err)

	return attrs
}

// GetApigeeKeystoreAliasAttrsE returns the settings Google Cloud holds for the given Apigee keystore alias.
// The ctx parameter supports cancellation and timeouts.
func GetApigeeKeystoreAliasAttrsE(t testing.TestingT, ctx context.Context, orgID string, envID string, keystoreID string, aliasID string) (*apigee.GoogleCloudApigeeV1Alias, error) {
	logger.Default.Logf(t, "Getting settings for Apigee keystore alias %s %s %s in organization %s", aliasID, keystoreID, envID, orgID)

	service, err := NewApigeeServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetApigeeKeystoreAliasAttrsWithClient(ctx, service, orgID, envID, keystoreID, aliasID)
}

// GetApigeeKeystoreAliasAttrsWithClient returns the settings Google Cloud holds for the given Apigee keystore alias using the supplied
// *apigee.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see apigee_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetApigeeKeystoreAliasAttrsWithClient(ctx context.Context, service *apigee.Service, orgID string, envID string, keystoreID string, aliasID string) (*apigee.GoogleCloudApigeeV1Alias, error) {
	name := fmt.Sprintf("organizations/%s/environments/%s/keystores/%s/aliases/%s", orgID, envID, keystoreID, aliasID)

	attrs, err := service.Organizations.Environments.Keystores.Aliases.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Apigee keystore alias %s %s %s in organization %s does not exist", aliasID, keystoreID, envID, orgID)
		}

		return nil, fmt.Errorf("failed to get settings for Apigee keystore alias %s %s %s in organization %s: %w", aliasID, keystoreID, envID, orgID, err)
	}

	return attrs, nil
}

// GetApigeeReferenceAttrs returns the settings Google Cloud holds for the given Apigee reference, so a test can assert on what was
// actually created rather than only that it exists.
// A reference lets a proxy name a keystore indirectly so it can be swapped without redeploying, so what it points at is the whole point.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetApigeeReferenceAttrs(t testing.TestingT, ctx context.Context, orgID string, envID string, referenceID string) *apigee.GoogleCloudApigeeV1Reference {
	attrs, err := GetApigeeReferenceAttrsE(t, ctx, orgID, envID, referenceID)
	require.NoError(t, err)

	return attrs
}

// GetApigeeReferenceAttrsE returns the settings Google Cloud holds for the given Apigee reference.
// The ctx parameter supports cancellation and timeouts.
func GetApigeeReferenceAttrsE(t testing.TestingT, ctx context.Context, orgID string, envID string, referenceID string) (*apigee.GoogleCloudApigeeV1Reference, error) {
	logger.Default.Logf(t, "Getting settings for Apigee reference %s %s in organization %s", referenceID, envID, orgID)

	service, err := NewApigeeServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetApigeeReferenceAttrsWithClient(ctx, service, orgID, envID, referenceID)
}

// GetApigeeReferenceAttrsWithClient returns the settings Google Cloud holds for the given Apigee reference using the supplied
// *apigee.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see apigee_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetApigeeReferenceAttrsWithClient(ctx context.Context, service *apigee.Service, orgID string, envID string, referenceID string) (*apigee.GoogleCloudApigeeV1Reference, error) {
	name := fmt.Sprintf("organizations/%s/environments/%s/references/%s", orgID, envID, referenceID)

	attrs, err := service.Organizations.Environments.References.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Apigee reference %s %s in organization %s does not exist", referenceID, envID, orgID)
		}

		return nil, fmt.Errorf("failed to get settings for Apigee reference %s %s in organization %s: %w", referenceID, envID, orgID, err)
	}

	return attrs, nil
}

// GetApigeeFlowHookAttrs returns the settings Google Cloud holds for the given Apigee flow hook, so a test can assert on what was
// actually created rather than only that it exists.
// A flow hook attaches a shared flow to every proxy in an environment, so the flow it names runs on traffic no proxy asked for.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetApigeeFlowHookAttrs(t testing.TestingT, ctx context.Context, orgID string, envID string, hookPoint string) *apigee.GoogleCloudApigeeV1FlowHook {
	attrs, err := GetApigeeFlowHookAttrsE(t, ctx, orgID, envID, hookPoint)
	require.NoError(t, err)

	return attrs
}

// GetApigeeFlowHookAttrsE returns the settings Google Cloud holds for the given Apigee flow hook.
// The ctx parameter supports cancellation and timeouts.
func GetApigeeFlowHookAttrsE(t testing.TestingT, ctx context.Context, orgID string, envID string, hookPoint string) (*apigee.GoogleCloudApigeeV1FlowHook, error) {
	logger.Default.Logf(t, "Getting settings for Apigee flow hook %s %s in organization %s", hookPoint, envID, orgID)

	service, err := NewApigeeServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetApigeeFlowHookAttrsWithClient(ctx, service, orgID, envID, hookPoint)
}

// GetApigeeFlowHookAttrsWithClient returns the settings Google Cloud holds for the given Apigee flow hook using the supplied
// *apigee.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see apigee_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetApigeeFlowHookAttrsWithClient(ctx context.Context, service *apigee.Service, orgID string, envID string, hookPoint string) (*apigee.GoogleCloudApigeeV1FlowHook, error) {
	name := fmt.Sprintf("organizations/%s/environments/%s/flowhooks/%s", orgID, envID, hookPoint)

	attrs, err := service.Organizations.Environments.Flowhooks.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Apigee flow hook %s %s in organization %s does not exist", hookPoint, envID, orgID)
		}

		return nil, fmt.Errorf("failed to get settings for Apigee flow hook %s %s in organization %s: %w", hookPoint, envID, orgID, err)
	}

	return attrs, nil
}

// GetApigeeEnvironmentKeyValueMapAttrs returns the settings Google Cloud holds for the given Apigee environment key value map, so a test can assert on what was
// actually created rather than only that it exists.
// A map is where a proxy keeps configuration it reads at runtime, so whether it is encrypted decides what may be put in it.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetApigeeEnvironmentKeyValueMapAttrs(t testing.TestingT, ctx context.Context, orgID string, envID string, mapID string) *apigee.GoogleCloudApigeeV1KeyValueMap {
	attrs, err := GetApigeeEnvironmentKeyValueMapAttrsE(t, ctx, orgID, envID, mapID)
	require.NoError(t, err)

	return attrs
}

// GetApigeeEnvironmentKeyValueMapAttrsE returns the settings Google Cloud holds for the given Apigee environment key value map.
// The ctx parameter supports cancellation and timeouts.
func GetApigeeEnvironmentKeyValueMapAttrsE(t testing.TestingT, ctx context.Context, orgID string, envID string, mapID string) (*apigee.GoogleCloudApigeeV1KeyValueMap, error) {
	logger.Default.Logf(t, "Getting settings for Apigee environment key value map %s %s in organization %s", mapID, envID, orgID)

	service, err := NewApigeeServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetApigeeEnvironmentKeyValueMapAttrsWithClient(ctx, service, orgID, envID, mapID)
}

// GetApigeeEnvironmentKeyValueMapAttrsWithClient returns the settings Google Cloud holds for the given Apigee environment key value map using the supplied
// *apigee.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see apigee_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetApigeeEnvironmentKeyValueMapAttrsWithClient(ctx context.Context, service *apigee.Service, orgID string, envID string, mapID string) (*apigee.GoogleCloudApigeeV1KeyValueMap, error) {
	name := fmt.Sprintf("organizations/%s/environments/%s/keyvaluemaps/%s", orgID, envID, mapID)

	attrs, err := service.Organizations.Environments.Keyvaluemaps.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Apigee environment key value map %s %s in organization %s does not exist", mapID, envID, orgID)
		}

		return nil, fmt.Errorf("failed to get settings for Apigee environment key value map %s %s in organization %s: %w", mapID, envID, orgID, err)
	}

	return attrs, nil
}

// GetApigeeEnvironmentKeyValueEntryAttrs returns the settings Google Cloud holds for the given Apigee environment key value entry, so a test can assert on what was
// actually created rather than only that it exists.
// An entry is one value in a map, so the name and the value are all there is to assert.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetApigeeEnvironmentKeyValueEntryAttrs(t testing.TestingT, ctx context.Context, orgID string, envID string, mapID string, entryID string) *apigee.GoogleCloudApigeeV1KeyValueEntry {
	attrs, err := GetApigeeEnvironmentKeyValueEntryAttrsE(t, ctx, orgID, envID, mapID, entryID)
	require.NoError(t, err)

	return attrs
}

// GetApigeeEnvironmentKeyValueEntryAttrsE returns the settings Google Cloud holds for the given Apigee environment key value entry.
// The ctx parameter supports cancellation and timeouts.
func GetApigeeEnvironmentKeyValueEntryAttrsE(t testing.TestingT, ctx context.Context, orgID string, envID string, mapID string, entryID string) (*apigee.GoogleCloudApigeeV1KeyValueEntry, error) {
	logger.Default.Logf(t, "Getting settings for Apigee environment key value entry %s %s %s in organization %s", entryID, mapID, envID, orgID)

	service, err := NewApigeeServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetApigeeEnvironmentKeyValueEntryAttrsWithClient(ctx, service, orgID, envID, mapID, entryID)
}

// GetApigeeEnvironmentKeyValueEntryAttrsWithClient returns the settings Google Cloud holds for the given Apigee environment key value entry using the supplied
// *apigee.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see apigee_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetApigeeEnvironmentKeyValueEntryAttrsWithClient(ctx context.Context, service *apigee.Service, orgID string, envID string, mapID string, entryID string) (*apigee.GoogleCloudApigeeV1KeyValueEntry, error) {
	name := fmt.Sprintf("organizations/%s/environments/%s/keyvaluemaps/%s/entries/%s", orgID, envID, mapID, entryID)

	attrs, err := service.Organizations.Environments.Keyvaluemaps.Entries.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Apigee environment key value entry %s %s %s in organization %s does not exist", entryID, mapID, envID, orgID)
		}

		return nil, fmt.Errorf("failed to get settings for Apigee environment key value entry %s %s %s in organization %s: %w", entryID, mapID, envID, orgID, err)
	}

	return attrs, nil
}

// GetApigeeEnvironmentAddonsConfigAttrs returns the settings Google Cloud holds for the given Apigee environment addons config, so a test can assert on what was
// actually created rather than only that it exists.
// The config says which paid add-ons an environment runs, so what is enabled here is what is billed.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetApigeeEnvironmentAddonsConfigAttrs(t testing.TestingT, ctx context.Context, orgID string, envID string) *apigee.GoogleCloudApigeeV1AddonsConfig {
	attrs, err := GetApigeeEnvironmentAddonsConfigAttrsE(t, ctx, orgID, envID)
	require.NoError(t, err)

	return attrs
}

// GetApigeeEnvironmentAddonsConfigAttrsE returns the settings Google Cloud holds for the given Apigee environment addons config.
// The ctx parameter supports cancellation and timeouts.
func GetApigeeEnvironmentAddonsConfigAttrsE(t testing.TestingT, ctx context.Context, orgID string, envID string) (*apigee.GoogleCloudApigeeV1AddonsConfig, error) {
	logger.Default.Logf(t, "Getting settings for Apigee environment addons config %s in organization %s", envID, orgID)

	service, err := NewApigeeServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetApigeeEnvironmentAddonsConfigAttrsWithClient(ctx, service, orgID, envID)
}

// GetApigeeEnvironmentAddonsConfigAttrsWithClient returns the settings Google Cloud holds for the given Apigee environment addons config using the supplied
// *apigee.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see apigee_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetApigeeEnvironmentAddonsConfigAttrsWithClient(ctx context.Context, service *apigee.Service, orgID string, envID string) (*apigee.GoogleCloudApigeeV1AddonsConfig, error) {
	name := fmt.Sprintf("organizations/%s/environments/%s/addonsConfig", orgID, envID)

	attrs, err := service.Organizations.Environments.GetAddonsConfig(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Apigee environment addons config %s in organization %s does not exist", envID, orgID)
		}

		return nil, fmt.Errorf("failed to get settings for Apigee environment addons config %s in organization %s: %w", envID, orgID, err)
	}

	return attrs, nil
}

// GetApigeeEnvironmentDebugMaskAttrs returns the settings Google Cloud holds for the given Apigee environment debug mask, so a test can assert on what was
// actually created rather than only that it exists.
// The mask is what keeps secrets out of trace output, so the fields it hides are the whole point of it.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetApigeeEnvironmentDebugMaskAttrs(t testing.TestingT, ctx context.Context, orgID string, envID string) *apigee.GoogleCloudApigeeV1DebugMask {
	attrs, err := GetApigeeEnvironmentDebugMaskAttrsE(t, ctx, orgID, envID)
	require.NoError(t, err)

	return attrs
}

// GetApigeeEnvironmentDebugMaskAttrsE returns the settings Google Cloud holds for the given Apigee environment debug mask.
// The ctx parameter supports cancellation and timeouts.
func GetApigeeEnvironmentDebugMaskAttrsE(t testing.TestingT, ctx context.Context, orgID string, envID string) (*apigee.GoogleCloudApigeeV1DebugMask, error) {
	logger.Default.Logf(t, "Getting settings for Apigee environment debug mask %s in organization %s", envID, orgID)

	service, err := NewApigeeServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetApigeeEnvironmentDebugMaskAttrsWithClient(ctx, service, orgID, envID)
}

// GetApigeeEnvironmentDebugMaskAttrsWithClient returns the settings Google Cloud holds for the given Apigee environment debug mask using the supplied
// *apigee.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see apigee_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetApigeeEnvironmentDebugMaskAttrsWithClient(ctx context.Context, service *apigee.Service, orgID string, envID string) (*apigee.GoogleCloudApigeeV1DebugMask, error) {
	name := fmt.Sprintf("organizations/%s/environments/%s/debugmask", orgID, envID)

	attrs, err := service.Organizations.Environments.GetDebugmask(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Apigee environment debug mask %s in organization %s does not exist", envID, orgID)
		}

		return nil, fmt.Errorf("failed to get settings for Apigee environment debug mask %s in organization %s: %w", envID, orgID, err)
	}

	return attrs, nil
}

// GetApigeeSecurityActionAttrs returns the settings Google Cloud holds for the given Apigee security action, so a test can assert on what was
// actually created rather than only that it exists.
// An action is what Apigee does to traffic that matches a condition, so the condition and whether it denies or flags are the point.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetApigeeSecurityActionAttrs(t testing.TestingT, ctx context.Context, orgID string, envID string, actionID string) *apigee.GoogleCloudApigeeV1SecurityAction {
	attrs, err := GetApigeeSecurityActionAttrsE(t, ctx, orgID, envID, actionID)
	require.NoError(t, err)

	return attrs
}

// GetApigeeSecurityActionAttrsE returns the settings Google Cloud holds for the given Apigee security action.
// The ctx parameter supports cancellation and timeouts.
func GetApigeeSecurityActionAttrsE(t testing.TestingT, ctx context.Context, orgID string, envID string, actionID string) (*apigee.GoogleCloudApigeeV1SecurityAction, error) {
	logger.Default.Logf(t, "Getting settings for Apigee security action %s %s in organization %s", actionID, envID, orgID)

	service, err := NewApigeeServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetApigeeSecurityActionAttrsWithClient(ctx, service, orgID, envID, actionID)
}

// GetApigeeSecurityActionAttrsWithClient returns the settings Google Cloud holds for the given Apigee security action using the supplied
// *apigee.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see apigee_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetApigeeSecurityActionAttrsWithClient(ctx context.Context, service *apigee.Service, orgID string, envID string, actionID string) (*apigee.GoogleCloudApigeeV1SecurityAction, error) {
	name := fmt.Sprintf("organizations/%s/environments/%s/securityActions/%s", orgID, envID, actionID)

	attrs, err := service.Organizations.Environments.SecurityActions.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Apigee security action %s %s in organization %s does not exist", actionID, envID, orgID)
		}

		return nil, fmt.Errorf("failed to get settings for Apigee security action %s %s in organization %s: %w", actionID, envID, orgID, err)
	}

	return attrs, nil
}

// GetApigeeSecurityFeedbackAttrs returns the settings Google Cloud holds for the given Apigee security feedback, so a test can assert on what was
// actually created rather than only that it exists.
// Feedback is how a human marks Apigee's own assessment right or wrong, so the verdict it carries is what it is for.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetApigeeSecurityFeedbackAttrs(t testing.TestingT, ctx context.Context, orgID string, feedbackID string) *apigee.GoogleCloudApigeeV1SecurityFeedback {
	attrs, err := GetApigeeSecurityFeedbackAttrsE(t, ctx, orgID, feedbackID)
	require.NoError(t, err)

	return attrs
}

// GetApigeeSecurityFeedbackAttrsE returns the settings Google Cloud holds for the given Apigee security feedback.
// The ctx parameter supports cancellation and timeouts.
func GetApigeeSecurityFeedbackAttrsE(t testing.TestingT, ctx context.Context, orgID string, feedbackID string) (*apigee.GoogleCloudApigeeV1SecurityFeedback, error) {
	logger.Default.Logf(t, "Getting settings for Apigee security feedback %s in organization %s", feedbackID, orgID)

	service, err := NewApigeeServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetApigeeSecurityFeedbackAttrsWithClient(ctx, service, orgID, feedbackID)
}

// GetApigeeSecurityFeedbackAttrsWithClient returns the settings Google Cloud holds for the given Apigee security feedback using the supplied
// *apigee.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see apigee_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetApigeeSecurityFeedbackAttrsWithClient(ctx context.Context, service *apigee.Service, orgID string, feedbackID string) (*apigee.GoogleCloudApigeeV1SecurityFeedback, error) {
	name := fmt.Sprintf("organizations/%s/securityFeedback/%s", orgID, feedbackID)

	attrs, err := service.Organizations.SecurityFeedback.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Apigee security feedback %s in organization %s does not exist", feedbackID, orgID)
		}

		return nil, fmt.Errorf("failed to get settings for Apigee security feedback %s in organization %s: %w", feedbackID, orgID, err)
	}

	return attrs, nil
}

// GetApigeeSecurityMonitoringConditionAttrs returns the settings Google Cloud holds for the given Apigee security monitoring condition, so a test can assert on what was
// actually created rather than only that it exists.
// A condition is what ties a security profile to an environment, so the profile and scope it names decide what is watched.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetApigeeSecurityMonitoringConditionAttrs(t testing.TestingT, ctx context.Context, orgID string, conditionID string) *apigee.GoogleCloudApigeeV1SecurityMonitoringCondition {
	attrs, err := GetApigeeSecurityMonitoringConditionAttrsE(t, ctx, orgID, conditionID)
	require.NoError(t, err)

	return attrs
}

// GetApigeeSecurityMonitoringConditionAttrsE returns the settings Google Cloud holds for the given Apigee security monitoring condition.
// The ctx parameter supports cancellation and timeouts.
func GetApigeeSecurityMonitoringConditionAttrsE(t testing.TestingT, ctx context.Context, orgID string, conditionID string) (*apigee.GoogleCloudApigeeV1SecurityMonitoringCondition, error) {
	logger.Default.Logf(t, "Getting settings for Apigee security monitoring condition %s in organization %s", conditionID, orgID)

	service, err := NewApigeeServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetApigeeSecurityMonitoringConditionAttrsWithClient(ctx, service, orgID, conditionID)
}

// GetApigeeSecurityMonitoringConditionAttrsWithClient returns the settings Google Cloud holds for the given Apigee security monitoring condition using the supplied
// *apigee.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see apigee_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetApigeeSecurityMonitoringConditionAttrsWithClient(ctx context.Context, service *apigee.Service, orgID string, conditionID string) (*apigee.GoogleCloudApigeeV1SecurityMonitoringCondition, error) {
	name := fmt.Sprintf("organizations/%s/securityMonitoringConditions/%s", orgID, conditionID)

	attrs, err := service.Organizations.SecurityMonitoringConditions.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Apigee security monitoring condition %s in organization %s does not exist", conditionID, orgID)
		}

		return nil, fmt.Errorf("failed to get settings for Apigee security monitoring condition %s in organization %s: %w", conditionID, orgID, err)
	}

	return attrs, nil
}

// GetApigeeSecurityProfileV2Attrs returns the settings Google Cloud holds for the given Apigee security profile, so a test can assert on what was
// actually created rather than only that it exists.
// A profile is the set of assessments run against traffic, so which assessments it carries is the whole point.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetApigeeSecurityProfileV2Attrs(t testing.TestingT, ctx context.Context, orgID string, profileID string) *apigee.GoogleCloudApigeeV1SecurityProfileV2 {
	attrs, err := GetApigeeSecurityProfileV2AttrsE(t, ctx, orgID, profileID)
	require.NoError(t, err)

	return attrs
}

// GetApigeeSecurityProfileV2AttrsE returns the settings Google Cloud holds for the given Apigee security profile.
// The ctx parameter supports cancellation and timeouts.
func GetApigeeSecurityProfileV2AttrsE(t testing.TestingT, ctx context.Context, orgID string, profileID string) (*apigee.GoogleCloudApigeeV1SecurityProfileV2, error) {
	logger.Default.Logf(t, "Getting settings for Apigee security profile %s in organization %s", profileID, orgID)

	service, err := NewApigeeServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetApigeeSecurityProfileV2AttrsWithClient(ctx, service, orgID, profileID)
}

// GetApigeeSecurityProfileV2AttrsWithClient returns the settings Google Cloud holds for the given Apigee security profile using the supplied
// *apigee.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see apigee_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetApigeeSecurityProfileV2AttrsWithClient(ctx context.Context, service *apigee.Service, orgID string, profileID string) (*apigee.GoogleCloudApigeeV1SecurityProfileV2, error) {
	name := fmt.Sprintf("organizations/%s/securityProfilesV2/%s", orgID, profileID)

	attrs, err := service.Organizations.SecurityProfilesV2.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Apigee security profile %s in organization %s does not exist", profileID, orgID)
		}

		return nil, fmt.Errorf("failed to get settings for Apigee security profile %s in organization %s: %w", profileID, orgID, err)
	}

	return attrs, nil
}

// GetApigeeControlPlaneAccessAttrs returns the settings Google Cloud holds for the given Apigee control plane access, so a test can assert on what was
// actually created rather than only that it exists.
// This is who may publish to the control plane, so the principals it lists are exactly the access it grants.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetApigeeControlPlaneAccessAttrs(t testing.TestingT, ctx context.Context, orgID string) *apigee.GoogleCloudApigeeV1ControlPlaneAccess {
	attrs, err := GetApigeeControlPlaneAccessAttrsE(t, ctx, orgID)
	require.NoError(t, err)

	return attrs
}

// GetApigeeControlPlaneAccessAttrsE returns the settings Google Cloud holds for the given Apigee control plane access.
// The ctx parameter supports cancellation and timeouts.
func GetApigeeControlPlaneAccessAttrsE(t testing.TestingT, ctx context.Context, orgID string) (*apigee.GoogleCloudApigeeV1ControlPlaneAccess, error) {
	logger.Default.Logf(t, "Getting settings for Apigee control plane access %s", orgID)

	service, err := NewApigeeServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetApigeeControlPlaneAccessAttrsWithClient(ctx, service, orgID)
}

// GetApigeeControlPlaneAccessAttrsWithClient returns the settings Google Cloud holds for the given Apigee control plane access using the supplied
// *apigee.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see apigee_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetApigeeControlPlaneAccessAttrsWithClient(ctx context.Context, service *apigee.Service, orgID string) (*apigee.GoogleCloudApigeeV1ControlPlaneAccess, error) {
	name := fmt.Sprintf("organizations/%s/controlPlaneAccess", orgID)

	attrs, err := service.Organizations.GetControlPlaneAccess(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Apigee control plane access %s does not exist", orgID)
		}

		return nil, fmt.Errorf("failed to get settings for Apigee control plane access %s: %w", orgID, err)
	}

	return attrs, nil
}

// NewApigeeServiceE creates an Apigee service authenticated the same way every other client in this
// module is.
// The ctx parameter supports cancellation and timeouts.
func NewApigeeServiceE(t testing.TestingT, ctx context.Context) (*apigee.Service, error) {
	return apigee.NewService(ctx, append(withOptions(), option.WithScopes(apigee.CloudPlatformScope))...)
}
