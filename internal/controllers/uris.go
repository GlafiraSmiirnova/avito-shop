package controllers

import (
	"avito-shop/internal/controllers/access"
	"avito-shop/internal/controllers/util"
	"avito-shop/internal/models/requests"
	"avito-shop/internal/repository"
	"avito-shop/internal/services"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
)

func Register(engine *gin.Engine) {

	userRepo := &repository.UserRepositoryPG{}
	merchRepo := &repository.MerchRepositoryPG{}
	txRepo := &repository.TransactionRepositoryPG{}

	authService := services.NewAuthService(userRepo)
	merchService := services.NewMerchService(userRepo, merchRepo)
	transferService := services.NewTransferService(userRepo, txRepo)
	infoService := services.NewInfoService(userRepo, txRepo)

	authController := NewAuthController(authService)
	merchController := NewMerchController(merchService)
	transferController := NewTransferController(transferService)
	infoController := NewInfoController(infoService)

	engine.StaticFile("/docs.yaml", "schema/docs.yaml")
	engine.GET("/swagger-ui/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("../docs.yaml")))

	engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api := engine.Group("/api")
	{
		api.POST("/auth", util.ValidateRequestData[requests.AuthRequest], authController.Auth)

		protected := api.Group("/")
		protected.Use(access.AuthMiddleware())

		protected.POST("/sendCoin", util.ValidateRequestData[requests.TransferRequest], transferController.SendCoin)
		protected.GET("/info", infoController.GetUserInfo)
		protected.GET("/buy/:item", merchController.BuyItem)
	}
}
