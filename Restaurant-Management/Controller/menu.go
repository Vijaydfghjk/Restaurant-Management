package controller

import (
	"context"
	"fmt"
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

type Menu_db struct {
	menucollection *mongo.Collection
	validate       *validator.Validate
}

func Menu_controll() *Menu_db {

	return &Menu_db{
		menucollection: database.Open_collection(database.Client, "Menu"),
		validate:       validator.New(),
	}
}

func (a *Menu_db) CreateMenu(c *gin.Context) {

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	defer cancel()
	var menu model.Menu

	if err := c.ShouldBindJSON(&menu); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if valerr := a.validate.Struct(menu); valerr != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": valerr.Error()})
		return
	}

	menu.ID = primitive.NewObjectID()
	menu.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	menu.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	menu.Menu_id = menu.ID.Hex()

	result, err := a.menucollection.InsertOne(ctx, menu)

	if err != nil {
		msg := fmt.Sprintf("Menu item is not created")
		c.JSON(http.StatusInternalServerError, gin.H{"Error": msg})
		return
	}

	c.JSON(http.StatusOK, result)

}

func (a *Menu_db) Getmenus(c *gin.Context) {

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	defer cancel()

	result, err := a.menucollection.Find(ctx, bson.M{})

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var all_menus []bson.M

	for result.Next(ctx) {

		var menu bson.M

		if err := result.Decode(&menu); err != nil {

			log.Fatal(err)
		}

		all_menus = append(all_menus, menu)

	}
	c.JSON(http.StatusOK, all_menus)
}

func (a *Menu_db) GetmenubyId(c *gin.Context) {

	menu_id := c.Param("menu_id")

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	defer cancel()

	var menu model.Menu
	err := a.menucollection.FindOne(ctx, bson.M{"menu_id": menu_id}).Decode(&menu)

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, menu)
}

func (a *Menu_db) Updatemenu(c *gin.Context) {

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	defer cancel()

	menu_id := c.Param("menu_id")

	myfiletr := bson.M{"menu_id": menu_id}

	var menu model.Menu

	if err := c.ShouldBindJSON(&menu); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	vali_err := a.validate.Struct(menu)
	if vali_err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": vali_err})
		return
	}

	var updatemenu primitive.D

	if menu.Category != "" {

		updatemenu = append(updatemenu, bson.E{"category", menu.Category})

	}
	if menu.Name != "" {

		updatemenu = append(updatemenu, bson.E{"name", menu.Name})
	}

	menu.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	updatemenu = append(updatemenu, bson.E{"updated_at", menu.Updated_at})

	upsert := true
	opt := options.UpdateOptions{

		Upsert: &upsert,
	}

	updateresult, err := a.menucollection.UpdateOne(

		ctx,
		myfiletr,
		bson.D{

			{"$set", updatemenu},
		},
		&opt,
	)

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updateresult)
}
