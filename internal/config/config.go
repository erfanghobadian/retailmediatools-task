package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config represents the application configuration
type Config struct {
	App      AppConfig      `split_words:"true"`
	Server   ServerConfig   `split_words:"true"`
	Database DatabaseConfig `split_words:"true"`
}

// AppConfig contains application-specific configuration
type AppConfig struct {
	Name        string `default:"Ad Bidding Service"`
	Environment string `default:"development"`
	LogLevel    string `default:"info" split_words:"true"`
	Version     string `default:"1.0.0"`
}

// ServerConfig contains HTTP server configuration
type ServerConfig struct {
	Port    int           `default:"8080"`
	Timeout time.Duration `default:"30s"`
}

type DatabaseConfig struct {
	Host     string `split_words:"true"`
	Port     int    `split_words:"true"`
	User     string `split_words:"true"`
	Password string `split_words:"true"`
	Database string `split_words:"true"`
}

// Load loads the configuration from environment variables
func Load() (*Config, error) {
	var config Config
	if err := envconfig.Process("app", &config); err != nil {
		return nil, err
	}

	return &config, nil
}
