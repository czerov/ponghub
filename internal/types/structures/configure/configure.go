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
		URL           string            `yaml:"url"`
		Method        string            `yaml:"method,omitempty"`
		Headers       map[string]string `yaml:"headers,omitempty"`
		Body          string            `yaml:"body,omitempty"`
		StatusCode    int               `yaml:"status_code,omitempty"`
		ResponseRegex string            `yaml:"response_regex,omitempty"`
		// Original template values for display
		OriginalURL           string            `yaml:"-"`
		OriginalBody          string            `yaml:"-"`
		OriginalResponseRegex string            `yaml:"-"`
		OriginalHeaders       map[string]string `yaml:"-"`
	}

	// Configure defines the overall configuration structure for the application
	Configure struct {
		Services       []Service `yaml:"services"`
		Timeout        int       `yaml:"timeout,omitempty"`
		MaxRetryTimes  int       `yaml:"max_retry_times,omitempty"`
		MaxLogDays     int       `yaml:"max_log_days,omitempty"`
		CertNotifyDays int       `yaml:"cert_notify_days,omitempty"`
	}
)
