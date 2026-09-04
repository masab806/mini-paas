package handlers

import (
	"mini-paas/internals/config"
	"mini-paas/internals/dto"
	"mini-paas/internals/services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

	value, exists := c.Get("claims")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Error": "Unauthorized",
		})
	}

	randomNum := uuid.NewString()

	claims := value.(*config.Claims)

	containerName := claims.Username + randomNum

	sanitizedName := strings.ReplaceAll(containerName, " ", "-")
	
	deploymentDetails, publicUrl, deployedPath, err := h.service.CloneRepository(c.Request.Context(), &req.RepoURL, &req.Branch, &req.ImageTag, &req.Framework, &sanitizedName, &req.Port, &claims.ID)


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
		"Details": deploymentDetails,
		"Public Url": publicUrl,
	})

}

