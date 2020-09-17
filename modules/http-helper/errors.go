package http_helper

import "fmt"

// ValidationFunctionFailed is an error that occurs if a validation function fails.
type ValidationFunctionFailed struct {
	Url    string
	Status int
	Body   string
}

func (err ValidationFunctionFailed) Error() string {
	return fmt.Sprintf("Validation failed for URL %s. Response status: %d. Response body:\n%s", err.Url, err.Status, err.Body)
}
