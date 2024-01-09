package middleware

import (
	"fmt"
	"net/http"

	helper "github.com/amitamrutiya/16-restaurant-management/helpers"
	"github.com/gin-gonic/gin"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Printf("Inside Authentication middleware\n")
		clientToken := c.Request.Header.Get("token")
		if clientToken == "" {
			fmt.Printf("Missing auth token\n")
			c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Errorf("Missing auth token")})
			c.Abort()
			return
		}
		claims, err := helper.ValidateToken(clientToken)
		if err != nil {
			fmt.Printf("Invalid auth token\n")
			c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Errorf("Invalid auth token")})
			c.Abort()
			return
		}
		c.Set("uid", claims.Uid)
		c.Set("email", claims.Email)
		c.Set("first_name", claims.First_name)
		c.Set("last_name", claims.Last_name)

		c.Next()
	}
}
