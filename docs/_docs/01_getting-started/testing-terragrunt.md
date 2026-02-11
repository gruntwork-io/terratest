---
layout: collection-browser-doc
title: Testing Terragrunt
category: getting-started
excerpt: >-
  Learn how to test Terragrunt configurations with Terratest.
tags: ["terragrunt", "testing", "quick-start"]
order: 104
nav_title: Documentation
nav_title_link: /docs/
---

## Overview

Terratest provides two approaches for testing Terragrunt configurations:

| Approach | Use Case | Package |
|----------|----------|---------|
| **Unit** | Testing individual units | `modules/terraform` with `TerraformBinary: "terragrunt"` |
| **Stack** | Testing a stack of units with `--all` commands | `modules/terragrunt` |

## Unit Testing

For testing a single Terragrunt unit, use the `terraform` package with `TerraformBinary` set to `"terragrunt"`:

```go
func TestTerragruntModule(t *testing.T) {
    t.Parallel()

    terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
        TerraformDir:    "../examples/my-module",
        TerraformBinary: "terragrunt",
    })

    defer terraform.Destroy(t, terraformOptions)
    terraform.Apply(t, terraformOptions)

    output := terraform.Output(t, terraformOptions, "my_output")
    assert.Equal(t, "expected_value", output)
}
```

## Stack Testing

For testing a stack of units with dependencies, use the dedicated `terragrunt` package:

```go
func TestStack(t *testing.T) {
    t.Parallel()

    testFolder, err := files.CopyTerragruntFolderToTemp("../live/prod", t.Name())
    require.NoError(t, err)

    options := &terragrunt.Options{
        TerragruntDir: testFolder,
    }

    defer terragrunt.DestroyAll(t, options)
    terragrunt.ApplyAll(t, options)

    exitCode := terragrunt.PlanAllExitCode(t, options)
    require.Equal(t, 0, exitCode)
}
```

### Available Functions

| Function | Description |
|----------|-------------|
| `Init` | Run `terragrunt init` |
| `ApplyAll` | Run `terragrunt apply --all` |
| `DestroyAll` | Run `terragrunt destroy --all` |
| `PlanAllExitCode` | Run `terragrunt plan --all`, return exit code |
| `ValidateAll` | Run `terragrunt validate --all` |
| `FormatAll` | Run `terragrunt fmt --all` |
| `RunAll` | Run any command with `--all` flag |

### Stack Functions

For Terragrunt [stacks](https://terragrunt.gruntwork.io/docs/features/stacks/):

| Function | Description |
|----------|-------------|
| `StackGenerate` | Generate stack from `terragrunt.stack.hcl` |
| `StackRun` | Run command on generated stack |
| `StackClean` | Remove `.terragrunt-stack` directory |
| `Output` | Get stack output value |
| `OutputAll` | Get all stack outputs as map |

## Further Reading

- [Terragrunt Documentation](https://terragrunt.gruntwork.io/)
- [Stack example](https://github.com/gruntwork-io/terratest/tree/main/examples/terragrunt-multi-module-example)
- [terragrunt package reference](https://pkg.go.dev/github.com/gruntwork-io/terratest/modules/terragrunt)
