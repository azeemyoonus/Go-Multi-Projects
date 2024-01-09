package routes

import (
	"github.com/amitamrutiya/10-go-authentication-jwt/middleware"
	controller "github.com/amitamrutiya/16-restaurant-management/controllers"
	"github.com/gin-gonic/gin"
)

// OrderItemRoutes : Routing for orderItem

func OrderItemRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middleware.Authenticate())

	incomingRoutes.GET("/orderItem", controller.GetOrderItems())
	incomingRoutes.GET("/orderItem/:orderItem_id", controller.GetOrderItem())
	incomingRoutes.GET("/orderItem/order/:order_id", controller.GetOrderItemsByOrder())
	incomingRoutes.POST("/orderItem", controller.CreateOrderItem())
	incomingRoutes.PUT("/orderItem/:orderItem_id", controller.UpdateOrderItem())
	incomingRoutes.DELETE("/orderItem/:orderItem_id", controller.DeleteOrderItem())
}
