package routes

import (
	"mini-paas/internals/handlers"
	"mini-paas/internals/middlewares"

	"github.com/gin-gonic/gin"
)

func DeployRoutes(r *gin.Engine, deployHandler *handlers.DeployHandler) {
	api := r.Group("/api/deploy")
	protected := api.Group("/")
	protected.Use(middlewares.AuthMiddleware())
	{
		protected.POST("/", deployHandler.DeployToServer)
	}
}