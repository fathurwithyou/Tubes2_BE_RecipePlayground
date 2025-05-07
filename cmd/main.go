package main

import (
	"log"
	"net/http"
	"strconv"

	"Tubes2_BE_RecipePlayground/internal/model"
	"Tubes2_BE_RecipePlayground/internal/scraper"
	"Tubes2_BE_RecipePlayground/internal/solver"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := scraper.Scrape("alchemy_elements.json"); err != nil {
		log.Fatalf("Gagal scrape elemen: %v", err)
	}
	appData, err := model.LoadElementsFromFile("alchemy_elements.json")
	if err != nil {
		log.Fatalf("Gagal memuat data elemen: %v", err)
	}
	recipetree.InitElementsMap(appData)

	router := gin.Default()

	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	router.GET("/solve/:target/:maxDepth", func(c *gin.Context) {
		target := c.Param("target")
		maxDepth, _ := strconv.Atoi(c.Param("maxDepth"))
		var resultData interface{}

		resultData = recipetree.GenerateRecipeTree(target, maxDepth)
		c.JSON(http.StatusOK, gin.H{"result": resultData})
	})

	log.Println("Server berjalan di port :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}
