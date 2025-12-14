package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Host string
	Port int

	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string

	ReadTimeout  int
	WriteTimeout int
	IdleTimeout  int
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
	config.Port = LoadInt("PORT", 8080)

	config.ReadTimeout = LoadInt("READ_TIMEOUT_SEC", 5)
	config.WriteTimeout = LoadInt("WRITE_TIMEOUT_SEC", 10)
	config.IdleTimeout = LoadInt("IDLE_TIMEOUT_SEC", 15)

	config.DBHost = LoadString("DB_HOST", "localhost")
	config.DBPort = LoadInt("DB_PORT", 5432)
	config.DBUser = LoadString("DB_USER", "postgres")
	config.DBPassword = LoadString("DB_PASSWORD", "secret")
	config.DBName = LoadString("DB_NAME", "hotel_db")
}

func (c *Config) GetDBConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName)
}
