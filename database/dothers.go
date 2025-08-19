package database

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"
)

func Deposit(userID, amount int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "update accounts set balance = balance + ? where userId = ?", amount, userID)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		_ = tx.Rollback()
		return err
	}

	redisErr := Redis.HIncrBy(ctx, "user:"+strconv.Itoa(userID), "balance", int64(amount)).Err()
	if redisErr != nil {
		return redisErr
	}
	return nil
}

func Withdraw(userID, amount int) error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}

	var res sql.Result
	res, err = tx.Exec(
		"update accounts set balance = balance - ? where userId = ? and balance >= ?",
		amount, userID, amount)
	if err != nil {
		_ = tx.Rollback()
		return err
	} else if eff, _ := res.RowsAffected(); eff == 0 {
		_ = tx.Rollback()
		return fmt.Errorf("insufficient balance")
	}

	if err = tx.Commit(); err != nil {
		_ = tx.Rollback()
		return err
	}

	err = Redis.HIncrBy(context.Background(), "user:"+strconv.Itoa(userID), "balance", -int64(amount)).Err()
	if err != nil {
		return err
	}

	return nil
}

func TransferMoney(userID, amount int, otherPerson string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	tx, err := DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	otherID, err := GetUserID(ctx, tx, otherPerson)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	if userID == otherID {
		_ = tx.Rollback()
		return fmt.Errorf("you cannot transfer yourself")
	}

	res, err := tx.ExecContext(ctx, `
		update accounts set balance = balance - ? where userId = ? and balance >= ?
   `, amount, userID, amount)
	if err != nil {
		_ = tx.Rollback()
		return err
	} else if eff, _ := res.RowsAffected(); eff == 0 {
		_ = tx.Rollback()
		return fmt.Errorf("insufficient money")
	}

	res, err = tx.ExecContext(ctx, `
		update accounts set balance = balance + ? where userId = ?
	`, amount, otherID)
	if err != nil {
		_ = tx.Rollback()
		return err
	} else if eff, _ := res.RowsAffected(); eff == 0 {
		_ = tx.Rollback()
		return fmt.Errorf("unknown error")
	}

	if err = tx.Commit(); err != nil {
		_ = tx.Rollback()
		return err
	}

	err = Redis.HIncrBy(ctx, "user:"+strconv.Itoa(userID), "balance", -int64(amount)).Err()
	if err != nil {
		return err
	}
	
	return nil
}
