package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	HTTPServer `yaml:"server"`
	Database   `yaml:"postgres"`
	Redis      `yaml:"redis"`
}

type HTTPServer struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type Database struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"user"`
	DBname   string `yaml:"dbname"`
	SSLMode  string `yaml:"ssl"`
}

type Redis struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

func LoadConfig() *Config {
	godotenv.Load()
	configPath := os.Getenv("CONFIG_PATH")

	if configPath == "" {
		log.Fatal("CONFIG_PATH isn't set in .env")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatal("Config file doesn't exists: ", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatal("Cannot read config: ", err)

	}

	return &cfg
}
