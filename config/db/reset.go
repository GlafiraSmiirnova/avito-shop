package db

import (
	"context"
	"go.uber.org/zap"
)

func ResetTestDB() {
	_, err := DB.Exec(context.Background(), `
		TRUNCATE TABLE transactions RESTART IDENTITY CASCADE;
		TRUNCATE TABLE users RESTART IDENTITY CASCADE;
	`)
	if err != nil {
		zap.L().Fatal("Ошибка очистки тестовой БД", zap.Error(err))
	}

	_, err = DB.Exec(context.Background(), `
		INSERT INTO users (username, password_hash) VALUES 
		('test_user1', 'hashed_password1'),
		('test_user2', 'hashed_password2');
	`)
	if err != nil {
		zap.L().Fatal("Ошибка наполнения тестовой БД", zap.Error(err))
	}
}
