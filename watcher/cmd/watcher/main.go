package main

import (
	"os"

	"github.com/example/go-streaming/watcher/domain/producer"
	"github.com/example/go-streaming/watcher/domain/watcher"
	"github.com/example/go-streaming/watcher/infra/output/kafka_sender"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	producer := producer.Producer{
		Sender: &kafka_sender.KafkaSender{},
	}
	watcher := watcher.Watcher{
		FolderPath: os.Getenv("WATCH_DIR"),
	}

	watcher.Start(&producer)
}
