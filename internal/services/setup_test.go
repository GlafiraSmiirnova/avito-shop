package services

import (
	"avito-shop/config/db"
	"avito-shop/internal/mocks"
	"context"
	"github.com/jackc/pgx/v5"
	"os"
	"testing"
)

var mockUserRepo *mocks.MockUserRepo
var mockMerchRepo *mocks.MockMerchRepo
var mockTxRepo *mocks.MockTxRepo

func TestMain(m *testing.M) {
	mockUserRepo = &mocks.MockUserRepo{}
	mockMerchRepo = &mocks.MockMerchRepo{}
	mockTxRepo = &mocks.MockTxRepo{}

	mockTx := new(mocks.MockTx)

	db.BeginTransaction = func(ctx context.Context) (pgx.Tx, error) {
		return mockTx, nil
	}

	db.CommitTransaction = func(ctx context.Context, tx pgx.Tx) error {
		return nil
	}

	code := m.Run()
	os.Exit(code)
}
