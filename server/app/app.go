package app

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
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
	// DROP TABLE if exists incomes;
	// DROP TABLE if exists users;
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		email TEXT NOT NULL,
		password TEXT NOT NULL,
		name TEXT NOT NULL
	);
	CREATE TABLE IF NOT EXISTS incomes (
		id SERIAL PRIMARY KEY,
		amount FLOAT NOT NULL,
		source TEXT NOT NULL,
		created_at TEXT NOT NULL,
		user_id INT,
		FOREIGN KEY (user_id) REFERENCES users(id)
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
	app.Router.HandleFunc("/users", app.getUsersHandler(userRepo)).Methods("GET")
	app.Router.HandleFunc("/login", app.loginHandler(userRepo)).Methods("POST")
	app.Router.HandleFunc("/income", app.createIncome(incomeRepo)).Methods("POST")
	app.Router.HandleFunc("/incomes", app.getAllIncomes(incomeRepo)).Methods("GET")
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	fmt.Println(payload)
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

		respondWithJSON(w, http.StatusCreated, "User Created")
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

		fmt.Println(fetchedUser.ID)

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"id":  fetchedUser.ID,
			"exp": time.Now().Add(time.Hour * 24).Unix(),
		})
		appConfig := config.GetConfig()
		appConfig.JwtSecret = os.Getenv("JWT_SECRET")
		tokenString, err := token.SignedString([]byte(appConfig.JwtSecret))

		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		respondWithJSON(w, http.StatusOK, map[string]string{"token": tokenString})
	}
}

func (app *App) createIncome(incomeRepo *repositories.IncomeRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		income := models.Income{}
		err := json.NewDecoder(r.Body).Decode(&income)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		// Extract the JWT token from the request header
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			respondWithError(w, http.StatusBadRequest, "Missing authorization token")
			return
		}
		appConfig := config.GetConfig()
		fmt.Println(appConfig.JwtSecret)
		// Parse and validate the JWT token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			// Return the secret key used to sign the token
			return []byte(appConfig.JwtSecret), nil
		})
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Invalid authorization token")
			return
		}

		// Extract the user ID from the JWT token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			respondWithError(w, http.StatusUnauthorized, "Invalid authorization token")
			return
		}
		userID := int(claims["id"].(float64))
		fmt.Println(userID)

		// Set the user ID in the income object
		income.UserID = userID

		d, err := time.Parse(time.RFC3339, income.Created_at)
		if err != nil {
			fmt.Println(err)
		}
		income.Created_at = d.Format("2006-01-02")

		// Create the income record in the database
		result, err := incomeRepo.CreateIncome(&income)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to create income record")
			return
		}

		// Return the created income record
		// t := time.Now()
		// income.Created_at = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
		respondWithJSON(w, http.StatusCreated, result)
	}
}

func (app *App) getAllIncomes(incomeRepo *repositories.IncomeRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			respondWithError(w, http.StatusBadRequest, "Missing authorization token")
			return
		}
		appConfig := config.GetConfig()
		fmt.Println(appConfig.JwtSecret)
		// Parse and validate the JWT token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			// Return the secret key used to sign the token
			return []byte(appConfig.JwtSecret), nil
		})

		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			respondWithError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		userID := int(claims["id"].(float64))
		fmt.Println(userID)
		incomes, err := incomeRepo.GetAllIncomes(userID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		fmt.Println(incomes)

		respondWithJSON(w, http.StatusOK, incomes)
	}
}
