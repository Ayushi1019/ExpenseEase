package models

import "time"

type Income struct {
	ID         int
	UserID     int
	Amount     float64
	Source     string
	Created_at time.Time
}
