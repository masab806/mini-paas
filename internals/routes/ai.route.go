package routes

import (
	"mini-paas/internals/handlers"

	"github.com/gin-gonic/gin"
)

func AIRoutes(r *gin.Engine, AIHandler *handlers.AIHandler){
	api := r.Group("/api/ai")
	{
		api.POST("/analyze", AIHandler.AnalyzeLogs)
	}
}