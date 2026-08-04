package routes

import (
	"mini-paas/internals/handlers"
	"mini-paas/internals/middlewares"

	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine, userHandler *handlers.UserHandler) {
	api := r.Group("/api/user")

	api.POST("/create", userHandler.UserRegistration)
	api.POST("/login", userHandler.UserLogin)

	protected := api.Group("/")
	protected.Use(middlewares.AuthMiddleware())
	{
		protected.GET("/getProfile", userHandler.GetUserProfile)
	}

}
