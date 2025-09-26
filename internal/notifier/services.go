package notifier

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/smtp"
	"os"
	"time"

	"github.com/wcy-dt/ponghub/internal/types/structures/configure"
)

// NotificationService defines the interface for notification services
type NotificationService interface {
	Send(title, message string) error
}

// EmailNotifier implements email notifications
type EmailNotifier struct {
	config *configure.EmailConfig
}

// NewEmailNotifier creates a new email notifier
func NewEmailNotifier(config *configure.EmailConfig) *EmailNotifier {
	return &EmailNotifier{config: config}
}

// Send sends an email notification
func (e *EmailNotifier) Send(title, message string) error {
	// Get SMTP credentials from environment variables
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")

	if username == "" || password == "" {
		return fmt.Errorf("SMTP credentials not found in environment variables")
	}

	auth := smtp.PlainAuth("", username, password, e.config.SMTPHost)

	subject := title
	if e.config.Subject != "" {
		subject = e.config.Subject
	}

	body := fmt.Sprintf("Subject: %s\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s", subject, message)

	addr := fmt.Sprintf("%s:%d", e.config.SMTPHost, e.config.SMTPPort)
	return smtp.SendMail(addr, auth, e.config.From, e.config.To, []byte(body))
}

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

// WebhookNotifier implements generic webhook notifications
type WebhookNotifier struct {
	config *configure.WebhookConfig
}

// NewWebhookNotifier creates a new generic webhook notifier
func NewWebhookNotifier(config *configure.WebhookConfig) *WebhookNotifier {
	return &WebhookNotifier{config: config}
}

// Send sends a generic webhook notification
func (w *WebhookNotifier) Send(title, message string) error {
	url := w.config.URL
	if url == "" {
		url = os.Getenv("WEBHOOK_URL")
	}

	if url == "" {
		return fmt.Errorf("Webhook URL not configured")
	}

	payload := map[string]interface{}{
		"title":     title,
		"message":   message,
		"timestamp": time.Now().Format(time.RFC3339),
	}

	method := "POST"
	if w.config.Method != "" {
		method = w.config.Method
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal webhook payload: %v", err)
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to create webhook request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Add custom headers
	for key, value := range w.config.Headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send webhook request: %v", err)
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			fmt.Println("Error closing response body:", err)
		}
	}(resp.Body)

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("webhook request failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// sendWebhookRequest is a helper function to send JSON webhook requests
func sendWebhookRequest(url string, payload map[string]interface{}) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			fmt.Println("Error closing response body:", err)
		}
	}(resp.Body)

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
