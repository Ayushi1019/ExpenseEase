package repositories

import (
	"database/sql"
	"errors"
	"fmt"
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

func (repo *IncomeRepository) CreateIncome(income *models.Income) (map[string]interface{}, error) {
	statement, err := repo.DB.Prepare("INSERT INTO incomes(amount, source, created_at,user_id) VALUES ($1, $2, $3, $4) RETURNING id, amount, source, created_at")
	if err != nil {
		log.Println("Error preparing SQL statement:", err)
		return nil, errors.New("could not create income")
	}
	defer statement.Close()

	var id int64
	var newAmount float64
	var newSource string
	var newDate string

	err = statement.QueryRow(income.Amount, income.Source, income.Created_at, income.UserID).Scan(&id, &newAmount, &newSource, &newDate)
	if err != nil {
		log.Println("Error executing SQL statement:", err)
		return nil, errors.New("could not create user")
	}

	fmt.Printf("New income added with id %d, amount %.2f, source %s, date %s\n", id, newAmount, newSource, newDate)

	return map[string]interface{}{
		"id":         id,
		"amount":     newAmount,
		"source":     newSource,
		"created_at": newDate,
	}, nil
}

func (repo *IncomeRepository) GetAllIncomes(user_id int) ([]models.Income, error) {
	rows, err := repo.DB.Query("SELECT * FROM incomes WHERE user_id=$1", user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var incomes []models.Income
	for rows.Next() {
		income := models.Income{}
		err := rows.Scan(&income.ID, &income.Amount, &income.Source, &income.Created_at, &income.UserID)
		if err != nil {
			return nil, err
		}
		incomes = append(incomes, income)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return incomes, nil
}

func (repo *IncomeRepository) UpdateIncome(incomeID int, income *models.Income) (*models.Income, error) {
	query := `UPDATE incomes SET amount = $1, source = $2, created_at = $3 WHERE id = $4`
	_, err := repo.DB.Exec(query, income.Amount, income.Source, income.Created_at, incomeID)
	if err != nil {
		return nil, err
	}

	// Query for the updated record
	query = `SELECT id, amount, source, created_at FROM incomes WHERE id = $1`
	row := repo.DB.QueryRow(query, income.ID)
	var updatedIncome models.Income
	err = row.Scan(&updatedIncome.ID, &updatedIncome.Amount, &updatedIncome.Source, &updatedIncome.Created_at)
	if err != nil {
		return nil, err
	}

	return &updatedIncome, nil
}
