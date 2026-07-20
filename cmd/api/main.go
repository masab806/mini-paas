package main

import (
	"log"
	"mini-paas/internals/handlers"
	"mini-paas/internals/routes"
	"mini-paas/internals/services"

	"github.com/gin-gonic/gin"
)

func main(){
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Working",
		})
	})

	deployService := services.NewDeployService()
	deployHandler := handlers.NewDeployHandler(deployService)

	routes.DeployRoutes(r, deployHandler)

	log.Fatal(r.Run(":8000"))

}