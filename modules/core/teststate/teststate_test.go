package teststate_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	gotesting "github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/gruntwork-io/terratest/modules/core/v2/teststate"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// tStringLogger captures everything written to the logger so a test can assert on what was, and was not, logged.
type tStringLogger struct {
	sb strings.Builder
}

func (l *tStringLogger) Logf(t gotesting.TestingT, format string, args ...any) {
	t.Helper()
	fmt.Fprintf(&l.sb, format, args...)
	l.sb.WriteRune('\n')
}

// captureLog swaps logger.Default for the duration of the test. Tests using it must not call t.Parallel.
func captureLog(t *testing.T) *tStringLogger {
	t.Helper()

	def, slogger := logger.Default, &tStringLogger{}
	logger.Default = logger.New(slogger)

	t.Cleanup(func() { logger.Default = def })

	return slogger
}

func TestFormatPath(t *testing.T) {
	t.Parallel()

	assert.Equal(t, filepath.Join("/foo", ".test-data", "Bar.json"), teststate.FormatPath("/foo", "Bar.json"))
}

func TestSaveAndLoadRoundTrip(t *testing.T) {
	t.Parallel()

	type payload struct {
		Name  string
		Count int
	}

	path := teststate.FormatPath(t.TempDir(), "Payload.json")
	expected := payload{Name: "terratest", Count: 3}

	teststate.Save(t, path, true, expected)

	var actual payload

	teststate.Load(t, path, &actual)
	assert.Equal(t, expected, actual)
}

func TestSaveOverwriteSemantics(t *testing.T) {
	t.Parallel()

	path := teststate.FormatPath(t.TempDir(), "Value.json")

	teststate.Save(t, path, true, "first")

	// overwrite=false must leave the existing value alone.
	teststate.Save(t, path, false, "second")

	var got string

	teststate.Load(t, path, &got)
	assert.Equal(t, "first", got, "overwrite=false must not clobber an existing value")

	// overwrite=true must replace it.
	teststate.Save(t, path, true, "third")
	teststate.Load(t, path, &got)
	assert.Equal(t, "third", got)
}

// TestSaveLogsValueAndSaveRedactedDoesNot is the behavioural contract that callers holding secrets depend on.
func TestSaveLogsValueAndSaveRedactedDoesNot(t *testing.T) {
	const secret = "-----BEGIN RSA PRIVATE KEY-----sentinel-----END RSA PRIVATE KEY-----"

	t.Run("Save logs the marshalled value", func(t *testing.T) {
		slogger := captureLog(t)
		teststate.Save(t, teststate.FormatPath(t.TempDir(), "Plain.json"), true, secret)
		assert.Contains(t, slogger.sb.String(), secret, "Save is expected to log the marshalled JSON")
	})

	t.Run("SaveRedacted does not", func(t *testing.T) {
		slogger := captureLog(t)
		teststate.SaveRedacted(t, teststate.FormatPath(t.TempDir(), "Secret.json"), true, secret)
		assert.NotContains(t, slogger.sb.String(), secret, "SaveRedacted must not log the marshalled JSON")
	})
}

// TestSaveRedactedDoesNotLeakViaOverwriteWarning covers the second log statement in the save path. The overwrite
// warning renders the value with %v, which is not suppressed by the redacted flag, so a redacted save over an
// existing file could still leak. Values reaching SaveRedacted are secrets by definition.
func TestSaveRedactedDoesNotLeakViaOverwriteWarning(t *testing.T) {
	const secret = "-----BEGIN RSA PRIVATE KEY-----sentinel-----END RSA PRIVATE KEY-----"

	path := teststate.FormatPath(t.TempDir(), "Secret.json")

	// First save creates the file, so the second one takes the overwrite-warning branch.
	teststate.SaveRedacted(t, path, true, secret)

	slogger := captureLog(t)
	teststate.SaveRedacted(t, path, true, secret)

	assert.NotContains(t, slogger.sb.String(), secret,
		"the overwrite warning must not render a redacted value")
}

func TestIsPresent(t *testing.T) {
	t.Parallel()

	path := teststate.FormatPath(t.TempDir(), "Maybe.json")
	assert.False(t, teststate.IsPresent(t, path), "a missing file is not present")

	teststate.Save(t, path, true, "value")
	assert.True(t, teststate.IsPresent(t, path))

	// An empty JSON value counts as absent, so a stage can re-create it.
	teststate.Save(t, path, true, "")
	assert.False(t, teststate.IsPresent(t, path), "an empty value counts as absent")
}

func TestIsEmptyJSON(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name  string
		bytes string
		empty bool
	}{
		{"no bytes", "", true},
		{"null", "null", true},
		{"false", "false", true},
		{"zero", "0", true},
		{"empty string", `""`, true},
		{"empty array", "[]", true},
		{"empty object", "{}", true},
		{"true", "true", false},
		{"non-zero", "42", false},
		{"string", `"value"`, false},
		{"array", `[1]`, false},
		{"object", `{"a":1}`, false},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, testCase.empty, teststate.IsEmptyJSON(t, []byte(testCase.bytes)))
		})
	}
}

func TestCleanup(t *testing.T) {
	t.Parallel()

	folder := t.TempDir()
	path := teststate.FormatPath(folder, "Value.json")

	teststate.Save(t, path, true, "value")
	require.FileExists(t, path)

	teststate.Cleanup(t, path)
	assert.NoFileExists(t, path)

	// Cleaning up an already-absent path is a no-op, not a failure.
	assert.NotPanics(t, func() { teststate.Cleanup(t, path) })
}

func TestCleanupFolder(t *testing.T) {
	t.Parallel()

	folder := t.TempDir()
	teststate.Save(t, teststate.FormatPath(folder, "One.json"), true, "1")
	teststate.Save(t, teststate.FormatPath(folder, "Two.json"), true, "2")

	require.NoError(t, teststate.CleanupFolderE(t, folder))
	assert.NoDirExists(t, filepath.Join(folder, ".test-data"))

	// Cleaning an absent folder is a no-op.
	require.NoError(t, teststate.CleanupFolderE(t, folder))
}

// TestSaveWritesOwnerOnlyPermissions pins the file mode. Saved state can hold private keys, so it must not be
// world readable.
func TestSaveWritesOwnerOnlyPermissions(t *testing.T) {
	t.Parallel()

	path := teststate.FormatPath(t.TempDir(), "Value.json")
	teststate.Save(t, path, true, "value")

	info, err := os.Stat(path)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0o600), info.Mode().Perm(), "saved test data must be owner read/write only")
}
