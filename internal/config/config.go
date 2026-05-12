package config

import (
	"fmt"
	"os"
	"slices"
	"strconv"
)

type Config struct {
	Host string
	Port string

	DatabaseUser     string
	DatabasePassword string
	DatabaseHost     string
	DatabasePort     string
	DatabaseName     string
	DatabaseSSLMode  string

	JSONWebTokenSecret        string
	JSONWebTokenExpireMinutes int
	JSONWebTokenIssuer        string

	RefreshTokenExpireDays int

	BcryptCost int
}

func Load() (*Config, error) {
	jwtExp, err := strconv.Atoi(getEnvOrDefault("JWT_EXP_MINUTES", "5"))
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_EXP_MINUTES: %w", err)
	}

	refreshExp, err := strconv.Atoi(getEnvOrDefault("REFRESH_EXP_DAYS", "7"))
	if err != nil {
		return nil, fmt.Errorf("invalid REFRESH_EXP_DAYS: %w", err)
	}

	bcryptCost, err := strconv.Atoi(getEnvOrDefault("BCRYPT_COST", "12"))
	if err != nil {
		return nil, fmt.Errorf("invalid BCRYPT_COST: %w", err)
	}

	cfg := &Config{
		Host: getEnvOrDefault("HOST", "0.0.0.0"),
		Port: getEnvOrDefault("PORT", "3000"),

		DatabaseUser:     getEnvOrDefault("DATABASE_USER", ""),
		DatabasePassword: getEnvOrDefault("DATABASE_PASSWORD", ""),
		DatabaseHost:     getEnvOrDefault("DATABASE_HOST", "localhost"),
		DatabasePort:     getEnvOrDefault("DATABASE_PORT", "5432"),
		DatabaseName:     getEnvOrDefault("DATABASE_NAME", "vivery"),
		DatabaseSSLMode:  getEnvOrDefault("DATABASE_SSLMODE", "prefer"),

		JSONWebTokenSecret:        getEnvOrDefault("JWT_SECRET", ""),
		JSONWebTokenExpireMinutes: jwtExp,
		JSONWebTokenIssuer:        getEnvOrDefault("JWT_ISSUER", "vivery"),

		RefreshTokenExpireDays: refreshExp,

		BcryptCost: bcryptCost,
	}

	if cfg.DatabaseUser == "" {
		return nil, fmt.Errorf("DATABASE_USER is required")
	}
	if cfg.DatabasePassword == "" {
		return nil, fmt.Errorf("DATABASE_PASSWORD is required")
	}
	if cfg.JSONWebTokenSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}

	sslModes := []string{"disable", "allow", "prefer", "require", "verify-ca", "verify-full"}
	if !slices.Contains(sslModes, cfg.DatabaseSSLMode) {
		return nil, fmt.Errorf("DATABASE_SSLMODE can only be %v", sslModes)
	}

	return cfg, nil
}

func getEnvOrDefault(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}
	return fallback
}
