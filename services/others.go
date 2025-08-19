package services

import (
	"bankSystem/database"
	"bankSystem/logs"
	"context"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"time"
)

func Welcome(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Welcome to the Training Bank!",
	})
}

func Logout(c *fiber.Ctx) error {
	claims, err := ParseToken(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	userid := int(claims["id"].(float64))
	userText := strconv.Itoa(userid)

	if err = database.Redis.Del(context.Background(), "user:"+userText).Err(); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	logoutCookie := new(fiber.Cookie)
	logoutCookie.Name = "jwt"
	logoutCookie.Value = ""
	logoutCookie.Path = "/"
	logoutCookie.Secure = false
	logoutCookie.HTTPOnly = true
	logoutCookie.Expires = time.Now().Add(-1 * time.Hour)
	logoutCookie.MaxAge = -1
	c.Cookie(logoutCookie)

	username := c.Locals("username").(string)
	logs.CreateLogObjects(4, username)

	return c.JSON(fiber.Map{
		"message": "you are logged out successfully.",
	})
}
