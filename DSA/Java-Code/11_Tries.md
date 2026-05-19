# Tries

> A Trie (prefix tree) is an N-ary tree where each node represents a character. Efficient for prefix-based search, autocomplete, and word lookups — O(L) per operation where L = word length.

---

## Table of Contents

1. [Implement Trie](#1-implement-trie)
2. [Word Search II](#2-word-search-ii)
3. [Auto-Complete System](#3-auto-complete-system)

---

## 1. Implement Trie

```java
class Trie {
    private TrieNode root;

    static class TrieNode {
        TrieNode[] children = new TrieNode[26];
        boolean isEnd = false;
    }

    public Trie() {
        root = new TrieNode();
    }

    // Insert a word — O(L)
    public void insert(String word) {
        TrieNode curr = root;
        for (char c : word.toCharArray()) {
            int idx = c - 'a';
            if (curr.children[idx] == null) {
                curr.children[idx] = new TrieNode();
            }
            curr = curr.children[idx];
        }
        curr.isEnd = true;
    }

    // Search for exact word — O(L)
    public boolean search(String word) {
        TrieNode node = findNode(word);
        return node != null && node.isEnd;
    }

    // Check if any word starts with prefix — O(L)
    public boolean startsWith(String prefix) {
        return findNode(prefix) != null;
    }

    private TrieNode findNode(String s) {
        TrieNode curr = root;
        for (char c : s.toCharArray()) {
            int idx = c - 'a';
            if (curr.children[idx] == null) return null;
            curr = curr.children[idx];
        }
        return curr;
    }
}

// Usage:
// Trie trie = new Trie();
// trie.insert("apple");
// trie.search("apple");    // true
// trie.search("app");      // false
// trie.startsWith("app");  // true
// trie.insert("app");
// trie.search("app");      // true
```

### Trie with Delete

```java
public boolean delete(String word) {
    return deleteHelper(root, word, 0);
}

private boolean deleteHelper(TrieNode curr, String word, int idx) {
    if (idx == word.length()) {
        if (!curr.isEnd) return false;  // word not found
        curr.isEnd = false;
        return isEmpty(curr);            // can delete node if no children
    }
    int i = word.charAt(idx) - 'a';
    if (curr.children[i] == null) return false;

    boolean shouldDelete = deleteHelper(curr.children[i], word, idx + 1);
    if (shouldDelete) {
        curr.children[i] = null;
        return !curr.isEnd && isEmpty(curr);  // delete this node if no other refs
    }
    return false;
}

private boolean isEmpty(TrieNode node) {
    for (TrieNode child : node.children) if (child != null) return false;
    return true;
}
```

> **Interview Q: What are the advantages of a Trie over a HashMap for string storage?**  
> A Trie allows **prefix search** (`startsWith`) in O(L) without iterating all keys, enables **lexicographic iteration**, and groups words sharing prefixes (memory efficient for many similar words). A HashMap needs O(1) for exact lookup but O(n*L) to find all words with a given prefix. Trie is preferred for autocomplete, spell checkers, and IP routing.

---

## 2. Word Search II

**Problem:** Given a board of characters and a list of words, find all words present in the board (connected path, no cell reuse).  
**Key:** Build a Trie from the word list, then DFS on the board pruning with the Trie.

```java
class WordSearchII {
    static class TrieNode {
        TrieNode[] children = new TrieNode[26];
        String word = null;  // store complete word at terminal node
    }

    public List<String> findWords(char[][] board, String[] words) {
        TrieNode root = buildTrie(words);
        List<String> result = new ArrayList<>();
        int rows = board.length, cols = board[0].length;

        for (int r = 0; r < rows; r++) {
            for (int c = 0; c < cols; c++) {
                dfs(board, r, c, root, result);
            }
        }
        return result;
    }

    private void dfs(char[][] board, int r, int c, TrieNode node, List<String> result) {
        if (r < 0 || c < 0 || r >= board.length || c >= board[0].length) return;
        char ch = board[r][c];
        if (ch == '#' || node.children[ch - 'a'] == null) return;  // visited or no prefix

        TrieNode next = node.children[ch - 'a'];
        if (next.word != null) {
            result.add(next.word);
            next.word = null;   // avoid duplicates
        }

        board[r][c] = '#';   // mark visited
        dfs(board, r+1, c, next, result);
        dfs(board, r-1, c, next, result);
        dfs(board, r, c+1, next, result);
        dfs(board, r, c-1, next, result);
        board[r][c] = ch;    // restore

        // Optimization: prune empty Trie branches
        if (isEmpty(next)) node.children[ch - 'a'] = null;
    }

    private TrieNode buildTrie(String[] words) {
        TrieNode root = new TrieNode();
        for (String word : words) {
            TrieNode curr = root;
            for (char c : word.toCharArray()) {
                int idx = c - 'a';
                if (curr.children[idx] == null) curr.children[idx] = new TrieNode();
                curr = curr.children[idx];
            }
            curr.word = word;
        }
        return root;
    }

    private boolean isEmpty(TrieNode node) {
        for (TrieNode c : node.children) if (c != null) return false;
        return true;
    }
}
```

> **Interview Q: Why use a Trie instead of a HashSet for Word Search II?**  
> With a HashSet, every DFS path of length L would do an O(L) hash lookup at each step, and you'd have no way to prune early. The Trie prunes entire DFS branches when no word in the list shares the current prefix — significantly reducing the search space on typical inputs.

---

## 3. Auto-Complete System

**Design:** A system that returns top-3 historically typed sentences matching a given prefix.

```java
class AutocompleteSystem {
    static class TrieNode {
        Map<Character, TrieNode> children = new HashMap<>();
        Map<String, Integer> counts = new HashMap<>();  // sentence → count at this node
    }

    private TrieNode root = new TrieNode();
    private TrieNode curr;
    private StringBuilder input = new StringBuilder();

    public AutocompleteSystem(String[] sentences, int[] times) {
        for (int i = 0; i < sentences.length; i++) insert(sentences[i], times[i]);
        curr = root;
    }

    // Called for each character typed. '#' means end of sentence.
    public List<String> input(char c) {
        if (c == '#') {
            insert(input.toString(), 1);
            input.setLength(0);
            curr = root;
            return new ArrayList<>();
        }

        input.append(c);
        if (curr != null) {
            curr = curr.children.get(c);
        }

        if (curr == null) return new ArrayList<>();

        // Return top 3: more frequent first, then lexicographic
        return curr.counts.entrySet().stream()
            .sorted((a, b) -> a.getValue().equals(b.getValue())
                ? a.getKey().compareTo(b.getKey())
                : b.getValue() - a.getValue())
            .limit(3)
            .map(Map.Entry::getKey)
            .collect(java.util.stream.Collectors.toList());
    }

    private void insert(String sentence, int count) {
        TrieNode node = root;
        for (char c : sentence.toCharArray()) {
            node.children.computeIfAbsent(c, k -> new TrieNode());
            node = node.children.get(c);
            node.counts.merge(sentence, count, Integer::sum);
        }
    }
}

// Usage:
// AutocompleteSystem sys = new AutocompleteSystem(
//     new String[]{"i love you", "island", "iroman", "i love leetcode"},
//     new int[]{5, 3, 2, 2}
// );
// sys.input('i');  // returns ["i love you", "island", "i love leetcode"]
// sys.input(' ');  // returns ["i love you", "i love leetcode"]
// sys.input('#');  // saves "i " — no results
```
