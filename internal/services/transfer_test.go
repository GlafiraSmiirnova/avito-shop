package services

import (
	"avito-shop/internal/mocks"
	"avito-shop/internal/models"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTransferCoins(t *testing.T) {
	assert.NotNil(t, mockUserRepo)
	assert.NotNil(t, mockTxRepo)

	service := TransferService{userRepo: mockUserRepo, txRepo: mockTxRepo}
	ctx := context.TODO()

	t.Run("Успешный перевод", func(t *testing.T) {
		tx := new(mocks.MockTx)
		mockUserRepo.On("GetUserByUsername", ctx, tx, "Alice").Return(&models.User{ID: 1, Balance: 200}, nil).Once()
		mockUserRepo.On("GetUserByUsername", ctx, tx, "Bob").Return(&models.User{ID: 2, Balance: 100}, nil).Once()
		mockUserRepo.On("UpdateBalance", ctx, tx, 1, -50).Return(nil).Once()
		mockUserRepo.On("UpdateBalance", ctx, tx, 2, 50).Return(nil).Once()
		mockTxRepo.On("RecordTransaction", ctx, tx, 1, 2, 50).Return(nil).Once()

		err := service.TransferCoins(ctx, "Alice", "Bob", 50)
		assert.NoError(t, err)

		mockUserRepo.AssertExpectations(t)
		mockTxRepo.AssertExpectations(t)
	})

	t.Run("Ошибка - недостаточно средств", func(t *testing.T) {
		tx := new(mocks.MockTx)
		mockUserRepo.On("GetUserByUsername", ctx, tx, "Alice").Return(&models.User{ID: 1, Balance: 20}, nil).Once()
		mockUserRepo.On("GetUserByUsername", ctx, tx, "Bob").Return(&models.User{ID: 2, Balance: 100}, nil).Once()

		err := service.TransferCoins(ctx, "Alice", "Bob", 50)
		assert.Error(t, err)
		assert.Equal(t, "недостаточно средств", err.Error())

		mockUserRepo.AssertExpectations(t)
	})

	t.Run("Ошибка при обновлении баланса отправителя", func(t *testing.T) {
		tx := new(mocks.MockTx)
		mockUserRepo.On("GetUserByUsername", ctx, tx, "Alice").Return(&models.User{ID: 1, Balance: 100}, nil).Once()
		mockUserRepo.On("GetUserByUsername", ctx, tx, "Bob").Return(&models.User{ID: 2, Balance: 50}, nil).Once()
		mockUserRepo.On("UpdateBalance", ctx, tx, 1, -50).Return(errors.New("db error")).Once()

		err := service.TransferCoins(ctx, "Alice", "Bob", 50)
		assert.Error(t, err)
		assert.Equal(t, "ошибка обновления баланса отправителя", err.Error())

		mockUserRepo.AssertExpectations(t)
		mockTxRepo.AssertExpectations(t)
	})

	t.Run("Ошибка - отправитель не найден", func(t *testing.T) {
		tx := new(mocks.MockTx)
		mockUserRepo.On("GetUserByUsername", ctx, tx, "UnknownUser").Return(nil, errors.New("отправитель не найден")).Once()

		err := service.TransferCoins(ctx, "UnknownUser", "Bob", 50)
		assert.Error(t, err)
		assert.Equal(t, "отправитель не найден", err.Error())

		mockUserRepo.AssertExpectations(t)
	})

	t.Run("Ошибка - получатель не найден", func(t *testing.T) {
		tx := new(mocks.MockTx)
		mockUserRepo.On("GetUserByUsername", ctx, tx, "Alice").Return(&models.User{ID: 1, Balance: 200}, nil).Once()
		mockUserRepo.On("GetUserByUsername", ctx, tx, "UnknownUser").Return(nil, errors.New("получатель не найден")).Once()

		err := service.TransferCoins(ctx, "Alice", "UnknownUser", 50)
		assert.Error(t, err)
		assert.Equal(t, "получатель не найден", err.Error())

		mockUserRepo.AssertExpectations(t)
	})

	t.Run("Ошибка - обновление баланса получателя", func(t *testing.T) {
		tx := new(mocks.MockTx)
		mockUserRepo.On("GetUserByUsername", ctx, tx, "Alice").Return(&models.User{ID: 1, Balance: 200}, nil).Once()
		mockUserRepo.On("GetUserByUsername", ctx, tx, "Bob").Return(&models.User{ID: 2, Balance: 100}, nil).Once()
		mockUserRepo.On("UpdateBalance", ctx, tx, 1, -50).Return(nil).Once()
		mockUserRepo.On("UpdateBalance", ctx, tx, 2, 50).Return(errors.New("ошибка обновления баланса получателя")).Once()

		err := service.TransferCoins(ctx, "Alice", "Bob", 50)
		assert.Error(t, err)
		assert.Equal(t, "ошибка обновления баланса получателя", err.Error())

		mockUserRepo.AssertExpectations(t)
	})

	t.Run("Ошибка - запись транзакции не удалась", func(t *testing.T) {
		tx := new(mocks.MockTx)
		mockUserRepo.On("GetUserByUsername", ctx, tx, "Alice").Return(&models.User{ID: 1, Balance: 200}, nil).Once()
		mockUserRepo.On("GetUserByUsername", ctx, tx, "Bob").Return(&models.User{ID: 2, Balance: 100}, nil).Once()
		mockUserRepo.On("UpdateBalance", ctx, tx, 1, -50).Return(nil).Once()
		mockUserRepo.On("UpdateBalance", ctx, tx, 2, 50).Return(nil).Once()
		mockTxRepo.On("RecordTransaction", ctx, tx, 1, 2, 50).Return(errors.New("ошибка записи транзакции")).Once()

		err := service.TransferCoins(ctx, "Alice", "Bob", 50)
		assert.Error(t, err)
		assert.Equal(t, "ошибка записи транзакции монеток", err.Error())

		mockUserRepo.AssertExpectations(t)
		mockTxRepo.AssertExpectations(t)
	})
}
