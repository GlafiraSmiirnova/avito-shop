package util

import (
	wrappedErrors "avito-shop/internal/errors"
	"avito-shop/internal/models/responses"
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func HandleErrors(c *gin.Context) {
	c.Next()

	if len(c.Errors) == 0 {
		return
	}

	if len(c.Errors) > 1 {
		zap.L().Error("error handling middleware cannot handle more than one error at once")
		return
	}

	var serviceError *wrappedErrors.Error
	ok := errors.As(c.Errors.Last(), &serviceError)
	if !ok {
		zap.L().Warn("unsolicited non-wrapped error", zap.Error(c.Errors[0]))
		return
	}

	c.AbortWithStatusJSON(serviceError.Code, &responses.ErrorResponse{Errors: serviceError.Err})
}
