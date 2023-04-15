package models

type Budget struct {
	ID         int
	Created_at string
	Amount     float64
	Category   string
	UserID     int
}
