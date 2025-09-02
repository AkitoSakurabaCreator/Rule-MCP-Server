package main

import (
	"os"
	"testing"
)

func TestLoadConfigDefault(t *testing.T) {
	// 環境変数をクリア
	os.Unsetenv("PORT")
	os.Unsetenv("HOST")
	os.Unsetenv("ENVIRONMENT")
	os.Unsetenv("LOG_LEVEL")

	config := LoadConfig()

	if config.Port != 8080 {
		t.Errorf("Expected default port 8080, got %d", config.Port)
	}

	if config.Host != "0.0.0.0" {
		t.Errorf("Expected default host '0.0.0.0', got '%s'", config.Host)
	}

	if config.Environment != "development" {
		t.Errorf("Expected default environment 'development', got '%s'", config.Environment)
	}

	if config.LogLevel != "info" {
		t.Errorf("Expected default log level 'info', got '%s'", config.LogLevel)
	}
}

func TestLoadConfigFromEnv(t *testing.T) {
	// 環境変数を設定
	os.Setenv("PORT", "3000")
	os.Setenv("HOST", "127.0.0.1")
	os.Setenv("ENVIRONMENT", "production")
	os.Setenv("LOG_LEVEL", "warn")

	defer func() {
		os.Unsetenv("PORT")
		os.Unsetenv("HOST")
		os.Unsetenv("ENVIRONMENT")
		os.Unsetenv("LOG_LEVEL")
	}()

	config := LoadConfig()

	if config.Port != 3000 {
		t.Errorf("Expected port 3000, got %d", config.Port)
	}

	if config.Host != "127.0.0.1" {
		t.Errorf("Expected host '127.0.0.1', got '%s'", config.Host)
	}

	if config.Environment != "production" {
		t.Errorf("Expected environment 'production', got '%s'", config.Environment)
	}

	if config.LogLevel != "warn" {
		t.Errorf("Expected log level 'warn', got '%s'", config.LogLevel)
	}
}

func TestLoadConfigInvalidPort(t *testing.T) {
	// 無効なポート番号を設定
	os.Setenv("PORT", "invalid")
	defer os.Unsetenv("PORT")

	config := LoadConfig()

	// 無効なポートの場合はデフォルト値を使用
	if config.Port != 8080 {
		t.Errorf("Expected default port 8080 for invalid port, got %d", config.Port)
	}
}

func TestConfigGetAddress(t *testing.T) {
	config := &Config{
		Port: 8080,
		Host: "0.0.0.0",
	}

	expected := "0.0.0.0:8080"
	if address := config.GetAddress(); address != expected {
		t.Errorf("Expected address '%s', got '%s'", expected, address)
	}
}

func TestConfigEnvironment(t *testing.T) {
	devConfig := &Config{Environment: "development"}
	prodConfig := &Config{Environment: "production"}

	if !devConfig.IsDevelopment() {
		t.Error("Expected development environment to be true")
	}

	if !prodConfig.IsProduction() {
		t.Error("Expected production environment to be true")
	}

	if devConfig.IsProduction() {
		t.Error("Expected development environment to not be production")
	}

	if prodConfig.IsDevelopment() {
		t.Error("Expected production environment to not be development")
	}
}
