package controllers

import (
	"avito-shop/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type MerchController struct {
	merchService *services.MerchService
}

func NewMerchController(merchService *services.MerchService) *MerchController {
	return &MerchController{merchService: merchService}
}

func (mc *MerchController) BuyItem(c *gin.Context) {
	username := c.GetString("username")
	itemName := c.Param("item")

	err := mc.merchService.BuyItem(c, username, itemName)
	if err != nil {
		zap.L().Error("Ошибка покупки предмета", zap.String("item", itemName), zap.Error(err))
		_ = c.Error(err)
		return
	}

	c.Status(http.StatusOK)
}
