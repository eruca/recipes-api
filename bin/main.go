package main

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"log"
	"os"

	"github.com/eruca/recipes-api/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	users := map[string]string{
		"admin":      "admin",
		"packt":      "packt",
		"mlabouardy": "mlabouardy",
	}

	ctx := context.Background()
	client, err := mongo.Connect(ctx,
		options.Client().ApplyURI(os.Getenv(utils.ENV_MONGO_URI)))
	if err != nil {
		panic(err)
	}
	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB")

	col := client.Database(os.Getenv(utils.ENV_MONGO_DATABASE)).Collection("users")

	h := sha256.New()
	for k, v := range users {
		pswd := base64.StdEncoding.EncodeToString(h.Sum([]byte(v)))
		if _, err := col.InsertOne(ctx, bson.M{
			"username": k,
			"password": string(pswd),
		}); err != nil {
			panic(err)
		}
	}
}
