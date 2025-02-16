package mocks

import (
	"avito-shop/internal/models"
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/mock"
)

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) GetUserByUsername(ctx context.Context, tx pgx.Tx, username string) (*models.User, error) {
	args := m.Called(ctx, tx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepo) CreateUser(ctx context.Context, tx pgx.Tx, username, passwordHash string) error {
	args := m.Called(ctx, tx, username, passwordHash)
	return args.Error(0)
}

func (m *MockUserRepo) UpdateBalance(ctx context.Context, tx pgx.Tx, userID, amount int) error {
	args := m.Called(ctx, tx, userID, amount)
	return args.Error(0)
}

func (m *MockUserRepo) GetUserInventory(ctx context.Context, tx pgx.Tx, userID int) ([]models.InventoryItem, error) {
	args := m.Called(ctx, tx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.InventoryItem), args.Error(1)
}
