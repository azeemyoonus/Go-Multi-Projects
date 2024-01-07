package routes

import (
	controllers "github.com/amitamrutiya/10-go-authentication-jwt/controllers"
	"github.com/gin-gonic/gin"
)

// AuthRoutes function
func AuthRoutes(route *gin.Engine) {
	auth := route.Group("/auth")
	{
		auth.POST("/login", controllers.Login())
		auth.POST("/register", controllers.Singup())
	}
}