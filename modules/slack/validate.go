package slack

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/slack-go/slack"
)

// ValidateExpectedSlackMessage validates whether a message containing the expected text was posted in the given channel
// ID, looking back historyLimit messages up to the given duration. For example, if you set (15*time.Minute) as the
// lookBack parameter with historyLimit set to 50, then this will look back the last 50 messages, up to 15 minutes ago.
// This expects a slack token to be provided. This returns MessageNotFoundErr when there is no match.
// For the purposes of matching, this only checks the following blocks:
// - Section block text
// - Header block text
// All other blocks are ignored in the validation.
// NOTE: This only looks for bot posted messages.
func ValidateExpectedSlackMessageE(
	t testing.TestingT,
	token,
	channelID,
	expectedText string,
	historyLimit int,
	lookBack time.Duration,
) error {
	lookBackTime := time.Now().Add(-1 * lookBack)
	slackClt := slack.New(token)
	params := slack.GetConversationHistoryParameters{
		ChannelID: channelID,
		Limit:     historyLimit,
		Oldest:    strconv.FormatInt(lookBackTime.Unix(), 10),
	}

	resp, err := slackClt.GetConversationHistory(&params)
	if err != nil {
		return err
	}

	for _, msg := range resp.Messages {
		if checkMessageContainsText(msg.Msg, expectedText) {
			return nil
		}

		if msg.SubMessage != nil {
			if checkMessageContainsText(*msg.SubMessage, expectedText) {
				return nil
			}
		}
	}
	return fmt.Errorf("still no message")
}

func checkMessageContainsText(msg slack.Msg, expectedText string) bool {
	// If this message is not a bot message, ignore.
	if msg.Type != slack.MsgSubTypeBotMessage && msg.BotID == "" {
		return false
	}

	// Check message text
	if strings.Contains(msg.Text, expectedText) {
		return true
	}

	// Check attachments
	for _, attachment := range msg.Attachments {
		if strings.Contains(attachment.Text, expectedText) {
			return true
		}
	}

	// Check blocks
	for _, block := range msg.Blocks.BlockSet {
		switch block.BlockType() {
		case slack.MBTSection:
			sectionBlk := block.(*slack.SectionBlock)
			if sectionBlk.Text != nil && strings.Contains(sectionBlk.Text.Text, expectedText) {
				return true
			}
		case slack.MBTHeader:
			headerBlk := block.(*slack.HeaderBlock)
			if headerBlk.Text != nil && strings.Contains(headerBlk.Text.Text, expectedText) {
				return true
			}
		}
	}

	return false
}

type MessageNotFoundErr struct{}

func (err MessageNotFoundErr) Error() string {
	return "Could not find the expected text in any of the messages posted in the given channel."
}
