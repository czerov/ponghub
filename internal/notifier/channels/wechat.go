package channels

import (
	"fmt"
	"os"

	"github.com/wcy-dt/ponghub/internal/types/structures/configure"
)

// WeChatNotifier implements WeChat Work webhook notifications
type WeChatNotifier struct {
	config *configure.WeChatConfig
}

// NewWeChatNotifier creates a new WeChat notifier
func NewWeChatNotifier(config *configure.WeChatConfig) *WeChatNotifier {
	return &WeChatNotifier{config: config}
}

// Send sends a WeChat Work webhook notification
func (w *WeChatNotifier) Send(title, message string) error {
	webhookURL := w.config.WebhookURL
	if webhookURL == "" {
		webhookURL = os.Getenv("WECHAT_WEBHOOK_URL")
	}

	if webhookURL == "" {
		return fmt.Errorf("WeChat webhook URL not configured")
	}

	payload := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]string{
			"content": fmt.Sprintf("%s\n%s", title, message),
		},
	}

	return sendWebhookRequest(webhookURL, payload)
}
