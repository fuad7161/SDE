# Trees

> Trees appear in ~25% of coding interviews. Master recursive thinking, BFS (level-order), and BST properties.

---

## Table of Contents

1. [Tree Node Definition](#1-tree-node-definition)
2. [Tree Traversals](#2-tree-traversals)
3. [Level Order / Zigzag Traversal](#3-level-order--zigzag-traversal)
4. [Height & Diameter](#4-height--diameter)
5. [Lowest Common Ancestor (LCA)](#5-lowest-common-ancestor-lca)
6. [Validate BST](#6-validate-bst)
7. [Serialize & Deserialize Binary Tree](#7-serialize--deserialize-binary-tree)
8. [Path Sum I, II, III](#8-path-sum-i-ii-iii)
9. [Right Side View](#9-right-side-view)
10. [Balanced Binary Tree](#10-balanced-binary-tree)

---

## 1. Tree Node Definition

```go
type TreeNode struct {
    Val   int
    Left  *TreeNode
    Right *TreeNode
}

// Helper — build tree from level-order slice (LeetCode style)
// e.g. [3, 9, 20, nil, nil, 15, 7]
```

---

## 2. Tree Traversals

```go
// ── Inorder (Left → Root → Right) ──
func inorder(root *TreeNode, result *[]int) {
    if root == nil { return }
    inorder(root.Left, result)
    *result = append(*result, root.Val)
    inorder(root.Right, result)
}

// ── Preorder (Root → Left → Right) ──
func preorder(root *TreeNode, result *[]int) {
    if root == nil { return }
    *result = append(*result, root.Val)
    preorder(root.Left, result)
    preorder(root.Right, result)
}

// ── Postorder (Left → Right → Root) ──
func postorder(root *TreeNode, result *[]int) {
    if root == nil { return }
    postorder(root.Left, result)
    postorder(root.Right, result)
    *result = append(*result, root.Val)
}

// ── Iterative Inorder (important for interviews) ──
func inorderIterative(root *TreeNode) []int {
    result := []int{}
    stack := []*TreeNode{}
    cur := root
    for cur != nil || len(stack) > 0 {
        for cur != nil {
            stack = append(stack, cur)
            cur = cur.Left
        }
        cur = stack[len(stack)-1]
        stack = stack[:len(stack)-1]
        result = append(result, cur.Val)
        cur = cur.Right
    }
    return result
}
```

> **Interview Q: When is iterative traversal better than recursive?**  
> Recursive uses the call stack (O(h) space). Iterative gives you explicit control — crucial for very deep trees and for Morris traversal (O(1) space) in follow-ups.

---

## 3. Level Order / Zigzag Traversal

```go
// Level Order BFS — Time: O(n), Space: O(n)
func levelOrder(root *TreeNode) [][]int {
    result := [][]int{}
    if root == nil { return result }

    queue := []*TreeNode{root}
    for len(queue) > 0 {
        size := len(queue)
        level := []int{}
        for i := 0; i < size; i++ {
            node := queue[0]
            queue = queue[1:]
            level = append(level, node.Val)
            if node.Left != nil  { queue = append(queue, node.Left) }
            if node.Right != nil { queue = append(queue, node.Right) }
        }
        result = append(result, level)
    }
    return result
}

// Zigzag Level Order
func zigzagLevelOrder(root *TreeNode) [][]int {
    result := [][]int{}
    if root == nil { return result }

    queue := []*TreeNode{root}
    leftToRight := true

    for len(queue) > 0 {
        size := len(queue)
        level := make([]int, size)
        for i := 0; i < size; i++ {
            node := queue[0]
            queue = queue[1:]
            pos := i
            if !leftToRight { pos = size - 1 - i }
            level[pos] = node.Val
            if node.Left != nil  { queue = append(queue, node.Left) }
            if node.Right != nil { queue = append(queue, node.Right) }
        }
        result = append(result, level)
        leftToRight = !leftToRight
    }
    return result
}
```

---

## 4. Height & Diameter

```go
// Height (max depth) — Time: O(n)
func maxDepth(root *TreeNode) int {
    if root == nil { return 0 }
    left := maxDepth(root.Left)
    right := maxDepth(root.Right)
    if left > right { return left + 1 }
    return right + 1
}

// Diameter — longest path (may not pass through root)
func diameterOfBinaryTree(root *TreeNode) int {
    maxDiam := 0
    var height func(*TreeNode) int
    height = func(node *TreeNode) int {
        if node == nil { return 0 }
        left := height(node.Left)
        right := height(node.Right)
        if left+right > maxDiam { maxDiam = left + right }
        if left > right { return left + 1 }
        return right + 1
    }
    height(root)
    return maxDiam
}
```

---

## 5. Lowest Common Ancestor (LCA)

```go
// LCA in Binary Tree — Time: O(n)
func lowestCommonAncestor(root, p, q *TreeNode) *TreeNode {
    if root == nil || root == p || root == q {
        return root
    }
    left := lowestCommonAncestor(root.Left, p, q)
    right := lowestCommonAncestor(root.Right, p, q)
    if left != nil && right != nil {
        return root // p and q on different sides
    }
    if left != nil { return left }
    return right
}

// LCA in BST — O(h) using BST property
func lowestCommonAncestorBST(root, p, q *TreeNode) *TreeNode {
    for root != nil {
        if p.Val < root.Val && q.Val < root.Val {
            root = root.Left
        } else if p.Val > root.Val && q.Val > root.Val {
            root = root.Right
        } else {
            return root
        }
    }
    return nil
}
```

---

## 6. Validate BST

```go
// Time: O(n), Space: O(h)
import "math"

func isValidBST(root *TreeNode) bool {
    var validate func(*TreeNode, int, int) bool
    validate = func(node *TreeNode, min, max int) bool {
        if node == nil { return true }
        if node.Val <= min || node.Val >= max { return false }
        return validate(node.Left, min, node.Val) &&
               validate(node.Right, node.Val, max)
    }
    return validate(root, math.MinInt, math.MaxInt)
}
```

> **Interview Q: Can you validate a BST using inorder traversal?**  
> Yes — an inorder traversal of a valid BST produces a strictly increasing sequence. Track the previous value and return false if the current node is not strictly greater.

---

## 7. Serialize & Deserialize Binary Tree

```go
import (
    "strconv"
    "strings"
)

type Codec struct{}

// Preorder — Time: O(n), Space: O(n)
func (c *Codec) serialize(root *TreeNode) string {
    if root == nil { return "N" }
    left := c.serialize(root.Left)
    right := c.serialize(root.Right)
    return strconv.Itoa(root.Val) + "," + left + "," + right
}

func (c *Codec) deserialize(data string) *TreeNode {
    tokens := strings.Split(data, ",")
    idx := 0
    var build func() *TreeNode
    build = func() *TreeNode {
        if tokens[idx] == "N" {
            idx++
            return nil
        }
        val, _ := strconv.Atoi(tokens[idx])
        idx++
        node := &TreeNode{Val: val}
        node.Left = build()
        node.Right = build()
        return node
    }
    return build()
}
```

---

## 8. Path Sum I, II, III

```go
// Path Sum I — any root-to-leaf path equals target
func hasPathSum(root *TreeNode, targetSum int) bool {
    if root == nil { return false }
    if root.Left == nil && root.Right == nil {
        return root.Val == targetSum
    }
    return hasPathSum(root.Left, targetSum-root.Val) ||
           hasPathSum(root.Right, targetSum-root.Val)
}

// Path Sum II — all root-to-leaf paths
func pathSum(root *TreeNode, targetSum int) [][]int {
    result := [][]int{}
    var dfs func(*TreeNode, int, []int)
    dfs = func(node *TreeNode, remain int, path []int) {
        if node == nil { return }
        path = append(path, node.Val)
        if node.Left == nil && node.Right == nil && remain == node.Val {
            // copy path to avoid mutation
            tmp := make([]int, len(path))
            copy(tmp, path)
            result = append(result, tmp)
            return
        }
        dfs(node.Left, remain-node.Val, path)
        dfs(node.Right, remain-node.Val, path)
    }
    dfs(root, targetSum, []int{})
    return result
}

// Path Sum III — any path (not just root-to-leaf) equals target
// Use prefix sum map — Time: O(n)
func pathSumIII(root *TreeNode, targetSum int) int {
    prefixCount := map[int]int{0: 1}
    count := 0
    var dfs func(*TreeNode, int)
    dfs = func(node *TreeNode, currentSum int) {
        if node == nil { return }
        currentSum += node.Val
        count += prefixCount[currentSum-targetSum]
        prefixCount[currentSum]++
        dfs(node.Left, currentSum)
        dfs(node.Right, currentSum)
        prefixCount[currentSum]-- // backtrack
    }
    dfs(root, 0)
    return count
}
```

> **Interview Q: Why do we need to copy the path slice in Path Sum II?**  
> In Go, slices share underlying arrays. Without copying, all appended paths would reference the same backing array and get overwritten as recursion continues.

---

## 9. Right Side View

```go
// BFS — take last element of each level
func rightSideView(root *TreeNode) []int {
    result := []int{}
    if root == nil { return result }

    queue := []*TreeNode{root}
    for len(queue) > 0 {
        size := len(queue)
        for i := 0; i < size; i++ {
            node := queue[0]
            queue = queue[1:]
            if i == size-1 { result = append(result, node.Val) }
            if node.Left != nil  { queue = append(queue, node.Left) }
            if node.Right != nil { queue = append(queue, node.Right) }
        }
    }
    return result
}
```

---

## 10. Balanced Binary Tree

```go
// Time: O(n) — check and compute height in one pass
func isBalanced(root *TreeNode) bool {
    var check func(*TreeNode) int
    check = func(node *TreeNode) int {
        if node == nil { return 0 }
        left := check(node.Left)
        if left == -1 { return -1 }
        right := check(node.Right)
        if right == -1 { return -1 }
        diff := left - right
        if diff < -1 || diff > 1 { return -1 }
        if left > right { return left + 1 }
        return right + 1
    }
    return check(root) != -1
}
```
