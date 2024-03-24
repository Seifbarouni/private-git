package handlers

import (
	"github.com/Seifbarouni/private-git/web-app/back/data"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var RepoService data.RepoServiceInterface = data.InitRepoService()

func GeRepos(c *fiber.Ctx) error {
	userId, err := getUserIdFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "error getting user id",
		})
	}

	repos, err := RepoService.GetReposByOwner(userId.Hex())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "error getting repos",
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
			"error": "error getting user id",
		})
	}

	repoId := c.Params("id")
	repo, err := RepoService.GetRepo(repoId, userId.Hex())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
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
			"error": "error getting user id",
		})
	}

	repo := new(data.Repo)
	if err := c.BodyParser(repo); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "error parsing request body",
		})
	}

	repo.Owner = userId
	repo.Status = "active"
	err = RepoService.CreateRepo(repo)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "error creating repo",
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
			"error": "error getting user id",
		})
	}
	repo := new(data.Repo)
	if err := c.BodyParser(repo); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "error parsing request body",
		})
	}
	if repo.Owner.Hex() != userId.Hex() {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	err = RepoService.UpdateRepo(repo)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "error updating repo",
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
			"error": "error getting user id",
		})
	}

	repoId := c.Params("id")
	repo, err := RepoService.GetRepo(repoId, userId.Hex())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err = RepoService.DeleteRepo(repoId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "error deleting repo",
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
