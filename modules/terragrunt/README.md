# Terragrunt Module

Testing library for Terragrunt configurations in Go. Provides helpers for running Terragrunt commands for single units, across multiple modules (run-all), and stack-based workflows.

## Requirements

- **Terragrunt** binary in PATH
- **OpenTofu** or **Terraform** binary in PATH (Terragrunt is a wrapper and requires one of these)

To specify which binary to use (terraform vs opentofu):
```go
// Option 1: Via environment variable
options := &terragrunt.Options{
    TerragruntDir: "/path/to/config",
    EnvVars: map[string]string{
        "TERRAGRUNT_TFPATH": "/usr/local/bin/tofu",  // or "TG_TF_PATH"
    },
}

// Option 2: Via command-line flag
options := &terragrunt.Options{
    TerragruntDir:  "/path/to/config",
    TerragruntArgs: []string{"--tf-path", "/usr/local/bin/tofu"},
}
```

## Quick Start

### Single Unit

```go
import (
    "testing"
    "github.com/gruntwork-io/terratest/modules/terragrunt"
    "github.com/stretchr/testify/assert"
)

func TestSingleUnit(t *testing.T) {
    t.Parallel()

    ctx := t.Context()

    options := &terragrunt.Options{
        TerragruntDir: "../path/to/terragrunt/unit",
    }

    defer terragrunt.DestroyContext(t, ctx, options)
    terragrunt.InitAndApplyContext(t, ctx, options)

    // Get a specific output as JSON
    vpcOutput := terragrunt.OutputJSONContext(t, ctx, options, "vpc_id")
    assert.Contains(t, vpcOutput, "vpc-")
}
```

### Multiple Modules (--all)

```go
func TestTerragruntApply(t *testing.T) {
    t.Parallel()

    ctx := t.Context()

    options := &terragrunt.Options{
        TerragruntDir: "../path/to/terragrunt/config",
    }

    defer terragrunt.DestroyAllContext(t, ctx, options)
    terragrunt.ApplyAllContext(t, ctx, options)
}
```

## Key Concepts

### Options Struct

The `Options` struct has two distinct parts:

1. **Test Framework Configuration** (NOT passed to terragrunt CLI):
   - `TerragruntDir` - where to run terragrunt (required)
   - `TerragruntBinary` - binary name (default: "terragrunt")
   - `EnvVars` - environment variables
   - `Logger` - custom logger for output
   - `MaxRetries`, `TimeBetweenRetries` - retry settings
   - `RetryableTerraformErrors` - map of error patterns to retry messages
   - `WarningsAsErrors` - map of warning patterns to treat as errors
   - `BackendConfig` - backend configuration passed to `init`
   - `PluginDir` - plugin directory passed to `init`
   - `Stdin` - stdin reader for commands

2. **Command-Line Arguments** (passed to terragrunt):
   - `TerragruntArgs` - global terragrunt flags (e.g., `--log-level`, `--no-color`)
   - `TerraformArgs` - command-specific OpenTofu/Terraform flags (e.g., `-upgrade`)

### Error-Returning Variants (E-suffix)

Every function has an `E`-suffix variant that returns an error instead of calling `t.Fatal` on failure. For example:

- `ApplyContext(t, ctx, options)` calls `t.Fatal` on error
- `ApplyContextE(t, ctx, options)` returns `(string, error)` for custom error handling

Use `E` variants when you need to test error cases or handle failures gracefully:
```go
_, err := terragrunt.ApplyContextE(t, t.Context(), options)
require.Error(t, err)
```

### TerragruntArgs vs TerraformArgs

Arguments are passed in this order:
```text
terragrunt [TerragruntArgs] --non-interactive run -- <command> [TerraformArgs]
```

**Example:**
```go
options := &terragrunt.Options{
    TerragruntDir:  "/path/to/config",
    TerragruntArgs: []string{"--log-level", "error"},  // Global TG flags
    TerraformArgs:  []string{"-upgrade"},              // OpenTofu/Terraform flags
}
// Executes: terragrunt --log-level error --non-interactive run -- init -upgrade
```

## Functions

### Single-Unit Commands

Run terragrunt commands against a single unit (one `terragrunt.hcl` directory):

- `InitContext(t, ctx, options)` - Initialize configuration
- `ApplyContext(t, ctx, options)` - Apply changes
- `DestroyContext(t, ctx, options)` - Destroy resources
- `PlanContext(t, ctx, options)` - Generate and show execution plan
- `PlanExitCodeContext(t, ctx, options)` - Plan and return exit code (0=no changes, 2=changes, other=error)
- `ValidateContext(t, ctx, options)` - Validate configuration
- `OutputJSONContext(t, ctx, options, key)` - Get output as JSON (specific key or all outputs)

### Convenience Wrappers

Run init + command in a single call:

- `InitAndApplyContext(t, ctx, options)` - Init then apply
- `InitAndPlanContext(t, ctx, options)` - Init then plan
- `InitAndValidateContext(t, ctx, options)` - Init then validate

### Run Command

- `RunContext(t, ctx, options, tgArgs, tfArgs)` - Run any OpenTofu/Terraform command via `terragrunt run [tgArgs...] -- [tfArgs...]`

The `--` separator disambiguates Terragrunt flags (like `--all`) from OpenTofu/Terraform flags. The OpenTofu/Terraform command (e.g. `"apply"`) should be the first element of `tfArgs`.

### Run --all Commands

Work with [implicit stacks](https://terragrunt.gruntwork.io/docs/features/stacks/#implicit-stacks) (multiple units in a directory):

- `ApplyAllContext(t, ctx, options)` - Apply all modules with dependencies
- `DestroyAllContext(t, ctx, options)` - Destroy all modules with dependencies
- `PlanAllExitCodeContext(t, ctx, options)` - Plan all and return exit code (0=no changes, 2=changes, other=error)
- `ValidateAllContext(t, ctx, options)` - Validate all modules
- `RunAllContext(t, ctx, options, command)` - *Deprecated: use `RunContext` with `--all` in tgArgs instead.* Run any OpenTofu/Terraform command with --all flag
- `OutputAllJSONContext(t, ctx, options)` - Get all outputs as raw JSON string (note: returns separate JSON objects per module)

### HCL Commands

Terragrunt HCL tooling commands:

- `FormatAllContext(t, ctx, options)` - Format all terragrunt.hcl files (`terragrunt hcl format`)
- `HclValidateContext(t, ctx, options)` - Validate terragrunt.hcl syntax and configuration (`terragrunt hcl validate`)

### Configuration Commands

- `RenderContext(t, ctx, options)` - Render resolved terragrunt configuration as HCL
- `RenderJSONContext(t, ctx, options)` - Render resolved terragrunt configuration as JSON
- `GraphContext(t, ctx, options)` - Output dependency graph in DOT format

### Stack Commands

Work with [explicit stacks](https://terragrunt.gruntwork.io/docs/features/stacks/#explicit-stacks) (a directory with a `terragrunt.stack.hcl` file):

- `StackGenerateContext(t, ctx, options)` - Generate stack from stack.hcl
- `StackRunContext(t, ctx, options)` - Run command on generated stack
- `StackCleanContext(t, ctx, options)` - Remove .terragrunt-stack directory
- `StackOutputContext(t, ctx, options, key)` - Get stack output value
- `StackOutputJSONContext(t, ctx, options, key)` - Get stack output as JSON
- `StackOutputAllContext(t, ctx, options)` - Get all stack outputs as map
- `StackOutputListAllContext(t, ctx, options)` - Get list of all output variable names

## Examples

See the [examples directory](../../examples/) for complete working examples:
- [terragrunt-example](../../examples/terragrunt-example/) - Single unit testing
- [terragrunt-multi-module-example](../../examples/terragrunt-multi-module-example/) - Multi-module testing
- [terragrunt-second-example](../../examples/terragrunt-second-example/) - Additional patterns

### Testing with Dependencies

```go
func TestStack(t *testing.T) {
    t.Parallel()

    ctx := t.Context()

    options := &terragrunt.Options{
        TerragruntDir: "../live/prod",
    }

    // Apply respects dependency order
    terragrunt.ApplyAllContext(t, ctx, options)
    defer terragrunt.DestroyAllContext(t, ctx, options)

    // Verify infrastructure
    // ... your assertions here
}
```

### Using Custom Arguments

```go
func TestWithCustomArgs(t *testing.T) {
    t.Parallel()

    options := &terragrunt.Options{
        TerragruntDir:  "../config",
        TerragruntArgs: []string{"--log-level", "error", "--no-color"},
        TerraformArgs:  []string{"-upgrade"},
    }

    terragrunt.InitContext(t, t.Context(), options)
}
```

### Testing Stack Outputs

```go
func TestStackOutput(t *testing.T) {
    t.Parallel()

    ctx := t.Context()

    options := &terragrunt.Options{
        TerragruntDir: "../stack",
    }

    applyOpts := &terragrunt.Options{
        TerragruntDir: "../stack",
        TerraformArgs: []string{"apply"},
    }
    destroyOpts := &terragrunt.Options{
        TerragruntDir: "../stack",
        TerraformArgs: []string{"destroy"},
    }

    terragrunt.StackRunContext(t, ctx, applyOpts)
    defer terragrunt.StackRunContext(t, ctx, destroyOpts)

    // Get specific output
    vpcID := terragrunt.StackOutputContext(t, ctx, options, "vpc_id")
    assert.NotEmpty(t, vpcID)

    // Get all outputs
    outputs := terragrunt.StackOutputAllContext(t, ctx, options)
    assert.Contains(t, outputs, "vpc_id")
}
```

### Checking Plan Exit Code

```go
func TestInfrastructureUpToDate(t *testing.T) {
    t.Parallel()

    ctx := t.Context()

    options := &terragrunt.Options{
        TerragruntDir: "../prod",
    }

    // First apply
    terragrunt.ApplyAllContext(t, ctx, options)
    defer terragrunt.DestroyAllContext(t, ctx, options)

    // Plan should show no changes (exit code 0)
    exitCode := terragrunt.PlanAllExitCodeContext(t, ctx, options)
    assert.Equal(t, 0, exitCode, "No changes expected")
}
```

### Using Run for Flexibility

```go
func TestCustomCommand(t *testing.T) {
    t.Parallel()

    ctx := t.Context()

    options := &terragrunt.Options{
        TerragruntDir: "../modules",
    }

    // Run any OpenTofu/Terraform command with --all
    terragrunt.RunContext(t, ctx, options, []string{"--all"}, []string{"refresh"})

    // Verify state is current
    output := terragrunt.RunContext(t, ctx, options, []string{"--all"}, []string{"show"})
    assert.Contains(t, output, "expected-resource")
}
```

### Validating Stack Output Keys

```go
func TestStackOutputKeys(t *testing.T) {
    t.Parallel()

    ctx := t.Context()

    options := &terragrunt.Options{
        TerragruntDir: "../stack",
    }

    applyOpts := &terragrunt.Options{
        TerragruntDir: "../stack",
        TerraformArgs: []string{"apply"},
    }
    destroyOpts := &terragrunt.Options{
        TerragruntDir: "../stack",
        TerraformArgs: []string{"destroy"},
    }

    terragrunt.StackRunContext(t, ctx, applyOpts)
    defer terragrunt.StackRunContext(t, ctx, destroyOpts)

    // Get list of all output keys
    keys := terragrunt.StackOutputListAllContext(t, ctx, options)

    // Verify required outputs exist
    assert.Contains(t, keys, "vpc_id")
    assert.Contains(t, keys, "subnet_ids")
}
```

### Using Filters (Terragrunt v0.97.0+)

```go
options := &terragrunt.Options{
    TerragruntDir:  "../live/prod",
    TerragruntArgs: []string{"--filter", "{./vpc}"},  // Only apply vpc
}
terragrunt.ApplyAllContext(t, t.Context(), options)
```

## Not Supported

This module does **NOT** have dedicated helpers for:
- `import`, `refresh`, `show`, `state`, `test` commands
- `backend`, `exec`, `catalog`, `scaffold` commands
- Discovery commands (`find`, `list`)
- Configuration commands (`info`)

For these commands, use `RunContext` / `RunContextE` or run terragrunt directly via the `shell` module.

## Compatibility

Tested with Terragrunt v1.0.x. Earlier v0.x versions may work but are not guaranteed.

### Migration from terraform Module

The following functions were previously in the `terraform` module and have been moved here. The deprecated versions have been removed from the `terraform` module.

| Removed (terraform module) | Replacement (terragrunt module) |
|----------------------------|----------------------------------|
| `TgApplyAll` / `TgApplyAllE` | `ApplyAll` / `ApplyAllE` |
| `TgDestroyAll` / `TgDestroyAllE` | `DestroyAll` / `DestroyAllE` |
| `TgPlanAllExitCode` / `TgPlanAllExitCodeE` | `PlanAllExitCode` / `PlanAllExitCodeE` |
| `ValidateInputs` / `ValidateInputsE` | `HclValidate` / `HclValidateE` |

> **Note:** `ValidateInputs` specifically checked input alignment. For equivalent behavior, pass `TerraformArgs: []string{"--inputs"}` to `HclValidate`.

## More Info

- [Terragrunt Documentation](https://terragrunt.gruntwork.io/)
- [Terratest Documentation](https://terratest.gruntwork.io/)
