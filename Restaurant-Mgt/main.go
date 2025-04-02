package main

import (
	middlware "restaurant/Middlware"
	routes "restaurant/Routes"

	"github.com/gin-gonic/gin"
)

func main() {

	app := gin.New()

	app.Use(gin.Logger())

	routes.UserRoutes(app)

	app.Use(middlware.Authentication())

	routes.FoodRoutes(app)
	routes.MenuRoutes(app)
	routes.TableRoutes(app)
	routes.InvoiceRoutes(app)
	routes.OrderRoutes(app)
	routes.OrderItemRouter(app)
	//routes.UserRoutes(app)

	app.Run(":8080")

}
