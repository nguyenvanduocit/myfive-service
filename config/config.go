package config

import (

	"github.com/joho/godotenv"
	"os"
	"strconv"
)

type Config struct {
	DatabaseName string
	DatabaseHost string
	DatabasePort int
	DatabaseUserName string
	DatabasePassword string
	Address string
}

func LoadConfig(filePath string)(*Config, error){
	err := godotenv.Load(filePath)
	if err != nil {
		return nil, err
	}
	port, err := strconv.Atoi(os.Getenv("DATABASE_PORT"));
	if err != nil {
		return nil, err
	}
	config := &Config{
		DatabaseName: os.Getenv("DATABASE_NAME"),
		DatabaseHost: os.Getenv("DATABASE_HOST"),
		DatabaseUserName: os.Getenv("DATABASE_USERNAME"),
		DatabasePassword: os.Getenv("DATABASE_PASSWORD"),
		Address: os.Getenv("ADDRESS"),
		DatabasePort:port,
	}
	return config, nil
}
