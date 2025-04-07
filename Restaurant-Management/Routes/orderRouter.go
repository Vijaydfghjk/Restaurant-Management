package routes

import (
	controller "restaurant-management/Controller"

	"github.com/gin-gonic/gin"
)

func OrderRoutes(incomingroutes *gin.Engine) {

	order_handler := controller.Ordercontroll()

	incomingroutes.GET("/orders", order_handler.Getorders)
	incomingroutes.GET("/orders/:order_id", order_handler.Getorderbyid)
	incomingroutes.POST("/orders", order_handler.Create_order)
	incomingroutes.PATCH("/orders/:order_id", order_handler.UpdateOrder)

}
