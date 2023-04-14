package models

type Expense struct {
	ID         int
	Amount     float64
	Category   string
	Created_at string
	UserID     int
}

// {
// 	"amount" : 200.00,
// 	"category": "Rent",
// 	"created_at" : "2023-04-14T10:00:00Z"
// }
