package handlers

import (
	"os"
	"time"

	"github.com/Seifbarouni/private-git/web-app/back/data"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var UserService data.UserServiceInterface = data.InitUserService()

func Login(c *fiber.Ctx) error {
	var user data.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": data.APIError{
				Message: err.Error(),
				Status:  fiber.StatusBadRequest,
			},
		})
	}
	userCheck, err := UserService.GetUserByEmail(user.Email)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": data.APIError{
				Message: "user not found",
				Status:  fiber.StatusNotFound,
			},
		})
	}

	if !checkPasswordHash(user.Password, userCheck.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": data.APIError{
				Message: "invalid password",
				Status:  fiber.StatusUnauthorized,
			},
		})
	}

	claims := jwt.MapClaims{
		"name":  userCheck.UserName,
		"email": userCheck.Email,
		"id":    userCheck.ID.Hex(),
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "login successful",
		"token":   t,
	})
}

func Register(c *fiber.Ctx) error {
	var user data.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": data.APIError{
				Message: err.Error(),
				Status:  fiber.StatusBadRequest,
			},
		})
	}

	hash, err := hashPassword(user.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": data.APIError{
				Message: err.Error(),
				Status:  fiber.StatusInternalServerError,
			},
		})
	}
	user.Password = hash
	user.Status = "active"

	if err := UserService.CreateUser(&user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": data.APIError{
				Message: err.Error(),
				Status:  fiber.StatusInternalServerError,
			},
		})
	}
	claims := jwt.MapClaims{
		"name":  user.UserName,
		"email": user.Email,
		"id":    user.ID.Hex(),
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "user created",
		"token":   t,
	})
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 5)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
