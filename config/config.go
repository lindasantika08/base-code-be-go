package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config adalah struct utama yang menampung semua konfigurasi
type Config struct {
	App      AppConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Redis    RedisConfig 
}

type AppConfig struct {
	Name           string
	Version        string
	Env            string 
	Port           string
	LogLevel       string 
	AllowedOrigins []string
}

type DatabaseConfig struct {
	Host            string
	Port            string
	Name            string
	User            string
	Password        string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime int 
}

type JWTConfig struct {
	Secret          string
	ExpiryHours     int
	RefreshSecret   string
	RefreshExpHours int
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

func Load() (*Config, error) {
	_ = godotenv.Load(".env")

	cfg := &Config{
		App: AppConfig{
			Name:    getEnv("APP_NAME", "go-base-project"),
			Version: getEnv("APP_VERSION", "1.0.0"),
			Env:     getEnv("APP_ENV", "development"),
			Port:    getEnv("APP_PORT", "8080"),
			LogLevel: getEnv("LOG_LEVEL", "debug"),
			AllowedOrigins: []string{
				getEnv("ALLOWED_ORIGINS", "*"),
			},
		},

		Database: DatabaseConfig{
			Host:            getEnv("DB_HOST", "localhost"),
			Port:            getEnv("DB_PORT", "3306"),
			Name:            getEnv("DB_NAME", "mydb"),
			User:            getEnv("DB_USER", "root"),
			Password:        getEnv("DB_PASSWORD", ""),
			SSLMode:         getEnv("DB_SSLMODE", "disable"),
			MaxOpenConns:    getEnvInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getEnvInt("DB_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: getEnvInt("DB_CONN_MAX_LIFETIME", 5),
		},

		JWT: JWTConfig{
			Secret:          getEnv("JWT_SECRET", "GANTI_INI_DENGAN_SECRET_YANG_KUAT"),
			ExpiryHours:     getEnvInt("JWT_EXPIRY_HOURS", 24),
			RefreshSecret:   getEnv("JWT_REFRESH_SECRET", "GANTI_INI_JUGA"),
			RefreshExpHours: getEnvInt("JWT_REFRESH_EXPIRY_HOURS", 168), // 7 hari
		},

		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvInt("REDIS_DB", 0),
		},
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) validate() error {
	if c.Database.Password == "" && c.App.Env == "production" {
		return fmt.Errorf("DB_PASSWORD wajib diisi di production")
	}
	if c.JWT.Secret == "GANTI_INI_DENGAN_SECRET_YANG_KUAT" && c.App.Env == "production" {
		return fmt.Errorf("JWT_SECRET wajib diganti di production")
	}
	return nil
}

func (d *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		d.Host, d.Port, d.Name, d.User, d.Password, d.SSLMode,
	)
}

// ── HELPER FUNCTIONS ─────────────────────────────────────────

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}
