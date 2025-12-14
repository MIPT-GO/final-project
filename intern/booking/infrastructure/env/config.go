package env

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"time"

	"final-project/pkg/booking/constants"
)

type DBConfig struct {
	DSN string
}

type HotelConfig struct {
	BaseURL string
	Timeout time.Duration
}

type ServerConfig struct {
	Port string
}

type Config struct {
	DB       DBConfig
	Hotel    HotelConfig
	Server   ServerConfig
	LogLevel string
}

func parseEnvFile(path string) (map[string]string, error) {
	f, err := os.Open(filepath.Clean(path))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	m := make(map[string]string)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		m[key] = val
	}
	return m, scanner.Err()
}

func lookup(env map[string]string, key, def string) string {
	if env != nil {
		if v, ok := env[key]; ok && v != "" {
			return v
		}
	}
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func NewConfigFromFile(path string) (*Config, error) {
	envMap := make(map[string]string)
	if path != "" {
		m, err := parseEnvFile(path)
		if err == nil {
			envMap = m
		}
	}
	host := lookup(envMap, constants.EnvPostgresHost, "localhost")
	port := lookup(envMap, constants.EnvPostgresPort, constants.DefaultPostgresPort)
	user := lookup(envMap, constants.EnvPostgresUser, "postgres")
	pass := lookup(envMap, constants.EnvPostgresPassword, "")
	dbname := lookup(envMap, constants.EnvPostgresDB, "booking")

	dsn := "postgres://" + user + ":" + pass + "@" + host + ":" + port + "/" + dbname + "?sslmode=disable"

	hotelHost := lookup(envMap, constants.EnvHotelServiceHost, "http://localhost")
	hotelPort := lookup(envMap, constants.EnvHotelServicePort, constants.DefaultHotelServicePort)
	hotelTimeoutStr := lookup(envMap, constants.EnvHotelServiceTimeout, "5s")
	hotelTimeout, _ := time.ParseDuration(hotelTimeoutStr)
	hotelBaseURL := strings.TrimRight(hotelHost, "/") + ":" + hotelPort

	srvPort := lookup(envMap, constants.EnvServerPort, constants.DefaultServerPort)

	logLevel := lookup(envMap, constants.EnvLogLevel, "info")

	cfg := &Config{
		DB: DBConfig{
			DSN: dsn,
		},
		Hotel: HotelConfig{
			BaseURL: hotelBaseURL,
			Timeout: hotelTimeout,
		},
		Server: ServerConfig{
			Port: srvPort,
		},
		LogLevel: logLevel,
	}

	return cfg, nil
}
