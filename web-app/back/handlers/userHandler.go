package handlers

import (
	"fmt"
	"os"
	"time"

	"github.com/Seifbarouni/private-git/web-app/back/data"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var userService data.UserServiceInterface = data.InitUserService()

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

	userCheck, err := userService.GetUserByEmail(user.Email)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": data.APIError{
				Message: fmt.Sprintf("user not found: %s", err.Error()),
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

	validate := validator.New()
	if err := validate.Struct(user); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": data.APIError{
				Message: err.Error(),
				Status:  fiber.StatusNotFound,
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

	if err := userService.CreateUser(&user); err != nil {
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

func AddSSHKey(c *fiber.Ctx) error {
	userId, err := getUserIdFromToken(c)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": data.APIError{
				Message: fmt.Sprintf("error getting user id: %s", err.Error()),
				Status:  fiber.StatusInternalServerError,
			},
		})
	}

	var sshk data.SSHKey

	if err := c.BodyParser(&sshk); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": data.APIError{
				Message: err.Error(),
				Status:  fiber.StatusBadRequest,
			},
		})
	}

	if err = userService.AddPublicKey(userId, sshk.Key); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": data.APIError{
				Message: err.Error(),
				Status:  fiber.StatusBadRequest,
			},
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "added ssh public key",
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
