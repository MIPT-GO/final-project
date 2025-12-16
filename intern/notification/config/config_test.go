package config_test

import (
	"testing"

	"final-project/intern/notification/config"

	"github.com/stretchr/testify/assert"
)

func TestLoadFromEnv_Defaults(t *testing.T) {
	t.Setenv("LOG_LEVEL", "")
	t.Setenv("HOST", "")
	t.Setenv("KAFKA_BROKERS", "")
	t.Setenv("KAFKA_TOPIC", "")
	t.Setenv("KAFKA_GROUP_ID", "")
	t.Setenv("SMTP_HOST", "")
	t.Setenv("SMTP_PORT", "")
	t.Setenv("SMTP_USERNAME", "")
	t.Setenv("SMTP_PASSWORD", "")
	t.Setenv("SMTP_FROM", "")

	cfg := config.Config{}
	cfg.LoadFromEnv()

	assert.Equal(t, "info", cfg.LogLevel)
	assert.Equal(t, "0.0.0.0", cfg.Host)
	assert.Equal(t, []string{"kafka:29092"}, cfg.KafkaBrokers)
	assert.Equal(t, "notifications", cfg.KafkaTopic)
	assert.Equal(t, "notification-svc", cfg.KafkaGroupID)
	assert.Equal(t, "smtp.gmail.com", cfg.SMTPHost)
	assert.Equal(t, 587, cfg.SMTPPort)
	assert.Equal(t, "", cfg.SMTPUsername)
	assert.Equal(t, "", cfg.SMTPPassword)
	assert.Equal(t, "", cfg.SMTPFrom)
}

func TestLoadFromEnv_KafkaHostList(t *testing.T) {
	t.Setenv("KAFKA_BROKERS", "kafka:29092, kafka2:29092,, ")

	cfg := config.Config{}
	cfg.LoadFromEnv()

	assert.Equal(t, []string{"kafka:29092", "kafka2:29092"}, cfg.KafkaBrokers)
}

func TestLoadFromEnv_SMTPFromDefaultsToUsername(t *testing.T) {
	t.Setenv("SMTP_USERNAME", "sender@example.com")
	t.Setenv("SMTP_FROM", "")

	cfg := config.Config{}
	cfg.LoadFromEnv()

	assert.Equal(t, "sender@example.com", cfg.SMTPFrom)
}

func TestLoadFromEnv_SMTPPortInvalidUsesDefault(t *testing.T) {
	t.Setenv("SMTP_PORT", "not-a-number")

	cfg := config.Config{}
	cfg.LoadFromEnv()

	assert.Equal(t, 587, cfg.SMTPPort)
}
