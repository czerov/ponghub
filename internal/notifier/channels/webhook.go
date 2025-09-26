package channels

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
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
			TLSClientConfig: &tls.Config{},
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
