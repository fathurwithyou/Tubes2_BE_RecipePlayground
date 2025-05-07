package recipetree

import (
	"runtime"
	"sync"
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

func Bfs(rootElementName string, maxDepth int) interface{} {
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

	for depth := 0; depth < maxDepth; depth++ {
		if len(currentLevel) == 0 {
			break
		}

		eventsPerWorker := make([][]bfsEvent, workers)
		for i := range eventsPerWorker {
			eventsPerWorker[i] = []bfsEvent{}
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

			go func(id int, nodes []*node) {
				defer wg.Done()
				localEvents := eventsPerWorker[id]
				for _, parent := range nodes {

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

					for _, rec := range combos {
						c1 := &node{Name: rec[0]}
						c2 := &node{Name: rec[1]}
						localEvents = append(localEvents, bfsEvent{parent: parent, child1: c1, child2: c2})
					}
				}
				eventsPerWorker[id] = localEvents
			}(w, chunk)
		}

		wg.Wait()

		nextLevel := make([]*node, 0)
		events := make([]bfsEvent, 0)
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
	pairs := make([]interface{}, 0, len(n.Children))
	for _, children := range n.Children {
		p1 := convertNode(children[0])
		p2 := convertNode(children[1])
		pairs = append(pairs, []interface{}{p1, p2})
	}
	return map[string]interface{}{n.Name: pairs}
}
