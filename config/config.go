package config

import (

	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	Address string
	SlackToken string
}

func LoadConfig(filePath string)(*Config, error){
	err := godotenv.Load(filePath)
	if err != nil {
		return nil, err
	}
	config := &Config{
		Address: os.Getenv("ADDRESS"),
		SlackToken: os.Getenv("SLACK_TOKEN"),
	}
	return config, nil
}
