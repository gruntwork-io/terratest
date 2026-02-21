# Terragrunt Module

Testing library for Terragrunt configurations in Go. Provides helpers for running Terragrunt commands for single units, across multiple modules (run-all), and stack-based workflows.

## Requirements

- **Terragrunt** binary in PATH
- **Terraform** or **OpenTofu** binary in PATH (Terragrunt is a wrapper and requires one of these)

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

    options := &terragrunt.Options{
        TerragruntDir: "../path/to/terragrunt/unit",
    }

    defer terragrunt.Destroy(t, options)
    terragrunt.InitAndApply(t, options)

    // Get a specific output as JSON
    vpcOutput := terragrunt.OutputJson(t, options, "vpc_id")
    assert.Contains(t, vpcOutput, "vpc-")
}
```

### Multiple Modules (--all)

```go
func TestTerragruntApply(t *testing.T) {
    t.Parallel()

    options := &terragrunt.Options{
        TerragruntDir: "../path/to/terragrunt/config",
    }

    defer terragrunt.DestroyAll(t, options)
    terragrunt.ApplyAll(t, options)
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
   - `TerraformArgs` - command-specific terraform flags (e.g., `-upgrade`)

### Error-Returning Variants (E-suffix)

Every function has an `E`-suffix variant that returns an error instead of calling `t.Fatal` on failure. For example:

- `Apply(t, options)` calls `t.Fatal` on error
- `ApplyE(t, options)` returns `(string, error)` for custom error handling

Use `E` variants when you need to test error cases or handle failures gracefully:
```go
_, err := terragrunt.ApplyE(t, options)
require.Error(t, err)
```

### TerragruntArgs vs TerraformArgs

Arguments are passed in this order:
```
terragrunt [TerragruntArgs] --non-interactive <command> [TerraformArgs]
```

**Example:**
```go
options := &terragrunt.Options{
    TerragruntDir:  "/path/to/config",
    TerragruntArgs: []string{"--log-level", "error"},  // Global TG flags
    TerraformArgs:  []string{"-upgrade"},              // Terraform flags
}
// Executes: terragrunt --log-level error --non-interactive init -upgrade
```

## Functions

### Single-Unit Commands

Run terragrunt commands against a single unit (one `terragrunt.hcl` directory):

- `Init(t, options)` - Initialize configuration
- `Apply(t, options)` - Apply changes
- `Destroy(t, options)` - Destroy resources
- `Plan(t, options)` - Generate and show execution plan
- `PlanExitCode(t, options)` - Plan and return exit code (0=no changes, 2=changes, other=error)
- `Validate(t, options)` - Validate configuration
- `OutputJson(t, options, key)` - Get output as JSON (specific key or all outputs)

### Convenience Wrappers

Run init + command in a single call:

- `InitAndApply(t, options)` - Init then apply
- `InitAndPlan(t, options)` - Init then plan
- `InitAndValidate(t, options)` - Init then validate

### Run --all Commands

Work with [implicit stacks](https://terragrunt.gruntwork.io/docs/features/stacks/#implicit-stacks) (multiple units in a directory):

- `ApplyAll(t, options)` - Apply all modules with dependencies
- `DestroyAll(t, options)` - Destroy all modules with dependencies
- `PlanAllExitCode(t, options)` - Plan all and return exit code (0=no changes, 2=changes, other=error)
- `ValidateAll(t, options)` - Validate all modules
- `RunAll(t, options, command)` - Run any terraform command with --all flag
- `OutputAllJson(t, options)` - Get all outputs as raw JSON string (note: returns separate JSON objects per module)

### HCL Commands

Terragrunt HCL tooling commands:

- `FormatAll(t, options)` - Format all terragrunt.hcl files (`terragrunt hcl format`)
- `HclValidate(t, options)` - Validate terragrunt.hcl syntax and configuration (`terragrunt hcl validate`)

### Configuration Commands

- `Render(t, options)` - Render resolved terragrunt configuration as HCL
- `RenderJson(t, options)` - Render resolved terragrunt configuration as JSON
- `Graph(t, options)` - Output dependency graph in DOT format

### Stack Commands

Work with [explicit stacks](https://terragrunt.gruntwork.io/docs/features/stacks/#explicit-stacks) (a directory with a `terragrunt.stack.hcl` file):

- `StackGenerate(t, options)` - Generate stack from stack.hcl
- `StackRun(t, options)` - Run command on generated stack
- `StackClean(t, options)` - Remove .terragrunt-stack directory
- `StackOutput(t, options, key)` - Get stack output value
- `StackOutputJson(t, options, key)` - Get stack output as JSON
- `StackOutputAll(t, options)` - Get all stack outputs as map
- `StackOutputListAll(t, options)` - Get list of all output variable names

## Examples

See the [examples directory](../../examples/) for complete working examples:
- [terragrunt-example](../../examples/terragrunt-example/) - Single unit testing
- [terragrunt-multi-module-example](../../examples/terragrunt-multi-module-example/) - Multi-module testing
- [terragrunt-second-example](../../examples/terragrunt-second-example/) - Additional patterns

### Testing with Dependencies

```go
func TestStack(t *testing.T) {
    t.Parallel()

    options := &terragrunt.Options{
        TerragruntDir: "../live/prod",
    }

    // Apply respects dependency order
    terragrunt.ApplyAll(t, options)
    defer terragrunt.DestroyAll(t, options)

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

    terragrunt.Init(t, options)
}
```

### Testing Stack Outputs

```go
func TestStackOutput(t *testing.T) {
    t.Parallel()

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

    terragrunt.StackRun(t, applyOpts)
    defer terragrunt.StackRun(t, destroyOpts)

    // Get specific output
    vpcID := terragrunt.StackOutput(t, options, "vpc_id")
    assert.NotEmpty(t, vpcID)

    // Get all outputs
    outputs := terragrunt.StackOutputAll(t, options)
    assert.Contains(t, outputs, "vpc_id")
}
```

### Checking Plan Exit Code

```go
func TestInfrastructureUpToDate(t *testing.T) {
    t.Parallel()

    options := &terragrunt.Options{
        TerragruntDir: "../prod",
    }

    // First apply
    terragrunt.ApplyAll(t, options)
    defer terragrunt.DestroyAll(t, options)

    // Plan should show no changes (exit code 0)
    exitCode := terragrunt.PlanAllExitCode(t, options)
    assert.Equal(t, 0, exitCode, "No changes expected")
}
```

### Using RunAll for Flexibility

```go
func TestCustomCommand(t *testing.T) {
    t.Parallel()

    options := &terragrunt.Options{
        TerragruntDir: "../modules",
    }

    // Run any terraform command with --all
    terragrunt.RunAll(t, options, "refresh")

    // Verify state is current
    output := terragrunt.RunAll(t, options, "show")
    assert.Contains(t, output, "expected-resource")
}
```

### Validating Stack Output Keys

```go
func TestStackOutputKeys(t *testing.T) {
    t.Parallel()

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

    terragrunt.StackRun(t, applyOpts)
    defer terragrunt.StackRun(t, destroyOpts)

    // Get list of all output keys
    keys := terragrunt.StackOutputListAll(t, options)

    // Verify required outputs exist
    assert.Contains(t, keys, "vpc_id")
    assert.Contains(t, keys, "subnet_ids")
}
```

### Using Filters (v0.97.0+)

```go
options := &terragrunt.Options{
    TerragruntDir:  "../live/prod",
    TerragruntArgs: []string{"--filter", "{./vpc}"},  // Only apply vpc
}
terragrunt.ApplyAll(t, options)
```

## Not Supported

This module does **NOT** support:
- `import`, `refresh`, `show`, `state`, `test` commands
- `backend`, `exec`, `catalog`, `scaffold` commands
- Discovery commands (`find`, `list`)
- Configuration commands (`info`)

For unsupported commands, run terragrunt directly via the `shell` module.

## Compatibility

Tested with Terragrunt v0.80.4+, v0.93.5+, and v0.99.x. Earlier versions may work but are not guaranteed.

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
