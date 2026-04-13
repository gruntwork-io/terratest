package slack_test

import (
	"os"
	"testing"
	"time"

	"github.com/slack-go/slack"
	"github.com/stretchr/testify/require"

	"github.com/gruntwork-io/terratest/modules/environment"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/retry"
	terratestslack "github.com/gruntwork-io/terratest/modules/slack"
)

const (
	slackTokenEnv     = "SLACK_TOKEN_FOR_TEST"
	slackChannelIDEnv = "SLACK_CHANNEL_ID_FOR_TEST"
)

func TestValidateSlackMessage(t *testing.T) {
	t.Parallel()

	environment.RequireEnvVar(t, slackTokenEnv)
	environment.RequireEnvVar(t, slackChannelIDEnv)

	token := os.Getenv(slackTokenEnv)
	channelID := os.Getenv(slackChannelIDEnv)

	uniqueID := random.UniqueID()
	msgTxt := "Test message from terratest: " + uniqueID

	slackClt := slack.New(token)

	_, _, err := slackClt.PostMessage(
		channelID,
		slack.MsgOptionText(msgTxt, false),
	)
	require.NoError(t, err)

	retry.DoWithRetryContext(
		t,
		t.Context(),
		"wait for slack message",
		10, 10*time.Second,
		func() (string, error) {
			err := terratestslack.ValidateExpectedSlackMessageE(t, token, channelID, msgTxt, 10, 5*time.Minute)
			return "", err
		},
	)
}
