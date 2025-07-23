package config

import "os"

// Config holds all configurable environment variables
type Config struct {
	ListingServiceURL string
	UserServiceURL    string
}

// Load reads env vars and returns a Config struct
func Load() Config {
	return Config{
		ListingServiceURL: getEnv("LISTING_SERVICE_URL", "http://localhost:6000"),
		UserServiceURL:    getEnv("USER_SERVICE_URL", "http://localhost:6001"),
	}
}

func getEnv(key string, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
