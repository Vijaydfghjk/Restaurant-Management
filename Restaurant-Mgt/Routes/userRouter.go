package routes

import (
	"restaurant/controlles"

	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingroutes *gin.Engine) {

	incomingroutes.GET("/users", controlles.GetUsers())
	incomingroutes.GET("/users/:user_id", controlles.Getuser())
	incomingroutes.POST("/users/signup", controlles.Signup())
	incomingroutes.POST("/users/login", controlles.Login())
}
