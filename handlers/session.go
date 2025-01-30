package handlers

import (
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func createToken(id uuid.UUID, isAdm *bool) (string, error) {
  token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
    "uuid": id,
    "isAdm": isAdm,
    "exp": time.Now().Add(time.Hour * 24).Unix(),
  })

  tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
  if err != nil {
    return "", err
  }

  return tokenString, nil
}

func (cfg *Config) LoginHandler(c *fiber.Ctx) error {
	c.Accepts("application/json")

	type parameters struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	params := parameters{}
	if err := c.BodyParser(&params); err != nil {
		log.Println(err)
		return &fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: "Error parsing body in request",
		}
	}

	errors := cfg.Validator.ValidateData(params)
	if errors != nil {
		c.Status(fiber.StatusBadRequest).JSON(errors)
		return nil
	}

	// Finding user in DB
	user := cfg.DB.ReturnUserByEmail(params.Email)

	if user.Email == "" {
		return &fiber.Error{
			Code:    fiber.StatusNotFound,
			Message: "User does not exist in the DB!",
		}
	}

  if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Password)); err != nil {
    return &fiber.Error{
      Code: fiber.StatusBadRequest,
      Message: "Wrong password!",
    }
  }

  // Generating token
  token, err := createToken(user.ID, &user.IsAdm)
  if err != nil {
    return &fiber.Error{
      Code: fiber.StatusInternalServerError,
      Message: "Couldn't generate JWT Token",
    }
  }

  tokenReturn := struct {
    Token string `json:"token"`
  }{
    Token: token,
  }

  c.Status(fiber.StatusOK).JSON(tokenReturn)
	return nil
}
