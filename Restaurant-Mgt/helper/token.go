package helper

import (
	"context"
	"fmt"
	"log"
	"os"
	database "restaurant/Database"
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

var SECRET_KEY string = os.Getenv("SECRET_KEY")

var usercollection *mongo.Collection = database.OpenCollection(database.Client, "user")

func GenerateallToken(email string, firstName string, lastName string, uid string) (signedToken string, signedRefreshToken string, err error) {

	claims := &SignedDetails{

		Email:      email,
		First_name: firstName,
		Last_name:  lastName,
		Uid:        uid,
		StandardClaims: jwt.StandardClaims{

			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(), // ExpiresAt value to 24 hours from the current local time.
		},
	}

	refereshclaims := &SignedDetails{

		StandardClaims: jwt.StandardClaims{

			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(), // ExpiresAt value to 7 days from the current local time.
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))

	refereshtoken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refereshclaims).SignedString([]byte(SECRET_KEY))

	if err != nil {

		log.Panic(err)
	}
	return token, refereshtoken, nil
}

func UpdateallToken(signedToken string, signedRefreshToken string, userId string) {

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	var updateobj primitive.D

	updateobj = append(updateobj, bson.E{"token", signedToken})
	updateobj = append(updateobj, bson.E{"referesh_token", signedRefreshToken})

	updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	updateobj = append(updateobj, bson.E{"updated_at", updated_at})

	upsert := true

	filter := bson.M{"user_id": userId}

	opt := options.UpdateOptions{

		Upsert: &upsert,
	}

	_, err := usercollection.UpdateOne(

		ctx,
		filter,
		bson.D{

			{"$set", updateobj},
		},
		&opt,
	)

	defer cancel()

	if err != nil {

		log.Panic(err)
		return
	}
	return

}
func ValidateToken(signedToken string) (claims *SignedDetails, msg string) {

	token, err := jwt.ParseWithClaims( // SECRET_KEY is the same key that was used when the JWT was initially signed.
		signedToken,
		&SignedDetails{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
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
