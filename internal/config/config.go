package config

import (
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	Addr     string
	LogLevel string
	Storage
}
type Storage struct {
	User     string
	Password string
	Host     string
	Port     string
	DBName   string
	SSLMode  string
}

func MustLoad(path string) *Config {
	if err := godotenv.Load(path); err != nil {
		panic(".env file not found")
	}
	stor := Storage{
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}
	return &Config{
		Addr:     os.Getenv("ADDR"),
		LogLevel: os.Getenv("LOG_LEVEL"),
		Storage:  stor,
	}
}
