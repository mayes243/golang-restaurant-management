package routes

import (
	controller "golang-restaurant-management/controllers"

	"github.com/gin-gonic/gin"
)

func FoodRoutes(router *gin.RouterGroup) {
	router.GET("/foods", controller.GetFoods())
	router.GET("/foods/:food_id", controller.GetFood())
	router.POST("/foods", controller.CreateFood())
	router.PATCH("/foods/:food_id", controller.UpdateFood())
}
