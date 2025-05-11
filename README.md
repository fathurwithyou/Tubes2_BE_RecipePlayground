# Tubes2_BE_RecipePlayground
![Go 1.2](https://img.shields.io/badge/Go-1.2-blue.svg)

This project is the backend for a recipe playground application, likely focusing on generating or exploring element combinations based on provided recipes. It implements both Depth-First Search (DFS) and Breadth-First Search (BFS) algorithms to solve the recipe generation problem.

## DFS Approach

The DFS approach explores recipe combinations by going as deep as possible down one path before backtracking. This implementation includes logic to limit the total number of recipes explored based on a `maxRecipes` parameter, specifically counting recipes where both child elements are Tier 0.

## BFS Approach

The BFS approach explores recipe combinations level by level, ensuring all nodes at a certain depth are visited before moving to the next depth. Similar to the DFS approach, this implementation also respects a `maxRecipes` limit, counting recipes where both child elements are Tier 0.

## Multithreading Approach

Both DFS and BFS implementations in this project utilize aspects of concurrency, primarily for thread-safe counting and parallel processing where applicable.

### DFS Multithreading

The core DFS traversal (`dfsChunked`) is recursive and primarily single-threaded in its exploration path. However, it uses atomic operations (`atomic.AddInt64`, `atomic.LoadInt64`) for updating shared counters like `totalCount` and `visitedNodeCount`. This makes the counting mechanism thread-safe, which is important if the `Dfs` function itself were to be called concurrently or integrated into a larger system that uses multiple goroutines interacting with these counters. The traversal logic itself is not parallelized across multiple threads/goroutines in the provided implementation.

### BFS Multithreading

The BFS implementation leverages multithreading to process nodes at the *same level* of the search tree in parallel. It uses:
- `runtime.NumCPU()` to determine the number of available CPU cores, which dictates the number of worker goroutines.
- `sync.WaitGroup` to synchronize the worker goroutines, ensuring that the main thread waits for all workers to complete processing a level before moving to the next.
- Goroutines (`go func(...)`) to handle chunks of nodes from the `currentLevel`. Each goroutine processes its assigned nodes, identifies valid recipes, creates child nodes, and collects these results in a local buffer (`localEvents`).
- After all workers finish, the main thread collects the results from the local buffers (`eventsPerWorker`), builds the `nextLevel` slice, and constructs the `Children` relationships for the parent nodes.
Atomic operations (`atomic.LoadInt64`, `atomic.AddInt64`) are also used here for thread-safe access to the shared `totalCount` and `visitedNodeCount`.

This parallel processing of levels significantly speeds up the BFS traversal compared to a single-threaded approach, especially on multi-core processors.

## Folder Structure

```
Tubes2_BE_RecipePlayground/
├── Dockerfile
├── LICENSE
├── README.md
├── api
│   └── index.go
├── data
│   └── alchemy_elements.json
├── go.mod
├── go.sum
├── main.go
├── server
├── service
│   ├── model
│   │   └── element.go
│   ├── scraper
│   │   └── main.go
│   └── solver
│       ├── bfs.go
│       ├── dfs.go
│       └── init.go
└── vercel.json
```

## Usage
Run this following command in root directory, ensure the Go version is over 1.2.
```bash
go run main.go 
```

## Contributors

| Nama | NIM |
|---|---|
| Adiel Rum | 10123004 |
| Muhammad Fathur Rizky | 13523105 |
| Ahmad Wafi Idzharulhaqq | 13523131 |