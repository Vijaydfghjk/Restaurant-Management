package routes

import (
	"restaurant/controlles"

	"github.com/gin-gonic/gin"
)

func FoodRoutes(incomingroute *gin.Engine) {

	incomingroute.GET("/foods", controlles.GetFoods())
	incomingroute.GET("/foods/:food_id", controlles.Getfood())
	incomingroute.POST("/foods", controlles.Createfood())
	incomingroute.PATCH("/foods/:food_id", controlles.Updatefood())
}
