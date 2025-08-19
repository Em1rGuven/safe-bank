package services

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

func ParseToken(c *fiber.Ctx) (map[string]interface{}, error) {
	tokenString := c.Cookies("jwt")
	if tokenString == "" {
		return nil, errors.New("token is empty")
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

func CreateToken(c *fiber.Ctx, username string, id int) error {
	content := jwt.MapClaims{
		"username": username,
		"id":       id,
		"exp":      time.Now().Add(time.Hour * 3).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, content)

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return err
	}

	CreateCookie(c, tokenString)

	return nil
}

func CreateCookie(c *fiber.Ctx, tokenString string) {
	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    tokenString,
		Expires:  time.Now().Add(3 * time.Hour),
		MaxAge:   int((3 * time.Hour).Seconds()),
		HTTPOnly: true,
		Secure:   false, // for edu, normally this MUST be true
		SameSite: "Lax",
		Path:     "/",
	})
}
