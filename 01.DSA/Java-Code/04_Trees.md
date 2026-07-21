# Trees

> Binary trees and BSTs are core interview topics. Master DFS (recursive), BFS (queue), and the "left/right recursion" pattern.

---

## Table of Contents

1. [Tree Node Definition](#tree-node-definition)
2. [Inorder / Preorder / Postorder Traversal](#1-inorder--preorder--postorder-traversal)
3. [Level Order Traversal (BFS)](#2-level-order-traversal-bfs)
4. [Height & Diameter of Tree](#3-height--diameter-of-tree)
5. [Lowest Common Ancestor (LCA)](#4-lowest-common-ancestor-lca)
6. [Validate BST](#5-validate-bst)
7. [Serialize and Deserialize Binary Tree](#6-serialize-and-deserialize-binary-tree)
8. [Path Sum Problems](#7-path-sum-problems)
9. [Right / Left Side View](#8-right--left-side-view)
10. [Balanced Binary Tree](#9-balanced-binary-tree)

---

## Tree Node Definition

```java
class TreeNode {
    int val;
    TreeNode left, right;
    TreeNode(int val) { this.val = val; }
}
```

---

## 1. Inorder / Preorder / Postorder Traversal

### Recursive

```java
// Inorder: Left → Root → Right  (gives sorted order for BST)
public void inorder(TreeNode root, List<Integer> result) {
    if (root == null) return;
    inorder(root.left, result);
    result.add(root.val);
    inorder(root.right, result);
}

// Preorder: Root → Left → Right  (used for tree copying/serialization)
public void preorder(TreeNode root, List<Integer> result) {
    if (root == null) return;
    result.add(root.val);
    preorder(root.left, result);
    preorder(root.right, result);
}

// Postorder: Left → Right → Root  (used for tree deletion, computing subtree info)
public void postorder(TreeNode root, List<Integer> result) {
    if (root == null) return;
    postorder(root.left, result);
    postorder(root.right, result);
    result.add(root.val);
}
```

### Iterative Inorder (stack-based)

```java
// Time: O(n), Space: O(h) where h = height
public List<Integer> inorderIterative(TreeNode root) {
    List<Integer> result = new ArrayList<>();
    Deque<TreeNode> stack = new ArrayDeque<>();
    TreeNode curr = root;

    while (curr != null || !stack.isEmpty()) {
        while (curr != null) {
            stack.push(curr);
            curr = curr.left;      // go as far left as possible
        }
        curr = stack.pop();
        result.add(curr.val);      // process node
        curr = curr.right;         // move to right subtree
    }
    return result;
}
```

### Iterative Preorder (stack-based)

```java
public List<Integer> preorderIterative(TreeNode root) {
    List<Integer> result = new ArrayList<>();
    if (root == null) return result;
    Deque<TreeNode> stack = new ArrayDeque<>();
    stack.push(root);

    while (!stack.isEmpty()) {
        TreeNode node = stack.pop();
        result.add(node.val);
        if (node.right != null) stack.push(node.right);  // push right first
        if (node.left != null)  stack.push(node.left);   // left processed first
    }
    return result;
}
```

> **Interview Q: What are the use cases of each traversal?**  
> **Inorder** — prints BST in sorted order. **Preorder** — used for serialization (root before children). **Postorder** — compute subtree values (child results used by parent), tree deletion.

---

## 2. Level Order Traversal (BFS)

```java
// Time: O(n), Space: O(n)
public List<List<Integer>> levelOrder(TreeNode root) {
    List<List<Integer>> result = new ArrayList<>();
    if (root == null) return result;

    Queue<TreeNode> queue = new LinkedList<>();
    queue.offer(root);

    while (!queue.isEmpty()) {
        int size = queue.size();           // number of nodes at current level
        List<Integer> level = new ArrayList<>();

        for (int i = 0; i < size; i++) {
            TreeNode node = queue.poll();
            level.add(node.val);
            if (node.left  != null) queue.offer(node.left);
            if (node.right != null) queue.offer(node.right);
        }
        result.add(level);
    }
    return result;
}

// ── Zigzag level order (alternate left-right) ──
public List<List<Integer>> zigzagLevelOrder(TreeNode root) {
    List<List<Integer>> result = new ArrayList<>();
    if (root == null) return result;
    Queue<TreeNode> queue = new LinkedList<>();
    queue.offer(root);
    boolean leftToRight = true;

    while (!queue.isEmpty()) {
        int size = queue.size();
        Deque<Integer> level = new ArrayDeque<>();
        for (int i = 0; i < size; i++) {
            TreeNode node = queue.poll();
            if (leftToRight) level.addLast(node.val);
            else             level.addFirst(node.val);
            if (node.left  != null) queue.offer(node.left);
            if (node.right != null) queue.offer(node.right);
        }
        result.add(new ArrayList<>(level));
        leftToRight = !leftToRight;
    }
    return result;
}
```

---

## 3. Height & Diameter of Tree

### Height (max depth)

```java
// Time: O(n), Space: O(h)
public int maxDepth(TreeNode root) {
    if (root == null) return 0;
    return 1 + Math.max(maxDepth(root.left), maxDepth(root.right));
}
```

### Diameter (longest path between any two nodes)

```java
// The diameter at each node = leftHeight + rightHeight
// It might not pass through root!
// Time: O(n), Space: O(h)
private int diameter = 0;

public int diameterOfBinaryTree(TreeNode root) {
    diameter = 0;
    height(root);
    return diameter;
}

private int height(TreeNode node) {
    if (node == null) return 0;
    int left  = height(node.left);
    int right = height(node.right);
    diameter = Math.max(diameter, left + right);  // update global max
    return 1 + Math.max(left, right);             // return height to parent
}
```

---

## 4. Lowest Common Ancestor (LCA)

### LCA of Binary Tree (general)

```java
// Time: O(n), Space: O(h)
public TreeNode lowestCommonAncestor(TreeNode root, TreeNode p, TreeNode q) {
    if (root == null || root == p || root == q) return root;

    TreeNode left  = lowestCommonAncestor(root.left, p, q);
    TreeNode right = lowestCommonAncestor(root.right, p, q);

    if (left != null && right != null) return root;  // p in left, q in right
    return (left != null) ? left : right;             // both in same subtree
}
```

### LCA of BST (use BST property)

```java
// For BST: if both p,q < root → go left; if both > root → go right; else root is LCA
public TreeNode lcaBST(TreeNode root, TreeNode p, TreeNode q) {
    if (p.val < root.val && q.val < root.val)  return lcaBST(root.left, p, q);
    if (p.val > root.val && q.val > root.val)  return lcaBST(root.right, p, q);
    return root;
}
```

---

## 5. Validate BST

**Rule:** Every node in the left subtree must be strictly less than the node, and every node in the right subtree must be strictly greater — checked against the full valid range, not just the parent.

```java
// Time: O(n), Space: O(h)
public boolean isValidBST(TreeNode root) {
    return validate(root, Long.MIN_VALUE, Long.MAX_VALUE);
}

private boolean validate(TreeNode node, long min, long max) {
    if (node == null) return true;
    if (node.val <= min || node.val >= max) return false;
    return validate(node.left,  min,       node.val) &&
           validate(node.right, node.val,  max);
}
// Use Long to avoid edge cases with Integer.MIN_VALUE / Integer.MAX_VALUE
```

> **Interview Q: Why pass min/max bounds rather than just checking parent?**  
> Checking only against the immediate parent misses deeper violations. E.g., in a tree where root=5, left=3, and left.right=7 — checking only parent says 7 > 3 (valid), but 7 > 5 (root) makes it an invalid BST. The bounds propagate the full valid range down.

---

## 6. Serialize and Deserialize Binary Tree

```java
// Preorder serialization with null markers
public class Codec {
    private static final String NULL = "#";
    private static final String SEP = ",";

    // Serialize: preorder DFS
    public String serialize(TreeNode root) {
        StringBuilder sb = new StringBuilder();
        serializeHelper(root, sb);
        return sb.toString();
    }

    private void serializeHelper(TreeNode node, StringBuilder sb) {
        if (node == null) { sb.append(NULL).append(SEP); return; }
        sb.append(node.val).append(SEP);
        serializeHelper(node.left, sb);
        serializeHelper(node.right, sb);
    }

    // Deserialize: reconstruct from preorder string
    public TreeNode deserialize(String data) {
        Deque<String> queue = new ArrayDeque<>(Arrays.asList(data.split(SEP)));
        return deserializeHelper(queue);
    }

    private TreeNode deserializeHelper(Deque<String> queue) {
        String val = queue.poll();
        if (NULL.equals(val)) return null;
        TreeNode node = new TreeNode(Integer.parseInt(val));
        node.left  = deserializeHelper(queue);
        node.right = deserializeHelper(queue);
        return node;
    }
}
// Tree: 1 → 2,3 → 4,5
// Serialized: "1,2,4,#,#,5,#,#,3,#,#,"
```

---

## 7. Path Sum Problems

### Path Sum I — Does any root-to-leaf path sum equal target?

```java
public boolean hasPathSum(TreeNode root, int targetSum) {
    if (root == null) return false;
    if (root.left == null && root.right == null) return root.val == targetSum;
    return hasPathSum(root.left,  targetSum - root.val) ||
           hasPathSum(root.right, targetSum - root.val);
}
```

### Path Sum II — All root-to-leaf paths with target sum

```java
public List<List<Integer>> pathSum(TreeNode root, int target) {
    List<List<Integer>> result = new ArrayList<>();
    dfs(root, target, new ArrayList<>(), result);
    return result;
}

private void dfs(TreeNode node, int remaining, List<Integer> path, List<List<Integer>> result) {
    if (node == null) return;
    path.add(node.val);
    if (node.left == null && node.right == null && remaining == node.val) {
        result.add(new ArrayList<>(path));   // deep copy — backtrack after
    }
    dfs(node.left,  remaining - node.val, path, result);
    dfs(node.right, remaining - node.val, path, result);
    path.remove(path.size() - 1);           // backtrack
}
```

### Path Sum III — Any path (not just root-to-leaf) equal to target

```java
// Prefix sum + HashMap — Time: O(n), Space: O(n)
public int pathSumIII(TreeNode root, int targetSum) {
    Map<Long, Integer> prefixCount = new HashMap<>();
    prefixCount.put(0L, 1);  // empty path
    return dfs(root, 0L, targetSum, prefixCount);
}

private int dfs(TreeNode node, long currSum, int target, Map<Long, Integer> map) {
    if (node == null) return 0;
    currSum += node.val;
    int count = map.getOrDefault(currSum - target, 0);
    map.merge(currSum, 1, Integer::sum);
    count += dfs(node.left, currSum, target, map);
    count += dfs(node.right, currSum, target, map);
    map.merge(currSum, -1, Integer::sum);  // backtrack
    return count;
}
```

---

## 8. Right / Left Side View

```java
// BFS — take the last element of each level
public List<Integer> rightSideView(TreeNode root) {
    List<Integer> result = new ArrayList<>();
    if (root == null) return result;
    Queue<TreeNode> queue = new LinkedList<>();
    queue.offer(root);

    while (!queue.isEmpty()) {
        int size = queue.size();
        for (int i = 0; i < size; i++) {
            TreeNode node = queue.poll();
            if (i == size - 1) result.add(node.val);  // last in level = rightmost
            if (node.left  != null) queue.offer(node.left);
            if (node.right != null) queue.offer(node.right);
        }
    }
    return result;
}

// DFS approach — visit right before left
public void rightSideViewDFS(TreeNode node, int depth, List<Integer> result) {
    if (node == null) return;
    if (depth == result.size()) result.add(node.val);  // first node at this depth
    rightSideViewDFS(node.right, depth + 1, result);   // right first
    rightSideViewDFS(node.left,  depth + 1, result);
}
```

---

## 9. Balanced Binary Tree

**A tree is balanced if** for every node, `|height(left) - height(right)| <= 1`.

```java
// Time: O(n) — compute height once per node and check balance simultaneously
public boolean isBalanced(TreeNode root) {
    return checkHeight(root) != -1;
}

// Returns -1 if unbalanced, otherwise returns height
private int checkHeight(TreeNode node) {
    if (node == null) return 0;

    int leftH  = checkHeight(node.left);
    if (leftH == -1) return -1;  // early exit — already unbalanced

    int rightH = checkHeight(node.right);
    if (rightH == -1) return -1;

    if (Math.abs(leftH - rightH) > 1) return -1;  // unbalanced at this node
    return 1 + Math.max(leftH, rightH);            // return height
}
```

> **Interview Q: Why is the naive O(n²) approach wrong?**  
> Computing `height()` from scratch at every node costs O(n) per node = O(n²) total. The O(n) approach checks balance and computes height in the **same DFS pass**, using -1 as a sentinel for "already unbalanced."
