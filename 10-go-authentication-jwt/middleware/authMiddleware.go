package middleware

import (
	"fmt"
	"net/http"

	helper "github.com/amitamrutiya/10-go-authentication-jwt/helpers"
	"github.com/gin-gonic/gin"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("token")
		if clientToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Errorf("Missing auth token")})
			c.Abort()
			return
		}
		claims, err := helper.ValidateToken(clientToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		c.Set("uid", claims.ID)
		c.Set("email", claims.Email)
		c.Set("username", claims.Username)
		c.Set("user_type", claims.User_type)
		c.Next()
	}
}
