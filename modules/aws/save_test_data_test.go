package aws_test

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/aws/v2"
	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	gotesting "github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/gruntwork-io/terratest/modules/ssh/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type tStringLogger struct {
	sb strings.Builder
}

func (l *tStringLogger) Logf(t gotesting.TestingT, format string, args ...any) {
	t.Helper()
	fmt.Fprintf(&l.sb, format, args...)
	l.sb.WriteRune('\n')
}

// Not parallel: this test swaps the package-level logger.Default, which every other test in this package reads.
// Running it alongside them is a data race.
func TestSaveAndLoadEC2KeyPair(t *testing.T) {
	def, slogger := logger.Default, &tStringLogger{}
	logger.Default = logger.New(slogger)

	t.Cleanup(func() {
		logger.Default = def
	})

	keyPair, err := ssh.GenerateRSAKeyPairE(t, 2048) //nolint:mnd // RSA key size for testing
	require.NoError(t, err)

	ec2KeyPair := &aws.Ec2Keypair{
		KeyPair: keyPair,
		Name:    "test-ec2-key-pair",
		Region:  "us-east-1",
	}

	storedEC2KeyPair, err := json.Marshal(ec2KeyPair) //nolint:musttag // aws.Ec2Keypair does not have json tags
	require.NoError(t, err)

	tmpFolder := t.TempDir()
	aws.SaveEc2KeyPair(t, tmpFolder, ec2KeyPair)
	loadedEC2KeyPair := aws.LoadEc2KeyPair(t, tmpFolder)
	assert.Equal(t, ec2KeyPair, loadedEC2KeyPair)

	assert.NotContains(t, slogger.sb.String(), string(storedEC2KeyPair), "stored ec2 key pair should not be logged")
}
