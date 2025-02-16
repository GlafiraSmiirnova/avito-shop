package services

import (
	"avito-shop/internal/models"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestAuthenticateUser(t *testing.T) {
	assert.NotNil(t, mockUserRepo)
	service := AuthService{userRepo: mockUserRepo}
	ctx := context.TODO()

	t.Run("Создание нового пользователя", func(t *testing.T) {
		mockUserRepo.On("GetUserByUsername", ctx, nil, "newUser").Return(nil, errors.New("not found")).Once()
		mockUserRepo.On("CreateUser", ctx, nil, "newUser", mock.Anything).Return(nil).Once()

		_, err := service.AuthenticateUser(ctx, "newUser", "password")
		assert.NoError(t, err)

		mockUserRepo.AssertExpectations(t)
	})

	t.Run("Ошибка при неверном пароле", func(t *testing.T) {
		mockUserRepo.On("GetUserByUsername", ctx, nil, "existingUser").Return(&models.User{
			ID:       1,
			Username: "existingUser",
			Password: "wrongHashedPassword",
		}, nil).Once()

		_, err := service.AuthenticateUser(ctx, "existingUser", "wrongPassword")
		assert.Error(t, err)
		assert.Equal(t, "неверный пароль", err.Error())

		mockUserRepo.AssertExpectations(t)
	})

	t.Run("Ошибка при создании пользователя", func(t *testing.T) {
		mockUserRepo.On("GetUserByUsername", ctx, nil, "newUser").Return(nil, errors.New("not found")).Once()
		mockUserRepo.On("CreateUser", ctx, nil, "newUser", mock.Anything).Return(errors.New("ошибка создания пользователя")).Once()

		_, err := service.AuthenticateUser(ctx, "newUser", "password")
		assert.Error(t, err)
		assert.Equal(t, "ошибка создания пользователя", err.Error())

		mockUserRepo.AssertExpectations(t)
	})
}
