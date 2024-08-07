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
	"context"
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

var ctx context.Context
var err error
var client *mongo.Client

// swagger:parameters recipes newRecipe
type Recipe struct {
	//swagger:ignore
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	Name         string             `json:"name" bson:"name"`
	Tags         []string           `json:"tags" bson:"tags"`
	Ingredients  []string           `json:"ingredients" bson:"ingredients"`
	Instructions []string           `json:"instructions" bson:"instructions"`
	PublishedAt  time.Time          `json:"publishedAt" bson:"publishedAt"`
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

	//recipe.ID = xid.New().String()
	recipe.ID = primitive.NewObjectID()
	recipe.PublishedAt = time.Now()
	//recipes = append(recipes, recipe)
	collection := client.Database(os.Getenv(
		"MONGO_DATABASE")).Collection("recipes")
	_, err = collection.InsertOne(ctx, recipe)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": "Error while inserting a new recipe"})
		return
	}
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
	collection := client.Database(os.Getenv(
		"MONGO_DATABASE")).Collection("recipes")
	// Passing bson.M{{}} as the filter matches all documents in the collection
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(ctx)
	recipes := make([]Recipe, 0)
	for cursor.Next(ctx) {
		var recipe Recipe
		cursor.Decode(&recipe)
		recipes = append(recipes, recipe)
	}
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
	objectId, _ := primitive.ObjectIDFromHex(id)
	collection := client.Database(os.Getenv(
		"MONGO_DATABASE")).Collection("recipes")
	_, err := collection.UpdateOne(ctx,
		bson.M{"_id": objectId},
		bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "name", Value: recipe.Name},
				{Key: "instructions", Value: recipe.Instructions},
				{Key: "ingredients", Value: recipe.Ingredients},
				{Key: "tags", Value: recipe.Tags},
			}},
		},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Recipe has been updated"})
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
	objectId, _ := primitive.ObjectIDFromHex(id)
	collection := client.Database(os.Getenv(
		"MONGO_DATABASE")).Collection("recipes")
	_, err := collection.DeleteOne(ctx, bson.M{
		"_id": objectId,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Recipe has been deleted"})
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
/*
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
*/

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
	objectId, _ := primitive.ObjectIDFromHex(id)
	collection := client.Database(os.Getenv(
		"MONGO_DATABASE")).Collection("recipes")

	cur := collection.FindOne(ctx, bson.M{
		"_id": objectId,
	})
	var recipe Recipe
	err := cur.Decode(&recipe)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, recipe)
}

func init() {
	/*
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
	*/
	ctx = context.Background()
	client, err = mongo.Connect(ctx,
		options.Client().ApplyURI(os.Getenv(("MONGO_URI"))))
	if err != nil {
		log.Fatal(err)
	}
	// Check the connection
	err = client.Ping(context.TODO(), readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB!")

	/*
		var documents []interface{}
		for _, recipe := range recipes {
			documents = append(documents, recipe)
		}
		collection := client.Database(os.Getenv(
			"MONGO_DATABASE")).Collection("recipes")
		insertManyResult, err := collection.InsertMany(ctx, documents)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Inserted recipes: ",
			len(insertManyResult.InsertedIDs))
	*/
}

func main() {
	r := gin.Default()
	r.POST("/recipes", NewRecipeHandler)
	r.GET("/recipes", ListRecipesHandler)
	r.PUT("/recipes:id", UpdateRecipeHandler)
	r.DELETE("/recipes:id", DeleteRecipeHandler)
	//r.GET("/recipes/search", SearchRecipesHandler)
	r.GET(("recipes/:id"), GetRecipeHandler)
	r.Run() // listen and serve on 8080
}
