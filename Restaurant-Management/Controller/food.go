package controller

import (
	"context"
	"fmt"
	"log"
	"math"
	"net/http"
	database "restaurant-management/Database"
	model "restaurant-management/Model"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Food_db struct {
	foodcollection *mongo.Collection
	validate       *validator.Validate
}

func Food_controll() *Food_db {
	return &Food_db{
		foodcollection: database.Open_collection(database.Client, "Food"),
		validate:       validator.New(),
	}
}

func (a *Food_db) GetFoods(c *gin.Context) {

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	defer cancel()

	recordperpage, err := strconv.Atoi(c.Query("recordPerPage"))

	if err != nil && recordperpage < 1 {

		recordperpage = 10
	}

	page, err := strconv.Atoi(c.Query("page"))

	if err != nil && page < 1 {

		recordperpage = 1
	}

	startIndex := (page - 1) * recordperpage

	mathstage := bson.D{{"$match", bson.D{}}}

	groupstage := bson.D{{"$group", bson.D{{"_id", bson.D{{"_id", "null"}}}, {"total_count", bson.D{{"$sum", 1}}}, {"data", bson.D{{"$push", "$$ROOT"}}}}}}

	projectstage := bson.D{{

		"$project", bson.D{

			{"_id", 0},
			{"total_count", 1},
			{"food_items", bson.D{{"$slice", []interface{}{"$data", startIndex, recordperpage}}}},
		},
	}}

	result, err := a.foodcollection.Aggregate(ctx, mongo.Pipeline{mathstage, groupstage, projectstage})

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing food items"})
		return
	}

	var allfoods []bson.M

	if err = result.All(ctx, &allfoods); err != nil {

		log.Fatal(err)
	}

	c.JSON(http.StatusOK, allfoods[0])
}

func (a *Food_db) GetFood(c *gin.Context) {

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	foodID := c.Param("foof_id")

	var food model.Food

	err := a.foodcollection.FindOne(ctx, bson.M{"food_id": foodID}).Decode(&food)

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while fetching the food item"})
		return
	}
	c.JSON(http.StatusOK, food)
}

func (a *Food_db) Createfood(c *gin.Context) {

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	defer cancel()

	var food model.Food

	if err := c.ShouldBindJSON(&food); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	valerr := a.validate.Struct(food)

	if valerr != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": valerr.Error()})
		return
	}

	food.Id = primitive.NewObjectID()
	food.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	food.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	food.Food_id = food.Id.Hex()

	var nun = toFixed(*food.Price, 2)

	food.Price = &nun

	result, err := a.foodcollection.InsertOne(ctx, food)

	if err != nil {
		msg := fmt.Sprintf("Food item is not created")
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusOK, result)
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {

	output := math.Pow(10, float64(precision))

	return float64(round(num*output)) / output
}

func (a *Food_db) Updatefood(c *gin.Context) {

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	defer cancel()
	foordId := c.Param("food_id")

	var food model.Food

	if err := c.ShouldBindJSON(&food); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	vali_err := a.validate.Struct(food)

	if vali_err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": vali_err.Error()})
		return
	}

	var updateobj primitive.D

	if food.Name != nil {

		updateobj = append(updateobj, bson.E{"name", food.Name})
	}

	if food.Price != nil {

		updateobj = append(updateobj, bson.E{"price", food.Price})
	}

	if food.Food_image != nil {

		updateobj = append(updateobj, bson.E{"food_image", food.Food_image})
	}

	if food.Menu_id != nil {

		updateobj = append(updateobj, bson.E{"menu_id", food.Menu_id})
	}

	food.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateobj = append(updateobj, bson.E{"updated_at", food.Updated_at})

	upser := true

	opt := options.UpdateOptions{

		Upsert: &upser,
	}

	filter := bson.M{"food_id": foordId}

	result, err := a.foodcollection.UpdateOne(
		ctx,
		filter,
		bson.D{

			{"$set", updateobj},
		},
		&opt,
	)

	if err != nil {
		msg := fmt.Sprint("Unable to update")
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusOK, result)

}
