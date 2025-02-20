package routes

import (
	controller "restaurant-management/Controller"

	"github.com/gin-gonic/gin"
)

func FoodRoutes(incomingroute *gin.Engine) {

	food_handler := controller.Food_controll()

	incomingroute.GET("/foods", food_handler.GetFoods)
	incomingroute.GET("/foods/:food_id", food_handler.GetFood)
	incomingroute.POST("/foods", food_handler.Createfood)
	incomingroute.PATCH("/foods/:food_id", food_handler.Updatefood)

}
