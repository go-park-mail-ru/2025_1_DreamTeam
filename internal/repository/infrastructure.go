package repository

import (
	"fmt"
	"log"
	"skillForce/config"
	"skillForce/internal/repository/mail"
	"skillForce/internal/repository/minio"
	"skillForce/internal/repository/postgres"
)

type Infrastructure struct {
	Database *postgres.Database
	Minio    *minio.Minio
	Mail     *mail.Mail
}

func NewInfrastructure(conf *config.Config) *Infrastructure {
	mn, err := minio.NewMinio(conf.Minio.Endpoint, conf.Minio.AccessKey, conf.Minio.SecretAccessKey, conf.Minio.UseSSL, conf.Minio.BucketName, conf.Minio.VideoBucket)
	if err != nil {
		log.Fatalf("Failed to connect to MinIO: %v", err)
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", conf.Database.Host, conf.Database.Port, conf.Database.User, conf.Database.Password, conf.Database.Name)
	database, err := postgres.NewDatabase(dsn, conf.Secrets.JwtSessionSecret)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	mail := mail.NewMail(conf.Mail.From, conf.Mail.Password, conf.Mail.Host, conf.Mail.Port)
	if err != nil {
		log.Fatalf("Failed to connect to mail: %v", err)
	}

	return &Infrastructure{
		Database: database,
		Minio:    mn,
		Mail:     mail,
	}
}

func (i *Infrastructure) Close() {
	i.Database.Close()
}
