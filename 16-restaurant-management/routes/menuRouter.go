package routes

import (
	"github.com/amitamrutiya/10-go-authentication-jwt/middleware"
	controller "github.com/amitamrutiya/16-restaurant-management/controllers"
	"github.com/gin-gonic/gin"
)

// MenuRoutes : Routing for menu

func MenuRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middleware.Authenticate())

	incomingRoutes.GET("/menu", controller.GetMenus())
	incomingRoutes.GET("/menu/:menu_id", controller.GetMenu())
	incomingRoutes.POST("/menu", controller.CreateMenu())
	incomingRoutes.PUT("/menu/:menu_id", controller.UpdateMenu())
	incomingRoutes.DELETE("/menu/:menu_id", controller.DeleteMenu())
}
