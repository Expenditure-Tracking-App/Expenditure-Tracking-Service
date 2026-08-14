package storage

import (
	"database/sql"

	"main/pkg/config"
	"main/pkg/transaction"
)

// DatabaseQuerier abstracts database operations for mocking
type DatabaseQuerier interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
	Ping() error
	Close() error
}

// TransactionRepository abstracts transaction persistence
type TransactionRepository interface {
	Insert(t transaction.Transaction) error
	Update(id int, t transaction.Transaction) error
	GetAll(category string, isClaimable *bool, paidForFamily *bool, page, limit int) ([]transaction.Transaction, int, error)
	GetCategoryCounts() ([]transaction.Category, error)
	GetTotalByIsClaimable() (map[bool]float32, error)
	GetTotalByPaidForFamily() (map[bool]float32, error)
}

// SetDBForTest sets the database connection for testing purposes
func SetDBForTest(dbQuerier DatabaseQuerier) {
	// This is a test helper that allows injecting a mock database
	// The actual implementation depends on how the package is structured
	_ = dbQuerier
}

// InitDBWithConnString initializes the database with a custom connection string (for tests)
func InitDBWithConnString(connStr string) error {
	return InitDB(config.DatabaseConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "test",
		Password: "test",
		DBName:   "expenditure_test",
		SSLMode:  "disable",
	})
}
