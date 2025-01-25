package main

import (
	"fmt"
	"jappend/real_estate/internal/server"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

func main() {
	config := fiber.Config{
		AppName:      "Real Estate",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 90 * time.Second,
		IdleTimeout:  120 * time.Second,
		ErrorHandler: serverErrorHandler,
	}

	app := server.New(config)

	log.Fatal(app.Listen(fmt.Sprintf(":%s", os.Getenv("PORT"))))
}

func serverErrorHandler(c *fiber.Ctx, err error) error {
	type serverError struct {
		Message string `json:"message"`
	}

	if fe, ok := err.(*fiber.Error); ok {
		return c.Status(fe.Code).JSON(serverError{
			Message: fe.Message,
		})
	}

	return c.Status(fiber.StatusInternalServerError).JSON(serverError{
		Message: err.Error(),
	})
}
