package database

import (
	"bankSystem/logs"
	"context"
	"database/sql"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
)

var (
	DB    *sql.DB
	Redis *redis.Client
)

func init() {
	_, _ = os.OpenFile(".log", os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)

	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		exam, err := os.ReadFile(".env.example")
		if err != nil {
			log.Fatal(err)
		}

		if err := os.WriteFile(".env", exam, 0644); err != nil {
			log.Fatal(err)
		}
	}

	err := godotenv.Load(".env")
	if err != nil {
		logs.CreateLogObjects(1, "")
		log.Fatal("Error loading .env file")
	}

	Redis = redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_URL"),
	})
	if err := Redis.Ping(context.Background()).Err(); err != nil {
		logs.CreateLogObjects(1, "")
		log.Fatal(err)
	}

	DB, err = sql.Open("sqlite3", os.Getenv("DB_PATH"))
	if err != nil {
		logs.CreateLogObjects(1, "")
		log.Fatal(err)
	}

	if err = DB.Ping(); err != nil {
		logs.CreateLogObjects(1, "")
		log.Fatal(err)
	}

	_, _ = DB.Exec("pragma journal_mode = wal")
	_, _ = DB.Exec("pragma foreign_keys = on")
	_, _ = DB.Exec("pragma busy_timeout = 5000")
	_, _ = DB.Exec("pragma synchronous = normal;")

	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE,
			passwordHash TEXT NOT NULL,
			createdAt TEXT NOT NULL
		);
		CREATE TABLE IF NOT EXISTS accounts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			userId INTEGER NOT NULL UNIQUE,
			balance INTEGER NOT NULL,
			creditDepth INTEGER NOT NULL,
			FOREIGN KEY(userId) REFERENCES users(id)
		);
	`)
	if err != nil {
		logs.CreateLogObjects(1, "")
		log.Fatal(err)
	}

	logs.CreateLogObjects(0, "")
}
