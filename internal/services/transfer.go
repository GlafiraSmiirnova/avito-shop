package services

import (
	"avito-shop/config/db"
	wrappedErrors "avito-shop/internal/errors"
	"avito-shop/internal/repository"
	"context"
	"net/http"
)

type TransferService struct {
	userRepo repository.UserRepository
	txRepo   repository.TransactionRepository
}

func NewTransferService(userRepo repository.UserRepository, txRepo repository.TransactionRepository) *TransferService {
	return &TransferService{userRepo, txRepo}
}

func (s *TransferService) TransferCoins(ctx context.Context, fromUser, toUser string, amount int) error {
	tx, err := db.BeginTransaction(ctx)
	if err != nil {
		return wrappedErrors.New(http.StatusInternalServerError, "ошибка начала транзакции")
	}
	defer tx.Rollback(ctx)

	from, err := s.userRepo.GetUserByUsername(ctx, tx, fromUser)
	if err != nil || from == nil {
		return wrappedErrors.New(http.StatusInternalServerError, "отправитель не найден")
	}

	to, err := s.userRepo.GetUserByUsername(ctx, tx, toUser)
	if err != nil || to == nil {
		return wrappedErrors.New(http.StatusBadRequest, "получатель не найден")
	}

	if from.Balance < amount {
		return wrappedErrors.New(http.StatusBadRequest, "недостаточно средств")
	}

	if err := s.userRepo.UpdateBalance(ctx, tx, from.ID, -amount); err != nil {
		return wrappedErrors.New(http.StatusInternalServerError, "ошибка обновления баланса отправителя")
	}
	if err := s.userRepo.UpdateBalance(ctx, tx, to.ID, amount); err != nil {
		return wrappedErrors.New(http.StatusInternalServerError, "ошибка обновления баланса получателя")
	}

	if err := s.txRepo.RecordTransaction(ctx, tx, from.ID, to.ID, amount); err != nil {
		return wrappedErrors.New(http.StatusInternalServerError, "ошибка записи транзакции монеток")
	}

	if err := db.CommitTransaction(ctx, tx); err != nil {
		return wrappedErrors.New(http.StatusInternalServerError, "ошибка фиксации транзакции")
	}

	return nil
}
