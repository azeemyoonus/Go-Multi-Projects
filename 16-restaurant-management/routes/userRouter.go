package routes

import (
	controller "github.com/amitamrutiya/16-restaurant-management/controllers"
	"github.com/amitamrutiya/16-restaurant-management/middleware"
	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middleware.Authentication())
	incomingRoutes.GET("/users", controller.GetUsers())
	incomingRoutes.GET("/users/:user_id", controller.GetUser())
	incomingRoutes.DELETE("/users/:user_id", controller.DeleteUser())
}
