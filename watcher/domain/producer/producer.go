package producer

import (
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/radovskyb/watcher"
)

type Producer struct {
	Sender SendOutput
}

type SendOutput interface {
	Send(message *MessageEvent) error
}

type MessageEvent struct {
	ID       string `json:"id"`
	Event    string `json:"event"`
	FileName string `json:"fileName"`
	Size     int64  `json:"size"`
	Path     string `json:"path"`
}

func (p *Producer) SendEvent(event watcher.Event) error {
	filename := event.Path
	log.Println("Processing file:", filename)
	return p.send(filename)
}

func (p *Producer) send(filename string) error {
	file, err := os.Open(filename)

	if err != nil {
		log.Println("Error opening file:", err)
		return err
	}

	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		log.Println("Error getting file info:", err)
		return err
	}

	if !fileInfo.IsDir() {
		log.Println("Skipping directory:", filename)
		return err
	}

	messageEvent := MessageEvent{
		ID:       generateUUID(),
		Event:    "file_created",
		FileName: fileInfo.Name(),
		Size:     fileInfo.Size(),
		Path:     filename,
	}

	return p.Sender.Send(&messageEvent)
}

func generateUUID() string {
	return uuid.New().String()
}
