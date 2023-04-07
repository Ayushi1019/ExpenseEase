package models

import "time"

type Income struct {
	ID     int
	UserID int
	Amount float64
	Source string
	Date   time.Time
}
