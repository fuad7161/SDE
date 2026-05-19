# Recursion & Backtracking

> Backtracking = DFS on a decision tree with pruning. Always think: **choose → explore → un-choose**.

---

## Table of Contents

1. [Backtracking Template](#1-backtracking-template)
2. [Subsets](#2-subsets)
3. [Permutations](#3-permutations)
4. [Combinations](#4-combinations)
5. [N-Queens](#5-n-queens)
6. [Sudoku Solver](#6-sudoku-solver)
7. [Word Search](#7-word-search)
8. [Palindrome Partitioning](#8-palindrome-partitioning)
9. [Rat in a Maze](#9-rat-in-a-maze)

---

## 1. Backtracking Template

```go
// Canonical backtracking template in Go
func backtrack(result *[][]int, current []int, /* other state */ candidates []int, start int) {
    // Base case — record result
    if /* solution found */ len(current) == someTarget {
        tmp := make([]int, len(current))
        copy(tmp, current)           // ALWAYS copy before appending to result
        *result = append(*result, tmp)
        return
    }

    for i := start; i < len(candidates); i++ {
        // Choose
        current = append(current, candidates[i])

        // Explore
        backtrack(result, current, candidates, i+1) // i+1 = no reuse; i = reuse allowed

        // Un-choose (backtrack)
        current = current[:len(current)-1]
    }
}
```

> **Critical in Go**: Slices share underlying arrays. Always `copy` before saving to `result`, or use `append(*result, current...)`.

---

## 2. Subsets

```go
// Time: O(n * 2^n), Space: O(n)
func subsets(nums []int) [][]int {
    result := [][]int{}
    var backtrack func(start int, current []int)
    backtrack = func(start int, current []int) {
        tmp := make([]int, len(current))
        copy(tmp, current)
        result = append(result, tmp)

        for i := start; i < len(nums); i++ {
            current = append(current, nums[i])
            backtrack(i+1, current)
            current = current[:len(current)-1]
        }
    }
    backtrack(0, []int{})
    return result
}

// Subsets II — with duplicates, no duplicate subsets
import "sort"

func subsetsWithDup(nums []int) [][]int {
    sort.Ints(nums) // sort to group duplicates
    result := [][]int{}
    var backtrack func(start int, current []int)
    backtrack = func(start int, current []int) {
        tmp := make([]int, len(current))
        copy(tmp, current)
        result = append(result, tmp)

        for i := start; i < len(nums); i++ {
            if i > start && nums[i] == nums[i-1] { continue } // skip dup
            current = append(current, nums[i])
            backtrack(i+1, current)
            current = current[:len(current)-1]
        }
    }
    backtrack(0, []int{})
    return result
}
```

---

## 3. Permutations

```go
// Permutations I — all permutations (distinct elements)
// Time: O(n! * n)
func permute(nums []int) [][]int {
    result := [][]int{}
    var backtrack func(current []int, used []bool)
    backtrack = func(current []int, used []bool) {
        if len(current) == len(nums) {
            tmp := make([]int, len(current))
            copy(tmp, current)
            result = append(result, tmp)
            return
        }
        for i, num := range nums {
            if used[i] { continue }
            used[i] = true
            current = append(current, num)
            backtrack(current, used)
            current = current[:len(current)-1]
            used[i] = false
        }
    }
    backtrack([]int{}, make([]bool, len(nums)))
    return result
}

// Permutations II — with duplicates
import "sort"

func permuteUnique(nums []int) [][]int {
    sort.Ints(nums)
    result := [][]int{}
    used := make([]bool, len(nums))

    var backtrack func(current []int)
    backtrack = func(current []int) {
        if len(current) == len(nums) {
            tmp := make([]int, len(current))
            copy(tmp, current)
            result = append(result, tmp)
            return
        }
        for i, num := range nums {
            if used[i] { continue }
            if i > 0 && nums[i] == nums[i-1] && !used[i-1] { continue }
            used[i] = true
            current = append(current, num)
            backtrack(current)
            current = current[:len(current)-1]
            used[i] = false
        }
    }
    backtrack([]int{})
    return result
}
```

---

## 4. Combinations

```go
// All combinations of k numbers from 1..n
// Time: O(C(n,k) * k)
func combine(n int, k int) [][]int {
    result := [][]int{}
    var backtrack func(start int, current []int)
    backtrack = func(start int, current []int) {
        if len(current) == k {
            tmp := make([]int, k)
            copy(tmp, current)
            result = append(result, tmp)
            return
        }
        // Pruning: need (k-len(current)) more — don't start too late
        for i := start; i <= n-(k-len(current))+1; i++ {
            current = append(current, i)
            backtrack(i+1, current)
            current = current[:len(current)-1]
        }
    }
    backtrack(1, []int{})
    return result
}

// Combination Sum — numbers may reuse, no duplicates in input
func combinationSum(candidates []int, target int) [][]int {
    result := [][]int{}
    var backtrack func(start, remain int, current []int)
    backtrack = func(start, remain int, current []int) {
        if remain == 0 {
            tmp := make([]int, len(current))
            copy(tmp, current)
            result = append(result, tmp)
            return
        }
        for i := start; i < len(candidates); i++ {
            if candidates[i] > remain { continue } // pruning (if sorted)
            current = append(current, candidates[i])
            backtrack(i, remain-candidates[i], current) // i = reuse allowed
            current = current[:len(current)-1]
        }
    }
    backtrack(0, target, []int{})
    return result
}
```

---

## 5. N-Queens

```go
// Time: O(n!), Space: O(n)
func solveNQueens(n int) [][]string {
    result := [][]string{}
    board := make([]int, n) // board[row] = column of queen
    for i := range board { board[i] = -1 }

    cols := make([]bool, n)
    diag1 := make([]bool, 2*n-1) // row - col + n - 1
    diag2 := make([]bool, 2*n-1) // row + col

    var backtrack func(row int)
    backtrack = func(row int) {
        if row == n {
            snapshot := make([]string, n)
            for r := 0; r < n; r++ {
                row2 := make([]byte, n)
                for c := 0; c < n; c++ {
                    if board[r] == c {
                        row2[c] = 'Q'
                    } else {
                        row2[c] = '.'
                    }
                }
                snapshot[r] = string(row2)
            }
            result = append(result, snapshot)
            return
        }
        for col := 0; col < n; col++ {
            d1, d2 := row-col+n-1, row+col
            if cols[col] || diag1[d1] || diag2[d2] { continue }
            board[row] = col
            cols[col], diag1[d1], diag2[d2] = true, true, true
            backtrack(row + 1)
            board[row] = -1
            cols[col], diag1[d1], diag2[d2] = false, false, false
        }
    }
    backtrack(0)
    return result
}
```

---

## 6. Sudoku Solver

```go
// Time: O(9^81) worst case — pruning makes it fast in practice
func solveSudoku(board [][]byte) {
    var solve func() bool
    solve = func() bool {
        for r := 0; r < 9; r++ {
            for c := 0; c < 9; c++ {
                if board[r][c] != '.' { continue }
                for d := byte('1'); d <= '9'; d++ {
                    if isValid(board, r, c, d) {
                        board[r][c] = d
                        if solve() { return true }
                        board[r][c] = '.'
                    }
                }
                return false // no valid digit found
            }
        }
        return true // all cells filled
    }
    solve()
}

func isValid(board [][]byte, row, col int, d byte) bool {
    box := (row/3)*3 + col/3
    for i := 0; i < 9; i++ {
        if board[row][i] == d { return false }
        if board[i][col] == d { return false }
        if board[(box/3)*3+i/3][(box%3)*3+i%3] == d { return false }
    }
    return true
}
```

---

## 7. Word Search

```go
// Time: O(m * n * 4^L), Space: O(L) where L = word length
func exist(board [][]byte, word string) bool {
    rows, cols := len(board), len(board[0])

    var dfs func(r, c, idx int) bool
    dfs = func(r, c, idx int) bool {
        if idx == len(word) { return true }
        if r < 0 || r >= rows || c < 0 || c >= cols { return false }
        if board[r][c] != word[idx] { return false }

        tmp := board[r][c]
        board[r][c] = '#' // mark visited

        found := dfs(r+1, c, idx+1) || dfs(r-1, c, idx+1) ||
                 dfs(r, c+1, idx+1) || dfs(r, c-1, idx+1)

        board[r][c] = tmp // restore
        return found
    }

    for r := 0; r < rows; r++ {
        for c := 0; c < cols; c++ {
            if dfs(r, c, 0) { return true }
        }
    }
    return false
}
```

---

## 8. Palindrome Partitioning

```go
// Time: O(n * 2^n), Space: O(n)
func partition(s string) [][]string {
    result := [][]string{}

    isPalin := func(str string, l, r int) bool {
        for l < r {
            if str[l] != str[r] { return false }
            l++; r--
        }
        return true
    }

    var backtrack func(start int, current []string)
    backtrack = func(start int, current []string) {
        if start == len(s) {
            tmp := make([]string, len(current))
            copy(tmp, current)
            result = append(result, tmp)
            return
        }
        for end := start; end < len(s); end++ {
            if isPalin(s, start, end) {
                current = append(current, s[start:end+1])
                backtrack(end+1, current)
                current = current[:len(current)-1]
            }
        }
    }
    backtrack(0, []string{})
    return result
}
```

---

## 9. Rat in a Maze

```go
// Find all paths from (0,0) to (n-1,n-1) in a binary maze
// Time: O(4^(n²)), Space: O(n²)
func findPaths(maze [][]int) []string {
    n := len(maze)
    result := []string{}
    if maze[0][0] == 0 { return result }

    visited := make([][]bool, n)
    for i := range visited { visited[i] = make([]bool, n) }

    dirChar := []byte{'D', 'L', 'R', 'U'}
    dr := []int{1, 0, 0, -1}
    dc := []int{0, -1, 1, 0}

    var dfs func(r, c int, path []byte)
    dfs = func(r, c int, path []byte) {
        if r == n-1 && c == n-1 {
            result = append(result, string(path))
            return
        }
        for i := 0; i < 4; i++ {
            nr, nc := r+dr[i], c+dc[i]
            if nr >= 0 && nr < n && nc >= 0 && nc < n &&
               maze[nr][nc] == 1 && !visited[nr][nc] {
                visited[nr][nc] = true
                path = append(path, dirChar[i])
                dfs(nr, nc, path)
                path = path[:len(path)-1]
                visited[nr][nc] = false
            }
        }
    }

    visited[0][0] = true
    dfs(0, 0, []byte{})
    return result
}
```
