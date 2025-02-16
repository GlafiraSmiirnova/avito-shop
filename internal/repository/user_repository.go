package repository

import (
	"avito-shop/internal/models"
	"context"
	"github.com/jackc/pgx/v5"
)

type UserRepository interface {
	GetUserByUsername(ctx context.Context, tx pgx.Tx, username string) (*models.User, error)
	CreateUser(ctx context.Context, tx pgx.Tx, username, passwordHash string) error
	UpdateBalance(ctx context.Context, tx pgx.Tx, userID, amount int) error
	GetUserInventory(ctx context.Context, tx pgx.Tx, userID int) ([]models.InventoryItem, error)
}
