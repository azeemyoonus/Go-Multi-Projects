package main

import (
	"os"

	"github.com/amitamrutiya/10-go-authentication-jwt/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	router := gin.New()
	router.Use(gin.Logger())

	routes.AuthRoutes(router)
	routes.UserRoutes(router)

	// router.Get("/api", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{
	// 		"message": "Welcome to the API v1",
	// 	})
	// })

	router.Run(":" + port)
}
