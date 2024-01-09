package routes

import (
	controller "golang-restaurant-management/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.RouterGroup) {

	// group
	v1UserRoute := router.Group("/v1")

	v1UserRoute.GET("/users", controller.GetUsers())
	v1UserRoute.GET("/users/:user_id", controller.GetUser())
	v1UserRoute.POST("/users/signup", controller.SignUp())
	v1UserRoute.POST("/users/login", controller.Login())

}
