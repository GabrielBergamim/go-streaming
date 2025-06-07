package main

import (
    "log"
    "github.com/gofiber/fiber/v2"
)

func main() {
    app := fiber.New()

    // Serve files from the "videos" directory
    app.Static("/videos", "./videos")

    log.Println("Starting API on :3000")
    if err := app.Listen(":3000"); err != nil {
        log.Fatal(err)
    }
}
