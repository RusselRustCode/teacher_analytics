package http

import (
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, handler *HTTPHandler) {
	api := router.Group("/api")
	{
		api.POST("/log", handler.SendLog)
		api.GET("/analytics/:id", handler.GetAnalytics)
	}
}