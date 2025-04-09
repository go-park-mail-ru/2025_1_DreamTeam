package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Database struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Name     string `yaml:"name"`
	} `yaml:"database"`

	Minio struct {
		Endpoint        string `yaml:"endpoint"`
		AccessKey       string `yaml:"access_key"`
		SecretAccessKey string `yaml:"secret_access_key"`
		BucketName      string `yaml:"bucket_name"`
		UseSSL          bool   `yaml:"use_ssl"`
	} `yaml:"minio"`

	Secrets struct {
		JwtSessionSecret string `yaml:"jwt_session_secret"`
	} `yaml:"secrets"`
}

// LoadConfig загружает конфигурацию из YAML-файла
func LoadConfig() *Config {
	data, err := os.ReadFile("../config/config.yaml")
	if err != nil {
		log.Fatalf("не удалось прочитать файл: %w", err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("ошибка парсинга YAML: %w", err)
	}

	return &config
}
