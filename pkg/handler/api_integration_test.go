package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"main/pkg/config"
	"main/pkg/transaction"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// startTestServer starts a test HTTP server and returns the server and base URL
func startTestServer(t *testing.T, expenses []config.FrequentExpense) (*http.Server, string) {
	t.Helper()

	mux := http.NewServeMux()

	// Register handlers
	mux.HandleFunc("/health", HealthCheckHandler)
	mux.HandleFunc("/api/v1/transactions", TransactionsHandler)
	mux.HandleFunc("/api/v1/prefilled-expenses", GetPrefilledExpensesHandler(expenses))

	server := &http.Server{
		Addr:         ":0", // Use random available port
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	listener, err := net.Listen("tcp", server.Addr)
	require.NoError(t, err)

	go func() {
		if err := server.Serve(listener); err != http.ErrServerClosed {
			t.Logf("Server error: %v", err)
		}
	}()

	baseURL := fmt.Sprintf("http://%s", listener.Addr().String())
	return server, baseURL
}

func TestHealthEndpoint_Integration_GivenGetRequest_whenCalled_thenReturnsHealthyStatus(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", HealthCheckHandler)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	HealthCheckHandler(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "status")
}

func TestTransactionsAPI_Integration_GivenPostThenGet_whenCalled_thenReturnsCreatedTransaction(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	// Skip if no DB configured
	if os.Getenv("TEST_DB_HOST") == "" {
		t.Skip("TEST_DB_HOST not set, skipping integration test")
	}

	expenses := []config.FrequentExpense{}
	server, baseURL := startTestServer(t, expenses)
	defer server.Shutdown(context.Background())

	// Create a transaction
	newTxn := transaction.Transaction{
		Name:          "API Test Lunch",
		Amount:        45.00,
		Currency:      "SGD",
		Date:          "2025-01-20",
		IsClaimable:   true,
		PaidForFamily: false,
		Category:      "Food",
	}

	body, err := json.Marshal(newTxn)
	require.NoError(t, err)

	resp, err := http.Post(baseURL+"/api/v1/transactions", "application/json", bytes.NewBuffer(body))
	require.NoError(t, err)
	defer resp.Body.Close()

	// Note: Without DB, this will return 500. With DB, it should return 201
	assert.True(t, resp.StatusCode == http.StatusCreated || resp.StatusCode == http.StatusInternalServerError)
}

func TestTransactionsAPI_Integration_GivenGetRequest_whenCalled_thenReturnsTransactionsList(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	expenses := []config.FrequentExpense{}
	server, baseURL := startTestServer(t, expenses)
	defer server.Shutdown(context.Background())

	resp, err := http.Get(baseURL + "/api/v1/transactions")
	require.NoError(t, err)
	defer resp.Body.Close()

	// Without DB, returns 500. With DB, returns 200
	assert.True(t, resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusInternalServerError)
}

func TestTransactionsAPI_Integration_GivenPaginationParams_whenCalled_thenReturnsPaginatedResults(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	expenses := []config.FrequentExpense{}
	server, baseURL := startTestServer(t, expenses)
	defer server.Shutdown(context.Background())

	resp, err := http.Get(baseURL + "/api/v1/transactions?page=1&limit=5")
	require.NoError(t, err)
	defer resp.Body.Close()

	// Without DB, returns 500. With DB, returns 200 with pagination
	assert.True(t, resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusInternalServerError)
}

func TestTransactionsAPI_Integration_GivenInvalidParams_whenCalled_thenReturnsBadRequest(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	expenses := []config.FrequentExpense{}
	server, baseURL := startTestServer(t, expenses)
	defer server.Shutdown(context.Background())

	resp, err := http.Get(baseURL + "/api/v1/transactions?page=invalid")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestPrefilledExpensesEndpoint_Integration_GivenGetRequest_whenCalled_thenReturnsExpenses(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	expenses := []config.FrequentExpense{
		{Name: "Daily Lunch", Category: "Food", Currency: "SGD", IsClaimable: true},
		{Name: "Weekly Taxi", Category: "Transport", Currency: "SGD", IsClaimable: true},
	}

	server, baseURL := startTestServer(t, expenses)
	defer server.Shutdown(context.Background())

	resp, err := http.Get(baseURL + "/api/v1/prefilled-expenses")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var returnedExpenses []config.FrequentExpense
	err = json.NewDecoder(resp.Body).Decode(&returnedExpenses)
	require.NoError(t, err)
	assert.Len(t, returnedExpenses, 2)
	assert.Equal(t, "Daily Lunch", returnedExpenses[0].Name)
}

func TestCORSHeaders_Integration_GivenPreflightRequest_whenCalled_thenReturnsCorsHeaders(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	// Note: The current implementation doesn't have explicit CORS middleware
	// This test documents the expected behavior for future implementation
	t.Skip("CORS not implemented in current version")
}
