package configure

import (
	"log"
	"os"

	"github.com/wcy-dt/ponghub/internal/common"
	"github.com/wcy-dt/ponghub/internal/types/structures/configure"
	"github.com/wcy-dt/ponghub/internal/types/types/default_config"

	"gopkg.in/yaml.v3"
)

// setDefaultConfigs sets default values for the configuration fields
func setDefaultConfigs(cfg *configure.Configure) {
	default_config.SetDefaultTimeout(&cfg.Timeout)
	default_config.SetDefaultMaxRetryTimes(&cfg.MaxRetryTimes)
	default_config.SetDefaultMaxLogDays(&cfg.MaxLogDays)
	default_config.SetDefaultCertNotifyDays(&cfg.CertNotifyDays)
	default_config.SetDefaultDisplayNum(&cfg.DisplayNum)

	for i := range cfg.Services {
		default_config.SetDefaultTimeout(&cfg.Services[i].Timeout)
		default_config.SetDefaultMaxRetryTimes(&cfg.Services[i].MaxRetryTimes)
	}
}

// resolveConfigParameters resolves dynamic parameters in configuration
func resolveConfigParameters(cfg *configure.Configure) {
	resolver := common.NewParameterResolver()

	for i := range cfg.Services {
		for j := range cfg.Services[i].Endpoints {
			endpoint := &cfg.Services[i].Endpoints[j]

			// Save original template values
			endpoint.OriginalURL = endpoint.URL
			endpoint.OriginalBody = endpoint.Body
			endpoint.OriginalResponseRegex = endpoint.ResponseRegex
			if endpoint.Headers != nil {
				endpoint.OriginalHeaders = make(map[string]string)
				for key, value := range endpoint.Headers {
					endpoint.OriginalHeaders[key] = value
				}
			}

			// Resolve parameters
			endpoint.URL = resolver.ResolveParameters(endpoint.URL)
			endpoint.Body = resolver.ResolveParameters(endpoint.Body)
			endpoint.ResponseRegex = resolver.ResolveParameters(endpoint.ResponseRegex)

			// Resolve headers
			for key, value := range endpoint.Headers {
				endpoint.Headers[key] = resolver.ResolveParameters(value)
			}
		}
	}
}

// ReadConfigs loads the configuration from a YAML file at the specified path
func ReadConfigs(path string) (*configure.Configure, error) {
	// Read the configuration file
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func(f *os.File) {
		if err := f.Close(); err != nil {
			log.Println("Error closing config file:", err)
		}
	}(f)

	// Decode the YAML configuration
	cfg := new(configure.Configure)
	decoder := yaml.NewDecoder(f)
	if err := decoder.Decode(cfg); err != nil {
		log.Fatalln("Failed to decode YAML config:", err)
	}

	// Resolve dynamic parameters
	resolveConfigParameters(cfg)

	// Set default values for the configuration
	setDefaultConfigs(cfg)

	if len(cfg.Services) == 0 {
		log.Fatalln("No services defined in the configuration file")
	}
	return cfg, nil
}
