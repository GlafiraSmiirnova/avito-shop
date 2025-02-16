package controllers

import (
	"avito-shop/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type InfoController struct {
	infoService *services.InfoService
}

func NewInfoController(infoService *services.InfoService) *InfoController {
	return &InfoController{infoService: infoService}
}

func (ic *InfoController) GetUserInfo(c *gin.Context) {
	username := c.GetString("username")

	info, err := ic.infoService.GetUserInfo(c, username)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, info)
}
