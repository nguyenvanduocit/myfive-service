package config

import (
	"github.com/joho/godotenv"
	"os"
	"time"
)

type Config struct {
	Address       string
	SlackToken    string
	CrawlInterval time.Duration
}

func LoadConfig(filePath string) (*Config, error) {
	err := godotenv.Load(filePath)
	if err != nil {
		return nil, err
	}
	data := &Config{
		Address:       os.Getenv("ADDRESS"),
		SlackToken:    os.Getenv("SLACK_TOKEN"),
		CrawlInterval: 15 * time.Minute,
	}
	return data, nil
}
