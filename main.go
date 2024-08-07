// Recipes API
//
// This is a sample recipes API. You can find out more about the API at https://github.com/PacktPublishing/Building-Distributed-Applications-in-Gin.
//
//	Schemes: http
//	Host: localhost:8080
//	BasePath: /
//	Version: 1.0.0
//	Contact: Mohamed Labouardy <mohamed@labouardy.com> https://labouardy.com
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
//	swagger:meta
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

var recipes []Recipe

// swagger:parameters recipes newRecipe
type Recipe struct {
	//swagger:ignore
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Tags         []string  `json:"tags"`
	Ingredients  []string  `json:"ingredients"`
	Instructions []string  `json:"instructions"`
	PublishedAt  time.Time `json:"publishedAt"`
}

// swagger:operation POST /recipes recipes newRecipe
// Create a new recipe
// ---
// produces:
// - application/json
// responses:
//
//		'200':
//	    	description: Successful operation
//		'400':
//	    	description: Invalid input
func NewRecipeHandler(c *gin.Context) {
	var recipe Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	recipe.ID = xid.New().String()
	recipe.PublishedAt = time.Now()
	recipes = append(recipes, recipe)
	c.JSON(http.StatusOK, recipe)
}

// swagger:operation GET /recipes recipes listRecipes
// Returns list of recipes
// ---
// produces:
// - application/json
// responses:
//
//		'200':
//	    	description: Successful operation
func ListRecipesHandler(c *gin.Context) {
	c.JSON(http.StatusOK, recipes)
}

// swagger:operation PUT /recipes/{id} recipes updateRecipe
// Update an existing recipe
// ---
// parameters:
//   - name: id
//     in: path
//     description: ID of the recipe
//     required: true
//     type: string
//
// produces:
// - application/json
// responses:
//
//		'200':
//	    	description: Successful operation
//		'400':
//	    	description: Invalid input
//		'404':
//	    	description: Invalid recipe ID
func UpdateRecipeHandler(c *gin.Context) {
	var recipe Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := c.Param("id")
	for i, r := range recipes {
		if r.ID == id {
			recipes[i] = recipe
			c.JSON(http.StatusOK, recipe)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "recipe not found"})
}

// swagger:operation DELETE /recipes/{id} recipes deleteRecipe
// Delete an existing recipe
// ---
// produces:
// - application/json
// parameters:
//   - name: id
//     in: path
//     description: ID of the recipe
//     required: true
//     type: string
//
// responses:
//
//		'200':
//	    	description: Successful operation
//		'404':
//	    	description: Invalid recipe ID
func DeleteRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	for i, r := range recipes {
		if r.ID == id {
			recipes = append(recipes[:i], recipes[i+1:]...)
			c.JSON(http.StatusOK, gin.H{"message": "recipe deleted"})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "recipe not found"})
}

// swagger:operation GET /recipes/search recipes findRecipe
// Search recipes based on tags
// ---
// produces:
// - application/json
// parameters:
//   - name: tag
//     in: query
//     description: recipe tag
//     required: true
//     type: string
//
// responses:
//
//		'200':
//	    	description: Successful operation
func SearchRecipesHandler(c *gin.Context) {
	tag := c.Query("tag")
	if tag == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing tag query parameter"})
		return
	}

	result := make([]Recipe, 0)
	for _, r := range recipes {
		for _, t := range r.Tags {
			if strings.EqualFold(t, tag) {
				result = append(result, r)
				break
			}
		}
	}

	c.JSON(http.StatusOK, result)
}

// swagger:operation GET /recipes/{id} recipes getRecipe
// Get a recipe by ID
// ---
// produces:
// - application/json
// parameters:
//   - name: id
//     in: path
//     description: ID of the recipe
//     required: true
//     type: string
//
// responses:
//
//		'200':
//	    	description: Successful operation
//		'404':
//	    	description: Invalid recipe ID
func GetRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	for _, r := range recipes {
		if r.ID == id {
			c.JSON(http.StatusOK, r)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "recipe not found"})
}

func init() {
	recipes = make([]Recipe, 0)
	//read file recipes.json
	b, err := os.ReadFile("recipes.json")
	if err != nil {
		log.Fatal(err)
	}
	//unmarshal json to recipes
	err = json.Unmarshal(b, &recipes)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	r := gin.Default()
	r.POST("/recipes", NewRecipeHandler)
	r.GET("/recipes", ListRecipesHandler)
	r.PUT("/recipes:id", UpdateRecipeHandler)
	r.DELETE("/recipes:id", DeleteRecipeHandler)
	r.GET("/recipes/search", SearchRecipesHandler)
	r.GET(("recipes/:id"), GetRecipeHandler)
	r.Run() // listen and serve on 8080
}
