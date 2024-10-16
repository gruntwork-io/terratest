package aws

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/stretchr/testify/assert"
)

func TestNewAuthenticatedSessionFromDefaultCredentialsShouldReturnCustomEnpointURLForS3Service(t *testing.T) {
	t.Parallel()

	const signingRegionExpected = "custom-signing-region"
	region := GetRandomStableRegion(t, nil, nil)

	SetAwsEndpointsOverrides(localEndpoints)
	session, errors := NewAuthenticatedSessionFromDefaultCredentials(region)
	resolver := session.Config.EndpointResolver
	resolverExpected, errors := resolver.EndpointFor(endpoints.S3ServiceID, signingRegionExpected)

	assert.NotNil(t, session)
	assert.NoError(t, errors)
	assert.Equal(t, signingRegionExpected, resolverExpected.SigningRegion)
	assert.Equal(t, localEndpoints[endpoints.S3ServiceID], resolverExpected.URL)
}
