package config

import "os"

type MailerConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

func LoadMailerConfig() *MailerConfig {
	return &MailerConfig{
		Host:     getEnvOrDefault("SMTP_HOST", "smtp.gmail.com"),
		Port:     587, // Default TLS port
		Username: os.Getenv("SMTP_USERNAME"),
		Password: os.Getenv("SMTP_PASSWORD"),
		From:     getEnvOrDefault("SMTP_FROM", "noreply@yourdomain.com"),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
