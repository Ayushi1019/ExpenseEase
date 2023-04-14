package models

type Budget struct {
	ID          int    `json:"id"`
	Income_ids  []int  `json:"income_ids"`
	Expense_ids []int  `json:"expense_ids"`
	Month       string `json:"month"`
	User_id     int    `json:"user_id"`
}
