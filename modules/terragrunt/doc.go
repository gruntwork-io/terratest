// Package terragrunt provides test helpers for running Terragrunt commands.
//
// This package wraps the Terragrunt CLI to simplify integration testing of
// Terragrunt configurations. It supports both single-unit testing and
// multi-unit stack testing with dependency management.
//
// For single-unit testing, you can use either this package or the terraform
// package with TerraformBinary set to "terragrunt". For stack testing with
// --all commands, use the dedicated functions in this package such as
// [ApplyAllContextE] and [DestroyAllContextE].
package terragrunt
