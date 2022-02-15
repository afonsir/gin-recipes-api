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
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

type Recipe struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Tags         []string  `json:"tags"`
	Ingredients  []string  `json:"ingredients"`
	Instructions []string  `json:"instructions"`
	PublishedAt  time.Time `json:"publishedAt"`
}

var recipes []Recipe

func init() {
	recipes = make([]Recipe, 0)
	file, _ := ioutil.ReadFile("recipes.json")
	_ = json.Unmarshal([]byte(file), &recipes)
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
