package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port             int    `envconfig:"port" default:"9001"`
	Url              string `envconfig:"url" default:"http://localhost"`
	PostgresHost     string `envconfig:"postgres_host" default:"localhost"`
	PostgresPort     int    `envconfig:"postgres_port" default:"5434"`
	PostgresDatabase string `envconfig:"postgres_database" default:"postgres"`
	PostgresUser     string `envconfig:"postgres_user" default:"root"`
	PostgresPassword string `envconfig:"postgres_password" default:"root-is-not-used"`
}

func NewConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Warning: Error loading .env file: %v", err)
	}

	var conf Config
	err = envconfig.Process("", &conf)
	if err != nil {
		log.Fatalf("fail to proceed the config: %v", err)
	}
	return &conf
}
