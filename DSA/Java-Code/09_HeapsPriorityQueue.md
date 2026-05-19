# Heaps / Priority Queue

> A heap is a complete binary tree satisfying the heap property. Java's `PriorityQueue` is a min-heap by default. Use `Collections.reverseOrder()` for max-heap.

---

## Table of Contents

1. [PriorityQueue Basics](#1-priorityqueue-basics)
2. [Kth Largest Element](#2-kth-largest-element)
3. [Top K Frequent Elements](#3-top-k-frequent-elements)
4. [Merge K Sorted Lists](#4-merge-k-sorted-lists)
5. [Find Median from Data Stream](#5-find-median-from-data-stream)
6. [Task Scheduler](#6-task-scheduler)

---

## 1. PriorityQueue Basics

```java
// Min-heap (default)
PriorityQueue<Integer> minHeap = new PriorityQueue<>();
minHeap.offer(5); minHeap.offer(1); minHeap.offer(3);
minHeap.peek();   // 1 — smallest at top
minHeap.poll();   // 1 — removes and returns smallest

// Max-heap
PriorityQueue<Integer> maxHeap = new PriorityQueue<>(Collections.reverseOrder());
// or: new PriorityQueue<>((a, b) -> b - a)

// Custom object min-heap (by length)
PriorityQueue<String> byLength = new PriorityQueue<>(Comparator.comparingInt(String::length));

// Key operations — all O(log n)
// offer(e)  — insert
// poll()    — remove + return min/max
// peek()    — view min/max without removing
// size()
// isEmpty()
```

---

## 2. Kth Largest Element

**Approach:** Use a **min-heap of size k**. After processing all elements, the root is the kth largest.

```java
// Time: O(n log k), Space: O(k)
public int findKthLargest(int[] nums, int k) {
    PriorityQueue<Integer> minHeap = new PriorityQueue<>();

    for (int num : nums) {
        minHeap.offer(num);
        if (minHeap.size() > k) {
            minHeap.poll();   // remove the smallest — keep only k largest
        }
    }
    return minHeap.peek();   // root = kth largest
}

// ── QuickSelect — O(n) average, O(n²) worst ──
public int findKthLargestQuickSelect(int[] nums, int k) {
    return quickSelect(nums, 0, nums.length - 1, nums.length - k);
}

private int quickSelect(int[] nums, int lo, int hi, int targetIdx) {
    int pivot = nums[hi];
    int i = lo;
    for (int j = lo; j < hi; j++) {
        if (nums[j] <= pivot) { int tmp = nums[i]; nums[i] = nums[j]; nums[j] = tmp; i++; }
    }
    int tmp = nums[i]; nums[i] = nums[hi]; nums[hi] = tmp;  // place pivot

    if (i == targetIdx)   return nums[i];
    if (i < targetIdx)    return quickSelect(nums, i + 1, hi, targetIdx);
    return quickSelect(nums, lo, i - 1, targetIdx);
}
```

> **Interview Q: When to use min-heap vs QuickSelect for Kth largest?**  
> Use **min-heap of size k** when the input is a stream or very large (online algorithm, O(n log k)). Use **QuickSelect** for arrays when average O(n) is acceptable and you don't need worst-case guarantees. QuickSelect modifies the input array.

---

## 3. Top K Frequent Elements

```java
// Time: O(n log k), Space: O(n)
public int[] topKFrequent(int[] nums, int k) {
    // Step 1: Count frequencies
    Map<Integer, Integer> freq = new HashMap<>();
    for (int num : nums) freq.merge(num, 1, Integer::sum);

    // Step 2: Min-heap of size k (keeps top k by frequency)
    PriorityQueue<Map.Entry<Integer, Integer>> minHeap =
        new PriorityQueue<>(Comparator.comparingInt(Map.Entry::getValue));

    for (Map.Entry<Integer, Integer> entry : freq.entrySet()) {
        minHeap.offer(entry);
        if (minHeap.size() > k) minHeap.poll();   // remove least frequent
    }

    int[] result = new int[k];
    for (int i = k - 1; i >= 0; i--) result[i] = minHeap.poll().getKey();
    return result;
}

// ── Bucket Sort approach — O(n) ──
public int[] topKFrequentBucket(int[] nums, int k) {
    Map<Integer, Integer> freq = new HashMap<>();
    for (int num : nums) freq.merge(num, 1, Integer::sum);

    // Buckets indexed by frequency (max frequency = nums.length)
    List<List<Integer>> buckets = new ArrayList<>();
    for (int i = 0; i <= nums.length; i++) buckets.add(new ArrayList<>());
    for (Map.Entry<Integer, Integer> e : freq.entrySet())
        buckets.get(e.getValue()).add(e.getKey());

    int[] result = new int[k];
    int idx = 0;
    for (int i = nums.length; i >= 0 && idx < k; i--)
        for (int num : buckets.get(i))
            if (idx < k) result[idx++] = num;

    return result;
}
```

---

## 4. Merge K Sorted Lists

**Approach:** Use a min-heap. Start with the head of each list. Always pop the smallest node, add to result, then push its next node.

```java
class ListNode { int val; ListNode next; ListNode(int val) { this.val = val; } }

// Time: O(N log k) where N = total nodes, k = number of lists
public ListNode mergeKLists(ListNode[] lists) {
    PriorityQueue<ListNode> minHeap = new PriorityQueue<>(Comparator.comparingInt(n -> n.val));

    // Initialize heap with head of each list
    for (ListNode node : lists) {
        if (node != null) minHeap.offer(node);
    }

    ListNode dummy = new ListNode(0);
    ListNode curr = dummy;

    while (!minHeap.isEmpty()) {
        ListNode node = minHeap.poll();
        curr.next = node;
        curr = curr.next;
        if (node.next != null) minHeap.offer(node.next);   // push next node of same list
    }
    return dummy.next;
}
```

---

## 5. Find Median from Data Stream

**Design:** Support `addNum(int)` and `findMedian()`. Use two heaps: max-heap for lower half, min-heap for upper half.

```java
class MedianFinder {
    private PriorityQueue<Integer> lower;  // max-heap (lower half)
    private PriorityQueue<Integer> upper;  // min-heap (upper half)

    public MedianFinder() {
        lower = new PriorityQueue<>(Collections.reverseOrder());
        upper = new PriorityQueue<>();
    }

    // O(log n)
    public void addNum(int num) {
        lower.offer(num);
        upper.offer(lower.poll());   // balance: push largest of lower to upper

        // Keep sizes balanced: |lower| == |upper| or |lower| == |upper| + 1
        if (lower.size() < upper.size()) {
            lower.offer(upper.poll());
        }
    }

    // O(1)
    public double findMedian() {
        if (lower.size() > upper.size()) return lower.peek();
        return (lower.peek() + upper.peek()) / 2.0;
    }
}

// Example:
// addNum(1): lower=[1], upper=[]       median=1.0
// addNum(2): lower=[1], upper=[2]      median=1.5
// addNum(3): lower=[2,1], upper=[3]    median=2.0
```

> **Interview Q: Why two heaps for median finding?**  
> The median is always at the boundary between two sorted halves. A max-heap gives the largest of the lower half in O(1), and a min-heap gives the smallest of the upper half in O(1). By keeping them balanced in size, the median is always one of these two tops.

---

## 6. Task Scheduler

**Problem:** Given tasks with cooldown `n` (same task must wait `n` intervals), find minimum time to finish all tasks.

```java
// Time: O(n log n or 26 log 26), Space: O(1) — 26 task types
public int leastInterval(char[] tasks, int n) {
    int[] freq = new int[26];
    for (char t : tasks) freq[t - 'A']++;
    Arrays.sort(freq);

    int maxFreq = freq[25];
    int idleSlots = (maxFreq - 1) * n;   // idle slots needed around most frequent task

    // Fill idle slots with other tasks
    for (int i = 24; i >= 0; i--) {
        idleSlots -= Math.min(freq[i], maxFreq - 1);
    }
    return tasks.length + Math.max(0, idleSlots);
}
// tasks = ["A","A","A","B","B","B"], n = 2
// A → B → idle → A → B → idle → A → B = 8 intervals
// Formula: tasks.length + max(0, idle) = 6 + max(0, 2) = 8

// ── Simulation with max-heap (handles ordering explicitly) ──
public int leastIntervalHeap(char[] tasks, int n) {
    int[] freq = new int[26];
    for (char t : tasks) freq[t - 'A']++;

    PriorityQueue<Integer> maxHeap = new PriorityQueue<>(Collections.reverseOrder());
    for (int f : freq) if (f > 0) maxHeap.offer(f);

    int time = 0;
    while (!maxHeap.isEmpty()) {
        List<Integer> window = new ArrayList<>();
        for (int i = 0; i <= n; i++) {
            if (!maxHeap.isEmpty()) window.add(maxHeap.poll());
        }
        for (int f : window) if (f - 1 > 0) maxHeap.offer(f - 1);
        time += maxHeap.isEmpty() ? window.size() : n + 1;
    }
    return time;
}
```
