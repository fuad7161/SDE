# Tries

> A Trie (prefix tree) provides O(L) insert, search, and prefix-check, where L is word length. Essential for autocomplete, spell-check, and IP routing.

---

## Table of Contents

1. [Trie Node in Go](#1-trie-node-in-go)
2. [Implement Trie (Insert, Search, StartsWith, Delete)](#2-implement-trie-insert-search-startswith-delete)
3. [Word Search II](#3-word-search-ii)
4. [Auto-Complete System](#4-auto-complete-system)

---

## 1. Trie Node in Go

```go
// Array-based children (fixed 26 lowercase letters)
type TrieNode struct {
    children [26]*TrieNode
    isEnd    bool
}

// Map-based children (flexible — supports any character set)
type TrieNodeMap struct {
    children map[rune]*TrieNodeMap
    isEnd    bool
}

func newTrieNodeMap() *TrieNodeMap {
    return &TrieNodeMap{children: make(map[rune]*TrieNodeMap)}
}
```

---

## 2. Implement Trie (Insert, Search, StartsWith, Delete)

```go
// Time: O(L) per operation, Space: O(total chars)
type Trie struct {
    root *TrieNode
}

func TrieConstructor() Trie {
    return Trie{root: &TrieNode{}}
}

func (t *Trie) Insert(word string) {
    node := t.root
    for _, c := range word {
        idx := c - 'a'
        if node.children[idx] == nil {
            node.children[idx] = &TrieNode{}
        }
        node = node.children[idx]
    }
    node.isEnd = true
}

func (t *Trie) Search(word string) bool {
    node := t.root
    for _, c := range word {
        idx := c - 'a'
        if node.children[idx] == nil {
            return false
        }
        node = node.children[idx]
    }
    return node.isEnd
}

func (t *Trie) StartsWith(prefix string) bool {
    node := t.root
    for _, c := range prefix {
        idx := c - 'a'
        if node.children[idx] == nil {
            return false
        }
        node = node.children[idx]
    }
    return true
}

// Delete — mark isEnd=false and prune empty nodes bottom-up
func (t *Trie) Delete(word string) {
    var del func(node *TrieNode, depth int) bool
    del = func(node *TrieNode, depth int) bool {
        if depth == len(word) {
            if !node.isEnd { return false } // not found
            node.isEnd = false
            // can delete this node if it has no children
            return t.isEmpty(node)
        }
        idx := word[depth] - 'a'
        if node.children[idx] == nil { return false }
        shouldDelete := del(node.children[idx], depth+1)
        if shouldDelete {
            node.children[idx] = nil
            return !node.isEnd && t.isEmpty(node)
        }
        return false
    }
    del(t.root, 0)
}

func (t *Trie) isEmpty(node *TrieNode) bool {
    for _, child := range node.children {
        if child != nil { return false }
    }
    return true
}
```

> **Interview Q: When should you use map-based vs array-based children?**  
> Array `[26]*TrieNode` is faster (O(1) index) and simpler for ASCII lowercase-only problems. Map-based is better when the character set is large or variable (Unicode, mixed case, digits).

---

## 3. Word Search II

```go
// Find all words from a dictionary that appear in a board
// Time: O(M * N * 4^L) with trie pruning, Space: O(total dict chars)
func findWords(board [][]byte, words []string) []string {
    // Build trie from all words
    type Node struct {
        children [26]*Node
        word     string // non-empty means a word ends here
    }
    root := &Node{}
    for _, w := range words {
        node := root
        for _, c := range w {
            idx := c - 'a'
            if node.children[idx] == nil {
                node.children[idx] = &Node{}
            }
            node = node.children[idx]
        }
        node.word = w
    }

    rows, cols := len(board), len(board[0])
    result := []string{}

    var dfs func(node *Node, r, c int)
    dfs = func(node *Node, r, c int) {
        if r < 0 || r >= rows || c < 0 || c >= cols { return }
        ch := board[r][c]
        if ch == '#' { return } // visited
        idx := ch - 'a'
        next := node.children[idx]
        if next == nil { return } // no matching prefix

        if next.word != "" {
            result = append(result, next.word)
            next.word = "" // deduplicate
        }

        board[r][c] = '#'
        dfs(next, r+1, c)
        dfs(next, r-1, c)
        dfs(next, r, c+1)
        dfs(next, r, c-1)
        board[r][c] = ch

        // Prune dead branches (optimization)
        if t.hasNoChildren(next) {
            node.children[idx] = nil
        }
    }

    for r := 0; r < rows; r++ {
        for c := 0; c < cols; c++ {
            dfs(root, r, c)
        }
    }
    return result
}

// ── Self-contained version without pruning helper ──
func findWordsSafe(board [][]byte, words []string) []string {
    type Node struct {
        children [26]*Node
        word     string
    }
    root := &Node{}
    for _, w := range words {
        node := root
        for i := 0; i < len(w); i++ {
            idx := w[i] - 'a'
            if node.children[idx] == nil {
                node.children[idx] = &Node{}
            }
            node = node.children[idx]
        }
        node.word = w
    }

    rows, cols := len(board), len(board[0])
    found := []string{}

    var dfs func(node *Node, r, c int)
    dfs = func(node *Node, r, c int) {
        if r < 0 || r >= rows || c < 0 || c >= cols || board[r][c] == '#' {
            return
        }
        next := node.children[board[r][c]-'a']
        if next == nil { return }
        if next.word != "" {
            found = append(found, next.word)
            next.word = ""
        }
        tmp := board[r][c]
        board[r][c] = '#'
        dfs(next, r+1, c); dfs(next, r-1, c)
        dfs(next, r, c+1); dfs(next, r, c-1)
        board[r][c] = tmp
    }

    for r := 0; r < rows; r++ {
        for c := 0; c < cols; c++ {
            dfs(root, r, c)
        }
    }
    return found
}
```

---

## 4. Auto-Complete System

```go
// Design a search autocomplete — returns top 3 historical queries by frequency
type AutocompleteNode struct {
    children map[rune]*AutocompleteNode
    counts   map[string]int // word → count for subtree
}

func newACNode() *AutocompleteNode {
    return &AutocompleteNode{
        children: make(map[rune]*AutocompleteNode),
        counts:   make(map[string]int),
    }
}

type AutocompleteSystem struct {
    root    *AutocompleteNode
    current *AutocompleteNode
    input   []rune
}

func AutocompleteConstructor(sentences []string, times []int) AutocompleteSystem {
    root := newACNode()
    sys := AutocompleteSystem{root: root, current: root}
    for i, s := range sentences {
        sys.addSentence(s, times[i])
    }
    return sys
}

func (s *AutocompleteSystem) addSentence(sentence string, count int) {
    node := s.root
    for _, c := range sentence {
        if node.children[c] == nil {
            node.children[c] = newACNode()
        }
        node = node.children[c]
        node.counts[sentence] += count
    }
}

func (s *AutocompleteSystem) Input(c rune) []string {
    if c == '#' {
        // Commit current input
        word := string(s.input)
        s.addSentence(word, 1)
        s.input = s.input[:0]
        s.current = s.root
        return []string{}
    }

    s.input = append(s.input, c)
    if s.current != nil {
        s.current = s.current.children[c]
    }
    if s.current == nil {
        return []string{}
    }

    // Sort by frequency desc, then lexicographically asc
    type entry struct {
        word  string
        count int
    }
    entries := make([]entry, 0, len(s.current.counts))
    for w, cnt := range s.current.counts {
        entries = append(entries, entry{w, cnt})
    }
    sort.Slice(entries, func(i, j int) bool {
        if entries[i].count != entries[j].count {
            return entries[i].count > entries[j].count
        }
        return entries[i].word < entries[j].word
    })

    result := []string{}
    for i := 0; i < len(entries) && i < 3; i++ {
        result = append(result, entries[i].word)
    }
    return result
}
```
