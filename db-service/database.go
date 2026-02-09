package main

import (
	"database/sql"
	"fmt"
	"time"
)

type Database struct {
	conn *sql.DB
}

func NewDatabase(dsn string) (*Database, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &Database{conn: db}, nil
}

func (d *Database) Close() error {
	return d.conn.Close()
}

func (d *Database) GetUserByCredentials(username, password string) (int, error) {
	var userID int
	err := d.conn.QueryRow(
		"SELECT id FROM users WHERE username = ? AND password = ?",
		username, password,
	).Scan(&userID)

	if err != nil {
		return 0, fmt.Errorf("invalid credentials")
	}
	return userID, nil
}

func (d *Database) SaveToken(token string, userID int, expiresAt time.Time) error {
	_, err := d.conn.Exec(
		"INSERT INTO tokens (token, user_id, expires_at) VALUES (?, ?, ?)",
		token, userID, expiresAt,
	)
	return err
}

func (d *Database) ValidateToken(token string) (int, error) {
	var userID int
	var expiresAtStr string

	err := d.conn.QueryRow(
		"SELECT user_id, expires_at FROM tokens WHERE token = ?",
		token,
	).Scan(&userID, &expiresAtStr)

	if err != nil {
		return 0, fmt.Errorf("invalid token")
	}

	expiresAt, err := time.Parse("2006-01-02 15:04:05", expiresAtStr)
	if err != nil {
		return 0, fmt.Errorf("could not parse expires_at")
	}

	if time.Now().After(expiresAt) {
		return 0, fmt.Errorf("token expired")
	}

	return userID, nil
}

func (d *Database) GetBalance(userID int) (float64, error) {
	var balance float64
	err := d.conn.QueryRow(
		"SELECT balance FROM accounts WHERE user_id = ?",
		userID,
	).Scan(&balance)

	return balance, err
}

func (d *Database) UpdateBalance(userID int, amount float64) error {
	_, err := d.conn.Exec(
		"UPDATE accounts SET balance = balance + ? WHERE user_id = ?",
		amount, userID,
	)
	return err
}

func (d *Database) SaveTransaction(userID int, transactionID, txType string, amount, balanceAfter float64, timestamp int64) error {
	_, err := d.conn.Exec(
		"INSERT INTO transactions (user_id, transaction_id, type, amount, balance_after, timestamp) VALUES (?, ?, ?, ?, ?, ?)",
		userID, transactionID, txType, amount, balanceAfter, timestamp,
	)
	return err
}
