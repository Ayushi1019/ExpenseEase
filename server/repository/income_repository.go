package repositories

import (
	"database/sql"
	"errors"
	"log"

	models "ExpenseEase/server/model"

	_ "github.com/lib/pq"
)

// IncomeRepository represents the repository for the users table
type IncomeRepository struct {
	DB *sql.DB
}

func NewIncomeRepository(db *sql.DB) *IncomeRepository {
	return &IncomeRepository{DB: db}
}

func (repo *IncomeRepository) CreateIncome(income *models.Income) error {
	statement, err := repo.DB.Prepare("INSERT INTO incomes (user_id, amount, source, date) VALUES ($1, $2, $3, DATE.NOW())")
	if err != nil {
		log.Println("Error preparing SQL statement:", err)
		return errors.New("could not create income")
	}
	defer statement.Close()

	_, err = statement.Exec(income.UserID, income.Amount, income.Source, income.Date)
	if err != nil {
		log.Println("Error executing SQL statement:", err)
		return errors.New("could not create user")
	}

	return nil
}
