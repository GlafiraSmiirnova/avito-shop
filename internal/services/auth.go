package services

import (
	"avito-shop/config"
	wrappedErrors "avito-shop/internal/errors"
	"avito-shop/internal/repository"
	"context"
	"go.uber.org/zap"
	"net/http"
)

type AuthService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

func (s *AuthService) AuthenticateUser(ctx context.Context, username, password string) (string, error) {
	user, err := s.userRepo.GetUserByUsername(ctx, nil, username)
	if err != nil {
		zap.L().Info("Пользователь не найден, создаем нового", zap.String("username", username))

		hashedPassword, _ := config.HashPassword(password)
		if err := s.userRepo.CreateUser(ctx, nil, username, hashedPassword); err != nil {
			return "", wrappedErrors.New(http.StatusInternalServerError, "ошибка создания пользователя")

		}
		return config.GenerateJWT(username)
	}

	if !config.ComparePasswords(user.Password, password) {
		zap.L().Warn("Неверный пароль", zap.String("username", username))
		return "", wrappedErrors.New(http.StatusBadRequest, "неверный пароль")
	}

	return config.GenerateJWT(username)
}
