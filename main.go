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

type Recipe struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Tags         []string  `json:"tags"`
	Ingredients  []string  `json:"ingredients"`
	Instructions []string  `json:"instructions"`
	PublishedAt  time.Time `json:"publishedAt"`
}

// NewRecipeHandler creates a new recipe
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

// ListRecipesHandler lists all recipes
func ListRecipesHandler(c *gin.Context) {
	c.JSON(http.StatusOK, recipes)
}

// UpdateRecipeHandler updates a recipe
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

// DeleteRecipeHandler deletes a recipe
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

// SearchRecipesHandler searches recipes by tag
func SearchRecipesHandler(c *gin.Context) {
	tag := c.Query("tag")
	if tag == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing tag query parameter"})
		return
	}

	result := make([]Recipe, 0)
	for _, r := range recipes {
		for _, t := range r.Tags {
			if strings.EqualFold(t, tag)  {
				result = append(result, r)
				break
			}
		}
	}

	c.JSON(http.StatusOK, result)
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
	r.Run() // listen and serve on 8080
}
