package database

import (
	"database/sql"
	"log"
)

func Initialize(db *sql.DB) {
  for _, query := range tableQueries {
    if _, err := db.Exec(query); err != nil {
      log.Fatal("Error creating tables: ", err)
    }
  } 
} 
