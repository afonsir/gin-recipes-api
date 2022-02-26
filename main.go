// Recipes API
//
// This is a sample recipes API. You can find out more about the API at https://github.com/afonsir/gin-recipes-api
//
// Schemes: http
// Host: api.recipes.io:8080
// BasePath: /
// Version: 1.0.0
// Contact: Afonso Costa <afonso@mail.com>
//
// SecurityDefinitions:
//   api_key:
//     type: apiKey
//     name: Authorization
//     in: header
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
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	redisStore "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	handlers "github.com/afonsir/gin-recipes-api/handlers"
)

var authHandler *handlers.AuthHandler
var recipesHandler *handlers.RecipesHandler
var originURL string

func init() {
	ctx := context.Background()

	client, err := mongo.Connect(ctx,
		options.Client().ApplyURI(os.Getenv("MONGODB_URI")))

	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDB")

	recipesCol := client.Database(os.Getenv("MONGODB_DATABASE")).Collection("recipes")
	usersCol := client.Database(os.Getenv("MONGODB_DATABASE")).Collection("users")

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	authMechanism := os.Getenv("AUTH_MECHANISM")

	authHandler = handlers.NewAuthHandler(ctx, usersCol, authMechanism)
	recipesHandler = handlers.NewRecipesHandler(ctx, recipesCol, redisClient)

	originURL = os.Getenv("ORIGIN_URL")
}

func main() {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{originURL},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	if authHandler.AuthMechanism == "COOKIE" {
		store, _ := redisStore.NewStore(10, "tcp", "localhost:6379", "", []byte("secret"))
		router.Use(sessions.Sessions("recipe_api", store))
	}

	router.GET("/recipes", recipesHandler.ListRecipesHandler)
	router.GET("/recipes/:id", recipesHandler.GetOneRecipeHandler)
	router.GET("/recipes/search", recipesHandler.SearchRecipesHandler)

	if authHandler.AuthMechanism == "COOKIE" {
		router.POST("/signin", authHandler.SignInWithCookieHandler)
		router.POST("/signout", authHandler.SignOutHandler)
	} else {
		router.POST("/signin", authHandler.SignInHandler)
		router.POST("/refresh", authHandler.RefreshHandler)
	}

	authorized := router.Group("/")
	authorized.Use(authHandler.AuthMiddleware())

	authorized.POST("/recipes", recipesHandler.NewRecipeHandler)
	authorized.PUT("/recipes/:id", recipesHandler.UpdateRecipesHandler)
	authorized.DELETE("/recipes/:id", recipesHandler.DeleteRecipesHandler)

	router.RunTLS(":4430", "certs/localhost.crt", "certs/localhost.key")
}
