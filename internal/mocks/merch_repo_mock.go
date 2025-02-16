package mocks

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/mock"
)

type MockMerchRepo struct {
	mock.Mock
}

func (m *MockMerchRepo) GetItemByName(ctx context.Context, tx pgx.Tx, itemName string) (int, int, error) {
	args := m.Called(ctx, tx, itemName)
	return args.Int(0), args.Int(1), args.Error(2)
}

func (m *MockMerchRepo) AddToInventory(ctx context.Context, tx pgx.Tx, userID, itemID int) error {
	args := m.Called(ctx, tx, userID, itemID)
	return args.Error(0)
}
