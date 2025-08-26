package main

import (
	"log"
	"os"

	"github.com/example/go-streaming/api/infra/server"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Println("Error loading .env file")
	}

	server := server.NewServer(os.Getenv("PORT"), os.Getenv("POSTGRES_DSN"))
	err = server.Start()

	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
