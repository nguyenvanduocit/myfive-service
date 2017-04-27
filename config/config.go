package config

import (
	"github.com/joho/godotenv"
	"os"
	"time"
)

type Config struct {
	Address    string
	SlackToken string
	CrawlInterval time.Duration
}

func LoadConfig(filePath string) (*Config, error) {
	err := godotenv.Load(filePath)
	if err != nil {
		return nil, err
	}
	config := &Config{
		Address:       os.Getenv("ADDRESS"),
		SlackToken:    os.Getenv("SLACK_TOKEN"),
		CrawlInterval: 10 * time.Minute,
	}
	return config, nil
}
