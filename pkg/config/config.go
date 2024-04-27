package config

import (
	"log"
	"os"
)

type Config struct {
	Port          string
	DatabaseUrl   string
	AdminUsername string
	AdminPassword string
}

func NewConfig() *Config {
	return &Config{
		Port:          getenv("PORT"),
		DatabaseUrl:   getenv("DATABASE_URL"),
		AdminUsername: getenv("ADMIN_USERNAME"),
		AdminPassword: getenv("ADMIN_PASSWORD"),
	}
}

func getenv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("missing required environment variable: %v", key)
	}
	return v
}
