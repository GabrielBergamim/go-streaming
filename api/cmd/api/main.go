package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Video struct {
	ID   string `json:"id" gorm:"column:id;primaryKey"`
	Name string `json:"name" gorm:"column:file_name"`
}

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Println("Error loading .env file")
	}

	dsn := os.Getenv("POSTGRES_DSN")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	db.AutoMigrate(&Video{})

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))

	app.Static("/videos", "./public/videos")

	app.Get("/videos", func(c *fiber.Ctx) error {
		name := c.Query("name")
		var videos []Video

		query := db.Model(&Video{})
		if name != "" {
			query = query.Where("file_name = ?", name)
		}
		if err := query.Find(&videos).Error; err != nil {
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
