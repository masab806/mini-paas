package routes

import (
	"mini-paas/internals/handlers"

	"github.com/gin-gonic/gin"
)

func ContainerRoutes(r *gin.Engine, ContainerHandler *handlers.ContainerHandler){
	api := r.Group("/api/container")
	{
		api.POST("/run", ContainerHandler.RunContainer)
		api.POST("/stop", ContainerHandler.StopContainer)
		api.POST("/delete", ContainerHandler.DeleteContainer)
	}
}

