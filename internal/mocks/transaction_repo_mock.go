package mocks

import (
	"avito-shop/internal/models"
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/mock"
)

type MockTxRepo struct {
	mock.Mock
}

func (m *MockTxRepo) GetReceivedTransactions(ctx context.Context, tx pgx.Tx, userID int) ([]models.ReceivedTransaction, error) {
	args := m.Called(ctx, tx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.ReceivedTransaction), args.Error(1)
}

func (m *MockTxRepo) GetSentTransactions(ctx context.Context, tx pgx.Tx, userID int) ([]models.SentTransaction, error) {
	args := m.Called(ctx, tx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.SentTransaction), args.Error(1)
}

func (m *MockTxRepo) RecordTransaction(ctx context.Context, tx pgx.Tx, fromUserID, toUserID, amount int) error {
	args := m.Called(ctx, tx, fromUserID, toUserID, amount)
	return args.Error(0)
}
