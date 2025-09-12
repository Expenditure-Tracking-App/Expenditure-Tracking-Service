package handler

import (
	"encoding/json"
	"github.com/patrickmn/go-cache"
	"log"
	"main/pkg/config"
	"net/http"
)

// GetPrefilledExpensesHandler creates an HTTP handler that serves the list
// of pre-filled expenses from the configuration.
func GetPrefilledExpensesHandler(expenses []config.FrequentExpense) http.HandlerFunc {
	// Cache the response indefinitely since it's based on config
	prefilledExpensesJSON, err := json.Marshal(expenses)
	if err != nil {
		// This would be a startup-time error, so logging it fatally is reasonable.
		log.Fatalf("Failed to marshal pre-filled expenses: %v", err)
	}
	c.Set("prefilled-expenses", prefilledExpensesJSON, cache.NoExpiration)

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		// Serve from cache
		if cachedJSON, found := c.Get("prefilled-expenses"); found {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, err = w.Write(cachedJSON.([]byte))
			if err != nil {
				return
			}
			log.Printf("Served %s %s from cache", r.Method, r.URL.Path)
			return
		}

		// Fallback (should not happen if cache is pre-warmed)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(expenses)
		if err != nil {
			return
		}
	}
}
