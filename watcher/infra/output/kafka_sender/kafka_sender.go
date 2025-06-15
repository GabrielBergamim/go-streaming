package kafka_sender

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"

	"github.com/example/go-streaming/watcher/domain/producer"
	"github.com/segmentio/kafka-go"
)

type KafkaSender struct{}

func (kf KafkaSender) Send(messageEvent *producer.MessageEvent) error {
	writer, err := getKafkaInstance()

	if err != nil {
		return err
	}

    defer writer.Close()

	jsonValue, err := json.Marshal(&messageEvent)

	if err != nil {
		log.Fatalf("failed to marshal message: %v", err)
		return err
	}

	message := kafka.Message{
		Key:   []byte("key"),
		Value: jsonValue,
	}

	return writer.WriteMessages(context.Background(), message)
}

func getKafkaInstance() (*kafka.Writer, error) {
	kafkaTopic := os.Getenv("KAFKA_TOPIC")
	kafkaURL := os.Getenv("KAFKA_BROKER")

	if kafkaTopic == "" || kafkaURL == "" {
		return nil, errors.New("Kafka topic or broker URL not set")
	}

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{kafkaURL},
		Topic:   kafkaTopic,
		Balancer: &kafka.LeastBytes{},
	})

	return writer, nil
}
