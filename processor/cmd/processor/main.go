package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/segmentio/kafka-go"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type MessageEvent struct {
	ID       string `json:"id"`
	Event    string `json:"event"`
	FileName string `json:"fileName"`
	Size     int64  `json:"size"`
	Path     string `json:"path"`
}

type Video struct {
	ID       string `json:"id" gorm:"column:id;primaryKey"`
	Event    string `json:"event" gorm:"column:event"`
	FileName string `json:"fileName" gorm:"column:file_name"`
	Size     int64  `json:"size" gorm:"column:size"`
	Path     string `json:"path" gorm:"column:path"`
}

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Println("Error loading .env file")
	}

	kafkaTopic := os.Getenv("KAFKA_TOPIC")
	kafkaURL := os.Getenv("KAFKA_BROKER")
	outDir := os.Getenv("OUTPUT_DIR")
	dsn := os.Getenv("POSTGRES_DSN")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	db.AutoMigrate(&Video{})

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

		var messageValue MessageEvent
		err = json.Unmarshal(m.Value, &messageValue)

		if err != nil {
			log.Fatal(err)
		}

		folder := string(messageValue.Path)
		log.Println("processing", folder)
		outDir = filepath.Join(outDir, messageValue.ID)
		log.Println("output directory will be:", outDir)
		os.MkdirAll(outDir, os.ModePerm)

		processVideos(folder, outDir)
		processSubtitles(folder, outDir)

		video := Video{
			ID:       messageValue.ID,
			Event:    messageValue.Event,
			FileName: messageValue.FileName,
			Size:     messageValue.Size,
			Path:     messageValue.Path,
		}
		if err := db.Create(&video).Error; err != nil {
			log.Println("failed to store video:", err)
		}
	}
}

func processVideos(folder, outDir string) {
	extensions := []string{"mkv", "mp4", "avi"}

	for _, ext := range extensions {
		pattern := filepath.Join(folder, "*."+ext)
		files, _ := filepath.Glob(pattern)

		for _, f := range files {
			log.Println("processing video file:", f)
			out := filepath.Join(outDir, "video")
			log.Println("output file will be:", out+".m3u8")
			log.Println("segment filename will be:", fmt.Sprintf("%s_segment_%%03d.ts", out))
			cmd := exec.Command("ffmpeg",
				"-i", f,
				"-preset", "slow",
				"-c:v", "libx264",
				"-crf", "18",
				"-hls_time", "10",
				"-hls_list_size", "0",
				"-hls_segment_filename", fmt.Sprintf("%s_segment_%%03d.ts", out),
				out+".m3u8")

			if err := cmd.Run(); err != nil {
				log.Println("ffmpeg error:", err)
			}
		}
	}
}

func processSubtitles(folder, outDir string) {
	extensions := []string{"srt"}

	for _, ext := range extensions {
		pattern := filepath.Join(folder, "*."+ext)
		files, _ := filepath.Glob(pattern)

		for _, f := range files {
			log.Println("processing subtitle file:", f)
			out := filepath.Join(outDir, "pt-BR")
			log.Println("output file will be:", out+".vtt")
			cmd := exec.Command("ffmpeg", "-i", f, out+".vtt")

			if err := cmd.Run(); err != nil {
				log.Println("ffmpeg error:", err)
			}
		}
	}
}
