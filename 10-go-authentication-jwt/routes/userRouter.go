package routes

import (
	controllers "github.com/amitamrutiya/10-go-authentication-jwt/controllers"
	"github.com/amitamrutiya/10-go-authentication-jwt/middleware"
	"github.com/gin-gonic/gin"
)

// UserRoutes function
func UserRoutes(route *gin.Engine) {
	// use middelware
	route.Use(middleware.Authenticate())
	user := route.Group("/user")
	{
		user.GET("/", controllers.GetUsers())
		user.GET("/:id", controllers.GetUser())
		// user.PUT("/:id", controllers.UpdateUser)
		// user.DELETE("/:id", controllers.DeleteUser)
	}
}
