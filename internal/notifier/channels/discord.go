package channels

import (
	"fmt"
	"os"

	"github.com/wcy-dt/ponghub/internal/types/structures/configure"
)

// DiscordNotifier implements Discord webhook notifications
type DiscordNotifier struct {
	config *configure.DiscordConfig
}

// NewDiscordNotifier creates a new Discord notifier
func NewDiscordNotifier(config *configure.DiscordConfig) *DiscordNotifier {
	return &DiscordNotifier{config: config}
}

// Send sends a Discord webhook notification
func (d *DiscordNotifier) Send(title, message string) error {
	webhookURL := d.config.WebhookURL
	if webhookURL == "" {
		webhookURL = os.Getenv("DISCORD_WEBHOOK_URL")
	}

	if webhookURL == "" {
		return fmt.Errorf("Discord webhook URL not configured")
	}

	payload := map[string]interface{}{
		"content": fmt.Sprintf("**%s**\n```\n%s\n```", title, message),
	}

	if d.config.Username != "" {
		payload["username"] = d.config.Username
	}
	if d.config.AvatarURL != "" {
		payload["avatar_url"] = d.config.AvatarURL
	}

	return sendWebhookRequest(webhookURL, payload)
}
