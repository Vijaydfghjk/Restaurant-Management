package main

import (
	routes "restaurant-management/Routes"

	"github.com/gin-gonic/gin"
)

func main() {

	app := gin.Default()

	routes.FoodRoutes(app)
	routes.InvoiceRoutes(app)
	routes.MenuRoutes(app)
	routes.OrderItemRouter(app)
	routes.OrderRoutes(app)
	routes.TableRoutes(app)

	app.Run(":8000")
}
