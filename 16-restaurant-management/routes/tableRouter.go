package routes

import (
	"github.com/amitamrutiya/10-go-authentication-jwt/middleware"
	controller "github.com/amitamrutiya/16-restaurant-management/controllers"
	"github.com/gin-gonic/gin"
)

// TableRoutes : Routing for table

func TableRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middleware.Authenticate())

	incomingRoutes.GET("/table", controller.GetTables())
	incomingRoutes.GET("/table/:table_id", controller.GetTable())
	incomingRoutes.POST("/table", controller.CreateTable())
	incomingRoutes.PUT("/table/:table_id", controller.UpdateTable())
	incomingRoutes.DELETE("/table/:table_id", controller.DeleteTable())
}
