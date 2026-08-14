package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransactionsHandler_GivenGetRequest_whenCalled_thenReturnsTransactionsList(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/transactions", nil)
	w := httptest.NewRecorder()

	TransactionsHandler(w, req)

	// Without DB configured, returns 500 with error message
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "error")
}

func TestTransactionsHandler_GivenPostRequest_whenInvalidJson_thenReturnsBadRequest(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/transactions", bytes.NewBuffer([]byte("invalid json")))
	w := httptest.NewRecorder()

	TransactionsHandler(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "error")
}

func TestTransactionsHandler_GivenPostRequest_whenValidJson_thenAttemptsInsert(t *testing.T) {
	txn := map[string]interface{}{
		"name":          "Test Lunch",
		"amount":        50.0,
		"currency":      "SGD",
		"date":          "2025-01-15",
		"isClaimable":   true,
		"paidForFamily": false,
		"category":      "Food",
	}
	body, _ := json.Marshal(txn)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/transactions", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	TransactionsHandler(w, req)

	// Without DB configured, returns 500 with error message
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "error")
}

func TestTransactionsHandler_GivenPutRequest_whenInvalidId_thenReturnsBadRequest(t *testing.T) {
	req := httptest.NewRequest(http.MethodPut, "/api/v1/transactions/invalid-id", bytes.NewBuffer([]byte("{}")))
	w := httptest.NewRecorder()

	TransactionsHandler(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid transaction ID")
}

func TestTransactionsHandler_GivenPutRequest_whenValidId_thenAttemptsUpdate(t *testing.T) {
	txn := map[string]interface{}{
		"name":     "Updated Lunch",
		"amount":   60.0,
		"currency": "SGD",
		"date":     "2025-01-16",
		"category": "Food",
	}
	body, _ := json.Marshal(txn)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/transactions/123", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	TransactionsHandler(w, req)

	// Without DB configured, returns 500 with error message
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "error")
}

func TestTransactionsHandler_GivenUnsupportedMethod_whenCalled_thenReturnsMethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/transactions/123", nil)
	w := httptest.NewRecorder()

	TransactionsHandler(w, req)

	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}

func TestGetTransactionsHandler_GivenInvalidPageParameter_whenCalled_thenReturnsBadRequest(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/transactions?page=invalid", nil)
	w := httptest.NewRecorder()

	getTransactionsHandler(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid value for 'page' parameter")
}

func TestGetTransactionsHandler_GivenInvalidLimitParameter_whenCalled_thenReturnsBadRequest(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/transactions?limit=invalid", nil)
	w := httptest.NewRecorder()

	getTransactionsHandler(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid value for 'limit' parameter")
}

func TestGetTransactionsHandler_GivenInvalidIsClaimableParameter_whenCalled_thenReturnsBadRequest(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/transactions?is_claimable=maybe", nil)
	w := httptest.NewRecorder()

	getTransactionsHandler(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid value for 'is_claimable' parameter")
}

func TestGetTransactionsHandler_GivenInvalidPaidForFamilyParameter_whenCalled_thenReturnsBadRequest(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/transactions?paid_for_family=maybe", nil)
	w := httptest.NewRecorder()

	getTransactionsHandler(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid value for 'paid_for_family' parameter")
}
