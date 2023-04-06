package models

import (
	"database/sql"
	"fmt"
)

type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

func CreateTables(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			email TEXT NOT NULL,
			password TEXT NOT NULL,
			name TEXT NOT NULL
		);
	`)
	if err != nil {
		return fmt.Errorf("Failed to create users table: %v", err)
	}

	return nil
}
