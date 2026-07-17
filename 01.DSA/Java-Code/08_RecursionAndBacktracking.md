# Recursion & Backtracking

> Backtracking = DFS + pruning. The template: choose → explore → unchoose (backtrack).

---

## Table of Contents

1. [Subsets / Permutations / Combinations](#1-subsets--permutations--combinations)
2. [N-Queens](#2-n-queens)
3. [Sudoku Solver](#3-sudoku-solver)
4. [Word Search](#4-word-search)
5. [Palindrome Partitioning](#5-palindrome-partitioning)
6. [Rat in a Maze](#6-rat-in-a-maze)

---

## Backtracking Template

```java
void backtrack(state, choices) {
    if (baseCase) {
        addToResult(state);
        return;
    }
    for (choice : choices) {
        if (isValid(choice)) {
            makeChoice(choice);       // choose
            backtrack(state, next);   // explore
            undoChoice(choice);       // unchoose (backtrack)
        }
    }
}
```

---

## 1. Subsets / Permutations / Combinations

### Subsets (Power Set)

```java
// All subsets of nums (no duplicates in input)
// Time: O(n * 2^n), Space: O(n)
public List<List<Integer>> subsets(int[] nums) {
    List<List<Integer>> result = new ArrayList<>();
    backtrackSubsets(nums, 0, new ArrayList<>(), result);
    return result;
}

private void backtrackSubsets(int[] nums, int start, List<Integer> current, List<List<Integer>> result) {
    result.add(new ArrayList<>(current));   // add current subset (including empty)
    for (int i = start; i < nums.length; i++) {
        current.add(nums[i]);               // choose
        backtrackSubsets(nums, i + 1, current, result);  // explore (i+1 avoids reuse)
        current.remove(current.size() - 1); // unchoose
    }
}

// ── Subsets with Duplicates ──
public List<List<Integer>> subsetsWithDup(int[] nums) {
    Arrays.sort(nums);   // sort first to group duplicates
    List<List<Integer>> result = new ArrayList<>();
    backtrackSubsetsDup(nums, 0, new ArrayList<>(), result);
    return result;
}

private void backtrackSubsetsDup(int[] nums, int start, List<Integer> current, List<List<Integer>> result) {
    result.add(new ArrayList<>(current));
    for (int i = start; i < nums.length; i++) {
        if (i > start && nums[i] == nums[i-1]) continue;  // skip duplicate at same level
        current.add(nums[i]);
        backtrackSubsetsDup(nums, i + 1, current, result);
        current.remove(current.size() - 1);
    }
}
```

### Permutations

```java
// All permutations of nums (distinct)
// Time: O(n * n!), Space: O(n)
public List<List<Integer>> permute(int[] nums) {
    List<List<Integer>> result = new ArrayList<>();
    backtrackPermute(nums, new ArrayList<>(), new boolean[nums.length], result);
    return result;
}

private void backtrackPermute(int[] nums, List<Integer> current, boolean[] used, List<List<Integer>> result) {
    if (current.size() == nums.length) {
        result.add(new ArrayList<>(current));
        return;
    }
    for (int i = 0; i < nums.length; i++) {
        if (used[i]) continue;
        used[i] = true;
        current.add(nums[i]);
        backtrackPermute(nums, current, used, result);
        current.remove(current.size() - 1);
        used[i] = false;
    }
}
```

### Combinations

```java
// All combinations of k numbers from 1 to n
// Time: O(C(n,k) * k), Space: O(k)
public List<List<Integer>> combine(int n, int k) {
    List<List<Integer>> result = new ArrayList<>();
    backtrackCombine(n, k, 1, new ArrayList<>(), result);
    return result;
}

private void backtrackCombine(int n, int k, int start, List<Integer> current, List<List<Integer>> result) {
    if (current.size() == k) {
        result.add(new ArrayList<>(current));
        return;
    }
    // Pruning: remaining slots needed = k - current.size()
    // Only iterate if enough numbers remain: i <= n - (k - current.size()) + 1
    for (int i = start; i <= n - (k - current.size()) + 1; i++) {
        current.add(i);
        backtrackCombine(n, k, i + 1, current, result);
        current.remove(current.size() - 1);
    }
}
```

---

## 2. N-Queens

**Problem:** Place N queens on an N×N chessboard so no two queens attack each other.

```java
// Time: O(n!), Space: O(n)
public List<List<String>> solveNQueens(int n) {
    List<List<String>> result = new ArrayList<>();
    int[] queens = new int[n];    // queens[row] = col position of queen in that row
    Arrays.fill(queens, -1);

    Set<Integer> cols     = new HashSet<>();
    Set<Integer> diag1    = new HashSet<>();  // row - col (\ diagonals)
    Set<Integer> diag2    = new HashSet<>();  // row + col (/ diagonals)

    backtrackQueens(n, 0, queens, cols, diag1, diag2, result);
    return result;
}

private void backtrackQueens(int n, int row, int[] queens,
        Set<Integer> cols, Set<Integer> diag1, Set<Integer> diag2,
        List<List<String>> result) {

    if (row == n) {
        result.add(buildBoard(queens, n));
        return;
    }
    for (int col = 0; col < n; col++) {
        if (cols.contains(col) || diag1.contains(row - col) || diag2.contains(row + col))
            continue;  // conflict — skip

        queens[row] = col;
        cols.add(col); diag1.add(row - col); diag2.add(row + col);

        backtrackQueens(n, row + 1, queens, cols, diag1, diag2, result);

        queens[row] = -1;
        cols.remove(col); diag1.remove(row - col); diag2.remove(row + col);
    }
}

private List<String> buildBoard(int[] queens, int n) {
    List<String> board = new ArrayList<>();
    for (int row = 0; row < n; row++) {
        char[] line = new char[n];
        Arrays.fill(line, '.');
        line[queens[row]] = 'Q';
        board.add(new String(line));
    }
    return board;
}
```

> **Interview Q: How do you check if two queens are on the same diagonal?**  
> Two queens at `(r1,c1)` and `(r2,c2)` are on the same **`\` diagonal** if `r1 - c1 == r2 - c2`, and the same **`/` diagonal** if `r1 + c1 == r2 + c2`. Using Sets for these two diagonal values gives O(1) conflict checking.

---

## 3. Sudoku Solver

```java
public void solveSudoku(char[][] board) {
    solve(board);
}

private boolean solve(char[][] board) {
    for (int r = 0; r < 9; r++) {
        for (int c = 0; c < 9; c++) {
            if (board[r][c] != '.') continue;

            for (char num = '1'; num <= '9'; num++) {
                if (isValid(board, r, c, num)) {
                    board[r][c] = num;          // place

                    if (solve(board)) return true;  // recurse

                    board[r][c] = '.';          // backtrack
                }
            }
            return false;   // no valid number fits — dead end
        }
    }
    return true;   // all cells filled
}

private boolean isValid(char[][] board, int row, int col, char num) {
    for (int i = 0; i < 9; i++) {
        if (board[row][i] == num) return false;            // check row
        if (board[i][col] == num) return false;            // check col
        // check 3x3 box
        int boxRow = 3 * (row / 3) + i / 3;
        int boxCol = 3 * (col / 3) + i % 3;
        if (board[boxRow][boxCol] == num) return false;
    }
    return true;
}
```

---

## 4. Word Search

**Problem:** Given a 2D board of characters, check if a word exists as a connected path (no cell reused).

```java
// Time: O(m * n * 4^L) where L = word length, Space: O(L)
public boolean exist(char[][] board, String word) {
    int rows = board.length, cols = board[0].length;
    for (int r = 0; r < rows; r++) {
        for (int c = 0; c < cols; c++) {
            if (dfsWord(board, word, r, c, 0)) return true;
        }
    }
    return false;
}

private boolean dfsWord(char[][] board, String word, int r, int c, int idx) {
    if (idx == word.length()) return true;               // all chars matched
    if (r < 0 || c < 0 || r >= board.length || c >= board[0].length) return false;
    if (board[r][c] != word.charAt(idx)) return false;

    char temp = board[r][c];
    board[r][c] = '#';   // mark visited (in-place — no extra visited array)

    boolean found = dfsWord(board, word, r+1, c, idx+1) ||
                    dfsWord(board, word, r-1, c, idx+1) ||
                    dfsWord(board, word, r, c+1, idx+1) ||
                    dfsWord(board, word, r, c-1, idx+1);

    board[r][c] = temp;  // restore (backtrack)
    return found;
}
```

---

## 5. Palindrome Partitioning

**Problem:** Partition string `s` so every substring is a palindrome. Return all such partitions.

```java
// Time: O(n * 2^n), Space: O(n²) for palindrome cache
public List<List<String>> partition(String s) {
    List<List<String>> result = new ArrayList<>();
    // Precompute palindromes using DP
    boolean[][] isPalin = new boolean[s.length()][s.length()];
    for (int i = s.length() - 1; i >= 0; i--) {
        for (int j = i; j < s.length(); j++) {
            isPalin[i][j] = s.charAt(i) == s.charAt(j) &&
                            (j - i <= 2 || isPalin[i+1][j-1]);
        }
    }
    backtrackPalin(s, 0, isPalin, new ArrayList<>(), result);
    return result;
}

private void backtrackPalin(String s, int start, boolean[][] isPalin,
        List<String> current, List<List<String>> result) {
    if (start == s.length()) {
        result.add(new ArrayList<>(current));
        return;
    }
    for (int end = start; end < s.length(); end++) {
        if (isPalin[start][end]) {
            current.add(s.substring(start, end + 1));
            backtrackPalin(s, end + 1, isPalin, current, result);
            current.remove(current.size() - 1);
        }
    }
}
```

---

## 6. Rat in a Maze

**Problem:** A rat starts at `(0,0)` in an N×N grid (1=open, 0=blocked) and must reach `(N-1,N-1)`. Find all paths.

```java
public List<String> findPaths(int[][] maze) {
    List<String> result = new ArrayList<>();
    boolean[][] visited = new boolean[maze.length][maze.length];
    solveMaze(maze, 0, 0, "", visited, result);
    return result;
}

private void solveMaze(int[][] maze, int r, int c, String path,
        boolean[][] visited, List<String> result) {
    int n = maze.length;
    if (r == n-1 && c == n-1) {
        result.add(path);
        return;
    }
    int[][] dirs = {{1,0},{-1,0},{0,1},{0,-1}};
    char[] dirChar = {'D','U','R','L'};

    for (int d = 0; d < 4; d++) {
        int nr = r + dirs[d][0];
        int nc = c + dirs[d][1];
        if (nr >= 0 && nc >= 0 && nr < n && nc < n
                && maze[nr][nc] == 1 && !visited[nr][nc]) {
            visited[r][c] = true;
            solveMaze(maze, nr, nc, path + dirChar[d], visited, result);
            visited[r][c] = false;   // backtrack
        }
    }
}
```
