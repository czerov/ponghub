package channels

import (
	"fmt"
	"os"
	"time"

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

// Send sends a WeChat Work webhook notification with enhanced features
func (w *WeChatNotifier) Send(title, message string) error {
	webhookURL := w.config.WebhookURL
	if webhookURL == "" {
		webhookURL = os.Getenv("WECHAT_WEBHOOK_URL")
	}

	if webhookURL == "" {
		return fmt.Errorf("WeChat webhook URL not configured")
	}

	payload := w.buildPayload(title, message)

	// Execute request with retry logic
	return w.sendWithRetry(webhookURL, payload)
}

// buildPayload constructs the WeChat API payload based on message type
func (w *WeChatNotifier) buildPayload(title, message string) map[string]interface{} {
	msgType := w.config.MsgType
	if msgType == "" {
		msgType = "text" // Default to text
	}

	switch msgType {
	case "markdown":
		return w.buildMarkdownPayload(title, message)
	case "text":
		return w.buildTextPayload(title, message)
	default:
		return w.buildTextPayload(title, message) // Default to text
	}
}

// buildTextPayload constructs a text message payload for WeChat Work
func (w *WeChatNotifier) buildTextPayload(title, message string) map[string]interface{} {
	content := fmt.Sprintf("%s\n\n%s", title, message)

	// Add timestamp
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	content += fmt.Sprintf("\n\n⏰ 时间: %s", timestamp)

	payload := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]interface{}{
			"content": content,
		},
	}

	// Add mentions if configured
	if len(w.config.Mentions) > 0 {
		mentionedList := make([]string, 0, len(w.config.Mentions))
		mentionedMobileList := make([]string, 0, len(w.config.Mentions))

		for _, mention := range w.config.Mentions {
			if mention != "" {
				// Check if it's a mobile number (simple check)
				if len(mention) == 11 && mention[0] == '1' {
					mentionedMobileList = append(mentionedMobileList, mention)
				} else {
					mentionedList = append(mentionedList, mention)
				}
			}
		}

		if len(mentionedList) > 0 || len(mentionedMobileList) > 0 {
			textObj := payload["text"].(map[string]interface{})
			if len(mentionedList) > 0 {
				textObj["mentioned_list"] = mentionedList
			}
			if len(mentionedMobileList) > 0 {
				textObj["mentioned_mobile_list"] = mentionedMobileList
			}
		}
	}

	return payload
}

// buildMarkdownPayload constructs a markdown message payload for WeChat Work
func (w *WeChatNotifier) buildMarkdownPayload(title, message string) map[string]interface{} {
	// Format content as markdown
	content := fmt.Sprintf("## %s\n\n", title)
	content += fmt.Sprintf("```\n%s\n```\n\n", message)

	// Add timestamp
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	content += fmt.Sprintf("⏰ **时间**: %s", timestamp)

	payload := map[string]interface{}{
		"msgtype": "markdown",
		"markdown": map[string]interface{}{
			"content": content,
		},
	}

	return payload
}

// sendWithRetry sends the WeChat webhook with retry logic
func (w *WeChatNotifier) sendWithRetry(webhookURL string, payload map[string]interface{}) error {
	maxRetries := 0
	if w.config.Retries > 0 {
		maxRetries = w.config.Retries
	}

	timeout := 30
	if w.config.Timeout > 0 {
		timeout = w.config.Timeout
	}

	headers := make(map[string]string)
	if w.config.UserAgent != "" {
		headers["User-Agent"] = w.config.UserAgent
	} else {
		headers["User-Agent"] = "PongHub-WeChat-Notifier/1.0"
	}

	return sendHTTPRequest(webhookURL, "POST", payload, headers, maxRetries, timeout, false)
}

// ValidateConfig validates the WeChat configuration
func (w *WeChatNotifier) ValidateConfig() error {
	if w.config.WebhookURL == "" && os.Getenv("WECHAT_WEBHOOK_URL") == "" {
		return fmt.Errorf("WeChat webhook URL not configured")
	}

	// Validate message type
	validMsgTypes := []string{"", "text", "markdown"}
	valid := false
	for _, msgType := range validMsgTypes {
		if w.config.MsgType == msgType {
			valid = true
			break
		}
	}
	if !valid {
		return fmt.Errorf("invalid message type: %s", w.config.MsgType)
	}

	return nil
}
