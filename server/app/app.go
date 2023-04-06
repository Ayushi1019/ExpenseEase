package app

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"ExpenseEase/server/config"
	models "ExpenseEase/server/model"
	repositories "ExpenseEase/server/repository"

	_ "github.com/lib/pq"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
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

func (app *App) Initialize() {
	var err error
	app.DB, err = config.ConnectDB()
	if err != nil {
		log.Fatal("failed to connect database: ", err)
	}
	createTable()
	app.Router = mux.NewRouter()
	userRepo := repositories.UserRepository{DB: app.DB}
	app.initializeRoutes(&userRepo)
}

func (app *App) Run(addr string) {
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"Authorization", "Content-Type"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	})

	handler := c.Handler(app.Router)

	fmt.Println("Listening on port 8080...")
	log.Fatal(http.ListenAndServe(addr, handler))
}

func (app *App) initializeRoutes(userRepo *repositories.UserRepository) {
	app.Router.HandleFunc("/user", app.createUserHandler(userRepo)).Methods("POST")
	app.Router.HandleFunc("/users", app.getUsersHandler(userRepo)).Methods("Get")
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "JSON encoding error")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func (app *App) createUserHandler(userRepo *repositories.UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		err = userRepo.CreateUser(&user)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		respondWithJSON(w, http.StatusCreated, user)
	}
}

func (app *App) getUsersHandler(userRepo *repositories.UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := userRepo.GetUsers()
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondWithJSON(w, http.StatusOK, users)
	}
}
