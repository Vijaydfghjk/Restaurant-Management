package routes

import (
	"restaurant/controlles"

	"github.com/gin-gonic/gin"
)

func MenuRoutes(incomingroutes *gin.Engine) {

	incomingroutes.GET("/menus", controlles.GetMenus())
	incomingroutes.GET("/menus/:menu_id", controlles.Getmenu())
	incomingroutes.POST("/menus", controlles.CreateMenu())
	incomingroutes.PATCH("/menus/:menu_id", controlles.UpdateMenu())
}
