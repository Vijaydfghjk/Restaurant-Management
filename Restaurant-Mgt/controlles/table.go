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

var tablecollection *mongo.Collection = database.OpenCollection(database.Client, "table")

func Gettables() gin.HandlerFunc {

	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		result, err := tablecollection.Find(ctx, bson.M{})

		defer cancel()

		if err != nil {

			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing table items"})
			return
		}

		var Datas []bson.M

		for result.Next(ctx) {

			var data bson.M

			if err := result.Decode(&data); err != nil {

				log.Fatal(err)
			}

			Datas = append(Datas, data)
		}

		c.JSON(http.StatusOK, Datas)
	}
}

func GetTable() gin.HandlerFunc {

	return func(c *gin.Context) {

		table_id := c.Param("table_id")

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var mytable models.Table

		err := tablecollection.FindOne(ctx, bson.M{"table_id": table_id}).Decode(&mytable)
		defer cancel()

		if err != nil {

			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to get the data"})
			return
		}

		c.JSON(http.StatusOK, mytable)
	}
}

func CreateTable() gin.HandlerFunc {

	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var table models.Table

		if err := c.ShouldBindJSON(&table); err != nil {

			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		checkerr := validate.Struct(table)

		if checkerr != nil {

			c.JSON(http.StatusBadRequest, gin.H{"error": "While validating error"})
			return
		}

		table.ID = primitive.NewObjectID()
		table.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		table.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		table.Table_id = table.ID.Hex()

		newdata, err := tablecollection.InsertOne(ctx, table)

		if err != nil {

			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to insert the data"})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, newdata)
	}
}

func UpdateTable() gin.HandlerFunc {

	return func(c *gin.Context) {

		tableId := c.Param("table_id")

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var table models.Table

		if err := c.ShouldBindJSON(&table); err != nil {

			c.JSON(http.StatusBadRequest, gin.H{"error": "getting error in the request"})
			return
		}

		var updateobj primitive.D

		if table.Number_of_guests != nil {

			updateobj = append(updateobj, bson.E{"number_of_guests", table.Number_of_guests})
		}

		if table.Table_number != nil {
			updateobj = append(updateobj, bson.E{"table_number", table.Table_number})
		}

		table.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		filter := bson.M{"table_id": tableId}

		upsert := true
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		updated, err := tablecollection.UpdateOne(

			ctx,
			filter,
			bson.D{

				{"$set", updateobj},
			},
			&opt,
		)
		if err != nil {
			msg := fmt.Sprintf("Unable to update")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, updated)
	}
}
