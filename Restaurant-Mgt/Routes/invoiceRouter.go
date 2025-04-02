package routes

import (
	"restaurant/controlles"

	"github.com/gin-gonic/gin"
)

func InvoiceRoutes(incomingroutes *gin.Engine) {

	incomingroutes.GET("/invoices", controlles.GetInvoices())
	incomingroutes.GET("/invoices/:invoice_id", controlles.GetInvoice())
	incomingroutes.POST("/invoices", controlles.CreateInvoice())
	incomingroutes.PATCH("/invoices/:invoice_id", controlles.UpdateInvoice())

}
