package routes

import (
	controller "golang-restaurant-management/controllers"

	"github.com/gin-gonic/gin"
)

func OrderRoutes(router *gin.RouterGroup) {
	router.GET("/orders", controller.GetOrders())
	router.GET("/orders/:order_id", controller.GetOrder())
	router.POST("/orders", controller.CreateOrder())
	router.PATCH("/orders/:order_id", controller.UpdateOrder())
}
