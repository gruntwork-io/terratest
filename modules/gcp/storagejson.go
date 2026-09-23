package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
	storagev1 "google.golang.org/api/storage/v1"
)

// GetStorageFolderAttrs returns the settings Google Cloud holds for the given folder in a Cloud Storage bucket, so a test can
// assert on what was actually created rather than only that it exists. Folders exist only in a bucket with a hierarchical
// namespace, and their names end in a slash.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetStorageFolderAttrs(t testing.TestingT, ctx context.Context, bucket string, folderName string) *storagev1.Folder {
	folder, err := GetStorageFolderAttrsE(t, ctx, bucket, folderName)
	require.NoError(t, err)

	return folder
}

// GetStorageFolderAttrsE returns the settings Google Cloud holds for the given folder in a Cloud Storage bucket.
// The ctx parameter supports cancellation and timeouts.
func GetStorageFolderAttrsE(t testing.TestingT, ctx context.Context, bucket string, folderName string) (*storagev1.Folder, error) {
	logger.Default.Logf(t, "Getting settings for folder %s in bucket %s", folderName, bucket)

	service, err := NewStorageJSONServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetStorageFolderAttrsWithClient(ctx, service, bucket, folderName)
}

// GetStorageFolderAttrsWithClient returns the settings Google Cloud holds for the given folder in a Cloud Storage bucket using
// the supplied *storagev1.Service. Prefer this variant in unit tests where the service is backed by
// an httptest fake server (see storagejson_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetStorageFolderAttrsWithClient(ctx context.Context, service *storagev1.Service, bucket string, folderName string) (*storagev1.Folder, error) {
	folder, err := service.Folders.Get(bucket, folderName).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the folder %s does not exist in bucket %s", folderName, bucket)
		}

		return nil, fmt.Errorf("failed to get settings for folder %s in bucket %s: %w", folderName, bucket, err)
	}

	return folder, nil
}

// GetStorageManagedFolderAttrs returns the settings Google Cloud holds for the given managed folder in a Cloud Storage bucket, so a test can
// assert on what was actually created rather than only that it exists. A managed folder carries its own access policy, and its name ends in a slash.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetStorageManagedFolderAttrs(t testing.TestingT, ctx context.Context, bucket string, folderName string) *storagev1.ManagedFolder {
	managedFolder, err := GetStorageManagedFolderAttrsE(t, ctx, bucket, folderName)
	require.NoError(t, err)

	return managedFolder
}

// GetStorageManagedFolderAttrsE returns the settings Google Cloud holds for the given managed folder in a Cloud Storage bucket.
// The ctx parameter supports cancellation and timeouts.
func GetStorageManagedFolderAttrsE(t testing.TestingT, ctx context.Context, bucket string, folderName string) (*storagev1.ManagedFolder, error) {
	logger.Default.Logf(t, "Getting settings for managed folder %s in bucket %s", folderName, bucket)

	service, err := NewStorageJSONServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetStorageManagedFolderAttrsWithClient(ctx, service, bucket, folderName)
}

// GetStorageManagedFolderAttrsWithClient returns the settings Google Cloud holds for the given managed folder in a Cloud Storage bucket using
// the supplied *storagev1.Service. Prefer this variant in unit tests where the service is backed by
// an httptest fake server (see storagejson_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetStorageManagedFolderAttrsWithClient(ctx context.Context, service *storagev1.Service, bucket string, folderName string) (*storagev1.ManagedFolder, error) {
	managedFolder, err := service.ManagedFolders.Get(bucket, folderName).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the managed folder %s does not exist in bucket %s", folderName, bucket)
		}

		return nil, fmt.Errorf("failed to get settings for managed folder %s in bucket %s: %w", folderName, bucket, err)
	}

	return managedFolder, nil
}

// GetStorageHmacKeyAttrs returns the settings Google Cloud holds for the given HMAC key, so a test can
// assert on what was actually created rather than only that it exists. Google assigns the access id, so the caller passes the one it
// got back rather than a name it chose, and the secret is never returned by this call.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetStorageHmacKeyAttrs(t testing.TestingT, ctx context.Context, projectID string, accessID string) *storagev1.HmacKeyMetadata {
	key, err := GetStorageHmacKeyAttrsE(t, ctx, projectID, accessID)
	require.NoError(t, err)

	return key
}

// GetStorageHmacKeyAttrsE returns the settings Google Cloud holds for the given HMAC key.
// The ctx parameter supports cancellation and timeouts.
func GetStorageHmacKeyAttrsE(t testing.TestingT, ctx context.Context, projectID string, accessID string) (*storagev1.HmacKeyMetadata, error) {
	logger.Default.Logf(t, "Getting settings for HMAC key %s in project %s", accessID, projectID)

	service, err := NewStorageJSONServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetStorageHmacKeyAttrsWithClient(ctx, service, projectID, accessID)
}

// GetStorageHmacKeyAttrsWithClient returns the settings Google Cloud holds for the given HMAC key using
// the supplied *storagev1.Service. Prefer this variant in unit tests where the service is backed by
// an httptest fake server (see storagejson_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetStorageHmacKeyAttrsWithClient(ctx context.Context, service *storagev1.Service, projectID string, accessID string) (*storagev1.HmacKeyMetadata, error) {
	key, err := service.Projects.HmacKeys.Get(projectID, accessID).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the HMAC key %s does not exist in project %s", accessID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for HMAC key %s in project %s: %w", accessID, projectID, err)
	}

	return key, nil
}

// GetStorageNotificationAttrs returns the settings Google Cloud holds for the given bucket notification, so a test can
// assert on what was actually created rather than only that it exists. Google assigns the notification its id, so the caller passes the
// one it got back rather than a name it chose.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetStorageNotificationAttrs(t testing.TestingT, ctx context.Context, bucket string, notificationID string) *storagev1.Notification {
	notification, err := GetStorageNotificationAttrsE(t, ctx, bucket, notificationID)
	require.NoError(t, err)

	return notification
}

// GetStorageNotificationAttrsE returns the settings Google Cloud holds for the given bucket notification.
// The ctx parameter supports cancellation and timeouts.
func GetStorageNotificationAttrsE(t testing.TestingT, ctx context.Context, bucket string, notificationID string) (*storagev1.Notification, error) {
	logger.Default.Logf(t, "Getting settings for notification %s on bucket %s", notificationID, bucket)

	service, err := NewStorageJSONServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetStorageNotificationAttrsWithClient(ctx, service, bucket, notificationID)
}

// GetStorageNotificationAttrsWithClient returns the settings Google Cloud holds for the given bucket notification using
// the supplied *storagev1.Service. Prefer this variant in unit tests where the service is backed by
// an httptest fake server (see storagejson_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetStorageNotificationAttrsWithClient(ctx context.Context, service *storagev1.Service, bucket string, notificationID string) (*storagev1.Notification, error) {
	notification, err := service.Notifications.Get(bucket, notificationID).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the notification %s does not exist on bucket %s", notificationID, bucket)
		}

		return nil, fmt.Errorf("failed to get settings for notification %s on bucket %s: %w", notificationID, bucket, err)
	}

	return notification, nil
}

// GetStorageObjectAccessControlAttrs returns the settings Google Cloud holds for the given entity's access to an object, so a test can
// assert on what was actually created rather than only that it exists. Object ACLs exist only on a bucket that has not
// turned on uniform bucket level access.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetStorageObjectAccessControlAttrs(t testing.TestingT, ctx context.Context, bucket string, object string, entity string) *storagev1.ObjectAccessControl {
	accessControl, err := GetStorageObjectAccessControlAttrsE(t, ctx, bucket, object, entity)
	require.NoError(t, err)

	return accessControl
}

// GetStorageObjectAccessControlAttrsE returns the settings Google Cloud holds for the given entity's access to an object.
// The ctx parameter supports cancellation and timeouts.
func GetStorageObjectAccessControlAttrsE(t testing.TestingT, ctx context.Context, bucket string, object string, entity string) (*storagev1.ObjectAccessControl, error) {
	logger.Default.Logf(t, "Getting settings for the access control for %s on object %s in bucket %s", entity, object, bucket)

	service, err := NewStorageJSONServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetStorageObjectAccessControlAttrsWithClient(ctx, service, bucket, object, entity)
}

// GetStorageObjectAccessControlAttrsWithClient returns the settings Google Cloud holds for the given entity's access to an object using
// the supplied *storagev1.Service. Prefer this variant in unit tests where the service is backed by
// an httptest fake server (see storagejson_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetStorageObjectAccessControlAttrsWithClient(ctx context.Context, service *storagev1.Service, bucket string, object string, entity string) (*storagev1.ObjectAccessControl, error) {
	accessControl, err := service.ObjectAccessControls.Get(bucket, object, entity).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("%s has no access control on object %s in bucket %s", entity, object, bucket)
		}

		return nil, fmt.Errorf("failed to get settings for the access control for %s on object %s in bucket %s: %w", entity, object, bucket, err)
	}

	return accessControl, nil
}

// GetStorageDefaultObjectAccessControlAttrs returns the settings Google Cloud holds for the given entity's default access to a bucket's objects, so a test can
// assert on what was actually created rather than only that it exists. The default applies to every object written
// afterwards that does not carry an access control of its own.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetStorageDefaultObjectAccessControlAttrs(t testing.TestingT, ctx context.Context, bucket string, entity string) *storagev1.ObjectAccessControl {
	accessControl, err := GetStorageDefaultObjectAccessControlAttrsE(t, ctx, bucket, entity)
	require.NoError(t, err)

	return accessControl
}

// GetStorageDefaultObjectAccessControlAttrsE returns the settings Google Cloud holds for the given entity's default access to a bucket's objects.
// The ctx parameter supports cancellation and timeouts.
func GetStorageDefaultObjectAccessControlAttrsE(t testing.TestingT, ctx context.Context, bucket string, entity string) (*storagev1.ObjectAccessControl, error) {
	logger.Default.Logf(t, "Getting settings for the default object access control for %s in bucket %s", entity, bucket)

	service, err := NewStorageJSONServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetStorageDefaultObjectAccessControlAttrsWithClient(ctx, service, bucket, entity)
}

// GetStorageDefaultObjectAccessControlAttrsWithClient returns the settings Google Cloud holds for the given entity's default access to a bucket's objects using
// the supplied *storagev1.Service. Prefer this variant in unit tests where the service is backed by
// an httptest fake server (see storagejson_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetStorageDefaultObjectAccessControlAttrsWithClient(ctx context.Context, service *storagev1.Service, bucket string, entity string) (*storagev1.ObjectAccessControl, error) {
	accessControl, err := service.DefaultObjectAccessControls.Get(bucket, entity).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("%s has no default object access control in bucket %s", entity, bucket)
		}

		return nil, fmt.Errorf("failed to get settings for the default object access control for %s in bucket %s: %w", entity, bucket, err)
	}

	return accessControl, nil
}

// GetStorageManagedFolderIamPolicyAttrs returns the IAM policy Google Cloud holds for the given
// managed folder, so a test can assert on who was actually granted access to it. A managed folder
// carries a policy of its own, separate from the bucket's.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetStorageManagedFolderIamPolicyAttrs(t testing.TestingT, ctx context.Context, bucket string, folderName string) *storagev1.Policy {
	policy, err := GetStorageManagedFolderIamPolicyAttrsE(t, ctx, bucket, folderName)
	require.NoError(t, err)

	return policy
}

// GetStorageManagedFolderIamPolicyAttrsE returns the IAM policy Google Cloud holds for the given
// managed folder.
// The ctx parameter supports cancellation and timeouts.
func GetStorageManagedFolderIamPolicyAttrsE(t testing.TestingT, ctx context.Context, bucket string, folderName string) (*storagev1.Policy, error) {
	logger.Default.Logf(t, "Getting the IAM policy for managed folder %s in bucket %s", folderName, bucket)

	service, err := NewStorageJSONServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetStorageManagedFolderIamPolicyAttrsWithClient(ctx, service, bucket, folderName)
}

// GetStorageManagedFolderIamPolicyAttrsWithClient returns the IAM policy Google Cloud holds for the
// given managed folder using the supplied *storagev1.Service. Prefer this variant in unit tests
// where the service is backed by an httptest fake server (see storagejson_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetStorageManagedFolderIamPolicyAttrsWithClient(ctx context.Context, service *storagev1.Service, bucket string, folderName string) (*storagev1.Policy, error) {
	policy, err := service.ManagedFolders.GetIamPolicy(bucket, folderName).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the managed folder %s does not exist in bucket %s", folderName, bucket)
		}

		return nil, fmt.Errorf("failed to get the IAM policy for managed folder %s in bucket %s: %w", folderName, bucket, err)
	}

	return policy, nil
}

// NewStorageJSONServiceE creates a Cloud Storage JSON API service authenticated the same way every
// other client in this module is. The Cloud Storage client used by storage.go covers buckets and
// objects; folders, managed folders, HMAC keys, notifications and access controls are only on this
// one.
// The ctx parameter supports cancellation and timeouts.
func NewStorageJSONServiceE(t testing.TestingT, ctx context.Context) (*storagev1.Service, error) {
	return storagev1.NewService(ctx, append(withOptions(), option.WithScopes(storagev1.CloudPlatformScope))...)
}
