# Linked List

> Master pointer manipulation, fast/slow runners, and dummy node techniques.

---

## Table of Contents

1. [Linked List Basics in Go](#1-linked-list-basics-in-go)
2. [Reverse a Linked List](#2-reverse-a-linked-list)
3. [Floyd's Cycle Detection](#3-floyds-cycle-detection)
4. [Merge Two Sorted Lists](#4-merge-two-sorted-lists)
5. [Find Middle of Linked List](#5-find-middle-of-linked-list)
6. [Remove Nth Node From End](#6-remove-nth-node-from-end)
7. [LRU Cache](#7-lru-cache)
8. [Flatten a Multilevel Doubly Linked List](#8-flatten-a-multilevel-doubly-linked-list)

---

## 1. Linked List Basics in Go

```go
// Node definition — used for all problems
type ListNode struct {
    Val  int
    Next *ListNode
}

// Helper — build list from slice
func buildList(vals []int) *ListNode {
    dummy := &ListNode{}
    cur := dummy
    for _, v := range vals {
        cur.Next = &ListNode{Val: v}
        cur = cur.Next
    }
    return dummy.Next
}

// Helper — print list
func printList(head *ListNode) {
    for head != nil {
        fmt.Print(head.Val, " -> ")
        head = head.Next
    }
    fmt.Println("nil")
}
```

---

## 2. Reverse a Linked List

### Iterative

```go
// Time: O(n), Space: O(1)
func reverseList(head *ListNode) *ListNode {
    var prev *ListNode
    cur := head
    for cur != nil {
        next := cur.Next
        cur.Next = prev
        prev = cur
        cur = next
    }
    return prev
}
```

### Recursive

```go
// Time: O(n), Space: O(n) stack
func reverseListRec(head *ListNode) *ListNode {
    if head == nil || head.Next == nil {
        return head
    }
    newHead := reverseListRec(head.Next)
    head.Next.Next = head
    head.Next = nil
    return newHead
}
```

### Reverse Sublist [left, right]

```go
func reverseBetween(head *ListNode, left, right int) *ListNode {
    dummy := &ListNode{Next: head}
    pre := dummy
    for i := 1; i < left; i++ {
        pre = pre.Next
    }
    cur := pre.Next
    for i := 0; i < right-left; i++ {
        next := cur.Next
        cur.Next = next.Next
        next.Next = pre.Next
        pre.Next = next
    }
    return dummy.Next
}
```

> **Interview Q: Iterative vs recursive reversal?**  
> Iterative is O(1) space — always prefer it. Recursive is cleaner to write but uses O(n) call stack — risky on large lists.

---

## 3. Floyd's Cycle Detection

### Detect Cycle

```go
// Time: O(n), Space: O(1)
func hasCycle(head *ListNode) bool {
    slow, fast := head, head
    for fast != nil && fast.Next != nil {
        slow = slow.Next
        fast = fast.Next.Next
        if slow == fast {
            return true
        }
    }
    return false
}
```

### Find Cycle Start

```go
func detectCycle(head *ListNode) *ListNode {
    slow, fast := head, head
    for fast != nil && fast.Next != nil {
        slow = slow.Next
        fast = fast.Next.Next
        if slow == fast {
            slow = head
            for slow != fast {
                slow = slow.Next
                fast = fast.Next
            }
            return slow // cycle start
        }
    }
    return nil
}
```

> **Interview Q: Why does resetting one pointer to head find the cycle start?**  
> The distance from head to the cycle entry equals the distance from the meeting point to the entry (proven mathematically with the 2x speed ratio). Resetting one pointer to head and advancing both by 1 makes them meet at the cycle entry.

---

## 4. Merge Two Sorted Lists

```go
// Time: O(m+n), Space: O(1) iterative
func mergeTwoLists(l1 *ListNode, l2 *ListNode) *ListNode {
    dummy := &ListNode{}
    cur := dummy
    for l1 != nil && l2 != nil {
        if l1.Val <= l2.Val {
            cur.Next = l1
            l1 = l1.Next
        } else {
            cur.Next = l2
            l2 = l2.Next
        }
        cur = cur.Next
    }
    if l1 != nil {
        cur.Next = l1
    } else {
        cur.Next = l2
    }
    return dummy.Next
}

// Recursive variant — cleaner but O(m+n) stack
func mergeTwoListsRec(l1 *ListNode, l2 *ListNode) *ListNode {
    if l1 == nil { return l2 }
    if l2 == nil { return l1 }
    if l1.Val <= l2.Val {
        l1.Next = mergeTwoListsRec(l1.Next, l2)
        return l1
    }
    l2.Next = mergeTwoListsRec(l1, l2.Next)
    return l2
}
```

---

## 5. Find Middle of Linked List

```go
// Slow/fast pointers — stops at the second middle for even length
func middleNode(head *ListNode) *ListNode {
    slow, fast := head, head
    for fast != nil && fast.Next != nil {
        slow = slow.Next
        fast = fast.Next.Next
    }
    return slow
}
```

---

## 6. Remove Nth Node From End

```go
// Time: O(n), Space: O(1) — one-pass with two pointers
func removeNthFromEnd(head *ListNode, n int) *ListNode {
    dummy := &ListNode{Next: head}
    fast, slow := dummy, dummy
    // advance fast n+1 steps
    for i := 0; i <= n; i++ {
        fast = fast.Next
    }
    // move both until fast reaches end
    for fast != nil {
        slow = slow.Next
        fast = fast.Next
    }
    slow.Next = slow.Next.Next // unlink
    return dummy.Next
}
```

---

## 7. LRU Cache

```go
// Get: O(1), Put: O(1)
type LRUNode struct {
    key, val   int
    prev, next *LRUNode
}

type LRUCache struct {
    cap        int
    cache      map[int]*LRUNode
    head, tail *LRUNode // sentinel nodes
}

func Constructor(capacity int) LRUCache {
    head := &LRUNode{}
    tail := &LRUNode{}
    head.next = tail
    tail.prev = head
    return LRUCache{
        cap:   capacity,
        cache: make(map[int]*LRUNode),
        head:  head,
        tail:  tail,
    }
}

func (c *LRUCache) remove(node *LRUNode) {
    node.prev.next = node.next
    node.next.prev = node.prev
}

func (c *LRUCache) insertFront(node *LRUNode) {
    node.next = c.head.next
    node.prev = c.head
    c.head.next.prev = node
    c.head.next = node
}

func (c *LRUCache) Get(key int) int {
    if node, ok := c.cache[key]; ok {
        c.remove(node)
        c.insertFront(node)
        return node.val
    }
    return -1
}

func (c *LRUCache) Put(key int, value int) {
    if node, ok := c.cache[key]; ok {
        c.remove(node)
        node.val = value
        c.insertFront(node)
    } else {
        if len(c.cache) == c.cap {
            lru := c.tail.prev
            c.remove(lru)
            delete(c.cache, lru.key)
        }
        node := &LRUNode{key: key, val: value}
        c.insertFront(node)
        c.cache[key] = node
    }
}
```

> **Interview Q: Why use a doubly linked list for LRU?**  
> A doubly linked list allows O(1) removal from any position (we hold a pointer to the node). Combined with a hash map for O(1) lookup, both Get and Put are O(1).

---

## 8. Flatten a Multilevel Doubly Linked List

```go
type MLNode struct {
    Val   int
    Prev  *MLNode
    Next  *MLNode
    Child *MLNode
}

// Time: O(n), Space: O(n) recursion depth
func flatten(head *MLNode) *MLNode {
    if head == nil {
        return nil
    }
    cur := head
    for cur != nil {
        if cur.Child != nil {
            child := flatten(cur.Child)
            next := cur.Next
            cur.Next = child
            child.Prev = cur
            cur.Child = nil
            // find the tail of flattened child
            tail := child
            for tail.Next != nil {
                tail = tail.Next
            }
            tail.Next = next
            if next != nil {
                next.Prev = tail
            }
        }
        cur = cur.Next
    }
    return head
}
```
