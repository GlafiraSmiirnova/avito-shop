package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"os"
)

var DB *pgxpool.Pool

func ConnectDB() {
	dbName := os.Getenv("DATABASE_NAME")

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_PORT"),
		dbName,
	)

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		zap.L().Fatal("Ошибка подключения к БД", zap.Error(err))
	}
	DB = pool
	zap.L().Info("Подключение к БД успешно установлено")
}

func ConnectTestDB() {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_PORT"),
		os.Getenv("DATABASE_NAME"),
	)

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		zap.L().Fatal("Ошибка подключения к тестовой БД", zap.Error(err))
	}
	DB = pool
	zap.L().Info("Подключение к тестовой БД успешно установлено")
}

func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}
