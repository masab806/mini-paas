package handlers

import (
	"fmt"
	"mini-paas/internals/dto"
	"mini-paas/internals/services"
	"net/http"
	"strings"

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

func (h *UserHandler) UserLogin(c *gin.Context){
	var req dto.UserLoginRequest

	fmt.Println(req.Email, req.Password)

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	token, err := h.service.LoginUser(c.Request.Context(), req.Email, req.Password)

	if err != nil {
		 c.JSON(400, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(200, gin.H{
		"Message": "User Logged In!",
		"token": token,
	})
}

func (h *UserHandler) GetUserProfile(c *gin.Context){
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" {
		c. JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid Token",
		})

		return
	}

	const bearerPrefix = "Bearer "

	if !strings.HasPrefix(authHeader, bearerPrefix ){
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid Token",
		})

		return
	}

	token := strings.TrimPrefix(authHeader, bearerPrefix)

	profile, err := h.service.GetProfile(c.Request.Context(), token)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid Token",
		})

		return
	}

	payload := gin.H{
		"userId":   profile.ID,
		"email":    profile.Email,
		"username": profile.Username,
	}

	c.JSON(http.StatusOK, payload)


}
