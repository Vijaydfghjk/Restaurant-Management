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

type Table_db struct {
	tablecollection *mongo.Collection
	valiedate       *validator.Validate
}

func Tablecontroll() *Table_db {

	return &Table_db{
		tablecollection: database.Open_collection(database.Client, "Table"),
		valiedate:       validator.New(),
	}
}

func (a *Table_db) Create_table(c *gin.Context) {

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var table model.Table

	if err := c.ShouldBindJSON(&table); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	valerr := a.valiedate.Struct(table)

	if valerr != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": valerr.Error()})
		return
	}

	table.ID = primitive.NewObjectID()
	table.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	table.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	table.Table_id = table.ID.Hex()

	newdata, err := a.tablecollection.InsertOne(ctx, table)

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, newdata)

}

func (a *Table_db) Get_tables(c *gin.Context) {

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	result, err := a.tablecollection.Find(ctx, bson.M{})

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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

func (a *Table_db) Get_table(c *gin.Context) {

	table_id := c.Param("table_id")

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var mytable model.Table

	err := a.tablecollection.FindOne(ctx, bson.M{"table_id": table_id}).Decode(&mytable)

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, mytable)

}

func (a *Table_db) Update_table(c *gin.Context) {

	tableId := c.Param("table_id")

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var table model.Table

	if err := c.ShouldBindJSON(&table); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": "getting error in the request"})
		return
	}

	if vallerr := a.valiedate.Struct(table); vallerr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "getting error in the request"})
		return
	}

	var updateobj primitive.D
	if table.Number_of_guests != nil {

		updateobj = append(updateobj, bson.E{"number_of_guests", table.Number_of_guests})
	}

	if table.Table_number != "" {

		updateobj = append(updateobj, bson.E{"table_number", table.Table_number})

	}

	table.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	filter := bson.M{"table_id": tableId}

	upsert := true
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	updated, err := a.tablecollection.UpdateOne(

		ctx,
		filter,
		bson.D{

			{"$set", updateobj},
		},
		&opt,
	)
	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updated)
}
