package util

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

const requestDataKey = "request_data"

func ValidateRequestData[BodyType any](c *gin.Context) {
	var body BodyType

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errors": fmt.Sprintf("неверный формат данных: %s", err.Error()),
		})
		c.Abort()
		return
	}

	c.Set(requestDataKey, body)
	c.Next()
}

func GetRequestData[BodyType any](c *gin.Context) BodyType {
	return c.MustGet(requestDataKey).(BodyType)
}
