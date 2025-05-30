package repository

import (
	"context"
	"fmt"
	"log"
	"skillForce/config"
	billingpb "skillForce/internal/delivery/grpc/proto"
	coursemodels "skillForce/internal/models/course"
	usermodels "skillForce/internal/models/user"
	"skillForce/internal/repository/kafka"
	"skillForce/internal/repository/postgres"
	"skillForce/internal/repository/yookassa"
)

type BillingInfrastructure struct {
	Database      *postgres.Database
	Billing       *yookassa.BillingServer
	KafkaProducer *kafka.Producer
}

func NewBillingInfrastructure(conf *config.Config) *BillingInfrastructure {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", conf.Database.Host, conf.Database.Port, conf.Database.User, conf.Database.Password, conf.Database.Name)
	database, err := postgres.NewDatabase(dsn, conf.Secrets.JwtSessionSecret)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	billingServer := yookassa.NewBillingServer(conf.Yookassa.ShopID, conf.Yookassa.SecretKey)

	kafkaProducer := kafka.NewKafkaProducer()
	return &BillingInfrastructure{
		Database:      database,
		Billing:       billingServer,
		KafkaProducer: kafkaProducer,
	}
}

func (i *BillingInfrastructure) Close() {
	if err := i.Database.Close(); err != nil {
		log.Fatal(err)
	}
}

func (i *BillingInfrastructure) AddNewBilling(ctx context.Context, userID int, courseID int, billing_id string) error {
	return i.Database.AddNewBilling(ctx, userID, courseID, billing_id)
}

func (i *BillingInfrastructure) CreatePayment(returnUrl string, title string, userID int32, courseID int32, amount int) (string, *billingpb.CreatePaymentResponse, error) {
	return i.Billing.CreatePayment(returnUrl, title, userID, courseID, amount)
}

func (i *BillingInfrastructure) UpdateBilling(ctx context.Context, billing_id string) (int, int, error) {
	return i.Database.UpdateBilling(ctx, billing_id)
}

func (i *BillingInfrastructure) HandleWebhook(ctx context.Context, req *billingpb.YooKassaWebhook) (bool, error) {
	return i.Billing.HandleWebhook(ctx, req)
}

func (i *BillingInfrastructure) GetBillingInfo(ctx context.Context, courseID int) (string, int, error) {
	return i.Database.GetBillingInfo(ctx, courseID)
}

func (i *BillingInfrastructure) SendReceipt(ctx context.Context, user *usermodels.User, token string, course *coursemodels.Course) error {
	return i.KafkaProducer.SendReceipt(ctx, user, token, course)
}

func (i *BillingInfrastructure) GetUserById(ctx context.Context, userId int) (*usermodels.User, error) {
	return i.Database.GetUserById(ctx, userId)
}

func (i *BillingInfrastructure) GetCourseById(ctx context.Context, courseID int) (*coursemodels.Course, error) {
	return i.Database.GetCourseById(ctx, courseID)
}
