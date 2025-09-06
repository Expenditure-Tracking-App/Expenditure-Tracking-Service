package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// Define a struct for example JSON responses
type HealthStatus struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}

// HealthCheckHandler healthCheckHandler responds with the server's status.
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow GET requests for this endpoint
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	status := HealthStatus{
		Status:    "OK",
		Timestamp: time.Now().UTC(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err := json.NewEncoder(w).Encode(status)
	if err != nil {
		log.Printf("Error encoding health status JSON: %v", err)
	}
	log.Printf("Served %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
}
