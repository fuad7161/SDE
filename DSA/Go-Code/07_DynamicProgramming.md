# Dynamic Programming

> DP = Recursion + Memoization, or equivalently, careful bottom-up table filling. Identify overlapping subproblems and optimal substructure.

---

## Table of Contents

1. [Fibonacci / Climbing Stairs](#1-fibonacci--climbing-stairs)
2. [0/1 Knapsack](#2-01-knapsack)
3. [Longest Common Subsequence (LCS)](#3-longest-common-subsequence-lcs)
4. [Longest Increasing Subsequence (LIS)](#4-longest-increasing-subsequence-lis)
5. [Edit Distance](#5-edit-distance)
6. [Coin Change](#6-coin-change)
7. [House Robber I / II / III](#7-house-robber-i--ii--iii)
8. [Partition Equal Subset Sum](#8-partition-equal-subset-sum)
9. [Grid DP — Unique Paths & Minimum Path Sum](#9-grid-dp--unique-paths--minimum-path-sum)

---

## 1. Fibonacci / Climbing Stairs

```go
// Climbing Stairs — Time: O(n), Space: O(1)
// Distinct ways to climb n stairs (1 or 2 steps at a time)
func climbStairs(n int) int {
    if n <= 2 { return n }
    prev1, prev2 := 1, 2
    for i := 3; i <= n; i++ {
        curr := prev1 + prev2
        prev1 = prev2
        prev2 = curr
    }
    return prev2
}
```

---

## 2. 0/1 Knapsack

```go
// Time: O(n * W), Space: O(W) — 1D DP
// Maximize value with items of given weights, capacity W
func knapsack(weights []int, values []int, W int) int {
    dp := make([]int, W+1)
    for i := 0; i < len(weights); i++ {
        // Traverse right-to-left to avoid using item twice
        for w := W; w >= weights[i]; w-- {
            if dp[w-weights[i]]+values[i] > dp[w] {
                dp[w] = dp[w-weights[i]] + values[i]
            }
        }
    }
    return dp[W]
}
```

> **Interview Q: Why iterate backwards for 0/1 knapsack?**  
> Right-to-left ensures each item is used at most once. If we go left-to-right, `dp[w-weights[i]]` may already include the current item, effectively using it multiple times (which is the unbounded knapsack variant).

---

## 3. Longest Common Subsequence (LCS)

```go
// Time: O(m*n), Space: O(m*n)
func longestCommonSubsequence(text1 string, text2 string) int {
    m, n := len(text1), len(text2)
    dp := make([][]int, m+1)
    for i := range dp { dp[i] = make([]int, n+1) }

    for i := 1; i <= m; i++ {
        for j := 1; j <= n; j++ {
            if text1[i-1] == text2[j-1] {
                dp[i][j] = dp[i-1][j-1] + 1
            } else {
                if dp[i-1][j] > dp[i][j-1] {
                    dp[i][j] = dp[i-1][j]
                } else {
                    dp[i][j] = dp[i][j-1]
                }
            }
        }
    }
    return dp[m][n]
}
```

---

## 4. Longest Increasing Subsequence (LIS)

```go
// DP — Time: O(n²), Space: O(n)
func lisDP(nums []int) int {
    n := len(nums)
    dp := make([]int, n)
    for i := range dp { dp[i] = 1 }
    result := 1
    for i := 1; i < n; i++ {
        for j := 0; j < i; j++ {
            if nums[j] < nums[i] && dp[j]+1 > dp[i] {
                dp[i] = dp[j] + 1
            }
        }
        if dp[i] > result { result = dp[i] }
    }
    return result
}

// Patience Sort — Time: O(n log n), Space: O(n)
import "sort"

func lengthOfLIS(nums []int) int {
    tails := []int{}
    for _, num := range nums {
        // Find first tail >= num
        pos := sort.SearchInts(tails, num)
        if pos == len(tails) {
            tails = append(tails, num)
        } else {
            tails[pos] = num
        }
    }
    return len(tails)
}
```

---

## 5. Edit Distance

```go
// Time: O(m*n), Space: O(m*n)
func minDistance(word1 string, word2 string) int {
    m, n := len(word1), len(word2)
    dp := make([][]int, m+1)
    for i := range dp {
        dp[i] = make([]int, n+1)
        dp[i][0] = i // delete all chars of word1
    }
    for j := 0; j <= n; j++ { dp[0][j] = j } // insert all chars

    for i := 1; i <= m; i++ {
        for j := 1; j <= n; j++ {
            if word1[i-1] == word2[j-1] {
                dp[i][j] = dp[i-1][j-1]
            } else {
                dp[i][j] = 1 + min3(dp[i-1][j], dp[i][j-1], dp[i-1][j-1])
            }
        }
    }
    return dp[m][n]
}

func min3(a, b, c int) int {
    if a < b {
        if a < c { return a }
        return c
    }
    if b < c { return b }
    return c
}
```

---

## 6. Coin Change

```go
// Minimum coins to make amount — Time: O(amount * n), Space: O(amount)
import "math"

func coinChange(coins []int, amount int) int {
    dp := make([]int, amount+1)
    for i := 1; i <= amount; i++ { dp[i] = math.MaxInt }

    for i := 1; i <= amount; i++ {
        for _, coin := range coins {
            if coin <= i && dp[i-coin] != math.MaxInt {
                if dp[i-coin]+1 < dp[i] {
                    dp[i] = dp[i-coin] + 1
                }
            }
        }
    }
    if dp[amount] == math.MaxInt { return -1 }
    return dp[amount]
}

// Coin Change II — Number of combinations (unbounded knapsack)
func change(amount int, coins []int) int {
    dp := make([]int, amount+1)
    dp[0] = 1
    for _, coin := range coins {
        for i := coin; i <= amount; i++ {
            dp[i] += dp[i-coin]
        }
    }
    return dp[amount]
}
```

> **Interview Q: Coin Change vs Coin Change II — what's the key difference?**  
> Change I asks for minimum *count* → minimize. Change II asks for number of *ways* → count. Also note: in Change II, iterating coin-by-coin (outer loop) ensures we don't count permutations as different combinations.

---

## 7. House Robber I / II / III

```go
// House Robber I — no two adjacent houses
// Time: O(n), Space: O(1)
func rob(nums []int) int {
    prev2, prev1 := 0, 0
    for _, n := range nums {
        curr := prev2 + n
        if prev1 > curr { curr = prev1 }
        prev2 = prev1
        prev1 = curr
    }
    return prev1
}

// House Robber II — circular array
func robII(nums []int) int {
    n := len(nums)
    if n == 1 { return nums[0] }
    robRange := func(start, end int) int {
        prev2, prev1 := 0, 0
        for i := start; i <= end; i++ {
            curr := prev2 + nums[i]
            if prev1 > curr { curr = prev1 }
            prev2 = prev1
            prev1 = curr
        }
        return prev1
    }
    a := robRange(0, n-2)
    b := robRange(1, n-1)
    if a > b { return a }
    return b
}

// House Robber III — binary tree
func robIII(root *TreeNode) int {
    var dfs func(*TreeNode) (int, int) // (rob root, skip root)
    dfs = func(node *TreeNode) (int, int) {
        if node == nil { return 0, 0 }
        lRob, lSkip := dfs(node.Left)
        rRob, rSkip := dfs(node.Right)
        robRoot := node.Val + lSkip + rSkip
        skipRoot := max(lRob, lSkip) + max(rRob, rSkip)
        return robRoot, skipRoot
    }
    a, b := dfs(root)
    if a > b { return a }
    return b
}

func max(a, b int) int { if a > b { return a }; return b }
```

---

## 8. Partition Equal Subset Sum

```go
// Can we split nums into two subsets with equal sum?
// Time: O(n * sum), Space: O(sum)
func canPartition(nums []int) bool {
    total := 0
    for _, n := range nums { total += n }
    if total%2 != 0 { return false }
    target := total / 2

    dp := make([]bool, target+1)
    dp[0] = true

    for _, num := range nums {
        for j := target; j >= num; j-- {
            if dp[j-num] { dp[j] = true }
        }
    }
    return dp[target]
}
```

---

## 9. Grid DP — Unique Paths & Minimum Path Sum

```go
// Unique Paths — Time: O(m*n), Space: O(n)
func uniquePaths(m int, n int) int {
    dp := make([]int, n)
    for j := range dp { dp[j] = 1 }
    for i := 1; i < m; i++ {
        for j := 1; j < n; j++ {
            dp[j] += dp[j-1]
        }
    }
    return dp[n-1]
}

// Minimum Path Sum
func minPathSum(grid [][]int) int {
    m, n := len(grid), len(grid[0])
    dp := make([]int, n)
    dp[0] = grid[0][0]
    for j := 1; j < n; j++ { dp[j] = dp[j-1] + grid[0][j] }

    for i := 1; i < m; i++ {
        dp[0] += grid[i][0]
        for j := 1; j < n; j++ {
            if dp[j] < dp[j-1] {
                dp[j] = dp[j] + grid[i][j]
            } else {
                dp[j] = dp[j-1] + grid[i][j]
            }
        }
    }
    return dp[n-1]
}
```
