package terragrunt

import (
	"testing"
)

// Get calls terragrunt get and return stdout/stderr.
func Get(t *testing.T, options *Options) string {
	out, err := GetE(t, options)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// GetE calls terragrunt get and return stdout/stderr.
func GetE(t *testing.T, options *Options) (string, error) {
	return RunTerragruntCommandE(t, options, "get", "-update")
}
