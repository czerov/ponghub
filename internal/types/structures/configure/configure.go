package configure

type (
	// Service defines the configuration for a service, including its health and Endpoints ports
	Service struct {
		Name          string     `yaml:"name"`
		Endpoints     []Endpoint `yaml:"endpoints"`
		Timeout       int        `yaml:"timeout,omitempty"`
		MaxRetryTimes int        `yaml:"max_retry_times,omitempty"`
	}

	// Endpoint defines the configuration for a port
	Endpoint struct {
		URL                 string            `yaml:"url"`
		ParsedURL           string            `yaml:"-"`
		Method              string            `yaml:"method,omitempty"`
		Headers             map[string]string `yaml:"headers,omitempty"`
		ParsedHeaders       map[string]string `yaml:"-"`
		Body                string            `yaml:"body,omitempty"`
		ParsedBody          string            `yaml:"-"`
		StatusCode          int               `yaml:"status_code,omitempty"`
		ResponseRegex       string            `yaml:"response_regex,omitempty"`
		ParsedResponseRegex string            `yaml:"-"`
	}

	// NotificationConfig defines notification settings
	NotificationConfig struct {
		Enabled  bool            `yaml:"enabled"`
		Methods  []string        `yaml:"methods"`
		Email    *EmailConfig    `yaml:"email,omitempty"`
		Discord  *DiscordConfig  `yaml:"discord,omitempty"`
		Slack    *SlackConfig    `yaml:"slack,omitempty"`
		Telegram *TelegramConfig `yaml:"telegram,omitempty"`
		WeChat   *WeChatConfig   `yaml:"wechat,omitempty"`
		Webhook  *WebhookConfig  `yaml:"webhook,omitempty"`
		Default  *DefaultConfig  `yaml:"default,omitempty"`
	}

	// EmailConfig defines email notification settings
	EmailConfig struct {
		SMTPHost string   `yaml:"smtp_host"`
		SMTPPort int      `yaml:"smtp_port"`
		From     string   `yaml:"from"`
		To       []string `yaml:"to"`
		Subject  string   `yaml:"subject,omitempty"`
	}

	// DiscordConfig defines Discord webhook notification settings
	DiscordConfig struct {
		WebhookURL string `yaml:"webhook_url"`
		Username   string `yaml:"username,omitempty"`
		AvatarURL  string `yaml:"avatar_url,omitempty"`
	}

	// SlackConfig defines Slack webhook notification settings
	SlackConfig struct {
		WebhookURL string `yaml:"webhook_url"`
		Channel    string `yaml:"channel,omitempty"`
		Username   string `yaml:"username,omitempty"`
		IconEmoji  string `yaml:"icon_emoji,omitempty"`
	}

	// TelegramConfig defines Telegram bot notification settings
	TelegramConfig struct {
		BotToken string `yaml:"bot_token"`
		ChatID   string `yaml:"chat_id"`
	}

	// WeChatConfig defines WeChat Work webhook notification settings
	WeChatConfig struct {
		WebhookURL string `yaml:"webhook_url"`
	}

	// WebhookConfig defines generic webhook notification settings
	WebhookConfig struct {
		URL     string            `yaml:"url"`
		Method  string            `yaml:"method,omitempty"`
		Headers map[string]string `yaml:"headers,omitempty"`
	}

	// DefaultConfig defines default notification settings (GitHub Actions failure)
	DefaultConfig struct {
		Enabled bool `yaml:"enabled"`
	}

	// Configure defines the overall configuration structure for the application
	Configure struct {
		Services       []Service           `yaml:"services"`
		Timeout        int                 `yaml:"timeout,omitempty"`
		MaxRetryTimes  int                 `yaml:"max_retry_times,omitempty"`
		MaxLogDays     int                 `yaml:"max_log_days,omitempty"`
		CertNotifyDays int                 `yaml:"cert_notify_days,omitempty"`
		DisplayNum     int                 `yaml:"display_num,omitempty"`
		Notifications  *NotificationConfig `yaml:"notifications,omitempty"`
	}
)
