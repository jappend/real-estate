package database

import (
	"database/sql"
	"log"
)

type Queries struct {
	db *sql.DB
}

func New(db *sql.DB) *Queries {
	for _, query := range tableQueries {
		if _, err := db.Exec(query); err != nil {
			log.Fatal("Error creating tables: ", err)
		}
	}

	return &Queries{db: db}
}
