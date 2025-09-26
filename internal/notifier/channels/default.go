package channels

import (
	"fmt"
	"log"
	"os"

	"github.com/wcy-dt/ponghub/internal/types/structures/configure"
)

// DefaultNotifier implements the NotificationService interface for default GitHub Actions notifications
type DefaultNotifier struct {
	config *configure.DefaultConfig
}

// NewDefaultNotifier creates a new default notifier
func NewDefaultNotifier(config *configure.DefaultConfig) *DefaultNotifier {
	return &DefaultNotifier{
		config: config,
	}
}

// Send implements the NotificationService interface
// For default notifications, we write to stderr and set an exit flag
func (d *DefaultNotifier) Send(title, message string) error {
	if d.config == nil {
		return fmt.Errorf("default notifier config is nil")
	}

	log.Println("🚨 DEFAULT NOTIFICATION TRIGGERED 🚨")
	log.Printf("Title: %s", title)
	log.Printf("Message:\n%s", message)

	// Write to stderr for GitHub Actions to capture
	_, _ = fmt.Fprintf(os.Stderr, "\n=== PongHub Alert ===\n")
	_, _ = fmt.Fprintf(os.Stderr, "%s\n\n", title)
	_, _ = fmt.Fprintf(os.Stderr, "%s\n", message)
	_, _ = fmt.Fprintf(os.Stderr, "=====================\n\n")

	// Set environment variable to indicate failure should occur
	_ = os.Setenv("PONGHUB_HAS_ALERTS", "true")

	log.Println("Default notification sent - GitHub Actions will be notified of service issues")
	return nil
}
