package config

import (
	"os"
)

type Config struct {
	DynamoDBTable           string
	SessionTable            string
	RecaptchaSecretKey      string
	SESRegion               string
	NotificationDstEmail    string
	NotificationSrcEmail    string
	NotificationPhoneNumber string
	Environment             string
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func Load() *Config {
	return &Config{
		DynamoDBTable:        getEnv("COUNTERS_TABLE", ""),
		SessionTable:         getEnv("SESSION_TABLE", ""),
		RecaptchaSecretKey:   getEnv("RECAPTCHA_SECRET_KEY", ""),
		SESRegion:            getEnv("SES_REGION", ""),
		NotificationDstEmail: getEnv("NOTIFICATION_DST_EMAIL", "tinvuong2003@gmail.com"),
		NotificationSrcEmail: getEnv("NOTIFICATION_SRC_EMAIL", ""),

		NotificationPhoneNumber: getEnv("NOTIFICATION_DST_PHONE", "5139148401"),
		Environment:             getEnv("ENVIRONMENT", "dev"),
	}
}
