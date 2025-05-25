package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	billingpb "skillForce/internal/delivery/grpc/proto/billing"
	"skillForce/internal/delivery/http/response"
	models "skillForce/internal/models/user"
	"skillForce/pkg/logs"

	"google.golang.org/grpc"
)

type CookieManagerInterface interface {
	CheckCookie(r *http.Request) *models.UserProfile
}

type Handler struct {
	billingClient billingpb.BillingServiceClient
	cookieManager CookieManagerInterface
}

func NewHandler(cookieManager CookieManagerInterface) *Handler {
	conn, err := grpc.Dial("billing-service:8084", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect to user service: %v", err)
	}
	billingClient := billingpb.NewBillingServiceClient(conn)
	return &Handler{
		billingClient: billingClient,
		cookieManager: cookieManager,
	}
}

func (h *Handler) CreatePaymentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		logs.PrintLog(r.Context(), "CreatePaymentHandler", "method not allowed")
		response.SendErrorResponse("method not allowed", http.StatusMethodNotAllowed, w, r)
		return
	}

	userProfile := h.cookieManager.CheckCookie(r)

	if userProfile == nil {
		logs.PrintLog(r.Context(), "CreatePaymentHandler", "user not logged in")
		userProfile = &models.UserProfile{Id: -1}
		response.SendErrorResponse("user not logged in", http.StatusUnauthorized, w, r)
		return
	}

	var req struct {
		ReturnURL string `json:"return_url"`
		User_ID   int32
		CourseID  int32 `json:"course_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logs.PrintLog(r.Context(), "CreatePaymentHandler", "invalid JSON")
		response.SendErrorResponse("invalid JSON", http.StatusBadRequest, w, r)
		return
	}

	req.User_ID = int32(userProfile.Id)

	resp, err := h.billingClient.CreatePayment(context.Background(), &billingpb.CreatePaymentRequest{
		ReturnUrl: req.ReturnURL,
		UserId:    req.User_ID,
		CourseId:  req.CourseID,
	})
	if err != nil {
		logs.PrintLog(r.Context(), "CreatePaymentHandler", "gRPC error: "+err.Error())
		response.SendErrorResponse("gRPC error: "+err.Error(), http.StatusInternalServerError, w, r)
		return
	}

	response.SendBillingRedirect(w, r, resp.ConfirmationUrl)
}

func (h *Handler) WebhookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		logs.PrintLog(r.Context(), "WebhookHandler", "method not allowed")
		response.SendErrorResponse("method not allowed", http.StatusMethodNotAllowed, w, r)
		return
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		logs.PrintLog(r.Context(), "WebhookHandler", "invalid body")
		response.SendErrorResponse("invalid body", http.StatusBadRequest, w, r)
		return
	}
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	var data struct {
		Event  string `json:"event"`
		Object struct {
			ID     string `json:"id"`
			Status string `json:"status"`
		} `json:"object"`
	}
	if err := json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(&data); err != nil {
		logs.PrintLog(r.Context(), "WebhookHandler", "invalid JSON")
		response.SendErrorResponse("invalid JSON", http.StatusBadRequest, w, r)
		return
	}

	_, err = h.billingClient.HandleWebhook(context.Background(), &billingpb.YooKassaWebhook{
		Event:      data.Event,
		PaymentId:  data.Object.ID,
		Status:     data.Object.Status,
		RawPayload: string(bodyBytes),
	})
	if err != nil {
		logs.PrintLog(r.Context(), "WebhookHandler", "webhook error: "+err.Error())
		response.SendErrorResponse("webhook error: "+err.Error(), http.StatusBadRequest, w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
}
