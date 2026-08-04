package middlewares

import (
	"mini-paas/internals/config"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authheader := c.GetHeader("Authorization")

		if authheader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"Error": "Invalid Token!",
			})

			c.Abort()
			return
		}

		const bearerPrefix = "Bearer "

		if !strings.HasPrefix(authheader, bearerPrefix){
			c.JSON(http.StatusUnauthorized, gin.H{
				"Error": "Invalid Token!",
			})

			c.Abort()

			return
		}

		token := strings.TrimPrefix(authheader, bearerPrefix)


		config.ValidateToken(token)

		c.Next()

	}
}