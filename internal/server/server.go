package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func New(config fiber.Config) *fiber.App {
	app := fiber.New(config)

	// Middleware
	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))

	// Router Initializer
	routerInitializer(app)

	return app
}
