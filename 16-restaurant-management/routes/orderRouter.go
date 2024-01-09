package routes

import (
	"github.com/amitamrutiya/10-go-authentication-jwt/middleware"
	controller "github.com/amitamrutiya/16-restaurant-management/controllers"
	"github.com/gin-gonic/gin"
)

// OrderRoutes : Routing for order

func OrderRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middleware.Authenticate())

	incomingRoutes.GET("/order", controller.GetOrders())
	incomingRoutes.GET("/order/:order_id", controller.GetOrder())
	incomingRoutes.POST("/order", controller.CreateOrder())
	incomingRoutes.PUT("/order/:order_id", controller.UpdateOrder())
	incomingRoutes.DELETE("/order/:order_id", controller.DeleteOrder())
}
