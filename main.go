package main

import (
	"context"
	"log"
	"os"

	"github.com/eruca/recipes-api/handlers"
	"github.com/eruca/recipes-api/utils"
	"github.com/gin-gonic/gin"
	redis "github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	authHandler    *handlers.AuthHandler
	recipesHandler *handlers.RecipeHandler
)

func init() {
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
	collection := client.Database(os.Getenv(utils.ENV_MONGO_DATABASE)).Collection(utils.COLLECTION)

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	status := redisClient.Ping(ctx)
	log.Println("Redis Client Ping:", status.String())

	recipesHandler = handlers.NewRecipeHandler(ctx, collection, redisClient)

	collectionUser := client.Database(os.Getenv(utils.ENV_MONGO_DATABASE)).Collection("users")
	authHandler = handlers.NewAuthHandler(ctx, collectionUser)
}

func main() {

	router := gin.Default()
	router.Use(Cors())
	router.POST("/signin", authHandler.SignInHandler)
	router.POST("/refresh", authHandler.RefreshHandler)

	authrized := router.Group("/")
	authrized.Use(handlers.AuthMiddleware())
	{
		authrized.POST("/recipes", recipesHandler.NewRecipeHandler)
		authrized.GET("/recipes", recipesHandler.ListRecipesHandler)
		authrized.PUT("/recipes/:id", recipesHandler.UpdateRecipeHandler)
		authrized.DELETE("/recipes/:id", recipesHandler.DeleteRecipeHandler)
	}

	log.Fatal(router.Run(":8080"))
}

// CORSMiddleware 实现跨域
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
