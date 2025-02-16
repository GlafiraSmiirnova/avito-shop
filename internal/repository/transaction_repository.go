package repository

import (
	"avito-shop/internal/models"
	"context"
	"github.com/jackc/pgx/v5"
)

type TransactionRepository interface {
	RecordTransaction(ctx context.Context, tx pgx.Tx, fromUserID, toUserID, amount int) error
	GetReceivedTransactions(ctx context.Context, tx pgx.Tx, userID int) ([]models.ReceivedTransaction, error)
	GetSentTransactions(ctx context.Context, tx pgx.Tx, userID int) ([]models.SentTransaction, error)
}
