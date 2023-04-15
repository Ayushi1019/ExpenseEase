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

// BudgetRepository represents the repository for the users table
type BudgetRepository struct {
	DB *sql.DB
}

func NewBudgetRepository(db *sql.DB) *BudgetRepository {
	return &BudgetRepository{DB: db}
}

func (repo *BudgetRepository) CreateBudget(budget *models.Budget) (map[string]interface{}, error) {
	statement, err := repo.DB.Prepare("INSERT INTO budgets(amount, category, created_at,user_id) VALUES ($1, $2, $3, $4) RETURNING id, amount, category, created_at")
	if err != nil {
		log.Println("Error preparing SQL statement:", err)
		return nil, errors.New("could not create budget")
	}
	defer statement.Close()

	var id int64
	var newAmount float64
	var newCategory string
	var newDate string

	err = statement.QueryRow(budget.Amount, budget.Category, budget.Created_at, budget.UserID).Scan(&id, &newAmount, &newCategory, &newDate)
	if err != nil {
		log.Println("Error executing SQL statement:", err)
		return nil, errors.New("could not create user")
	}

	fmt.Printf("New budget added with id %d, amount %.2f, category %s, date %s\n", id, newAmount, newCategory, newDate)

	return map[string]interface{}{
		"id":         id,
		"amount":     newAmount,
		"category":   newCategory,
		"created_at": newDate,
	}, nil
}

func (repo *BudgetRepository) GetAllbudgets(user_id int) ([]models.Budget, error) {
	rows, err := repo.DB.Query("SELECT * FROM budgets WHERE user_id=$1", user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var budgets []models.Budget
	for rows.Next() {
		budget := models.Budget{}
		err := rows.Scan(&budget.ID, &budget.Amount, &budget.Category, &budget.Created_at, &budget.UserID)
		if err != nil {
			return nil, err
		}
		budgets = append(budgets, budget)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return budgets, nil
}

func (repo *BudgetRepository) UpdateBudget(budgetID int, budget *models.Budget) (map[string]interface{}, error) {
	query := `UPDATE budgets SET amount = $1, category = $2, created_at = $3 WHERE id = $4`
	_, err := repo.DB.Exec(query, budget.Amount, budget.Category, budget.Created_at, budgetID)
	if err != nil {
		return nil, err
	}

	// Query for the updated record
	query = `SELECT id, amount, category, created_at FROM budgets WHERE id = $1`
	row := repo.DB.QueryRow(query, budget.ID)
	var updatedBudget models.Budget
	err = row.Scan(&updatedBudget.ID, &updatedBudget.Amount, &updatedBudget.Category, &updatedBudget.Created_at)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"id":         &updatedBudget.ID,
		"amount":     &updatedBudget.Amount,
		"category":   &updatedBudget.Category,
		"created_at": &updatedBudget.Created_at,
	}, nil
}

func (repo *BudgetRepository) DeleteBudget(budgetID int) error {
	query := `DELETE FROM budgets WHERE id = $1`
	_, err := repo.DB.Exec(query, budgetID)
	if err != nil {
		return err
	}

	return nil
}
