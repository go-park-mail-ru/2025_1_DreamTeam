package grpc

import (
	"context"
	billingpb "skillForce/internal/delivery/grpc/proto"
	"skillForce/internal/usecase"

	"google.golang.org/protobuf/types/known/emptypb"
)

type BillingHandler struct {
	billingpb.UnimplementedBillingServiceServer
	usecase *usecase.BillingUsecase
}

func NewBillingHandler(uc *usecase.BillingUsecase) *BillingHandler {
	return &BillingHandler{
		usecase: uc,
	}
}

func (h *BillingHandler) CreatePayment(ctx context.Context, req *billingpb.CreatePaymentRequest) (*billingpb.CreatePaymentResponse, error) {
	response, err := h.usecase.CreatePayment(ctx, int(req.UserId), int(req.CourseId), req.ReturnUrl)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (h *BillingHandler) HandleWebhook(ctx context.Context, req *billingpb.YooKassaWebhook) (*emptypb.Empty, error) {
	_, err := h.usecase.UpdateBilling(ctx, req)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
