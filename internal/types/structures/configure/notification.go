package configure

type (
	// NotificationConfig defines the configuration for all notification channels
	NotificationConfig struct {
		Enabled  bool            `yaml:"enabled,omitempty"`
		Methods  []string        `yaml:"methods,omitempty"`
		Default  *DefaultConfig  `yaml:"default,omitempty"`
		Discord  *DiscordConfig  `yaml:"discord,omitempty"`
		Email    *EmailConfig    `yaml:"email,omitempty"`
		Slack    *SlackConfig    `yaml:"slack,omitempty"`
		Telegram *TelegramConfig `yaml:"telegram,omitempty"`
		WeChat   *WeChatConfig   `yaml:"wechat,omitempty"`
		Webhook  *WebhookConfig  `yaml:"webhook,omitempty"`
	}

	// DiscordConfig defines Discord webhook notification settings
	DiscordConfig struct {
		WebhookURL string   `yaml:"webhook_url,omitempty"`
		Username   string   `yaml:"username,omitempty"`
		AvatarURL  string   `yaml:"avatar_url,omitempty"`
		UseEmbeds  bool     `yaml:"use_embeds,omitempty"`
		Color      int      `yaml:"color,omitempty"`
		Mentions   []string `yaml:"mentions,omitempty"`
		Retries    int      `yaml:"retries,omitempty"`
		Timeout    int      `yaml:"timeout,omitempty"`
		UserAgent  string   `yaml:"user_agent,omitempty"`
	}

	// EmailConfig defines SMTP email notification settings
	EmailConfig struct {
		SMTPHost    string   `yaml:"smtp_host"`
		SMTPPort    int      `yaml:"smtp_port"`
		From        string   `yaml:"from"`
		To          []string `yaml:"to"`
		ReplyTo     string   `yaml:"reply_to,omitempty"`
		UseTLS      bool     `yaml:"use_tls,omitempty"`
		UseStartTLS bool     `yaml:"use_starttls,omitempty"`
		SkipVerify  bool     `yaml:"skip_verify,omitempty"`
	}

	// SlackConfig defines Slack webhook notification settings
	SlackConfig struct {
		WebhookURL string   `yaml:"webhook_url,omitempty"`
		Channel    string   `yaml:"channel,omitempty"`
		Username   string   `yaml:"username,omitempty"`
		IconEmoji  string   `yaml:"icon_emoji,omitempty"`
		IconURL    string   `yaml:"icon_url,omitempty"`
		UseBlocks  bool     `yaml:"use_blocks,omitempty"`
		Color      string   `yaml:"color,omitempty"`
		Mentions   []string `yaml:"mentions,omitempty"`
		Retries    int      `yaml:"retries,omitempty"`
		Timeout    int      `yaml:"timeout,omitempty"`
		UserAgent  string   `yaml:"user_agent,omitempty"`
	}

	// TelegramConfig defines Telegram bot notification settings
	TelegramConfig struct {
		BotToken              string `yaml:"bot_token,omitempty"`
		ChatID                string `yaml:"chat_id,omitempty"`
		ParseMode             string `yaml:"parse_mode,omitempty"`
		DisableWebPagePreview bool   `yaml:"disable_web_page_preview,omitempty"`
		DisableNotification   bool   `yaml:"disable_notification,omitempty"`
		ReplyToMessageID      int    `yaml:"reply_to_message_id,omitempty"`
		Retries               int    `yaml:"retries,omitempty"`
		Timeout               int    `yaml:"timeout,omitempty"`
		UserAgent             string `yaml:"user_agent,omitempty"`
	}

	// WeChatConfig defines WeChat Work webhook notification settings
	WeChatConfig struct {
		WebhookURL string   `yaml:"webhook_url,omitempty"`
		MsgType    string   `yaml:"msg_type,omitempty"`
		Mentions   []string `yaml:"mentions,omitempty"`
		Retries    int      `yaml:"retries,omitempty"`
		Timeout    int      `yaml:"timeout,omitempty"`
		UserAgent  string   `yaml:"user_agent,omitempty"`
	}

	// WebhookConfig defines generic webhook notification settings
	WebhookConfig struct {
		URL           string            `yaml:"url,omitempty"`
		Method        string            `yaml:"method,omitempty"`
		Headers       map[string]string `yaml:"headers,omitempty"`
		Template      string            `yaml:"template,omitempty"`
		Format        string            `yaml:"format,omitempty"`
		ContentType   string            `yaml:"content_type,omitempty"`
		AuthType      string            `yaml:"auth_type,omitempty"`
		AuthToken     string            `yaml:"auth_token,omitempty"`
		AuthUsername  string            `yaml:"auth_username,omitempty"`
		AuthPassword  string            `yaml:"auth_password,omitempty"`
		AuthHeader    string            `yaml:"auth_header,omitempty"`
		Retries       int               `yaml:"retries,omitempty"`
		Timeout       int               `yaml:"timeout,omitempty"`
		SkipTLSVerify bool              `yaml:"skip_tls_verify,omitempty"`
	}
)
