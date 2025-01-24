package handlers

import (
	_ "crypto"
	"jappend/real_estate/internal/database"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	IsAdm     bool      `json:"is_adm"`
	IsActive  bool      `json:"is_active"`
}

func (cfg *Config) UsersCreate(c *fiber.Ctx) error {
	c.Accepts("application/json")

	type parameters struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		IsAdm    bool   `json:"is_adm"`
		IsActive bool   `json:"is_active"`
	}

	params := parameters{}
	if err := c.BodyParser(&params); err != nil {
		return &fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: "Error parsing body in request",
		}
	}

	// Checking duplicate emails

	if cfg.DB.CheckDuplicatedEmail(params.Email) {
		return &fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: "Email already in use.",
		}
	}

	// Hashing password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), 12)
	if err != nil {
		return &fiber.Error{
			Code:    fiber.StatusInternalServerError,
			Message: "Couldn't encrypt",
		}
	}

	user, err := cfg.DB.CreateUser(database.CreateUserParam{
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Email:     params.Email,
		Password:  string(hashedPassword),
		IsAdm:     params.IsAdm,
		IsActive:  params.IsActive,
	})

	if err != nil {
		log.Println(err)
		return &fiber.Error{
			Code:    fiber.StatusInternalServerError,
			Message: "Unknown error, something went wrong",
		}
	}

	c.Status(fiber.StatusCreated).JSON(handleUserReturnOnCreation(user))
	return nil
}

func handleUserReturnOnCreation(user database.User) User {
	return User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Name:      user.Name,
		Email:     user.Email,
		IsAdm:     user.IsAdm,
		IsActive:  user.IsActive,
	}
}
