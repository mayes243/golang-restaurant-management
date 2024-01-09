package routes

import (
	controller "golang-restaurant-management/controllers"

	"github.com/gin-gonic/gin"
)

func InvoiceRoutes(router *gin.RouterGroup) {
	router.GET("/invoices", controller.GetInvoices())
	router.GET("/invoices/:invoice_id", controller.GetInvoice())
	router.POST("/invoices", controller.CreateInvoice())
	router.PATCH("/invoices/:invoice_id", controller.UpdateInvoice())
}
