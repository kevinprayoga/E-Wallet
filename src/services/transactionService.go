package services

import (
	"errors"
	"strings"
	"time"

	"application-wallet/models"
	"application-wallet/repositories"
  log "github.com/sirupsen/logrus"
)

type TransactionService struct {
	Repo *repositories.TransactionRepository
}

func (s *TransactionService) TopUp(userID string, amount float64, source string) error {
	// Validate amount
	if amount < 0 {
		log.WithFields(log.Fields{
			"user_id": userID,
			"amount":  amount,
		}).Warn("Top up amount must be greater than 0")
		return errors.New("amount must be greater than 0")
	}

	// Validate daily limit
	total, err := s.Repo.TotalTransactionOneDay(userID, "topup")
	if err != nil {
		log.WithFields(log.Fields{
			"user_id": userID,
			"amount":  amount,
		}).Error("failed to get total transaction")
		return errors.New("failed to get total transaction")
	}

	if (amount + total) > 10000000 {
		log.WithFields(log.Fields{
			"user_id": userID,
			"amount":  amount,
		}).Warn("Exceeded daily top up limit")
		return errors.New("maximum top up is Rp10,000,000 per day")
	}

	// Validate source
	finalSource := strings.ToLower(source)
	if finalSource != "bank_transfer" && finalSource != "credit_card" && finalSource != "e-wallet" {
		log.WithFields(log.Fields{
			"user_id": userID,
			"amount":  amount,
			"source":  source,
		}).Warn("Invalid top up source")
		return errors.New("invalid source")
	}

	// Validate balance
	balance, err := s.Repo.GetUserBalance(userID)
	if err != nil {	
		log.WithFields(log.Fields{
			"user_id": userID,
		}).Error("failed to get user balance")
		return errors.New("failed to get user balance")
	}

	// Top-up
	newBalance := balance + amount
	err = s.Repo.UpdateUserBalance(userID, newBalance)
	if err != nil {
		log.WithFields(log.Fields{
			"user_id": userID,
			"amount":  amount,
		}).Error("failed to update user balance")
		return errors.New("failed to update user balance")
	}

	// Create transaction
	transaction := models.Transaction{
		UserID: userID,
		Type: strings.ToLower("topup"),
		Source: finalSource,
		Amount: amount,
		BalanceBefore: balance,
		BalanceAfter: newBalance,
		Description: "Top-up menggunakan metode pembayaran " + source,
	}
	return s.Repo.CreateTransaction(transaction)
}

func (s *TransactionService) Withdraw(userID string, amount float64, pin, bankCode, accountNumber string, now time.Time) error {
	// Validate pin
	if err := s.Repo.ValidatePin(userID, pin); err != nil {
		log.WithFields(log.Fields{
			"user_id": userID,
		}).Error("failed to validate pin")
		return err
	}

	// Validate amount
	if amount < 0 {
		log.WithFields(log.Fields{
			"user_id": userID,
			"amount":  amount,
		}).Warn("Withdraw amount must be greater than 0")
		return errors.New("amount must be greater than 0")
	}

	// Validate balance
	balance, err := s.Repo.GetUserBalance(userID)
	if err != nil {
		log.WithFields(log.Fields{
			"user_id": userID,
		}).Error("failed to get user balance")
		return errors.New("failed to get user balance")
	}
	if balance < amount {
		log.WithFields(log.Fields{
			"user_id": userID,
			"amount":  amount,
		}).Error("insufficient balance")
		return errors.New("insufficient balance")
	}

	// Validate daily limit
	total, err := s.Repo.TotalTransactionOneDay(userID, "withdraw")
	if err != nil {
		log.WithFields(log.Fields{
			"user_id": userID,
			"amount":  amount,
		}).Error("failed to get total transaction")
		return errors.New("failed to get total transaction")
	}

	if (amount + total) > 10000000 {
		failedTransaction := models.WithdrawRequest{
			UserID: userID,
			Amount: amount,
			Status: strings.ToLower("rejected"),
			Reason: "Exceeded daily withdrawal limit",
		}
		_ = s.Repo.CreateWithdrawRequest(failedTransaction)
		log.WithFields(log.Fields{
			"user_id": userID,
			"amount":  amount,
		}).Warn("Exceeded daily withdraw limit")
		return errors.New("maximum withdraw is Rp10,000,000 per day")
	}

	// Validate bank code
	isBankAvailable, err := s.Repo.IsBankAvailable(bankCode)
	if err != nil || !isBankAvailable {
		log.WithFields(log.Fields{
			"user_id":   userID,
			"amount":    amount,
			"bank_code": bankCode,
		}).Warn("Bank is not available")
		return errors.New("bank is not available")
	}

	// Validate pending withdrawal request
	isPendingAvailable, err := s.Repo.IsPendingWithdrawRequest(userID)
	if err != nil || isPendingAvailable {
		log.WithFields(log.Fields{
			"user_id": userID,
			"amount":  amount,
		}).Warn("There is a pending withdrawal request")
		return errors.New("pending withdrawal request")
	}

	// Validate operational time
	startHour, endHour := 7, 15
	if now.Hour() < startHour || now.Hour() > endHour {
		withdrawRequest := models.WithdrawRequest{
			UserID: userID,
			Amount: amount,
			Status: strings.ToLower("pending"),
			Reason: "Outside operational time",
		}
		err := s.Repo.CreateWithdrawRequest(withdrawRequest)
		if err != nil {
			log.WithFields(log.Fields{
				"user_id": userID,
				"amount":  amount,
			}).Error("failed to create pending withdrawal request")
			return errors.New("failed to create pending withdrawal request")
		}
		log.WithFields(log.Fields{
			"user_id": userID,
			"amount":  amount,
		}).Info("Withdrawal is pending due to outside operational time")
		return errors.New("withdrawal is pending due to outside operational time")
	}

	// Withdraw
	newBalance := balance - amount
	err = s.Repo.UpdateUserBalance(userID, newBalance)
	if err != nil {
		failedTransaction := models.WithdrawRequest{
			UserID: userID,
			Amount: amount,
			Status: strings.ToLower("rejected"),
			Reason: "Failed to update user balance",
		}
		_ = s.Repo.CreateWithdrawRequest(failedTransaction)
		log.WithFields(log.Fields{
			"user_id": userID,
			"amount":  amount,
		}).Error("failed to update user balance")
		return errors.New("failed to update user balance")
	}

	// Log successful withdrawal
	withdrawRequest := models.WithdrawRequest{
		UserID: userID,
		Amount: amount,
		Status: strings.ToLower("approved"),
		Reason: "Verified and processed",
	}
	err = s.Repo.CreateWithdrawRequest(withdrawRequest)
	if err != nil {
		log.WithFields(log.Fields{
			"user_id": userID,
			"amount":  amount,
		}).Error("failed to log successful withdrawal")
		return errors.New("failed to log successful withdrawal")
	}

	// Create transaction
	transaction := models.Transaction{
		UserID: userID,
		Type: strings.ToLower("withdrawal"),
		Source: strings.ToLower("bank_transfer"),
		Amount: amount,
		BalanceBefore: balance,
		BalanceAfter: newBalance,
		Description: "Withdrawal to " + bankCode + " - " + accountNumber,
	}
	return s.Repo.CreateTransaction(transaction)
}

func (s *TransactionService) UpdatePendingTransaction(userID, pin string) error {
	// Validate pin
	if err := s.Repo.ValidatePin(userID, pin); err != nil {
		log.WithFields(log.Fields{
			"user_id": userID,
		}).Error("failed to validate pin")
		return err
	}

	// Validate admin
	roleName, err := s.Repo.GetUserRole(userID)
	if err != nil || roleName != "admin" {
		log.WithFields(log.Fields{
			"user_id": userID,
		}).Warn("unauthorized access to update pending transaction")
		return errors.New("unauthorized")
	}

	// Get pending withdrawal request
	dataPendingWithdraw, err := s.Repo.GetPendingWithdrawRequest()
	if err != nil {
		log.WithFields(log.Fields{
			"user_id": userID,
		}).Error("failed to get pending withdrawal request")
		return errors.New("failed to get pending withdrawal request")
	}

	var withdrawRequestArray []models.WithdrawRequest
	for dataPendingWithdraw.Next() {
		var withdrawRequest models.WithdrawRequest
		dataPendingWithdraw.Scan(&withdrawRequest.ID, &withdrawRequest.UserID, &withdrawRequest.Amount, &withdrawRequest.Status, &withdrawRequest.RequestedAt, &withdrawRequest.ProcessedAt, &withdrawRequest.Reason)

		// Validate balance
		balance, err := s.Repo.GetUserBalance(withdrawRequest.UserID)
		if err != nil {
			log.WithFields(log.Fields{
				"user_id": withdrawRequest.UserID,
			}).Error("failed to get user balance")
			return errors.New("failed to get user balance")
		}
		if balance < withdrawRequest.Amount {
			log.WithFields(log.Fields{
				"user_id": withdrawRequest.UserID,
				"amount":  withdrawRequest.Amount,
			}).Error("insufficient balance")
			return errors.New("insufficient balance")
		}
		
		// Withdraw
		newBalance := balance - withdrawRequest.Amount
		err = s.Repo.UpdateUserBalance(withdrawRequest.UserID, newBalance)
		if err != nil {
			log.WithFields(log.Fields{
				"user_id": withdrawRequest.UserID,
				"amount":  withdrawRequest.Amount,
			}).Error("failed to update user balance")
			return errors.New("failed to update user balance")
		}

		// Create transaction
		transaction := models.Transaction{
			UserID: withdrawRequest.UserID,
			Type: strings.ToLower("withdrawal"),
			Source: strings.ToLower("bank_transfer"),
			Amount: withdrawRequest.Amount,
			BalanceBefore: balance,
			BalanceAfter: newBalance,
			Description: "Withdrawal to bank account",
		}
		err = s.Repo.CreateTransaction(transaction)
		if err != nil {
			log.WithFields(log.Fields{
				"user_id": withdrawRequest.UserID,
				"amount":  withdrawRequest.Amount,
			}).Error("failed to create transaction")
			return errors.New("failed to create transaction")
		}

		withdrawRequestArray = append(withdrawRequestArray, withdrawRequest)
	}

	if len(withdrawRequestArray) == 0 {
		log.WithFields(log.Fields{
			"user_id": userID,
		}).Warn("no pending withdrawal request")
		return errors.New("no pending withdrawal request")
	}

	status := strings.ToLower("approved")
	reason := "Verified and processed"

	err = s.Repo.UpdatePendingTransaction(status, reason)
	if err != nil {
		log.WithFields(log.Fields{
			"user_id": userID,
		}).Error("failed to update pending transaction")
		return errors.New("failed to update pending transaction")
	}

	return nil
}

func (s *TransactionService) GetTransactions(userID string) ([]models.Transaction, error) {
	// Validate admin
	roleName, err := s.Repo.GetUserRole(userID)
	if err != nil || roleName != "admin" {
		log.WithFields(log.Fields{
			"user_id": userID,
		}).Warn("unauthorized access to get all transactions")
		return nil, errors.New("unauthorized")
	}

	transactions, err := s.Repo.GetTransactionsData(userID)
	if err != nil {
		log.WithFields(log.Fields{
			"user_id": userID,
		}).Error("failed to get transactions")
		return nil, errors.New("failed to get transactions")
	}

	var transactionsArr []models.Transaction
	for transactions.Next() {
		var t models.Transaction
		err := transactions.Scan(&t.ID, &t.UserID, &t.Type, &t.Source, &t.Amount, &t.BalanceBefore, &t.BalanceAfter, &t.Date, &t.Description)
		if err != nil {
			log.WithFields(log.Fields{
				"user_id": userID,
			}).Error("failed to scan transaction")
			return nil, err
		}
		transactionsArr = append(transactionsArr, t)
	}

	return transactionsArr, nil
}
