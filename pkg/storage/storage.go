package storage

import (
	"encoding/json"
	"fmt"
	"log"
	"main/pkg/transaction"
	"os"
)

// SaveFilePath File to save responses
const SaveFilePath = "responses.txt"

// SaveResponseToFile saves the transaction to a file.
func SaveResponseToFile(response transaction.Transaction) {
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

// Replace the old saveResponseToFile function
func SaveTransactionToDB(response transaction.TransactionV2) error {
	insertSQL := `
        INSERT INTO transactions (name, amount, currency, date, is_claimable, paid_for_family)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id; -- Optional: get the ID of the inserted row
        `
	var insertedID int64

	err := db.QueryRow(
		insertSQL,
		response.Name,
		response.Amount,
		response.Currency,
		response.Date, // Pass the time.Time object directly
		response.IsClaimable,
		response.PaidForFamily,
	).Scan(&insertedID) // Scan the returned ID

	if err != nil {
		log.Printf("Error inserting transaction into database: %v", err)
		return err
	}

	log.Printf("Successfully inserted transaction with ID: %d", insertedID)
	return nil
}
