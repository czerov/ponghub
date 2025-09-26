package channels

import (
	"fmt"
	"os"

	"github.com/wcy-dt/ponghub/internal/types/structures/configure"
)

// TelegramNotifier implements Telegram bot notifications
type TelegramNotifier struct {
	config *configure.TelegramConfig
}

// NewTelegramNotifier creates a new Telegram notifier
func NewTelegramNotifier(config *configure.TelegramConfig) *TelegramNotifier {
	return &TelegramNotifier{config: config}
}

// Send sends a Telegram bot notification
func (t *TelegramNotifier) Send(title, message string) error {
	botToken := t.config.BotToken
	if botToken == "" {
		botToken = os.Getenv("TELEGRAM_BOT_TOKEN")
	}

	chatID := t.config.ChatID
	if chatID == "" {
		chatID = os.Getenv("TELEGRAM_CHAT_ID")
	}

	if botToken == "" || chatID == "" {
		return fmt.Errorf("Telegram bot token or chat ID not configured")
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)
	text := fmt.Sprintf("*%s*\n```\n%s\n```", title, message)

	payload := map[string]interface{}{
		"chat_id":    chatID,
		"text":       text,
		"parse_mode": "Markdown",
	}

	return sendWebhookRequest(url, payload)
}
