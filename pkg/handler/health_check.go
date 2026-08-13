package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// ErrorResponse represents a standard JSON error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// HealthStatus represents the health check response
type HealthStatus struct {
	Status     string    `json:"status"`
	Timestamp  time.Time `json:"timestamp"`
	Database   string    `json:"database,omitempty"`
	DatabaseAt time.Time `json:"database_at,omitempty"`
}

// NewErrorResponse creates a standard error response
func NewErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	resp := ErrorResponse{
		Error:   http.StatusText(statusCode),
		Message: message,
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error encoding error response JSON: %v", err)
	}
}

// HealthCheckHandler checks server and database health.
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		NewErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	status := HealthStatus{
		Status:    "ok",
		Timestamp: time.Now().UTC(),
	}

	// Check database connectivity
	db := getDatabase()
	if db != nil {
		err := db.Ping()
		if err != nil {
			status.Status = "degraded"
			status.Database = "unreachable"
		} else {
			status.Database = "connected"
			status.DatabaseAt = time.Now().UTC()
		}
	} else {
		status.Database = "not configured"
	}

	w.Header().Set("Content-Type", "application/json")
	if status.Status == "degraded" {
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	if err := json.NewEncoder(w).Encode(status); err != nil {
		log.Printf("Error encoding health status JSON: %v", err)
	}
	log.Printf("Served %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
}

