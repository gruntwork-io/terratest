package http_helper //nolint:staticcheck // package name determined by directory

import "fmt"

// ValidationFunctionFailed is an error that occurs if a validation function fails.
type ValidationFunctionFailed struct {
	Url    string //nolint:staticcheck // preserving existing field name
	Body   string
	Status int
}

func (err ValidationFunctionFailed) Error() string {
	return fmt.Sprintf("Validation failed for URL %s. Response status: %d. Response body:\n%s", err.Url, err.Status, err.Body)
}
