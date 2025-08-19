package database

import (
	"context"
	"database/sql"
)

func GetUserInfo(id int) (int, int, error) {
	balance, depth := 0, 0
	err := DB.QueryRow(`
		select balance, creditDepth from accounts
		where userId = ?
	`, id).Scan(&balance, &depth)
	if err != nil {
		return -1, -1, err
	}
	return balance, depth, nil
}

func GetUserID(ctx context.Context, tx *sql.Tx, username string) (int, error) {
	var id int

	err := tx.QueryRowContext(ctx, "SELECT id FROM users WHERE name = ?", username).Scan(&id)
	if err != nil {
		return -1, err
	}

	return id, nil
}
