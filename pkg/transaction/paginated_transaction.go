package transaction

// PaginatedTransactionsResponse defines the structure for paginated transaction results.
type PaginatedTransactionsResponse struct {
	Transactions []Transaction `json:"transactions"`
	CurrentPage  int           `json:"currentPage"`
	PageSize     int           `json:"pageSize"`
	TotalItems   int           `json:"totalItems"`
	TotalPages   int           `json:"totalPages"`
}
