package recipetree

import (
	"Tubes2_BE_RecipePlayground/internal/model"
	"fmt"
	"sync"
)

var (
	elementsMapGlobal map[string]model.Element
	cache             map[string]interface{}
	cacheMu           sync.RWMutex
)

func InitElementsMap(data model.Data) {
	elementsMapGlobal = make(map[string]model.Element)
	for _, el := range data.Elements {
		elementsMapGlobal[el.Name] = el
	}
}

func GenerateRecipeTree(rootElementName string, maxDepth int) interface{} {
	if elementsMapGlobal == nil {
		return map[string]string{"error": "Element data not initialized"}
	}
	// initialize cache per generation
	cacheMu.Lock()
	cache = make(map[string]interface{})
	cacheMu.Unlock()
	return generateNodeRecursive(rootElementName, 0, maxDepth)
}

func generateNodeRecursive(elementName string, currentDepth, maxDepth int) interface{} {
	key := fmt.Sprintf("%s|%d", elementName, currentDepth)
	// check cache
	cacheMu.RLock()
	if val, ok := cache[key]; ok {
		cacheMu.RUnlock()
		return val
	}
	cacheMu.RUnlock()

	eData, exists := elementsMapGlobal[elementName]
	if !exists || len(eData.Recipes) == 0 || currentDepth >= maxDepth {
		cacheMu.Lock()
		cache[key] = elementName
		cacheMu.Unlock()
		return elementName
	}

	var validRecipes [][]string
	for _, rec := range eData.Recipes {
		if len(rec) == 2 {
			validRecipes = append(validRecipes, rec)
		}
	}
	results := make([][]interface{}, len(validRecipes))
	var wg sync.WaitGroup
	wg.Add(len(validRecipes))
	for i, rec := range validRecipes {
		i, rec := i, rec
		go func() {
			defer wg.Done()
			n1 := generateNodeRecursive(rec[0], currentDepth+1, maxDepth)
			n2 := generateNodeRecursive(rec[1], currentDepth+1, maxDepth)
			results[i] = []interface{}{n1, n2}
		}()
	}
	wg.Wait()
	res := map[string]interface{}{elementName: results}
	cacheMu.Lock()
	cache[key] = res
	cacheMu.Unlock()
	return res
}
