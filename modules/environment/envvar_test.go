package environment_test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/environment"
	"github.com/stretchr/testify/assert"
)

// MockT is used to test that the function under test will fail the test under certain circumstances.
type MockT struct {
	Failed bool
}

func (t *MockT) Fail() {
	t.Failed = true
}

func (t *MockT) FailNow() {
	t.Failed = true
}

func (t *MockT) Error(args ...any) {
	t.Failed = true
}

func (t *MockT) Errorf(format string, args ...any) {
	t.Failed = true
}

func (t *MockT) Fatal(args ...any) {
	t.Failed = true
}

func (t *MockT) Fatalf(format string, args ...any) {
	t.Failed = true
}

func (t *MockT) Helper() {}

func (t *MockT) Name() string {
	return "mockT"
}

// End MockT

var envvarList = []string{
	"TERRATEST_TEST_ENVIRONMENT",
	"TERRATESTTESTENVIRONMENT",
	"TERRATESTENVIRONMENT",
}

//nolint:paralleltest // These tests manipulate env vars and cannot run in parallel.
func TestGetFirstNonEmptyEnvVarOrEmptyStringChecksInOrder(t *testing.T) {
	t.Setenv("TERRATESTTESTENVIRONMENT", "test")
	t.Setenv("TERRATESTENVIRONMENT", "circleCI")

	value := environment.GetFirstNonEmptyEnvVarOrEmptyString(t, envvarList)
	assert.Equal(t, "test", value)
}

//nolint:paralleltest // These tests manipulate env vars and cannot run in parallel.
func TestGetFirstNonEmptyEnvVarOrEmptyStringReturnsEmpty(t *testing.T) {
	value := environment.GetFirstNonEmptyEnvVarOrEmptyString(t, envvarList)
	assert.Empty(t, value)
}

//nolint:paralleltest // These tests manipulate env vars and cannot run in parallel.
func TestRequireEnvVarFails(t *testing.T) {
	envVarName := "TERRATESTTESTENVIRONMENT"
	mockT := new(MockT)

	// Make sure the check fails when env var is not set
	environment.RequireEnvVar(mockT, envVarName)
	assert.True(t, mockT.Failed)
}

//nolint:paralleltest // These tests manipulate env vars and cannot run in parallel.
func TestRequireEnvVarPasses(t *testing.T) {
	envVarName := "TERRATESTTESTENVIRONMENT"

	// Make sure the check passes when env var is set
	t.Setenv(envVarName, "test")

	environment.RequireEnvVar(t, envVarName)
}
