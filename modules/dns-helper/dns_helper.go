// Package dns_helper contains helpers to interact with the Domain Name System.
package dns_helper //nolint:staticcheck // package name determined by directory

import (
	"context"
	"fmt"
	"net"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/retry"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/miekg/dns"
	"github.com/stretchr/testify/require"
)

// DNSQuery represents a DNS query with a record type and name.
type DNSQuery struct {
	Type, Name string
}

// DNSAnswer represents a single DNS answer with a record type and value.
type DNSAnswer struct {
	Type, Value string
}

// String returns a human-readable representation of the DNS answer.
func (a DNSAnswer) String() string {
	return fmt.Sprintf("%s %s", a.Type, a.Value)
}

// DNSAnswers is a collection of DNS answers.
type DNSAnswers []DNSAnswer

// Sort sorts the answers by type and value.
func (a DNSAnswers) Sort() {
	sort.Slice(a, func(i, j int) bool {
		return a[i].Type < a[j].Type || a[i].Value < a[j].Value
	})
}

// DNSLookupContext sends a DNS query for the specified record and type using the given resolvers.
// Fails on any error.
// Supported record types: A, AAAA, CNAME, MX, NS, TXT
func DNSLookupContext(t testing.TestingT, ctx context.Context, query DNSQuery, resolvers []string) DNSAnswers {
	res, err := DNSLookupContextE(t, ctx, query, resolvers)
	require.NoError(t, err)

	return res
}

// DNSLookupContextE sends a DNS query for the specified record and type using the given resolvers.
// Returns [QueryTypeError] when record type is not supported.
// Returns [NoResolversError] when no resolvers are provided.
// Returns any underlying error.
// Supported record types: A, AAAA, CNAME, MX, NS, TXT
func DNSLookupContextE(t testing.TestingT, ctx context.Context, query DNSQuery, resolvers []string) (DNSAnswers, error) {
	if len(resolvers) == 0 {
		return nil, &NoResolversError{}
	}

	var dnsAnswers DNSAnswers

	var err error

	for _, resolver := range resolvers {
		dnsAnswers, err = dnsLookupContext(t, ctx, query, resolver)
		if err == nil {
			return dnsAnswers, nil
		}
	}

	return nil, err
}

// DNSLookup sends a DNS query for the specified record and type using the given resolvers.
// Fails on any error.
// Supported record types: A, AAAA, CNAME, MX, NS, TXT
//
// Deprecated: Use [DNSLookupContext] instead.
func DNSLookup(t testing.TestingT, query DNSQuery, resolvers []string) DNSAnswers {
	return DNSLookupContext(t, context.Background(), query, resolvers)
}

// DNSLookupE sends a DNS query for the specified record and type using the given resolvers.
// Returns [QueryTypeError] when record type is not supported.
// Returns [NoResolversError] when no resolvers are provided.
// Returns any underlying error.
// Supported record types: A, AAAA, CNAME, MX, NS, TXT
//
// Deprecated: Use [DNSLookupContextE] instead.
func DNSLookupE(t testing.TestingT, query DNSQuery, resolvers []string) (DNSAnswers, error) {
	return DNSLookupContextE(t, context.Background(), query, resolvers)
}

// DNSFindNameserversContext tries to find the NS record for the given FQDN, iterating down the domain hierarchy
// until it finds the NS records and returns them. Fails if there's any error or no NS record is found up to the apex domain.
func DNSFindNameserversContext(t testing.TestingT, ctx context.Context, fqdn string, resolvers []string) []string {
	nameservers, err := DNSFindNameserversContextE(t, ctx, fqdn, resolvers)
	require.NoError(t, err)

	return nameservers
}

// DNSFindNameserversContextE tries to find the NS record for the given FQDN, iterating down the domain hierarchy
// until it finds the NS records and returns them. Returns [NSNotFoundError] if the apex domain is reached with no result.
func DNSFindNameserversContextE(t testing.TestingT, ctx context.Context, fqdn string, resolvers []string) ([]string, error) {
	var lookupFunc func(domain string) ([]string, error)

	if resolvers == nil {
		resolver := &net.Resolver{}

		lookupFunc = func(domain string) ([]string, error) {
			res, err := resolver.LookupNS(ctx, domain)
			nameservers := make([]string, 0, len(res))

			for _, ns := range res {
				nameservers = append(nameservers, ns.Host)
			}

			return nameservers, err
		}
	} else {
		lookupFunc = func(domain string) ([]string, error) {
			res, err := DNSLookupContextE(t, ctx, DNSQuery{Type: "NS", Name: domain}, resolvers)
			nameservers := make([]string, 0, len(res))

			for _, r := range res {
				if r.Type == "NS" {
					nameservers = append(nameservers, r.Value)
				}
			}

			return nameservers, err
		}
	}

	parts := strings.Split(fqdn, ".")

	var domain string

	for i := range parts[:len(parts)-1] {
		domain = strings.Join(parts[i:], ".")

		res, err := lookupFunc(domain)

		if len(res) > 0 {
			nameservers := make([]string, 0, len(res))

			for _, ns := range res {
				nameservers = append(nameservers, strings.TrimSuffix(ns, "."))
			}

			logger.Default.Logf(t, "FQDN %s belongs to domain %s, found NS record: %s", fqdn, domain, nameservers)

			return nameservers, nil
		}

		if err != nil {
			logger.Default.Logf(t, "%s", err.Error())
		}
	}

	return nil, &NSNotFoundError{FQDN: fqdn, Nameserver: domain}
}

// DNSFindNameservers tries to find the NS record for the given FQDN, iterating down the domain hierarchy
// until it finds the NS records and returns them. Fails if there's any error or no NS record is found up to the apex domain.
//
// Deprecated: Use [DNSFindNameserversContext] instead.
func DNSFindNameservers(t testing.TestingT, fqdn string, resolvers []string) []string {
	return DNSFindNameserversContext(t, context.Background(), fqdn, resolvers)
}

// DNSFindNameserversE tries to find the NS record for the given FQDN, iterating down the domain hierarchy
// until it finds the NS records and returns them. Returns [NSNotFoundError] if the apex domain is reached with no result.
//
// Deprecated: Use [DNSFindNameserversContextE] instead.
func DNSFindNameserversE(t testing.TestingT, fqdn string, resolvers []string) ([]string, error) {
	return DNSFindNameserversContextE(t, context.Background(), fqdn, resolvers)
}

// DNSLookupAuthoritativeContext gets authoritative answers for the specified record and type.
// If resolvers are defined, uses them instead of the default system ones to find the authoritative nameservers.
// Fails on any error from [DNSLookupAuthoritativeContextE].
func DNSLookupAuthoritativeContext(t testing.TestingT, ctx context.Context, query DNSQuery, resolvers []string) DNSAnswers {
	res, err := DNSLookupAuthoritativeContextE(t, ctx, query, resolvers)
	require.NoError(t, err)

	return res
}

// DNSLookupAuthoritativeContextE gets authoritative answers for the specified record and type.
// If resolvers are defined, uses them instead of the default system ones to find the authoritative nameservers.
// Returns [NotFoundError] when no answer found in any authoritative nameserver.
// Returns any underlying error from individual lookups.
func DNSLookupAuthoritativeContextE(t testing.TestingT, ctx context.Context, query DNSQuery, resolvers []string) (DNSAnswers, error) {
	nameservers, err := DNSFindNameserversContextE(t, ctx, query.Name, resolvers)
	if err != nil {
		return nil, err
	}

	return DNSLookupContextE(t, ctx, query, nameservers)
}

// DNSLookupAuthoritative gets authoritative answers for the specified record and type.
// If resolvers are defined, uses them instead of the default system ones to find the authoritative nameservers.
// Fails on any error from [DNSLookupAuthoritativeContextE].
//
// Deprecated: Use [DNSLookupAuthoritativeContext] instead.
func DNSLookupAuthoritative(t testing.TestingT, query DNSQuery, resolvers []string) DNSAnswers {
	return DNSLookupAuthoritativeContext(t, context.Background(), query, resolvers)
}

// DNSLookupAuthoritativeE gets authoritative answers for the specified record and type.
// If resolvers are defined, uses them instead of the default system ones to find the authoritative nameservers.
// Returns [NotFoundError] when no answer found in any authoritative nameserver.
// Returns any underlying error from individual lookups.
//
// Deprecated: Use [DNSLookupAuthoritativeContextE] instead.
func DNSLookupAuthoritativeE(t testing.TestingT, query DNSQuery, resolvers []string) (DNSAnswers, error) {
	return DNSLookupAuthoritativeContextE(t, context.Background(), query, resolvers)
}

// DNSLookupAuthoritativeWithRetryContext repeatedly gets authoritative answers for the specified record and type
// until ANY of the authoritative nameservers found replies with a non-empty answer,
// or until max retries has been exceeded.
// If resolvers are defined, uses them instead of the default system ones to find the authoritative nameservers.
// Fails on any error from [DNSLookupAuthoritativeWithRetryContextE].
func DNSLookupAuthoritativeWithRetryContext(t testing.TestingT, ctx context.Context, query DNSQuery, resolvers []string, maxRetries int, sleepBetweenRetries time.Duration) DNSAnswers {
	res, err := DNSLookupAuthoritativeWithRetryContextE(t, ctx, query, resolvers, maxRetries, sleepBetweenRetries)
	require.NoError(t, err)

	return res
}

// DNSLookupAuthoritativeWithRetryContextE repeatedly gets authoritative answers for the specified record and type
// until ANY of the authoritative nameservers found replies with a non-empty answer,
// or until max retries has been exceeded.
// If resolvers are defined, uses them instead of the default system ones to find the authoritative nameservers.
func DNSLookupAuthoritativeWithRetryContextE(t testing.TestingT, ctx context.Context, query DNSQuery, resolvers []string, maxRetries int, sleepBetweenRetries time.Duration) (DNSAnswers, error) {
	res, err := retry.DoWithRetryInterfaceContextE(
		t, ctx, fmt.Sprintf("DNSLookupAuthoritativeContextE %s record for %s using authoritative nameservers", query.Type, query.Name),
		maxRetries, sleepBetweenRetries,
		func() (any, error) {
			return DNSLookupAuthoritativeContextE(t, ctx, query, resolvers)
		})

	return res.(DNSAnswers), err
}

// DNSLookupAuthoritativeWithRetry repeatedly gets authoritative answers for the specified record and type
// until ANY of the authoritative nameservers found replies with a non-empty answer,
// or until max retries has been exceeded.
// If resolvers are defined, uses them instead of the default system ones to find the authoritative nameservers.
// Fails on any error from [DNSLookupAuthoritativeWithRetryContextE].
//
// Deprecated: Use [DNSLookupAuthoritativeWithRetryContext] instead.
func DNSLookupAuthoritativeWithRetry(t testing.TestingT, query DNSQuery, resolvers []string, maxRetries int, sleepBetweenRetries time.Duration) DNSAnswers {
	return DNSLookupAuthoritativeWithRetryContext(t, context.Background(), query, resolvers, maxRetries, sleepBetweenRetries)
}

// DNSLookupAuthoritativeWithRetryE repeatedly gets authoritative answers for the specified record and type
// until ANY of the authoritative nameservers found replies with a non-empty answer,
// or until max retries has been exceeded.
// If resolvers are defined, uses them instead of the default system ones to find the authoritative nameservers.
//
// Deprecated: Use [DNSLookupAuthoritativeWithRetryContextE] instead.
func DNSLookupAuthoritativeWithRetryE(t testing.TestingT, query DNSQuery, resolvers []string, maxRetries int, sleepBetweenRetries time.Duration) (DNSAnswers, error) {
	return DNSLookupAuthoritativeWithRetryContextE(t, context.Background(), query, resolvers, maxRetries, sleepBetweenRetries)
}

// DNSLookupAuthoritativeAllContext gets authoritative answers for the specified record and type.
// All the authoritative nameservers found must give the same answers.
// If resolvers are defined, uses them instead of the default system ones to find the authoritative nameservers.
// Fails on any error from [DNSLookupAuthoritativeAllContextE].
func DNSLookupAuthoritativeAllContext(t testing.TestingT, ctx context.Context, query DNSQuery, resolvers []string) DNSAnswers {
	res, err := DNSLookupAuthoritativeAllContextE(t, ctx, query, resolvers)
	require.NoError(t, err)

	return res
}

// DNSLookupAuthoritativeAllContextE gets authoritative answers for the specified record and type.
// All the authoritative nameservers found must give the same answers.
// If resolvers are defined, uses them instead of the default system ones to find the authoritative nameservers.
// Returns [InconsistentAuthoritativeError] when any authoritative nameserver gives a different answer.
// Returns any underlying error.
func DNSLookupAuthoritativeAllContextE(t testing.TestingT, ctx context.Context, query DNSQuery, resolvers []string) (DNSAnswers, error) {
	nameservers, err := DNSFindNameserversContextE(t, ctx, query.Name, resolvers)
	if err != nil {
		return nil, err
	}

	var answers DNSAnswers

	for _, ns := range nameservers {
		res, err := DNSLookupContextE(t, ctx, query, []string{ns})
		if err != nil {
			return nil, err
		}

		if len(answers) > 0 {
			if !reflect.DeepEqual(answers, res) {
				return nil, &InconsistentAuthoritativeError{Query: query, Answers: res, Nameserver: ns, PreviousAnswers: answers}
			}
		} else {
			answers = res
		}
	}

	return answers, nil
}

// DNSLookupAuthoritativeAll gets authoritative answers for the specified record and type.
// All the authoritative nameservers found must give the same answers.
// If resolvers are defined, uses them instead of the default system ones to find the authoritative nameservers.
// Fails on any error from [DNSLookupAuthoritativeAllContextE].
//
// Deprecated: Use [DNSLookupAuthoritativeAllContext] instead.
func DNSLookupAuthoritativeAll(t testing.TestingT, query DNSQuery, resolvers []string) DNSAnswers {
	return DNSLookupAuthoritativeAllContext(t, context.Background(), query, resolvers)
}

// DNSLookupAuthoritativeAllE gets authoritative answers for the specified record and type.
// All the authoritative nameservers found must give the same answers.
// If resolvers are defined, uses them instead of the default system ones to find the authoritative nameservers.
// Returns [InconsistentAuthoritativeError] when any authoritative nameserver gives a different answer.
// Returns any underlying error.
//
// Deprecated: Use [DNSLookupAuthoritativeAllContextE] instead.
func DNSLookupAuthoritativeAllE(t testing.TestingT, query DNSQuery, resolvers []string) (DNSAnswers, error) {
	return DNSLookupAuthoritativeAllContextE(t, context.Background(), query, resolvers)
}

// DNSLookupAuthoritativeAllWithRetryContext repeatedly sends DNS requests for the specified record and type,
// until ALL authoritative nameservers reply with the exact same non-empty answers or until max retries has been exceeded.
// If defined, uses the given resolvers instead of the default system ones to find the authoritative nameservers.
// Fails when max retries has been exceeded.
func DNSLookupAuthoritativeAllWithRetryContext(t testing.TestingT, ctx context.Context, query DNSQuery, resolvers []string, maxRetries int, sleepBetweenRetries time.Duration) {
	_, err := DNSLookupAuthoritativeAllWithRetryContextE(t, ctx, query, resolvers, maxRetries, sleepBetweenRetries)
	require.NoError(t, err)
}

// DNSLookupAuthoritativeAllWithRetryContextE repeatedly sends DNS requests for the specified record and type,
// until ALL authoritative nameservers reply with the exact same non-empty answers or until max retries has been exceeded.
// If defined, uses the given resolvers instead of the default system ones to find the authoritative nameservers.
func DNSLookupAuthoritativeAllWithRetryContextE(t testing.TestingT, ctx context.Context, query DNSQuery, resolvers []string, maxRetries int, sleepBetweenRetries time.Duration) (DNSAnswers, error) {
	res, err := retry.DoWithRetryInterfaceContextE(
		t, ctx, fmt.Sprintf("DNSLookupAuthoritativeAllContextE %s record for %s using authoritative nameservers", query.Type, query.Name),
		maxRetries, sleepBetweenRetries,
		func() (any, error) {
			return DNSLookupAuthoritativeAllContextE(t, ctx, query, resolvers)
		})

	return res.(DNSAnswers), err
}

// DNSLookupAuthoritativeAllWithRetry repeatedly sends DNS requests for the specified record and type,
// until ALL authoritative nameservers reply with the exact same non-empty answers or until max retries has been exceeded.
// If defined, uses the given resolvers instead of the default system ones to find the authoritative nameservers.
// Fails when max retries has been exceeded.
//
// Deprecated: Use [DNSLookupAuthoritativeAllWithRetryContext] instead.
func DNSLookupAuthoritativeAllWithRetry(t testing.TestingT, query DNSQuery, resolvers []string, maxRetries int, sleepBetweenRetries time.Duration) {
	DNSLookupAuthoritativeAllWithRetryContext(t, context.Background(), query, resolvers, maxRetries, sleepBetweenRetries)
}

// DNSLookupAuthoritativeAllWithRetryE repeatedly sends DNS requests for the specified record and type,
// until ALL authoritative nameservers reply with the exact same non-empty answers or until max retries has been exceeded.
// If defined, uses the given resolvers instead of the default system ones to find the authoritative nameservers.
//
// Deprecated: Use [DNSLookupAuthoritativeAllWithRetryContextE] instead.
func DNSLookupAuthoritativeAllWithRetryE(t testing.TestingT, query DNSQuery, resolvers []string, maxRetries int, sleepBetweenRetries time.Duration) (DNSAnswers, error) {
	return DNSLookupAuthoritativeAllWithRetryContextE(t, context.Background(), query, resolvers, maxRetries, sleepBetweenRetries)
}

// DNSLookupAuthoritativeAllWithValidationContext gets authoritative answers for the specified record and type.
// All the authoritative nameservers found must give the same answers and match the expectedAnswers.
// If resolvers are defined, uses them instead of the default system ones to find the authoritative nameservers.
// Fails on any underlying error from [DNSLookupAuthoritativeAllWithValidationContextE].
func DNSLookupAuthoritativeAllWithValidationContext(t testing.TestingT, ctx context.Context, query DNSQuery, resolvers []string, expectedAnswers DNSAnswers) {
	err := DNSLookupAuthoritativeAllWithValidationContextE(t, ctx, query, resolvers, expectedAnswers)
	require.NoError(t, err)
}

// DNSLookupAuthoritativeAllWithValidationContextE gets authoritative answers for the specified record and type.
// All the authoritative nameservers found must give the same answers and match the expectedAnswers.
// If resolvers are defined, uses them instead of the default system ones to find the authoritative nameservers.
// Returns [ValidationError] when expectedAnswers differ from the obtained ones.
// Returns any underlying error from [DNSLookupAuthoritativeAllContextE].
func DNSLookupAuthoritativeAllWithValidationContextE(t testing.TestingT, ctx context.Context, query DNSQuery, resolvers []string, expectedAnswers DNSAnswers) error {
	expectedAnswers.Sort()

	answers, err := DNSLookupAuthoritativeAllContextE(t, ctx, query, resolvers)
	if err != nil {
		return err
	}

	if !reflect.DeepEqual(answers, expectedAnswers) {
		return &ValidationError{Query: query, Answers: answers, ExpectedAnswers: expectedAnswers}
	}

	return nil
}

// DNSLookupAuthoritativeAllWithValidation gets authoritative answers for the specified record and type.
// All the authoritative nameservers found must give the same answers and match the expectedAnswers.
// If resolvers are defined, uses them instead of the default system ones to find the authoritative nameservers.
// Fails on any underlying error from [DNSLookupAuthoritativeAllWithValidationContextE].
//
// Deprecated: Use [DNSLookupAuthoritativeAllWithValidationContext] instead.
func DNSLookupAuthoritativeAllWithValidation(t testing.TestingT, query DNSQuery, resolvers []string, expectedAnswers DNSAnswers) {
	DNSLookupAuthoritativeAllWithValidationContext(t, context.Background(), query, resolvers, expectedAnswers)
}

// DNSLookupAuthoritativeAllWithValidationE gets authoritative answers for the specified record and type.
// All the authoritative nameservers found must give the same answers and match the expectedAnswers.
// If resolvers are defined, uses them instead of the default system ones to find the authoritative nameservers.
// Returns [ValidationError] when expectedAnswers differ from the obtained ones.
// Returns any underlying error from [DNSLookupAuthoritativeAllContextE].
//
// Deprecated: Use [DNSLookupAuthoritativeAllWithValidationContextE] instead.
func DNSLookupAuthoritativeAllWithValidationE(t testing.TestingT, query DNSQuery, resolvers []string, expectedAnswers DNSAnswers) error {
	return DNSLookupAuthoritativeAllWithValidationContextE(t, context.Background(), query, resolvers, expectedAnswers)
}

// DNSLookupAuthoritativeAllWithValidationRetryContext repeatedly gets authoritative answers for the specified record and type
// until ALL the authoritative nameservers found give the same answers and match the expectedAnswers,
// or until max retries has been exceeded.
// If resolvers are defined, uses them instead of the default system ones to find the authoritative nameservers.
// Fails when max retries has been exceeded.
func DNSLookupAuthoritativeAllWithValidationRetryContext(t testing.TestingT, ctx context.Context, query DNSQuery, resolvers []string, expectedAnswers DNSAnswers, maxRetries int, sleepBetweenRetries time.Duration) {
	err := DNSLookupAuthoritativeAllWithValidationRetryContextE(t, ctx, query, resolvers, expectedAnswers, maxRetries, sleepBetweenRetries)
	require.NoError(t, err)
}

// DNSLookupAuthoritativeAllWithValidationRetryContextE repeatedly gets authoritative answers for the specified record and type
// until ALL the authoritative nameservers found give the same answers and match the expectedAnswers,
// or until max retries has been exceeded.
// If resolvers are defined, uses them instead of the default system ones to find the authoritative nameservers.
func DNSLookupAuthoritativeAllWithValidationRetryContextE(t testing.TestingT, ctx context.Context, query DNSQuery, resolvers []string, expectedAnswers DNSAnswers, maxRetries int, sleepBetweenRetries time.Duration) error {
	_, err := retry.DoWithRetryInterfaceContextE(
		t, ctx, fmt.Sprintf("DNSLookupAuthoritativeAllWithValidationRetryContextE %s record for %s using authoritative nameservers", query.Type, query.Name),
		maxRetries, sleepBetweenRetries,
		func() (any, error) {
			return nil, DNSLookupAuthoritativeAllWithValidationContextE(t, ctx, query, resolvers, expectedAnswers)
		})

	return err
}

// DNSLookupAuthoritativeAllWithValidationRetry repeatedly gets authoritative answers for the specified record and type
// until ALL the authoritative nameservers found give the same answers and match the expectedAnswers,
// or until max retries has been exceeded.
// If resolvers are defined, uses them instead of the default system ones to find the authoritative nameservers.
// Fails when max retries has been exceeded.
//
// Deprecated: Use [DNSLookupAuthoritativeAllWithValidationRetryContext] instead.
func DNSLookupAuthoritativeAllWithValidationRetry(t testing.TestingT, query DNSQuery, resolvers []string, expectedAnswers DNSAnswers, maxRetries int, sleepBetweenRetries time.Duration) {
	DNSLookupAuthoritativeAllWithValidationRetryContext(t, context.Background(), query, resolvers, expectedAnswers, maxRetries, sleepBetweenRetries)
}

// DNSLookupAuthoritativeAllWithValidationRetryE repeatedly gets authoritative answers for the specified record and type
// until ALL the authoritative nameservers found give the same answers and match the expectedAnswers,
// or until max retries has been exceeded.
// If resolvers are defined, uses them instead of the default system ones to find the authoritative nameservers.
//
// Deprecated: Use [DNSLookupAuthoritativeAllWithValidationRetryContextE] instead.
func DNSLookupAuthoritativeAllWithValidationRetryE(t testing.TestingT, query DNSQuery, resolvers []string, expectedAnswers DNSAnswers, maxRetries int, sleepBetweenRetries time.Duration) error {
	return DNSLookupAuthoritativeAllWithValidationRetryContextE(t, context.Background(), query, resolvers, expectedAnswers, maxRetries, sleepBetweenRetries)
}

// dnsLookupContext sends a DNS query for the specified record and type using the given resolver.
// Returns DNSAnswers to the DNSQuery.
// If no records found, returns [NotFoundError].
func dnsLookupContext(t testing.TestingT, ctx context.Context, query DNSQuery, resolver string) (DNSAnswers, error) {
	switch query.Type {
	case "A", "AAAA", "CNAME", "MX", "NS", "TXT":
	default:
		return nil, &QueryTypeError{Type: query.Type}
	}

	qType, ok := dns.StringToType[strings.ToUpper(query.Type)]
	if !ok {
		return nil, &QueryTypeError{Type: query.Type}
	}

	if strings.LastIndex(resolver, ":") <= strings.LastIndex(resolver, "]") {
		resolver += ":53"
	}

	c := new(dns.Client)
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(query.Name), qType)

	in, _, err := c.ExchangeContext(ctx, m, resolver)
	if err != nil {
		logger.Default.Logf(t, "Error sending DNS query %s: %s", query, err)

		return nil, err
	}

	if len(in.Answer) == 0 {
		return nil, &NotFoundError{Query: query, Nameserver: resolver}
	}

	var dnsAnswers DNSAnswers

	for _, a := range in.Answer {
		switch at := a.(type) {
		case *dns.A:
			dnsAnswers = append(dnsAnswers, DNSAnswer{Type: "A", Value: at.A.String()})
		case *dns.AAAA:
			dnsAnswers = append(dnsAnswers, DNSAnswer{Type: "AAAA", Value: at.AAAA.String()})
		case *dns.CNAME:
			dnsAnswers = append(dnsAnswers, DNSAnswer{Type: "CNAME", Value: at.Target})
		case *dns.NS:
			dnsAnswers = append(dnsAnswers, DNSAnswer{Type: "NS", Value: at.Ns})
		case *dns.MX:
			dnsAnswers = append(dnsAnswers, DNSAnswer{Type: "MX", Value: fmt.Sprintf("%d %s", at.Preference, at.Mx)})
		case *dns.TXT:
			for _, txt := range at.Txt {
				dnsAnswers = append(dnsAnswers, DNSAnswer{Type: "TXT", Value: fmt.Sprintf(`"%s"`, txt)})
			}
		}
	}

	dnsAnswers.Sort()

	return dnsAnswers, nil
}
