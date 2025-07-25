package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/segmentio/kafka-go"

	"github.com/example/go-streaming/processor/application"
	"github.com/example/go-streaming/processor/domain/message"
	"github.com/example/go-streaming/processor/infra/persistence"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file")
	}

	kafkaTopic := os.Getenv("KAFKA_TOPIC")
	kafkaURL := os.Getenv("KAFKA_BROKER")
	outDir := os.Getenv("OUTPUT_DIR")
	dsn := os.Getenv("POSTGRES_DSN")

	log.Println("Kafka Topic:", kafkaTopic)
	log.Println("Kafka Broker:", kafkaURL)

	db, err := persistence.NewDB(dsn)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	repo := persistence.NewGormVideoRepository(db)
	processor := application.Processor{Repo: repo}

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{kafkaURL},
		Topic:   kafkaTopic,
		GroupID: "video-processor",
	})
	defer r.Close()

	log.Println("Processor waiting for messages...")

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Println("Failed to connect")
			time.Sleep(5 * time.Second)
			continue
		}

		var msg message.Event
		if err := json.Unmarshal(m.Value, &msg); err != nil {
			log.Println("failed to decode message:", err)
			continue
		}

		log.Printf("Received message: ID=%s, Event=%s, FileName=%s, Size=%d, Path=%s\n",
			msg.ID, msg.Event, msg.FileName, msg.Size, msg.Path)
		processor.HandleEvent(&msg, outDir)
	}
}
