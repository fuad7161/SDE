# Graphs

> Graphs model networks, maps, and dependencies. BFS for shortest paths (unweighted), DFS for connectivity and topological order.

---

## Table of Contents

1. [Graph Representation in Go](#1-graph-representation-in-go)
2. [BFS & DFS Templates](#2-bfs--dfs-templates)
3. [Number of Islands](#3-number-of-islands)
4. [Clone Graph](#4-clone-graph)
5. [Cycle Detection](#5-cycle-detection)
6. [Topological Sort](#6-topological-sort)
7. [Dijkstra's Algorithm](#7-dijkstras-algorithm)
8. [Union-Find (Disjoint Set)](#8-union-find-disjoint-set)
9. [Course Schedule](#9-course-schedule)
10. [Word Ladder](#10-word-ladder)

---

## 1. Graph Representation in Go

```go
// Adjacency list (most common in interviews)
graph := make(map[int][]int)
graph[0] = append(graph[0], 1)
graph[1] = append(graph[1], 0) // undirected

// Adjacency list with weights
type Edge struct {
    to, weight int
}
wGraph := make(map[int][]Edge)
wGraph[0] = append(wGraph[0], Edge{1, 5})

// Grid as implicit graph — 4 directions
dirs := [][2]int{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}
```

---

## 2. BFS & DFS Templates

```go
// ── BFS — shortest path in unweighted graph ──
func bfs(graph map[int][]int, start int) map[int]int {
    dist := map[int]int{start: 0}
    queue := []int{start}
    for len(queue) > 0 {
        node := queue[0]
        queue = queue[1:]
        for _, neighbor := range graph[node] {
            if _, visited := dist[neighbor]; !visited {
                dist[neighbor] = dist[node] + 1
                queue = append(queue, neighbor)
            }
        }
    }
    return dist
}

// ── DFS — recursive ──
func dfs(graph map[int][]int, node int, visited map[int]bool) {
    visited[node] = true
    for _, neighbor := range graph[node] {
        if !visited[neighbor] {
            dfs(graph, neighbor, visited)
        }
    }
}

// ── DFS — iterative ──
func dfsIterative(graph map[int][]int, start int) {
    visited := make(map[int]bool)
    stack := []int{start}
    for len(stack) > 0 {
        node := stack[len(stack)-1]
        stack = stack[:len(stack)-1]
        if visited[node] { continue }
        visited[node] = true
        for _, neighbor := range graph[node] {
            if !visited[neighbor] {
                stack = append(stack, neighbor)
            }
        }
    }
}
```

---

## 3. Number of Islands

```go
// Time: O(m*n), Space: O(m*n) call stack
func numIslands(grid [][]byte) int {
    if len(grid) == 0 { return 0 }
    rows, cols := len(grid), len(grid[0])
    count := 0

    var dfs func(r, c int)
    dfs = func(r, c int) {
        if r < 0 || r >= rows || c < 0 || c >= cols || grid[r][c] != '1' {
            return
        }
        grid[r][c] = '0' // mark visited (in-place)
        dfs(r+1, c)
        dfs(r-1, c)
        dfs(r, c+1)
        dfs(r, c-1)
    }

    for r := 0; r < rows; r++ {
        for c := 0; c < cols; c++ {
            if grid[r][c] == '1' {
                count++
                dfs(r, c)
            }
        }
    }
    return count
}
```

> **Interview Q: BFS or DFS for islands? Does it matter?**  
> Both work — same O(m*n) complexity. DFS is shorter to write. BFS avoids deep recursion on large grids. If asked "shortest path to expand an island", always use BFS.

---

## 4. Clone Graph

```go
type GraphNode struct {
    Val       int
    Neighbors []*GraphNode
}

// Time: O(n), Space: O(n)
func cloneGraph(node *GraphNode) *GraphNode {
    if node == nil { return nil }
    clones := make(map[*GraphNode]*GraphNode)

    var dfs func(*GraphNode) *GraphNode
    dfs = func(n *GraphNode) *GraphNode {
        if clone, ok := clones[n]; ok {
            return clone
        }
        clone := &GraphNode{Val: n.Val}
        clones[n] = clone
        for _, neighbor := range n.Neighbors {
            clone.Neighbors = append(clone.Neighbors, dfs(neighbor))
        }
        return clone
    }
    return dfs(node)
}
```

---

## 5. Cycle Detection

### Undirected Graph (DFS)

```go
func hasCycleUndirected(graph map[int][]int, n int) bool {
    visited := make([]bool, n)

    var dfs func(node, parent int) bool
    dfs = func(node, parent int) bool {
        visited[node] = true
        for _, neighbor := range graph[node] {
            if !visited[neighbor] {
                if dfs(neighbor, node) { return true }
            } else if neighbor != parent {
                return true // back edge
            }
        }
        return false
    }

    for i := 0; i < n; i++ {
        if !visited[i] {
            if dfs(i, -1) { return true }
        }
    }
    return false
}
```

### Directed Graph (DFS with recursion stack)

```go
func hasCycleDirected(graph map[int][]int, n int) bool {
    visited := make([]int, n) // 0=unvisited, 1=in-stack, 2=done

    var dfs func(node int) bool
    dfs = func(node int) bool {
        visited[node] = 1
        for _, neighbor := range graph[node] {
            if visited[neighbor] == 1 { return true }
            if visited[neighbor] == 0 && dfs(neighbor) { return true }
        }
        visited[node] = 2
        return false
    }

    for i := 0; i < n; i++ {
        if visited[i] == 0 && dfs(i) { return true }
    }
    return false
}
```

---

## 6. Topological Sort

### DFS-based (post-order)

```go
func topoSortDFS(graph map[int][]int, n int) []int {
    visited := make([]bool, n)
    result := []int{}

    var dfs func(node int)
    dfs = func(node int) {
        visited[node] = true
        for _, neighbor := range graph[node] {
            if !visited[neighbor] {
                dfs(neighbor)
            }
        }
        result = append(result, node) // post-order
    }

    for i := 0; i < n; i++ {
        if !visited[i] { dfs(i) }
    }

    // reverse
    for l, r := 0, len(result)-1; l < r; l, r = l+1, r-1 {
        result[l], result[r] = result[r], result[l]
    }
    return result
}
```

### Kahn's Algorithm (BFS with in-degree)

```go
func topoSortKahn(graph map[int][]int, n int) []int {
    inDegree := make([]int, n)
    for node := range graph {
        for _, neighbor := range graph[node] {
            inDegree[neighbor]++
        }
    }

    queue := []int{}
    for i := 0; i < n; i++ {
        if inDegree[i] == 0 {
            queue = append(queue, i)
        }
    }

    result := []int{}
    for len(queue) > 0 {
        node := queue[0]
        queue = queue[1:]
        result = append(result, node)
        for _, neighbor := range graph[node] {
            inDegree[neighbor]--
            if inDegree[neighbor] == 0 {
                queue = append(queue, neighbor)
            }
        }
    }

    if len(result) != n { return nil } // cycle detected
    return result
}
```

---

## 7. Dijkstra's Algorithm

```go
// Time: O((V + E) log V) with min-heap
import (
    "container/heap"
    "math"
)

type Item struct{ node, dist int }
type PQ []Item

func (pq PQ) Len() int            { return len(pq) }
func (pq PQ) Less(i, j int) bool  { return pq[i].dist < pq[j].dist }
func (pq PQ) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PQ) Push(x interface{}) { *pq = append(*pq, x.(Item)) }
func (pq *PQ) Pop() interface{} {
    old := *pq; n := len(old); x := old[n-1]; *pq = old[:n-1]; return x
}

func dijkstra(graph map[int][]Item, n, src int) []int {
    dist := make([]int, n)
    for i := range dist { dist[i] = math.MaxInt }
    dist[src] = 0

    pq := &PQ{{src, 0}}
    heap.Init(pq)

    for pq.Len() > 0 {
        curr := heap.Pop(pq).(Item)
        if curr.dist > dist[curr.node] { continue }
        for _, edge := range graph[curr.node] {
            newDist := dist[curr.node] + edge.dist
            if newDist < dist[edge.node] {
                dist[edge.node] = newDist
                heap.Push(pq, Item{edge.node, newDist})
            }
        }
    }
    return dist
}
```

---

## 8. Union-Find (Disjoint Set)

```go
// Path compression + union by rank — near O(1) amortized
type UnionFind struct {
    parent []int
    rank   []int
}

func NewUnionFind(n int) *UnionFind {
    parent := make([]int, n)
    rank := make([]int, n)
    for i := range parent { parent[i] = i }
    return &UnionFind{parent, rank}
}

func (uf *UnionFind) Find(x int) int {
    if uf.parent[x] != x {
        uf.parent[x] = uf.Find(uf.parent[x]) // path compression
    }
    return uf.parent[x]
}

func (uf *UnionFind) Union(x, y int) bool {
    px, py := uf.Find(x), uf.Find(y)
    if px == py { return false } // already connected
    if uf.rank[px] < uf.rank[py] { px, py = py, px }
    uf.parent[py] = px
    if uf.rank[px] == uf.rank[py] { uf.rank[px]++ }
    return true
}

func (uf *UnionFind) Connected(x, y int) bool {
    return uf.Find(x) == uf.Find(y)
}
```

---

## 9. Course Schedule

```go
// Can finish all courses? (Cycle detection in directed graph)
// Time: O(V + E)
func canFinish(numCourses int, prerequisites [][]int) bool {
    graph := make(map[int][]int)
    for _, pre := range prerequisites {
        graph[pre[1]] = append(graph[pre[1]], pre[0])
    }

    // 0 = unvisited, 1 = visiting, 2 = visited
    state := make([]int, numCourses)
    var dfs func(node int) bool
    dfs = func(node int) bool {
        if state[node] == 1 { return false } // cycle
        if state[node] == 2 { return true }
        state[node] = 1
        for _, next := range graph[node] {
            if !dfs(next) { return false }
        }
        state[node] = 2
        return true
    }

    for i := 0; i < numCourses; i++ {
        if !dfs(i) { return false }
    }
    return true
}
```

---

## 10. Word Ladder

```go
// Time: O(M² × N) where M = word length, N = number of words
import "strings"

func ladderLength(beginWord string, endWord string, wordList []string) int {
    wordSet := make(map[string]bool)
    for _, w := range wordList { wordSet[w] = true }
    if !wordSet[endWord] { return 0 }

    queue := []string{beginWord}
    visited := map[string]bool{beginWord: true}
    steps := 1

    for len(queue) > 0 {
        size := len(queue)
        for i := 0; i < size; i++ {
            word := queue[0]
            queue = queue[1:]
            if word == endWord { return steps }
            chars := []byte(word)
            for j := 0; j < len(chars); j++ {
                orig := chars[j]
                for c := byte('a'); c <= byte('z'); c++ {
                    if c == orig { continue }
                    chars[j] = c
                    next := string(chars)
                    if wordSet[next] && !visited[next] {
                        visited[next] = true
                        queue = append(queue, next)
                    }
                    chars[j] = orig
                }
            }
        }
        steps++
    }
    return 0
}
```

> **Interview Q: Why BFS instead of DFS for Word Ladder?**  
> BFS guarantees the *shortest* transformation sequence because it explores level by level. DFS may find a path, but not necessarily the shortest one.
