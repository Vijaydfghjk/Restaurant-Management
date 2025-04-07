package tokenservice

import (
	"context"
	"fmt"
	"log"
	"os"
	database "restaurant-management/Database"
	"time"

	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SignedDetails struct {
	Email      string
	First_name string
	Last_name  string
	Uid        string
	jwt.StandardClaims
}

var Secret string = os.Getenv("SECRET_KEY")
var usercollection *mongo.Collection = database.Open_collection(database.Client, "User")

func GenerateToken(email string, firstName string, lastName string, uid string) (signedToken string, signedRefreshToken string, err error) {

	clims := &SignedDetails{

		Email:      email,
		First_name: firstName,
		Last_name:  lastName,
		Uid:        uid,
		StandardClaims: jwt.StandardClaims{

			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	refereshclaims := &SignedDetails{

		StandardClaims: jwt.StandardClaims{

			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, clims).SignedString([]byte(Secret))

	refereshtoken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refereshclaims).SignedString([]byte(Secret))

	if err != nil {

		log.Panic(err)
	}

	return token, refereshtoken, nil
}

func UpdateToken(signedToken string, signedRefreshToken string, userId string) {

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	defer cancel()
	var updateobject primitive.D

	updateobject = append(updateobject, bson.E{"token", signedToken})
	updateobject = append(updateobject, bson.E{"referesh_token", signedRefreshToken})

	updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	updateobject = append(updateobject, bson.E{"updated_at", updated_at})

	filter := bson.M{"user_id": userId}

	upsert := true

	opt := options.UpdateOptions{

		Upsert: &upsert,
	}
	_, err := usercollection.UpdateOne(

		ctx,
		filter,
		bson.D{

			{"$set", updateobject},
		},
		&opt,
	)

	if err != nil {

		log.Panic(err)
		return
	}

}

func Validatetoken(signedtoken string) (claims *SignedDetails, msg string) {

	token, err := jwt.ParseWithClaims(

		signedtoken,
		&SignedDetails{},
		func(t *jwt.Token) (interface{}, error) {

			return []byte(Secret), nil
		},
	)

	claims, ok := token.Claims.(*SignedDetails)

	if !ok {

		msg = fmt.Sprintf("The token is invalid")
		msg = err.Error()
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = fmt.Sprint("token is expired")
		msg = err.Error()
		return
	}

	return claims, msg
}
