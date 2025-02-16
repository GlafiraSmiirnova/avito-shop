package main

import (
	"avito-shop/config"
	"avito-shop/config/db"
	"avito-shop/internal/controllers"
	"avito-shop/internal/controllers/util"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Не удалось загрузить .env")
	}

	config.InitLogger()
	config.InitJWT()

	log.Println("подключаемся к БД")

	db.ConnectDB()
	defer db.CloseDB()

	engine := gin.Default()
	config.SetupLogger(engine)

	engine.Use(util.HandleErrors)

	controllers.Register(engine)

	server := &http.Server{
		Addr:    ":8080",
		Handler: engine,
	}

	sig := make(chan os.Signal, 1)

	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sig
		zap.S().Info("Выключение сервера...")
		server.Shutdown(nil)
	}()

	zap.S().Infow("starting", "addr", server.Addr)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		zap.L().Fatal("listen error", zap.Error(err))
	}
}
