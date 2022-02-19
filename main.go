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
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// swagger:parameters recipes newRecipe
type Recipe struct {
	// swagger:ignore
	ID           primitive.ObjectID `json:"id"           bson:"_id"`
	Name         string             `json:"name"         bson:"name"`
	Tags         []string           `json:"tags"         bson:"tags"`
	Ingredients  []string           `json:"ingredients"  bson:"ingredients"`
	Instructions []string           `json:"instructions" bson:"instructions"`
	PublishedAt  time.Time          `json:"publishedAt"  bson:"publishedAt"`
}

var recipes []Recipe
var ctx context.Context
var err error
var client *mongo.Client
var collection *mongo.Collection

func init() {
	ctx = context.Background()

	client, err = mongo.Connect(ctx,
		options.Client().ApplyURI(os.Getenv("MONGODB_URI")))

	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDB")

	collection = client.Database(os.Getenv("MONGODB_DATABASE")).Collection("recipes")
}

func main() {
	router := gin.Default()

	router.POST("/recipes", NewRecipeHandler)
	router.GET("/recipes", ListRecipesHandler)
	router.GET("/recipes/search", SearchRecipesHandler)
	router.PUT("/recipes/:id", UpdateRecipesHandler)
	router.DELETE("/recipes/:id", DeleteRecipesHandler)

	router.Run()
}

// swagger:operation POST /recipes recipes createRecipe
//
// Create a new recipe
//
// ---
// produces:
// - application/json
//
// responses:
//   '200':
//     description: Successful operation
//   '400':
//     description: Invalid input
func NewRecipeHandler(c *gin.Context) {
	var recipe Recipe

	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	recipe.ID = xid.New().String()
	recipe.PublishedAt = time.Now()
	recipes = append(recipes, recipe)

	c.JSON(http.StatusOK, recipe)
}

// swagger:operation GET /recipes recipes listRecipes
//
// Returns list of recipes
//
// ---
// produces:
// - application/json
//
// responses:
//   '200':
//     description: Successful operation
func ListRecipesHandler(c *gin.Context) {
	c.JSON(http.StatusOK, recipes)
}

// swagger:operation GET /recipes/search recipes searchRecipes
//
// Search recipes by tag
//
// ---
// parameters:
// - name: tag
//   in: query
//   description: tag to filter by
//   required: false
//   type: string
//
// produces:
// - application/json
//
// responses:
//   '200':
//     description: Successful operation
func SearchRecipesHandler(c *gin.Context) {
	tag := c.Query("tag")
	listOfRecipes := make([]Recipe, 0)

	for i := 0; i < len(recipes); i++ {
		found := false

		for _, t := range recipes[i].Tags {
			if strings.EqualFold(t, tag) {
				found = true
			}
		}

		if found {
			listOfRecipes = append(listOfRecipes, recipes[i])
		}
	}

	c.JSON(http.StatusOK, listOfRecipes)
}

// swagger:operation PUT /recipes/{id} recipes updateRecipe
//
// Update an existing recipe
//
// ---
// parameters:
// - name: id
//   in: path
//   description: ID of the recipe
//   required: true
//   type: string
//
// produces:
// - application/json
//
// responses:
//   '200':
//     description: Successful operation
//   '400':
//     description: Invalid input
//   '404':
//     description: Invalid recipe ID
func UpdateRecipesHandler(c *gin.Context) {
	id := c.Param("id")

	var recipe Recipe

	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	index := -1

	for i := 0; i < len(recipes); i++ {
		if recipes[i].ID == id {
			index = i
		}
	}

	if index == -1 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Recipe not found",
		})
		return
	}

	recipes[index] = recipe

	c.JSON(http.StatusOK, recipe)
}

// swagger:operation DELETE /recipes/{id} recipes deleteRecipe
//
// Delete an existing recipe
//
// ---
// parameters:
// - name: id
//   in: path
//   description: ID of the recipe
//   required: true
//   type: string
//
// produces:
// - application/json
//
// responses:
//   '200':
//     description: Successful operation
//   '404':
//     description: Invalid recipe ID
func DeleteRecipesHandler(c *gin.Context) {
	id := c.Param("id")
	index := -1

	for i := 0; i < len(recipes); i++ {
		if recipes[i].ID == id {
			index = i
		}
	}

	if index == -1 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Recipe not found",
		})
		return
	}

	recipes = append(recipes[:index], recipes[index+1:]...)

	c.JSON(http.StatusOK, gin.H{
		"message": "Recipe has been deleted",
	})
}
