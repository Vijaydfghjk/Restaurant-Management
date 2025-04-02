package routes

import (
	"restaurant/controlles"

	"github.com/gin-gonic/gin"
)

func OrderRoutes(incomingroutes *gin.Engine) {

	incomingroutes.GET("/orders", controlles.Getorders())
	incomingroutes.GET("/orders/:order_id", controlles.Getorder())
	incomingroutes.POST("/orders", controlles.CreateOrder())
	incomingroutes.PATCH("/orders/:order_id", controlles.Updateorder())
}
