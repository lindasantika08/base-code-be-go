package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Log      LogConfig
}

type AppConfig struct {
	Name    string
	Env     string
	Port    string
	Version string
}

type DatabaseConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	Name            string
	MaxOpenConn     int
	MaxIdleConn     int
	ConnMaxLifetime time.Duration
}

type JWTConfig struct {
	Secret      string
	ExpiryHours int
}

type LogConfig struct {
	Level  string
	Format string
}

func Load() (*Config, error) {
	// Load .env file (ignore error if not found — use OS env directly)
	_ = godotenv.Load()

	cfg := &Config{
		App: AppConfig{
			Name:    getEnv("APP_NAME", "banking-app"),
			Env:     getEnv("APP_ENV", "development"),
			Port:    getEnv("APP_PORT", "8080"),
			Version: getEnv("APP_VERSION", "1.0.0"),
		},
		Database: DatabaseConfig{
			Host:            getEnv("DB_HOST", "localhost"),
			Port:            getEnv("DB_PORT", "3306"),
			User:            getEnv("DB_USER", "root"),
			Password:        getEnv("DB_PASSWORD", ""),
			Name:            getEnv("DB_NAME", "banking_db"),
			MaxOpenConn:     getEnvInt("DB_MAX_OPEN_CONN", 25),
			MaxIdleConn:     getEnvInt("DB_MAX_IDLE_CONN", 5),
			ConnMaxLifetime: getEnvDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute),
		},
		JWT: JWTConfig{
			Secret:      getEnv("JWT_SECRET", "change-me"),
			ExpiryHours: getEnvInt("JWT_EXPIRY_HOURS", 24),
		},
		Log: LogConfig{
			Level:  getEnv("LOG_LEVEL", "debug"),
			Format: getEnv("LOG_FORMAT", "text"),
		},
	}

	return cfg, cfg.validate()
}

func (c *Config) validate() error {
	if c.Database.Password == "" && c.App.Env == "production" {
		return fmt.Errorf("DB_PASSWORD is required in production")
	}
	if c.JWT.Secret == "change-me" && c.App.Env == "production" {
		return fmt.Errorf("JWT_SECRET must be changed in production")
	}
	return nil
}

func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local",
		c.User, c.Password, c.Host, c.Port, c.Name)
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return fallback
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return fallback
}
