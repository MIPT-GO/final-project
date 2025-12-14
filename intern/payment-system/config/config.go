package config

import (
	"os"
	"strconv"
)

type Config struct {
	Host string
	Port int

	RedisHost     string
	RedisPort     int
	RedisDb       int
	RedisUser     string
	RedisPassword string

	LinkPrefix     string
	RouterPrefix   string
	PaymentTimeout int
	StorageTimeout int
	TemplateDir    string
}

func LoadInt(key string, def int) int {
	str_value := os.Getenv(key)
	if str_value == "" {
		return def
	}

	value, err := strconv.Atoi(str_value)

	if err != nil {
		return def
	}

	return value
}

func LoadString(key string, def string) string {
	value := os.Getenv(key)
	if value == "" {
		return def
	}

	return value
}

func (config *Config) LoadFromEnv() {
	config.Host = LoadString("HOST", "0.0.0.0")
	config.Port = LoadInt("PORT", 8000)
	config.RedisHost = LoadString("REDIS_HOST", "localhost")
	config.RedisPort = LoadInt("REDIS_PORT", 6379)
	config.RedisDb = LoadInt("REDIS_DB", 0)
	config.RedisUser = LoadString("REDIS_USER", "")
	config.RedisPassword = LoadString("REDIS_PASSWORD", "")
	config.PaymentTimeout = LoadInt("PAYMENT_TIMEOUT", 10)
	config.StorageTimeout = LoadInt("STORAGE_TIMEOUT", 100)
	config.LinkPrefix = LoadString("LINK_PREFIX", "")
	config.RouterPrefix = LoadString("ROUTER_PREFIX", "/v1/payment")
	config.TemplateDir = LoadString("TEMPLATE_DIR", "intern/payment-system/infrastructure/templates")
}
