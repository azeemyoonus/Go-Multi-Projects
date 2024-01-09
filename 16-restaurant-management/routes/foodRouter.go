package routes

import (
	"github.com/amitamrutiya/10-go-authentication-jwt/middleware"
	controller "github.com/amitamrutiya/16-restaurant-management/controllers"
	"github.com/gin-gonic/gin"
)

// FoodRoutes : Routing for food

func FoodRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middleware.Authenticate())

	incomingRoutes.GET("/food", controller.GetFoods())
	incomingRoutes.GET("/food/:food_id", controller.GetFood())
	incomingRoutes.POST("/food", controller.CreateFood())
	incomingRoutes.PUT("/food/:food_id", controller.UpdateFood())
	incomingRoutes.DELETE("/food/:food_id", controller.DeleteFood())
}
