package routes

import (
	"mini-paas/internals/handlers"

	"github.com/gin-gonic/gin"
)

func DeployRoutes(r *gin.Engine, deployHandler *handlers.DeployHandler) {
	api := r.Group("/api/deploy")
	{
		api.POST("/", deployHandler.DeployToServer)
	}
}