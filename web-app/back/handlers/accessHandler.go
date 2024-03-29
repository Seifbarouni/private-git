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
	_ = userID

	return nil

}

func RevokeAccess(c *fiber.Ctx) error {
	userID, err := getUserIdFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "error getting user id",
		})
	}
	_ = userID

	return nil
}
