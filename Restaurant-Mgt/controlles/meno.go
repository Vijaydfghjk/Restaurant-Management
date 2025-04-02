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

var menucollection *mongo.Collection = database.OpenCollection(database.Client, "menu")

func GetMenus() gin.HandlerFunc {

	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		result, err := menucollection.Find(context.TODO(), bson.M{})
		defer cancel()

		if err != nil {

			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing the menu items"})
		}

		var allmenus []bson.M

		for result.Next(ctx) {

			var menu bson.M

			if err := result.Decode(&menu); err != nil {

				log.Fatal(err)
			}

			allmenus = append(allmenus, menu)
		}
		c.JSON(http.StatusOK, allmenus)
	}
}

func Getmenu() gin.HandlerFunc {

	return func(c *gin.Context) {

		menuId := c.Param("menu_id")

		var menu models.Menu

		var ctx, calcel = context.WithTimeout(context.Background(), 100*time.Second)

		err := menucollection.FindOne(ctx, bson.M{"menu_id": menuId}).Decode(&menu)
		defer calcel()

		if err != nil {

			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while fetching the menu"})
		}
		c.JSON(http.StatusOK, menu)
	}
}

func CreateMenu() gin.HandlerFunc {

	return func(c *gin.Context) {

		var menu models.Menu
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		if err := c.ShouldBindJSON(&menu); err != nil {

			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(menu)

		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return

		}

		menu.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		menu.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		menu.ID = primitive.NewObjectID()
		menu.Menu_id = menu.ID.Hex()

		result, err := menucollection.InsertOne(ctx, menu)

		if err != nil {
			msg := fmt.Sprintf("Menu item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		}
		defer cancel()
		c.JSON(http.StatusOK, result)
		defer cancel()
	}
}

func inTimeSpan(start, end, check time.Time) bool {

	return start.After(time.Now()) && end.After(start)
}

func UpdateMenu() gin.HandlerFunc {

	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var menu models.Menu

		if err := c.ShouldBindJSON(&menu); err != nil {

			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		menuId := c.Param("menu_id")

		filter := bson.M{"menu_id": menuId}

		vali_err := validate.Struct(menu)

		if vali_err != nil {

			c.JSON(http.StatusBadRequest, gin.H{"error": vali_err})
			return
		}

		var updateobj primitive.D

		updateobj = append(updateobj, bson.E{"start_date", menu.Start_Date})
		updateobj = append(updateobj, bson.E{"end_date", menu.End_Date})

		if menu.Name != "" {

			updateobj = append(updateobj, bson.E{"name", menu.Name})
		}
		if menu.Category != "" {
			updateobj = append(updateobj, bson.E{"category", menu.Category})
		}

		menu.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateobj = append(updateobj, bson.E{"updated_at", menu.Updated_at})

		log.Println("chech", updateobj)
		upsert := true
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		result, err := menucollection.UpdateOne(

			ctx,
			filter,
			bson.D{

				{"$set", updateobj},
			},
			&opt,
		)

		if err != nil {
			msg := fmt.Sprintf("Menu update failed")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		}

		c.JSON(http.StatusOK, result)

	}
}
