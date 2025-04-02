package routes

import (
	"restaurant/controlles"

	"github.com/gin-gonic/gin"
)

func OrderItemRouter(incomingRoutes *gin.Engine) {

	incomingRoutes.GET("/OrderItems", controlles.GetorderItems())
	incomingRoutes.GET("/OrderItems/:order_item_id", controlles.GetOrderItem())
	incomingRoutes.GET("/OrderItems-order/:order_id", controlles.GetOrderItemsByOrder())
	incomingRoutes.POST("/OrderItems", controlles.CreateOrderItem())
	incomingRoutes.PATCH("OrderItems/:order_item_id", controlles.UpdateOrderItem())

}
