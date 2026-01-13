package http

import (
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	_ "github.com/RusselRustCode/teacher_analytics/core-service/docs" 
)

func SetupRoutes(router *gin.Engine, handler *HTTPHandler) {
	api := router.Group("/api")
	{
		api.POST("/log", handler.SendLog)
		api.GET("/analytics/:student_id", handler.GetAnalytics)
		api.GET("/students/:student_id/logs", handler.GetStudentLogs)
        api.GET("/students", handler.GetStudents)
	}
	router.GET("/ping-swagger", func(c *gin.Context) {
		c.String(200, "Router is working")
	})
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}