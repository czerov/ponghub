package channels

import (
	"fmt"
	"os"

	"github.com/wcy-dt/ponghub/internal/types/structures/configure"
)

// SlackNotifier implements Slack webhook notifications
type SlackNotifier struct {
	config *configure.SlackConfig
}

// NewSlackNotifier creates a new Slack notifier
func NewSlackNotifier(config *configure.SlackConfig) *SlackNotifier {
	return &SlackNotifier{config: config}
}

// Send sends a Slack webhook notification
func (s *SlackNotifier) Send(title, message string) error {
	webhookURL := s.config.WebhookURL
	if webhookURL == "" {
		webhookURL = os.Getenv("SLACK_WEBHOOK_URL")
	}

	if webhookURL == "" {
		return fmt.Errorf("Slack webhook URL not configured")
	}

	payload := map[string]interface{}{
		"text": fmt.Sprintf("*%s*\n```%s```", title, message),
	}

	if s.config.Channel != "" {
		payload["channel"] = s.config.Channel
	}
	if s.config.Username != "" {
		payload["username"] = s.config.Username
	}
	if s.config.IconEmoji != "" {
		payload["icon_emoji"] = s.config.IconEmoji
	}

	return sendWebhookRequest(webhookURL, payload)
}
