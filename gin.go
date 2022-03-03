package log

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

func newGinLogger(config *Config) *zap.Logger {
	accessPriority := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level == zapcore.InfoLevel
	})
	errorPriority := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level == zap.ErrorLevel
	})

	var cores []zapcore.Core
	var encoder zapcore.Encoder
	switch config.Format {
	case jsonFormat:
		encoder = zapcore.NewJSONEncoder(jsonEncoderConfig())
	case textFormat:
		encoder = zapcore.NewConsoleEncoder(consoleEncoderConfig(false))
	}

	cores = append(
		cores,
		zapcore.NewCore(encoder, config.fileWriteSyncer("gin.access"), accessPriority),
		zapcore.NewCore(encoder, config.fileWriteSyncer("gin.recovery"), errorPriority),
	)

	logger := zap.New(zapcore.NewTee(cores...)).WithOptions(
		zap.AddStacktrace(errorPriority),
	)

	if config.Development {
		logger = logger.WithOptions(zap.Development())
	}

	return logger
}

func AccessLoggerWithConfigFile(configFile string) gin.HandlerFunc {
	config, err := parseConfigFile(configFile)
	if err != nil {
		fmt.Println("parse config error, using default config")
		config = defaultConfig
	}
	return AccessLogger(config)
}

func AccessLogger(config *Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		mergeConfig(config, defaultConfig)
		logger := newGinLogger(config)
		startTime := time.Now()
		method := c.Request.Method
		c.Next()

		cost := time.Since(startTime)
		switch {
		case config.Format == "text":
			logger.Sugar().Infof("[%s] %v \"%v\" %v %v %v %v",
				config.Application,
				c.ClientIP(),
				c.Request.URL.RequestURI(),
				c.Writer.Status(),
				cost,
				method,
				c.Request.UserAgent(),
			)
		case config.Format == "json":
			logger.Sugar().With(zap.String("app", config.Application),
				zap.String("remote_addr", c.ClientIP()),
				zap.String("uri", c.Request.URL.RequestURI()),
				zap.Int("status", c.Writer.Status()),
				zap.Duration("request_time", cost),
				zap.String("method", method),
				zap.String("user_agent", c.Request.UserAgent())).Info()
		}
	}
}
