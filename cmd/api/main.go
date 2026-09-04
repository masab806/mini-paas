package main

import (
	"context"
	"log"
	"os"

	"mini-paas/internals/config"
	"mini-paas/internals/database"
	"mini-paas/internals/handlers"
	"mini-paas/internals/repositories"
	"mini-paas/internals/routes"
	"mini-paas/internals/services"
	"mini-paas/internals/utils"

	"github.com/gin-gonic/gin"
)

func main() {
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

	containerService := services.NewContainerService()
	containerHandler := handlers.NewContainerHandler(containerService)

	userRepo := repositories.NewUserRepository(client)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	mailService := services.NewMailService()

	AIService, err := services.NewAIService(context.Background(), os.Getenv("GEMINI_KEY"), userRepo, mailService)
	if err != nil {
		log.Fatal(err)
	}

	if err := utils.SetupNginx(); err != nil {
		log.Fatalf(
			"failed to setup nginx: %v",
			err,
		)
	}

	AIHandler := handlers.NewAIHandler(AIService, logService)

	routes.DeployRoutes(r, deployHandler)
	routes.UserRoutes(r, userHandler)
	routes.LogRoutes(r, logHandler)
	routes.ContainerRoutes(r, containerHandler)
	routes.AIRoutes(r, AIHandler)

	log.Fatal(r.Run(":8000"))
}
