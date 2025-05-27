package yookassa

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	billingpb "skillForce/internal/delivery/grpc/proto"
	"time"

	"github.com/google/uuid"
)

type BillingServer struct {
	billingpb.UnimplementedBillingServiceServer
	shopID    string
	secretKey string
}

func NewBillingServer(shopID, secretKey string) *BillingServer {
	return &BillingServer{
		shopID:    shopID,
		secretKey: secretKey,
	}
}

func (s *BillingServer) CreatePayment(returnUrl string, title string, userID int32, courseID int32, amount int) (string, *billingpb.CreatePaymentResponse, error) {
	payment := map[string]interface{}{
		"amount": map[string]string{
			"value":    fmt.Sprintf("%.2f", float64(amount)),
			"currency": "RUB",
		},
		"confirmation": map[string]string{
			"type":       "redirect",
			"return_url": returnUrl,
		},
		"capture":     true,
		"description": fmt.Sprintf("Оплата курса %s", title),
	}

	body, _ := json.Marshal(payment)
	auth := base64.StdEncoding.EncodeToString([]byte(s.shopID + ":" + s.secretKey))
	idempotenceKey := uuid.New().String()

	reqHTTP, _ := http.NewRequest("POST", "https://api.yookassa.ru/v3/payments", bytes.NewBuffer(body))
	reqHTTP.Header.Set("Authorization", "Basic "+auth)
	reqHTTP.Header.Set("Content-Type", "application/json")
	reqHTTP.Header.Set("Idempotence-Key", idempotenceKey)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(reqHTTP)
	if err != nil {
		return "", nil, fmt.Errorf("error in request to yookassa: %v", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("failed to close response body: %v", err)
		}
	}()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", nil, fmt.Errorf("reading yookassa error: %v", err)
	}

	billing_id, ok := result["id"].(string)
	if !ok {
		return "", nil, fmt.Errorf("cannot get billing id")
	}

	fmt.Print("OKKKKK")
	confirmation := result["confirmation"].(map[string]interface{})
	confirmation_url := confirmation["confirmation_url"].(string)
	return billing_id, &billingpb.CreatePaymentResponse{
		ConfirmationUrl: confirmation_url,
	}, nil
}

func (s *BillingServer) HandleWebhook(ctx context.Context, req *billingpb.YooKassaWebhook) (bool, error) {
	log.Printf("[Webhook] Event: %s, PaymentID: %s, Status: %s", req.Event, req.PaymentId, req.Status)
	if req.Status == "succeeded" {
		return true, nil
	}

	return false, nil
}
