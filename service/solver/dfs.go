package solver

import (
	"sync"
	"sync/atomic"
)

var (
	cache   map[string]interface{}
	cacheMu sync.RWMutex
)

func Dfs(rootElementName string, maxRecipes int64) interface{} {
	if elementsMapGlobal == nil {
		return map[string]string{"error": "Element data not initialized"}
	}
	cacheMu.Lock()
	cache = make(map[string]interface{})
	cacheMu.Unlock()

	var totalCount int64
	visitedNodeCount = 0
	return dfsChunked(rootElementName, &totalCount, maxRecipes)
}

func dfsChunked(elementName string, totalCount *int64, maxRecipes int64) interface{} {

	cacheMu.RLock()
	if val, ok := cache[elementName]; ok {

		atomic.AddInt64(&visitedNodeCount, 1)
		cacheMu.RUnlock()
		return val
	}
	cacheMu.RUnlock()

	if atomic.LoadInt64(totalCount) >= maxRecipes {
		cacheMu.Lock()
		cache[elementName] = elementName
		cacheMu.Unlock()
		atomic.AddInt64(&visitedNodeCount, 1)
		return elementName
	}

	eData, exists := elementsMapGlobal[elementName]
	if !exists || len(eData.Recipes) == 0 || eData.Tier == 0 {
		cacheMu.Lock()
		cache[elementName] = elementName
		cacheMu.Unlock()
		atomic.AddInt64(&visitedNodeCount, 1)
		return elementName
	}

	atomic.AddInt64(&visitedNodeCount, 1)

	currentTier := eData.Tier
	recipes := make([][]string, 0, len(eData.Recipes))

	for _, rec := range eData.Recipes {
		if len(rec) != 2 {
			continue
		}
		c1, c2 := rec[0], rec[1]
		child1, ok1 := elementsMapGlobal[c1]
		child2, ok2 := elementsMapGlobal[c2]

		if !ok1 || !ok2 || child1.Tier >= currentTier || child2.Tier >= currentTier {
			continue
		}
		recipes = append(recipes, []string{c1, c2})
	}

	var results [][]interface{}
	for _, rec := range recipes {
		if atomic.LoadInt64(totalCount) >= maxRecipes {
			break
		}
		atomic.AddInt64(totalCount, 1)
		left := dfsChunked(rec[0], totalCount, maxRecipes)
		right := dfsChunked(rec[1], totalCount, maxRecipes)

		if left == nil || right == nil {
			continue
		}
		results = append(results, []interface{}{left, right})
	}

	res := map[string]interface{}{elementName: results}
	cacheMu.Lock()
	cache[elementName] = res
	cacheMu.Unlock()

	return res
}

func GetVisitedNodeCount() int64 {
	return atomic.LoadInt64(&visitedNodeCount)
}
