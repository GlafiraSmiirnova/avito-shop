package services

import (
	"avito-shop/config/db"
	wrappedErrors "avito-shop/internal/errors"
	"avito-shop/internal/repository"
	"context"
	"net/http"
)

type MerchService struct {
	userRepo  repository.UserRepository
	merchRepo repository.MerchRepository
}

func NewMerchService(userRepo repository.UserRepository, merchRepo repository.MerchRepository) *MerchService {
	return &MerchService{userRepo, merchRepo}
}

func (s *MerchService) BuyItem(ctx context.Context, username, itemName string) error {
	tx, err := db.BeginTransaction(ctx)
	if err != nil {
		return wrappedErrors.New(http.StatusInternalServerError, "ошибка начала транзакции")
	}
	defer tx.Rollback(ctx)

	user, err := s.userRepo.GetUserByUsername(ctx, tx, username)
	if err != nil {
		return wrappedErrors.New(http.StatusInternalServerError, "пользователь не найден")
	}

	itemID, itemPrice, err := s.merchRepo.GetItemByName(ctx, tx, itemName)
	if err != nil {
		return wrappedErrors.New(http.StatusBadRequest, "предмет не найден")
	}

	if user.Balance < itemPrice {
		return wrappedErrors.New(http.StatusBadRequest, "недостаточно средств")
	}

	if err := s.userRepo.UpdateBalance(ctx, tx, user.ID, -itemPrice); err != nil {
		return wrappedErrors.New(http.StatusInternalServerError, "ошибка при списании")
	}

	if err := s.merchRepo.AddToInventory(ctx, tx, user.ID, itemID); err != nil {
		return wrappedErrors.New(http.StatusInternalServerError, "ошибка при добавлении в инвентарь")
	}

	if err := db.CommitTransaction(ctx, tx); err != nil {
		return wrappedErrors.New(http.StatusInternalServerError, "ошибка фиксации транзакции")
	}

	return nil
}
