package main

import (
    "context"
    "log"
    "os"

    "github.com/fsnotify/fsnotify"
    "github.com/segmentio/kafka-go"
)

func main() {
    watchDir := os.Getenv("WATCH_DIR")
    if watchDir == "" {
        watchDir = "./watch"
    }

    kafkaTopic := os.Getenv("KAFKA_TOPIC")
    if kafkaTopic == "" {
        kafkaTopic = "folders"
    }

    kafkaURL := os.Getenv("KAFKA_BROKER")
    if kafkaURL == "" {
        kafkaURL = "localhost:9092"
    }

    writer := kafka.NewWriter(kafka.WriterConfig{
        Brokers: []string{kafkaURL},
        Topic:   kafkaTopic,
    })
    defer writer.Close()

    watcher, err := fsnotify.NewWatcher()
    if err != nil {
        log.Fatal(err)
    }
    defer watcher.Close()

    if err := watcher.Add(watchDir); err != nil {
        log.Fatal(err)
    }

    log.Println("Watching", watchDir)
    for {
        select {
        case event, ok := <-watcher.Events:
            if !ok {
                return
            }
            if event.Op&fsnotify.Create == fsnotify.Create {
                // Send folder path to Kafka
                msg := kafka.Message{Value: []byte(event.Name)}
                if err := writer.WriteMessages(context.Background(), msg); err != nil {
                    log.Println("failed to write message:", err)
                } else {
                    log.Println("sent event for", event.Name)
                }
            }
        case err, ok := <-watcher.Errors:
            if !ok {
                return
            }
            log.Println("watch error:", err)
        }
    }
}
