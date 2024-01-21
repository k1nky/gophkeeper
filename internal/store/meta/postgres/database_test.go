package database

import (
	"context"
	"os"
	"testing"
)

func openTestDB() (*PostgresStorage, error) {
	a := New("postgres://postgres:postgres@localhost:5432/praktikum?sslmode=disable")
	if err := a.Open(context.TODO()); err != nil {
		return nil, err
	}
	return a, nil
}

func shouldSkipDBTest(t *testing.T) bool {
	// TODO:
	// return false
	if len(os.Getenv("TEST_DB_READY")) == 0 {
		t.Skip()
		return true
	}
	return false
}

func TestAdapter(t *testing.T) {
	// suite.Run(t, new(usersTestSuite))
	// suite.Run(t, new(ordersTestSuite))
}
