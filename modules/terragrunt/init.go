package terragrunt

import (
	"fmt"
	"testing"
)

// Init calls terragrunt init and return stdout/stderr.
func Init(t *testing.T, options *Options) string {
	out, err := InitE(t, options)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// InitE calls terragrunt init and return stdout/stderr.
func InitE(t *testing.T, options *Options) (string, error) {
	args := []string{"init", fmt.Sprintf("-upgrade=%t", options.Upgrade)}
	args = append(args, FormatTerragruntBackendConfigAsArgs(options.BackendConfig)...)
	return RunTerragruntCommandE(t, options, args...)
}
