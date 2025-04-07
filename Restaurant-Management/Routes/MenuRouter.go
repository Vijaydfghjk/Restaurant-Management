package routes

import (
	controller "restaurant-management/Controller"

	"github.com/gin-gonic/gin"
)

func MenuRoutes(incomingroutes *gin.Engine) {

	menu_handler := controller.Menu_controll()

	incomingroutes.GET("/menus", menu_handler.Getmenus)
	incomingroutes.GET("/menus/:menu_id", menu_handler.GetmenubyId)
	incomingroutes.POST("/menus", menu_handler.CreateMenu)
	incomingroutes.PATCH("/menus/:menu_id", menu_handler.Updatemenu)

}
