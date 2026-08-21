package handlers

import (
	"mini-paas/internals/dto"
	"mini-paas/internals/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LogHandler struct {
	service services.LogService
}

func NewLogHandler(service *services.LogService) *LogHandler {
	return &LogHandler{
		service: *service,
	}
}

func (h *LogHandler) GetAllLogs(c *gin.Context) {
	var req dto.LogRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})

		return
	}

	result, err := h.service.GetContainerLogs(c.Request.Context(), req.ContainerName, req.LineNo)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
	}

	c.JSON(200, gin.H{
		"Result": result,
	})
}