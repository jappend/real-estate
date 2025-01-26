package database

var createUUIDExtensionQuery = `
  CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
`

var usersTableQuery = `
  CREATE TABLE IF NOT EXISTS users(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR NOT NULL,
    email VARCHAR UNIQUE NOT NULL,
    is_adm BOOLEAN NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    password VARCHAR NOT NULL,
    created_at DATE NOT NULL,
    updated_at DATE NOT NULL
  );	  
  `

var tableQueries = []string{createUUIDExtensionQuery, usersTableQuery}
