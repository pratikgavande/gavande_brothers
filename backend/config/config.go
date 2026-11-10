package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Port      string
	MySQLDSN  string
	JWTSecret string
	SMTP      SMTPConfig
}

type SMTPConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

func (s SMTPConfig) Configured() bool {
	return strings.TrimSpace(s.Host) != "" &&
		strings.TrimSpace(s.Port) != "" &&
		strings.TrimSpace(s.Username) != "" &&
		s.Password != "" &&
		strings.TrimSpace(s.From) != ""
}

func Load() (Config, error) {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return Config{}, fmt.Errorf("load .env file: %w", err)
	}

	port := strings.TrimSpace(os.Getenv("PORT"))
	if port == "" {
		port = "8080"
	}

	dsn := strings.TrimSpace(os.Getenv("MYSQL_DSN"))
	if dsn == "" {
		return Config{}, errors.New("MYSQL_DSN must be set")
	}

	jwtSecret := strings.TrimSpace(os.Getenv("JWT_SECRET"))
	if len(jwtSecret) < 32 {
		return Config{}, errors.New("JWT_SECRET must be at least 32 characters")
	}

	return Config{
		Port:      port,
		MySQLDSN:  dsn,
		JWTSecret: jwtSecret,
		SMTP: SMTPConfig{
			Host:     strings.TrimSpace(os.Getenv("SMTP_HOST")),
			Port:     strings.TrimSpace(os.Getenv("SMTP_PORT")),
			Username: strings.TrimSpace(os.Getenv("SMTP_USERNAME")),
			Password: os.Getenv("SMTP_PASSWORD"),
			From:     strings.TrimSpace(os.Getenv("SMTP_FROM")),
		},
	}, nil
}
