package services

import (
	"bankSystem/database"
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

func LoginCheck(c *fiber.Ctx) error {
	tokenString := c.Cookies("jwt")
	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "you can not access this resource.",
		})
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "you can not access this resource.",
		})
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		c.Locals("userid", claims["id"].(float64))
		c.Locals("username", claims["username"].(string))
	}

	return c.Next()
}

func RateLimiter(c *fiber.Ctx) error {
	amount, _ := database.Redis.Incr(context.Background(), "rate:"+c.IP()).Result()
	if amount == 1 {
		database.Redis.Expire(context.Background(), "rate:"+c.IP(), time.Second*60)
	}
	if amount > 120 {
		return c.Status(429).JSON(fiber.Map{
			"error": "too many requests, please wait.",
		})
	}
	return c.Next()
}
