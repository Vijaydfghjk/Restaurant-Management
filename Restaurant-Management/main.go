package main

import (
	middleware "restaurant-management/Middleware"
	routes "restaurant-management/Routes"

	"github.com/gin-gonic/gin"
)

func main() {

	app := gin.New()

	app.Use(gin.Logger())

	routes.UserRoute(app)
	app.Use(middleware.Authentication)

	routes.FoodRoutes(app)
	routes.InvoiceRoutes(app)
	routes.MenuRoutes(app)
	routes.OrderItemRouter(app)
	routes.OrderRoutes(app)
	routes.TableRoutes(app)

	app.Run(":8000") // :8000
}
