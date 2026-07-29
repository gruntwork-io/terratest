package ssh_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	gotesting "github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/gruntwork-io/terratest/modules/ssh/v2"
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

func TestSaveAndLoadSSHKeyPair(t *testing.T) {
	t.Parallel()

	expectedData, err := ssh.GenerateRSAKeyPairE(t, 2048) //nolint:mnd // RSA key size for testing
	require.NoError(t, err)

	tmpFolder := t.TempDir()
	ssh.SaveSSHKeyPair(t, tmpFolder, expectedData)

	actualData := ssh.LoadSSHKeyPair(t, tmpFolder)
	assert.Equal(t, expectedData, actualData)
}

// TestSaveSSHKeyPairDoesNotLogPrivateKey pins that SaveSSHKeyPair redacts the marshalled value. KeyPair.PrivateKey
// is a PEM private key, and the test log is routinely captured by CI, so it must never appear there.
//
// Not parallel: this test swaps the package-level logger.Default, which other tests in this package read.
func TestSaveSSHKeyPairDoesNotLogPrivateKey(t *testing.T) {
	def, slogger := logger.Default, &tStringLogger{}
	logger.Default = logger.New(slogger)

	t.Cleanup(func() {
		logger.Default = def
	})

	keyPair, err := ssh.GenerateRSAKeyPairE(t, 2048) //nolint:mnd // RSA key size for testing
	require.NoError(t, err)

	tmpFolder := t.TempDir()
	ssh.SaveSSHKeyPair(t, tmpFolder, keyPair)

	logged := slogger.sb.String()
	assert.NotContains(t, logged, keyPair.PrivateKey, "the private key must never be logged")
	assert.NotContains(t, logged, "PRIVATE KEY", "no PEM block should reach the log")

	// Confirm the logger really was wired up, so the assertions above are not vacuous.
	assert.Contains(t, logged, "Storing test data in", "the save operation should still be logged")

	// The key pair must still round-trip; redaction applies to the log, not to the file.
	assert.Equal(t, keyPair, ssh.LoadSSHKeyPair(t, tmpFolder))
}
