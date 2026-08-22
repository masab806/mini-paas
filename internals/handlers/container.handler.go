package handlers

import (
	"mini-paas/internals/dto"
	"mini-paas/internals/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ContainerHandler struct {
	service services.ContainerService
}

func NewContainerHandler(service *services.ContainerService) *ContainerHandler {
	return &ContainerHandler{
		service: *service,
	}
}

func (h ContainerHandler) RunContainer(c *gin.Context){
	var req dto.ContainerRequest

	err := c.ShouldBindJSON(&req); if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})

		return
	}

	containerName, containerErr := h.service.RunContainer(c.Request.Context(), req.ContainerName)

	if containerErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": containerErr.Error(),
		})
	}

	c.JSON(200, gin.H{
		"Container Started": containerName,
	})
}

func (h ContainerHandler) StopContainer(c *gin.Context){
	var req dto.ContainerRequest

	err := c.ShouldBindJSON(&req); if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
	}

	containerName, containerErr := h.service.StopContainer(c.Request.Context(), req.ContainerName)

	if containerErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": containerErr.Error(),
		})
	}

	c.JSON(200, gin.H{
		"Container Stopped": containerName,
	})
}

func (h ContainerHandler) DeleteContainer(c *gin.Context){
	var req dto.ContainerRequest

	err := c.ShouldBindJSON(&req); if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
	}

	containerName, containerErr := h.service.DeleteContainer(c.Request.Context(), req.ContainerName)

	if containerErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": containerErr.Error(),
		})
	}

	c.JSON(200, gin.H{
		"Container Deleted": containerName,
	})


}