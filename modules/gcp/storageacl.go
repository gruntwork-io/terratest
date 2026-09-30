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

// The rest of this package reads Cloud Storage through cloud.google.com/go/storage, which does not
// expose per-object ACL entries as a resource a test can read by entity. These reads use the REST
// client instead, which is why the package is aliased here rather than in storage.go.

// GetDefaultObjectACLAttrs returns the default object ACL entry Google Cloud holds for the given
// bucket and entity, so a test can assert on the role a new object would inherit. A bucket with
// uniform access carries no such entry, so this only answers for one that keeps fine grained ACLs.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDefaultObjectACLAttrs(t testing.TestingT, ctx context.Context, bucketName string, entity string) *storagev1.ObjectAccessControl {
	acl, err := GetDefaultObjectACLAttrsE(t, ctx, bucketName, entity)
	require.NoError(t, err)

	return acl
}

// GetDefaultObjectACLAttrsE returns the default object ACL entry Google Cloud holds for the given
// bucket and entity.
// The ctx parameter supports cancellation and timeouts.
func GetDefaultObjectACLAttrsE(t testing.TestingT, ctx context.Context, bucketName string, entity string) (*storagev1.ObjectAccessControl, error) {
	logger.Default.Logf(t, "Getting the default object ACL entry for %s on bucket %s", entity, bucketName)

	service, err := NewStorageRESTServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetDefaultObjectACLAttrsWithClient(ctx, service, bucketName, entity)
}

// GetDefaultObjectACLAttrsWithClient returns the default object ACL entry Google Cloud holds for the
// given bucket and entity using the supplied *storagev1.Service. Prefer this variant in unit tests
// where the service is backed by an httptest fake server (see storageacl_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetDefaultObjectACLAttrsWithClient(ctx context.Context, service *storagev1.Service, bucketName string, entity string) (*storagev1.ObjectAccessControl, error) {
	acl, err := service.DefaultObjectAccessControls.Get(bucketName, entity).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("bucket %s has no default object ACL entry for %s", bucketName, entity)
		}

		return nil, fmt.Errorf("failed to get the default object ACL entry for %s on bucket %s: %w", entity, bucketName, err)
	}

	return acl, nil
}

// GetObjectACLAttrs returns the ACL entry Google Cloud holds for the given object and entity, so a
// test can assert on who may read one object rather than on the bucket's defaults.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetObjectACLAttrs(t testing.TestingT, ctx context.Context, bucketName string, objectName string, entity string) *storagev1.ObjectAccessControl {
	acl, err := GetObjectACLAttrsE(t, ctx, bucketName, objectName, entity)
	require.NoError(t, err)

	return acl
}

// GetObjectACLAttrsE returns the ACL entry Google Cloud holds for the given object and entity.
// The ctx parameter supports cancellation and timeouts.
func GetObjectACLAttrsE(t testing.TestingT, ctx context.Context, bucketName string, objectName string, entity string) (*storagev1.ObjectAccessControl, error) {
	logger.Default.Logf(t, "Getting the ACL entry for %s on object %s in bucket %s", entity, objectName, bucketName)

	service, err := NewStorageRESTServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetObjectACLAttrsWithClient(ctx, service, bucketName, objectName, entity)
}

// GetObjectACLAttrsWithClient returns the ACL entry Google Cloud holds for the given object and
// entity using the supplied *storagev1.Service. Prefer this variant in unit tests where the service
// is backed by an httptest fake server (see storageacl_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetObjectACLAttrsWithClient(ctx context.Context, service *storagev1.Service, bucketName string, objectName string, entity string) (*storagev1.ObjectAccessControl, error) {
	acl, err := service.ObjectAccessControls.Get(bucketName, objectName, entity).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("object %s in bucket %s has no ACL entry for %s", objectName, bucketName, entity)
		}

		return nil, fmt.Errorf("failed to get the ACL entry for %s on object %s in bucket %s: %w", entity, objectName, bucketName, err)
	}

	return acl, nil
}

// NewStorageRESTServiceE creates a Cloud Storage REST service authenticated the same way every other
// client in this module is. It serves the ACL calls that the idiomatic client does not expose.
// The ctx parameter supports cancellation and timeouts.
func NewStorageRESTServiceE(t testing.TestingT, ctx context.Context) (*storagev1.Service, error) {
	return storagev1.NewService(ctx, append(withOptions(), option.WithScopes(storagev1.CloudPlatformScope))...)
}
