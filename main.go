package main

import (
	"os"
	"strings"

	_ "golang-restaurant-management/docs"
	middleware "golang-restaurant-management/middleware"
	routes "golang-restaurant-management/routes"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title 	Restaurant Service API
// @version	1.0
// @description A Restaurant Service API in Go using Gin framework
// @host 	localhost:4000
// @BasePath /api
func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	router := gin.New()
	router.ForwardedByClientIP = true

	trustedProxies := strings.Split(os.Getenv("TRUSTED_PROXIES"), ",")
	router.SetTrustedProxies(trustedProxies)

	// prefix for all routes
	apiGroup := router.Group("/api")

	apiGroup.GET("/docs/*.html", ginSwagger.WrapHandler(swaggerFiles.Handler))

	apiGroup.Use(gin.Logger())
	routes.UserRoutes(apiGroup)
	apiGroup.Use(middleware.Authentication())

	// routes.FoodRoutes(apiGroup)
	routes.MenuRoutes(apiGroup)
	routes.TableRoutes(apiGroup)
	routes.OrderRoutes(apiGroup)
	routes.OrderItemRoutes(apiGroup)
	routes.InvoiceRoutes(apiGroup)

	router.Run(":" + port)
}
