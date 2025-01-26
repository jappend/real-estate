package server

import (
	"jappend/real_estate/handlers"
	"jappend/real_estate/internal/database"
	"jappend/real_estate/internal/validation"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func routerInitializer(app *fiber.App) {
	// Initializing environment
	envInitializer()

	// Initializing our database
	db := dbInitializer()

	// Initializing the validator
	validate := validator.New(validator.WithRequiredStructEnabled())

	validator := validation.New(validate)
	dbQueries := database.New(db)
	handlersConfig := handlers.Config{
		DB:        dbQueries,
		Validator: validator,
	}

	// Routes
	// Users
	app.Post("/users", handlersConfig.UsersCreate)
	app.Get("/users", handlersConfig.UsersListAllinDB)
}
