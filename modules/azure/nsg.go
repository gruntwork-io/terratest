package azure

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-09-01/network"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// NsgRuleSummaryList holds a colleciton of NsgRuleSummary rules
type NsgRuleSummaryList struct {
	SummarizedRules []NsgRuleSummary
}

// NsgRuleSummary is a string-based (non-pointer) summary of an NSG rule with several helper methods attached
// to help with verification of rule configuratoin.
type NsgRuleSummary struct {
	Name                     string
	Description              string
	Protocol                 string
	SourcePortRange          string
	DestinationPortRange     string
	SourceAddressPrefix      string
	DestinationAddressPrefix string
	Access                   string
	Priority                 int32
	Direction                string
}

// GetDefaultNsgRulesClient returns a rules client which can be used to read the list of *default* security rules
// defined on an network security group. Note that the "default" rules are those provided implicitly
// by the Azure platform.
// This function would fail the test if there is an error.
func GetDefaultNsgRulesClient(t *testing.T, subscriptionID string) network.DefaultSecurityRulesClient {
	client, err := GetDefaultNsgRulesClientE(subscriptionID)
	require.NoError(t, err)
	return client
}

// GetDefaultNsgRulesClientE returns a rules client which can be used to read the list of *default* security rules
// defined on an network security group. Note that the "default" rules are those provided implicitly
// by the Azure platform.
func GetDefaultNsgRulesClientE(subscriptionID string) (network.DefaultSecurityRulesClient, error) {
	// Get new default client from client factory
	nsgClient, err := CreateNsgDefaultRulesClientE(subscriptionID)
	if err != nil {
		return network.DefaultSecurityRulesClient{}, err
	}

	// Get an authorizer
	auth, err := NewAuthorizer()
	if err != nil {
		return network.DefaultSecurityRulesClient{}, err
	}

	nsgClient.Authorizer = *auth
	return *nsgClient, nil
}

// GetCustomNsgRulesClient returns a rules client which can be used to read the list of *custom* security rules
// defined on an network security group. Note that the "custom" rules are those defined by
// end users.
// This function would fail the test if there is an error.
func GetCustomNsgRulesClient(t *testing.T, subscriptionID string) network.SecurityRulesClient {
	client, err := GetCustomNsgRulesClientE(subscriptionID)
	require.NoError(t, err)
	return client
}

// GetCustomNsgRulesClientE returns a rules client which can be used to read the list of *custom* security rules
// defined on an network security group. Note that the "custom" rules are those defined by
// end users.
func GetCustomNsgRulesClientE(subscriptionID string) (network.SecurityRulesClient, error) {
	// Get new custom rules client from client factory
	nsgClient, err := CreateNsgCustomRulesClientE(subscriptionID)
	if err != nil {
		return network.SecurityRulesClient{}, err
	}

	// Get an authorizer
	auth, err := NewAuthorizer()
	if err != nil {
		return network.SecurityRulesClient{}, err
	}

	nsgClient.Authorizer = *auth
	return *nsgClient, nil
}

// GetAllNSGRules returns an NsgRuleSummaryList instance containing the combined "default" and "custom" rules from a network
// security group.
// This function would fail the test if there is an error.
func GetAllNSGRules(t *testing.T, resourceGroupName, nsgName, subscriptionID string) NsgRuleSummaryList {
	results, err := GetAllNSGRulesE(resourceGroupName, nsgName, subscriptionID)
	require.NoError(t, err)
	return results
}

// GetAllNSGRulesE returns an NsgRuleSummaryList instance containing the combined "default" and "custom" rules from a network
// security group.
func GetAllNSGRulesE(resourceGroupName, nsgName, subscriptionID string) (NsgRuleSummaryList, error) {
	defaultRulesClient, err := GetDefaultNsgRulesClientE(subscriptionID)
	if err != nil {
		return NsgRuleSummaryList{}, err
	}

	// Get a client instance
	customRulesClient, err := GetCustomNsgRulesClientE(subscriptionID)
	if err != nil {
		return NsgRuleSummaryList{}, err
	}

	// Read all default (platform) rules.
	defaultRuleList, err := defaultRulesClient.ListComplete(context.Background(), resourceGroupName, nsgName)
	if err != nil {
		return NsgRuleSummaryList{}, err
	}

	// Read any custom (user provided) rules
	customRuleList, err := customRulesClient.ListComplete(context.Background(), resourceGroupName, nsgName)
	if err != nil {
		return NsgRuleSummaryList{}, err
	}

	// Convert the default list to our summary type
	boundDefaultRules, err := bindRuleList(defaultRuleList)
	if err != nil {
		return NsgRuleSummaryList{}, err
	}

	// Convert the custom list to our summary type
	boundCustomRules, err := bindRuleList(customRuleList)
	if err != nil {
		return NsgRuleSummaryList{}, err
	}

	// Join the summarized lists and wrap in NsgRuleSummaryList struct
	allRules := append(boundDefaultRules, boundCustomRules...)
	ruleList := NsgRuleSummaryList{}
	ruleList.SummarizedRules = allRules
	return ruleList, nil
}

// bindRuleList takes a raw list of security rules from the SDK and converts them into a string-based
// summary struct.
func bindRuleList(source network.SecurityRuleListResultIterator) ([]NsgRuleSummary, error) {
	rules := make([]NsgRuleSummary, 0)
	for source.NotDone() {
		v := source.Value()
		rules = append(rules, convertToNsgRuleSummary(v.Name, v.SecurityRulePropertiesFormat))
		err := source.NextWithContext(context.Background())
		if err != nil {
			return []NsgRuleSummary{}, err
		}
	}
	return rules, nil
}

// convertToNsgRuleSummary converts the raw SDK security rule type into a summarized struct, flattening the
// rules properties and name into a single, string-based struct.
func convertToNsgRuleSummary(name *string, rule *network.SecurityRulePropertiesFormat) NsgRuleSummary {
	summary := NsgRuleSummary{}
	summary.Description = safePtrToString(rule.Description)
	summary.Name = safePtrToString(name)
	summary.Protocol = string(rule.Protocol)
	summary.SourcePortRange = safePtrToString(rule.SourcePortRange)
	summary.DestinationPortRange = safePtrToString(rule.DestinationPortRange)
	summary.SourceAddressPrefix = safePtrToString(rule.SourceAddressPrefix)
	summary.DestinationAddressPrefix = safePtrToString(rule.DestinationAddressPrefix)
	summary.Access = string(rule.Access)
	summary.Priority = safePtrToInt32(rule.Priority)
	summary.Direction = string(rule.Direction)
	return summary
}

// FindRuleByName looks for a matching rule by name within the current collection of rules.
func (summarizedRules *NsgRuleSummaryList) FindRuleByName(name string) NsgRuleSummary {
	for _, r := range summarizedRules.SummarizedRules {
		if r.Name == name {
			return r
		}
	}

	return NsgRuleSummary{}
}

// AllowsDestinationPort checks to see if the rule allows a specific destination port. This is helpful when verifying
// that a given rule is configured properly for a given port.
func (summarizedRule *NsgRuleSummary) AllowsDestinationPort(t *testing.T, port string) bool {
	allowed, err := portRangeAllowsPort(summarizedRule.DestinationPortRange, port)
	assert.NoError(t, err)
	return allowed && (summarizedRule.Access == "Allow")
}

// AllowsSourcePort checks to see if the rule allows a specific source port. This is helpful when verifying
// that a given rule is configured properly for a given port.
func (summarizedRule *NsgRuleSummary) AllowsSourcePort(t *testing.T, port string) bool {
	allowed, err := portRangeAllowsPort(summarizedRule.SourcePortRange, port)
	assert.NoError(t, err)
	return allowed && (summarizedRule.Access == "Allow")
}

// portRangeAllowsPort is the internal impelmentation of AllowsSourcePort and AllowsDestinationPort.
func portRangeAllowsPort(portRange string, port string) (bool, error) {
	if portRange == "*" {
		return true, nil
	}

	// Decode the provided port range
	low, high, parseErr := parsePortRangeString(portRange)
	if parseErr != nil {
		return false, parseErr
	}

	// Decode user-provided port
	portAsInt, parseErr := strconv.ParseInt(port, 10, 16)
	if (parseErr != nil) && (port != "*") {
		return false, parseErr
	}

	// If the user wants to check "all", make sure we parsed input range to include all ports.
	if (port == "*") && (low == 0) && (high == 65535) {
		return true, nil
	}

	// Evaluate and return
	return ((uint16(portAsInt) >= low) && (uint16(portAsInt) <= high)), nil
}

// parsePortRangeString decodes a range string ("2-100") or a single digit ("22") and returns
// a tuple in [low, hi] form. Note that if a single digit is supplied, both members of the
// return tuple will be the same value (e.g., "22" returns (22, 22))
func parsePortRangeString(rangeString string) (uint16, uint16, error) {
	// An asterisk means all ports
	if rangeString == "*" {
		return uint16(0), uint16(65535), nil
	}

	// Check for range string that contains hyphen separator
	if !strings.Contains(rangeString, "-") {
		val, parseErr := strconv.ParseInt(rangeString, 10, 16)
		if parseErr != nil {
			return 0, 0, parseErr
		}
		return uint16(val), uint16(val), nil
	}

	// Split the rang into parts and validate
	parts := strings.Split(rangeString, "-")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("Invalid port range specified; must be of the format '{low port}-{high port}'")
	}

	// Assume the low port is listed first; parse it
	lowVal, parseErr := strconv.ParseInt(parts[0], 10, 16)
	if parseErr != nil {
		return 0, 0, parseErr
	}

	// Assume the hi port is listed first; parse it
	highVal, parseErr := strconv.ParseInt(parts[1], 10, 16)
	if parseErr != nil {
		return 0, 0, parseErr
	}

	// Normalize ordering in the case that low and hi were reversed.
	// This should _never_ happen, as the Azure API's won't allow it, but
	// we shouldn't fail if it's the case.
	if lowVal > highVal {
		temp := lowVal
		lowVal = highVal
		highVal = temp
	}

	// Return values
	return uint16(lowVal), uint16(highVal), nil
}
