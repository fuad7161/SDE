# Graphs

> Master BFS/DFS templates, cycle detection, and shortest path. Most graph problems reduce to one of these patterns.

---

## Table of Contents

1. [BFS & DFS Templates](#1-bfs--dfs-templates)
2. [Number of Islands](#2-number-of-islands)
3. [Clone Graph](#3-clone-graph)
4. [Detect Cycle](#4-detect-cycle)
5. [Topological Sort](#5-topological-sort)
6. [Dijkstra's Shortest Path](#6-dijkstras-shortest-path)
7. [Union-Find / Disjoint Set](#7-union-find--disjoint-set)
8. [Course Schedule (I & II)](#8-course-schedule-i--ii)
9. [Word Ladder](#9-word-ladder)

---

## Graph Representation

```java
// Adjacency List (most common for interviews)
int V = 5;
List<List<Integer>> adj = new ArrayList<>();
for (int i = 0; i < V; i++) adj.add(new ArrayList<>());

adj.get(0).add(1);  // edge 0 → 1
adj.get(0).add(2);  // edge 0 → 2
```

---

## 1. BFS & DFS Templates

### BFS (Breadth-First Search)

```java
// Time: O(V + E), Space: O(V)
public void bfs(List<List<Integer>> adj, int start) {
    boolean[] visited = new boolean[adj.size()];
    Queue<Integer> queue = new LinkedList<>();

    visited[start] = true;
    queue.offer(start);

    while (!queue.isEmpty()) {
        int node = queue.poll();
        System.out.print(node + " ");

        for (int neighbor : adj.get(node)) {
            if (!visited[neighbor]) {
                visited[neighbor] = true;
                queue.offer(neighbor);
            }
        }
    }
}
```

### DFS (Depth-First Search)

```java
// Recursive DFS
public void dfs(List<List<Integer>> adj, int node, boolean[] visited) {
    visited[node] = true;
    System.out.print(node + " ");

    for (int neighbor : adj.get(node)) {
        if (!visited[neighbor]) {
            dfs(adj, neighbor, visited);
        }
    }
}

// Iterative DFS (using explicit stack)
public void dfsIterative(List<List<Integer>> adj, int start) {
    boolean[] visited = new boolean[adj.size()];
    Deque<Integer> stack = new ArrayDeque<>();
    stack.push(start);

    while (!stack.isEmpty()) {
        int node = stack.pop();
        if (visited[node]) continue;
        visited[node] = true;
        System.out.print(node + " ");
        for (int neighbor : adj.get(node)) {
            if (!visited[neighbor]) stack.push(neighbor);
        }
    }
}
```

> **Interview Q: BFS vs DFS — when to use which?**  
> Use **BFS** for shortest path in unweighted graphs, level-order problems, or when the answer is close to the source. Use **DFS** for cycle detection, topological sort, pathfinding in mazes, and when exploring all paths. BFS uses more memory (queue can hold a whole level), DFS uses O(h) stack space.

---

## 2. Number of Islands

**Problem:** Count connected components of '1's in a 2D grid.

```java
// Time: O(m*n), Space: O(m*n)
public int numIslands(char[][] grid) {
    if (grid == null || grid.length == 0) return 0;
    int count = 0;
    int rows = grid.length, cols = grid[0].length;

    for (int r = 0; r < rows; r++) {
        for (int c = 0; c < cols; c++) {
            if (grid[r][c] == '1') {
                count++;
                dfsSink(grid, r, c);  // "sink" the island to avoid revisiting
            }
        }
    }
    return count;
}

private void dfsSink(char[][] grid, int r, int c) {
    if (r < 0 || c < 0 || r >= grid.length || c >= grid[0].length || grid[r][c] != '1')
        return;
    grid[r][c] = '0';           // mark visited by sinking
    dfsSink(grid, r + 1, c);
    dfsSink(grid, r - 1, c);
    dfsSink(grid, r, c + 1);
    dfsSink(grid, r, c - 1);
}
```

---

## 3. Clone Graph

**Problem:** Deep copy a graph where each node has a value and a list of neighbors.

```java
class Node {
    int val;
    List<Node> neighbors;
    Node(int val) { this.val = val; this.neighbors = new ArrayList<>(); }
}

// HashMap to map original node → cloned node (handles cycles)
// Time: O(V + E), Space: O(V)
public Node cloneGraph(Node node) {
    if (node == null) return null;
    Map<Node, Node> visited = new HashMap<>();
    return dfsClone(node, visited);
}

private Node dfsClone(Node node, Map<Node, Node> visited) {
    if (visited.containsKey(node)) return visited.get(node);

    Node clone = new Node(node.val);
    visited.put(node, clone);                          // register before recursing (handles cycles)

    for (Node neighbor : node.neighbors) {
        clone.neighbors.add(dfsClone(neighbor, visited));
    }
    return clone;
}
```

---

## 4. Detect Cycle

### Undirected Graph — DFS with parent tracking

```java
public boolean hasCycleUndirected(List<List<Integer>> adj, int V) {
    boolean[] visited = new boolean[V];
    for (int i = 0; i < V; i++) {
        if (!visited[i] && dfsCycleUndirected(adj, i, -1, visited)) return true;
    }
    return false;
}

private boolean dfsCycleUndirected(List<List<Integer>> adj, int node, int parent, boolean[] visited) {
    visited[node] = true;
    for (int neighbor : adj.get(node)) {
        if (!visited[neighbor]) {
            if (dfsCycleUndirected(adj, neighbor, node, visited)) return true;
        } else if (neighbor != parent) {
            return true;   // back edge to a non-parent visited node = cycle
        }
    }
    return false;
}
```

### Directed Graph — DFS with recursion stack (color states)

```java
// 0 = unvisited, 1 = in current DFS path, 2 = fully processed
public boolean hasCycleDirected(List<List<Integer>> adj, int V) {
    int[] color = new int[V];
    for (int i = 0; i < V; i++) {
        if (color[i] == 0 && dfsCycleDirected(adj, i, color)) return true;
    }
    return false;
}

private boolean dfsCycleDirected(List<List<Integer>> adj, int node, int[] color) {
    color[node] = 1;   // mark as in-progress
    for (int neighbor : adj.get(node)) {
        if (color[neighbor] == 1) return true;   // back edge = cycle
        if (color[neighbor] == 0 && dfsCycleDirected(adj, neighbor, color)) return true;
    }
    color[node] = 2;   // fully processed
    return false;
}
```

---

## 5. Topological Sort

**Only for Directed Acyclic Graphs (DAG).** Ordering where every edge `u → v` has `u` before `v`.

### Kahn's Algorithm (BFS — in-degree based)

```java
// Time: O(V + E), Space: O(V)
public int[] topologicalSort(List<List<Integer>> adj, int V) {
    int[] inDegree = new int[V];
    for (int u = 0; u < V; u++)
        for (int v : adj.get(u)) inDegree[v]++;

    Queue<Integer> queue = new LinkedList<>();
    for (int i = 0; i < V; i++) if (inDegree[i] == 0) queue.offer(i);

    int[] order = new int[V];
    int idx = 0;
    while (!queue.isEmpty()) {
        int node = queue.poll();
        order[idx++] = node;
        for (int neighbor : adj.get(node)) {
            inDegree[neighbor]--;
            if (inDegree[neighbor] == 0) queue.offer(neighbor);
        }
    }
    return (idx == V) ? order : new int[0];   // empty if cycle detected
}
```

### DFS-based Topological Sort

```java
public List<Integer> topologicalSortDFS(List<List<Integer>> adj, int V) {
    boolean[] visited = new boolean[V];
    Deque<Integer> stack = new ArrayDeque<>();

    for (int i = 0; i < V; i++) {
        if (!visited[i]) dfsTopoSort(adj, i, visited, stack);
    }

    List<Integer> order = new ArrayList<>();
    while (!stack.isEmpty()) order.add(stack.pop());
    return order;
}

private void dfsTopoSort(List<List<Integer>> adj, int node, boolean[] visited, Deque<Integer> stack) {
    visited[node] = true;
    for (int neighbor : adj.get(node)) {
        if (!visited[neighbor]) dfsTopoSort(adj, neighbor, visited, stack);
    }
    stack.push(node);   // push AFTER processing all descendants
}
```

---

## 6. Dijkstra's Shortest Path

**Shortest path in a weighted graph with non-negative edges.**

```java
// Time: O((V + E) log V) with priority queue, Space: O(V)
public int[] dijkstra(List<List<int[]>> adj, int src, int V) {
    int[] dist = new int[V];
    Arrays.fill(dist, Integer.MAX_VALUE);
    dist[src] = 0;

    // Priority queue: [distance, node] — min-heap by distance
    PriorityQueue<int[]> pq = new PriorityQueue<>((a, b) -> a[0] - b[0]);
    pq.offer(new int[]{0, src});

    while (!pq.isEmpty()) {
        int[] curr = pq.poll();
        int d = curr[0], u = curr[1];

        if (d > dist[u]) continue;    // stale entry — skip

        for (int[] edge : adj.get(u)) {
            int v = edge[0], weight = edge[1];
            if (dist[u] + weight < dist[v]) {
                dist[v] = dist[u] + weight;
                pq.offer(new int[]{dist[v], v});
            }
        }
    }
    return dist;
}
```

> **Interview Q: Why doesn't Dijkstra work with negative weights?**  
> Dijkstra assumes that once a node is popped from the min-heap, its shortest path is finalized. A negative edge could create a shorter path to an already-processed node. Use **Bellman-Ford** (O(VE)) for graphs with negative edges.

---

## 7. Union-Find / Disjoint Set

**Efficiently answers "are two nodes in the same connected component?"** Supports union and find in near O(1) with path compression + union by rank.

```java
class UnionFind {
    private int[] parent, rank;
    private int components;

    public UnionFind(int n) {
        parent = new int[n];
        rank = new int[n];
        components = n;
        for (int i = 0; i < n; i++) parent[i] = i;  // each node is its own root
    }

    // Find with path compression
    public int find(int x) {
        if (parent[x] != x) parent[x] = find(parent[x]);  // path compress
        return parent[x];
    }

    // Union by rank
    public boolean union(int x, int y) {
        int px = find(x), py = find(y);
        if (px == py) return false;   // already connected
        if (rank[px] < rank[py])      parent[px] = py;
        else if (rank[px] > rank[py]) parent[py] = px;
        else { parent[py] = px; rank[px]++; }
        components--;
        return true;
    }

    public int getComponents() { return components; }
    public boolean connected(int x, int y) { return find(x) == find(y); }
}

// Usage: Number of connected components
int countComponents(int n, int[][] edges) {
    UnionFind uf = new UnionFind(n);
    for (int[] e : edges) uf.union(e[0], e[1]);
    return uf.getComponents();
}
```

---

## 8. Course Schedule (I & II)

### Course Schedule I — Can all courses be finished? (cycle detection in directed graph)

```java
public boolean canFinish(int numCourses, int[][] prerequisites) {
    List<List<Integer>> adj = new ArrayList<>();
    for (int i = 0; i < numCourses; i++) adj.add(new ArrayList<>());
    int[] inDegree = new int[numCourses];

    for (int[] pre : prerequisites) {
        adj.get(pre[1]).add(pre[0]);
        inDegree[pre[0]]++;
    }

    Queue<Integer> queue = new LinkedList<>();
    for (int i = 0; i < numCourses; i++) if (inDegree[i] == 0) queue.offer(i);

    int completed = 0;
    while (!queue.isEmpty()) {
        int course = queue.poll();
        completed++;
        for (int next : adj.get(course)) {
            if (--inDegree[next] == 0) queue.offer(next);
        }
    }
    return completed == numCourses;   // if cycle exists, some nodes never reach in-degree 0
}
```

### Course Schedule II — Return topological order

```java
public int[] findOrder(int numCourses, int[][] prerequisites) {
    List<List<Integer>> adj = new ArrayList<>();
    int[] inDegree = new int[numCourses];
    for (int i = 0; i < numCourses; i++) adj.add(new ArrayList<>());
    for (int[] pre : prerequisites) { adj.get(pre[1]).add(pre[0]); inDegree[pre[0]]++; }

    Queue<Integer> queue = new LinkedList<>();
    for (int i = 0; i < numCourses; i++) if (inDegree[i] == 0) queue.offer(i);

    int[] order = new int[numCourses];
    int idx = 0;
    while (!queue.isEmpty()) {
        int course = queue.poll();
        order[idx++] = course;
        for (int next : adj.get(course)) if (--inDegree[next] == 0) queue.offer(next);
    }
    return idx == numCourses ? order : new int[0];
}
```

---

## 9. Word Ladder

**Problem:** Transform `beginWord` to `endWord` one letter at a time (each intermediate word must be in `wordList`). Return length of shortest transformation.

```java
// BFS — each level = one transformation step
// Time: O(M² * N) where M = word length, N = wordList size
public int ladderLength(String beginWord, String endWord, List<String> wordList) {
    Set<String> wordSet = new HashSet<>(wordList);
    if (!wordSet.contains(endWord)) return 0;

    Queue<String> queue = new LinkedList<>();
    queue.offer(beginWord);
    int steps = 1;

    while (!queue.isEmpty()) {
        int size = queue.size();
        for (int i = 0; i < size; i++) {
            String word = queue.poll();
            char[] chars = word.toCharArray();

            for (int j = 0; j < chars.length; j++) {
                char original = chars[j];
                for (char c = 'a'; c <= 'z'; c++) {
                    chars[j] = c;
                    String next = new String(chars);
                    if (next.equals(endWord)) return steps + 1;
                    if (wordSet.contains(next)) {
                        queue.offer(next);
                        wordSet.remove(next);  // mark visited by removing from set
                    }
                }
                chars[j] = original;  // restore
            }
        }
        steps++;
    }
    return 0;
}
```
