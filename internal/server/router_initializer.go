package server

import (
	"jappend/real_estate/handlers"
	"jappend/real_estate/internal/database"

	"github.com/gofiber/fiber/v2"
)

func routerInitializer(app *fiber.App) {
	// Initializing environment
	envInitializer()

	// Initializing our database
	db := dbInitializer()

	dbQueries := database.New(db)
	handlersConfig := handlers.Config{
		DB: dbQueries,
	}

	// Routes
	// Users
	app.Post("/users", handlersConfig.UsersCreate)
	app.Get("/users", handlersConfig.UsersListAllinDB)
}
