package controllers

import (
	"avito-shop/internal/controllers/util"
	"avito-shop/internal/models/requests"
	"avito-shop/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService *services.AuthService
}

func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

func (ac *AuthController) Auth(c *gin.Context) {
	data := util.GetRequestData[requests.AuthRequest](c)

	token, err := ac.authService.AuthenticateUser(c, data.Username, data.Password)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
