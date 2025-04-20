package storage

import (
	"database/sql"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"main/pkg/session"
	"main/pkg/transaction"
	"os"
)

var db *sql.DB

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
func saveTransactionToDB(response transaction.Transaction) error {
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

// Modify finishSession to call the new function
func finishSession(bot *tgbotapi.BotAPI, chatID int64, userSession *session.UserSession) error {
	// Save the responses to the database
	err := saveTransactionToDB(userSession.Answers)
	if err != nil {
		// Inform the user if saving failed
		errMsg := tgbotapi.NewMessage(chatID, "Sorry, there was an error saving your transaction. Please try again later.")
		_, sendErr := bot.Send(errMsg)
		if sendErr != nil {
			log.Printf("Error sending save error message: %v", sendErr)
		}
		// Also return the original save error
		return fmt.Errorf("failed to save transaction to DB: %w", err)
	}

	// Send a thank-you message and confirmation
	// Format the date for display using response.Date.Format("2006-01-02")
	confirmationText := fmt.Sprintf(
		"Thank you! Transaction saved.\n\nName: %s\nAmount: %.2f\nCurrency: %s\nDate: %s\nClaimable: %t\nPaid for Family: %t",
		userSession.Answers.Name,
		userSession.Answers.Amount,
		userSession.Answers.Currency,
		userSession.Answers.Date, // Format the date for display
		userSession.Answers.IsClaimable,
		userSession.Answers.PaidForFamily,
	)
	msg := tgbotapi.NewMessage(chatID, confirmationText)
	_, err = bot.Send(msg)
	if err != nil {
		// Log this error, but the transaction was already saved
		log.Printf("Error sending confirmation message: %v", err)
	}

	return nil // Return nil even if confirmation sending failed
}
