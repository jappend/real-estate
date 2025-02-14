-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users(
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  name VARCHAR NOT NULL,
  email VARCHAR UNIQUE NOT NULL,
  is_adm BOOLEAN NOT NULL,
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  password VARCHAR NOT NULL,
  created_at DATE NOT NULL DEFAULT NOW(),
  updated_at DATE NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE users;

DROP EXTENSION "uuid-ossp";
