package repository

import (
	"avito-shop/config/db"
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type MerchRepositoryPG struct{}

func (r *MerchRepositoryPG) GetItemByName(ctx context.Context, tx pgx.Tx, itemName string) (int, int, error) {
	var itemID, price int

	zap.L().Info("Executing SQL", zap.String("itemName", itemName))
	err := db.GetExecutor(tx).QueryRow(ctx, "SELECT id, price FROM merch WHERE name=$1", itemName).Scan(&itemID, &price)
	if err != nil {
		zap.L().Error("Item not found", zap.String("itemName", itemName), zap.Error(err))
		return 0, 0, errors.New("предмет не найден")
	}
	return itemID, price, nil
}

func (r *MerchRepositoryPG) AddToInventory(ctx context.Context, tx pgx.Tx, userID, itemID int) error {
	_, err := db.GetExecutor(tx).Exec(ctx, `
		INSERT INTO inventory (user_id, merch_id, quantity) 
		VALUES ($1, $2, 1) 
		ON CONFLICT (user_id, merch_id) DO UPDATE SET quantity = inventory.quantity + 1
	`, userID, itemID)
	return err
}
