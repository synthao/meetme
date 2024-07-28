package config

import (
	"fmt"
	"net"
	"os"
)

type Config struct {
	ServerAddress     string
	AppVersion        string
	GrpcServerAddress string
	LogLevel          string
}

func LoadConfig() *Config {
	return &Config{
		AppVersion:        getEnv("APP_VERSION", ""),
		ServerAddress:     net.JoinHostPort(getEnv("SERVER_HOST", "0.0.0.0"), getEnv("SERVER_PORT", "8080")),
		GrpcServerAddress: os.Getenv("IMGPROC_GRPC_SERVER_ADDRESS"),
		LogLevel:          getEnv("LOG_LEVEL", "INFO"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	)
}
