package repository

import (
	"context"
	"github.com/jackc/pgx/v5"
)

type MerchRepository interface {
	GetItemByName(ctx context.Context, tx pgx.Tx, itemName string) (int, int, error)
	AddToInventory(ctx context.Context, tx pgx.Tx, userID, itemID int) error
}
