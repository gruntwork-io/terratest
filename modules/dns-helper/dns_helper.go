// Package dns_helper contains helpers to interact with the Domain Name System.
package dns_helper

import (
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

// DNSFindNameservers tries to find the NS record for the given FQDN, iterating down the domain hierarchy
// until it founds the NS records and returns it. Fails if there's any error or no NS record is found up to the apex domain.
func DNSFindNameservers(t testing.TestingT, fqdn string, resolvers []string) []string {
	nameservers, err := DNSFindNameserversE(t, fqdn, resolvers)
	require.NoError(t, err)
	return nameservers
}

// DNSFindNameserversE tries to find the NS record for the given FQDN, iterating down the domain hierarchy
// until it founds the NS records and returns it. Returns the last error if the apex domain is reached with no result.
func DNSFindNameserversE(t testing.TestingT, fqdn string, resolvers []string) ([]string, error) {
	var lookupFunc func(domain string) ([]string, error)

	if resolvers == nil {
		lookupFunc = func(domain string) ([]string, error) {
			var nameservers []string
			res, err := net.LookupNS(domain)
			for _, ns := range res {
				nameservers = append(nameservers, ns.Host)
			}
			return nameservers, err
		}
	} else {
		lookupFunc = func(domain string) ([]string, error) {
			var nameservers []string
			res, err := DNSLookupE(t, DNSQuery{"NS", domain}, resolvers)
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
			var nameservers []string

			for _, ns := range res {
				nameservers = append(nameservers, strings.TrimSuffix(ns, "."))
			}

			logger.Default.Logf(t, "FQDN %s belongs to domain %s, found NS record: %s", fqdn, domain, nameservers)
			return nameservers, nil
		}

		if err != nil {
			logger.Default.Logf(t, err.Error())
		}
	}

	err := &NSNotFoundError{fqdn, domain}
	return nil, err
}

// DNSLookupAuthoritative gets authoritative answers for the specified record and type.
// If resolvers are defined, uses them instead of the default system ones to find the authoritative nameservers.
// Fails on any error from DNSLookupAuthoritativeE.
func DNSLookupAuthoritative(t testing.TestingT, query DNSQuery, resolvers []string) DNSAnswers {
	res, err := DNSLookupAuthoritativeE(t, query, resolvers)
	require.NoError(t, err)
	return res
}

// DNSLookupAuthoritativeE gets authoritative answers for the specified record and type.
// If resolvers are defined, uses them instead of the default system ones to find the authoritative nameservers.
// Returns NotFoundError when no answer found in any authoritative nameserver.
// Returns any underlying error from individual lookups.
func DNSLookupAuthoritativeE(t testing.TestingT, query DNSQuery, resolvers []string) (DNSAnswers, error) {
	nameservers, err := DNSFindNameserversE(t, query.Name, resolvers)

	if err != nil {
		return nil, err
	}

	return DNSLookupE(t, query, nameservers)
}

// DNSLookupAuthoritativeWithRetry repeatedly gets authoritative answers for the specified record and type
// until ANY of the authoritative nameservers found replies with non-empty answer matching the expectedAnswers,
// or until max retries has been exceeded.
// If resolvers are defined, uses them instead of the default system ones to find the authoritative nameservers.
// Fails on any error from DNSLookupAuthoritativeWithRetryE.
func DNSLookupAuthoritativeWithRetry(t testing.TestingT, query DNSQuery, resolvers []string, maxRetries int, sleepBetweenRetries time.Duration) DNSAnswers {
	res, err := DNSLookupAuthoritativeWithRetryE(t, query, resolvers, maxRetries, sleepBetweenRetries)
	require.NoError(t, err)
	return res
}

// DNSLookupAuthoritativeWithRetryE repeatedly gets authoritative answers for the specified record and type
// until ANY of the authoritative nameservers found replies with non-empty answer matching the expectedAnswers,
// or until max retries has been exceeded.
// If resolvers are defined, uses them instead of the default system ones to find the authoritative nameservers.
func DNSLookupAuthoritativeWithRetryE(t testing.TestingT, query DNSQuery, resolvers []string, maxRetries int, sleepBetweenRetries time.Duration) (DNSAnswers, error) {
	res, err := retry.DoWithRetryInterfaceE(
		t, fmt.Sprintf("DNSLookupAuthoritativeE %s record for %s using authoritative nameservers", query.Type, query.Name),
		maxRetries, sleepBetweenRetries,
		func() (interface{}, error) {
			return DNSLookupAuthoritativeE(t, query, resolvers)
		})

	return res.(DNSAnswers), err
}

// DNSLookupAuthoritativeAll gets authoritative answers for the specified record and type.
// All the authoritative nameservers found must give the same answers.
// If resolvers are defined, uses them instead of the default system ones to find the authoritative nameservers.
// Fails on any error from DNSLookupAuthoritativeAllE.
func DNSLookupAuthoritativeAll(t testing.TestingT, query DNSQuery, resolvers []string) DNSAnswers {
	res, err := DNSLookupAuthoritativeAllE(t, query, resolvers)
	require.NoError(t, err)
	return res
}

// DNSLookupAuthoritativeAllE gets authoritative answers for the specified record and type.
// All the authoritative nameservers found must give the same answers.
// If resolvers are defined, uses them instead of the default system ones to find the authoritative nameservers.
// Returns InconsistentAuthoritativeError when any authoritative nameserver gives a different answer.
// Returns any underlying error.
func DNSLookupAuthoritativeAllE(t testing.TestingT, query DNSQuery, resolvers []string) (DNSAnswers, error) {
	nameservers, err := DNSFindNameserversE(t, query.Name, resolvers)

	if err != nil {
		return nil, err
	}

	var answers DNSAnswers

	for _, ns := range nameservers {
		res, err := DNSLookupE(t, query, []string{ns})

		if err != nil {
			return nil, err
		}

		if len(answers) > 0 {
			if !reflect.DeepEqual(answers, res) {
				err := &InconsistentAuthoritativeError{Query: query, Answers: res, Nameserver: ns, PreviousAnswers: answers}
				return nil, err
			}
		} else {
			answers = res
		}
	}

	return answers, nil
}

// DNSLookupAuthoritativeAllWithRetry repeatedly sends DNS requests for the specified record and type,
// until ALL authoritative nameservers reply with the exact same non-empty answers or until max retries has been exceeded.
// If defined, uses the given resolvers instead of the default system ones to find the authoritative nameservers.
// Fails when max retries has been exceeded.
func DNSLookupAuthoritativeAllWithRetry(t testing.TestingT, query DNSQuery, resolvers []string, maxRetries int, sleepBetweenRetries time.Duration) {
	_, err := DNSLookupAuthoritativeAllWithRetryE(t, query, resolvers, maxRetries, sleepBetweenRetries)
	require.NoError(t, err)
}

// DNSLookupAuthoritativeAllWithRetryE repeatedly sends DNS requests for the specified record and type,
// until ALL authoritative nameservers reply with the exact same non-empty answers or until max retries has been exceeded.
// If defined, uses the given resolvers instead of the default system ones to find the authoritative nameservers.
func DNSLookupAuthoritativeAllWithRetryE(t testing.TestingT, query DNSQuery, resolvers []string, maxRetries int, sleepBetweenRetries time.Duration) (DNSAnswers, error) {
	res, err := retry.DoWithRetryInterfaceE(
		t, fmt.Sprintf("DNSLookupAuthoritativeAllE %s record for %s using authoritative nameservers", query.Type, query.Name),
		maxRetries, sleepBetweenRetries,
		func() (interface{}, error) {
			return DNSLookupAuthoritativeAllE(t, query, resolvers)
		})

	return res.(DNSAnswers), err
}

// DNSLookupAuthoritativeAllWithValidation gets authoritative answers for the specified record and type.
// All the authoritative nameservers found must give the same answers and match the expectedAnswers.
// If resolvers are defined, uses them instead of the default system ones to find the authoritative nameservers.
// Fails on any underlying error from DNSLookupAuthoritativeAllWithValidationE.
func DNSLookupAuthoritativeAllWithValidation(t testing.TestingT, query DNSQuery, resolvers []string, expectedAnswers DNSAnswers) {
	err := DNSLookupAuthoritativeAllWithValidationE(t, query, resolvers, expectedAnswers)
	require.NoError(t, err)
}

// DNSLookupAuthoritativeAllWithValidationE gets authoritative answers for the specified record and type.
// All the authoritative nameservers found must give the same answers and match the expectedAnswers.
// If resolvers are defined, uses them instead of the default system ones to find the authoritative nameservers.
// Returns ValidationError when expectedAnswers differ from the obtained ones.
// Returns any underlying error from DNSLookupAuthoritativeAllE.
func DNSLookupAuthoritativeAllWithValidationE(t testing.TestingT, query DNSQuery, resolvers []string, expectedAnswers DNSAnswers) error {
	expectedAnswers.Sort()

	answers, err := DNSLookupAuthoritativeAllE(t, query, resolvers)

	if err != nil {
		return err
	}

	if !reflect.DeepEqual(answers, expectedAnswers) {
		err := &ValidationError{Query: query, Answers: answers, ExpectedAnswers: expectedAnswers}
		return err
	}

	return nil
}

// DNSLookupAuthoritativeAllWithValidationRetry repeatedly gets authoritative answers for the specified record and type
// until ALL the authoritative nameservers found give the same answers and match the expectedAnswers,
// or until max retries has been exceeded.
// If resolvers are defined, uses them instead of the default system ones to find the authoritative nameservers.
// Fails when max retries has been exceeded.
func DNSLookupAuthoritativeAllWithValidationRetry(t testing.TestingT, query DNSQuery, resolvers []string, expectedAnswers DNSAnswers, maxRetries int, sleepBetweenRetries time.Duration) {
	err := DNSLookupAuthoritativeAllWithValidationRetryE(t, query, resolvers, expectedAnswers, maxRetries, sleepBetweenRetries)
	require.NoError(t, err)
}

// DNSLookupAuthoritativeAllWithValidationRetryE repeatedly gets authoritative answers for the specified record and type
// until ALL the authoritative nameservers found give the same answers and match the expectedAnswers,
// or until max retries has been exceeded.
// If resolvers are defined, uses them instead of the default system ones to find the authoritative nameservers.
func DNSLookupAuthoritativeAllWithValidationRetryE(t testing.TestingT, query DNSQuery, resolvers []string, expectedAnswers DNSAnswers, maxRetries int, sleepBetweenRetries time.Duration) error {
	_, err := retry.DoWithRetryInterfaceE(
		t, fmt.Sprintf("DNSLookupAuthoritativeAllWithValidationRetryE %s record for %s using authoritative nameservers", query.Type, query.Name),
		maxRetries, sleepBetweenRetries,
		func() (interface{}, error) {
			return nil, DNSLookupAuthoritativeAllWithValidationE(t, query, resolvers, expectedAnswers)
		})

	return err
}

// DNSLookup sends a DNS query for the specified record and type using the given resolvers.
// Fails on any error.
// Supported record types: A, AAAA, CNAME, MX, NS, TXT
func DNSLookup(t testing.TestingT, query DNSQuery, resolvers []string) DNSAnswers {
	res, err := DNSLookupE(t, query, resolvers)
	require.NoError(t, err)
	return res
}

// DNSLookupE sends a DNS query for the specified record and type using the given resolvers.
// Returns QueryTypeError when record type is not supported.
// Returns any underlying error.
// Supported record types: A, AAAA, CNAME, MX, NS, TXT
func DNSLookupE(t testing.TestingT, query DNSQuery, resolvers []string) (DNSAnswers, error) {
	if len(resolvers) == 0 {
		err := &NoResolversError{}
		return nil, err
	}

	var dnsAnswers DNSAnswers
	var err error
	for _, resolver := range resolvers {
		dnsAnswers, err = dnsLookup(t, query, resolver)

		if err == nil {
			return dnsAnswers, nil
		}
	}

	return nil, err
}

// dnsLookup sends a DNS query for the specified record and type using the given resolver.
// Returns DNSAnswers to the DNSQuery.
// If no records found, returns NotFoundError.
func dnsLookup(t testing.TestingT, query DNSQuery, resolver string) (DNSAnswers, error) {
	switch query.Type {
	case "A", "AAAA", "CNAME", "MX", "NS", "TXT":
	default:
		err := &QueryTypeError{query.Type}
		return nil, err
	}

	qType, ok := dns.StringToType[strings.ToUpper(query.Type)]
	if !ok {
		err := &QueryTypeError{query.Type}
		return nil, err
	}

	if strings.LastIndex(resolver, ":") <= strings.LastIndex(resolver, "]") {
		resolver += ":53"
	}

	c := new(dns.Client)
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(query.Name), qType)

	in, _, err := c.Exchange(m, resolver)
	if err != nil {
		logger.Default.Logf(t, "Error sending DNS query %s: %s", query, err)
		return nil, err
	}

	if len(in.Answer) == 0 {
		err := &NotFoundError{query, resolver}
		return nil, err
	}

	var dnsAnswers DNSAnswers

	for _, a := range in.Answer {
		switch at := a.(type) {
		case *dns.A:
			dnsAnswers = append(dnsAnswers, DNSAnswer{"A", at.A.String()})
		case *dns.AAAA:
			dnsAnswers = append(dnsAnswers, DNSAnswer{"AAAA", at.AAAA.String()})
		case *dns.CNAME:
			dnsAnswers = append(dnsAnswers, DNSAnswer{"CNAME", at.Target})
		case *dns.NS:
			dnsAnswers = append(dnsAnswers, DNSAnswer{"NS", at.Ns})
		case *dns.MX:
			dnsAnswers = append(dnsAnswers, DNSAnswer{"MX", fmt.Sprintf("%d %s", at.Preference, at.Mx)})
		case *dns.TXT:
			for _, txt := range at.Txt {
				dnsAnswers = append(dnsAnswers, DNSAnswer{"TXT", fmt.Sprintf(`"%s"`, txt)})
			}
		}
	}

	dnsAnswers.Sort()

	return dnsAnswers, nil
}

// DNSQuery type
type DNSQuery struct {
	Type, Name string
}

// DNSAnswer type
type DNSAnswer struct {
	Type, Value string
}

func (a DNSAnswer) String() string {
	return fmt.Sprintf("%s %s", a.Type, a.Value)
}

// DNSAnswers type
type DNSAnswers []DNSAnswer

// Sort sorts the answers by type and value
func (a DNSAnswers) Sort() {
	sort.Slice(a, func(i, j int) bool {
		return a[i].Type < a[j].Type || a[i].Value < a[j].Value
	})
}
