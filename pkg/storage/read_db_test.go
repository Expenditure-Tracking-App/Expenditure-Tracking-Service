package storage

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetAllTransactionsFromDB_GivenNoFilters_whenCalled_thenReturnsAllTransactions(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer dbMock.Close()

	oldDB := db
	db = dbMock
	defer func() { db = oldDB }()

	// Mock the count query
	countRows := sqlmock.NewRows([]string{"count"}).AddRow(5)
	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM transactions`).WillReturnRows(countRows)

	// Mock the select query
	selectRows := sqlmock.NewRows([]string{"id", "name", "amount", "currency", "date", "is_claimable", "paid_for_family", "category", "created_at"}).
		AddRow(1, "Lunch", 50.0, "SGD", "2025-01-15", true, false, "Food", time.Now())

	mock.ExpectQuery(`SELECT id, name, amount, currency, date, is_claimable, paid_for_family, category, created_at`).
		WithArgs(10, 0).
		WillReturnRows(selectRows)

	transactions, total, err := GetAllTransactionsFromDB("", nil, nil, 1, 10)

	assert.NoError(t, err)
	assert.Equal(t, 5, total)
	assert.Len(t, transactions, 1)
	assert.Equal(t, "Lunch", transactions[0].Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAllTransactionsFromDB_GivenCategoryFilter_whenCalled_thenFiltersByCategory(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer dbMock.Close()

	oldDB := db
	db = dbMock
	defer func() { db = oldDB }()

	// Mock the count query with category filter
	countRows := sqlmock.NewRows([]string{"count"}).AddRow(2)
	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM transactions WHERE category =`).
		WithArgs("Food").
		WillReturnRows(countRows)

	// Mock the select query with category filter
	selectRows := sqlmock.NewRows([]string{"id", "name", "amount", "currency", "date", "is_claimable", "paid_for_family", "category", "created_at"}).
		AddRow(1, "Lunch", 50.0, "SGD", "2025-01-15", true, false, "Food", time.Now())

	mock.ExpectQuery(`SELECT id, name, amount, currency, date, is_claimable, paid_for_family, category, created_at`).
		WithArgs("Food", 10, 0).
		WillReturnRows(selectRows)

	transactions, total, err := GetAllTransactionsFromDB("Food", nil, nil, 1, 10)

	assert.NoError(t, err)
	assert.Equal(t, 2, total)
	assert.Len(t, transactions, 1)
	assert.Equal(t, "Food", transactions[0].Category)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAllTransactionsFromDB_GivenIsClaimableFilter_whenCalled_thenFiltersByClaimable(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer dbMock.Close()

	oldDB := db
	db = dbMock
	defer func() { db = oldDB }()

	claimable := true

	// Mock the count query with is_claimable filter
	countRows := sqlmock.NewRows([]string{"count"}).AddRow(3)
	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM transactions WHERE is_claimable =`).
		WithArgs(true).
		WillReturnRows(countRows)

	// Mock the select query with is_claimable filter
	selectRows := sqlmock.NewRows([]string{"id", "name", "amount", "currency", "date", "is_claimable", "paid_for_family", "category", "created_at"}).
		AddRow(1, "Lunch", 50.0, "SGD", "2025-01-15", true, false, "Food", time.Now())

	mock.ExpectQuery(`SELECT id, name, amount, currency, date, is_claimable, paid_for_family, category, created_at`).
		WithArgs(true, 10, 0).
		WillReturnRows(selectRows)

	transactions, total, err := GetAllTransactionsFromDB("", &claimable, nil, 1, 10)

	assert.NoError(t, err)
	assert.Equal(t, 3, total)
	assert.Len(t, transactions, 1)
	assert.True(t, transactions[0].IsClaimable)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAllTransactionsFromDB_GivenPaidForFamilyFilter_whenCalled_thenFiltersByPaidForFamily(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer dbMock.Close()

	oldDB := db
	db = dbMock
	defer func() { db = oldDB }()

	paidForFamily := false

	// Mock the count query with paid_for_family filter
	countRows := sqlmock.NewRows([]string{"count"}).AddRow(4)
	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM transactions WHERE paid_for_family =`).
		WithArgs(false).
		WillReturnRows(countRows)

	// Mock the select query with paid_for_family filter
	selectRows := sqlmock.NewRows([]string{"id", "name", "amount", "currency", "date", "is_claimable", "paid_for_family", "category", "created_at"}).
		AddRow(1, "Lunch", 50.0, "SGD", "2025-01-15", true, false, "Food", time.Now())

	mock.ExpectQuery(`SELECT id, name, amount, currency, date, is_claimable, paid_for_family, category, created_at`).
		WithArgs(false, 10, 0).
		WillReturnRows(selectRows)

	transactions, total, err := GetAllTransactionsFromDB("", nil, &paidForFamily, 1, 10)

	assert.NoError(t, err)
	assert.Equal(t, 4, total)
	assert.Len(t, transactions, 1)
	assert.False(t, transactions[0].PaidForFamily)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAllTransactionsFromDB_GivenNoMatchingResults_whenCalled_thenReturnsEmptySlice(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer dbMock.Close()

	oldDB := db
	db = dbMock
	defer func() { db = oldDB }()

	// Mock the count query with zero results
	countRows := sqlmock.NewRows([]string{"count"}).AddRow(0)
	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM transactions`).WillReturnRows(countRows)

	transactions, total, err := GetAllTransactionsFromDB("", nil, nil, 1, 10)

	assert.NoError(t, err)
	assert.Equal(t, 0, total)
	assert.Empty(t, transactions)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetTransactionCountByCategory_GivenDataExists_whenCalled_thenReturnsCategoryCounts(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer dbMock.Close()

	oldDB := db
	db = dbMock
	defer func() { db = oldDB }()

	// Mock the category count query
	rows := sqlmock.NewRows([]string{"category", "total_cost"}).
		AddRow("Food", 150.50).
		AddRow("Transport", 75.00)

	mock.ExpectQuery(`SELECT category, SUM\(amount\) AS total_cost FROM transactions GROUP BY category`).
		WillReturnRows(rows)

	counts, err := GetTransactionCountByCategory()

	assert.NoError(t, err)
	assert.Len(t, counts, 2)
	assert.Equal(t, "Food", counts[0].Name)
	assert.Equal(t, float32(150.50), counts[0].TotalAmount)
	assert.Equal(t, "Transport", counts[1].Name)
	assert.Equal(t, float32(75.00), counts[1].TotalAmount)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetTotalAmountByIsClaimable_GivenDataExists_whenCalled_thenReturnsGroupedTotals(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer dbMock.Close()

	oldDB := db
	db = dbMock
	defer func() { db = oldDB }()

	// Mock the is_claimable grouping query
	rows := sqlmock.NewRows([]string{"is_claimable", "total_amount"}).
		AddRow(true, 200.00).
		AddRow(false, 100.00)

	mock.ExpectQuery(`SELECT is_claimable, SUM\(amount\) AS total_amount FROM transactions GROUP BY is_claimable`).
		WillReturnRows(rows)

	totals, err := GetTotalAmountByIsClaimable()

	assert.NoError(t, err)
	assert.Len(t, totals, 2)
	assert.Equal(t, float32(200.00), totals[true])
	assert.Equal(t, float32(100.00), totals[false])
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetTotalAmountByPaidForFamily_GivenDataExists_whenCalled_thenReturnsGroupedTotals(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer dbMock.Close()

	oldDB := db
	db = dbMock
	defer func() { db = oldDB }()

	// Mock the paid_for_family grouping query
	rows := sqlmock.NewRows([]string{"paid_for_family", "total_amount"}).
		AddRow(true, 80.00).
		AddRow(false, 250.00)

	mock.ExpectQuery(`SELECT paid_for_family, SUM\(amount\) AS total_amount FROM transactions GROUP BY paid_for_family`).
		WillReturnRows(rows)

	totals, err := GetTotalAmountByPaidForFamily()

	assert.NoError(t, err)
	assert.Len(t, totals, 2)
	assert.Equal(t, float32(80.00), totals[true])
	assert.Equal(t, float32(250.00), totals[false])
	assert.NoError(t, mock.ExpectationsWereMet())
}
