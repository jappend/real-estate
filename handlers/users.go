package handlers

import (
	_ "crypto"
	"jappend/real_estate/internal/database"
	"log"
	"strconv"
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
		Name     string `json:"name" validate:"required"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
		IsAdm    bool   `json:"is_adm"`
		IsActive bool   `json:"is_active"`
	}

	params := parameters{}
	if err := c.BodyParser(&params); err != nil {
		log.Println(err)
		return &fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: "Unknown error",
		}
	}

	errors := cfg.Validator.ValidateData(params)
	if errors != nil {
		c.Status(fiber.StatusBadRequest).JSON(errors)
		return nil
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
		log.Println(err)
		return &fiber.Error{
			Code:    fiber.StatusInternalServerError,
			Message: "Unknown error",
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
			Message: "Unknown error",
		}
	}

	c.Status(fiber.StatusCreated).JSON(databaseUserToHandleUser(user))
	return nil
}

func databaseUserToHandleUser(user database.User) User {
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

func (cfg *Config) UsersListAllinDB(c *fiber.Ctx) error {
	c.Accepts("application/json")

	offset := c.Query("offset", "0")
	limit := c.Query("limit", "10")

	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		return &fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: "Offset query param has to be a valid integer",
		}
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		return &fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: "Limit query param has to be a valid integer",
		}
	}

	users, err := cfg.DB.ListAllUsersInDB(database.ListAllUsersParams{
		Offset: offsetInt,
		Limit:  limitInt,
	})
	if err != nil {
		log.Println(err)
		return &fiber.Error{
			Code:    fiber.StatusInternalServerError,
			Message: "Error getting users from database",
		}
	}

	c.Status(fiber.StatusOK).JSON(databaseUsersSliceToHandleUserSlice(users))
	return nil
}

func databaseUsersSliceToHandleUserSlice(users []database.User) []User {
	result := make([]User, len(users))
	for i, user := range users {
		result[i] = databaseUserToHandleUser(user)
	}

	return result
}
