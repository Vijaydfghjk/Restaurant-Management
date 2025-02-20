package routes

import (
	controller "restaurant-management/Controller"

	"github.com/gin-gonic/gin"
)

func TableRoutes(incomingroutes *gin.Engine) {

	table_handler := controller.Tablecontroll()

	incomingroutes.GET("/tables", table_handler.Get_tables)
	incomingroutes.GET("/tables/:table_id", table_handler.Get_table)
	incomingroutes.POST("/tables", table_handler.Create_table)
	incomingroutes.PATCH("tables/:table_id", table_handler.Update_table)
}
