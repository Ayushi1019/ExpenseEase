package app

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"ExpenseEase/server/config"
	models "ExpenseEase/server/model"
	repositories "ExpenseEase/server/repository"

	_ "github.com/lib/pq"
)

type App struct {
	Router           *mux.Router
	DB               *sql.DB
	Config           *config.Config
	UserRepository   *repositories.UserRepository
	IncomeRepository *repositories.IncomeRepository
}

func createTable() error {
	db, err := config.ConnectDB()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		email TEXT NOT NULL,
		password TEXT NOT NULL,
		name TEXT NOT NULL
	);
	CREATE TABLE IF NOT EXISTS income (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		amount FLOAT NOT NULL,
		userID INT,
		FOREIGN KEY (userID) REFERENCES users(id)
	);
	`)
	fmt.Println("createTable runs", err)
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

	app.Router = mux.NewRouter()
	userRepo := repositories.UserRepository{DB: app.DB}
	incomeRepo := repositories.IncomeRepository{DB: app.DB}
	app.initializeRoutes(&userRepo, &incomeRepo)
}

func (app *App) Run(addr string) {

	createTable()
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"Authorization", "Content-Type"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	})

	handler := c.Handler(app.Router)

	fmt.Println("Listening on port 8080...")
	log.Fatal(http.ListenAndServe(addr, handler))
}

func (app *App) initializeRoutes(userRepo *repositories.UserRepository, incomeRepo *repositories.IncomeRepository) {
	app.Router.HandleFunc("/user", app.createUserHandler(userRepo)).Methods("POST")
	app.Router.HandleFunc("/users", app.getUsersHandler(userRepo)).Methods("Get")
	app.Router.HandleFunc("/login", app.loginHandler(userRepo)).Methods("POST")
	// app.Router.HandleFunc("/income", app.createIncome(incomeRepo)).Methods("POST")
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

func (app *App) loginHandler(userRepo *repositories.UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		fetchedUser, err := userRepo.GetUserByEmail(user.Email)

		key := make([]byte, 32)
		_, err = rand.Read(key)
		if err != nil {
			fmt.Println(err)
		}

		if err != nil {
			if err == sql.ErrNoRows {
				respondWithError(w, http.StatusUnauthorized, "Invalid email or password")
			} else {
				respondWithError(w, http.StatusInternalServerError, err.Error())
			}
			return
		}

		if user.Password != fetchedUser.Password {
			respondWithError(w, http.StatusUnauthorized, "Invalid email or password")
			return
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"email": fetchedUser.Email,
			"exp":   time.Now().Add(time.Hour * 24).Unix(),
		})
		t := base64.StdEncoding.EncodeToString(key)
		appConfig := config.GetConfig()
		appConfig.JwtSecret = t
		tokenString, err := token.SignedString([]byte(t))

		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		respondWithJSON(w, http.StatusOK, map[string]string{"token": tokenString})
	}
}
