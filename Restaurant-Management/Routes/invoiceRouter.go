package routes

import (
	controller "restaurant-management/Controller"

	"github.com/gin-gonic/gin"
)

func InvoiceRoutes(incomingroutes *gin.Engine) {

	invoice_handler := controller.Invoiceccontrll()

	incomingroutes.GET("/invoices", invoice_handler.GetInvoices)
	incomingroutes.GET("/invoices/:invoice_id", invoice_handler.Getinvoice)
	incomingroutes.POST("/invoices", invoice_handler.Create_inoice)
	incomingroutes.PATCH("/invoices/:invoice_id", invoice_handler.Update_invoice)

}
 