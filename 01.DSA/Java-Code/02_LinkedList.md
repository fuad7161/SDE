# Linked List

> Classic interview topic. Master reversal, two-pointer (slow/fast), and dummy node patterns.

---

## Table of Contents

1. [Reverse a Linked List](#1-reverse-a-linked-list)
2. [Detect Cycle — Floyd's Algorithm](#2-detect-cycle--floyds-algorithm)
3. [Merge Two Sorted Lists](#3-merge-two-sorted-lists)
4. [Find Middle Element](#4-find-middle-element)
5. [Remove Nth Node from End](#5-remove-nth-node-from-end)
6. [LRU Cache](#6-lru-cache)
7. [Flatten a Multilevel Linked List](#7-flatten-a-multilevel-linked-list)

---

## Node Definition

```java
class ListNode {
    int val;
    ListNode next;
    ListNode(int val) { this.val = val; }
}
```

---

## 1. Reverse a Linked List

### Iterative

```java
// Time: O(n), Space: O(1)
public ListNode reverseList(ListNode head) {
    ListNode prev = null;
    ListNode curr = head;

    while (curr != null) {
        ListNode next = curr.next;   // save next
        curr.next = prev;            // reverse pointer
        prev = curr;                 // advance prev
        curr = next;                 // advance curr
    }
    return prev;   // prev is the new head
}
// null ← 1 ← 2 ← 3 ← 4 ← 5
//                           ^
//                         (new head = prev)
```

### Recursive

```java
// Time: O(n), Space: O(n) — call stack
public ListNode reverseListRecursive(ListNode head) {
    if (head == null || head.next == null) return head;

    ListNode newHead = reverseListRecursive(head.next);
    head.next.next = head;   // node ahead points back to current
    head.next = null;        // current stops pointing forward
    return newHead;
}
```

> **Interview Q: Iterative vs recursive reversal — which is preferred?**  
> **Iterative** is preferred in interviews — O(1) space, no stack overflow risk for large lists. Recursive is elegant but uses O(n) stack space.

---

## 2. Detect Cycle — Floyd's Algorithm

**Slow/Fast pointer (tortoise & hare):** If a cycle exists, fast pointer will eventually lap the slow pointer.

```java
// ── Detect if cycle exists ──
// Time: O(n), Space: O(1)
public boolean hasCycle(ListNode head) {
    ListNode slow = head, fast = head;
    while (fast != null && fast.next != null) {
        slow = slow.next;
        fast = fast.next.next;
        if (slow == fast) return true;   // cycle detected
    }
    return false;
}

// ── Find the start of the cycle ──
public ListNode detectCycle(ListNode head) {
    ListNode slow = head, fast = head;

    // Phase 1: detect meeting point
    while (fast != null && fast.next != null) {
        slow = slow.next;
        fast = fast.next.next;
        if (slow == fast) break;
    }
    if (fast == null || fast.next == null) return null;   // no cycle

    // Phase 2: find cycle start
    // Reset one pointer to head — both now move at speed 1
    slow = head;
    while (slow != fast) {
        slow = slow.next;
        fast = fast.next;
    }
    return slow;   // meeting point = cycle start
}
```

**Why does Phase 2 work?**  
If the distance from head to cycle start is `F`, and the cycle length is `C`, when the pointers met, slow had traveled `F + a` steps. It can be proved that `F ≡ C - a (mod C)`, meaning moving one pointer back to head and stepping both by 1 will make them meet exactly at the cycle entrance.

---

## 3. Merge Two Sorted Lists

```java
// ── Iterative (with dummy head) ──
// Time: O(m+n), Space: O(1)
public ListNode mergeTwoLists(ListNode l1, ListNode l2) {
    ListNode dummy = new ListNode(0);
    ListNode curr = dummy;

    while (l1 != null && l2 != null) {
        if (l1.val <= l2.val) {
            curr.next = l1;
            l1 = l1.next;
        } else {
            curr.next = l2;
            l2 = l2.next;
        }
        curr = curr.next;
    }
    curr.next = (l1 != null) ? l1 : l2;   // attach remaining
    return dummy.next;
}

// ── Recursive ──
public ListNode mergeTwoListsRecursive(ListNode l1, ListNode l2) {
    if (l1 == null) return l2;
    if (l2 == null) return l1;

    if (l1.val <= l2.val) {
        l1.next = mergeTwoListsRecursive(l1.next, l2);
        return l1;
    } else {
        l2.next = mergeTwoListsRecursive(l1, l2.next);
        return l2;
    }
}
```

> **Interview tip:** Always use a **dummy node** when building a linked list result — it eliminates the special case of the empty output list.

---

## 4. Find Middle Element

**Slow/Fast pointer:** When fast reaches the end, slow is at the middle.

```java
// Time: O(n), Space: O(1)
public ListNode middleNode(ListNode head) {
    ListNode slow = head, fast = head;
    while (fast != null && fast.next != null) {
        slow = slow.next;
        fast = fast.next.next;
    }
    return slow;
}
// For even-length list [1,2,3,4], returns node 3 (second middle)
// Adjust to return first middle: while (fast.next != null && fast.next.next != null)
```

---

## 5. Remove Nth Node from End

**Two-pointer trick:** Advance fast pointer `n` steps ahead, then move both until fast reaches the end.

```java
// Time: O(n), Space: O(1) — single pass with dummy node
public ListNode removeNthFromEnd(ListNode head, int n) {
    ListNode dummy = new ListNode(0);
    dummy.next = head;
    ListNode fast = dummy, slow = dummy;

    // Move fast n+1 steps ahead
    for (int i = 0; i <= n; i++) fast = fast.next;

    // Move both until fast hits end
    while (fast != null) {
        fast = fast.next;
        slow = slow.next;
    }

    // slow.next is the node to remove
    slow.next = slow.next.next;
    return dummy.next;
}
```

---

## 6. LRU Cache

**Design:** `get` and `put` both in O(1).  
**Data structure:** `HashMap<key, Node>` + **Doubly Linked List** (most recently used at head, LRU at tail).

```java
class LRUCache {
    private final int capacity;
    private final Map<Integer, Node> map;
    private final Node head, tail;   // sentinels (dummy nodes)

    class Node {
        int key, val;
        Node prev, next;
        Node(int key, int val) { this.key = key; this.val = val; }
    }

    public LRUCache(int capacity) {
        this.capacity = capacity;
        this.map = new HashMap<>();
        head = new Node(0, 0);   // dummy head (most recent side)
        tail = new Node(0, 0);   // dummy tail (LRU side)
        head.next = tail;
        tail.prev = head;
    }

    public int get(int key) {
        if (!map.containsKey(key)) return -1;
        Node node = map.get(key);
        moveToFront(node);       // recently accessed → move to front
        return node.val;
    }

    public void put(int key, int value) {
        if (map.containsKey(key)) {
            Node node = map.get(key);
            node.val = value;
            moveToFront(node);
        } else {
            Node node = new Node(key, value);
            map.put(key, node);
            addToFront(node);
            if (map.size() > capacity) {
                Node lru = tail.prev;   // evict least recently used
                remove(lru);
                map.remove(lru.key);
            }
        }
    }

    private void addToFront(Node node) {
        node.next = head.next;
        node.prev = head;
        head.next.prev = node;
        head.next = node;
    }

    private void remove(Node node) {
        node.prev.next = node.next;
        node.next.prev = node.prev;
    }

    private void moveToFront(Node node) {
        remove(node);
        addToFront(node);
    }
}

// Usage:
// LRUCache cache = new LRUCache(2);
// cache.put(1, 1);   // {1=1}
// cache.put(2, 2);   // {1=1, 2=2}
// cache.get(1);      // 1, now 1 is most recent: {2=2, 1=1}
// cache.put(3, 3);   // evict key 2 (LRU): {1=1, 3=3}
// cache.get(2);      // -1 (not found)
```

> **Interview Q: Why doubly linked list + hashmap?**  
> HashMap gives O(1) lookup. Doubly linked list gives O(1) insert/delete anywhere (we hold a direct node reference from the map). A singly linked list would require O(n) to find the previous node for deletion.

---

## 7. Flatten a Multilevel Linked List

**Problem:** A doubly linked list where some nodes have a `child` pointer to another list. Flatten it depth-first.

```java
class Node {
    int val;
    Node prev, next, child;
}

// Time: O(n), Space: O(depth) — call stack
public Node flatten(Node head) {
    if (head == null) return null;
    Node curr = head;

    while (curr != null) {
        if (curr.child != null) {
            Node child = curr.child;
            Node next = curr.next;

            // Connect curr → child
            curr.next = child;
            child.prev = curr;
            curr.child = null;

            // Find end of child list
            Node tail = child;
            while (tail.next != null) tail = tail.next;

            // Connect tail of child → next
            tail.next = next;
            if (next != null) next.prev = tail;
        }
        curr = curr.next;
    }
    return head;
}
```
