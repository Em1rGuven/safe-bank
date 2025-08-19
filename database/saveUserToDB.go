package database

import (
	"database/sql"
	"time"
)

func SaveUserToDB(username, passwordHash string) (int, error) {
	tx, err := DB.Begin()
	if err != nil {
		return -1, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	createdAt := time.Now().Format("2006-01-02 15:04:05")
	var id int
	err = tx.QueryRow(`
		INSERT INTO users (name, passwordHash, createdAt) VALUES (?, ?, ?)
		RETURNING id
	`, username, passwordHash, createdAt).Scan(&id)
	if err != nil {
		return -1, err
	}

	if err = CreateAccountLine(id, tx); err != nil {
		return -1, err
	}

	if err = tx.Commit(); err != nil {
		return -1, err
	}

	return id, nil
}

func CreateAccountLine(id int, tx *sql.Tx) error {
	_, err := tx.Exec("INSERT INTO accounts (userId, balance, creditDepth) VALUES (?, ?, ?)", id, 0, 0)
	return err
}
