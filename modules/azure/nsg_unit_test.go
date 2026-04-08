package azure

import (
	"context"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v6"
	networkfake "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v6/fake"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// convertToNsgRuleSummary tests
// ---------------------------------------------------------------------------

func TestConvertToNsgRuleSummary_FullyPopulated(t *testing.T) {
	t.Parallel()

	protocol := armnetwork.SecurityRuleProtocolTCP
	access := armnetwork.SecurityRuleAccessAllow
	direction := armnetwork.SecurityRuleDirectionInbound

	props := &armnetwork.SecurityRulePropertiesFormat{
		Protocol:                 &protocol,
		Access:                   &access,
		Direction:                &direction,
		Priority:                 to.Ptr[int32](100),
		Description:              to.Ptr("allow ssh"),
		SourcePortRange:          to.Ptr("*"),
		DestinationPortRange:     to.Ptr("22"),
		SourceAddressPrefix:      to.Ptr("10.0.0.0/8"),
		DestinationAddressPrefix: to.Ptr("VirtualNetwork"),
	}

	summary := convertToNsgRuleSummary(to.Ptr("AllowSSH"), props)

	assert.Equal(t, "AllowSSH", summary.Name)
	assert.Equal(t, "allow ssh", summary.Description)
	assert.Equal(t, string(armnetwork.SecurityRuleProtocolTCP), summary.Protocol)
	assert.Equal(t, string(armnetwork.SecurityRuleAccessAllow), summary.Access)
	assert.Equal(t, string(armnetwork.SecurityRuleDirectionInbound), summary.Direction)
	assert.Equal(t, int32(100), summary.Priority)
	assert.Equal(t, "*", summary.SourcePortRange)
	assert.Equal(t, "22", summary.DestinationPortRange)
	assert.Equal(t, "10.0.0.0/8", summary.SourceAddressPrefix)
	assert.Equal(t, "VirtualNetwork", summary.DestinationAddressPrefix)
}

func TestConvertToNsgRuleSummary_NilEnums(t *testing.T) {
	t.Parallel()

	props := &armnetwork.SecurityRulePropertiesFormat{
		Protocol:  nil,
		Access:    nil,
		Direction: nil,
		Priority:  to.Ptr[int32](200),
	}

	summary := convertToNsgRuleSummary(to.Ptr("SomeRule"), props)

	assert.Equal(t, "", summary.Protocol)
	assert.Equal(t, "", summary.Access)
	assert.Equal(t, "", summary.Direction)
	assert.Equal(t, int32(200), summary.Priority)
}

func TestConvertToNsgRuleSummary_NilName(t *testing.T) {
	t.Parallel()

	props := &armnetwork.SecurityRulePropertiesFormat{
		Priority: to.Ptr[int32](300),
	}

	summary := convertToNsgRuleSummary(nil, props)

	assert.Equal(t, "", summary.Name)
}

// ---------------------------------------------------------------------------
// safeDerefString tests
// ---------------------------------------------------------------------------

func TestSafeDerefString_Nil(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "", safeDerefString(nil))
}

func TestSafeDerefString_NonNil(t *testing.T) {
	t.Parallel()
	s := "hello"
	assert.Equal(t, "hello", safeDerefString(&s))
}

// ---------------------------------------------------------------------------
// Fake client helpers
// ---------------------------------------------------------------------------

func newFakeDefaultSecurityRulesClient(t *testing.T, srv networkfake.DefaultSecurityRulesServer) *armnetwork.DefaultSecurityRulesClient {
	t.Helper()
	client, err := armnetwork.NewDefaultSecurityRulesClient("fake-sub", &azfake.TokenCredential{},
		&arm.ClientOptions{ClientOptions: policy.ClientOptions{
			Transport: networkfake.NewDefaultSecurityRulesServerTransport(&srv),
		}})
	require.NoError(t, err)
	return client
}

func newFakeSecurityRulesClient(t *testing.T, srv networkfake.SecurityRulesServer) *armnetwork.SecurityRulesClient {
	t.Helper()
	client, err := armnetwork.NewSecurityRulesClient("fake-sub", &azfake.TokenCredential{},
		&arm.ClientOptions{ClientOptions: policy.ClientOptions{
			Transport: networkfake.NewSecurityRulesServerTransport(&srv),
		}})
	require.NoError(t, err)
	return client
}

// ---------------------------------------------------------------------------
// collectDefaultSecurityRules tests
// ---------------------------------------------------------------------------

func TestCollectDefaultSecurityRules_OnePage(t *testing.T) {
	t.Parallel()

	protocol := armnetwork.SecurityRuleProtocolTCP
	access := armnetwork.SecurityRuleAccessAllow
	direction := armnetwork.SecurityRuleDirectionInbound

	srv := networkfake.DefaultSecurityRulesServer{
		NewListPager: func(resourceGroupName, nsgName string, options *armnetwork.DefaultSecurityRulesClientListOptions) (resp azfake.PagerResponder[armnetwork.DefaultSecurityRulesClientListResponse]) {
			resp.AddPage(http.StatusOK, armnetwork.DefaultSecurityRulesClientListResponse{
				SecurityRuleListResult: armnetwork.SecurityRuleListResult{
					Value: []*armnetwork.SecurityRule{
						{
							Name: to.Ptr("AllowVnetInBound"),
							Properties: &armnetwork.SecurityRulePropertiesFormat{
								Protocol:                 &protocol,
								Access:                   &access,
								Direction:                &direction,
								Priority:                 to.Ptr[int32](65000),
								SourcePortRange:          to.Ptr("*"),
								DestinationPortRange:     to.Ptr("*"),
								SourceAddressPrefix:      to.Ptr("VirtualNetwork"),
								DestinationAddressPrefix: to.Ptr("VirtualNetwork"),
							},
						},
					},
				},
			}, nil)
			return
		},
	}

	client := newFakeDefaultSecurityRulesClient(t, srv)
	rules, err := collectDefaultSecurityRules(context.Background(), client, "rg", "nsg")
	require.NoError(t, err)
	require.Len(t, rules, 1)

	rule := rules[0]
	assert.Equal(t, "AllowVnetInBound", rule.Name)
	assert.Equal(t, string(armnetwork.SecurityRuleProtocolTCP), rule.Protocol)
	assert.Equal(t, string(armnetwork.SecurityRuleAccessAllow), rule.Access)
	assert.Equal(t, string(armnetwork.SecurityRuleDirectionInbound), rule.Direction)
	assert.Equal(t, int32(65000), rule.Priority)
	assert.Equal(t, "*", rule.SourcePortRange)
	assert.Equal(t, "*", rule.DestinationPortRange)
	assert.Equal(t, "VirtualNetwork", rule.SourceAddressPrefix)
	assert.Equal(t, "VirtualNetwork", rule.DestinationAddressPrefix)
}

func TestCollectDefaultSecurityRules_Empty(t *testing.T) {
	t.Parallel()

	srv := networkfake.DefaultSecurityRulesServer{
		NewListPager: func(resourceGroupName, nsgName string, options *armnetwork.DefaultSecurityRulesClientListOptions) (resp azfake.PagerResponder[armnetwork.DefaultSecurityRulesClientListResponse]) {
			resp.AddPage(http.StatusOK, armnetwork.DefaultSecurityRulesClientListResponse{
				SecurityRuleListResult: armnetwork.SecurityRuleListResult{
					Value: []*armnetwork.SecurityRule{},
				},
			}, nil)
			return
		},
	}

	client := newFakeDefaultSecurityRulesClient(t, srv)
	rules, err := collectDefaultSecurityRules(context.Background(), client, "rg", "nsg")
	require.NoError(t, err)
	assert.Empty(t, rules)
}

// ---------------------------------------------------------------------------
// collectCustomSecurityRules tests
// ---------------------------------------------------------------------------

func TestCollectCustomSecurityRules_OnePage(t *testing.T) {
	t.Parallel()

	protocol := armnetwork.SecurityRuleProtocolUDP
	access := armnetwork.SecurityRuleAccessDeny
	direction := armnetwork.SecurityRuleDirectionOutbound

	srv := networkfake.SecurityRulesServer{
		NewListPager: func(resourceGroupName, nsgName string, options *armnetwork.SecurityRulesClientListOptions) (resp azfake.PagerResponder[armnetwork.SecurityRulesClientListResponse]) {
			resp.AddPage(http.StatusOK, armnetwork.SecurityRulesClientListResponse{
				SecurityRuleListResult: armnetwork.SecurityRuleListResult{
					Value: []*armnetwork.SecurityRule{
						{
							Name: to.Ptr("DenyUDPOut"),
							Properties: &armnetwork.SecurityRulePropertiesFormat{
								Protocol:                 &protocol,
								Access:                   &access,
								Direction:                &direction,
								Priority:                 to.Ptr[int32](500),
								SourcePortRange:          to.Ptr("*"),
								DestinationPortRange:     to.Ptr("53"),
								SourceAddressPrefix:      to.Ptr("10.0.0.0/8"),
								DestinationAddressPrefix: to.Ptr("Internet"),
							},
						},
					},
				},
			}, nil)
			return
		},
	}

	client := newFakeSecurityRulesClient(t, srv)
	rules, err := collectCustomSecurityRules(context.Background(), client, "rg", "nsg")
	require.NoError(t, err)
	require.Len(t, rules, 1)

	rule := rules[0]
	assert.Equal(t, "DenyUDPOut", rule.Name)
	assert.Equal(t, string(armnetwork.SecurityRuleProtocolUDP), rule.Protocol)
	assert.Equal(t, string(armnetwork.SecurityRuleAccessDeny), rule.Access)
	assert.Equal(t, string(armnetwork.SecurityRuleDirectionOutbound), rule.Direction)
	assert.Equal(t, int32(500), rule.Priority)
	assert.Equal(t, "53", rule.DestinationPortRange)
	assert.Equal(t, "10.0.0.0/8", rule.SourceAddressPrefix)
	assert.Equal(t, "Internet", rule.DestinationAddressPrefix)
}
