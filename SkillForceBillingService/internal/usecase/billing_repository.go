package usecase

import (
	"context"
	billingpb "skillForce/internal/delivery/grpc/proto"
)

type BillingRepository interface {
	AddNewBilling(ctx context.Context, userID int, courseID int, billing_id string) error
	UpdateBilling(ctx context.Context, billing_id string) error
	GetBillingInfo(ctx context.Context, courseID int) (string, int, error)
	CreatePayment(returnUrl string, title string, userID int32, courseID int32, amount int) (string, *billingpb.CreatePaymentResponse, error)
	HandleWebhook(ctx context.Context, req *billingpb.YooKassaWebhook) (bool, error)
}
