package config

import (
	"time"

	"github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func InitLogger() {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, _ := config.Build()
	zap.ReplaceGlobals(logger)
	Logger = logger
}

func SetupLogger(engine *gin.Engine) {
	engine.Use(ginzap.Ginzap(Logger, time.RFC3339, true))
	engine.Use(ginzap.RecoveryWithZap(Logger, true))
}
