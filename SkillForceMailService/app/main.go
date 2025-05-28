package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"skillForce/config"
	"skillForce/mail"
	"skillForce/metrics"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	config := config.LoadConfig()

	mailClient := mail.NewMail(config.Mail.From, config.Mail.Password, config.Mail.Host, config.Mail.Port)
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "kafka:9092",
		"group.id":          "mail-service",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		log.Fatalf("Failed to create consumer: %s", err)
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			log.Printf("Failed to close Kafka consumer: %v", err)
		}
	}()

	metrics.Init()
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Println("Prometheus metrics available at :9083/metrics")
		if err := http.ListenAndServe(":9083", nil); err != nil {
			log.Fatalf("failed to start metrics HTTP server: %v", err)
		}
	}()

	topic := "mail"
	err = consumer.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		log.Fatalf("Failed to subscribe to topic: %s", err)
	}

	log.Println("Waiting for Kafka messages...")

	for {
		msg, err := consumer.ReadMessage(-1) // -1 = блокировать до получения сообщения
		if err == nil {
			handleMessage(mailClient, msg)
		} else {
			log.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}
}

func handleMessage(mailClient *mail.Mail, msg *kafka.Message) {
	var message mail.KafkaMessage
	if err := json.Unmarshal(msg.Value, &message); err != nil {
		log.Printf("Failed to unmarshal JSON: %v\n", err)
		return
	}

	log.Printf("Received message: %+v\n", message)

	ctx := context.Background()
	var sendErr error

	switch message.Method {
	case "send_confirm_mail":
		sendErr = mailClient.SendRegMail(ctx, message)
	case "send_welcome_course_mail":
		sendErr = mailClient.SendWelcomeCourseMail(ctx, message)
	case "send_middle_course_mail":
		// TODO: implement send_middle_course_mail
	}

	if sendErr != nil {
		log.Printf("Failed to send %s mail: %v", message.Method, sendErr)
	}
}
