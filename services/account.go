package services

import (
	"bankSystem/database"
	"bankSystem/logs"
	"bankSystem/types"
	"context"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"strings"
)

func ShowAccount(c *fiber.Ctx) error {
	id := int(c.Locals("userid").(float64))

	balance, depth, err := database.GetUserInfo(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error.",
		})
	}

	var account = &types.Account{
		ID:       id,
		Username: c.Locals("username").(string),
		Balance:  balance,
		Depth:    depth,
	}

	return c.JSON(account)
}

func DepositMoney(c *fiber.Ctx) error {
	amount, _ := strconv.Atoi(c.FormValue("amount"))
	if amount < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid amount",
		})
	}

	id := int(c.Locals("userid").(float64))
	username := c.Locals("username").(string)

	if err := database.Deposit(id, amount); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error.",
		})
	}

	logs.CreateLogObjects(5, username)

	return c.JSON(fiber.Map{
		"success": "deposit completed successfully.",
	})
}

func WithdrawMoney(c *fiber.Ctx) error {
	amount, _ := strconv.Atoi(c.FormValue("amount"))
	idText := strconv.Itoa(int(c.Locals("userid").(float64)))
	userBalance, _ := strconv.Atoi(database.Redis.HGet(context.Background(), "user:"+idText, "balance").Val())
	if amount < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid amount.",
		})
	} else if userBalance < amount {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "insufficient balance.",
		})
	}

	id := int(c.Locals("userid").(float64))
	if err := database.Withdraw(id, amount); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error.",
		})
	}

	username := c.Locals("username").(string)
	logs.CreateLogObjects(6, username)

	return c.JSON(fiber.Map{
		"success": "withdraw completed successfully.",
	})
}

func TransferMoney(c *fiber.Ctx) error {
	id := int(c.Locals("userid").(float64))
	idText := strconv.Itoa(id)
	currentBalance, _ := strconv.Atoi(database.Redis.HGet(context.Background(), "user:"+idText, "balance").Val())
	amount, _ := strconv.Atoi(c.FormValue("amount"))
	otherPerson := c.FormValue("otherPerson")
	otherPerson = strings.TrimSpace(otherPerson)
	if amount < 0 || otherPerson == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid inputs, please try again.",
		})
	} else if amount > currentBalance {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "insufficient balance.",
		})
	}

	if err := database.TransferMoney(id, amount, otherPerson); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error.",
		})
	}

	username := c.Locals("username").(string)
	logs.CreateLogObjects(7, username)

	return c.JSON(fiber.Map{
		"success": "transfer completed successfully.",
	})
}
