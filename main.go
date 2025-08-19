package main

import (
	_ "bankSystem/database"
	"bankSystem/logs"
	"bankSystem/services"
	"github.com/gofiber/fiber/v2"
	"log"
	"os"
)

func main() {
	logs.Logs()
	app := fiber.New()

	app.Use(services.RateLimiter)

	app.Get("/", services.Welcome)

	app.Post("/register", services.Register)
	app.Post("/login", services.Login)

	auth := app.Group("/", services.LoginCheck)
	auth.Get("/logout", services.Logout)
	auth.Get("/account", services.ShowAccount)
	auth.Post("/deposit", services.DepositMoney)
	auth.Post("/withdraw", services.WithdrawMoney)
	auth.Post("/transfer", services.TransferMoney)

	log.Fatal(app.Listen(":" + os.Getenv("PORT")))
}
