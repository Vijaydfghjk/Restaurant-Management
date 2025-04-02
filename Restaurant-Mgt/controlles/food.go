package controlles

import (
	"context"
	"fmt"
	"log"
	"math"
	"net/http"
	database "restaurant/Database"
	models "restaurant/Models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var foodcollection *mongo.Collection = database.OpenCollection(database.Client, "food")

var validate = validator.New()

func GetFoods() gin.HandlerFunc {

	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		recordperpage, err := strconv.Atoi(c.Query("recordPerPage"))

		if err != nil || recordperpage < 1 {

			recordperpage = 10
		}

		page, err := strconv.Atoi(c.Query("page"))

		if err != nil || page < 1 {

			page = 1
		}

		startIndex := (page - 1) * recordperpage
		log.Printf("start index %v type is %T\n", startIndex, startIndex)

		matchStage := bson.D{{"$match", bson.D{{}}}}
		groupstage := bson.D{{"$group", bson.D{{"_id", bson.D{{"_id", "null"}}}, {"total_count", bson.D{{"$sum", 1}}}, {"data", bson.D{{"$push", "$$ROOT"}}}}}}

		projectstage := bson.D{
			{
				"$project", bson.D{

					{"_id", 0},
					{"total_count", 1},
					{"food_items", bson.D{{"$slice", []interface{}{"$data", startIndex, recordperpage}}}},
				}}}

		result, err := foodcollection.Aggregate(ctx, mongo.Pipeline{matchStage, groupstage, projectstage})

		defer cancel()
		if err != nil {

			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing food items"})
		}

		var allfoods []bson.M

		if err = result.All(ctx, &allfoods); err != nil {

			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allfoods[0])
	}
}

func Getfood() gin.HandlerFunc {

	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		foodID := c.Param("foof_id")

		var food models.Food

		err := foodcollection.FindOne(ctx, bson.M{"food_id": foodID}).Decode(&food)
		defer cancel()
		if err != nil {

			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while fetching the food item"})
		}
		c.JSON(http.StatusOK, food)
	}
}

func Createfood() gin.HandlerFunc {

	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		//var menu models.Menu
		var food models.Food

		if err := c.ShouldBindJSON(&food); err != nil {

			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}
		validaterr := validate.Struct(food)

		if validaterr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validaterr.Error()})
			return
		}

		//err := menucollection.FindOne(ctx, bson.M{"menu_id": food.Menu_id}).Decode(&menu)
		defer cancel()

		food.Id = primitive.NewObjectID()
		food.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		food.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		food.Food_id = food.Id.Hex()

		var nun = toFixed(*food.Price, 2)

		food.Price = &nun

		result, err := foodcollection.InsertOne(ctx, food)
		defer cancel()
		if err != nil {
			msg := fmt.Sprintf("Food item is not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		}

		c.JSON(http.StatusOK, result)
	}
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {

	output := math.Pow(10, float64(precision))

	return float64(round(num*output)) / output
}

func Updatefood() gin.HandlerFunc {

	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var food models.Food
		if err := c.ShouldBindJSON(&food); err != nil {

			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		foordId := c.Param("food_id")

		vali_err := validate.Struct(food)

		if vali_err != nil {

			c.JSON(http.StatusBadRequest, gin.H{"error": vali_err})
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
		food.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateobj = append(updateobj, bson.E{"updated_at", food.Updated_at})

		upser := true

		opt := options.UpdateOptions{

			Upsert: &upser,
		}

		filter := bson.M{"food_id": foordId}

		result, err := foodcollection.UpdateOne(
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
}
