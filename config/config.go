package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort        int
	Environment       string
	LogLevel          string
	AllowedOrigins    []string
	PrometheusEnabled bool
	PrometheusPath    string
	ZendeskSubdomain  string
	ZendeskToken      string
	ZendeskEmail      string
	WebhookEndpoint   string
	WebhookHost       string
	WebhookPort       int
}

// Load reads the configuration from .env file and environment variables
func Load() *Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found, using environment variables")
	}

	serverPort := getEnvAsInt("SERVER_PORT", 8080)
	webhookEndpoint := getEnv("WEBHOOK_ENDPOINT", "/webhook")
	webhookPort := getEnvAsInt("WEBHOOK_PORT", serverPort)

	config := &Config{
		ServerPort:        serverPort,
		Environment:       getEnv("ENVIRONMENT", "development"),
		LogLevel:          getEnv("LOG_LEVEL", "info"),
		AllowedOrigins:    getEnvAsSlice("ALLOWED_ORIGINS", []string{"http://localhost:3000"}),
		PrometheusEnabled: getEnvAsBool("PROMETHEUS_ENABLED", true),
		PrometheusPath:    getEnv("PROMETHEUS_PATH", "/metrics"),
		ZendeskSubdomain:  getEnv("ZENDESK_SUBDOMAIN", ""),
		ZendeskToken:      getEnv("ZENDESK_TOKEN", ""),
		ZendeskEmail:      getEnv("ZENDESK_EMAIL", ""),
		WebhookEndpoint:   webhookEndpoint,
		WebhookHost:       getEnv("WEBHOOK_HOST", "0.0.0.0"),
		WebhookPort:       webhookPort,
	}

	return config
}

// getEnv reads an environment variable with a default value
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

// getEnvAsInt reads an environment variable as integer with a default value
func getEnvAsInt(key string, defaultVal int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultVal
}

// getEnvAsSlice reads an environment variable as a comma-separated list
func getEnvAsSlice(key string, defaultVal []string) []string {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultVal
	}
	return strings.Split(valueStr, ",")
}

// getEnvAsBool reads an environment variable as a boolean with a default value
func getEnvAsBool(key string, defaultVal bool) bool {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultVal
}
