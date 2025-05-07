package recipetree

import (
	"fmt"
	"runtime"
	"sync"
)

var (
	cache   map[string]interface{}
	cacheMu sync.RWMutex
)

func Dfs(rootElementName string, maxDepth int) interface{} {
	if elementsMapGlobal == nil {
		return map[string]string{"error": "Element data not initialized"}
	}
	cacheMu.Lock()
	cache = make(map[string]interface{})
	cacheMu.Unlock()
	return dfsChunked(rootElementName, 0, maxDepth)
}

func dfsChunked(elementName string, currentDepth, maxDepth int) interface{} {
	key := fmt.Sprintf("%s|%d", elementName, currentDepth)

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

	var recipes = make([][]string, 0, len(eData.Recipes))
	for _, rec := range eData.Recipes {
		if len(rec) == 2 {
			recCopy := make([]string, 2)
			copy(recCopy, rec)
			recipes = append(recipes, recCopy)
		}
	}

	workers := runtime.NumCPU()
	n := len(recipes)
	chunkSize := (n + workers - 1) / workers
	results := make([][]interface{}, n)

	var wg sync.WaitGroup
	for w := 0; w < workers; w++ {
		start := w * chunkSize
		if start >= n {
			break
		}
		end := start + chunkSize
		if end > n {
			end = n
		}
		wg.Add(1)

		go func(s, e int) {
			defer wg.Done()
			for i := s; i < e; i++ {
				r := recipes[i]
				n1 := dfsChunked(r[0], currentDepth+1, maxDepth)
				n2 := dfsChunked(r[1], currentDepth+1, maxDepth)
				results[i] = []interface{}{n1, n2}
			}
		}(start, end)
	}
	wg.Wait()

	res := map[string]interface{}{elementName: results}
	cacheMu.Lock()
	cache[key] = res
	cacheMu.Unlock()
	return res
}
