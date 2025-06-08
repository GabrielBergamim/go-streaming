package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
		return
	}

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))

	port := os.Getenv("PORT")

	app.Static("/videos", port)

	log.Printf("Starting API on %s\n", port)

	if err := app.Listen(":" + port); err != nil {
		log.Fatal(err)
	}
}
