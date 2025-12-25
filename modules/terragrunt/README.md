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
    "github.com/gruntwork-io/terratest/modules/terragrunt/v2"
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

### Non-Stack Commands
Work with standard terragrunt configurations (dependencies via `dependency` blocks):

- `Init(t, options)` - Initialize configuration
- `ApplyAll(t, options)` - Apply all modules with dependencies
- `DestroyAll(t, options)` - Destroy all modules with dependencies
- `PlanAllExitCode(t, options)` - Plan all and return exit code (0=no changes, 2=changes)
- `ValidateAll(t, options)` - Validate all modules
- `RunAll(t, options, command)` - Run any terraform command with --all flag
- `FormatAll(t, options)` - Format all terragrunt.hcl files

### Stack Commands
Work with `terragrunt.stack.hcl` configurations:

- `StackGenerate(t, options)` - Generate stack from stack.hcl
- `StackRun(t, options)` - Run command on generated stack
- `StackClean(t, options)` - Remove .terragrunt-stack directory
- `Output(t, options, key)` - Get stack output value
- `OutputJson(t, options, key)` - Get stack output as JSON
- `OutputAll(t, options)` - Get all stack outputs as map
- `OutputListAll(t, options)` - Get list of all output variable names

## Examples

### Testing with Dependencies

```go
func TestMultiModuleStack(t *testing.T) {
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

    terragrunt.ApplyAll(t, options)
    defer terragrunt.DestroyAll(t, options)

    // Get specific output
    vpcID := terragrunt.Output(t, options, "vpc_id")
    assert.NotEmpty(t, vpcID)

    // Get all outputs
    outputs := terragrunt.OutputAll(t, options)
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

### Validating Output Keys

```go
func TestOutputKeys(t *testing.T) {
    t.Parallel()

    options := &terragrunt.Options{
        TerragruntDir: "../stack",
    }

    terragrunt.ApplyAll(t, options)
    defer terragrunt.DestroyAll(t, options)

    // Get list of all output keys
    keys := terragrunt.OutputListAll(t, options)

    // Verify required outputs exist
    assert.Contains(t, keys, "vpc_id")
    assert.Contains(t, keys, "subnet_ids")
}
```

## Not Supported

This module does **NOT** support:
- Single-module commands (non-`--all` operations)
- `validate`, `graph`, `import`, `refresh`, `show`, `state`, `test` commands
- `backend`, `exec`, `catalog`, `scaffold` commands
- Discovery commands (`find`, `list`)
- Configuration commands (`dag`, `hcl`, `info`, `render`)

For single-module testing, consider using the `terraform` module instead, or run terragrunt commands directly via the `shell` module.

## Compatibility

Tested with Terragrunt v0.80.4+ and v0.93.5+. Earlier versions may work but are not guaranteed.

## More Info

- [Terragrunt Documentation](https://terragrunt.gruntwork.io/)
- [Terratest Documentation](https://terratest.gruntwork.io/)
