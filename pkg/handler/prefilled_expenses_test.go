package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"main/pkg/config"

	"github.com/stretchr/testify/assert"
)

func TestGetPrefilledExpensesHandler_GivenGetRequest_whenCalled_thenReturnsExpensesAsJson(t *testing.T) {
	expenses := []config.FrequentExpense{
		{Name: "Lunch", Category: "Food", Currency: "SGD", IsClaimable: true, PaidForFamily: false},
		{Name: "Taxi", Category: "Transport", Currency: "SGD", IsClaimable: true, PaidForFamily: false},
	}

	handler := GetPrefilledExpensesHandler(expenses)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/prefilled-expenses", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Lunch")
	assert.Contains(t, w.Body.String(), "Taxi")
}

func TestGetPrefilledExpensesHandler_GivenNonGetMethod_whenCalled_thenReturnsMethodNotAllowed(t *testing.T) {
	expenses := []config.FrequentExpense{}
	handler := GetPrefilledExpensesHandler(expenses)

	methods := []string{http.MethodPost, http.MethodPut, http.MethodDelete}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/api/v1/prefilled-expenses", nil)
			w := httptest.NewRecorder()

			handler(w, req)

			assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
		})
	}
}

func TestGetPrefilledExpensesHandler_GivenCachedResponse_whenCalled_thenServesFromCache(t *testing.T) {
	expenses := []config.FrequentExpense{
		{Name: "Coffee", Category: "Food", Currency: "USD"},
	}

	handler := GetPrefilledExpensesHandler(expenses)

	// First request
	req1 := httptest.NewRequest(http.MethodGet, "/api/v1/prefilled-expenses", nil)
	w1 := httptest.NewRecorder()
	handler(w1, req1)

	// Second request (should be served from cache)
	req2 := httptest.NewRequest(http.MethodGet, "/api/v1/prefilled-expenses", nil)
	w2 := httptest.NewRecorder()
	handler(w2, req2)

	assert.Equal(t, w1.Body.String(), w2.Body.String())
}
