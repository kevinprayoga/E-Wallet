package controllers

import (
	"database/sql"
	"fmt"
	"net/http"

	"application-wallet/utils"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	DB *sql.DB
}

func (a *AuthController) Login(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	var (
		userID 					string
		hashedPassword	string
	)

	query := `SELECT id, password_hash FROM users WHERE LOWER(email) = LOWER($1)`
	err := a.DB.QueryRow(query, email).Scan(&userID, &hashedPassword)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusUnauthorized, utils.Data(http.StatusUnauthorized, []interface{}{}, 0, "Invalid email"))
		return
	}

	// Validate password
	err = utils.ValidateHashedString(hashedPassword, password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.Data(http.StatusUnauthorized, []interface{}{}, 0, "Invalid password"))
		return
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Data(http.StatusInternalServerError, []interface{}{}, 0, "Failed to generate token"))
		return
	}

	c.JSON(http.StatusOK, utils.Data(http.StatusOK, gin.H{"token": token}, 1, "Login success"))
}
