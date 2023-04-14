package repositories

import (
	"database/sql"

	_ "github.com/lib/pq"
)

// BudgetRepository represents the repository for the users table
type BudgetRepository struct {
	DB *sql.DB
}

func NewBudgetRepository(db *sql.DB) *BudgetRepository {
	return &BudgetRepository{DB: db}
}
