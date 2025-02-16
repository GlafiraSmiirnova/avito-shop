package services

import (
	"avito-shop/internal/mocks"
	"avito-shop/internal/models"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuyItem(t *testing.T) {
	assert.NotNil(t, mockUserRepo)
	assert.NotNil(t, mockMerchRepo)

	service := MerchService{userRepo: mockUserRepo, merchRepo: mockMerchRepo}
	ctx := context.TODO()

	t.Run("Покупка успешна", func(t *testing.T) {
		tx := new(mocks.MockTx)
		mockUserRepo.On("GetUserByUsername", ctx, tx, "user").Return(&models.User{ID: 1, Balance: 100}, nil).Once()
		mockMerchRepo.On("GetItemByName", ctx, tx, "t-shirt").Return(1, 80, nil).Once()
		mockUserRepo.On("UpdateBalance", ctx, tx, 1, -80).Return(nil).Once()
		mockMerchRepo.On("AddToInventory", ctx, tx, 1, 1).Return(nil).Once()

		err := service.BuyItem(ctx, "user", "t-shirt")
		assert.NoError(t, err)

		mockUserRepo.AssertExpectations(t)
		mockMerchRepo.AssertExpectations(t)
	})

	t.Run("Ошибка - недостаточно средств", func(t *testing.T) {
		tx := new(mocks.MockTx)
		mockUserRepo.On("GetUserByUsername", ctx, tx, "user").Return(&models.User{ID: 1, Balance: 50}, nil).Once()
		mockMerchRepo.On("GetItemByName", ctx, tx, "hoody").Return(2, 300, nil).Once()

		err := service.BuyItem(ctx, "user", "hoody")
		assert.Error(t, err)
		assert.Equal(t, "недостаточно средств", err.Error())

		mockUserRepo.AssertExpectations(t)
		mockMerchRepo.AssertExpectations(t)
	})

	t.Run("Ошибка при обновлении баланса", func(t *testing.T) {
		tx := new(mocks.MockTx)
		mockUserRepo.On("GetUserByUsername", ctx, tx, "user").Return(&models.User{ID: 1, Balance: 100}, nil).Once()
		mockMerchRepo.On("GetItemByName", ctx, tx, "t-shirt").Return(1, 80, nil).Once()
		mockUserRepo.On("UpdateBalance", ctx, tx, 1, -80).Return(errors.New("db error")).Once()

		err := service.BuyItem(ctx, "user", "t-shirt")
		assert.Error(t, err)
		assert.Equal(t, "ошибка при списании", err.Error())

		mockUserRepo.AssertExpectations(t)
		mockMerchRepo.AssertExpectations(t)
	})

	t.Run("Ошибка - пользователь не найден", func(t *testing.T) {
		tx := new(mocks.MockTx)
		mockUserRepo.On("GetUserByUsername", ctx, tx, "unknown_user").Return(nil, errors.New("пользователь не найден")).Once()

		err := service.BuyItem(ctx, "unknown_user", "t-shirt")
		assert.Error(t, err)
		assert.Equal(t, "пользователь не найден", err.Error())

		mockUserRepo.AssertExpectations(t)
	})

	t.Run("Ошибка - предмет не найден", func(t *testing.T) {
		tx := new(mocks.MockTx)
		mockUserRepo.On("GetUserByUsername", ctx, tx, "user").Return(&models.User{ID: 1, Balance: 100}, nil).Once()
		mockMerchRepo.On("GetItemByName", ctx, tx, "unknown_item").Return(0, 0, errors.New("предмет не найден")).Once()

		err := service.BuyItem(ctx, "user", "unknown_item")
		assert.Error(t, err)
		assert.Equal(t, "предмет не найден", err.Error())

		mockUserRepo.AssertExpectations(t)
		mockMerchRepo.AssertExpectations(t)
	})

	t.Run("Ошибка при добавлении в инвентарь", func(t *testing.T) {
		tx := new(mocks.MockTx)
		mockUserRepo.On("GetUserByUsername", ctx, tx, "user").Return(&models.User{ID: 1, Balance: 100}, nil).Once()
		mockMerchRepo.On("GetItemByName", ctx, tx, "t-shirt").Return(1, 80, nil).Once()
		mockUserRepo.On("UpdateBalance", ctx, tx, 1, -80).Return(nil).Once()
		mockMerchRepo.On("AddToInventory", ctx, tx, 1, 1).Return(errors.New("ошибка при добавлении в инвентарь")).Once()

		err := service.BuyItem(ctx, "user", "t-shirt")
		assert.Error(t, err)
		assert.Equal(t, "ошибка при добавлении в инвентарь", err.Error())

		mockUserRepo.AssertExpectations(t)
		mockMerchRepo.AssertExpectations(t)
	})
}
