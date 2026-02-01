package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DBUri       string
	RedisUri    string
	GitHubToken string
}

func LoadConfig() (*Config, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	return &Config{
		Port:        os.Getenv("PORT"),
		DBUri:       os.Getenv("MONGODB_URI"),
		RedisUri:    os.Getenv("REDIS_URI"),
		GitHubToken: os.Getenv("GITHUB_TOKEN"),
	}, nil
}
