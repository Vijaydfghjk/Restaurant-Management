package routes

import (
	controller "restaurant-management/Controller"

	"github.com/gin-gonic/gin"
)

func UserRoute(incomingroutes *gin.Engine) {

	incomingroutes.POST("users/signup", controller.Signup)

	incomingroutes.POST("/users/login", controller.Login)

}
