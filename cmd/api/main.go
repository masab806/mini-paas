package main

import (
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

	userRepo := repositories.NewUserRepository(client)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	routes.DeployRoutes(r, deployHandler)
	routes.UserRoutes(r, userHandler)
	routes.LogRoutes(r, logHandler)

	log.Fatal(r.Run(":8000"))

}