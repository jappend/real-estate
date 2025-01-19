package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
  // Initializing environment variables
  err := godotenv.Load() 
  if err != nil {
    log.Fatal(err)
  }

  config := fiber.Config{
    AppName: "Real Estate",
    ReadTimeout: 30 * time.Second,
    WriteTimeout: 90 * time.Second,
    IdleTimeout: 120 * time.Second,
  } 

	app := fiber.New(config)
  app.Use(logger.New(logger.Config{
    Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
  }))

  // Routes
  app.Get("/helloworld", func(c *fiber.Ctx) error {
    c.Status(fiber.StatusOK).JSON("Hello World!")
    return nil
  })

  log.Fatal(app.Listen(fmt.Sprintf(":%s", os.Getenv("PORT"))))
}
