package handlers

import (
	"fmt"

	"github.com/Seifbarouni/private-git/web-app/back/data"

	"github.com/gofiber/fiber/v2"
)

var accessService data.AccessServiceInterface = data.InitAccessService()

func GrantAccess(c *fiber.Ctx) error {
	userID, err := getUserIdFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": data.APIError{
				Message: fmt.Sprintf("error getting user id: %s", err.Error()),
				Status:  fiber.StatusInternalServerError,
			},
		})
	}
	access := new(data.Access)

	if err := c.BodyParser(access); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": data.APIError{
				Message: fmt.Sprintf("error parsing request body: %s", err.Error()),
				Status:  fiber.StatusInternalServerError,
			},
		})
	}

	_, err = RepoService.GetRepo(access.RepoId.Hex(), userID.Hex())

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": data.APIError{
				Message: err.Error(),
				Status:  fiber.StatusInternalServerError,
			},
		})
	}

	err = accessService.GrantAccess(access)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": data.APIError{
				Message: fmt.Sprintf("error granting access: %s", err.Error()),
				Status:  fiber.StatusInternalServerError,
			},
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "access granted",
	})
}

func GetAccesses(c *fiber.Ctx) error {
	userID, err := getUserIdFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": data.APIError{
				Message: fmt.Sprintf("error getting user id: %s", err.Error()),
				Status:  fiber.StatusInternalServerError,
			},
		})
	}

	access, err := accessService.GetAccessesByUserId(userID.Hex())

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": data.APIError{
				Message: "error getting accesses",
				Status:  fiber.StatusInternalServerError,
			},
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
			"error": data.APIError{
				Message: fmt.Sprintf("error getting user id: %s", err.Error()),
				Status:  fiber.StatusInternalServerError,
			},
		})
	}

	repoID := c.Params("repo_id")

	if repoID == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": data.APIError{
				Message: "error getting repo id ",
				Status:  fiber.StatusInternalServerError,
			},
		})
	}

	rev_user_id := c.Params("user_id")

	if rev_user_id == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": data.APIError{
				Message: "error getting user id",
				Status:  fiber.StatusInternalServerError,
			},
		})
	}

	repo, err := RepoService.GetRepo(repoID, userID.Hex())

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": data.APIError{
				Message: fmt.Sprintf("error getting repo: %s", err.Error()),
				Status:  fiber.StatusInternalServerError,
			},
		})
	}

	err = accessService.RevokeAccess(rev_user_id, repo)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": data.APIError{
				Message: fmt.Sprintf("error revoking access: %s", err.Error()),
				Status:  fiber.StatusInternalServerError,
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "access revoked",
	})
}
