package main

import (
	"encoding/json"
	"log"

	"skillForce/config"
	"skillForce/mail"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
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

	defer consumer.Close()

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

func handleMessage(mailLient *mail.Mail, msg *kafka.Message) {
	var message mail.KafkaMessage
	err := json.Unmarshal(msg.Value, &message)
	if err != nil {
		log.Printf("Failed to unmarshal JSON: %v\n", err)
		return
	}

	log.Printf("Received message: %+v\n", message)

	switch message.Method {
	case "send_confirm_mail":
		mailLient.SendRegMail(nil, message)
	case "send_welcome_mail":
		// TODO: send_welcome_mail
	case "send_welcome_course_mail":
		// TODO: send_welcome_course_mail
	case "send_middle_course_mail":
		// TODO: send_middle_course_mail
	}
}
