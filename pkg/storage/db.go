package storage

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	// Import the PostgreSQL driver. The underscore means we only need its side effects (registering the driver).
	_ "github.com/lib/pq"
)

// db holds the database connection pool. It's a package-level variable.
var db = GetDB()

// InitDB initializes the database connection pool using environment variables.
// It should be called once when your application starts.
func InitDB() error {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbSSLMode := os.Getenv("DB_SSLMODE") // Often 'disable' for local, 'require' or 'verify-full' for prod

	// Set default SSL mode if not provided
	if dbSSLMode == "" {
		dbSSLMode = "disable"
	}

	// Construct the connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode)

	var err error
	// Open the database connection pool. Note: sql.Open doesn't establish any connections yet.
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Printf("Error opening database connection: %v", err)
		return fmt.Errorf("failed to open database connection: %w", err)
	}

	// Ping the database to verify the connection details are correct and the DB is reachable.
	err = db.Ping()
	if err != nil {
		// Close the pool if ping fails, as it's unusable.
		err := db.Close()
		if err != nil {
			panic(err)
		}
		log.Printf("Error connecting to database: %v", err)
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Successfully connected to the database!")

	// Optionally, ensure the necessary table exists
	return createTableIfNotExists()
}

// createTableIfNotExists creates the 'transactions' table if it doesn't already exist.
func createTableIfNotExists() error {
	// SQL statement to create the table. Adjust types/constraints as needed.
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS transactions (
		id SERIAL PRIMARY KEY,
		name TEXT,
		amount NUMERIC(10, 2), -- Example: 10 total digits, 2 after decimal
		currency VARCHAR(10),
		date DATE,
		is_claimable BOOLEAN,
		paid_for_family BOOLEAN,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := db.Exec(createTableSQL)
	if err != nil {
		log.Printf("Error creating transactions table: %v", err)
		return fmt.Errorf("failed to create transactions table: %w", err)
	}
	log.Println("Transactions table checked/created successfully.")
	return nil
}

// GetDB returns the initialized database connection pool.
// Other functions in this package (or other packages, if exported) can use this
// to interact with the database.
func GetDB() *sql.DB {
	if db == nil {
		log.Fatal("Database connection pool is not initialized. Call InitDB first.")
	}
	return db
}

// CloseDB closes the database connection pool.
// It should be called when the application is shutting down gracefully.
func CloseDB() {
	if db != nil {
		err := db.Close()
		if err != nil {
			log.Printf("Error closing database connection: %v", err)
		} else {
			log.Println("Database connection closed.")
		}
	}
}
