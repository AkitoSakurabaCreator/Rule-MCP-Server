package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port        int
	Host        string
	Environment string
	LogLevel    string
}

func LoadConfig() *Config {
	config := &Config{
		Port:        8080,
		Host:        "0.0.0.0",
		Environment: "development",
		LogLevel:    "info",
	}

	// Load configuration from environment variables
	if portStr := os.Getenv("PORT"); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil {
			config.Port = port
		}
	}

	if host := os.Getenv("HOST"); host != "" {
		config.Host = host
	}

	if env := os.Getenv("ENVIRONMENT"); env != "" {
		config.Environment = env
	}

	if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
		config.LogLevel = logLevel
	}

	return config
}

func (c *Config) GetAddress() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}
