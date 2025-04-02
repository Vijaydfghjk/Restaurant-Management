package routes

import (
	"restaurant/controlles"

	"github.com/gin-gonic/gin"
)

func TableRoutes(incomingroutes *gin.Engine) {

	incomingroutes.GET("/tables", controlles.Gettables())
	incomingroutes.GET("/tables/:table_id", controlles.GetTable())
	incomingroutes.POST("/tables", controlles.CreateTable())
	incomingroutes.PATCH("/tables/:table_id", controlles.UpdateTable())

}
