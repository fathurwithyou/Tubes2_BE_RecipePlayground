package solver

import (
	"runtime"
	"sync"
	"sync/atomic"
)

var (
	combosCache map[string][][]string
	combosMu    sync.RWMutex
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

	combosMu.Lock()
	if combosCache == nil {
		combosCache = make(map[string][][]string, len(elementsMapGlobal))
	}
	combosMu.Unlock()

	root := &node{Name: rootElementName}
	currentLevel := []*node{root}
	workers := runtime.NumCPU()
	var totalCount int64

	for len(currentLevel) > 0 && atomic.LoadInt64(&totalCount) < maxRecipes {

		eventsPerWorker := make([][]bfsEvent, workers)
		for i := range eventsPerWorker {
			eventsPerWorker[i] = nil
		}

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
				local := eventsPerWorker[id]
				for _, parent := range parents {
					if atomic.LoadInt64(&totalCount) >= maxRecipes {
						break
					}

					combosMu.RLock()
					combos, ok := combosCache[parent.Name]
					combosMu.RUnlock()

					if !ok {
						eData := elementsMapGlobal[parent.Name]
						var filtered [][]string
						for _, rec := range eData.Recipes {
							if len(rec) == 2 {
								filtered = append(filtered, rec)
							}
						}
						combosMu.Lock()
						combosCache[parent.Name] = filtered
						combosMu.Unlock()
						combos = filtered
					}

					parentTier := elementsMapGlobal[parent.Name].Tier
					for _, rec := range combos {
						if atomic.LoadInt64(&totalCount) >= maxRecipes {
							break
						}
						c1, c2 := rec[0], rec[1]

						if elementsMapGlobal[c1].Tier >= parentTier || elementsMapGlobal[c2].Tier >= parentTier {
							continue
						}

						atomic.AddInt64(&totalCount, 1)
						child1 := &node{Name: c1}
						child2 := &node{Name: c2}
						local = append(local, bfsEvent{parent: parent, child1: child1, child2: child2})
					}
					eventsPerWorker[id] = local
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
