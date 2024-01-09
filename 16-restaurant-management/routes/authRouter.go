package routes

import (
	controller "github.com/amitamrutiya/16-restaurant-management/controllers"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/auth/signup", controller.SignUp())
	incomingRoutes.POST("/auth/login", controller.Login())

}
