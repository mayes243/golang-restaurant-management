package routes

import (
	controller "golang-restaurant-management/controllers"

	"github.com/gin-gonic/gin"
)

func TableRoutes(router *gin.RouterGroup) {
	router.GET("/tables", controller.GetTables())
	router.GET("/tables/:table_id", controller.GetTable())
	router.POST("/tables", controller.CreateTable())
	router.PATCH("/tables/:table_id", controller.UpdateTable())
}
