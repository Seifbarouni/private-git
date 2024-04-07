package handlers

import (
	"fmt"

	"github.com/Seifbarouni/private-git/web-app/back/data"
	"github.com/go-playground/validator/v10"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var RepoService data.RepoServiceInterface = data.InitRepoService()

func GeRepos(c *fiber.Ctx) error {
	userId, err := getUserIdFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": data.APIError{
				Message: fmt.Sprintf("error getting user id: %s", err.Error()),
				Status:  fiber.StatusInternalServerError,
			},
		})
	}
	repos, err := RepoService.GetReposByOwner(userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": data.APIError{
				Message: fmt.Sprintf("error getting repos: %s", err.Error()),
				Status:  fiber.StatusInternalServerError,
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"repos": repos,
	})
}

func GetRepoById(c *fiber.Ctx) error {
	userId, err := getUserIdFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": data.APIError{
				Message: fmt.Sprintf("error getting user id: %s", err.Error()),
				Status:  fiber.StatusInternalServerError,
			},
		})
	}

	repoId := c.Params("id")
	repoIdHex, err := primitive.ObjectIDFromHex(repoId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": data.APIError{
				Message: "Could not parse repo id",
				Status:  fiber.StatusInternalServerError,
			},
		})
	}
	repo, err := RepoService.GetRepo(repoIdHex, userId.Hex())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": data.APIError{
				Message: err.Error(),
				Status:  fiber.StatusInternalServerError,
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"repo": repo,
	})

}

func CreateRepo(c *fiber.Ctx) error {
	userId, err := getUserIdFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": data.APIError{
				Message: fmt.Sprintf("error getting user id: %s", err.Error()),
				Status:  fiber.StatusInternalServerError,
			},
		})
	}

	repo := new(data.Repo)
	if err := c.BodyParser(repo); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": data.APIError{
				Message: fmt.Sprintf("error parsing request body: %s", err.Error()),
				Status:  fiber.StatusInternalServerError,
			},
		})
	}

	validate := validator.New()
	if err := validate.Struct(repo); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": data.APIError{
				Message: err.Error(),
				Status:  fiber.StatusNotFound,
			},
		})
	}

	repo.Owner = userId
	repo.Status = "active"
	err = RepoService.CreateRepo(repo)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": data.APIError{
				Message: err.Error(),
				Status:  fiber.StatusInternalServerError,
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"repo": repo,
	})
}

func UpdateRepo(c *fiber.Ctx) error {
	userId, err := getUserIdFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": data.APIError{
				Message: fmt.Sprintf("error getting user id: %s", err.Error()),
				Status:  fiber.StatusInternalServerError,
			},
		})
	}
	repo := new(data.Repo)
	if err := c.BodyParser(repo); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": data.APIError{
				Message: fmt.Sprintf("error parsing request body: %s", err.Error()),
				Status:  fiber.StatusInternalServerError,
			},
		})
	}

	validate := validator.New()
	if err := validate.Struct(repo); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": data.APIError{
				Message: err.Error(),
				Status:  fiber.StatusNotFound,
			},
		})
	}

	if repo.Owner.Hex() != userId.Hex() {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": data.APIError{
				Message: "unauthorized",
				Status:  fiber.StatusInternalServerError,
			},
		})
	}

	err = RepoService.UpdateRepo(repo)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": data.APIError{
				Message: fmt.Sprintf("error updating repo: %s", err.Error()),
				Status:  fiber.StatusInternalServerError,
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"repo": repo,
	})
}

func DeleteRepo(c *fiber.Ctx) error {
	userId, err := getUserIdFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": data.APIError{
				Message: fmt.Sprintf("error getting user id: %s", err.Error()),
				Status:  fiber.StatusInternalServerError,
			},
		})
	}

	repoId := c.Params("id")
	repoIdHex, err := primitive.ObjectIDFromHex(repoId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": data.APIError{
				Message: err.Error(),
				Status:  fiber.StatusInternalServerError,
			},
		})
	}
	repo, err := RepoService.GetRepo(repoIdHex, userId.Hex())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": data.APIError{
				Message: err.Error(),
				Status:  fiber.StatusInternalServerError,
			},
		})
	}

	err = RepoService.DeleteRepo(repoId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": data.APIError{
				Message: fmt.Sprintf("error deleting repo: %s", err.Error()),
				Status:  fiber.StatusInternalServerError,
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"repo": repo,
	})
}

func getUserIdFromToken(c *fiber.Ctx) (primitive.ObjectID, error) {
	user := c.Locals("user").(*jwt.Token)

	claims := user.Claims.(jwt.MapClaims)
	userIDStr := claims["id"].(string)
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return userID, nil
}
