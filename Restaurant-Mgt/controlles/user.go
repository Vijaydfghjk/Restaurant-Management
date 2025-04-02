package controlles

import (
	"context"
	"fmt"
	"log"
	"net/http"
	database "restaurant/Database"
	models "restaurant/Models"
	"restaurant/helper"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var usercollection *mongo.Collection = database.OpenCollection(database.Client, "user")

func GetUsers() gin.HandlerFunc {

	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))

		if err != nil || recordPerPage < 1 {

			recordPerPage = 10
		}

		page, err := strconv.Atoi(c.Query("page"))

		if err != nil || page < 1 {

			page = 1
		}

		startIndex := (page - 1) * recordPerPage
		startIndex, err = strconv.Atoi(c.Query("startIndex"))

		matchstage := bson.D{{"$match", bson.D{{}}}}

		projectstage := bson.D{{"$project", bson.D{

			{"_id", 0},
			{"total_count", 1},
			{"user_items", bson.D{{"$slice", []interface{}{"$data", startIndex, recordPerPage}}}},
		},
		}}

		result, err := usercollection.Aggregate(ctx, mongo.Pipeline{

			matchstage, projectstage})
		defer cancel()

		if err != nil {

			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing user items"})
		}

		var allusers []bson.M

		if err = result.All(ctx, &allusers); err != nil {

			log.Fatal(err)
		}

		c.JSON(http.StatusOK, allusers)
	}
}

func Getuser() gin.HandlerFunc {

	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		userid := c.Param("user_id")

		var user models.User

		err := usercollection.FindOne(ctx, bson.M{"user_id": userid}).Decode(&user)

		defer cancel()

		if err != nil {

			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing user items"})
		}

		c.JSON(http.StatusOK, user)
	}
}

func Signup() gin.HandlerFunc {

	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var user models.User

		if err := c.ShouldBindJSON(&user); err != nil {

			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		validationerr := validate.Struct(user)

		if validationerr != nil {

			c.JSON(http.StatusBadRequest, gin.H{"error": validationerr.Error()})
			return
		}

		count, err := usercollection.CountDocuments(ctx, bson.M{"email": user.Email})

		defer cancel()

		if err != nil {

			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the email"})
			return
		}

		password := Hashpassword(*user.Password)

		user.Password = &password

		count, err = usercollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		defer cancel()
		if err != nil {

			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the phone_number"})
			return
		}

		if count > 0 {

			c.JSON(http.StatusInternalServerError, gin.H{"error": "this email or phone number already exsits"})
			return
		}

		user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()

		token, refersh_token, _ := helper.GenerateallToken(*user.Email, *user.First_name, *user.Last_name, user.User_id)

		user.Token = &token
		user.Refresh_Token = &refersh_token

		resultInsertionNumber, insertErr := usercollection.InsertOne(ctx, user)

		if insertErr != nil {

			msg := fmt.Sprintf("User item was not created")

			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		defer cancel()

		c.JSON(http.StatusOK, resultInsertionNumber)
	}
}

func Login() gin.HandlerFunc {

	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User

		var founded_user models.User

		if err := c.ShouldBindJSON(&user); err != nil {

			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		err := usercollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&founded_user)
		defer cancel()
		if err != nil {

			c.JSON(http.StatusInternalServerError, gin.H{"Messsage": err.Error()})
			return
		}

		pwdcheck, msg := veifupassword(*user.Password, *founded_user.Password)

		if !pwdcheck {

			c.JSON(http.StatusInternalServerError, gin.H{"Messsage": msg})
			return
		}

		token, refersh_token, _ := helper.GenerateallToken(*founded_user.Email, *founded_user.First_name, *founded_user.Last_name, founded_user.User_id)

		helper.UpdateallToken(token, refersh_token, founded_user.User_id)

		c.JSON(http.StatusOK, founded_user)
	}
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
