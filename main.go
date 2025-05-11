package main

import (
	"log"
	"net/http"

	"github.com/fathurwithyou/Tubes2_BE_RecipePlayground/api"
)

func main() {
	log.Println("Server is running on http://localhost:5000")
	err := http.ListenAndServe(":5000", http.HandlerFunc(handler.Handler))
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
