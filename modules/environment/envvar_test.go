package environment

import (
	"os"
	"testing"

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

func (t *MockT) Error(args ...interface{}) {
	t.Failed = true
}

func (t *MockT) Errorf(format string, args ...interface{}) {
	t.Failed = true
}

func (t *MockT) Fatal(args ...interface{}) {
	t.Failed = true
}

func (t *MockT) Fatalf(format string, args ...interface{}) {
	t.Failed = true
}

func (t *MockT) Name() string {
	return "mockT"
}

// End MockT

var envvarList = []string{
	"TERRATEST_TEST_ENVIRONMENT",
	"TERRATESTTESTENVIRONMENT",
	"TERRATESTENVIRONMENT",
}

func TestGetFirstNonEmptyEnvVarOrEmptyStringChecksInOrder(t *testing.T) {
	// These tests can not run in parallel, since they manipulate env vars
	// DO NOT ADD THIS: t.Parallel()

	os.Setenv("TERRATESTTESTENVIRONMENT", "test")
	os.Setenv("TERRATESTENVIRONMENT", "circleCI")
	defer os.Setenv("TERRATESTTESTENVIRONMENT", "")
	defer os.Setenv("TERRATESTENVIRONMENT", "")
	value := GetFirstNonEmptyEnvVarOrEmptyString(t, envvarList)
	assert.Equal(t, value, "test")
}

func TestGetFirstNonEmptyEnvVarOrEmptyStringReturnsEmpty(t *testing.T) {
	// These tests can not run in parallel, since they manipulate env vars
	// DO NOT ADD THIS: t.Parallel()

	value := GetFirstNonEmptyEnvVarOrEmptyString(t, envvarList)
	assert.Equal(t, value, "")
}

func TestRequireEnvVarFails(t *testing.T) {
	// These tests can not run in parallel, since they manipulate env vars
	// DO NOT ADD THIS: t.Parallel()

	envVarName := "TERRATESTTESTENVIRONMENT"
	mockT := new(MockT)

	// Make sure the check fails when env var is not set
	RequireEnvVar(mockT, envVarName)
	assert.True(t, mockT.Failed)
}

func TestRequireEnvVarPasses(t *testing.T) {
	// These tests can not run in parallel, since they manipulate env vars
	// DO NOT ADD THIS: t.Parallel()

	envVarName := "TERRATESTTESTENVIRONMENT"

	// Make sure the check passes when env var is set
	os.Setenv(envVarName, "test")
	defer os.Setenv(envVarName, "")
	RequireEnvVar(t, envVarName)
}
