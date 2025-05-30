package usecase

import (
	"context"
	"fmt"
	billingpb "skillForce/internal/delivery/grpc/proto"
	"skillForce/pkg/logs"

	"google.golang.org/protobuf/types/known/emptypb"
)

type BillingUsecase struct {
	repo BillingRepository
}

func NewBillingUsecase(repo BillingRepository) *BillingUsecase {
	return &BillingUsecase{
		repo: repo,
	}
}

func (uc *BillingUsecase) CreatePayment(ctx context.Context, userId int, courseID int, returnUrl string) (*billingpb.CreatePaymentResponse, error) {
	title, amount, err := uc.repo.GetBillingInfo(ctx, courseID)
	if err != nil {
		logs.PrintLog(ctx, "GetBillingInfo", fmt.Sprintf("%+v", err))
		return nil, err
	}
	billing_id, response, err := uc.repo.CreatePayment(returnUrl, title, int32(userId), int32(courseID), amount)
	if err != nil {
		logs.PrintLog(ctx, "GetBillingInfo", fmt.Sprintf("%+v", err))
		return nil, err
	}
	err = uc.repo.AddNewBilling(ctx, userId, courseID, billing_id)
	if err != nil {
		logs.PrintLog(ctx, "GetBillingInfo", fmt.Sprintf("%+v", err))
		return nil, err
	}
	return response, nil
}

func (uc *BillingUsecase) UpdateBilling(ctx context.Context, req *billingpb.YooKassaWebhook) (*emptypb.Empty, error) {
	is_successed, err := uc.repo.HandleWebhook(ctx, req)
	if err != nil {
		logs.PrintLog(ctx, "UpdateBilling", fmt.Sprintf("%+v", err))
		return nil, err
	}

	if is_successed {
		userID, courseID, err := uc.repo.UpdateBilling(ctx, req.PaymentId)
		if err != nil {
			logs.PrintLog(ctx, "UpdateBilling", fmt.Sprintf("%+v", err))
			return nil, err
		}

		user, err := uc.repo.GetUserById(ctx, userID)
		if err != nil {
			logs.PrintLog(ctx, "UpdateBilling", fmt.Sprintf("%+v", err))
			return nil, err
		}

		course, err := uc.repo.GetCourseById(ctx, courseID)
		if err != nil {
			logs.PrintLog(ctx, "UpdateBilling", fmt.Sprintf("%+v", err))
			return nil, err
		}

		go func() {
			err = uc.repo.SendReceipt(ctx, user, req.PaymentId, course)
			if err != nil {
				logs.PrintLog(ctx, "UpdateBilling", fmt.Sprintf("problem with sending receipt: %+v", err))
			}
		}()
		return &emptypb.Empty{}, nil
	}
	return &emptypb.Empty{}, nil
}
