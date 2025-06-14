package main

import (
	"os"

	"github.com/example/go-streaming/watcher/domain/producer"
	"github.com/example/go-streaming/watcher/domain/watcher"
	"github.com/example/go-streaming/watcher/infra/output/kafka_sender"
	"github.com/joho/godotenv"
)

func main() {
	err:= godotenv.Load("./.env")

	if err != nil {
		panic("Error loading .env file")
	}

	producer := producer.Producer{
		Sender: &kafka_sender.KafkaSender{},
	}
	watcher := watcher.Watcher{
		FolderPath: os.Getenv("WATCH_DIR"),
	}

	watcher.Start(&producer)
}
