package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	models "ExpenseEase/server/model"
	"time"

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
	statement, err := repo.DB.Prepare("INSERT INTO incomes(user_id, amount, source, created_at) VALUES ($1, $2, $3, $4)")
	if err != nil {
		log.Println("Error preparing SQL statement:", err)
		return errors.New("could not create income")
	}
	defer statement.Close()

	t := time.Now()
	d := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	fmt.Println(d)
	_, err = statement.Exec(income.UserID, income.Amount, income.Source, d.Format("2023-01-01"))
	if err != nil {
		log.Println("Error executing SQL statement:", err)
		return errors.New("could not create user")
	}

	return nil
}
