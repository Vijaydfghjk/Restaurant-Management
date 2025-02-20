package routes

import (
	controller "restaurant-management/Controller"

	"github.com/gin-gonic/gin"
)

func OrderItemRouter(incomingRoutes *gin.Engine) {

	orderItem_handler := controller.Orderitemcontroll()

	incomingRoutes.GET("/OrderItems", orderItem_handler.GetOrderItems)
	incomingRoutes.GET("/OrderItems/:order_item_id", orderItem_handler.GetorderItem)
	incomingRoutes.POST("/OrderItems", orderItem_handler.CreateOrderItem)
	incomingRoutes.PATCH("/OrderItems/:order_item_id", orderItem_handler.UpdateOrderitem)
	incomingRoutes.GET("/OrderItemsbyorder_id/:order_id", orderItem_handler.GetOrderItemsByOrder)
}
