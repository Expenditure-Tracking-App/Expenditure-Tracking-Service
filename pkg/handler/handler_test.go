package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHealthCheckHandler_GivenGetRequest_whenDatabaseConnected_thenReturnsOkStatus(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	HealthCheckHandler(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "ok")
}

func TestHealthCheckHandler_GivenNonGetMethod_whenCalled_thenReturnsMethodNotAllowed(t *testing.T) {
	methods := []string{http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/health", nil)
			w := httptest.NewRecorder()

			HealthCheckHandler(w, req)

			assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
			assert.Contains(t, w.Body.String(), "Method not allowed")
		})
	}
}

func TestNewErrorResponse_GivenStatusAndMessage_whenEncoded_thenReturnsJsonError(t *testing.T) {
	w := httptest.NewRecorder()

	NewErrorResponse(w, http.StatusNotFound, "Resource not found")

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "error")
	assert.Contains(t, w.Body.String(), "Resource not found")
}

func TestHealthStatus_GivenHealthyDb_whenSerialized_thenContainsTimestamp(t *testing.T) {
	status := HealthStatus{
		Status:     "ok",
		Timestamp:  time.Now().UTC(),
		Database:   "connected",
		DatabaseAt: time.Now().UTC(),
	}

	assert.Equal(t, "ok", status.Status)
	assert.Equal(t, "connected", status.Database)
	assert.False(t, status.Timestamp.IsZero())
	assert.False(t, status.DatabaseAt.IsZero())
}
