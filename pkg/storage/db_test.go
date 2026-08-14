package storage

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateTableIfNotExists_GivenValidDb_whenCalled_thenCreatesTable(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer dbMock.Close()

	// Expect the table creation
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS transactions").
		WillReturnResult(sqlmock.NewResult(0, 1))

	oldDB := db
	db = dbMock
	defer func() { db = oldDB }()

	err = createTableIfNotExists()

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetDB_GivenNilConnection_whenCalled_thenReturnsError(t *testing.T) {
	// Ensure db is nil
	oldDB := db
	db = nil
	defer func() { db = oldDB }()

	result, err := GetDB()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database connection pool is not initialized")
	assert.Nil(t, result)
}

func TestGetDB_GivenValidConnection_whenCalled_thenReturnsDatabase(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer dbMock.Close()

	oldDB := db
	db = dbMock
	defer func() { db = oldDB }()

	result, err := GetDB()

	assert.NoError(t, err)
	assert.Equal(t, dbMock, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCloseDB_GivenValidConnection_whenClosed_thenCleansUp(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	assert.NoError(t, err)

	mock.ExpectClose()

	oldDB := db
	db = dbMock
	CloseDB()

	assert.NoError(t, mock.ExpectationsWereMet())
	db = oldDB
}

func TestCloseDB_GivenNilConnection_whenCalled_thenDoesNothing(t *testing.T) {
	oldDB := db
	db = nil
	defer func() { db = oldDB }()

	// Should not panic
	assert.NotPanics(t, func() { CloseDB() })
}
