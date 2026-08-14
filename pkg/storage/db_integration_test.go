package storage

import (
	"context"
	"fmt"
	"os/exec"
	"testing"
	"time"

	"main/pkg/transaction"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// setupTestDatabase starts a PostgreSQL container for integration tests
func setupTestDatabase(ctx context.Context, t *testing.T) (testcontainers.Container, string) {
	t.Helper()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:15",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_PASSWORD": "test",
			"POSTGRES_DB":       "expenditure_test",
			"POSTGRES_USER":     "test",
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").
			WithOccurrence(2).
			WithStartupTimeout(60 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)

	// Get connection string
	host, err := container.Host(ctx)
	require.NoError(t, err)

	port, err := container.MappedPort(ctx, "5432")
	require.NoError(t, err)

	connStr := fmt.Sprintf("host=%s port=%s user=test password=test dbname=expenditure_test sslmode=disable",
		host, port.Port())

	return container, connStr
}

func teardownTestDatabase(ctx context.Context, container testcontainers.Container, t *testing.T) {
	t.Helper()
	err := container.Terminate(ctx)
	assert.NoError(t, err)
}

func TestInitDB_Integration_GivenPostgresContainer_whenInitialized_thenConnectsSuccessfully(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	// Check if docker is available first
	if _, err := exec.LookPath("docker"); err != nil {
		t.Skip("Docker not available, skipping integration test")
	}

	ctx := context.Background()
	container, connStr := setupTestDatabase(ctx, t)
	defer teardownTestDatabase(ctx, container, t)

	// Initialize DB with test container
	err := initDBFromConnString(connStr)
	require.NoError(t, err)
	defer CloseDB()

	// Verify connection
	db, err := GetDB()
	require.NoError(t, err)
	assert.NoError(t, db.Ping())
}

func TestInsertTransaction_Integration_GivenValidTransaction_whenInserted_thenCanBeRetrieved(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	if _, err := exec.LookPath("docker"); err != nil {
		t.Skip("Docker not available, skipping integration test")
	}

	ctx := context.Background()
	container, connStr := setupTestDatabase(ctx, t)
	defer teardownTestDatabase(ctx, container, t)

	err := initDBFromConnString(connStr)
	require.NoError(t, err)
	defer CloseDB()

	// Insert test transaction
	testTxn := transaction.Transaction{
		Name:          "Integration Test Lunch",
		Amount:        25.50,
		Currency:      "SGD",
		Date:          "2025-01-15",
		IsClaimable:   true,
		PaidForFamily: false,
		Category:      "Food",
	}

	err = InsertTransaction(testTxn)
	require.NoError(t, err)

	// Verify it can be retrieved
	txns, total, err := GetAllTransactionsFromDB("", nil, nil, 1, 10)
	require.NoError(t, err)
	assert.Equal(t, 1, total)
	require.Len(t, txns, 1)
	assert.Equal(t, "Integration Test Lunch", txns[0].Name)
	assert.Equal(t, float32(25.50), txns[0].Amount)
}

func TestGetAllTransactionsFromDB_Integration_GivenFilters_whenApplied_thenReturnsFilteredResults(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	if _, err := exec.LookPath("docker"); err != nil {
		t.Skip("Docker not available, skipping integration test")
	}

	ctx := context.Background()
	container, connStr := setupTestDatabase(ctx, t)
	defer teardownTestDatabase(ctx, container, t)

	err := initDBFromConnString(connStr)
	require.NoError(t, err)
	defer CloseDB()

	// Insert test data
	testTxns := []transaction.Transaction{
		{Name: "Lunch", Amount: 20.0, Currency: "SGD", Date: "2025-01-15", IsClaimable: true, PaidForFamily: false, Category: "Food"},
		{Name: "Family Dinner", Amount: 80.0, Currency: "SGD", Date: "2025-01-16", IsClaimable: false, PaidForFamily: true, Category: "Food"},
		{Name: "Taxi", Amount: 15.0, Currency: "SGD", Date: "2025-01-17", IsClaimable: true, PaidForFamily: false, Category: "Transport"},
	}

	for _, txn := range testTxns {
		err = InsertTransaction(txn)
		require.NoError(t, err)
	}

	// Test category filter
	txns, total, err := GetAllTransactionsFromDB("Food", nil, nil, 1, 10)
	require.NoError(t, err)
	assert.Equal(t, 2, total)
	assert.Len(t, txns, 2)

	// Test is_claimable filter
	claimable := true
	txns, total, err = GetAllTransactionsFromDB("", &claimable, nil, 1, 10)
	require.NoError(t, err)
	assert.Equal(t, 2, total)

	// Test paid_for_family filter
	paidForFamily := true
	txns, total, err = GetAllTransactionsFromDB("", nil, &paidForFamily, 1, 10)
	require.NoError(t, err)
	assert.Equal(t, 1, total)
}

func TestGetTransactionCountByCategory_Integration_GivenDataExists_whenCalled_thenReturnsCounts(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	if _, err := exec.LookPath("docker"); err != nil {
		t.Skip("Docker not available, skipping integration test")
	}

	ctx := context.Background()
	container, connStr := setupTestDatabase(ctx, t)
	defer teardownTestDatabase(ctx, container, t)

	err := initDBFromConnString(connStr)
	require.NoError(t, err)
	defer CloseDB()

	// Insert test data
	testTxns := []transaction.Transaction{
		{Name: "Lunch", Amount: 20.0, Currency: "SGD", Date: "2025-01-15", IsClaimable: true, PaidForFamily: false, Category: "Food"},
		{Name: "Dinner", Amount: 30.0, Currency: "SGD", Date: "2025-01-16", IsClaimable: true, PaidForFamily: false, Category: "Food"},
		{Name: "Taxi", Amount: 15.0, Currency: "SGD", Date: "2025-01-17", IsClaimable: true, PaidForFamily: false, Category: "Transport"},
	}

	for _, txn := range testTxns {
		err = InsertTransaction(txn)
		require.NoError(t, err)
	}

	counts, err := GetTransactionCountByCategory()
	require.NoError(t, err)
	assert.Len(t, counts, 2)

	// Find Food category
	var foodCount float32
	for _, cat := range counts {
		if cat.Name == "Food" {
			foodCount = cat.TotalAmount
		}
	}
	assert.Equal(t, float32(50.0), foodCount)
}

func TestGetTotalAmountByIsClaimable_Integration_GivenDataExists_whenCalled_thenReturnsGroupedTotals(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	if _, err := exec.LookPath("docker"); err != nil {
		t.Skip("Docker not available, skipping integration test")
	}

	ctx := context.Background()
	container, connStr := setupTestDatabase(ctx, t)
	defer teardownTestDatabase(ctx, container, t)

	err := initDBFromConnString(connStr)
	require.NoError(t, err)
	defer CloseDB()

	// Insert test data with different claimable statuses
	testTxns := []transaction.Transaction{
		{Name: "Lunch", Amount: 20.0, Currency: "SGD", Date: "2025-01-15", IsClaimable: true, PaidForFamily: false, Category: "Food"},
		{Name: "Family Dinner", Amount: 80.0, Currency: "SGD", Date: "2025-01-16", IsClaimable: false, PaidForFamily: true, Category: "Food"},
	}

	for _, txn := range testTxns {
		err = InsertTransaction(txn)
		require.NoError(t, err)
	}

	totals, err := GetTotalAmountByIsClaimable()
	require.NoError(t, err)
	assert.Len(t, totals, 2)
	assert.Equal(t, float32(20.0), totals[true])
	assert.Equal(t, float32(80.0), totals[false])
}

func TestGetTotalAmountByPaidForFamily_Integration_GivenDataExists_whenCalled_thenReturnsGroupedTotals(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	if _, err := exec.LookPath("docker"); err != nil {
		t.Skip("Docker not available, skipping integration test")
	}

	ctx := context.Background()
	container, connStr := setupTestDatabase(ctx, t)
	defer teardownTestDatabase(ctx, container, t)

	err := initDBFromConnString(connStr)
	require.NoError(t, err)
	defer CloseDB()

	// Insert test data with different paid_for_family statuses
	testTxns := []transaction.Transaction{
		{Name: "Lunch", Amount: 20.0, Currency: "SGD", Date: "2025-01-15", IsClaimable: true, PaidForFamily: false, Category: "Food"},
		{Name: "Family Dinner", Amount: 80.0, Currency: "SGD", Date: "2025-01-16", IsClaimable: false, PaidForFamily: true, Category: "Food"},
	}

	for _, txn := range testTxns {
		err = InsertTransaction(txn)
		require.NoError(t, err)
	}

	totals, err := GetTotalAmountByPaidForFamily()
	require.NoError(t, err)
	assert.Len(t, totals, 2)
	assert.Equal(t, float32(80.0), totals[true])
	assert.Equal(t, float32(20.0), totals[false])
}

// Helper functions
func initDBFromConnString(connStr string) error {
	// Parse connection string and initialize DB
	// This is a simplified version - in production you'd parse the connStr properly
	return InitDBWithConnString(connStr)
}
