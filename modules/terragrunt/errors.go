package terragrunt

import "errors"

// ErrNilOptions is returned when a nil Options pointer is passed to a function
// that requires a valid configuration.
var ErrNilOptions = errors.New("options cannot be nil")

// ErrMissingTerragruntDir is returned when the required TerragruntDir field
// is empty in the provided Options.
var ErrMissingTerragruntDir = errors.New("TerragruntDir is required")

// ErrEmptyTfArgs is returned when tfArgs is empty in a call that requires at
// least one OpenTofu/Terraform command argument (e.g. "apply", "plan").
var ErrEmptyTfArgs = errors.New("tfArgs cannot be empty; at minimum, an OpenTofu/Terraform command (e.g. \"apply\") is required")
