package main

import (
    "context"
    "log"
    "os/exec"
    "path/filepath"

    "github.com/segmentio/kafka-go"
)

func main() {
    kafkaTopic := getEnv("KAFKA_TOPIC", "folders")
    kafkaURL := getEnv("KAFKA_BROKER", "localhost:9092")
    outDir := getEnv("OUT_DIR", "./out")

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
            log.Fatal(err)
        }
        folder := string(m.Value)
        log.Println("processing", folder)

        // Example ffmpeg command: transcode all mp4 files in folder
        pattern := filepath.Join(folder, "*.mp4")
        files, _ := filepath.Glob(pattern)
        for _, f := range files {
            out := filepath.Join(outDir, filepath.Base(f))
            cmd := exec.Command("ffmpeg", "-i", f, out)
            if err := cmd.Run(); err != nil {
                log.Println("ffmpeg error:", err)
            }
        }
    }
}

func getEnv(key, def string) string {
    if v := os.Getenv(key); v != "" {
        return v
    }
    return def
}
