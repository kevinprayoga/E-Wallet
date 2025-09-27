package controllers

import (
	"net/http"
	"strconv"
	"time"

	"application-wallet/services"
	"application-wallet/utils"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type TransactionController struct {
	Service *services.TransactionService
}

func (co *TransactionController) TopUp(c *gin.Context) {
	sessionUserID, exists := c.Get("user_id")
	if !exists {
		log.WithFields(log.Fields{
			"user_id": sessionUserID,
		}).Error("failed to get user ID from session")
		c.JSON(http.StatusUnauthorized, utils.Data(http.StatusUnauthorized, []interface{}{}, 0, "Unauthorized"))
		return
	}

	userID := c.Param("userID")
	if sessionUserID != userID {
		log.WithFields(log.Fields{
			"user_id": sessionUserID,
		}).Warn("unauthorized access to top up")
		c.JSON(http.StatusForbidden, utils.Data(http.StatusForbidden, []interface{}{}, 0, "Access denied"))
		return
	}

	amount,  _ := strconv.ParseFloat(c.PostForm("amount"), 64)
	source := c.PostForm("source")

	err := co.Service.TopUp(userID, amount, source)
	if err != nil {
		log.WithFields(log.Fields{
			"user_id": userID,
			"amount":  amount,
			"source":  source,
			"error":   err.Error(),
		}).Error("failed to top up")
		c.JSON(http.StatusBadRequest, utils.Data(http.StatusBadRequest, []interface{}{}, 0, err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.Data(http.StatusOK, []interface{}{}, 0, "Top up success"))
}

func (co *TransactionController) Withdraw(c *gin.Context) {
	sessionUserID, exists := c.Get("user_id")
	if !exists {
		log.WithFields(log.Fields{
			"user_id": sessionUserID,
		}).Error("failed to get user ID from session")
		c.JSON(http.StatusUnauthorized, utils.Data(http.StatusUnauthorized, []interface{}{}, 0, "Unauthorized"))
		return
	}

	userID := c.Param("userID")
	if sessionUserID != userID {
		log.WithFields(log.Fields{
			"user_id": sessionUserID,
		}).Warn("unauthorized access to withdraw")
		c.JSON(http.StatusForbidden, utils.Data(http.StatusForbidden, []interface{}{}, 0, "Access denied"))
		return
	}

	amount, _ := strconv.ParseFloat(c.PostForm("amount"), 64)
	pin := c.PostForm("pin")
	bankCode := c.PostForm("bank_code")
	accountNumber := c.PostForm("account_number")

	now := time.Now()

	err := co.Service.Withdraw(userID, amount, pin, bankCode, accountNumber, now)
	if err != nil {
		if err.Error() == "withdrawal is pending due to outside operational time" {
			log.WithFields(log.Fields{
				"user_id": userID,
				"amount":  amount,
			}).Info("Withdrawal is pending due to outside operational time")
			c.JSON(http.StatusAccepted, utils.Data(http.StatusAccepted, []interface{}{}, 0, err.Error()))
			return
		}
		log.WithFields(log.Fields{
			"user_id": userID,
			"amount":  amount,
			"error":   err.Error(),
		}).Error("failed to withdraw")
		c.JSON(http.StatusBadRequest, utils.Data(http.StatusBadRequest, []interface{}{}, 0, err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.Data(http.StatusOK, []interface{}{}, 0, "Withdraw success"))
}

func (co *TransactionController) UpdatePendingTransaction(c *gin.Context) {
	sessionUserID, exists := c.Get("user_id")
	if !exists {
		log.WithFields(log.Fields{
			"user_id": sessionUserID,
		}).Error("failed to get user ID from session")
		c.JSON(http.StatusUnauthorized, utils.Data(http.StatusUnauthorized, []interface{}{}, 0, "Unauthorized"))
		return
	}

	userID := sessionUserID.(string)
	pin := c.PostForm("pin")

	err := co.Service.UpdatePendingTransaction(userID, pin)
	if err != nil {
		if err.Error() == "no pending withdrawal request" {
			log.WithFields(log.Fields{
				"user_id": userID,
			}).Warn("no pending withdrawal request")
			c.JSON(http.StatusAccepted, utils.Data(http.StatusAccepted, []interface{}{}, 0, err.Error()))
			return
		}
		log.WithFields(log.Fields{
			"user_id": userID,
			"error":   err.Error(),
		}).Error("failed to update pending transaction")
		c.JSON(http.StatusBadRequest, utils.Data(http.StatusBadRequest, []interface{}{}, 0, err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.Data(http.StatusOK, []interface{}{}, 0, "Pending transaction updated"))
}

func (co *TransactionController) GetAllTransactionHistory(c *gin.Context) {
	sessionUserID, exists := c.Get("user_id")
	if !exists {
		log.WithFields(log.Fields{
			"user_id": sessionUserID,
		}).Error("failed to get user ID from session")
		c.JSON(http.StatusUnauthorized, utils.Data(http.StatusUnauthorized, []interface{}{}, 0, "Unauthorized"))
		return
	}

	userID := c.Param("userID")
	if sessionUserID != userID {
		log.WithFields(log.Fields{
			"user_id": sessionUserID,
		}).Warn("unauthorized access to get all transaction history")
		c.JSON(http.StatusForbidden, utils.Data(http.StatusForbidden, []interface{}{}, 0, "Access denied"))
		return
	}

	transactions, err := co.Service.GetTransactions(userID)
	if err != nil {
		log.WithFields(log.Fields{
			"user_id": userID,
			"error":   err.Error(),
		}).Error("failed to get all transaction history")
		c.JSON(http.StatusBadRequest, utils.Data(http.StatusBadRequest, []interface{}{}, 0, err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.Data(http.StatusOK, transactions, len(transactions), "Success"))
}
