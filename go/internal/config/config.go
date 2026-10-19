package config

import (
	"os"
	"strconv"
)

type Config struct {
	GatewayPort  string
	GuidePort    string
	MatchingPort string
	ChatPort     string
	AIPort       string
	BookingPort  string

	DatabaseURL string
	RedisURL    string
	GeminiAPIKey string

	GuideServiceURL    string
	MatchingServiceURL string
	ChatServiceURL     string
	AIServiceURL       string
	BookingServiceURL  string
}

func Load() *Config {
	return &Config{
		GatewayPort:  getEnv("GATEWAY_PORT", "8080"),
		GuidePort:    getEnv("GUIDE_PORT", "8081"),
		MatchingPort: getEnv("MATCHING_PORT", "8082"),
		ChatPort:     getEnv("CHAT_PORT", "8083"),
		AIPort:       getEnv("AI_PORT", "8084"),
		BookingPort:  getEnv("BOOKING_PORT", "8085"),

		DatabaseURL:  getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/guidenow?sslmode=disable"),
		RedisURL:     getEnv("REDIS_URL", "localhost:6379"),
		GeminiAPIKey: getEnv("GEMINI_API_KEY", ""),

		GuideServiceURL:    getEnv("GUIDE_SERVICE_URL", "http://localhost:8081"),
		MatchingServiceURL: getEnv("MATCHING_SERVICE_URL", "http://localhost:8082"),
		ChatServiceURL:     getEnv("CHAT_SERVICE_URL", "http://localhost:8083"),
		AIServiceURL:       getEnv("AI_SERVICE_URL", "http://localhost:8084"),
		BookingServiceURL:  getEnv("BOOKING_SERVICE_URL", "http://localhost:8085"),
	}
}

func getEnv(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return defaultVal
}

func GetEnvInt(key string, defaultVal int) int {
	if val, ok := os.LookupEnv(key); ok {
		if intVal, err := strconv.Atoi(val); err == nil {
			return intVal
		}
	}
	return defaultVal
}
