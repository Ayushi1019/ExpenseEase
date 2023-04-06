package repositories

import (
	"database/sql"
	"errors"
	"log"

	models "ExpenseEase/server/model"

	_ "github.com/lib/pq"
)

// UserRepository represents the repository for the users table
type UserRepository struct {
	DB *sql.DB
}

// NewUserRepository returns a new instance of UserRepository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// CreateUser creates a new user in the database
func (repo *UserRepository) CreateUser(user *models.User) error {
	statement, err := repo.DB.Prepare("INSERT INTO users(id, name, email, password) VALUES ($1, $2, $3, $4)")
	if err != nil {
		log.Println("Error preparing SQL statement:", err)
		return errors.New("could not create user")
	}
	defer statement.Close()

	_, err = statement.Exec(user.ID, user.Name, user.Email, user.Password)
	if err != nil {
		log.Println("Error executing SQL statement:", err)
		return errors.New("could not create user")
	}

	return nil
}

func (repo *UserRepository) GetUsers() ([]models.User, error) {
	var users []models.User

	rows, err := repo.DB.Query("SELECT id, name, email FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Name, &user.Email)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
