package terragrunt

import (
	"testing"
)

// InitAndApply runs terragrunt init and apply with the given options and return stdout/stderr from the apply command. Note that this
// method does NOT call destroy and assumes the caller is responsible for cleaning up any resources created by running
// apply.
func InitAndApply(t *testing.T, options *Options) string {
	out, err := InitAndApplyE(t, options)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// InitAndApplyE runs terragrunt init and apply with the given options and return stdout/stderr from the apply command. Note that this
// method does NOT call destroy and assumes the caller is responsible for cleaning up any resources created by running
// apply.
func InitAndApplyE(t *testing.T, options *Options) (string, error) {
	if _, err := InitE(t, options); err != nil {
		return "", err
	}

	if _, err := GetE(t, options); err != nil {
		return "", err
	}

	return ApplyE(t, options)
}

// Apply runs terragrunt apply with the given options and return stdout/stderr. Note that this method does NOT call destroy and
// assumes the caller is responsible for cleaning up any resources created by running apply.
func Apply(t *testing.T, options *Options) string {
	out, err := ApplyE(t, options)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// ApplyE runs terragrunt apply with the given options and return stdout/stderr. Note that this method does NOT call destroy and
// assumes the caller is responsible for cleaning up any resources created by running apply.
func ApplyE(t *testing.T, options *Options) (string, error) {
	return RunTerragruntCommandE(t, options, FormatArgs(options.Vars, "apply", "-input=false", "-lock=false", "-auto-approve")...)
}
