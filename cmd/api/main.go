package main

import (
	"context"
	"log"
	"mini-paas/internals/config"
	"mini-paas/internals/database"
	"mini-paas/internals/handlers"
	"mini-paas/internals/repositories"
	"mini-paas/internals/routes"
	"mini-paas/internals/services"
	"os"

	"github.com/gin-gonic/gin"
)

func main(){
	config.LoadConfig()

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Working",
		})
	})

	client, err := database.Connect(os.Getenv("DATABASE_URL"))

	if err != nil {
		log.Fatal(err)
	}

	deploymentRepo := repositories.NewDeploymentRepository(client)
	deployService := services.NewDeployService(client, deploymentRepo)
	deployHandler := handlers.NewDeployHandler(deployService)

	logService := services.NewLogService()
	logHandler := handlers.NewLogHandler(logService)

	AIService, err := services.NewAIService(context.Background(), os.Getenv("GEMINI_KEY"))
	if err != nil {
		log.Fatal(err)
	}
	
	AIHandler := handlers.NewAIHandler(AIService, logService)

	containerService := services.NewContainerService()
	containerHandler := handlers.NewContainerHandler(containerService)

	userRepo := repositories.NewUserRepository(client)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	routes.DeployRoutes(r, deployHandler)
	routes.UserRoutes(r, userHandler)
	routes.LogRoutes(r, logHandler)
	routes.ContainerRoutes(r, containerHandler)
	routes.AIRoutes(r, AIHandler)

	log.Fatal(r.Run(":8000"))

}