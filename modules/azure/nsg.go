package azure

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v6"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	// allPortsMax is the maximum port number, representing "all ports" when paired with 0.
	allPortsMax = 65535

	// portRangeParts is the expected number of parts when splitting a port range string on "-".
	portRangeParts = 2
)

// NsgRuleSummaryList holds a collection of NsgRuleSummary rules.
type NsgRuleSummaryList struct {
	SummarizedRules []NsgRuleSummary
}

// NsgRuleSummary is a string-based (non-pointer) summary of an NSG rule with several helper methods attached
// to help with verification of rule configuration.
type NsgRuleSummary struct {
	SourceAddressPrefix        string
	Direction                  string
	SourcePortRange            string
	Access                     string
	Name                       string
	Description                string
	DestinationPortRange       string
	Protocol                   string
	DestinationAddressPrefix   string
	SourcePortRanges           []string
	DestinationPortRanges      []string
	DestinationAddressPrefixes []string
	SourceAdresssPrefixes      []string
	Priority                   int32
}

// GetDefaultNsgRulesClient returns a rules client which can be used to read the list of *default* security rules
// defined on a network security group. Note that the "default" rules are those provided implicitly
// by the Azure platform.
// This function would fail the test if there is an error.
func GetDefaultNsgRulesClient(t testing.TestingT, subscriptionID string) *armnetwork.DefaultSecurityRulesClient {
	t.Helper()

	client, err := GetDefaultNsgRulesClientE(subscriptionID)
	require.NoError(t, err)

	return client
}

// GetDefaultNsgRulesClientE returns a rules client which can be used to read the list of *default* security rules
// defined on a network security group. Note that the "default" rules are those provided implicitly
// by the Azure platform.
func GetDefaultNsgRulesClientE(subscriptionID string) (*armnetwork.DefaultSecurityRulesClient, error) {
	return CreateNsgDefaultRulesClientE(subscriptionID)
}

// GetCustomNsgRulesClient returns a rules client which can be used to read the list of *custom* security rules
// defined on a network security group. Note that the "custom" rules are those defined by
// end users.
// This function would fail the test if there is an error.
func GetCustomNsgRulesClient(t testing.TestingT, subscriptionID string) *armnetwork.SecurityRulesClient {
	t.Helper()

	client, err := GetCustomNsgRulesClientE(subscriptionID)
	require.NoError(t, err)

	return client
}

// GetCustomNsgRulesClientE returns a rules client which can be used to read the list of *custom* security rules
// defined on a network security group. Note that the "custom" rules are those defined by
// end users.
func GetCustomNsgRulesClientE(subscriptionID string) (*armnetwork.SecurityRulesClient, error) {
	return CreateNsgCustomRulesClientE(subscriptionID)
}

// GetAllNSGRulesContext returns an NsgRuleSummaryList instance containing the combined "default" and "custom" rules
// from a network security group. The ctx parameter supports cancellation and timeouts.
// This function would fail the test if there is an error.
func GetAllNSGRulesContext(t testing.TestingT, ctx context.Context, resourceGroupName, nsgName, subscriptionID string) NsgRuleSummaryList {
	t.Helper()

	results, err := GetAllNSGRulesContextE(ctx, resourceGroupName, nsgName, subscriptionID)
	require.NoError(t, err)

	return results
}

// GetAllNSGRulesContextE returns an NsgRuleSummaryList instance containing the combined "default" and "custom" rules
// from a network security group. The ctx parameter supports cancellation and timeouts.
func GetAllNSGRulesContextE(ctx context.Context, resourceGroupName, nsgName, subscriptionID string) (NsgRuleSummaryList, error) {
	defaultRulesClient, err := GetDefaultNsgRulesClientE(subscriptionID)
	if err != nil {
		return NsgRuleSummaryList{}, err
	}

	// Get a client instance
	customRulesClient, err := GetCustomNsgRulesClientE(subscriptionID)
	if err != nil {
		return NsgRuleSummaryList{}, err
	}

	// Read all default (platform) rules using the pager pattern.
	boundDefaultRules, err := collectDefaultSecurityRules(ctx, defaultRulesClient, resourceGroupName, nsgName)
	if err != nil {
		return NsgRuleSummaryList{}, err
	}

	// Read any custom (user provided) rules using the pager pattern.
	boundCustomRules, err := collectCustomSecurityRules(ctx, customRulesClient, resourceGroupName, nsgName)
	if err != nil {
		return NsgRuleSummaryList{}, err
	}

	// Join the summarized lists and wrap in NsgRuleSummaryList struct
	allRules := make([]NsgRuleSummary, 0, len(boundDefaultRules)+len(boundCustomRules))
	allRules = append(allRules, boundDefaultRules...)
	allRules = append(allRules, boundCustomRules...)

	return NsgRuleSummaryList{SummarizedRules: allRules}, nil
}

// collectDefaultSecurityRules uses the pager pattern to iterate over all default security rules and
// convert them into NsgRuleSummary instances.
func collectDefaultSecurityRules(ctx context.Context, client *armnetwork.DefaultSecurityRulesClient, resourceGroupName, nsgName string) ([]NsgRuleSummary, error) {
	var rules []NsgRuleSummary

	pager := client.NewListPager(resourceGroupName, nsgName, nil)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return []NsgRuleSummary{}, err
		}

		for _, rule := range page.Value {
			rules = append(rules, convertToNsgRuleSummary(rule.Name, rule.Properties))
		}
	}

	return rules, nil
}

// collectCustomSecurityRules uses the pager pattern to iterate over all custom security rules and
// convert them into NsgRuleSummary instances.
func collectCustomSecurityRules(ctx context.Context, client *armnetwork.SecurityRulesClient, resourceGroupName, nsgName string) ([]NsgRuleSummary, error) {
	var rules []NsgRuleSummary

	pager := client.NewListPager(resourceGroupName, nsgName, nil)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return []NsgRuleSummary{}, err
		}

		for _, rule := range page.Value {
			rules = append(rules, convertToNsgRuleSummary(rule.Name, rule.Properties))
		}
	}

	return rules, nil
}

// convertToNsgRuleSummary converts the raw SDK security rule type into a summarized struct, flattening the
// rules properties and name into a single, string-based struct.
func convertToNsgRuleSummary(name *string, rule *armnetwork.SecurityRulePropertiesFormat) NsgRuleSummary {
	summary := NsgRuleSummary{Name: safePtrToString(name)}

	if rule == nil {
		return summary
	}

	summary.Description = safePtrToString(rule.Description)
	summary.Protocol = safeDerefString((*string)(rule.Protocol))
	summary.SourcePortRange = safePtrToString(rule.SourcePortRange)
	summary.SourcePortRanges = safePtrToList(rule.SourcePortRanges)
	summary.DestinationPortRange = safePtrToString(rule.DestinationPortRange)
	summary.DestinationPortRanges = safePtrToList(rule.DestinationPortRanges)
	summary.SourceAddressPrefix = safePtrToString(rule.SourceAddressPrefix)
	summary.SourceAdresssPrefixes = safePtrToList(rule.SourceAddressPrefixes)
	summary.DestinationAddressPrefix = safePtrToString(rule.DestinationAddressPrefix)
	summary.DestinationAddressPrefixes = safePtrToList(rule.DestinationAddressPrefixes)
	summary.Access = safeDerefString((*string)(rule.Access))
	summary.Priority = safePtrToInt32(rule.Priority)
	summary.Direction = safeDerefString((*string)(rule.Direction))

	return summary
}

// safeDerefString safely dereferences a *string, returning "" if nil.
func safeDerefString(s *string) string {
	if s == nil {
		return ""
	}

	return *s
}

// FindRuleByName looks for a matching rule by name within the current collection of rules.
func (summarizedRules *NsgRuleSummaryList) FindRuleByName(name string) NsgRuleSummary {
	for i := range summarizedRules.SummarizedRules {
		if summarizedRules.SummarizedRules[i].Name == name {
			return summarizedRules.SummarizedRules[i]
		}
	}

	return NsgRuleSummary{}
}

// AllowsDestinationPort checks to see if the rule allows a specific destination port. This is helpful when verifying
// that a given rule is configured properly for a given port.
func (summarizedRule *NsgRuleSummary) AllowsDestinationPort(t testing.TestingT, port string) bool {
	t.Helper()

	allowed, err := portRangeAllowsPort(summarizedRule.DestinationPortRange, port)
	assert.NoError(t, err)

	return allowed && (summarizedRule.Access == "Allow")
}

// AllowsSourcePort checks to see if the rule allows a specific source port. This is helpful when verifying
// that a given rule is configured properly for a given port.
func (summarizedRule *NsgRuleSummary) AllowsSourcePort(t testing.TestingT, port string) bool {
	t.Helper()

	allowed, err := portRangeAllowsPort(summarizedRule.SourcePortRange, port)
	assert.NoError(t, err)

	return allowed && (summarizedRule.Access == "Allow")
}

// portRangeAllowsPort is the internal implementation of AllowsSourcePort and AllowsDestinationPort.
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
	if (port == "*") && (low == 0) && (high == allPortsMax) {
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
		return uint16(0), uint16(allPortsMax), nil
	}

	// Check for range string that contains hyphen separator
	if !strings.Contains(rangeString, "-") {
		val, parseErr := strconv.ParseInt(rangeString, 10, 16)
		if parseErr != nil {
			return 0, 0, parseErr
		}

		return uint16(val), uint16(val), nil
	}

	// Split the range into parts and validate
	parts := strings.Split(rangeString, "-")

	if len(parts) != portRangeParts {
		return 0, 0, errors.New("invalid port range specified; must be of the format '{low port}-{high port}'")
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
		lowVal, highVal = highVal, lowVal
	}

	// Return values
	return uint16(lowVal), uint16(highVal), nil
}
