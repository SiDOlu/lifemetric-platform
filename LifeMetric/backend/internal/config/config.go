package config

import (
	"os"
)

// Config holds our environment variable properties
type Config struct {
	Port              string
	DatabaseURL       string
	JWTSecret         string
	MyBalanceSecret   string
}

// LoadConfig initializes application configurations with secure default values
func LoadConfig() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgresql://postgres:postgres@localhost:5432/lifemetric?sslmode=disable"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "sec_jwt_lifemetric_ambient_intel_39485"
	}

	StripeWebhookSecret := os.Getenv("MY_BALANCE_SECRET")
	if StripeWebhookSecret == "" {
		StripeWebhookSecret = "whsec_mock_lifemetric_test_key"
	}

	return &Config{
		Port:            port,
		DatabaseURL:     dbURL,
		JWTSecret:       jwtSecret,
		MyBalanceSecret: StripeWebhookSecret,
	}
}
