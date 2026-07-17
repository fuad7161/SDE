# Dynamic Programming

> Heavy in FAANG interviews. Always think: define state, recurrence relation, and base cases. Start with brute force → memoization → tabulation.

---

## Table of Contents

1. [Fibonacci / Climbing Stairs](#1-fibonacci--climbing-stairs)
2. [0/1 Knapsack](#2-01-knapsack)
3. [Longest Common Subsequence (LCS)](#3-longest-common-subsequence-lcs)
4. [Longest Increasing Subsequence (LIS)](#4-longest-increasing-subsequence-lis)
5. [Edit Distance](#5-edit-distance)
6. [Coin Change (I & II)](#6-coin-change-i--ii)
7. [House Robber (I, II, III)](#7-house-robber-i-ii-iii)
8. [Partition Equal Subset Sum](#8-partition-equal-subset-sum)
9. [Grid DP](#9-grid-dp)

---

## DP Approach

```
1. Identify if it's a DP problem (optimal substructure + overlapping subproblems)
2. Define state: dp[i] means "..."
3. Write recurrence: dp[i] = f(dp[i-1], dp[i-2], ...)
4. Set base cases
5. Determine order of computation
6. Optimize space if possible
```

---

## 1. Fibonacci / Climbing Stairs

```java
// ── Fibonacci ──
// dp[i] = dp[i-1] + dp[i-2]

// Naive recursion — O(2^n) time (exponential!)
int fib(int n) {
    if (n <= 1) return n;
    return fib(n-1) + fib(n-2);
}

// Memoization (top-down DP) — O(n) time, O(n) space
int fib(int n, int[] memo) {
    if (n <= 1) return n;
    if (memo[n] != 0) return memo[n];
    return memo[n] = fib(n-1, memo) + fib(n-2, memo);
}

// Tabulation (bottom-up DP) — O(n) time, O(n) space
int fib(int n) {
    if (n <= 1) return n;
    int[] dp = new int[n+1];
    dp[0] = 0; dp[1] = 1;
    for (int i = 2; i <= n; i++) dp[i] = dp[i-1] + dp[i-2];
    return dp[n];
}

// Space-optimized — O(n) time, O(1) space
int fib(int n) {
    if (n <= 1) return n;
    int prev2 = 0, prev1 = 1;
    for (int i = 2; i <= n; i++) {
        int curr = prev1 + prev2;
        prev2 = prev1;
        prev1 = curr;
    }
    return prev1;
}

// ── Climbing Stairs ──
// n steps, can climb 1 or 2 steps at a time. How many distinct ways?
// Exactly Fibonacci! dp[i] = dp[i-1] + dp[i-2]
public int climbStairs(int n) {
    if (n <= 2) return n;
    int one = 1, two = 2;
    for (int i = 3; i <= n; i++) {
        int curr = one + two;
        one = two;
        two = curr;
    }
    return two;
}
```

---

## 2. 0/1 Knapsack

**Problem:** Given items with weights and values, fill a knapsack of capacity W to maximize value. Each item used at most once.

```java
// Time: O(n*W), Space: O(n*W)
public int knapsack(int[] weights, int[] values, int W) {
    int n = weights.length;
    int[][] dp = new int[n + 1][W + 1];
    // dp[i][w] = max value using first i items with capacity w

    for (int i = 1; i <= n; i++) {
        for (int w = 0; w <= W; w++) {
            dp[i][w] = dp[i-1][w];  // don't take item i
            if (weights[i-1] <= w) {
                dp[i][w] = Math.max(dp[i][w],
                    dp[i-1][w - weights[i-1]] + values[i-1]);  // take item i
            }
        }
    }
    return dp[n][W];
}

// ── Space-optimized: O(W) space ──
// Traverse W from right to left to avoid using same item twice
public int knapsackOptimized(int[] weights, int[] values, int W) {
    int[] dp = new int[W + 1];
    for (int i = 0; i < weights.length; i++) {
        for (int w = W; w >= weights[i]; w--) {
            dp[w] = Math.max(dp[w], dp[w - weights[i]] + values[i]);
        }
    }
    return dp[W];
}
```

> **Interview Q: Why iterate W from right to left in the space-optimized version?**  
> In the 0/1 knapsack, each item can be used at most once. Iterating right to left ensures that when we update `dp[w]`, `dp[w - weight[i]]` still refers to the state **before** including item `i` (i.e., the previous row). Left-to-right would use `dp[w - weight[i]]` updated in the current row, allowing the same item to be counted multiple times (like unbounded knapsack).

---

## 3. Longest Common Subsequence (LCS)

**Problem:** Length of the longest subsequence present in both strings (not necessarily contiguous).

```java
// Time: O(m*n), Space: O(m*n)
// dp[i][j] = LCS of s1[0..i-1] and s2[0..j-1]
public int longestCommonSubsequence(String s1, String s2) {
    int m = s1.length(), n = s2.length();
    int[][] dp = new int[m + 1][n + 1];

    for (int i = 1; i <= m; i++) {
        for (int j = 1; j <= n; j++) {
            if (s1.charAt(i-1) == s2.charAt(j-1)) {
                dp[i][j] = dp[i-1][j-1] + 1;           // chars match — extend LCS
            } else {
                dp[i][j] = Math.max(dp[i-1][j], dp[i][j-1]);  // skip one char
            }
        }
    }
    return dp[m][n];
}

// ── Reconstruct the LCS string ──
public String lcString(String s1, String s2) {
    int m = s1.length(), n = s2.length();
    int[][] dp = new int[m + 1][n + 1];
    for (int i = 1; i <= m; i++)
        for (int j = 1; j <= n; j++)
            dp[i][j] = s1.charAt(i-1) == s2.charAt(j-1)
                ? dp[i-1][j-1] + 1
                : Math.max(dp[i-1][j], dp[i][j-1]);

    StringBuilder sb = new StringBuilder();
    int i = m, j = n;
    while (i > 0 && j > 0) {
        if (s1.charAt(i-1) == s2.charAt(j-1)) { sb.append(s1.charAt(i-1)); i--; j--; }
        else if (dp[i-1][j] > dp[i][j-1])       i--;
        else                                      j--;
    }
    return sb.reverse().toString();
}
```

---

## 4. Longest Increasing Subsequence (LIS)

```java
// ── O(n²) DP ──
// dp[i] = length of LIS ending at index i
public int lengthOfLIS(int[] nums) {
    int n = nums.length, maxLen = 1;
    int[] dp = new int[n];
    Arrays.fill(dp, 1);

    for (int i = 1; i < n; i++) {
        for (int j = 0; j < i; j++) {
            if (nums[j] < nums[i]) {
                dp[i] = Math.max(dp[i], dp[j] + 1);
            }
        }
        maxLen = Math.max(maxLen, dp[i]);
    }
    return maxLen;
}

// ── O(n log n) — patience sorting with binary search ──
public int lengthOfLISFast(int[] nums) {
    List<Integer> tails = new ArrayList<>();  // tails[i] = smallest tail for LIS of length i+1

    for (int num : nums) {
        int lo = 0, hi = tails.size();
        while (lo < hi) {                     // binary search for first tail >= num
            int mid = lo + (hi - lo) / 2;
            if (tails.get(mid) < num) lo = mid + 1;
            else hi = mid;
        }
        if (lo == tails.size()) tails.add(num);  // extend LIS
        else tails.set(lo, num);                 // replace to keep smallest tail
    }
    return tails.size();
}
// nums = [10,9,2,5,3,7,101,18] → LIS length = 4 ([2,3,7,18] or [2,5,7,101])
```

---

## 5. Edit Distance

**Problem:** Minimum number of operations (insert, delete, replace) to convert `word1` to `word2`.

```java
// Time: O(m*n), Space: O(m*n)
// dp[i][j] = min ops to convert word1[0..i-1] to word2[0..j-1]
public int minDistance(String word1, String word2) {
    int m = word1.length(), n = word2.length();
    int[][] dp = new int[m + 1][n + 1];

    // Base cases: converting to/from empty string
    for (int i = 0; i <= m; i++) dp[i][0] = i;  // delete all chars
    for (int j = 0; j <= n; j++) dp[0][j] = j;  // insert all chars

    for (int i = 1; i <= m; i++) {
        for (int j = 1; j <= n; j++) {
            if (word1.charAt(i-1) == word2.charAt(j-1)) {
                dp[i][j] = dp[i-1][j-1];            // no operation needed
            } else {
                dp[i][j] = 1 + Math.min(dp[i-1][j-1],   // replace
                               Math.min(dp[i-1][j],       // delete from word1
                                        dp[i][j-1]));     // insert into word1
            }
        }
    }
    return dp[m][n];
}
// "horse" → "ros": dp[5][3] = 3
// replace h→r, delete r, delete e
```

---

## 6. Coin Change (I & II)

### Coin Change I — Minimum coins to make amount

```java
// Unbounded knapsack: each coin can be used multiple times
// Time: O(amount * coins), Space: O(amount)
public int coinChange(int[] coins, int amount) {
    int[] dp = new int[amount + 1];
    Arrays.fill(dp, amount + 1);  // initialize to "infinity"
    dp[0] = 0;

    for (int i = 1; i <= amount; i++) {
        for (int coin : coins) {
            if (coin <= i) {
                dp[i] = Math.min(dp[i], dp[i - coin] + 1);
            }
        }
    }
    return dp[amount] > amount ? -1 : dp[amount];
}
```

### Coin Change II — Number of combinations to make amount

```java
// Time: O(amount * coins), Space: O(amount)
// ORDER: iterate coins in outer loop to count COMBINATIONS (not permutations)
public int change(int amount, int[] coins) {
    int[] dp = new int[amount + 1];
    dp[0] = 1;  // one way to make amount 0 (use no coins)

    for (int coin : coins) {
        for (int i = coin; i <= amount; i++) {
            dp[i] += dp[i - coin];
        }
    }
    return dp[amount];
}
// If we swapped loops (amount outer, coins inner), we'd count permutations
```

> **Interview Q: Why does loop order matter in Coin Change II?**  
> Iterating coins in the outer loop ensures each coin is "considered" for extension of previously computed sub-amounts, preventing duplicate counting of the same combination in different order. If amount is outer, `(1,2)` and `(2,1)` would be counted as different — giving permutations instead of combinations.

---

## 7. House Robber (I, II, III)

### House Robber I — Linear street

```java
// Cannot rob adjacent houses
// dp[i] = max money robbing houses 0..i
public int rob(int[] nums) {
    if (nums.length == 1) return nums[0];
    int prev2 = 0, prev1 = 0;
    for (int num : nums) {
        int curr = Math.max(prev1, prev2 + num);
        prev2 = prev1;
        prev1 = curr;
    }
    return prev1;
}
```

### House Robber II — Circular street (first and last houses are adjacent)

```java
// Rob either houses[0..n-2] or houses[1..n-1] (can't rob both first and last)
public int robCircular(int[] nums) {
    if (nums.length == 1) return nums[0];
    return Math.max(robRange(nums, 0, nums.length - 2),
                    robRange(nums, 1, nums.length - 1));
}

private int robRange(int[] nums, int lo, int hi) {
    int prev2 = 0, prev1 = 0;
    for (int i = lo; i <= hi; i++) {
        int curr = Math.max(prev1, prev2 + nums[i]);
        prev2 = prev1;
        prev1 = curr;
    }
    return prev1;
}
```

### House Robber III — Binary tree (no parent-child rob at same time)

```java
// At each node, return [maxRobbedWithNode, maxRobbedWithoutNode]
public int robTree(TreeNode root) {
    int[] res = dfs(root);
    return Math.max(res[0], res[1]);
}

private int[] dfs(TreeNode node) {
    if (node == null) return new int[]{0, 0};
    int[] left  = dfs(node.left);
    int[] right = dfs(node.right);

    int rob    = node.val + left[1] + right[1];       // rob this node (skip children)
    int noRob  = Math.max(left[0], left[1])
               + Math.max(right[0], right[1]);         // skip this node (children optional)
    return new int[]{rob, noRob};
}
```

---

## 8. Partition Equal Subset Sum

**Problem:** Can the array be partitioned into two subsets with equal sum?  
**Equivalent to:** Does any subset sum to `totalSum / 2`?

```java
// Time: O(n * sum/2), Space: O(sum/2)
public boolean canPartition(int[] nums) {
    int total = Arrays.stream(nums).sum();
    if (total % 2 != 0) return false;
    int target = total / 2;

    boolean[] dp = new boolean[target + 1];
    dp[0] = true;  // sum 0 is always achievable (empty subset)

    for (int num : nums) {
        // Traverse right to left: 0/1 knapsack (each number used once)
        for (int j = target; j >= num; j--) {
            dp[j] = dp[j] || dp[j - num];
        }
    }
    return dp[target];
}
```

---

## 9. Grid DP

### Unique Paths

```java
// Robot in m×n grid, move only right or down. Count paths from top-left to bottom-right.
// dp[i][j] = paths to reach (i,j)
public int uniquePaths(int m, int n) {
    int[][] dp = new int[m][n];
    for (int i = 0; i < m; i++) dp[i][0] = 1;   // first column: only one way (go down)
    for (int j = 0; j < n; j++) dp[0][j] = 1;   // first row: only one way (go right)

    for (int i = 1; i < m; i++)
        for (int j = 1; j < n; j++)
            dp[i][j] = dp[i-1][j] + dp[i][j-1];

    return dp[m-1][n-1];
}
```

### Minimum Path Sum

```java
// Find path from top-left to bottom-right minimizing sum of numbers along the path
public int minPathSum(int[][] grid) {
    int m = grid.length, n = grid[0].length;
    int[][] dp = new int[m][n];
    dp[0][0] = grid[0][0];

    for (int i = 1; i < m; i++) dp[i][0] = dp[i-1][0] + grid[i][0];
    for (int j = 1; j < n; j++) dp[0][j] = dp[0][j-1] + grid[0][j];

    for (int i = 1; i < m; i++)
        for (int j = 1; j < n; j++)
            dp[i][j] = Math.min(dp[i-1][j], dp[i][j-1]) + grid[i][j];

    return dp[m-1][n-1];
}
```
