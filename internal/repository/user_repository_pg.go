package repository

import (
	"avito-shop/config/db"
	"avito-shop/internal/models"
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
)

type UserRepositoryPG struct{}

func (r *UserRepositoryPG) GetUserByUsername(ctx context.Context, tx pgx.Tx, username string) (*models.User, error) {
	var user models.User
	err := db.GetExecutor(tx).QueryRow(ctx, "SELECT id, username, password_hash, balance FROM users WHERE username=$1", username).
		Scan(&user.ID, &user.Username, &user.Password, &user.Balance)

	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryPG) CreateUser(ctx context.Context, tx pgx.Tx, username, passwordHash string) error {
	_, err := db.GetExecutor(tx).Exec(ctx, "INSERT INTO users (username, password_hash, balance) VALUES ($1, $2, 1000)", username, passwordHash)
	return err
}

func (r *UserRepositoryPG) UpdateBalance(ctx context.Context, tx pgx.Tx, userID, amount int) error {
	_, err := db.GetExecutor(tx).Exec(ctx, "UPDATE users SET balance = balance + $1 WHERE id=$2", amount, userID)
	return err
}

func (r *UserRepositoryPG) GetUserInventory(ctx context.Context, tx pgx.Tx, userID int) ([]models.InventoryItem, error) {
	rows, err := db.GetExecutor(tx).Query(ctx, `
		SELECT m.name, i.quantity 
		FROM inventory i 
		JOIN merch m ON i.merch_id = m.id 
		WHERE i.user_id=$1
	`, userID)

	if err != nil {
		return nil, errors.New("ошибка получения инвентаря")
	}
	defer rows.Close()

	var inventory []models.InventoryItem
	for rows.Next() {
		var item models.InventoryItem
		if err := rows.Scan(&item.Type, &item.Quantity); err == nil {
			inventory = append(inventory, item)
		}
	}
	return inventory, nil
}
