package terraform

import (
	"fmt"
	"reflect"
)

// TgInvalidBinary occurs when a terragrunt function is called and the TerraformBinary is
// set to a value other than terragrunt
type TgInvalidBinary string

func (err TgInvalidBinary) Error() string {
	return fmt.Sprintf("terragrunt must be set as TerraformBinary to use this function. [ TerraformBinary : %s ]", string(err))
}

// OutputKeyNotFound occurs when terraform output does not contain a value for the key
// specified in the function call
type OutputKeyNotFound string

func (err OutputKeyNotFound) Error() string {
	return fmt.Sprintf("output doesn't contain a value for the key %q", string(err))
}

// OutputValueNotMap occures when casting a found output value to a map fails
type OutputValueNotMap struct {
	Value interface{}
}

func (err OutputValueNotMap) Error() string {
	return fmt.Sprintf("Output value %q is not a map", err.Value)
}

// OutputValueNotList occurs when casting a found output value to a
// list of interfaces fails
type OutputValueNotList struct {
	Value interface{}
}

func (err OutputValueNotList) Error() string {
	return fmt.Sprintf("Output value %q is not a list", err.Value)
}

// EmptyOutput is an error that occurs when an output is empty.
type EmptyOutput string

func (outputName EmptyOutput) Error() string {
	return fmt.Sprintf("Required output %s was empty", string(outputName))
}

// UnexpectedOutputType is an error that occurs when the output is not of the type we expect
type UnexpectedOutputType struct {
	Key          string
	ExpectedType string
	ActualType   string
}

func (err UnexpectedOutputType) Error() string {
	return fmt.Sprintf("Expected output '%s' to be of type '%s' but got '%s'", err.Key, err.ExpectedType, err.ActualType)
}

// VarFileNotFound is an error that occurs when a var file cannot be found in an option's VarFile list
type VarFileNotFound struct {
	Path string
}

func (err VarFileNotFound) Error() string {
	return fmt.Sprintf("Var file '%s' not found", err.Path)
}

// InputFileKeyNotFound occurs when tfvar file does not contain a value for the key
// specified in the function call
type InputFileKeyNotFound struct {
	FilePath string
	Key      string
}

func (err InputFileKeyNotFound) Error() string {
	return fmt.Sprintf("tfvar file %q doesn't contain a value for the key %q", err.FilePath, err.Key)
}

// PanicWhileParsingVarFile is returned when the HCL parsing routine panics due to errors.
type PanicWhileParsingVarFile struct {
	ConfigFile     string
	RecoveredValue interface{}
}

func (err PanicWhileParsingVarFile) Error() string {
	return fmt.Sprintf("Recovering panic while parsing '%s'. Got error of type '%v': %v", err.ConfigFile, reflect.TypeOf(err.RecoveredValue), err.RecoveredValue)
}

// UnsupportedDefaultWorkspaceDeletion is returned when user tries to delete the workspace "default"
type UnsupportedDefaultWorkspaceDeletion struct{}

func (err *UnsupportedDefaultWorkspaceDeletion) Error() string {
	return "Deleting the workspace 'default' is not supported"
}

// WorkspaceDoesNotExist is returned when user tries to delete a workspace which does not exist
type WorkspaceDoesNotExist string

func (err WorkspaceDoesNotExist) Error() string {
	return fmt.Sprintf("The workspace %q does not exist.", string(err))
}
