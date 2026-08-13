package handler

import (
	"encoding/json"
	"github.com/patrickmn/go-cache"
	"log"
	"main/pkg/storage"
	"main/pkg/transaction"
	"net/http"
	"strconv"
	"strings"
)

// MessageResponse represents a simple message response
type MessageResponse struct {
	Message string `json:"message"`
}

// getTransactionsHandler retrieves transactions, allowing filtering and pagination.
func getTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	cacheKey := r.URL.String() // Use the full URL as the cache key

	// Check cache first
	if cachedResponse, found := c.Get(cacheKey); found {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(cachedResponse)
		if err != nil {
			return
		}
		log.Printf("Served %s %s from cache", r.Method, r.URL.Path)
		return
	}

	// Filtering parameters
	categoryFilter := queryParams.Get("category")

	var isClaimableFilter *bool
	if claimableStr := queryParams.Get("is_claimable"); claimableStr != "" {
		val, err := strconv.ParseBool(claimableStr)
		if err != nil {
			log.Printf("Invalid boolean value for is_claimable: %s. Error: %v", claimableStr, err)
			NewErrorResponse(w, http.StatusBadRequest, "Invalid value for 'is_claimable' parameter. Use 'true' or 'false'.")
			return
		}
		isClaimableFilter = &val
	}

	var paidForFamilyFilter *bool
	if paidForFamilyStr := queryParams.Get("paid_for_family"); paidForFamilyStr != "" {
		val, err := strconv.ParseBool(paidForFamilyStr)
		if err != nil {
			log.Printf("Invalid boolean value for paid_for_family: %s. Error: %v", paidForFamilyStr, err)
			NewErrorResponse(w, http.StatusBadRequest, "Invalid value for 'paid_for_family' parameter. Use 'true' or 'false'.")
			return
		}
		paidForFamilyFilter = &val
	}

	// Pagination parameters
	pageStr := queryParams.Get("page")
	limitStr := queryParams.Get("limit")

	page := 1   // Default page
	limit := 10 // Default limit (items per page)
	var err error

	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			log.Printf("Invalid value for 'page' parameter: %s. Must be a positive integer.", pageStr)
			NewErrorResponse(w, http.StatusBadRequest, "Invalid value for 'page' parameter. Must be a positive integer.")
			return
		}
	}

	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit < 1 {
			log.Printf("Invalid value for 'limit' parameter: %s. Must be a positive integer.", limitStr)
			NewErrorResponse(w, http.StatusBadRequest, "Invalid value for 'limit' parameter. Must be a positive integer.")
			return
		}
	}

	transactions, totalItems, err := storage.GetAllTransactionsFromDB(
		categoryFilter,
		isClaimableFilter,
		paidForFamilyFilter,
		page,
		limit,
	)
	if err != nil {
		log.Printf("Error fetching transactions: %v", err)
		NewErrorResponse(w, http.StatusInternalServerError, "Internal Server Error while fetching transactions.")
		return
	}

	totalPages := 0
	if totalItems > 0 && limit > 0 {
		totalPages = (totalItems + limit - 1) / limit // Ceiling division
	}

	categories, err := storage.GetTransactionCountByCategory()
	if err != nil {
		log.Printf("Error fetching total expenses by category: %v", err)
		return
	}

	response := transaction.PaginatedTransactionsResponse{
		Categories:   categories,
		Transactions: transactions,
		CurrentPage:  page,
		PageSize:     limit,
		TotalItems:   totalItems,
		TotalPages:   totalPages,
	}

	// Store in cache
	c.Set(cacheKey, response, cache.DefaultExpiration)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		return
	}
	log.Printf("Served %s %s with %d transactions (page %d, limit %d, total %d) from %s",
		r.Method, r.URL.Path, len(transactions), page, limit, totalItems, r.RemoteAddr)
}

// createTransactionHandler handles the creation of a new transaction.
func createTransactionHandler(w http.ResponseWriter, r *http.Request) {
	var newTransaction transaction.Transaction
	if err := json.NewDecoder(r.Body).Decode(&newTransaction); err != nil {
		NewErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := storage.InsertTransaction(newTransaction); err != nil {
		log.Printf("Error inserting transaction: %v", err)
		NewErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	// Invalidate cache
	c.Flush()
	log.Println("Cache flushed due to new transaction")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response := MessageResponse{Message: "Transaction created successfully"}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response JSON: %v", err)
	}
}

func updateTransactionHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/transactions/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		NewErrorResponse(w, http.StatusBadRequest, "Invalid transaction ID")
		return
	}

	var updatedTransaction transaction.Transaction
	if err := json.NewDecoder(r.Body).Decode(&updatedTransaction); err != nil {
		NewErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := storage.UpdateTransaction(id, updatedTransaction); err != nil {
		log.Printf("Error updating transaction: %v", err)
		NewErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	// Invalidate cache
	c.Flush()
	log.Println("Cache flushed due to updated transaction")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := MessageResponse{Message: "Transaction updated successfully"}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response JSON: %v", err)
	}
}

// TransactionsHandler routes to different handlers based on the HTTP method.
func TransactionsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getTransactionsHandler(w, r)
	case http.MethodPost:
		createTransactionHandler(w, r)
	case http.MethodPut:
		updateTransactionHandler(w, r)
	default:
		NewErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}
