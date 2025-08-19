package services

import (
	"bankSystem/database"
	"bankSystem/logs"
	"context"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"strings"
)

const (
	MinPasswordLength = 6
	MaxPasswordLength = 18
)

func Register(c *fiber.Ctx) error {
	username := c.FormValue("username")
	username = strings.TrimSpace(username)
	password := c.FormValue("password")

	if username == "" || password == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "username or password is empty, please try again.",
		})
	} else if len(password) < MinPasswordLength || len(password) > MaxPasswordLength {
		return c.Status(400).JSON(fiber.Map{
			"error": "password length must be between 6 and 18, please try again.",
		})
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var id int
	if id, err = database.SaveUserToDB(username, string(passwordHash)); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err = CreateToken(c, username, id); err != nil { // JWT
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	userId := strconv.Itoa(id)
	err = database.Redis.HSet(context.Background(), "user:"+userId, map[string]interface{}{ // Redis
		"id":          id,
		"name":        username,
		"balance":     0,
		"creditDepth": 0,
	}).Err()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "internal server error.",
		})
	}

	logs.CreateLogObjects(2, username)

	return c.Status(201).JSON(fiber.Map{
		"message": "you are registered successfully.",
	})
}
