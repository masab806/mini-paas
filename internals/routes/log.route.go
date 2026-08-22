package routes

import (
	"mini-paas/internals/handlers"

	"github.com/gin-gonic/gin"
)

func LogRoutes(r *gin.Engine, logHandler *handlers.LogHandler){
	api := r.Group("/api/logs")
	{
		api.POST("/getLogs", logHandler.GetAllLogs)
		api.GET("/apps/:name/logs/stream", logHandler.StreamLogs)
	}
}