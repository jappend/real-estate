package main

import (
	"database/sql"
	"fmt"
	"jappend/real_estate/handlers"
	"jappend/real_estate/internal/database"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Initializing environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	// Opening a connection to the database
	var (
		user     string = os.Getenv("PGUSER")
		password string = os.Getenv("PGPASSWORD")
		host     string = os.Getenv("PGHOST")
		port     string = os.Getenv("PGPORT")
		dbname   string = os.Getenv("PGDATABASE")
	)
	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbname)

	log.Printf("Opening connection with database %s on port %s...\n", dbname, port)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error opening database connection: ", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Error pinging database: ", err)
	}
	log.Println("Database connection opened successefuly!")

	dbQueries := database.New(db)
	handlersConfig := handlers.Config{
		DB: dbQueries,
	}

	config := fiber.Config{
		AppName:      "Real Estate",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 90 * time.Second,
		IdleTimeout:  120 * time.Second,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
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
		},
	}

	app := fiber.New(config)
	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))

	// Routes
	// Users
	app.Post("/users", handlersConfig.UsersCreate)

	log.Fatal(app.Listen(fmt.Sprintf(":%s", os.Getenv("PORT"))))
}
