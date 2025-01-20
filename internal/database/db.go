package database

import "database/sql"

type Queries struct {
  db *sql.DB
}

func New(db *sql.DB) *Queries {
  Initialize(db)

  return &Queries{db: db}
}
