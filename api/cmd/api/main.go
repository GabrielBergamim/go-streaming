package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"

	"github.com/example/go-streaming/api/infra/persistence"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Println("Error loading .env file")
	}

	dsn := os.Getenv("POSTGRES_DSN")

	db, err := persistence.NewDB(dsn)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	repo := persistence.NewGormVideoRepository(db)

	app := fiber.New()

	app.Use(cors.New(cors.Config{AllowOrigins: "*"}))

	app.Use("/api/video", func(c *fiber.Ctx) error {
	    c.Set("Access-Control-Allow-Origin", "*")
	    return c.Next()
	})

	app.Static("/api/video", "./public/videos", fiber.Static{
		Compress: false, ByteRange: true,})

	app.Get("/api/videos", func(c *fiber.Ctx) error {
		name := c.Query("name")
		videos, err := repo.FindByName(name)
		if err != nil {
			return err
		}
		return c.JSON(videos)
	})

	port := os.Getenv("PORT")
	log.Printf("Starting API on %s\n", port)

	if err := app.Listen(":" + port); err != nil {
		log.Fatal(err)
	}
}
