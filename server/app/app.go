package app

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
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
	Router            *mux.Router
	DB                *sql.DB
	Config            *config.Config
	UserRepository    *repositories.UserRepository
	IncomeRepository  *repositories.IncomeRepository
	ExpenseRepository *repositories.ExpenseRepository
	BudgetRepository  *repositories.BudgetRepository
}

func createTable() error {
	db, err := config.ConnectDB()
	if err != nil {
		return err
	}
	defer db.Close()

	fmt.Printf("createTable")

	_, err = db.Exec(`
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
	CREATE TABLE IF NOT EXISTS expenses (
		id SERIAL PRIMARY KEY,
		amount FLOAT NOT NULL,
		category TEXT NOT NULL,
		created_at TEXT NOT NULL,
		user_id INT,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);
	CREATE TABLE IF NOT EXISTS budgets (
		id SERIAL PRIMARY KEY,
		amount FLOAT NOT NULL,
		category TEXT NOT NULL,
		created_at TEXT NOT NULL,
		user_id INT,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);
	`)

	fmt.Println(err)
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
	expenseRepo := repositories.ExpenseRepository{DB: app.DB}
	budgetRepo := repositories.BudgetRepository{DB: app.DB}
	app.initializeRoutes(&userRepo, &incomeRepo, &expenseRepo, &budgetRepo)
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

func (app *App) initializeRoutes(userRepo *repositories.UserRepository, incomeRepo *repositories.IncomeRepository, expenseRepo *repositories.ExpenseRepository, budgetRepo *repositories.BudgetRepository) {
	app.Router.HandleFunc("/user", app.createUserHandler(userRepo)).Methods("POST")
	app.Router.HandleFunc("/users", app.getUsersHandler(userRepo)).Methods("GET")
	app.Router.HandleFunc("/login", app.loginHandler(userRepo)).Methods("POST")
	app.Router.HandleFunc("/signout", app.signoutHandler(userRepo)).Methods("POST")
	app.Router.HandleFunc("/income", app.createIncome(incomeRepo)).Methods("POST")
	app.Router.HandleFunc("/incomes", app.getAllIncomes(incomeRepo)).Methods("GET")
	app.Router.HandleFunc("/income/{incomeID}", app.editIncome(incomeRepo)).Methods("PUT")
	app.Router.HandleFunc("/income/{incomeID}", app.deleteIncome(incomeRepo)).Methods("DELETE")
	app.Router.HandleFunc("/expense", app.createExpense(expenseRepo)).Methods("POST")
	app.Router.HandleFunc("/expenses", app.getAllExpenses(expenseRepo)).Methods("GET")
	app.Router.HandleFunc("/expense/{expenseID}", app.editExpense(expenseRepo)).Methods("PUT")
	app.Router.HandleFunc("/expense/{expenseID}", app.deleteExpense(expenseRepo)).Methods("DELETE")
	app.Router.HandleFunc("/expense_by_month", app.getExpensesByMonthAndCategory(expenseRepo)).Methods("GET")
	app.Router.HandleFunc("/budget", app.createBudget(budgetRepo)).Methods("POST")
	app.Router.HandleFunc("/budgets", app.getAllBudgets(budgetRepo)).Methods("GET")
	app.Router.HandleFunc("/budget/{budgetID}", app.editBudget(budgetRepo)).Methods("PUT")
	app.Router.HandleFunc("/budget/{budgetID}", app.deleteBudget(budgetRepo)).Methods("DELETE")
	app.Router.HandleFunc("/budget_by_month", app.getBudgetsByMonthAndCategory(budgetRepo)).Methods("GET")
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

func extractTokenFromHeader(r *http.Request) string {
	// Get the authorization header value
	authHeader := r.Header.Get("Authorization")
	// Check if the authorization header is empty
	if authHeader == "" {
		return ""
	}
	// Return the token string
	return authHeader
}

func (app *App) signoutHandler(userRepo *repositories.UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := extractTokenFromHeader(r)
		if tokenString == "" {
			respondWithError(w, http.StatusBadRequest, "Missing token")
			return
		}

		respondWithJSON(w, http.StatusOK, map[string]string{"message": "Successfully signed out"})
	}

}

//Income

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

		d, err := time.Parse("2006-01-02", income.Created_at)
		if err != nil {
			fmt.Println(err)
		}
		income.Created_at = d.Format("2006-01-02")
		fmt.Println(income)

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

func (a *App) editIncome(incomeRepo *repositories.IncomeRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			respondWithError(w, http.StatusBadRequest, "Missing authorization token")
			return
		}
		appConfig := config.GetConfig()
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

		// Decode request body into Income struct
		var income *models.Income
		err = json.NewDecoder(r.Body).Decode(&income)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		userID := int(claims["id"].(float64))
		fmt.Println(userID, "----------userID")

		// Get income ID from URL parameters
		incomeID := mux.Vars(r)["incomeID"]
		fmt.Println(incomeID)

		// Parse income ID to int
		id, err := strconv.Atoi(incomeID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid income ID")
			return
		}

		income.UserID = userID
		income.ID, err = strconv.Atoi(incomeID)

		fmt.Println("income---------", income)

		if err != nil {
			fmt.Println("error with incomeID")
			return
		}
		// Update income in database
		updatedIncome, err := incomeRepo.UpdateIncome(id, income)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		respondWithJSON(w, http.StatusOK, updatedIncome)
	}
}

func (a *App) deleteIncome(incomeRepo *repositories.IncomeRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get income ID from URL parameters
		incomeID := mux.Vars(r)["incomeID"]
		fmt.Println(incomeID)

		// Parse income ID to int
		id, err := strconv.Atoi(incomeID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid income ID")
			return
		}

		// Delete income from database
		if err := incomeRepo.DeleteIncome(id); err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		respondWithJSON(w, http.StatusOK, map[string]string{"message": "Income deleted successfully"})
	}
}

//Expense

func (app *App) createExpense(expenseRepo *repositories.ExpenseRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		expense := models.Expense{}
		err := json.NewDecoder(r.Body).Decode(&expense)
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

		// Set the user ID in the expense object
		expense.UserID = userID

		d, err := time.Parse("2006-01-02", expense.Created_at)
		if err != nil {
			fmt.Println(err)
		}
		expense.Created_at = d.Format("2006-01-02")
		fmt.Println(expense)

		// Create the expense record in the database
		result, err := expenseRepo.CreateExpense(&expense)

		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to create income record")
			return
		}

		// Return the created expense record
		respondWithJSON(w, http.StatusCreated, result)
	}
}

func (app *App) getAllExpenses(expenseRepo *repositories.ExpenseRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			respondWithError(w, http.StatusBadRequest, "Missing authorization token")
			return
		}
		appConfig := config.GetConfig()
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
		expenses, err := expenseRepo.GetAllExpenses(userID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		fmt.Println(expenses)

		respondWithJSON(w, http.StatusOK, expenses)
	}
}

func (a *App) editExpense(expenseRepo *repositories.ExpenseRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			respondWithError(w, http.StatusBadRequest, "Missing authorization token")
			return
		}
		appConfig := config.GetConfig()
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

		// Decode request body into Expense struct
		var expense *models.Expense
		err = json.NewDecoder(r.Body).Decode(&expense)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		userID := int(claims["id"].(float64))
		fmt.Println(userID, "----------userID")

		// Get expense ID from URL parameters
		expenseID := mux.Vars(r)["expenseID"]
		fmt.Println(expenseID)

		// Parse expense ID to int
		id, err := strconv.Atoi(expenseID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid expense ID")
			return
		}

		expense.UserID = userID
		expense.ID, err = strconv.Atoi(expenseID)

		fmt.Println("expense---------", expense)

		if err != nil {
			fmt.Println("error with expenseID")
			return
		}
		// Update expense in database
		updatedExpense, err := expenseRepo.UpdateExpense(id, expense)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		respondWithJSON(w, http.StatusOK, updatedExpense)
	}
}

func (a *App) deleteExpense(expenseRepo *repositories.ExpenseRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get expense ID from URL parameters
		expenseID := mux.Vars(r)["expenseID"]
		fmt.Println(expenseID)

		// Parse expense ID to int
		id, err := strconv.Atoi(expenseID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid expense ID")
			return
		}

		// Delete expense from database
		if err := expenseRepo.DeleteExpense(id); err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		respondWithJSON(w, http.StatusOK, map[string]string{"message": "Expense deleted successfully"})
	}
}

//Budget

func (a *App) createBudget(budgetRepo *repositories.BudgetRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		budget := models.Budget{}
		err := json.NewDecoder(r.Body).Decode(&budget)
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

		// Set the user ID in the budget object
		budget.UserID = userID

		d, err := time.Parse("2006-01-02", budget.Created_at)
		if err != nil {
			fmt.Println(err)
		}
		budget.Created_at = d.Format("2006-01-02")
		fmt.Println(budget)

		// Create the budget record in the database
		result, err := budgetRepo.CreateBudget(&budget)

		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to create income record")
			return
		}

		// Return the created budget record
		respondWithJSON(w, http.StatusCreated, result)
	}
}

func (app *App) getAllBudgets(budgetRepo *repositories.BudgetRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			respondWithError(w, http.StatusBadRequest, "Missing authorization token")
			return
		}
		appConfig := config.GetConfig()
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
		budgets, err := budgetRepo.GetAllbudgets(userID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		fmt.Println(budgets)

		respondWithJSON(w, http.StatusOK, budgets)
	}
}

func (a *App) editBudget(budgetRepo *repositories.BudgetRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			respondWithError(w, http.StatusBadRequest, "Missing authorization token")
			return
		}
		appConfig := config.GetConfig()
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

		// Decode request body into Budget struct
		var budget *models.Budget
		err = json.NewDecoder(r.Body).Decode(&budget)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		userID := int(claims["id"].(float64))
		fmt.Println(userID, "----------userID")

		// Get budget ID from URL parameters
		budgetID := mux.Vars(r)["budgetID"]
		fmt.Println(budgetID)

		// Parse budget ID to int
		id, err := strconv.Atoi(budgetID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid budget ID")
			return
		}

		budget.UserID = userID
		budget.ID, err = strconv.Atoi(budgetID)

		fmt.Println("budget---------", budget)

		if err != nil {
			fmt.Println("error with budgetID")
			return
		}
		// Update budget in database
		updatedBudget, err := budgetRepo.UpdateBudget(id, budget)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		respondWithJSON(w, http.StatusOK, updatedBudget)
	}
}

func (a *App) deleteBudget(budgetRepo *repositories.BudgetRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get budget ID from URL parameters
		budgetID := mux.Vars(r)["budgetID"]
		fmt.Println(budgetID)

		// Parse budget ID to int
		id, err := strconv.Atoi(budgetID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid budget ID")
			return
		}

		// Delete budget from database
		if err := budgetRepo.DeleteBudget(id); err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		respondWithJSON(w, http.StatusOK, map[string]string{"message": "Budget deleted successfully"})
	}
}

func (a *App) getBudgetsByMonthAndCategory(budgetRepo *repositories.BudgetRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Retrieve all budgets from the database
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			respondWithError(w, http.StatusBadRequest, "Missing authorization token")
			return
		}
		appConfig := config.GetConfig()
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

		budgets, err := budgetRepo.GetAllbudgets(userID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		// Create a map to store the budgets grouped by month and category
		budgetMap := make(map[string]map[string][]models.Budget)

		// Group budgets by month and category
		for _, budget := range budgets {
			// Convert the budget's created_at time to a string in the format "YYYY-MM"

			d, err := time.Parse("2006-01-02", budget.Created_at)
			if err != nil {
				fmt.Println(err)
			}
			monthStr := d.Format("2006-01-02")

			// Check if a map entry exists for the month
			_, ok := budgetMap[monthStr]
			if !ok {
				// Create a new map entry for the month
				budgetMap[monthStr] = make(map[string][]models.Budget)
			}

			// Check if a map entry exists for the budget's category
			_, ok = budgetMap[monthStr][budget.Category]
			if !ok {
				// Create a new map entry for the category
				budgetMap[monthStr][budget.Category] = []models.Budget{}
			}

			// Add the budget to the appropriate category list
			budgetMap[monthStr][budget.Category] = append(budgetMap[monthStr][budget.Category], budget)
		}

		// Return the budget map as JSON
		respondWithJSON(w, http.StatusOK, budgetMap)
	}
}

func (a *App) getExpensesByMonthAndCategory(expenseRepo *repositories.ExpenseRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Retrieve all expenses from the database
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			respondWithError(w, http.StatusBadRequest, "Missing authorization token")
			return
		}
		appConfig := config.GetConfig()
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

		expenses, err := expenseRepo.GetAllExpenses(userID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		// Create a map to store the expenses grouped by month and category
		expenseMap := make(map[string]map[string][]models.Expense)

		// Group expenses by month and category
		for _, expense := range expenses {
			// Convert the expense's created_at time to a string in the format "YYYY-MM"

			d, err := time.Parse("2006-01-02", expense.Created_at)
			if err != nil {
				fmt.Println(err)
			}
			monthStr := d.Format("2006-01")

			// Check if a map entry exists for the month
			_, ok := expenseMap[monthStr]
			if !ok {
				// Create a new map entry for the month
				expenseMap[monthStr] = make(map[string][]models.Expense)
			}

			// Check if a map entry exists for the expense's category
			_, ok = expenseMap[monthStr][expense.Category]
			if !ok {
				// Create a new map entry for the category
				expenseMap[monthStr][expense.Category] = []models.Expense{}
			}

			// Add the expense to the appropriate category list
			expenseMap[monthStr][expense.Category] = append(expenseMap[monthStr][expense.Category], expense)
		}

		// Return the expense map as JSON
		respondWithJSON(w, http.StatusOK, expenseMap)
	}
}
