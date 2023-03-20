package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"ExpenseEase/server/config"
)

// User is a model that represents a user in the system
type User struct {
	Id       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"Name"`
}

var users []User

func getUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db, err := config.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, email, password, name FROM users")
	if err != nil {
		log.Fatalf("Failed to fetch users: %v", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Id, &user.Email, &user.Password, &user.Name); err != nil {
			log.Fatalf("Failed to scan user row: %v", err)
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		log.Fatalf("Failed to iterate over users rows: %v", err)
	}

	json.NewEncoder(w).Encode(users)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	db, err := config.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	var user User
	err = db.QueryRow("SELECT id, email, password, name FROM users WHERE id=$1", params["id"]).Scan(&user.Id, &user.Email, &user.Password, &user.Name)
	if err != nil {
		log.Fatalf("Failed to fetch user: %v", err)
	}

	json.NewEncoder(w).Encode(user)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user User
	_ = json.NewDecoder(r.Body).Decode(&user)
	user.Id = strconv.Itoa(rand.Intn(1000000))
	db, err := config.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO users(id, email, password, name) VALUES($1, $2, $3, $4)", user.Id, user.Email, user.Password, user.Name)
	if err != nil {
		log.Fatalf("Failed to create user: %v", err)
	}

	json.NewEncoder(w).Encode(user)
}

func createTable() error {
	db, err := config.ConnectDB()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			email TEXT NOT NULL,
			password TEXT NOT NULL,
			name TEXT NOT NULL
		);
	`)
	if err != nil {
		return err
	}
	return nil
}


func main() {
	createTable()
	r := mux.NewRouter()
	r.HandleFunc("/users", getUsers).Methods("GET")
	r.HandleFunc("/users/{id}", getUser).Methods("GET")
	r.HandleFunc("/users", createUser).Methods("POST")

	fmt.Printf("Starting server")
	log.Fatal(http.ListenAndServe(":8000", r))
}