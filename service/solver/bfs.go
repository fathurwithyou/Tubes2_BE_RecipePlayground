package solver

import (
	"runtime"
	"sync"
	"sync/atomic"
)

var (
	visitedNodeCount int64
)

type node struct {
	Name     string
	Children [][]*node
}

type bfsEvent struct {
	parent *node
	child1 *node
	child2 *node
}

func Bfs(rootElementName string, maxRecipes int64) interface{} {

	if elementsMapGlobal == nil {
		return map[string]string{"error": "Element data not initialized"}
	}

	root := &node{Name: rootElementName}
	currentLevel := []*node{root}
	workers := runtime.NumCPU()
	var totalCount int64

	atomic.StoreInt64(&visitedNodeCount, 0)

	for len(currentLevel) > 0 && atomic.LoadInt64(&totalCount) < maxRecipes {

		eventsPerWorker := make([][]bfsEvent, workers)

		var wg sync.WaitGroup
		wg.Add(workers)

		n := len(currentLevel)

		chunkSize := (n + workers - 1) / workers

		for w := 0; w < workers; w++ {
			start := w * chunkSize
			end := start + chunkSize

			if start >= n {
				wg.Done()
				continue
			}

			if end > n {
				end = n
			}
			chunk := currentLevel[start:end]

			go func(id int, parents []*node) {
				defer wg.Done()
				localEvents := eventsPerWorker[id]

				for _, parent := range parents {

					if atomic.LoadInt64(&totalCount) >= maxRecipes {
						break
					}

					atomic.AddInt64(&visitedNodeCount, 1)

					eData, exists := elementsMapGlobal[parent.Name]

					if !exists || len(eData.Recipes) == 0 || eData.Tier == 0 {

						continue
					}

					parentTier := eData.Tier

					for _, rec := range eData.Recipes {

						if atomic.LoadInt64(&totalCount) >= maxRecipes {
							break
						}

						if len(rec) != 2 {
							continue
						}
						c1Name, c2Name := rec[0], rec[1]

						child1Data, ok1 := elementsMapGlobal[c1Name]
						child2Data, ok2 := elementsMapGlobal[c2Name]

						if !ok1 || !ok2 || child1Data.Tier >= parentTier || child2Data.Tier >= parentTier {
							continue
						}

						child1 := &node{Name: c1Name}
						child2 := &node{Name: c2Name}

						localEvents = append(localEvents, bfsEvent{parent: parent, child1: child1, child2: child2})

						if child1Data.Tier == 0 && child2Data.Tier == 0 {
							atomic.AddInt64(&totalCount, 1)
						}

						if atomic.LoadInt64(&totalCount) >= maxRecipes {
							break
						}
					}
					eventsPerWorker[id] = localEvents
				}
			}(w, chunk)
		}
		wg.Wait()

		nextLevel := make([]*node, 0)
		events := []bfsEvent{}
		for _, buf := range eventsPerWorker {
			events = append(events, buf...)
		}

		for _, ev := range events {

			ev.parent.Children = append(ev.parent.Children, []*node{ev.child1, ev.child2})

			nextLevel = append(nextLevel, ev.child1, ev.child2)
		}

		currentLevel = nextLevel
	}

	return convertNode(root)
}

func convertNode(n *node) interface{} {

	if len(n.Children) == 0 {
		return n.Name
	}

	pairs := make([]interface{}, len(n.Children))
	for i, children := range n.Children {

		p1 := convertNode(children[0])
		p2 := convertNode(children[1])

		pairs[i] = []interface{}{p1, p2}
	}

	return map[string]interface{}{n.Name: pairs}
}
