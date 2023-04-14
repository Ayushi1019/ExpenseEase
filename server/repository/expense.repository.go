package repositories

// package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	models "ExpenseEase/server/model"

	_ "github.com/lib/pq"
)

// ExpenseRepository represents the repository for the users table
type ExpenseRepository struct {
	DB *sql.DB
}

func NewExpenseRepository(db *sql.DB) *ExpenseRepository {
	return &ExpenseRepository{DB: db}
}

func (repo *ExpenseRepository) CreateExpense(expense *models.Expense) (map[string]interface{}, error) {
	statement, err := repo.DB.Prepare("INSERT INTO expenses(amount, category, created_at,user_id) VALUES ($1, $2, $3, $4) RETURNING id, amount, category, created_at")
	if err != nil {
		log.Println("Error preparing SQL statement:", err)
		return nil, errors.New("could not create expense")
	}
	defer statement.Close()

	var id int64
	var newAmount float64
	var newCategory string
	var newDate string

	err = statement.QueryRow(expense.Amount, expense.Category, expense.Created_at, expense.UserID).Scan(&id, &newAmount, &newCategory, &newDate)
	if err != nil {
		log.Println("Error executing SQL statement:", err)
		return nil, errors.New("could not create user")
	}

	fmt.Printf("New expense added with id %d, amount %.2f, category %s, date %s\n", id, newAmount, newCategory, newDate)

	return map[string]interface{}{
		"id":         id,
		"amount":     newAmount,
		"category":   newCategory,
		"created_at": newDate,
	}, nil
}

func (repo *ExpenseRepository) GetAllExpenses(user_id int) ([]models.Expense, error) {
	rows, err := repo.DB.Query("SELECT * FROM expenses WHERE user_id=$1", user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expenses []models.Expense
	for rows.Next() {
		expense := models.Expense{}
		err := rows.Scan(&expense.ID, &expense.Amount, &expense.Category, &expense.Created_at, &expense.UserID)
		if err != nil {
			return nil, err
		}
		expenses = append(expenses, expense)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return expenses, nil
}

func (repo *ExpenseRepository) UpdateExpense(expenseID int, expense *models.Expense) (map[string]interface{}, error) {
	query := `UPDATE expenses SET amount = $1, category = $2, created_at = $3 WHERE id = $4`
	_, err := repo.DB.Exec(query, expense.Amount, expense.Category, expense.Created_at, expenseID)
	if err != nil {
		return nil, err
	}

	// Query for the updated record
	query = `SELECT id, amount, category, created_at FROM expenses WHERE id = $1`
	row := repo.DB.QueryRow(query, expense.ID)
	var updatedexpense models.Expense
	err = row.Scan(&updatedexpense.ID, &updatedexpense.Amount, &updatedexpense.Category, &updatedexpense.Created_at)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"id":         &updatedexpense.ID,
		"amount":     &updatedexpense.Amount,
		"category":   &updatedexpense.Category,
		"created_at": &updatedexpense.Created_at,
	}, nil
}

func (repo *ExpenseRepository) DeleteExpense(expenseID int) error {
	query := `DELETE FROM expenses WHERE id = $1`
	_, err := repo.DB.Exec(query, expenseID)
	if err != nil {
		return err
	}

	return nil
}
