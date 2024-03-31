package handlers

import (
	"github.com/Seifbarouni/private-git/web-app/back/data"

	"github.com/gofiber/fiber/v2"
)

var AccessService data.AccessServiceInterface = data.InitAccessService()

func GrantAccess(c *fiber.Ctx) error {
	userID, err := getUserIdFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "error getting user id",
		})
	}
	access := new(data.Access)

	if err := c.BodyParser(access); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "error parsing request",
		})
	}

	_, err = RepoService.GetRepo(access.RepoId.Hex(), userID.Hex())

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "error getting repo",
		})
	}

	err = AccessService.GrantAccess(access)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "error granting access",
		})
	}

	return nil
}

func GetAccesses(c *fiber.Ctx) error {
	userID, err := getUserIdFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "error getting user id",
		})
	}

	access, err := AccessService.GetAccessesByUserId(userID.Hex())

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "error getting accesses",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"accesses": access,
	})
}

func RevokeAccess(c *fiber.Ctx) error {
	userID, err := getUserIdFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "error getting user id",
		})
	}

	repoID := c.Params("repo_id")

	if repoID == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "error getting repo id",
		})
	}

	rev_user_id := c.Params("user_id")

	if rev_user_id == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "error getting user id",
		})
	}

	repo, err := RepoService.GetRepo(repoID, userID.Hex())

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "error getting repo",
		})
	}

	err = AccessService.RevokeAccess(rev_user_id, repo)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "error revoking access",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "access revoked",
	})
}
