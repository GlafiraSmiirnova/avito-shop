package services

import (
	"avito-shop/internal/models"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetUserInfo(t *testing.T) {
	service := InfoService{userRepo: mockUserRepo, txRepo: mockTxRepo}
	ctx := context.TODO()

	t.Run("Успешное получение информации о пользователе", func(t *testing.T) {
		mockUserRepo.On("GetUserByUsername", ctx, nil, "Alice").
			Return(&models.User{ID: 1, Balance: 500}, nil).Once()

		mockUserRepo.On("GetUserInventory", ctx, nil, 1).
			Return([]models.InventoryItem{
				{Type: "t-shirt", Quantity: 1},
				{Type: "cup", Quantity: 2},
			}, nil).Once()

		mockTxRepo.On("GetReceivedTransactions", ctx, nil, 1).
			Return([]models.ReceivedTransaction{
				{FromUser: "Bob", Amount: 50},
			}, nil).Once()

		mockTxRepo.On("GetSentTransactions", ctx, nil, 1).
			Return([]models.SentTransaction{
				{ToUser: "Charlie", Amount: 100},
			}, nil).Once()

		info, err := service.GetUserInfo(ctx, "Alice")
		assert.NoError(t, err)
		assert.Equal(t, 500, info.Coins)
		assert.Len(t, info.Inventory, 2)
		assert.Len(t, info.CoinHistory.Received, 1)
		assert.Len(t, info.CoinHistory.Sent, 1)

		mockUserRepo.AssertExpectations(t)
		mockTxRepo.AssertExpectations(t)
	})

	t.Run("Ошибка - пользователь не найден", func(t *testing.T) {
		mockUserRepo.On("GetUserByUsername", ctx, nil, "Alice").
			Return(nil, errors.New("пользователь не найден")).Once()

		info, err := service.GetUserInfo(ctx, "Alice")
		assert.Error(t, err)
		assert.Nil(t, info)
		assert.Equal(t, "пользователь не найден", err.Error())

		mockUserRepo.AssertExpectations(t)
	})

	t.Run("Ошибка - не удалось получить инвентарь", func(t *testing.T) {
		mockUserRepo.On("GetUserByUsername", ctx, nil, "Alice").
			Return(&models.User{ID: 1, Balance: 500}, nil).Once()

		mockUserRepo.On("GetUserInventory", ctx, nil, 1).
			Return(nil, errors.New("ошибка получения инвентаря")).Once()

		info, err := service.GetUserInfo(ctx, "Alice")
		assert.Error(t, err)
		assert.Nil(t, info)
		assert.Equal(t, "ошибка получения инвентаря", err.Error())

		mockUserRepo.AssertExpectations(t)
	})

	t.Run("Ошибка - не удалось получить входящие транзакции", func(t *testing.T) {
		mockUserRepo.On("GetUserByUsername", ctx, nil, "Alice").
			Return(&models.User{ID: 1, Balance: 500}, nil).Once()

		mockUserRepo.On("GetUserInventory", ctx, nil, 1).
			Return([]models.InventoryItem{}, nil).Once()

		mockTxRepo.On("GetReceivedTransactions", ctx, nil, 1).
			Return(nil, errors.New("ошибка получения входящих транзакций")).Once()

		info, err := service.GetUserInfo(ctx, "Alice")
		assert.Error(t, err)
		assert.Nil(t, info)
		assert.Equal(t, "ошибка получения входящих транзакций", err.Error())

		mockUserRepo.AssertExpectations(t)
		mockTxRepo.AssertExpectations(t)
	})

	t.Run("Ошибка - не удалось получить исходящие транзакции", func(t *testing.T) {
		mockUserRepo.On("GetUserByUsername", ctx, nil, "Alice").
			Return(&models.User{ID: 1, Balance: 500}, nil).Once()

		mockUserRepo.On("GetUserInventory", ctx, nil, 1).
			Return([]models.InventoryItem{}, nil).Once()

		mockTxRepo.On("GetReceivedTransactions", ctx, nil, 1).
			Return([]models.ReceivedTransaction{}, nil).Once()

		mockTxRepo.On("GetSentTransactions", ctx, nil, 1).
			Return(nil, errors.New("ошибка получения исходящих транзакций")).Once()

		info, err := service.GetUserInfo(ctx, "Alice")
		assert.Error(t, err)
		assert.Nil(t, info)
		assert.Equal(t, "ошибка получения исходящих транзакций", err.Error())

		mockUserRepo.AssertExpectations(t)
		mockTxRepo.AssertExpectations(t)
	})
}
