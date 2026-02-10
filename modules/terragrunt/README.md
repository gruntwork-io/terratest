# Terragrunt Module

Testing library for Terragrunt configurations in Go. Provides helpers for running Terragrunt commands across multiple modules (run-all) and stack-based workflows.

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

```go
import (
    "testing"
    "github.com/gruntwork-io/terratest/modules/terragrunt"
)

func TestTerragruntApply(t *testing.T) {
    t.Parallel()

    options := &terragrunt.Options{
        TerragruntDir: "../path/to/terragrunt/config",
    }

    // Apply all modules
    terragrunt.ApplyAll(t, options)

    // Clean up
    defer terragrunt.DestroyAll(t, options)
}
```

## Key Concepts

### Options Struct

The `Options` struct has two distinct parts:

1. **Test Framework Configuration** (NOT passed to terragrunt CLI):
   - `TerragruntDir` - where to run terragrunt
   - `TerragruntBinary` - binary name (default: "terragrunt")
   - `EnvVars` - environment variables
   - `Logger`, `MaxRetries`, `TimeBetweenRetries` - test framework settings

2. **Command-Line Arguments** (passed to terragrunt):
   - `TerragruntArgs` - global terragrunt flags (e.g., `--log-level`, `--no-color`)
   - `TerraformArgs` - command-specific terraform flags (e.g., `-upgrade`)

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

### Run --all Commands

Work with [implicit stacks](https://terragrunt.gruntwork.io/docs/features/stacks/#implicit-stacks) (multiple units in a directory):

- `Init(t, options)` - Initialize configuration
- `ApplyAll(t, options)` - Apply all modules with dependencies
- `DestroyAll(t, options)` - Destroy all modules with dependencies
- `PlanAllExitCode(t, options)` - Plan all and return exit code (0=no changes, 2=changes, other=error)
- `ValidateAll(t, options)` - Validate all modules
- `RunAll(t, options, command)` - Run any terraform command with --all flag
- `OutputAllJson(t, options)` - Get all outputs as raw JSON string (note: returns separate JSON objects per module)
- `FormatAll(t, options)` - Format all terragrunt.hcl files
- `HclValidate(t, options)` - Validate terragrunt.hcl syntax and configuration

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

    terragrunt.StackRun(t, &terragrunt.Options{
        TerragruntDir: "../stack",
        TerraformArgs: []string{"apply"},
    })
    defer terragrunt.StackRun(t, &terragrunt.Options{
        TerragruntDir: "../stack",
        TerraformArgs: []string{"destroy"},
    })

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

    terragrunt.StackRun(t, &terragrunt.Options{
        TerragruntDir: "../stack",
        TerraformArgs: []string{"apply"},
    })
    defer terragrunt.StackRun(t, &terragrunt.Options{
        TerragruntDir: "../stack",
        TerraformArgs: []string{"destroy"},
    })

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
- Single-unit commands (non-`--all` operations)
- `validate`, `graph`, `import`, `refresh`, `show`, `state`, `test` commands
- `backend`, `exec`, `catalog`, `scaffold` commands
- Discovery commands (`find`, `list`)
- Configuration commands (`dag`, `hcl`, `info`, `render`)

For single-unit testing, consider using the `terraform` module instead, or run terragrunt commands directly via the `shell` module.

## Compatibility

Tested with Terragrunt v0.80.4+, v0.93.5+, and v0.99.x. Earlier versions may work but are not guaranteed.

### Migration from terraform Module

| Deprecated (terraform module) | Replacement (terragrunt module) |
|-------------------------------|----------------------------------|
| `TgApplyAll` / `TgApplyAllE` | `ApplyAll` / `ApplyAllE` |
| `TgDestroyAll` / `TgDestroyAllE` | `DestroyAll` / `DestroyAllE` |
| `TgPlanAllExitCode` / `TgPlanAllExitCodeE` | `PlanAllExitCode` / `PlanAllExitCodeE` |
| `ValidateInputs` / `ValidateInputsE` | `HclValidate` / `HclValidateE` |

## More Info

- [Terragrunt Documentation](https://terragrunt.gruntwork.io/)
- [Terratest Documentation](https://terratest.gruntwork.io/)
