package main

import (
	"log"
	"os"

	"github.com/Seifbarouni/private-git/web-app/back/data"
	"github.com/Seifbarouni/private-git/web-app/back/db"
	h "github.com/Seifbarouni/private-git/web-app/back/handlers"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading .env file, will us the default environment. If I need a variable and don't find it I will error out")
	}

	err = db.Init()
	if err != nil {
		log.Fatal(err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT must be set")
	}

	jwt_secret := os.Getenv("JWT_SECRET")
	if jwt_secret == "" {
		log.Fatal("JWT_SECRET must be set")
	}

	app := fiber.New(fiber.Config{
		AppName:               "Private Git",
		DisableStartupMessage: true,
	})

	// unrestricted routes
	app.Post("/api/v1/register", h.Register)
	app.Post("/api/v1/login", h.Login)

	// restricted routes
	v1 := app.Group("/api/v1")
	v1.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(jwt_secret)},
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusUnauthorized).JSON(map[string]data.APIError{"error": {
				Message: "unauthorized",
				Status:  fiber.StatusUnauthorized,
			},
			})
		},
	}))
	// user routes
	v1.Post("/user/sshk", h.AddSSHKey)
	// repos routes
	v1.Get("/repos", h.GeRepos)
	v1.Get("/repos/:id", h.GetRepoById)
	v1.Post("/repo", h.CreateRepo)
	v1.Put("/repo", h.UpdateRepo)
	v1.Delete("/repo/:id", h.DeleteRepo)

	// access routes
	v1.Get("/accesses", h.GetAccesses)
	v1.Post("/access", h.GrantAccess)
	v1.Delete("/access/:repo_id/:user_id", h.RevokeAccess)

	log.Fatal(app.Listen(":" + port))
}
