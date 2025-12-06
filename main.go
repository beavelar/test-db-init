package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

var logger *log.Logger

const (
	numUsers           = 1000
	minMessagesPerSeed = 1000
	maxMessagesPerSeed = 10000
)

func init() {
	logger = log.NewWithOptions(os.Stdout, log.Options{
		TimeFormat:   time.DateTime,
		Level:        log.DebugLevel,
		ReportCaller: true,
		Formatter:    log.JSONFormatter,
	})
}

func main() {
	dbConnStr := os.Getenv("DATABASE_URL")
	if dbConnStr == "" {
		logger.Fatal("DATABASE_URL environment variable not set. Please set it to your PostgreSQL connection string.")
	}

	db, err := sql.Open("postgres", dbConnStr)
	if err != nil {
		logger.Fatal("Error opening database", "error", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		logger.Fatal("Error connecting to the database", "error", err)
	}
	logger.Info("Successfully connected to the database!")

	logger.Info("Executing init.sql")
	sqlFile, err := os.ReadFile("init.sql")
	if err != nil {
		logger.Fatal("Error reading SQL file", "error", err)
	}

	_, err = db.Exec(string(sqlFile))
	if err != nil {
		logger.Fatal("Error executing SQL script", "error", err)
	}

	logger.Info("SQL script (init.sql) executed successfully!")

	logger.Info("Generating and inserting users", "numUsers", numUsers)
	userIDs := make([]uuid.UUID, numUsers)

	stmt, err := db.Prepare("INSERT INTO users (id, username) VALUES ($1, $2)")
	if err != nil {
		logger.Fatal("Failed to prepare statement for users", "error", err)
	}
	defer stmt.Close()

	for i := range numUsers {
		userID := uuid.New()
		username := fmt.Sprintf("user_%04d", i+1)
		_, err = stmt.Exec(userID, username)
		if err != nil {
			logger.Fatal("Error inserting user", "user", username, "error", err)
		}
		userIDs[i] = userID
	}

	logger.Info("Users inserted successfully!", "numUsers", numUsers)

	msgStmt, err := db.Prepare("INSERT INTO messages (user_id, message) VALUES ($1, $2)")
	if err != nil {
		logger.Fatal("Failed to prepare statement for messages", "error", err)
	}
	defer msgStmt.Close()

	for _, user := range userIDs {
		numMessages := rand.Intn(maxMessagesPerSeed-minMessagesPerSeed+1) + minMessagesPerSeed
		logger.Info("Generating and inserting messages", "numMessages", numMessages, "user", user)

		for i := range numMessages {
			randomUserIndex := rand.Intn(numUsers)
			userID := userIDs[randomUserIndex]
			messageContent := fmt.Sprintf("Hello from user %s! This is message number %d.", userID.String()[:8], i+1)

			_, err = msgStmt.Exec(userID, messageContent)
			if err != nil {
				logger.Fatal("Error inserting message", "messageIdx", i, "error", err)
			}
		}

		logger.Info("Messages inserted successfully!", "numMessages", numMessages)
	}
}
