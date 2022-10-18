//go:build azure
// +build azure

// NOTE: We use build tags to differentiate azure testing because we currently do not have azure access setup for
// CircleCI.

package azure

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetTargetAzureSubscription(t *testing.T) {
	t.Parallel()

	//Check that ARM_SUBSCRIPTION_ID env variable is set, CI requires this value to run all test.
	require.NotEmpty(t, os.Getenv(AzureSubscriptionID), "ARM_SUBSCRIPTION_ID environment variable not set.")

	type args struct {
		subID string
	}

	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{name: "subIDProvidedAsArg", args: args{subID: "test"}, want: "test", wantErr: false},
		{name: "subIDNotProvidedFallbackToEnv", args: args{subID: ""}, want: os.Getenv(AzureSubscriptionID), wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetTargetAzureSubscription(tt.args.subID)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.Equal(t, tt.want, got)
			}
		})
	}
}

func TestGetTargetAzureResourceGroupName(t *testing.T) {
	t.Parallel()

	type args struct {
		rgName string
	}

	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{name: "rgNameProvidedAsArg", args: args{rgName: "test"}, want: "test", wantErr: false},
		{name: "rgNameNotProvided", args: args{rgName: ""}, want: "", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetTargetAzureResourceGroupName(tt.args.rgName)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.Equal(t, tt.want, got)
			}
		})
	}
}

func TestSafePtrToString(t *testing.T) {
	// When given a nil, should always return an empty string
	var nilPtr *string = nil
	nilResult := safePtrToString(nilPtr)
	assert.Equal(t, "", nilResult)

	// When given a string, should just de-ref and return
	stringPtr := "Test"
	stringResult := safePtrToString(&stringPtr)
	assert.Equal(t, "Test", stringResult)
}

func TestSafePtrToInt32(t *testing.T) {
	// When given a nil, should always return an zero value int32
	var nilPtr *int32 = nil
	nilResult := safePtrToInt32(nilPtr)
	assert.Equal(t, int32(0), nilResult)

	// When given a string, should just de-ref and return
	intPtr := int32(42)
	intResult := safePtrToInt32(&intPtr)
	assert.Equal(t, int32(42), intResult)
}
