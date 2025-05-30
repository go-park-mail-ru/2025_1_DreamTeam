package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	coursemodels "skillForce/internal/models/course"
	usermodels "skillForce/internal/models/user"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Producer struct {
	producer *kafka.Producer
	topik    string
}

type KafkaMessage struct {
	Method     string
	Token      string
	UserEmail  string
	UserName   string
	CourseName string
	CourseId   int
	Url        string
}

func NewKafkaProducer() *Producer {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "kafka:9092",
	})
	if err != nil {
		log.Fatalf("Failed to create producer: %s", err)
	}

	return &Producer{producer: producer, topik: "mail"}
}

func (p *Producer) Close() {
	p.producer.Close()
}

func (p *Producer) SendReceipt(ctx context.Context, user *usermodels.User, token string, course *coursemodels.Course) error {
	msg := KafkaMessage{
		Method:     "send_receipt_mail",
		Token:      token,
		UserEmail:  user.Email,
		UserName:   user.Name,
		CourseName: course.Title,
	}

	value, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("SendRegMail", err.Error())
	}

	err = p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &p.topik, Partition: kafka.PartitionAny},
		Value:          value,
	}, nil)

	if err != nil {
		fmt.Println("SendRegMail", err.Error())
	}

	// ожидаем подтверждение доставки
	e := <-p.producer.Events()
	switch ev := e.(type) {
	case *kafka.Message:
		if ev.TopicPartition.Error != nil {
			fmt.Printf("Delivery failed: %v\n", ev.TopicPartition.Error)
		} else {
			fmt.Printf("Message delivered to %v\n", ev.TopicPartition)
		}
	}
	return nil
}
