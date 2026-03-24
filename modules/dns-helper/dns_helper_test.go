package dns_helper_test //nolint:staticcheck // package name determined by directory

import (
	"testing"
	"time"

	dnshelper "github.com/gruntwork-io/terratest/modules/dns-helper"
	"github.com/gruntwork-io/terratest/modules/retry"
	"github.com/stretchr/testify/require"
)

// These are the current public nameservers for gruntwork.io domain
// They should be updated whenever they change to pass the tests
// relying on the public DNS infrastructure
var publicDomainNameservers = []string{
	"ns-1499.awsdns-59.org",
	"ns-190.awsdns-23.com",
	"ns-1989.awsdns-56.co.uk",
	"ns-853.awsdns-42.net",
}

var testDNSDatabase = dnsDatabase{
	dnshelper.DNSQuery{Type: "A", Name: "a." + testDomain}: dnshelper.DNSAnswers{
		{Type: "A", Value: "2.2.2.2"},
		{Type: "A", Value: "1.1.1.1"},
	},

	dnshelper.DNSQuery{Type: "AAAA", Name: "aaaa." + testDomain}: dnshelper.DNSAnswers{
		{Type: "AAAA", Value: "2001:db8::aaaa"},
	},

	dnshelper.DNSQuery{Type: "CNAME", Name: "terratest." + testDomain}: dnshelper.DNSAnswers{
		{Type: "CNAME", Value: "gruntwork-io.github.io."},
	},

	dnshelper.DNSQuery{Type: "CNAME", Name: "cname1." + testDomain}: dnshelper.DNSAnswers{
		{Type: "CNAME", Value: "cname2." + testDomain + "."},
	},

	dnshelper.DNSQuery{Type: "A", Name: "cname1." + testDomain}: dnshelper.DNSAnswers{
		{Type: "CNAME", Value: "cname2." + testDomain + "."},
		{Type: "CNAME", Value: "cname3." + testDomain + "."},
		{Type: "CNAME", Value: "cname4." + testDomain + "."},
		{Type: "CNAME", Value: "cnamefinal." + testDomain + "."},
		{Type: "A", Value: "1.1.1.1"},
	},

	dnshelper.DNSQuery{Type: "TXT", Name: "txt." + testDomain}: dnshelper.DNSAnswers{
		{Type: "TXT", Value: `"This is a text."`},
	},

	dnshelper.DNSQuery{Type: "MX", Name: testDomain}: dnshelper.DNSAnswers{
		{Type: "MX", Value: "10 mail." + testDomain + "."},
	},
}

// Lookup should succeed in finding the nameservers of the public domain
// Uses system resolver config
func TestOkDNSFindNameservers(t *testing.T) {
	t.Parallel()

	fqdn := "terratest.gruntwork.io"
	expectedNameservers := publicDomainNameservers

	nameservers, err := dnshelper.DNSFindNameserversContextE(t, t.Context(), fqdn, nil)
	require.NoError(t, err)
	require.ElementsMatch(t, nameservers, expectedNameservers)
}

// Lookup should fail because of inexistent domain
// Uses system resolver config
func TestErrorDNSFindNameservers(t *testing.T) {
	t.Parallel()

	fqdn := "this.domain.doesnt.exist"

	nameservers, err := dnshelper.DNSFindNameserversContextE(t, t.Context(), fqdn, nil)
	require.Error(t, err)
	require.Nil(t, nameservers)
}

// Lookup should succeed with answers from just one authoritative nameserver
// Uses system resolver config to lookup a public domain and its public nameservers
func TestOkTerratestDNSLookupAuthoritative(t *testing.T) {
	t.Parallel()

	dnsQuery := dnshelper.DNSQuery{Type: "CNAME", Name: "terratest." + testDomain}
	expected := dnshelper.DNSAnswers{{Type: "CNAME", Value: "gruntwork-io.github.io."}}

	res, err := dnshelper.DNSLookupAuthoritativeContextE(t, t.Context(), dnsQuery, nil)
	require.NoError(t, err)
	require.ElementsMatch(t, res, expected)
}

// ***********************************
// Tests that use local dnsTestServers

// Lookup should succeed with answers from just one authoritative nameserver
func TestOkLocalDNSLookupAuthoritative(t *testing.T) {
	t.Parallel()

	s1, s2 := setupTestDNSServers(t)
	defer shutDownServers(t, s1, s2)

	for dnsQuery, expected := range testDNSDatabase {
		s1.AddEntryToDNSDatabase(dnsQuery, expected)

		res, err := dnshelper.DNSLookupAuthoritativeContextE(t, t.Context(), dnsQuery, []string{s1.Address(), s2.Address()})
		require.NoError(t, err)
		require.ElementsMatch(t, res, expected)
	}
}

// Lookup should fail because of missing answers from all authoritative nameservers
func TestErrorLocalDNSLookupAuthoritative(t *testing.T) {
	t.Parallel()

	s1, s2 := setupTestDNSServers(t)
	defer shutDownServers(t, s1, s2)

	dnsQuery := dnshelper.DNSQuery{Type: "A", Name: "txt." + testDomain}

	_, err := dnshelper.DNSLookupAuthoritativeContextE(t, t.Context(), dnsQuery, []string{s1.Address(), s2.Address()})

	var notFoundErr *dnshelper.NotFoundError

	require.ErrorAs(t, err, &notFoundErr)
}

// Lookup should succeed with consistent answers from all authoritative nameservers
func TestOkLocalDNSLookupAuthoritativeAll(t *testing.T) {
	t.Parallel()

	s1, s2 := setupTestDNSServers(t)
	defer shutDownServers(t, s1, s2)

	for dnsQuery, expected := range testDNSDatabase {
		s1.AddEntryToDNSDatabase(dnsQuery, expected)
		s2.AddEntryToDNSDatabase(dnsQuery, expected)

		res, err := dnshelper.DNSLookupAuthoritativeContextE(t, t.Context(), dnsQuery, []string{s1.Address(), s2.Address()})
		require.NoError(t, err)
		require.ElementsMatch(t, res, expected)
	}
}

// Lookup should fail because of missing answers from all authoritative nameservers
func TestError1DNSLookupAuthoritativeAll(t *testing.T) {
	t.Parallel()

	s1, s2 := setupTestDNSServers(t)
	defer shutDownServers(t, s1, s2)

	dnsQuery := dnshelper.DNSQuery{Type: "A", Name: "txt." + testDomain}

	_, err := dnshelper.DNSLookupAuthoritativeAllContextE(t, t.Context(), dnsQuery, []string{s1.Address(), s2.Address()})

	var notFoundErr *dnshelper.NotFoundError

	require.ErrorAs(t, err, &notFoundErr)
}

// Lookup should fail because of missing answers from one authoritative nameserver
func TestError2DNSLookupAuthoritativeAll(t *testing.T) {
	t.Parallel()

	s1, s2 := setupTestDNSServers(t)
	defer shutDownServers(t, s1, s2)

	dnsQuery := dnshelper.DNSQuery{Type: "A", Name: "a." + testDomain}
	s1.AddEntryToDNSDatabase(dnsQuery, dnshelper.DNSAnswers{{Type: "A", Value: "1.1.1.1"}})

	_, err := dnshelper.DNSLookupAuthoritativeAllContextE(t, t.Context(), dnsQuery, []string{s1.Address(), s2.Address()})

	var notFoundErr *dnshelper.NotFoundError

	require.ErrorAs(t, err, &notFoundErr)
}

// Lookup should fail because of inconsistent answers from authoritative nameservers
func TestError3DNSLookupAuthoritativeAll(t *testing.T) {
	t.Parallel()

	s1, s2 := setupTestDNSServers(t)
	defer shutDownServers(t, s1, s2)

	dnsQuery := dnshelper.DNSQuery{Type: "A", Name: "a." + testDomain}
	s1.AddEntryToDNSDatabase(dnsQuery, dnshelper.DNSAnswers{{Type: "A", Value: "1.1.1.1"}})
	s2.AddEntryToDNSDatabase(dnsQuery, dnshelper.DNSAnswers{{Type: "A", Value: "2.2.2.2"}})

	_, err := dnshelper.DNSLookupAuthoritativeAllContextE(t, t.Context(), dnsQuery, []string{s1.Address(), s2.Address()})

	var inconsistentErr *dnshelper.InconsistentAuthoritativeError

	require.ErrorAs(t, err, &inconsistentErr)
}

// Lookup should fail because of inexistent domain
func TestError4DNSLookupAuthoritativeAll(t *testing.T) {
	t.Parallel()

	s1, s2 := setupTestDNSServers(t)
	defer shutDownServers(t, s1, s2)

	dnsQuery := dnshelper.DNSQuery{Type: "A", Name: "this.domain.doesnt.exist"}

	_, err := dnshelper.DNSLookupAuthoritativeAllContextE(t, t.Context(), dnsQuery, []string{s1.Address(), s2.Address()})

	var nsNotFoundErr *dnshelper.NSNotFoundError

	require.ErrorAs(t, err, &nsNotFoundErr)
}

// First lookups should fail because of missing answers from all authoritative nameservers
// Retry lookups should succeed with answers from just one authoritative nameserver
func TestOkDNSLookupAuthoritativeWithRetry(t *testing.T) {
	t.Parallel()

	s1, s2 := setupTestDNSServersRetry(t)
	defer shutDownServers(t, s1, s2)

	dnsQuery := dnshelper.DNSQuery{Type: "A", Name: "a." + testDomain}
	expectedRes := dnshelper.DNSAnswers{{Type: "A", Value: "1.1.1.1"}, {Type: "A", Value: "2.2.2.2"}}
	s1.AddEntryToDNSDatabaseRetry(dnsQuery, expectedRes)

	res, err := dnshelper.DNSLookupAuthoritativeWithRetryContextE(t, t.Context(), dnsQuery, []string{s1.Address(), s2.Address()}, 5, time.Second)
	require.NoError(t, err)
	require.ElementsMatch(t, res, expectedRes)
}

// First lookups should fail because of missing answers from all authoritative nameservers
// Retry lookups should fail because of missing answers from all authoritative nameservers
func TestErrorDNSLookupAuthoritativeWithRetry(t *testing.T) {
	t.Parallel()

	s1, s2 := setupTestDNSServersRetry(t)
	defer shutDownServers(t, s1, s2)

	dnsQuery := dnshelper.DNSQuery{Type: "A", Name: "txt." + testDomain}

	_, err := dnshelper.DNSLookupAuthoritativeWithRetryContextE(t, t.Context(), dnsQuery, []string{s1.Address(), s2.Address()}, 5, time.Second)
	require.Error(t, err)

	var maxRetriesErr retry.MaxRetriesExceeded

	require.ErrorAs(t, err, &maxRetriesErr)
}

// First lookups should fail because of missing answers from one authoritative nameservers
// Retry lookups should succeed with consistent answers
func TestOkDNSLookupAuthoritativeAllWithRetryNotfound(t *testing.T) {
	t.Parallel()

	s1, s2 := setupTestDNSServersRetry(t)
	defer shutDownServers(t, s1, s2)

	dnsQuery := dnshelper.DNSQuery{Type: "A", Name: "a." + testDomain}
	expectedRes := dnshelper.DNSAnswers{{Type: "A", Value: "1.1.1.1"}, {Type: "A", Value: "2.2.2.2"}}
	s1.AddEntryToDNSDatabase(dnsQuery, expectedRes)
	s1.AddEntryToDNSDatabaseRetry(dnsQuery, expectedRes)
	s2.AddEntryToDNSDatabaseRetry(dnsQuery, expectedRes)

	res, err := dnshelper.DNSLookupAuthoritativeAllWithRetryContextE(t, t.Context(), dnsQuery, []string{s1.Address(), s2.Address()}, 5, time.Second)
	require.NoError(t, err)
	require.ElementsMatch(t, res, expectedRes)
}

// First lookups should fail because of inconsistent answers from authoritative nameservers
// Retry lookups should succeed with consistent answers
func TestOkDNSLookupAuthoritativeAllWithRetryInconsistent(t *testing.T) {
	t.Parallel()

	s1, s2 := setupTestDNSServersRetry(t)
	defer shutDownServers(t, s1, s2)

	dnsQuery := dnshelper.DNSQuery{Type: "A", Name: "a." + testDomain}
	expectedRes := dnshelper.DNSAnswers{{Type: "A", Value: "1.1.1.1"}, {Type: "A", Value: "2.2.2.2"}}
	s1.AddEntryToDNSDatabase(dnsQuery, expectedRes)
	s2.AddEntryToDNSDatabase(dnsQuery, dnshelper.DNSAnswers{{Type: "A", Value: "2.2.2.2"}})
	s1.AddEntryToDNSDatabaseRetry(dnsQuery, expectedRes)
	s2.AddEntryToDNSDatabaseRetry(dnsQuery, expectedRes)

	res, err := dnshelper.DNSLookupAuthoritativeAllWithRetryContextE(t, t.Context(), dnsQuery, []string{s1.Address(), s2.Address()}, 5, time.Second)
	require.NoError(t, err)
	require.ElementsMatch(t, res, expectedRes)
}

// First lookups should fail because of missing answer from one authoritative nameserver
// Retry lookups should fail because of inconsistent answers from authoritative nameservers
func TestErrorDNSLookupAuthoritativeAllWithRetry(t *testing.T) {
	t.Parallel()

	s1, s2 := setupTestDNSServersRetry(t)
	defer shutDownServers(t, s1, s2)

	dnsQuery := dnshelper.DNSQuery{Type: "A", Name: "a." + testDomain}
	s1.AddEntryToDNSDatabase(dnsQuery, dnshelper.DNSAnswers{{Type: "A", Value: "2.2.2.2"}})
	s1.AddEntryToDNSDatabaseRetry(dnsQuery, dnshelper.DNSAnswers{{Type: "A", Value: "1.1.1.1"}, {Type: "A", Value: "2.2.2.2"}})
	s2.AddEntryToDNSDatabaseRetry(dnsQuery, dnshelper.DNSAnswers{{Type: "A", Value: "1.1.1.1"}})

	_, err := dnshelper.DNSLookupAuthoritativeAllWithRetryContextE(t, t.Context(), dnsQuery, []string{s1.Address(), s2.Address()}, 5, time.Second)
	require.Error(t, err)

	var maxRetriesErr retry.MaxRetriesExceeded

	require.ErrorAs(t, err, &maxRetriesErr)
}

// Lookup should succeed with consistent and validated replies
func TestOkDNSLookupAuthoritativeAllWithValidation(t *testing.T) {
	t.Parallel()

	s1, s2 := setupTestDNSServers(t)
	defer shutDownServers(t, s1, s2)

	dnsQuery := dnshelper.DNSQuery{Type: "A", Name: "a." + testDomain}
	expectedRes := dnshelper.DNSAnswers{{Type: "A", Value: "1.1.1.1"}, {Type: "A", Value: "2.2.2.2"}}
	s1.AddEntryToDNSDatabase(dnsQuery, expectedRes)
	s2.AddEntryToDNSDatabase(dnsQuery, expectedRes)

	err := dnshelper.DNSLookupAuthoritativeAllWithValidationContextE(t, t.Context(), dnsQuery, []string{s1.Address(), s2.Address()}, expectedRes)
	require.NoError(t, err)
}

// Lookup should fail because of missing answers from all authoritative nameservers
func TestErrorDNSLookupAuthoritativeAllWithValidation(t *testing.T) {
	t.Parallel()

	s1, s2 := setupTestDNSServers(t)
	defer shutDownServers(t, s1, s2)

	dnsQuery := dnshelper.DNSQuery{Type: "A", Name: "a." + testDomain}
	expectedRes := dnshelper.DNSAnswers{{Type: "A", Value: "1.1.1.1"}, {Type: "A", Value: "2.2.2.2"}}

	err := dnshelper.DNSLookupAuthoritativeAllWithValidationContextE(t, t.Context(), dnsQuery, []string{s1.Address(), s2.Address()}, expectedRes)
	require.Error(t, err)

	var notFoundErr *dnshelper.NotFoundError

	require.ErrorAs(t, err, &notFoundErr)
}

// Lookup should fail because of missing answers from one authoritative nameservers
func TestError2DNSLookupAuthoritativeAllWithValidation(t *testing.T) {
	t.Parallel()

	s1, s2 := setupTestDNSServers(t)
	defer shutDownServers(t, s1, s2)

	dnsQuery := dnshelper.DNSQuery{Type: "A", Name: "a." + testDomain}
	expectedRes := dnshelper.DNSAnswers{{Type: "A", Value: "1.1.1.1"}, {Type: "A", Value: "2.2.2.2"}}
	s1.AddEntryToDNSDatabase(dnsQuery, expectedRes)

	err := dnshelper.DNSLookupAuthoritativeAllWithValidationContextE(t, t.Context(), dnsQuery, []string{s1.Address(), s2.Address()}, expectedRes)
	require.Error(t, err)

	var notFoundErr *dnshelper.NotFoundError

	require.ErrorAs(t, err, &notFoundErr)
}

// Lookup should fail because of inconsistent authoritative replies
func TestError3DNSLookupAuthoritativeAllWithValidation(t *testing.T) {
	t.Parallel()

	s1, s2 := setupTestDNSServers(t)
	defer shutDownServers(t, s1, s2)

	dnsQuery := dnshelper.DNSQuery{Type: "A", Name: "a." + testDomain}
	expectedRes := dnshelper.DNSAnswers{{Type: "A", Value: "1.1.1.1"}, {Type: "A", Value: "2.2.2.2"}}
	s1.AddEntryToDNSDatabase(dnsQuery, expectedRes)
	s2.AddEntryToDNSDatabase(dnsQuery, dnshelper.DNSAnswers{{Type: "A", Value: "2.2.2.2"}})

	err := dnshelper.DNSLookupAuthoritativeAllWithValidationContextE(t, t.Context(), dnsQuery, []string{s1.Address(), s2.Address()}, expectedRes)
	require.Error(t, err)

	var inconsistentErr *dnshelper.InconsistentAuthoritativeError

	require.ErrorAs(t, err, &inconsistentErr)
}

// First lookups should fail because of missing answers from all authoritative nameservers
// Retry lookups should succeed with consistent and validated replies
func TestOkDNSLookupAuthoritativeAllWithValidationRetry(t *testing.T) {
	t.Parallel()

	s1, s2 := setupTestDNSServersRetry(t)
	defer shutDownServers(t, s1, s2)

	dnsQuery := dnshelper.DNSQuery{Type: "A", Name: "a." + testDomain}
	expectedRes := dnshelper.DNSAnswers{{Type: "A", Value: "1.1.1.1"}, {Type: "A", Value: "2.2.2.2"}}
	s1.AddEntryToDNSDatabaseRetry(dnsQuery, expectedRes)
	s2.AddEntryToDNSDatabaseRetry(dnsQuery, expectedRes)

	err := dnshelper.DNSLookupAuthoritativeAllWithValidationRetryContextE(t, t.Context(), dnsQuery, []string{s1.Address(), s2.Address()}, expectedRes, 5, time.Second)
	require.NoError(t, err)
}

// First lookups should fail because of missing answer from one authoritative nameserver
// Retry lookups should succeed with consistent and validated replies
func TestOk2DNSLookupAuthoritativeAllWithValidationRetry(t *testing.T) {
	t.Parallel()

	s1, s2 := setupTestDNSServersRetry(t)
	defer shutDownServers(t, s1, s2)

	dnsQuery := dnshelper.DNSQuery{Type: "A", Name: "a." + testDomain}
	expectedRes := dnshelper.DNSAnswers{{Type: "A", Value: "1.1.1.1"}, {Type: "A", Value: "2.2.2.2"}}
	s1.AddEntryToDNSDatabase(dnsQuery, dnshelper.DNSAnswers{{Type: "A", Value: "2.2.2.2"}})
	s1.AddEntryToDNSDatabaseRetry(dnsQuery, expectedRes)
	s2.AddEntryToDNSDatabaseRetry(dnsQuery, expectedRes)

	err := dnshelper.DNSLookupAuthoritativeAllWithValidationRetryContextE(t, t.Context(), dnsQuery, []string{s1.Address(), s2.Address()}, expectedRes, 5, time.Second)
	require.NoError(t, err)
}

// First lookups should fail because of inconsistent authoritative replies
// Retry lookups should succeed with consistent and validated replies
func TestOk3DNSLookupAuthoritativeAllWithValidationRetry(t *testing.T) {
	t.Parallel()

	s1, s2 := setupTestDNSServersRetry(t)
	defer shutDownServers(t, s1, s2)

	dnsQuery := dnshelper.DNSQuery{Type: "A", Name: "a." + testDomain}
	expectedRes := dnshelper.DNSAnswers{{Type: "A", Value: "1.1.1.1"}, {Type: "A", Value: "2.2.2.2"}}
	s1.AddEntryToDNSDatabase(dnsQuery, expectedRes)
	s2.AddEntryToDNSDatabase(dnsQuery, dnshelper.DNSAnswers{{Type: "A", Value: "2.2.2.2"}})
	s1.AddEntryToDNSDatabaseRetry(dnsQuery, expectedRes)
	s2.AddEntryToDNSDatabaseRetry(dnsQuery, expectedRes)

	err := dnshelper.DNSLookupAuthoritativeAllWithValidationRetryContextE(t, t.Context(), dnsQuery, []string{s1.Address(), s2.Address()}, expectedRes, 5, time.Second)
	require.NoError(t, err)
}

// First lookups should fail because of inconsistent authoritative replies
// Retry lookups should fail also because of inconsistent authoritative replies
func TestErrorDNSLookupAuthoritativeAllWithValidationRetry(t *testing.T) {
	t.Parallel()

	s1, s2 := setupTestDNSServersRetry(t)
	defer shutDownServers(t, s1, s2)

	dnsQuery := dnshelper.DNSQuery{Type: "A", Name: "a." + testDomain}
	expectedRes := dnshelper.DNSAnswers{{Type: "A", Value: "1.1.1.1"}, {Type: "A", Value: "2.2.2.2"}}
	s1.AddEntryToDNSDatabase(dnsQuery, expectedRes)
	s2.AddEntryToDNSDatabase(dnsQuery, dnshelper.DNSAnswers{{Type: "A", Value: "2.2.2.2"}})
	s1.AddEntryToDNSDatabaseRetry(dnsQuery, expectedRes)
	s2.AddEntryToDNSDatabaseRetry(dnsQuery, dnshelper.DNSAnswers{{Type: "A", Value: "2.2.2.2"}})

	err := dnshelper.DNSLookupAuthoritativeAllWithValidationRetryContextE(t, t.Context(), dnsQuery, []string{s1.Address(), s2.Address()}, expectedRes, 5, time.Second)

	var maxRetriesErr retry.MaxRetriesExceeded

	require.ErrorAs(t, err, &maxRetriesErr)
}

func shutDownServers(t *testing.T, s1, s2 *dnsTestServer) {
	t.Helper()

	err := s1.Server.Shutdown()
	require.NoError(t, err)

	err = s2.Server.Shutdown()
	require.NoError(t, err)
}
