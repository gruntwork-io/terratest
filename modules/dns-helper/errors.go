package dns_helper

import "fmt"

// NoResolversError is an error that occurs if no resolvers have been set for DNSLookupE
type NoResolversError struct{}

func (err NoResolversError) Error() string {
	return "No resolvers set for DNSLookupE call"
}

// QueryTypeError is an error that occurs if the DNS query type is not supported
type QueryTypeError struct {
	Type string
}

func (err QueryTypeError) Error() string {
	return fmt.Sprintf("Wrong DNS query type: %s", err.Type)
}

// NotFoundError is an error that occurs if no answer found
type NotFoundError struct {
	Query      DNSQuery
	Nameserver string
}

func (err NotFoundError) Error() string {
	return fmt.Sprintf("No %s record found for %s querying nameserver %s", err.Query.Type, err.Query.Name, err.Nameserver)
}

// InconsistentAuthoritativeError is an error that occurs if an authoritative answer is different from another
type InconsistentAuthoritativeError struct {
	Query           DNSQuery
	Answers         DNSAnswers
	Nameserver      string
	PreviousAnswers DNSAnswers
}

func (err InconsistentAuthoritativeError) Error() string {
	return fmt.Sprintf("Inconsistent authoritative answer from %s to DNS query %s. Got: %s Previous: %s", err.Nameserver, err.Query, err.Answers, err.PreviousAnswers)
}

// NSNotFoundError is an error that occurs if no NS records found
type NSNotFoundError struct {
	FQDN       string
	Nameserver string
}

func (err NSNotFoundError) Error() string {
	return fmt.Sprintf("No NS record found for %s up to apex domain %s", err.FQDN, err.Nameserver)
}

// MaxRetriesExceeded is an error that occurs when the maximum amount of retries is exceeded.
type MaxRetriesExceeded struct {
	Description string
	MaxRetries  int
}

func (err MaxRetriesExceeded) Error() string {
	return fmt.Sprintf("'%s' unsuccessful after %d retries", err.Description, err.MaxRetries)
}

// ValidationError is an error that occurs when answers validation fails
type ValidationError struct {
	Query           DNSQuery
	Answers         DNSAnswers
	ExpectedAnswers DNSAnswers
}

func (err ValidationError) Error() string {
	return fmt.Sprintf("Unexpected answer to DNS query %s. Got: %s Expected: %s", err.Query, err.Answers, err.ExpectedAnswers)
}
