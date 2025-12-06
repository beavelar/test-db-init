package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

const (
	numUsers           = 1000
	minMessagesPerSeed = 1000
	maxMessagesPerSeed = 10000
)

func main() {
	dbConnStr := os.Getenv("DATABASE_URL")
	if dbConnStr == "" {
		log.Fatal("DATABASE_URL environment variable not set. Please set it to your PostgreSQL connection string.")
	}

	db, err := sql.Open("postgres", dbConnStr)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	fmt.Println("Successfully connected to the database!")

	fmt.Println("Executing init.sql...")
	sqlFile, err := os.ReadFile("init.sql")
	if err != nil {
		log.Fatalf("Error reading SQL file: %v", err)
	}

	_, err = db.Exec(string(sqlFile))
	if err != nil {
		log.Fatalf("Error executing SQL script: %v", err)
	}

	fmt.Println("SQL script (init.sql) executed successfully!")

	fmt.Printf("Generating and inserting %d users...\n", numUsers)
	userIDs := make([]uuid.UUID, numUsers)

	stmt, err := db.Prepare("INSERT INTO users (id, username) VALUES ($1, $2)")
	if err != nil {
		log.Fatalf("Failed to prepare statement for users: %v", err)
	}
	defer stmt.Close()

	for i := 0; i < numUsers; i++ {
		userID := uuid.New()
		username := fmt.Sprintf("user_%04d", i+1)
		_, err = stmt.Exec(userID, username)
		if err != nil {
			log.Fatalf("Error inserting user %s: %v", username, err)
		}
		userIDs[i] = userID
	}

	fmt.Printf("%d users inserted successfully!\n", numUsers)

	numMessages := rand.Intn(maxMessagesPerSeed-minMessagesPerSeed+1) + minMessagesPerSeed
	fmt.Printf("Generating and inserting %d messages...\n", numMessages)

	msgStmt, err := db.Prepare("INSERT INTO messages (user_id, message) VALUES ($1, $2)")
	if err != nil {
		log.Fatalf("Failed to prepare statement for messages: %v", err)
	}
	defer msgStmt.Close()

	for i := 0; i < numMessages; i++ {
		randomUserIndex := rand.Intn(numUsers)
		userID := userIDs[randomUserIndex]
		messageContent := fmt.Sprintf("Hello from user %s! This is message number %d.", userID.String()[:8], i+1)

		_, err = msgStmt.Exec(userID, messageContent)
		if err != nil {
			log.Fatalf("Error inserting message %d: %v", i, err)
		}
	}

	fmt.Printf("%d messages inserted successfully!\n", numMessages)
}
