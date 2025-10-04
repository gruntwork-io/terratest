// Package terragrunt provides testing helpers for Terragrunt, a thin wrapper for Terraform
// that provides extra tools for keeping configurations DRY, working with multiple modules,
// and managing remote state.
//
// This package contains functions for testing Terragrunt-specific features, including:
//   - run-all commands (ApplyAll, DestroyAll, PlanAllExitCode)
//   - stack operations (TgStackGenerate, TgStackRun, TgStackClean, TgStackOutput)
//   - standard terragrunt operations (TgInit)
//
// # When to Use This Package
//
// Use this package when:
//   - Testing Terragrunt stacks with multiple units
//   - Using run-all commands to apply/destroy across multiple modules
//   - Testing Terragrunt-specific features like dependency management
//
// Use the terraform package instead when:
//   - Testing a single Terraform module
//   - Not using Terragrunt-specific features
//
// # Function Naming Conventions
//
// This package uses two naming patterns:
//   - Functions prefixed with "Tg" (e.g., TgInit, TgStackGenerate) are older
//     functions maintained for backward compatibility
//   - Functions without prefix (e.g., ApplyAll, DestroyAll) are newer functions
//     introduced during the module refactoring
//
// All functions follow the standard Terratest pattern:
//   - Non-E suffix: Calls require.NoError and fails the test on error
//   - E suffix: Returns error for manual handling
//
// # Example Usage
//
//	func TestTerragruntStack(t *testing.T) {
//	    options := &terragrunt.Options{
//	        TerragruntBinary: "terragrunt",
//	        TerragruntDir:    "../examples/my-stack",
//	    }
//
//	    // Initialize the stack
//	    terragrunt.TgInit(t, options)
//
//	    // Apply all modules in the stack
//	    terragrunt.ApplyAll(t, options)
//
//	    // Cleanup
//	    defer terragrunt.DestroyAll(t, options)
//	}
package terragrunt
