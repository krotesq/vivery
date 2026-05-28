package config

import (
	"encoding/base64"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
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

	RedisHost            string
	RedisPort            string
	RedisPassword        string
	RedisProtocolVersion int

	AllowedOrigins []string

	JSONWebTokenSecret        []byte
	JSONWebTokenExpireMinutes int
	JSONWebTokenIssuer        string

	RefreshTokenExpireDays int

	BcryptCost int
}

func Load() (cfg *Config, err error) {

	// convert, decode and validate special env vars
	jwtExp, err := strconv.Atoi(getEnvOrDefault("JWT_EXP_MINUTES", "5"))
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_EXP_MINUTES: %s", err.Error())
	}

	refreshExp, err := strconv.Atoi(getEnvOrDefault("REFRESH_EXP_DAYS", "7"))
	if err != nil {
		return nil, fmt.Errorf("invalid REFRESH_EXP_DAYS: %s", err.Error())
	}

	bcryptCost, err := strconv.Atoi(getEnvOrDefault("BCRYPT_COST", "12"))
	if err != nil {
		return nil, fmt.Errorf("invalid BCRYPT_COST: %s", err.Error())
	}

	redisProtocolVersion, err := strconv.Atoi(getEnvOrDefault("REDIS_PROTOCOL_VERSION", "2"))
	if err != nil {
		return nil, fmt.Errorf("invalid REDIS_PROTOCOL_VERSION: %s", err.Error())
	}
	if redisProtocolVersion != 2 && redisProtocolVersion != 3 {
		return nil, fmt.Errorf("invalid REDIS_PROTOCOL_VERSION: version needs to be 2 or 3")
	}

	secret, err := base64.StdEncoding.DecodeString(getEnvOrDefault("JWT_SECRET", ""))
	if err != nil {
		return nil, fmt.Errorf("could not decode jwt secret: %s", err.Error())
	}
	if len(secret) == 0 {
		return nil, fmt.Errorf("JWT_SECRET can not be empty")
	}

	allowedOrigins := strings.Split(getEnvOrDefault("ALLOWED_ORIGINS", "http://localhost:3000"), ",")

	sslModes := []string{"disable", "allow", "prefer", "require", "verify-ca", "verify-full"}
	sslMode := getEnvOrDefault("DATABASE_SSLMODE", "prefer")
	if !slices.Contains(sslModes, sslMode) {
		return nil, fmt.Errorf("DATABASE_SSLMODE can only be %v", sslModes)
	}

	cfg = &Config{
		Host: getEnvOrDefault("HOST", "0.0.0.0"),
		Port: getEnvOrDefault("PORT", "3000"),

		DatabaseUser:     getEnvOrDefault("DATABASE_USER", "vivery"),
		DatabasePassword: getEnvOrDefault("DATABASE_PASSWORD", "vivery"),
		DatabaseHost:     getEnvOrDefault("DATABASE_HOST", "localhost"),
		DatabasePort:     getEnvOrDefault("DATABASE_PORT", "5432"),
		DatabaseName:     getEnvOrDefault("DATABASE_NAME", "vivery"),
		DatabaseSSLMode:  sslMode,

		RedisHost:            getEnvOrDefault("REDIS_HOST", "localhost"),
		RedisPort:            getEnvOrDefault("REDIS_PORT", "6379"),
		RedisPassword:        getEnvOrDefault("REDIS_PASSWORD", ""),
		RedisProtocolVersion: redisProtocolVersion,

		JSONWebTokenSecret:        secret,
		JSONWebTokenExpireMinutes: jwtExp,
		JSONWebTokenIssuer:        getEnvOrDefault("JWT_ISSUER", "vivery"),

		AllowedOrigins: allowedOrigins,

		RefreshTokenExpireDays: refreshExp,

		BcryptCost: bcryptCost,
	}

	return cfg, nil
}

func getEnvOrDefault(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}
	return fallback
}
