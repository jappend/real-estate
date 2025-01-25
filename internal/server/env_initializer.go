package server

import (
	"log"

	"github.com/joho/godotenv"
)

func envInitializer() {
	// Initializing environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
}
