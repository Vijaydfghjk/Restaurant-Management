package controller

import (
	"context"
	"fmt"
	"log"
	"net/http"
	database "restaurant-management/Database"
	model "restaurant-management/Model"
	tokenservice "restaurant-management/Token_service"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var usercollection *mongo.Collection = database.Open_collection(database.Client, "User")
var validate *validator.Validate = validator.New()

func Signup(c *gin.Context) {

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	defer cancel()

	var user model.User

	if err := c.ShouldBindJSON(&user); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if valerr := validate.Struct(user); valerr != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": valerr.Error()})
		return
	}

	count, err := usercollection.CountDocuments(ctx, bson.M{"email": user.Email})

	if err != nil {

		log.Panic(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the email"})
		return
	}

	count, err = usercollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
	password := Hashpassword(*user.Password)

	if count > 0 {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "this email or phone number already exsits"})
		return
	}

	user.Password = &password
	user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.ID = primitive.NewObjectID()
	user.User_id = user.ID.Hex()

	inseruser, err := usercollection.InsertOne(ctx, user)

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, inseruser)

}
func Login(c *gin.Context) {

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	defer cancel()

	var user model.User

	var found_user model.User

	if err := c.ShouldBindJSON(&user); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := usercollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&found_user)

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"Messsage": err.Error()})
		return
	}

	ok, msg := veifupassword(*user.Password, *found_user.Password)

	if !ok {

		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}

	signedtoken, refershtoken, _ := tokenservice.GenerateToken(*found_user.Email, *found_user.First_name, *found_user.Last_name, found_user.User_id)

	found_user.Token = &signedtoken
	found_user.Refresh_Token = &refershtoken

	tokenservice.UpdateToken(signedtoken, refershtoken, found_user.User_id)

	c.JSON(http.StatusOK, found_user)
}
func Hashpassword(s string) string {

	bytes, err := bcrypt.GenerateFromPassword([]byte(s), 14)

	if err != nil {

		log.Panic(err)
	}

	return string(bytes)
}

func veifupassword(user_password, provided_password string) (bool, string) {

	err := bcrypt.CompareHashAndPassword([]byte(provided_password), []byte(user_password))

	check := true

	msg := ""

	if err != nil {

		msg = fmt.Sprintf("Email or password is incorrect")
		check = false
	}

	return check, msg

}
