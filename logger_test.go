package log

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"testing"
)

func TestLogger(t *testing.T) {
	logger := NewLogger(&Config{
		Color:       false,
		Writers:     "file,stdout",
		Development: false,
	})
	logger.Debug("debug")
	logger.Info("info")
	logger.Warn("warn")
	logger.Error("error")
	r := gin.New()
	r.Use(AccessLoggerWithConfigFile("./example.config.yaml"))
	r.GET("/", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{"hello": "gin"})
	})
	r.Run(":8080")
}
