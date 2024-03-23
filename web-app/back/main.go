package main

import (
	"log"
	"os"

	"github.com/Seifbarouni/private-git/web-app/back/db"
	h "github.com/Seifbarouni/private-git/web-app/back/handlers"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

func main() {
	err := db.Init()
	if err != nil {
		log.Fatal(err)
	}

	app := fiber.New(fiber.Config{
		AppName:               "Private Git",
		DisableStartupMessage: true,
	})

	// unrestricted routes
	app.Post("/api/v1/register", h.Register)
	app.Post("/api/v1/login", h.Login)
	// app.Get("/refresh", nil)

	// restricted routes
	v1 := app.Group("/api/v1")
	v1.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(os.Getenv("JWT_SECRET"))},
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusUnauthorized).JSON(map[string]string{"error": "unauthorized"})
		},
	}))

	log.Fatal(app.Listen(":" + os.Getenv("PORT")))
}
