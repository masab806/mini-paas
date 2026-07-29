package handlers

import (
	"mini-paas/internals/dto"
	"mini-paas/internals/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service services.UserService
}

func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{
		service: *service,
	}
}


func (h *UserHandler) UserRegistration(c *gin.Context){
	var req dto.UserRequestType

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	newUser, err := h.service.RegisterUser(c.Request.Context(), req.Email, req.Username, req.Password)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(200, gin.H{
		"message": "User Registered",
		"newUser": newUser,
	})


}
