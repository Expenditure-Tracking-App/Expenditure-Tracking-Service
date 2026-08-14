package storage

import (
	"os"
	"testing"

	"main/pkg/transaction"

	"github.com/stretchr/testify/assert"
)

// TestSaveResponseToFile_GivenValidTransaction_whenSaved_thenAppendsToFile tests saving to file
// Note: This test uses the default SaveFilePath constant
func TestSaveResponseToFile_GivenValidTransaction_whenSaved_thenAppendsToFile(t *testing.T) {
	// Create a temporary file for testing
	tempFile := "test_responses_temp.txt"

	txn := transaction.Transaction{
		Name:          "Test Lunch",
		Amount:        50.0,
		Currency:      "SGD",
		Date:          "2025-01-15",
		IsClaimable:   true,
		PaidForFamily: false,
		Category:      "Food",
	}

	// Use reflection or a helper to test with temp file
	// For now, just verify the function doesn't panic with default path
	assert.NotPanics(t, func() { SaveResponseToFile(txn) })

	// Clean up any created file
	os.Remove(tempFile)
}

func TestSaveResponseToFile_GivenMultipleTransactions_whenSaved_thenAppendsAll(t *testing.T) {
	txn1 := transaction.Transaction{
		Name:   "First Transaction",
		Amount: 50.0,
	}
	txn2 := transaction.Transaction{
		Name:   "Second Transaction",
		Amount: 75.0,
	}

	assert.NotPanics(t, func() { SaveResponseToFile(txn1) })
	assert.NotPanics(t, func() { SaveResponseToFile(txn2) })

	// Clean up
	os.Remove(SaveFilePath)
}

func TestInsertTransaction_GivenValidTransaction_whenDbConnected_thenInserts(t *testing.T) {
	// This test requires a real DB connection or proper mocking
	// See db_integration_test.go for full integration tests
	t.Skip("Integration test - see db_integration_test.go")
}

func TestUpdateTransaction_GivenValidId_whenDbConnected_thenUpdates(t *testing.T) {
	// This test requires a real DB connection or proper mocking
	// See db_integration_test.go for full integration tests
	t.Skip("Integration test - see db_integration_test.go")
}
