package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/fathurwithyou/Tubes2_BE_RecipePlayground/service/model"
	"github.com/fathurwithyou/Tubes2_BE_RecipePlayground/service/scraper"
	"github.com/fathurwithyou/Tubes2_BE_RecipePlayground/service/solver"

	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func init() {

	if err := scraper.Scrape("alchemy_elements.json"); err != nil {
		log.Fatalf("Failed to scrape elements: %v", err)
	}
	appData, err := model.LoadElementsFromFile("alchemy_elements.json")
	if err != nil {
		log.Fatalf("Failed to load element data: %v", err)
	}
	solver.InitElementsMap(appData)

	router = gin.New()

	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers",
			"Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	
	router.GET("/solve/:method/:target/:maxRecipe", func(c *gin.Context) {
		method := c.Param("method")
		target := c.Param("target")
		maxRecipeStr := c.Param("maxRecipe")
		maxRecipeInt, err := strconv.Atoi(maxRecipeStr)
		if err != nil || maxRecipeInt < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid maxRecipe"})
			return
		}

		var resultData interface{}
		switch method {
		case "bfs":
			resultData = solver.Bfs(target, int64(maxRecipeInt))
		case "dfs":
			resultData = solver.Dfs(target, int64(maxRecipeInt))
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid method"})
			return
		}
		if resultData == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Element not found"})
			return
		}

		visitedNodeCount := solver.GetVisitedNodeCount()

		c.JSON(http.StatusOK, gin.H{
			"result":              resultData,
			"visited_node_count":  visitedNodeCount,
		})
	})

	router.GET("/elements", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"elements": solver.GetAllElements()})
	})
}

func Handler(w http.ResponseWriter, r *http.Request) {
	router.ServeHTTP(w, r)
}
