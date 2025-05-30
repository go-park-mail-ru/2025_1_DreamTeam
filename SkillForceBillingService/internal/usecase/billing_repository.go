package usecase

import (
	"context"
	billingpb "skillForce/internal/delivery/grpc/proto"
	coursemodels "skillForce/internal/models/course"
	usermodels "skillForce/internal/models/user"
)

type BillingRepository interface {
	AddNewBilling(ctx context.Context, userID int, courseID int, billing_id string) error
	UpdateBilling(ctx context.Context, billing_id string) (int, int, error)
	GetBillingInfo(ctx context.Context, courseID int) (string, int, error)
	CreatePayment(returnUrl string, title string, userID int32, courseID int32, amount int) (string, *billingpb.CreatePaymentResponse, error)
	HandleWebhook(ctx context.Context, req *billingpb.YooKassaWebhook) (bool, error)

	GetUserById(ctx context.Context, userId int) (*usermodels.User, error)
	GetCourseById(ctx context.Context, courseID int) (*coursemodels.Course, error)

	SendReceipt(ctx context.Context, user *usermodels.User, token string, course *coursemodels.Course) error
}
