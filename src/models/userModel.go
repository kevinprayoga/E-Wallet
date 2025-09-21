package models

type User struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Email    string  `json:"email"`
	Balance  float64 `json:"balance"`
	PinHash  string  `json:"-"`
	Password string  `json:"-"`
}

type Transaction struct {
	ID            string  `json:"id"`
	UserID        string  `json:"user_id"`
	Type          string  `json:"type"`
	Source				string  `json:"source"`
	Amount        float64 `json:"amount"`
	BalanceBefore float64 `json:"balance_before"`
	BalanceAfter  float64 `json:"balance_after"`
	Date          string  `json:"date"`
	Description   string  `json:"description"`
}

type WithdrawRequest struct {
	ID          string  `json:"id"`
	UserID      string  `json:"user_id"`
	Amount      float64 `json:"amount"`
	Status      string  `json:"status"`
	RequestedAt string  `json:"requested_at"`
	ProcessedAt string  `json:"processed_at"`
	Reason      string  `json:"reason"`
}
