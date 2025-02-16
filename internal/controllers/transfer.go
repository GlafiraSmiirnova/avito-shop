package controllers

import (
	"avito-shop/internal/controllers/util"
	"avito-shop/internal/models/requests"
	"avito-shop/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type TransferController struct {
	transferService *services.TransferService
}

func NewTransferController(transferService *services.TransferService) *TransferController {
	return &TransferController{transferService: transferService}
}

func (tc *TransferController) SendCoin(c *gin.Context) {
	username := c.GetString("username")

	data := util.GetRequestData[requests.TransferRequest](c)

	err := tc.transferService.TransferCoins(c, username, data.ToUser, data.Amount)
	if err != nil {
		zap.L().Error("Ошибка перевода монет", zap.Error(err))
		_ = c.Error(err)
		return
	}

	c.Status(http.StatusOK)
}
