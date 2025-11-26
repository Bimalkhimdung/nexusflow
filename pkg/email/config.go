package email

import (
	"fmt"
	"os"
	"strconv"
)

// LoadConfigFromEnv loads email configuration from environment variables
func LoadConfigFromEnv() (*Config, error) {
	smtpPort, err := strconv.Atoi(getEnv("SMTP_PORT", "587"))
	if err != nil {
		return nil, fmt.Errorf("invalid SMTP_PORT: %w", err)
	}

	config := &Config{
		SMTPHost:     getEnv("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:     smtpPort,
		SMTPUser:     getEnv("SMTP_USER", ""),
		SMTPPassword: getEnv("SMTP_PASSWORD", ""),
		FromAddress:  getEnv("SMTP_FROM_ADDRESS", "noreply@nexusflow.io"),
		FromName:     getEnv("SMTP_FROM_NAME", "NexusFlow"),
		TemplateDir:  getEnv("EMAIL_TEMPLATE_DIR", "templates/email"),
	}

	// Validate required fields
	if config.SMTPUser == "" {
		return nil, fmt.Errorf("SMTP_USER is required")
	}
	if config.SMTPPassword == "" {
		return nil, fmt.Errorf("SMTP_PASSWORD is required")
	}

	return config, nil
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
