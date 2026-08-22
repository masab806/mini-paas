package handlers

import (
	"bufio"
	"io"
	"mini-paas/internals/dto"
	"mini-paas/internals/services"
	"net/http"
	"strings"

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

func (h *LogHandler) StreamLogs(c *gin.Context) {
	containerName := c.Param("name")
	tail := c.DefaultQuery("tail", "50")

	streamReader, err := h.service.StreamContainerLogs(c.Request.Context(), containerName, tail)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	defer streamReader.Close()

	// Use event-stream so Postman & Browsers render lines as they arrive
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Transfer-Encoding", "chunked")

	scanner := bufio.NewScanner(streamReader)

	c.Stream(func(w io.Writer) bool {
		if scanner.Scan() {
			line := scanner.Text()
			line = strings.TrimRight(line, "\r")

			if len(strings.TrimSpace(line)) == 0 {
				return true
			}

			// Format as SSE data line to force Postman to display instantly
			c.Writer.WriteString("data: " + line + "\n\n")

			if flusher, ok := c.Writer.(http.Flusher); ok {
				flusher.Flush()
			}
			return true
		}
		return false
	})
}