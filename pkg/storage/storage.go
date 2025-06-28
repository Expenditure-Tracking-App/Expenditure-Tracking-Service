package storage

import (
	"encoding/json"
	"fmt"
	"log"
	"main/pkg/transaction" // Assuming Transaction is here
	"os"
)

// SaveFilePath File to save responses
const SaveFilePath = "responses.txt"

// SaveResponseToFile saves the transaction to a file.
func SaveResponseToFile(response transaction.Transaction) { // Assuming Transaction is an older version
	file, err := os.OpenFile(SaveFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Error opening file: %v", err)
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Printf("Error closing file: %v", err)
		}
	}(file)

	data, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error marshaling response: %v", err)
		return
	}

	_, err = file.WriteString(fmt.Sprintf("%s\n", data))
	if err != nil {
		log.Printf("Error writing to file: %v", err)
	}
}

// SaveTransactionToDB saves the transaction to the database.
func SaveTransactionToDB(response transaction.Transaction) error {
	insertSQL := `
        INSERT INTO transactions (name, amount, currency, date, is_claimable, paid_for_family, category)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id;
        `
	var insertedID int64

	currentDB, err := GetDB()
	if err != nil {
		log.Printf("Error getting DB connection for insert: %v", err)
		return fmt.Errorf("failed to get DB connection: %w", err)
	}

	err = currentDB.QueryRow(
		insertSQL,
		response.Name,
		response.Amount,
		response.Currency,
		response.Date,
		response.IsClaimable,
		response.PaidForFamily,
		response.Category,
	).Scan(&insertedID)

	if err != nil {
		log.Printf("Error inserting transaction into database: %v", err)
		return fmt.Errorf("database insert failed: %w", err)
	}

	log.Printf("Successfully inserted transaction with ID: %d", insertedID)
	return nil
}
