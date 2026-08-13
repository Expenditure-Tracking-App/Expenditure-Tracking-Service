# Test Plan - Expenditure Tracking Bot

This document outlines the comprehensive testing strategy for the Expenditure Tracking Bot application, covering unit tests, integration tests, and the use of mocks.

---

## Table of Contents

1. [Overview](#overview)
2. [Testing Tools & Setup](#testing-tools--setup)
3. [Unit Tests](#unit-tests)
4. [Integration Tests](#integration-tests)
5. [Mock Strategy](#mock-strategy)
6. [Test File Structure](#test-file-structure)
7. [Priority Matrix](#priority-matrix)

---

## Overview

### Application Architecture

The application consists of two entry points sharing common packages:

```
cmd/server/          cmd/bot/
    │                    │
    └────────┬───────────┘
             │
    ┌────────┴────────┐
    │  pkg/handler/   │◄── HTTP handlers (REST API)
    │  pkg/bot/       │◄── Telegram bot logic
    │  pkg/storage/   │◄── Database/file persistence
    │  pkg/session/   │◄── Bot conversation state
    │  pkg/transaction│◄── Domain models
    └────────┬────────┘
             │
    ┌────────┴────────┐
    │  pkg/config/    │◄── Configuration loading
    └─────────────────┘
```

### External Dependencies

| Dependency | Package | Usage |
|------------|---------|-------|
| PostgreSQL | `pkg/storage` | Transaction persistence |
| Telegram Bot API | `pkg/bot` | Bot communication |
| File system | `pkg/storage` | Fallback storage |
| HTTP Server | `cmd/server` | REST API |

---

## Testing Tools & Setup

### Recommended Tools

| Tool | Purpose | Installation |
|------|---------|--------------|
| `testing` (stdlib) | Base testing framework | Built-in |
| `testify/assert` | Assertions | `go get github.com/stretchr/testify` |
| `testify/mock` | Mock utilities | `go get github.com/stretchr/testify` |
| `mockery` | Interface mock generation | `go install github.com/vektra/mockery/v2@latest` |
| `testcontainers-go` | Containerized integration tests | `go get github.com/testcontainers/testcontainers-go` |
| `sqlmock` | SQL mock expectations | `go get github.com/DATA-DOG/go-sqlmock` |

### Mockery Setup

Create interfaces for external dependencies, then generate mocks:

```bash
# Install mockery
go install github.com/vektra/mockery/v2@latest

# Generate mocks for all interfaces
mockery --dir=./pkg/storage --name=DatabaseQuerier --output=./pkg/storage/mocks
mockery --dir=./pkg/bot --name=TelegramAPI --output=./pkg/bot/mocks
```

### Test Configuration

Create `pkg/config/test_config.go`:

```go
package config

// TestConfig returns a configuration suitable for tests
func TestConfig() Config {
    return Config{
        Server: ServerConfig{Port: 0}, // Use random port
        FeaturesConfig: FeaturesConfig{
            SaveToDB: false, // Use file storage for unit tests
        },
        Database: DatabaseConfig{
            Host:     "localhost",
            Port:     5432,
            User:     "test",
            Password: "test",
            DBName:   "expenditure_test",
            SSLMode:  "disable",
        },
        ExpenseCategories:   []string{"Food", "Transport", "Utilities"},
        FrequentExpenses:    []FrequentExpense{},
        SupportedCurrencies: []string{"USD", "EUR", "SGD"},
    }
}
```

---

## Unit Tests

### 1. `pkg/transaction/`

**File**: `pkg/transaction/transaction_test.go`

| Function | Test Cases | Priority |
|----------|-----------|----------|
| `ValidateAmount` | - Valid integer ("100")<br>- Valid decimal ("99.99")<br>- Invalid string ("abc")<br>- Negative value<br>- Empty string | P0 |
| `ValidateBool` | - "true", "false"<br>- "yes", "no" (if supported)<br>- Invalid input ("maybe")<br>- Empty string | P0 |
| `ProcessDate` | - "t" (today)<br>- Valid DD.MM.YY ("15.04.25")<br>- Invalid format<br>- Edge cases (29.02.24) | P0 |

```go
func TestValidateAmount(t *testing.T) {
    tests := []struct {
        name        string
        input       string
        want        float32
        wantErr     bool
        errContains string
    }{
        {"valid integer", "100", 100.0, false, ""},
        {"valid decimal", "99.99", 99.99, false, ""},
        {"invalid string", "abc", 0, true, "invalid amount"},
        {"negative", "-50", -50.0, false, ""},
        {"empty", "", 0, true, "invalid amount"},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := ValidateAmount(tt.input)
            if tt.wantErr {
                assert.Error(t, err)
                if tt.errContains != "" {
                    assert.Contains(t, err.Error(), tt.errContains)
                }
            } else {
                assert.NoError(t, err)
                assert.Equal(t, tt.want, got)
            }
        })
    }
}
```

### 2. `pkg/session/`

**File**: `pkg/session/session_test.go`

| Function | Test Cases | Priority |
|----------|-----------|----------|
| `NewUserSession` | - Returns session with QuestionName as current<br>- Answers initialized to zero values | P0 |
| `IsSessionComplete` | - Returns false for questions < QuestionCount<br>- Returns true for QuestionCount | P0 |
| `HandleAnswer` | - Each question type with valid input<br>- Each question type with invalid input<br>- Pre-filled expense auto-advance | P0 |

**File**: `pkg/session/handle_answer_test.go`

```go
func TestHandleAnswer_ValidAmount(t *testing.T) {
    session := NewUserSession()
    session.CurrentQuestion = QuestionAmount
    
    err := session.HandleAnswer("50.00")
    
    assert.NoError(t, err)
    assert.Equal(t, float32(50.00), session.Answers.Amount)
    assert.Equal(t, QuestionCurrency, session.CurrentQuestion)
}

func TestHandleAnswer_InvalidAmount(t *testing.T) {
    session := NewUserSession()
    session.CurrentQuestion = QuestionAmount
    
    err := session.HandleAnswer("not-a-number")
    
    assert.Error(t, err)
    assert.Equal(t, QuestionAmount, session.CurrentQuestion) // Should not advance
}
```

**File**: `pkg/session/constants_test.go`

| Test | Description |
|------|-------------|
| `TestQuestionCount` | Verify QuestionCount equals length of Questions array |
| `TestQuestionOrder` | Verify iota-based constants are sequential |

### 3. `pkg/config/`

**File**: `pkg/config/config_test.go`

| Function | Test Cases | Priority |
|----------|-----------|----------|
| `LoadConfig` | - Valid config file<br>- Missing file<br>- Invalid YAML<br>- Missing required fields | P1 |

### 4. `pkg/bot/`

**File**: `pkg/bot/bot_test.go`

Requires mocking Telegram API. Create interface first:

```go
// pkg/bot/telegram_api.go
type TelegramAPI interface {
    Send(msg tgbotapi.MessageConfig) (tgbotapi.Message, error)
    Request(chattable tgbotapi.Chattable) (*tgbotapi.APIResponse, error)
    GetUpdatesChan(config tgbotapi.UpdateConfig) tgbotapi.UpdatesChannel
}
```

| Function | Test Cases | Priority |
|----------|-----------|----------|
| `NewBot` | - Valid token<br>- Invalid token<br>- Empty token | P1 |
| `handleTextMessage` | - /add command<br>- /summary command<br>- Answer during session<br>- Unknown message (no session) | P0 |
| `handleCallbackQuery` | - Yes/No for boolean questions<br>- Category selection<br>- Currency selection<br>- Expired session | P0 |
| `startSession` | - Creates new session<br>- Sends first question | P1 |
| `completeSession` | - SaveToDB=true path<br>- SaveToDB=false path<br>- Save failure handling | P1 |

### 5. `pkg/handler/`

**File**: `pkg/handler/handler_test.go`

| Function | Test Cases | Priority |
|----------|-----------|----------|
| `HealthCheckHandler` | - GET request, DB connected<br>- GET request, DB disconnected<br>- Non-GET method | P0 |
| `NewErrorResponse` | - Various status codes | P2 |

**File**: `pkg/handler/transaction_test.go`

| Function | Test Cases | Priority |
|----------|-----------|----------|
| `getTransactionsHandler` | - No filters<br>- Category filter<br>- is_claimable filter<br>- paid_for_family filter<br>- Pagination<br>- Invalid parameters<br>- Cache hit | P0 |
| `createTransactionHandler` | - Valid transaction<br>- Invalid JSON<br>- Missing required fields | P0 |
| `updateTransactionHandler` | - Valid update<br>- Invalid ID<br>- Invalid JSON<br>- Non-existent ID | P1 |
| `TransactionsHandler` | - GET, POST, PUT, unsupported method | P0 |

**File**: `pkg/handler/prefilled_expenses_test.go`

| Test | Description |
|------|-------------|
| `GetPrefilledExpensesHandler` | Returns configured expenses as JSON |
| `GetPrefilledExpensesHandler_MethodNotAllowed` | Non-GET returns 405 |

### 6. `pkg/storage/`

**File**: `pkg/storage/storage_test.go`

| Function | Test Cases | Priority |
|----------|-----------|----------|
| `SaveResponseToFile` | - Valid transaction<br>- File permission error | P1 |
| `InsertTransaction` | - Valid insert (with sqlmock)<br>- DB connection error | P0 |
| `UpdateTransaction` | - Valid update<br>- Non-existent ID<br>- DB connection error | P1 |

---

## Integration Tests

### 1. Database Integration Tests

**File**: `pkg/storage/db_integration_test.go`

Use testcontainers-go for PostgreSQL:

```go
func TestGetAllTransactionsFromDB_Integration(t *testing.T) {
    // Skip if not running integration tests
    if testing.Short() {
        t.Skip("skipping integration test")
    }
    
    // Start PostgreSQL container
    ctx := context.Background()
    container, connStr := setupTestDatabase(ctx, t)
    defer teardownTestDatabase(ctx, container, t)
    
    // Initialize storage with test DB
    err := InitDBWithConnString(connStr)
    require.NoError(t, err)
    defer CloseDB()
    
    // Insert test data
    testTxn := transaction.Transaction{
        Name:          "Test Lunch",
        Amount:        25.50,
        Currency:      "SGD",
        Category:      "Food",
        IsClaimable:   true,
        PaidForFamily: false,
    }
    err = InsertTransaction(testTxn)
    require.NoError(t, err)
    
    // Test retrieval
    txns, total, err := GetAllTransactionsFromDB("", nil, nil, 1, 10)
    require.NoError(t, err)
    assert.Equal(t, 1, total)
    assert.Len(t, txns, 1)
    assert.Equal(t, "Test Lunch", txns[0].Name)
}
```

| Test | Description |
|------|-------------|
| `TestInitDB` | Connection success/failure |
| `TestCreateTableIfNotExists` | Table creation on init |
| `TestInsertTransaction_Integration` | Full insert cycle |
| `TestGetAllTransactionsFromDB_Integration` | Query with various filters |
| `TestGetTransactionCountByCategory_Integration` | Aggregation query |
| `TestGetTotalAmountByIsClaimable_Integration` | Boolean grouping |
| `TestGetTotalAmountByPaidForFamily_Integration` | Boolean grouping |

### 2. HTTP API Integration Tests

**File**: `pkg/handler/api_integration_test.go`

```go
func TestTransactionsAPI_Integration(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test")
    }
    
    // Start test server on random port
    server := startTestServer(t)
    defer server.Shutdown(context.Background())
    
    baseURL := fmt.Sprintf("http://%s", server.Addr)
    
    // Test GET /api/v1/transactions
    resp, err := http.Get(baseURL + "/api/v1/transactions")
    require.NoError(t, err)
    assert.Equal(t, http.StatusOK, resp.StatusCode)
    
    // Test POST /api/v1/transactions
    txn := transaction.Transaction{Name: "API Test", Amount: 100}
    body, _ := json.Marshal(txn)
    resp, err = http.Post(baseURL+"/api/v1/transactions", "application/json", bytes.NewReader(body))
    require.NoError(t, err)
    assert.Equal(t, http.StatusCreated, resp.StatusCode)
}
```

| Test | Description |
|------|-------------|
| `TestHealthEndpoint` | GET /health returns status |
| `TestTransactionsAPI_CRUD` | Full CRUD cycle |
| `TestTransactionsAPI_Filtering` | Query parameter filtering |
| `TestTransactionsAPI_Pagination` | Page/limit parameters |
| `TestTransactionsAPI_Caching` | Cache invalidation on write |
| `TestPrefilledExpensesEndpoint` | GET /api/v1/prefilled-expenses |
| `TestCORSHeaders` | CORS preflight and actual requests |

### 3. Bot Integration Tests

**File**: `pkg/bot/bot_integration_test.go`

Note: Telegram Bot API integration tests require a test bot token. These can be run conditionally:

```go
func TestBotEndToEnd(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test")
    }
    
    token := os.Getenv("TEST_TELEGRAM_BOT_TOKEN")
    if token == "" {
        t.Skip("TEST_TELEGRAM_BOT_TOKEN not set")
    }
    
    // Create bot with file storage (no DB dependency)
    cfg := config.TestConfig()
    cfg.FeaturesConfig.SaveToDB = false
    
    bot, err := NewBot(token, cfg.FeaturesConfig, cfg.FrequentExpenses, 
                       cfg.ExpenseCategories, cfg.SupportedCurrencies)
    require.NoError(t, err)
    
    // Test would require simulating Telegram updates
    // This is complex and may be better suited for manual testing
}
```

---

## Mock Strategy

### Interfaces to Create

#### 1. Storage Layer (`pkg/storage/interfaces.go`)

```go
// DatabaseQuerier abstracts database operations for mocking
type DatabaseQuerier interface {
    Query(query string, args ...interface{}) (*sql.Rows, error)
    QueryRow(query string, args ...interface{}) *sql.Row
    Exec(query string, args ...interface{}) (sql.Result, error)
    Ping() error
    Close() error
}

// TransactionRepository abstracts transaction persistence
type TransactionRepository interface {
    Insert(t transaction.Transaction) error
    Update(id int, t transaction.Transaction) error
    GetAll(category string, isClaimable *bool, paidForFamily *bool, page, limit int) ([]transaction.Transaction, int, error)
    GetCategoryCounts() ([]transaction.Category, error)
    GetTotalByIsClaimable() (map[bool]float32, error)
    GetTotalByPaidForFamily() (map[bool]float32, error)
}
```

Generate mock with mockery:

```bash
mockery --dir=./pkg/storage --name=DatabaseQuerier --output=./pkg/storage/mocks --outpkg=mocks
mockery --dir=./pkg/storage --name=TransactionRepository --output=./pkg/storage/mocks --outpkg=mocks
```

#### 2. Bot Layer (`pkg/bot/interfaces.go`)

```go
// TelegramSender abstracts Telegram message sending
type TelegramSender interface {
    Send(msg tgbotapi.MessageConfig) (tgbotapi.Message, error)
    Request(chattable tgbotapi.Chattable) (*tgbotapi.APIResponse, error)
}

// SessionStore manages user sessions
type SessionStore interface {
    Get(chatID int64) (*session.UserSession, bool)
    Set(chatID int64, s *session.UserSession)
    Delete(chatID int64)
}
```

#### 3. HTTP Handler Layer

For HTTP handlers, use `httptest.ResponseRecorder` and `http.NewRequest`:

```go
func TestGetTransactionsHandler(t *testing.T) {
    req := httptest.NewRequest(http.MethodGet, "/api/v1/transactions?page=1&limit=10", nil)
    w := httptest.NewRecorder()
    
    // Use mocked storage
    mockRepo := &mocks.MockTransactionRepository{}
    mockRepo.On("GetAll", "", (*bool)(nil), (*bool)(nil), 1, 10).
        Return([]transaction.Transaction{}, 0, nil)
    
    handler := getTransactionsHandlerWithRepo(mockRepo)
    handler.ServeHTTP(w, req)
    
    assert.Equal(t, http.StatusOK, w.Code)
}
```

### Using sqlmock for Database Tests

```go
import (
    "github.com/DATA-DOG/go-sqlmock"
)

func TestInsertTransaction(t *testing.T) {
    db, mock, err := sqlmock.New()
    require.NoError(t, err)
    defer db.Close()
    
    // Set global db for the package (may require refactoring)
    SetDBForTest(db)
    
    // Expect the INSERT query
    mock.ExpectExec(`INSERT INTO transactions`).
        WithArgs("Test", 100.0, "SGD", sqlmock.AnyArg(), true, false, "Food").
        WillReturnResult(sqlmock.NewResult(1, 1))
    
    txn := transaction.Transaction{
        Name: "Test", Amount: 100, Currency: "SGD",
        IsClaimable: true, PaidForFamily: false, Category: "Food",
    }
    
    err = InsertTransaction(txn)
    assert.NoError(t, err)
    
    // Verify all expectations were met
    assert.NoError(t, mock.ExpectationsWereMet())
}
```

---

## Test File Structure

```
pkg/
├── config/
│   ├── config.go
│   ├── config_test.go           # Unit: YAML parsing
│   └── test_config.go           # Test helpers
│
├── transaction/
│   ├── transaction.go
│   ├── transaction_test.go      # Unit: ValidateAmount, ValidateBool, ProcessDate
│   └── transaction_categories.go
│
├── session/
│   ├── session.go
│   ├── session_test.go          # Unit: NewUserSession, IsSessionComplete
│   ├── handle_answer.go
│   ├── handle_answer_test.go    # Unit: HandleAnswer for each question
│   ├── constants.go
│   └── constants_test.go        # Unit: QuestionCount validation
│
├── storage/
│   ├── db.go
│   ├── db_test.go               # Unit: InitDB error cases
│   ├── db_integration_test.go   # Integration: PostgreSQL tests
│   ├── storage.go
│   ├── storage_test.go          # Unit: SaveResponseToFile
│   ├── read_db.go
│   ├── read_db_test.go          # Unit: Query building (with sqlmock)
│   ├── interfaces.go            # Interfaces for mocking
│   └── mocks/
│       └── mock_querier.go      # Generated by mockery
│
├── bot/
│   ├── bot.go
│   ├── bot_test.go              # Unit: Message handling (with mocks)
│   ├── bot_integration_test.go  # Integration: Real Telegram API
│   └── interfaces.go            # TelegramAPI interface
│
└── handler/
    ├── handler.go
    ├── handler_test.go          # Unit: HealthCheckHandler
    ├── transaction.go
    ├── transaction_test.go      # Unit: CRUD handlers
    ├── api_integration_test.go  # Integration: HTTP server
    └── prefilled_expenses.go
    └── prefilled_expenses_test.go
```

---

## Priority Matrix

### P0 - Critical Path (Write First)

| Package | Test File | Rationale |
|---------|-----------|-----------|
| `transaction` | `transaction_test.go` | Core validation logic |
| `session` | `handle_answer_test.go` | Bot conversation flow |
| `handler` | `transaction_test.go` | API correctness |
| `storage` | `read_db_test.go` | Data access with sqlmock |

### P1 - Important Features

| Package | Test File | Rationale |
|---------|-----------|-----------|
| `bot` | `bot_test.go` | Telegram interaction |
| `handler` | `handler_test.go` | Health checks, error handling |
| `storage` | `db_integration_test.go` | End-to-end DB tests |
| `session` | `session_test.go` | Session lifecycle |

### P2 - Edge Cases & Completeness

| Package | Test File | Rationale |
|---------|-----------|-----------|
| `config` | `config_test.go` | Configuration loading |
| `storage` | `storage_test.go` | File fallback storage |
| `handler` | `prefilled_expenses_test.go` | Config endpoint |
| `bot` | `bot_integration_test.go` | Real API testing |

---

## Implementation Checklist

### Phase 1: Foundation (P0)

- [ ] Create `pkg/config/test_config.go`
- [ ] Create `pkg/storage/interfaces.go`
- [ ] Create `pkg/bot/interfaces.go`
- [ ] Install testify: `go get github.com/stretchr/testify`
- [ ] Install sqlmock: `go get github.com/DATA-DOG/go-sqlmock`
- [ ] Write `pkg/transaction/transaction_test.go`
- [ ] Write `pkg/session/handle_answer_test.go`

### Phase 2: Core Testing (P0 continued)

- [ ] Write `pkg/handler/transaction_test.go`
- [ ] Write `pkg/storage/read_db_test.go` (with sqlmock)
- [ ] Generate mocks with mockery
- [ ] Write `pkg/bot/bot_test.go` (with mocks)

### Phase 3: Integration (P1)

- [ ] Write `pkg/storage/db_integration_test.go` (testcontainers)
- [ ] Write `pkg/handler/api_integration_test.go`
- [ ] Write `pkg/session/session_test.go`
- [ ] Write `pkg/handler/handler_test.go`

### Phase 4: Completeness (P2)

- [ ] Write `pkg/config/config_test.go`
- [ ] Write `pkg/storage/storage_test.go`
- [ ] Write `pkg/handler/prefilled_expenses_test.go`
- [ ] Write `pkg/bot/bot_integration_test.go`

---

## Running Tests

```bash
# Run all unit tests (fast, no external dependencies)
go test ./...

# Run all tests including integration (requires Docker)
go test ./... -tags=integration

# Run specific package tests
go test ./pkg/transaction/...
go test ./pkg/session/...

# Run with coverage
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out

# Run with verbose output
go test ./... -v

# Run specific test function
go test ./pkg/transaction -run TestValidateAmount

# Skip integration tests
go test ./... -short
```

---

## CI/CD Integration

Add to `.github/workflows/test.yml`:

```yaml
name: Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    
    services:
      postgres:
        image: postgres:15
        env:
          POSTGRES_PASSWORD: test
          POSTGRES_DB: expenditure_test
        ports:
          - 5432:5432
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
    
    steps:
      - uses: actions/checkout@v4
      
      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.26'
      
      - name: Install dependencies
        run: go mod download
      
      - name: Run unit tests
        run: go test ./... -short -coverprofile=coverage.out
      
      - name: Run integration tests
        run: go test ./... -tags=integration -v
      
      - name: Upload coverage
        uses: codecov/codecov-action@v4
        with:
          file: ./coverage.out
```

---

## Notes

1. **Mock Generation**: After creating any new interface, run mockery to generate/update mocks:
   ```bash
   mockery --dir=./pkg/<package> --name=<InterfaceName> --output=./pkg/<package>/mocks
   ```

2. **Test Data**: Consider creating a `testdata/` directory for:
   - Sample YAML configurations
   - Expected JSON responses
   - Database migration scripts

3. **Table-Driven Tests**: Follow the table-driven test pattern shown in examples for consistency and maintainability.

4. **Test Naming**: Use format `Test<Function>_<Scenario>_<ExpectedResult>` for clarity.
