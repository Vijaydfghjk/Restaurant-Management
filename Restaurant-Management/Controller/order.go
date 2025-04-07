package controller

import (
	"context"
	"log"
	"net/http"
	database "restaurant-management/Database"
	model "restaurant-management/Model"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Order_db struct {
	ordercollection *mongo.Collection
	validate        *validator.Validate
}

func Ordercontroll() *Order_db {

	return &Order_db{
		ordercollection: database.Open_collection(database.Client, "Order"),
		validate:        validator.New(),
	}
}

func (a *Order_db) Create_order(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	defer cancel()
	var order model.Order

	if err := c.ShouldBindJSON(&order); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if valerror := a.validate.Struct(order); valerror != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": valerror.Error()})
		return
	}

	order.ID = primitive.NewObjectID()
	order.Order_id = order.ID.Hex()
	order.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	result, err := a.ordercollection.InsertOne(ctx, order)

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (a *Order_db) Getorders(c *gin.Context) {

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	defer cancel()

	results, err := a.ordercollection.Find(ctx, bson.M{})

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var orders []bson.M
	for results.Next(ctx) {
		var order bson.M
		if err := results.Decode(&order); err != nil {
			log.Fatal(err)
			return
		}
		orders = append(orders, order)
	}
	c.JSON(http.StatusOK, orders)
}

func (a *Order_db) Getorderbyid(c *gin.Context) {

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	var order model.Order

	order_id := c.Param("order_id")
	defer cancel()

	err := a.ordercollection.FindOne(ctx, bson.M{"order_id": order_id}).Decode(&order)

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

func (a *Order_db) UpdateOrder(c *gin.Context) {

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	defer cancel()

	order_id := c.Param("order_id")

	myfilter := bson.M{"order_id": order_id}

	var order model.Order

	if err := c.ShouldBind(&order); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	if valerr := a.validate.Struct(order); valerr != nil {

		c.JSON(http.StatusBadRequest, gin.H{"Error": valerr.Error()})
		return
	}

	var updateOrder primitive.D

	if order.Table_id != nil {

		updateOrder = append(updateOrder, bson.E{"table_id", order.Table_id})

	}

	order.Order_Date, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateOrder = append(updateOrder, bson.E{"order_date", order.Order_Date})

	order.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	upsert := true
	opt := options.UpdateOptions{

		Upsert: &upsert,
	}

	updated_data, err := a.ordercollection.UpdateOne(

		ctx,
		myfilter,
		bson.D{

			{"$set", updateOrder},
		},
		&opt,
	)

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updated_data)
}
