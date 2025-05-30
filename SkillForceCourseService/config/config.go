package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Database struct {
		Host     string
		Port     string
		User     string
		Password string
		Name     string
	}

	Minio struct {
		Endpoint          string
		AccessKey         string
		SecretAccessKey   string
		BucketName        string
		CoursesBucketName string
		UseSSL            bool
	}

	Secrets struct {
		JwtSessionSecret string
	}

	Mail struct {
		From     string
		Password string
		Host     string
		Port     string
	}
}

type yamlConfig struct {
	Database struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
		Name string `yaml:"name"`
	} `yaml:"database"`

	Minio struct {
		Endpoint          string `yaml:"endpoint"`
		BucketName        string `yaml:"bucket_name"`
		CoursesBucketName string `yaml:"courses_bucket_name"`
		UseSSL            bool   `yaml:"use_ssl"`
	} `yaml:"minio"`
}

func LoadConfig() *Config {
	err := godotenv.Load("./config/.env")
	if err != nil {
		log.Fatalf("ошибка загрузки .env файла: %v", err)
	}

	data, err := os.ReadFile("./config/config.yaml")
	if err != nil {
		log.Fatalf("не удалось прочитать YAML файл: %v", err)
	}

	var ycfg yamlConfig
	err = yaml.Unmarshal(data, &ycfg)
	if err != nil {
		log.Fatalf("ошибка парсинга YAML: %v", err)
	}

	return &Config{
		Database: struct {
			Host     string
			Port     string
			User     string
			Password string
			Name     string
		}{
			Host:     ycfg.Database.Host,
			Port:     ycfg.Database.Port,
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     ycfg.Database.Name,
		},
		Minio: struct {
			Endpoint          string
			AccessKey         string
			SecretAccessKey   string
			BucketName        string
			CoursesBucketName string
			UseSSL            bool
		}{
			Endpoint:          ycfg.Minio.Endpoint,
			AccessKey:         os.Getenv("MINIO_ACCESS_KEY"),
			SecretAccessKey:   os.Getenv("MINIO_SECRET_KEY"),
			BucketName:        ycfg.Minio.BucketName,
			CoursesBucketName: ycfg.Minio.CoursesBucketName,
			UseSSL:            ycfg.Minio.UseSSL,
		},
		Secrets: struct{ JwtSessionSecret string }{
			JwtSessionSecret: os.Getenv("JWT_SESSION_SECRET"),
		},
		Mail: struct {
			From     string
			Password string
			Host     string
			Port     string
		}{
			From:     os.Getenv("MAIL_FROM"),
			Password: os.Getenv("MAIL_PASSWORD"),
			Host:     os.Getenv("MAIL_HOST"),
			Port:     os.Getenv("MAIL_PORT"),
		},
	}
}
