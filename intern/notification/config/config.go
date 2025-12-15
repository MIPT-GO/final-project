package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	LogLevel string

	Host string

	KafkaBrokers []string
	KafkaTopic   string
	KafkaGroupID string

	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
	SMTPFrom     string
}

func getEnvString(key, defaultValue string) string {
	value := os.Getenv(key)
	if value != "" {
		return value
	}

	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	rawValue := os.Getenv(key)
	if rawValue == "" {
		return defaultValue
	}

	intValue, err := strconv.Atoi(rawValue)
	if err != nil {
		return defaultValue
	}

	return intValue
}

func (config *Config) LoadFromEnv() {
	config.LogLevel = getEnvString("LOG_LEVEL", "info")
	config.Host = getEnvString("HOST", "0.0.0.0")

	brokersList := getEnvString("KAFKA_BROKERS", "kafka:29092")
	config.KafkaBrokers = make([]string, 0)
	for _, brokerRaw := range strings.Split(brokersList, ",") {
		broker := strings.TrimSpace(brokerRaw)
		if broker != "" {
			config.KafkaBrokers = append(config.KafkaBrokers, broker)
		}
	}

	config.KafkaTopic = getEnvString("KAFKA_TOPIC", "notifications")
	config.KafkaGroupID = getEnvString("KAFKA_GROUP_ID", "notification-svc")

	config.SMTPHost = getEnvString("SMTP_HOST", "smtp.gmail.com")
	config.SMTPPort = getEnvInt("SMTP_PORT", 587)
	config.SMTPUsername = getEnvString("SMTP_USERNAME", "")
	config.SMTPPassword = getEnvString("SMTP_PASSWORD", "")
	config.SMTPFrom = getEnvString("SMTP_FROM", config.SMTPUsername)
}
