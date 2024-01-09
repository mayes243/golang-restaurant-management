package routes

import (
	controller "golang-restaurant-management/controllers"

	"github.com/gin-gonic/gin"
)

func MenuRoutes(router *gin.RouterGroup) {
	router.GET("/menus", controller.GetMenus())
	router.GET("/menus/:menu_id", controller.GetMenu())
	router.POST("/menus", controller.CreateMenu())
	router.PATCH("/menus/:menu_id", controller.UpdateMenu())
}
