package server

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

func dbInitializer() *sql.DB {
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

	if err := db.Ping(); err != nil {
		log.Fatal("Error pinging database: ", err)
	}
	log.Println("Database connection opened successefuly!")

	return db
}
