package main

import (
	"log"
	"os"

	"application-wallet/config"

	"application-wallet/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	// err := godotenv.Load("../.env")
	// if err != nil {
	// 	log.Fatalf("Error loading .env file: %v", err)
	// }

	if err := godotenv.Load("../.env"); err != nil {
		log.Println("No .env file found, using environment variables from container/host")
	}

	port := os.Getenv("PORT_URL")

	db, err := config.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	r := gin.Default()
	routes.SetupRoutes(r, db)

	r.Run(":" + port)
}

// import (
// 	"database/sql"
// 	"fmt"
// 	"log"
// 	"os"

// 	_ "github.com/lib/pq"
// 	"application-wallet/utils"
// 	"github.com/joho/godotenv"
// )

// func main() {
// 	// Load .env file
// 	err := godotenv.Load()
// 	if err != nil {
// 		log.Fatalf("Error loading .env file: %v", err)
// 	}

// 	user := os.Getenv("DB_USER")
// 	pw := os.Getenv("DB_PASSWORD")
// 	host := os.Getenv("DB_HOST")
// 	port := os.Getenv("DB_PORT")
// 	dbname := os.Getenv("DB_NAME")

// 	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, pw, host, port, dbname)
// 	fmt.Println(dsn)

// 	db, err := sql.Open("postgres", dsn)
// 	if err != nil {
// 		log.Fatalf("Failed to connect to database: %v", err)
// 	}
// 	defer db.Close()

	// // Data seed admin
	// name := "Kevin Prayoga"
	// email := "kevinp@example.com"
	// password := "mypassword123" // plaintext password
	// pin := "123456"            // plaintext PIN
	// balance := 100000.00

// 	// Data seed user
// 	name := "User 1"
// 	email := "user1@example.com"
// 	password := "mypassword456" // plaintext password
// 	pin := "456789"            // plaintext PIN
// 	balance := 100000.00

// 	// Hash password dan PIN
// 	hashedPassword := utils.HashString(password)
// 	hashedPin := utils.HashString(pin)

// 	query := `
// 		INSERT INTO users (name, email, password_hash, pin_hash, balance)
// 		VALUES ($1, $2, $3, $4, $5)
// 	`
// 	fmt.Println(name, email, hashedPassword, hashedPin, balance)
// 	_, err = db.Exec(query, name, email, hashedPassword, hashedPin, balance)
// 	if err != nil {
// 		log.Fatalf("Failed to seed user data: %v", err)
// 	}

// 	fmt.Println("Seed data successfully inserted!")
// }
