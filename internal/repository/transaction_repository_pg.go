package repository

import (
	"avito-shop/config/db"
	"avito-shop/internal/models"
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
)

type TransactionRepositoryPG struct{}

func (r *TransactionRepositoryPG) RecordTransaction(ctx context.Context, tx pgx.Tx, fromUserID, toUserID, amount int) error {
	_, err := db.GetExecutor(tx).Exec(ctx, "INSERT INTO transactions (from_user, to_user, amount) VALUES ($1, $2, $3)", fromUserID, toUserID, amount)
	return err
}

func (r *TransactionRepositoryPG) GetReceivedTransactions(ctx context.Context, tx pgx.Tx, userID int) ([]models.ReceivedTransaction, error) {
	rows, err := db.GetExecutor(tx).Query(ctx, `
		SELECT u.username, t.amount 
		FROM transactions t 
		JOIN users u ON t.from_user = u.id 
		WHERE t.to_user=$1
	`, userID)

	if err != nil {
		return nil, errors.New("ошибка получения входящих транзакций")
	}
	defer rows.Close()

	var received []models.ReceivedTransaction
	for rows.Next() {
		var tx models.ReceivedTransaction
		if err := rows.Scan(&tx.FromUser, &tx.Amount); err == nil {
			received = append(received, tx)
		}
	}
	return received, nil
}

func (r *TransactionRepositoryPG) GetSentTransactions(ctx context.Context, tx pgx.Tx, userID int) ([]models.SentTransaction, error) {
	rows, err := db.GetExecutor(tx).Query(ctx, `
		SELECT u.username, t.amount 
		FROM transactions t 
		JOIN users u ON t.to_user = u.id 
		WHERE t.from_user=$1
	`, userID)

	if err != nil {
		return nil, errors.New("ошибка получения исходящих транзакций")
	}
	defer rows.Close()

	var sent []models.SentTransaction
	for rows.Next() {
		var tx models.SentTransaction
		if err := rows.Scan(&tx.ToUser, &tx.Amount); err == nil {
			sent = append(sent, tx)
		}
	}
	return sent, nil
}
