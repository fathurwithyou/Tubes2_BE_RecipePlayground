package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/fathurwithyou/Tubes2_BE_RecipePlayground/service/model"
	"github.com/fathurwithyou/Tubes2_BE_RecipePlayground/service/scraper"
	"github.com/fathurwithyou/Tubes2_BE_RecipePlayground/service/solver"

	"github.com/gin-gonic/gin"
)

var (
	app *gin.Engine
)

func main() {
	if err := scraper.Scrape("alchemy_elements.json"); err != nil {
		log.Fatalf("Failed to scrape elements: %v", err)
	}
	appData, err := model.LoadElementsFromFile("alchemy_elements.json")
	if err != nil {
		log.Fatalf("Failed to load element data: %v", err)
	}
	solver.InitElementsMap(appData)

	app = gin.New()

	app.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	app.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	app.GET("/solve/:method/:target/:maxRecipe", func(c *gin.Context) {
		method := c.Param("method")
		target := c.Param("target")
		maxRecipeStr := c.Param("maxRecipe")
		maxRecipeInt, err := strconv.Atoi(maxRecipeStr)
		if err != nil || maxRecipeInt < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid maxRecipe"})
			return
		}
		maxRecipes := int64(maxRecipeInt)
		var resultData interface{}

		switch method {
		case "bfs":
			resultData = solver.Bfs(target, maxRecipes)
		case "dfs":
			resultData = solver.Dfs(target, maxRecipes)
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid method"})
			return
		}

		if resultData == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Element not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"result": resultData})
	})

	app.GET("/elements", func(c *gin.Context) {
		elements := solver.GetAllElements()
		c.JSON(http.StatusOK, gin.H{"elements": elements})
	})

	log.Println("Server running on port :8080")
	if err := app.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
