package controlles

import (
	"context"
	"fmt"
	"log"
	"net/http"
	database "restaurant/Database"
	models "restaurant/Models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

var ordercollection *mongo.Collection = database.OpenCollection(database.Client, "order")

func Getorders() gin.HandlerFunc {

	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		orders, err := ordercollection.Find(ctx, bson.M{})

		if err != nil {

			c.JSON(http.StatusInternalServerError, gin.H{"error": "unbale to fetch the data"})
			return
		}

		var datas = []bson.M{}
		for orders.Next(ctx) {

			var data bson.M

			if err := orders.Decode(&data); err != nil {
				log.Fatal(err)
			}
			datas = append(datas, data)
		}
		defer cancel()
		c.JSON(http.StatusOK, datas)
	}
}

func Getorder() gin.HandlerFunc {

	return func(c *gin.Context) {

		orderid := c.Param("order_id")

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var order models.Order
		err := ordercollection.FindOne(ctx, bson.M{"order_id": orderid}).Decode(&order)
		defer cancel()
		if err != nil {

			c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to get the data"})
			return
		}

		c.JSON(http.StatusOK, order)
	}
}

func CreateOrder() gin.HandlerFunc {

	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var order models.Order
		var table models.Table

		if err := c.ShouldBindJSON(&order); err != nil {

			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if valerr := validate.Struct(order); valerr != nil {

			c.JSON(http.StatusBadRequest, gin.H{"error": valerr.Error()})
			return
		}
		if order.Table_id != nil {
			err := tablecollection.FindOne(ctx, bson.M{"table_id": order.Table_id}).Decode(&table)
			defer cancel()

			if err != nil {

				msg := fmt.Sprintf("message:Table was not found")
				c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
				return
			}

		}

		order.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		order.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		order.ID = primitive.NewObjectID()
		order.Order_id = order.ID.Hex()

		newdata, err := ordercollection.InsertOne(ctx, order)

		if err != nil {

			c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to insert the data"})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, newdata)
	}
}

func Updateorder() gin.HandlerFunc {

	return func(c *gin.Context) {

		orderid := c.Param("order_id")

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var order models.Order
		var table models.Table

		if err := c.ShouldBindJSON(&order); err != nil {

			c.JSON(http.StatusBadRequest, gin.H{"error": "Badrequest"})
			return
		}

		var updateobj primitive.D

		/*
			if order.Table_id != nil {

				updateobj = append(updateobj, bson.E{"table_id", order.Order_id})
			}
		*/

		if order.Table_id != nil {

			err := tablecollection.FindOne(ctx, bson.M{"table_id": order.Table_id}).Decode(&table)
			defer cancel()
			if err != nil {

				msg := fmt.Sprintf("message:Menu was not found")
				c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
				return
			}
			updateobj = append(updateobj, bson.E{"table_id", order.Table_id})
		}

		order.Order_Date, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateobj = append(updateobj, bson.E{"updated_at", order.Order_Date})

		filter := bson.M{"order_id": orderid}

		upsert := true

		opt := options.UpdateOptions{

			Upsert: &upsert,
		}

		result, err := ordercollection.UpdateOne(
			ctx,
			filter,
			bson.D{

				{"$set", updateobj},
			},
			&opt,
		)
		if err != nil {

			msg := fmt.Sprintf("order item update failed")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

func OrderItemOrderCreator(order models.Order) string {

	order.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.ID = primitive.NewObjectID()
	order.Order_id = order.ID.Hex()
	orderItemCollection.InsertOne(ctx, order)

	defer cancel()

	return order.Order_id
}
