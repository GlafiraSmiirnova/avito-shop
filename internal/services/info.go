package services

import (
	wrappedErrors "avito-shop/internal/errors"
	"avito-shop/internal/models"
	"avito-shop/internal/models/responses"
	"avito-shop/internal/repository"
	"context"
	"net/http"
)

type InfoService struct {
	userRepo repository.UserRepository
	txRepo   repository.TransactionRepository
}

func NewInfoService(userRepo repository.UserRepository, txRepo repository.TransactionRepository) *InfoService {
	return &InfoService{userRepo: userRepo, txRepo: txRepo}
}

func (s *InfoService) GetUserInfo(ctx context.Context, username string) (*responses.InfoResponse, error) {

	user, err := s.userRepo.GetUserByUsername(ctx, nil, username)
	if err != nil {
		return nil, wrappedErrors.New(http.StatusInternalServerError, "пользователь не найден")
	}

	inventory, err := s.userRepo.GetUserInventory(ctx, nil, user.ID)
	if err != nil {
		return nil, wrappedErrors.New(http.StatusInternalServerError, "ошибка получения инвентаря")
	}
	if inventory == nil {
		inventory = []models.InventoryItem{}
	}

	received, err := s.txRepo.GetReceivedTransactions(ctx, nil, user.ID)
	if err != nil {
		return nil, wrappedErrors.New(http.StatusInternalServerError, "ошибка получения входящих транзакций")
	}
	if received == nil {
		received = []models.ReceivedTransaction{}
	}

	sent, err := s.txRepo.GetSentTransactions(ctx, nil, user.ID)
	if err != nil {
		return nil, wrappedErrors.New(http.StatusInternalServerError, "ошибка получения исходящих транзакций")
	}
	if sent == nil {
		sent = []models.SentTransaction{}
	}

	return &responses.InfoResponse{
		Coins:     user.Balance,
		Inventory: inventory,
		CoinHistory: models.CoinHistory{
			Received: received,
			Sent:     sent,
		},
	}, nil
}
