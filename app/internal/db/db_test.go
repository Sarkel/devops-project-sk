package db

import (
	"database/sql"
	"testing"

	genDb "devops/app/internal/db/gen"
	cDB "devops/common/db"

	"github.com/stretchr/testify/assert"
)

type MockConManager struct {
	db *sql.DB
}

func (m *MockConManager) GetDB() *sql.DB {
	return m.db
}

func (m *MockConManager) Close() error {
	return nil
}

func TestWithQ_CreatesQueryInstance(t *testing.T) {
	// Test that WithQ creates a valid Queries instance
	conManager := &cDB.ConManager{}

	// Note: This would normally need a real DB connection
	// For now, we just test that the function signature is correct
	// In a real test environment, you'd use a test database

	// Test that nil doesn't panic
	assert.NotPanics(t, func() {
		_ = WithQ(conManager)
	})
}

func TestWithQ_ReturnsQueriesType(t *testing.T) {
	// Test that WithQ returns the correct type
	conManager := &cDB.ConManager{}

	queries := WithQ(conManager)

	// Verify that a valid Queries instance is created
	assert.NotNil(t, queries)
	assert.IsType(t, &genDb.Queries{}, queries)
}
