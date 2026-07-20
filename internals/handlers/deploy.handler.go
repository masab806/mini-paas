package handlers

import (
	"mini-paas/internals/dto"
	"mini-paas/internals/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DeployHandler struct {
	service services.DeployService
}

func NewDeployHandler(service *services.DeployService) *DeployHandler {
	return &DeployHandler{service: *service}
}

func (h *DeployHandler) DeployToServer(c *gin.Context) {
	var req dto.DeployRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	

	deployedPath, err := h.service.CloneRepository(c.Request.Context(), req.RepoURL, req.Branch)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Message": "Internal Server Error!",
			"Error": err.Error(),
		})

		return
	}

	c.JSON(200, gin.H{
		"Message": "Deployed!",
		"Path": deployedPath,
	})

}

