package database

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
  ID uuid.UUID
  CreatedAt time.Time
  UpdatedAt time.Time
  Name string
  Email string
  Password string
  IsAdm bool
  IsActive bool
}

var usersTableQuery = `
  CREATE TABLE IF NOT EXISTS users(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR NOT NULL,
    email VARCHAR UNIQUE NOT NULL,
    is_adm BOOLEAN NOT NULL,
    Is_active BOOLEAN NOT NULL,
    password VARCHAR NOT NULL,
    created_at DATE NOT NULL,
    updated_at DATE NOT NULL
  );	  
  `

var tableQueries = []string{usersTableQuery}
