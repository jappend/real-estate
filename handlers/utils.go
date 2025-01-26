package handlers

import (
	"jappend/real_estate/internal/database"
	"jappend/real_estate/internal/validation"
)

type Config struct {
	DB        *database.Queries
	Validator *validation.Validator
}
