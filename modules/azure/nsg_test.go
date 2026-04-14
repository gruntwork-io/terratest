package azure //nolint:testpackage // tests access unexported functions

import (
	"context"
	"fmt"
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

// Fake client helpers

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

// Port range parsing

func TestPortRangeParsing(t *testing.T) {
	t.Parallel()

	cases := []struct {
		portRange    string
		expectedLo   int
		expectedHi   int
		expectsError bool
	}{
		{"22", 22, 22, false},
		{"22-80", 22, 80, false},
		{"*", 0, 65535, false},
		{"*-*", 0, 0, true},
		{"22-", 0, 0, true},
		{"-80", 0, 0, true},
		{"-", 0, 0, true},
		{"80-22", 22, 80, false},
	}

	for _, tt := range cases {
		t.Run(tt.portRange, func(t *testing.T) {
			t.Parallel()

			lo, hi, err := parsePortRangeString(tt.portRange)
			if !tt.expectsError {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.expectedLo, int(lo))
			assert.Equal(t, tt.expectedHi, int(hi))
		})
	}
}

// Rule summary conversion

func TestNsgRuleSummaryConversion(t *testing.T) {
	t.Parallel()

	t.Run("FullyPopulated", func(t *testing.T) {
		t.Parallel()

		protocol := armnetwork.SecurityRuleProtocolTCP
		access := armnetwork.SecurityRuleAccessAllow
		direction := armnetwork.SecurityRuleDirectionInbound
		props := &armnetwork.SecurityRulePropertiesFormat{
			Protocol: &protocol, Access: &access, Direction: &direction,
			Priority: to.Ptr[int32](100), Description: to.Ptr("allow ssh"),
			SourcePortRange: to.Ptr("*"), DestinationPortRange: to.Ptr("22"),
			SourceAddressPrefix: to.Ptr("10.0.0.0/8"), DestinationAddressPrefix: to.Ptr("VirtualNetwork"),
		}
		summary := convertToNsgRuleSummary(to.Ptr("AllowSSH"), props)
		assert.Equal(t, "AllowSSH", summary.Name)
		assert.Equal(t, "Tcp", summary.Protocol)
		assert.Equal(t, "Allow", summary.Access)
		assert.Equal(t, int32(100), summary.Priority)
	})

	t.Run("NilEnums", func(t *testing.T) {
		t.Parallel()

		props := &armnetwork.SecurityRulePropertiesFormat{Priority: to.Ptr[int32](200)}
		summary := convertToNsgRuleSummary(to.Ptr("Rule"), props)
		assert.Empty(t, summary.Protocol)
		assert.Empty(t, summary.Access)
		assert.Empty(t, summary.Direction)
	})

	t.Run("NilName", func(t *testing.T) {
		t.Parallel()

		props := &armnetwork.SecurityRulePropertiesFormat{}
		summary := convertToNsgRuleSummary(nil, props)
		assert.Empty(t, summary.Name)
	})
}

// Port allow/deny

func TestAllowSourcePort(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name, portRange, access, testPort string
		result                            bool
	}{
		{"22 allowed", "22", "Allow", "22", true},
		{"22 denied", "22", "Deny", "22", false},
		{"Any allows any", "*", "Allow", "*", true},
		{"Range allows", "80-90", "Allow", "85", true},
		{"Range denies", "80-90", "Deny", "85", false},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := NsgRuleSummary{SourcePortRange: tt.portRange, Access: tt.access}
			assert.Equal(t, tt.result, s.AllowsSourcePort(t, tt.testPort))
		})
	}
}

func TestAllowDestinationPort(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name, portRange, access, testPort string
		result                            bool
	}{
		{"22 allowed", "22", "Allow", "22", true},
		{"22 denied", "22", "Deny", "22", false},
		{"Any allows any", "*", "Allow", "*", true},
		{"Range allows", "80-90", "Allow", "85", true},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := NsgRuleSummary{DestinationPortRange: tt.portRange, Access: tt.access}
			assert.Equal(t, tt.result, s.AllowsDestinationPort(t, tt.testPort))
		})
	}
}

// Rule finding

func TestFindSummarizedRule(t *testing.T) {
	t.Parallel()

	rules := make([]NsgRuleSummary, 5)
	for i := range rules {
		rules[i].Name = fmt.Sprintf("rule_%d", i+1)
	}

	ruleList := NsgRuleSummaryList{SummarizedRules: rules}

	assert.Equal(t, "rule_1", ruleList.FindRuleByName("rule_1").Name)
	assert.Equal(t, NsgRuleSummary{}, ruleList.FindRuleByName("nonexistent"))
}

// safeDerefString

func TestSafeDerefString(t *testing.T) {
	t.Parallel()

	assert.Empty(t, safeDerefString(nil))

	s := "hello"
	assert.Equal(t, "hello", safeDerefString(&s))
}

// GetDefaultNSGRulesWithClient / GetCustomNSGRulesWithClient

func TestGetDefaultNSGRulesWithClient(t *testing.T) {
	t.Parallel()

	protocol := armnetwork.SecurityRuleProtocolTCP
	access := armnetwork.SecurityRuleAccessAllow
	direction := armnetwork.SecurityRuleDirectionInbound
	srv := networkfake.DefaultSecurityRulesServer{
		NewListPager: func(_, _ string, _ *armnetwork.DefaultSecurityRulesClientListOptions) (resp azfake.PagerResponder[armnetwork.DefaultSecurityRulesClientListResponse]) {
			resp.AddPage(http.StatusOK, armnetwork.DefaultSecurityRulesClientListResponse{
				SecurityRuleListResult: armnetwork.SecurityRuleListResult{
					Value: []*armnetwork.SecurityRule{{
						Name: to.Ptr("AllowVnetInBound"),
						Properties: &armnetwork.SecurityRulePropertiesFormat{
							Protocol: &protocol, Access: &access, Direction: &direction,
							Priority: to.Ptr[int32](65000),
						},
					}},
				},
			}, nil)

			return
		},
	}
	client := newFakeDefaultSecurityRulesClient(t, srv)
	rules, err := GetDefaultNSGRulesWithClient(context.Background(), client, "rg", "nsg")

	require.NoError(t, err)
	require.Len(t, rules, 1)
	assert.Equal(t, "AllowVnetInBound", rules[0].Name)
}

func TestGetCustomNSGRulesWithClient(t *testing.T) {
	t.Parallel()

	protocol := armnetwork.SecurityRuleProtocolUDP
	access := armnetwork.SecurityRuleAccessDeny
	direction := armnetwork.SecurityRuleDirectionOutbound
	srv := networkfake.SecurityRulesServer{
		NewListPager: func(_, _ string, _ *armnetwork.SecurityRulesClientListOptions) (resp azfake.PagerResponder[armnetwork.SecurityRulesClientListResponse]) {
			resp.AddPage(http.StatusOK, armnetwork.SecurityRulesClientListResponse{
				SecurityRuleListResult: armnetwork.SecurityRuleListResult{
					Value: []*armnetwork.SecurityRule{{
						Name: to.Ptr("DenyUDPOut"),
						Properties: &armnetwork.SecurityRulePropertiesFormat{
							Protocol: &protocol, Access: &access, Direction: &direction,
							Priority: to.Ptr[int32](500), DestinationPortRange: to.Ptr("53"),
						},
					}},
				},
			}, nil)

			return
		},
	}
	client := newFakeSecurityRulesClient(t, srv)
	rules, err := GetCustomNSGRulesWithClient(context.Background(), client, "rg", "nsg")

	require.NoError(t, err)
	require.Len(t, rules, 1)
	assert.Equal(t, "DenyUDPOut", rules[0].Name)
}

// Collect security rules with fake servers

func TestCollectDefaultSecurityRules(t *testing.T) {
	t.Parallel()

	t.Run("OnePage", func(t *testing.T) {
		t.Parallel()

		protocol := armnetwork.SecurityRuleProtocolTCP
		access := armnetwork.SecurityRuleAccessAllow
		direction := armnetwork.SecurityRuleDirectionInbound
		srv := networkfake.DefaultSecurityRulesServer{
			NewListPager: func(_, _ string, _ *armnetwork.DefaultSecurityRulesClientListOptions) (resp azfake.PagerResponder[armnetwork.DefaultSecurityRulesClientListResponse]) {
				resp.AddPage(http.StatusOK, armnetwork.DefaultSecurityRulesClientListResponse{
					SecurityRuleListResult: armnetwork.SecurityRuleListResult{
						Value: []*armnetwork.SecurityRule{{
							Name: to.Ptr("AllowVnetInBound"),
							Properties: &armnetwork.SecurityRulePropertiesFormat{
								Protocol: &protocol, Access: &access, Direction: &direction,
								Priority: to.Ptr[int32](65000),
							},
						}},
					},
				}, nil)

				return
			},
		}
		client := newFakeDefaultSecurityRulesClient(t, srv)
		rules, err := collectDefaultSecurityRules(context.Background(), client, "rg", "nsg")

		require.NoError(t, err)
		require.Len(t, rules, 1)
		assert.Equal(t, "AllowVnetInBound", rules[0].Name)
	})

	t.Run("Empty", func(t *testing.T) {
		t.Parallel()

		srv := networkfake.DefaultSecurityRulesServer{
			NewListPager: func(_, _ string, _ *armnetwork.DefaultSecurityRulesClientListOptions) (resp azfake.PagerResponder[armnetwork.DefaultSecurityRulesClientListResponse]) {
				resp.AddPage(http.StatusOK, armnetwork.DefaultSecurityRulesClientListResponse{
					SecurityRuleListResult: armnetwork.SecurityRuleListResult{Value: []*armnetwork.SecurityRule{}},
				}, nil)

				return
			},
		}
		client := newFakeDefaultSecurityRulesClient(t, srv)
		rules, err := collectDefaultSecurityRules(context.Background(), client, "rg", "nsg")

		require.NoError(t, err)
		assert.Empty(t, rules)
	})
}

func TestCollectCustomSecurityRules(t *testing.T) {
	t.Parallel()

	protocol := armnetwork.SecurityRuleProtocolUDP
	access := armnetwork.SecurityRuleAccessDeny
	direction := armnetwork.SecurityRuleDirectionOutbound
	srv := networkfake.SecurityRulesServer{
		NewListPager: func(_, _ string, _ *armnetwork.SecurityRulesClientListOptions) (resp azfake.PagerResponder[armnetwork.SecurityRulesClientListResponse]) {
			resp.AddPage(http.StatusOK, armnetwork.SecurityRulesClientListResponse{
				SecurityRuleListResult: armnetwork.SecurityRuleListResult{
					Value: []*armnetwork.SecurityRule{{
						Name: to.Ptr("DenyUDPOut"),
						Properties: &armnetwork.SecurityRulePropertiesFormat{
							Protocol: &protocol, Access: &access, Direction: &direction,
							Priority: to.Ptr[int32](500), DestinationPortRange: to.Ptr("53"),
						},
					}},
				},
			}, nil)

			return
		},
	}
	client := newFakeSecurityRulesClient(t, srv)
	rules, err := collectCustomSecurityRules(context.Background(), client, "rg", "nsg")

	require.NoError(t, err)
	require.Len(t, rules, 1)
	assert.Equal(t, "DenyUDPOut", rules[0].Name)
}
