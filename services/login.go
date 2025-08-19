package services

import (
	"bankSystem/database"
	"bankSystem/logs"
	"context"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"strings"
)

func Login(c *fiber.Ctx) error {
	username := c.FormValue("username")
	username = strings.TrimSpace(username)
	password := c.FormValue("password")
	if username == "" || password == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "username or password is empty, please try again.",
		})
	}

	var id int
	var err error
	if id, err = database.VerifyUser(username, password); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid username or password, please try again.",
		})
	}

	if err = CreateToken(c, username, id); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "internal server error.",
		})
	}

	userBalance, userDepth, err := database.GetUserInfo(id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "internal server error.",
		})
	}

	userId := strconv.Itoa(id)
	err = database.Redis.HSet(context.Background(), "user:"+userId, map[string]interface{}{ // Redis
		"id":          id,
		"name":        username,
		"balance":     userBalance,
		"creditDepth": userDepth,
	}).Err()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "internal server error.",
		})
	}

	logs.CreateLogObjects(3, username)

	return c.JSON(fiber.Map{
		"success": "you are logged in successfully.",
	})
}
