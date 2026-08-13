package handler

import (
	"database/sql"
	"github.com/patrickmn/go-cache"
	"main/pkg/storage"
	"time"
)

// --- Global Cache ---
var c = cache.New(5*time.Minute, 10*time.Minute)

// getDatabase returns the database connection pool if initialized
func getDatabase() *sql.DB {
	db, err := storage.GetDB()
	if err != nil {
		return nil
	}
	return db
}
