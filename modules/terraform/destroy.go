package terraform

import (
	"testing"
)

// Destroy runs terraform destroy with the given options and return stdout/stderr.
func Destroy(t *testing.T, options *Options) string {
	out, err := DestroyE(t, options)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// DestroyE runs terraform destroy with the given options and return stdout/stderr.
func DestroyE(t *testing.T, options *Options) (string, error) {
	return RunTerraformCommandE(t, options, FormatArgs(options.Vars, "destroy", "-force", "-input=false", "-lock=false")...)
}
