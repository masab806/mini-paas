package routes

import (
	"mini-paas/internals/handlers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine, userHandler *handlers.UserHandler){
	api := r.Group("/api/user")
	{
		api.POST("/create", userHandler.UserRegistration)
		api.POST("/login", userHandler.UserLogin)
		api.GET("/getProfile", userHandler.GetUserProfile)
	}
}