// Recipes API
//
// This is a sample recipes API. You can find out more about the API at https://github.com/afonsir/gin-recipes-api
//
// Schemes: http
// Host: localhost:8080
// BasePath: /
// Version: 1.0.0
// Contact: Afonso Costa <afonso@mail.com>
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
// swagger:meta
package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	handlers "github.com/afonsir/gin-recipes-api/handlers"
)

var recipesHandler *handlers.RecipesHandler

func init() {
	ctx := context.Background()

	client, err := mongo.Connect(ctx,
		options.Client().ApplyURI(os.Getenv("MONGODB_URI")))

	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDB")

	collection := client.Database(os.Getenv("MONGODB_DATABASE")).Collection("recipes")

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	recipesHandler = handlers.NewRecipesHandler(ctx, collection, redisClient)
}

func main() {
	router := gin.Default()

	router.POST("/recipes", recipesHandler.NewRecipeHandler)
	router.GET("/recipes", recipesHandler.ListRecipesHandler)
	router.GET("/recipes/search", recipesHandler.SearchRecipesHandler)
	router.PUT("/recipes/:id", recipesHandler.UpdateRecipesHandler)
	router.DELETE("/recipes/:id", recipesHandler.DeleteRecipesHandler)

	router.Run()
}
