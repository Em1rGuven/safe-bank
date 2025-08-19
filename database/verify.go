package database

import (
	"database/sql"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func VerifyUser(username, password string) (int, error) {
	var (
		id           int
		passwordHash string
		errorMessage = fmt.Errorf("invalid username or password")
	)
	const dummyHash = "$2a$10$7EqJtq98hPqEX7fNZaFWoOHi8A7C6iY6xkT7Y8v8v0iE6YfK1Lw5K"
	err := DB.QueryRow("select id, passwordHash from users where name = ? limit 1", username).Scan(&id, &passwordHash)
	if err != nil {
		_ = bcrypt.CompareHashAndPassword([]byte(dummyHash), []byte(password))
		if errors.Is(err, sql.ErrNoRows) {
			return -1, errorMessage
		}
		return -1, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)); err != nil {
		return -1, errorMessage
	}

	return id, nil
}
