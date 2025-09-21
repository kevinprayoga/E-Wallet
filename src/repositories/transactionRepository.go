package repositories

import (
	"application-wallet/models"
	"application-wallet/utils"
	"database/sql"
)

type TransactionRepository struct {
	DB *sql.DB
}

func (r *TransactionRepository) TotalTransactionOneDay(userID string, typeTransaction string) (float64, error) {
	var total float64
	query := `
		SELECT COALESCE(SUM(amount), 0) 
		FROM transactions 
		WHERE user_id = $1 
			AND type = $2 
			AND transaction_date >= CURRENT_DATE
			AND transaction_date < CURRENT_DATE + INTERVAL '1 day'
	`
	err := r.DB.QueryRow(query, userID, typeTransaction).Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}

func (r *TransactionRepository) GetUserBalance(userID string) (float64, error) {
	var balance float64
	query := `SELECT balance FROM users WHERE id = $1`
	err := r.DB.QueryRow(query, userID).Scan(&balance)
	return balance, err
}

func (r *TransactionRepository) UpdateUserBalance(userID string, amount float64) error {
	query := `UPDATE users SET balance = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.DB.Exec(query, amount, userID)
	return err
}

func (r *TransactionRepository) CreateTransaction(t models.Transaction) error {
	query	:= `
		INSERT INTO transactions (user_id, type, source, amount, balance_before, balance_after, transaction_date, description)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), $7)
	`
	_, err := r.DB.Exec(query, t.UserID, t.Type, t.Source, t.Amount, t.BalanceBefore, t.BalanceAfter, t.Description)
	return err
}

func (r *TransactionRepository) CreateWithdrawRequest(w models.WithdrawRequest) error {
	query := `
		INSERT INTO withdraw_requests (user_id, amount, status, requested_at, processed_at, reason)
		VALUES ($1, $2, $3, NOW(), NOW(), $4)
	`
	_, err := r.DB.Exec(query, w.UserID, w.Amount, w.Status, w.Reason)
	return err
}

func (r *TransactionRepository) CreateWithdrawPendingRequest(w models.WithdrawRequest) error {
	query := `
		INSERT INTO withdraw_requests (user_id, amount, status, requested_at, reason)
		VALUES ($1, $2, $3, NOW(), $4)
	`
	_, err := r.DB.Exec(query, w.UserID, w.Amount, w.Status, w.Reason)
	return err
}

func (r *TransactionRepository) ValidatePin(userID, pin string) error {
	var pinHash string
	query := `SELECT pin_hash FROM users WHERE id = $1`
	err := r.DB.QueryRow(query, userID).Scan(&pinHash)
	if err != nil {
		return err
	}
	return utils.ValidateHashedString(pinHash, pin)
}

func (r *TransactionRepository) IsBankAvailable(bankCode string) (bool, error) {
	var isActive bool
	query := "SELECT is_active FROM banks WHERE LOWER(code) = LOWER($1)"
	err := r.DB.QueryRow(query, bankCode).Scan(&isActive)
	return isActive, err
}

func (r *TransactionRepository) IsPendingWithdrawRequest(userID string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM withdraw_requests WHERE user_id = $1 AND status = 'pending'`
	err := r.DB.QueryRow(query, userID).Scan(&count)
	return count > 0, err
}

func (r *TransactionRepository) UpdatePendingTransaction(status, reason string) error {
	query := `UPDATE withdraw_requests SET status = $1, processed_at = NOW(), reason = $2 WHERE status = 'pending'`
	_, err := r.DB.Exec(query, status, reason)
	return err
}

func (r *TransactionRepository) GetPendingWithdrawRequest() (*sql.Rows, error) {
	query := `SELECT id, user_id, amount, status, requested_at, processed_at, reason FROM withdraw_requests WHERE status = 'pending'`
	rows, err := r.DB.Query(query)
	return rows, err
}

func (r *TransactionRepository) GetUserRole(userID string) (string, error) {
	var roleName string
	query := `SELECT name FROM user_roles u INNER JOIN roles r on u.role_id = r.id WHERE u.user_id = $1`
	err := r.DB.QueryRow(query, userID).Scan(&roleName)
	return roleName, err
}

func (r *TransactionRepository) GetTransactionsData(userID string) (*sql.Rows, error) {
	query := `
		SELECT id, user_id, type, source, amount, balance_before, balance_after, transaction_date, description
		FROM transactions
		WHERE user_id = $1
		ORDER BY transaction_date DESC
	`
	rows, err := r.DB.Query(query, userID)
	return rows, err
}
