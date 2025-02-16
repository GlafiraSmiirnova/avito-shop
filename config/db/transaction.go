package db

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

var BeginTransaction = func(ctx context.Context) (pgx.Tx, error) {
	if DB == nil {
		return nil, errors.New("DB не инициализирована")
	}
	tx, err := DB.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		zap.L().Error("Ошибка начала транзакции", zap.Error(err))
		return nil, err
	}
	return tx, nil
}

var CommitTransaction = func(ctx context.Context, tx pgx.Tx) error {
	err := tx.Commit(ctx)
	if err != nil {
		zap.L().Error("Ошибка коммита транзакции", zap.Error(err))
		return err
	}
	return nil
}
