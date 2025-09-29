package channels

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/wcy-dt/ponghub/internal/types/structures/configure"
)

// WebhookNotifier implements generic webhook notifications
type WebhookNotifier struct {
	config *configure.WebhookConfig
}

// NewWebhookNotifier creates a new generic webhook notifier
func NewWebhookNotifier(config *configure.WebhookConfig) *WebhookNotifier {
	return &WebhookNotifier{config: config}
}

// Send sends a generic webhook notification with enhanced configuration support
func (w *WebhookNotifier) Send(title, message string) error {
	url := w.config.URL
	if url == "" {
		url = os.Getenv("WEBHOOK_URL")
	}

	if url == "" {
		return fmt.Errorf("webhook URL not configured")
	}

	// Prepare the payload
	payload, contentType, err := w.buildPayload(title, message)
	if err != nil {
		return fmt.Errorf("failed to build webhook payload: %v", err)
	}

	// Determine method
	method := "POST"
	if w.config.Method != "" {
		method = strings.ToUpper(w.config.Method)
	}

	// Prepare headers
	headers := make(map[string]string)
	for key, value := range w.config.Headers {
		headers[key] = value
	}

	// Set authentication if configured
	if w.config.AuthType != "" {
		w.setAuthentication(headers)
	}

	// Execute request with retry logic
	return w.sendWithRetry(url, method, payload, contentType, headers)
}

// buildPayload constructs the webhook payload based on configuration
func (w *WebhookNotifier) buildPayload(title, message string) (interface{}, string, error) {
	data := map[string]interface{}{
		"title":     title,
		"message":   message,
		"timestamp": time.Now().Format(time.RFC3339),
		"service":   "ponghub",
	}

	// Use custom template if configured
	if w.config.Template != "" {
		return w.buildTemplatePayload(data)
	}

	// Use predefined format if configured
	if w.config.Format != "" {
		return w.buildFormattedPayload(data)
	}

	// Default JSON payload
	return data, "application/json", nil
}

// buildTemplatePayload builds payload using custom template
func (w *WebhookNotifier) buildTemplatePayload(data map[string]interface{}) (interface{}, string, error) {
	tmpl, err := template.New("webhook").Parse(w.config.Template)
	if err != nil {
		return nil, "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, "", fmt.Errorf("failed to execute template: %w", err)
	}

	// Try to parse as JSON first
	var jsonData interface{}
	if err := json.Unmarshal(buf.Bytes(), &jsonData); err == nil {
		return jsonData, "application/json", nil
	}

	// Return as string if not valid JSON
	contentType := "text/plain"
	if w.config.ContentType != "" {
		contentType = w.config.ContentType
	}
	return buf.String(), contentType, nil
}

// buildFormattedPayload builds payload using predefined format
func (w *WebhookNotifier) buildFormattedPayload(data map[string]interface{}) (interface{}, string, error) {
	switch strings.ToLower(w.config.Format) {
	case "slack":
		return w.buildSlackFormat(data), "application/json", nil
	case "discord":
		return w.buildDiscordFormat(data), "application/json", nil
	case "teams":
		return w.buildTeamsFormat(data), "application/json", nil
	case "mattermost":
		return w.buildMattermostFormat(data), "application/json", nil
	default:
		return data, "application/json", nil
	}
}

// buildSlackFormat builds Slack-compatible payload
func (w *WebhookNotifier) buildSlackFormat(data map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"text": fmt.Sprintf("*%s*", data["title"]),
		"attachments": []map[string]interface{}{
			{
				"color":     "danger",
				"text":      data["message"],
				"timestamp": time.Now().Unix(),
				"fields": []map[string]interface{}{
					{
						"title": "Service",
						"value": data["service"],
						"short": true,
					},
				},
			},
		},
	}
}

// buildDiscordFormat builds Discord-compatible payload
func (w *WebhookNotifier) buildDiscordFormat(data map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"embeds": []map[string]interface{}{
			{
				"title":       data["title"],
				"description": data["message"],
				"color":       0xFF0000, // Red
				"timestamp":   data["timestamp"],
				"fields": []map[string]interface{}{
					{
						"name":   "Service",
						"value":  data["service"],
						"inline": true,
					},
				},
			},
		},
	}
}

// buildTeamsFormat builds Microsoft Teams-compatible payload
func (w *WebhookNotifier) buildTeamsFormat(data map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"@type":      "MessageCard",
		"@context":   "http://schema.org/extensions",
		"themeColor": "FF0000",
		"summary":    data["title"],
		"sections": []map[string]interface{}{
			{
				"activityTitle": data["title"],
				"activityText":  data["message"],
				"facts": []map[string]interface{}{
					{
						"name":  "Service",
						"value": data["service"],
					},
					{
						"name":  "Timestamp",
						"value": data["timestamp"],
					},
				},
			},
		},
	}
}

// buildMattermostFormat builds Mattermost-compatible payload
func (w *WebhookNotifier) buildMattermostFormat(data map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"text": fmt.Sprintf("## %s\n\n%s\n\n**Service:** %s\n**Time:** %s",
			data["title"], data["message"], data["service"], data["timestamp"]),
	}
}

// setAuthentication sets authentication headers based on configuration
func (w *WebhookNotifier) setAuthentication(headers map[string]string) {
	switch strings.ToLower(w.config.AuthType) {
	case "bearer":
		if w.config.AuthToken != "" {
			headers["Authorization"] = "Bearer " + w.config.AuthToken
		}
	case "basic":
		if w.config.AuthUsername != "" && w.config.AuthPassword != "" {
			auth := fmt.Sprintf("%s:%s", w.config.AuthUsername, w.config.AuthPassword)
			headers["Authorization"] = "Basic " + w.base64Encode(auth)
		}
	case "apikey":
		if w.config.AuthToken != "" {
			if w.config.AuthHeader != "" {
				headers[w.config.AuthHeader] = w.config.AuthToken
			} else {
				headers["X-API-Key"] = w.config.AuthToken
			}
		}
	}
}

// base64Encode encodes string to base64
func (w *WebhookNotifier) base64Encode(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

// sendWithRetry sends the webhook with retry logic
func (w *WebhookNotifier) sendWithRetry(url, method string, payload interface{}, contentType string, headers map[string]string) error {
	maxRetries := 0
	if w.config.Retries > 0 {
		maxRetries = w.config.Retries
	}

	timeout := 30
	if w.config.Timeout > 0 {
		timeout = w.config.Timeout
	}

	// Handle different payload types
	var body io.Reader
	if payload != nil {
		switch v := payload.(type) {
		case string:
			body = strings.NewReader(v)
		default:
			jsonData, err := json.Marshal(payload)
			if err != nil {
				return fmt.Errorf("failed to marshal payload: %w", err)
			}
			body = bytes.NewReader(jsonData)
		}
	}

	return sendHTTPRequestWithCustomBody(url, method, body, contentType, headers, maxRetries, timeout, w.config.SkipTLSVerify)
}

// WebhookError represents a webhook-specific error
type WebhookError struct {
	StatusCode int
	Body       string
	Retryable  bool
	Message    string
}
