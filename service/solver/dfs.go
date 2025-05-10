package solver

import (
	"sync/atomic"
)

func Dfs(rootElementName string, maxRecipes int64) interface{} {

	if elementsMapGlobal == nil {
		return map[string]string{"error": "Element data not initialized"}
	}

	var totalCount int64
	atomic.StoreInt64(&visitedNodeCount, 0)
	return dfsChunked(rootElementName, &totalCount, maxRecipes)
}

func dfsChunked(elementName string, totalCount *int64, maxRecipes int64) interface{} {
	if atomic.LoadInt64(totalCount) >= maxRecipes {
		atomic.AddInt64(&visitedNodeCount, 1)
		return elementName
	}

	eData, exists := elementsMapGlobal[elementName]

	if !exists || len(eData.Recipes) == 0 || eData.Tier == 0 {
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

		left := dfsChunked(rec[0], totalCount, maxRecipes)
		right := dfsChunked(rec[1], totalCount, maxRecipes)

		leftData, leftOk := elementsMapGlobal[rec[0]]
		rightData, rightOk := elementsMapGlobal[rec[1]]
		if leftOk && rightOk && leftData.Tier == 0 && rightData.Tier == 0 {
			atomic.AddInt64(totalCount, 1)
		}

		results = append(results, []interface{}{left, right})
	}

	return map[string]interface{}{elementName: results}
}

func GetVisitedNodeCount() int64 {
	return atomic.LoadInt64(&visitedNodeCount)
}
