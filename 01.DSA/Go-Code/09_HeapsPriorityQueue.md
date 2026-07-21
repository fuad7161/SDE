# Heaps & Priority Queue

> Go does not have a built-in heap type. You must implement the `heap.Interface` from `container/heap`. This is a common Go interview requirement.

---

## Table of Contents

1. [heap.Interface in Go](#1-heapinterface-in-go)
2. [Kth Largest Element](#2-kth-largest-element)
3. [Top K Frequent Elements](#3-top-k-frequent-elements)
4. [Merge K Sorted Lists](#4-merge-k-sorted-lists)
5. [Find Median from Data Stream](#5-find-median-from-data-stream)
6. [Task Scheduler](#6-task-scheduler)

---

## 1. heap.Interface in Go

```go
import "container/heap"

// ── Min-Heap ──
type MinHeap []int

func (h MinHeap) Len() int           { return len(h) }
func (h MinHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h MinHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *MinHeap) Push(x interface{}) {
    *h = append(*h, x.(int))
}
func (h *MinHeap) Pop() interface{} {
    old := *h
    n := len(old)
    x := old[n-1]
    *h = old[:n-1]
    return x
}

// ── Max-Heap — just flip the Less comparison ──
type MaxHeap []int

func (h MaxHeap) Len() int           { return len(h) }
func (h MaxHeap) Less(i, j int) bool { return h[i] > h[j] } // reversed
func (h MaxHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *MaxHeap) Push(x interface{}) { *h = append(*h, x.(int)) }
func (h *MaxHeap) Pop() interface{} {
    old := *h; n := len(old); x := old[n-1]; *h = old[:n-1]; return x
}

// Usage
func heapExample() {
    h := &MinHeap{3, 1, 4, 1, 5}
    heap.Init(h)
    heap.Push(h, 2)
    min := heap.Pop(h).(int) // 1
    _ = min
    peek := (*h)[0]          // peek without removing
    _ = peek
}
```

> **Interview Q: Why does Go require implementing an interface for heaps?**  
> Go's design philosophy favors composition and interfaces over built-in magic. `container/heap` operates on any slice type that satisfies `heap.Interface`, making it generic without generics syntax.

---

## 2. Kth Largest Element

```go
// Min-heap of size k — Time: O(n log k), Space: O(k)
import "container/heap"

func findKthLargest(nums []int, k int) int {
    h := &MinHeap{}
    heap.Init(h)
    for _, num := range nums {
        heap.Push(h, num)
        if h.Len() > k {
            heap.Pop(h) // remove the smallest
        }
    }
    return (*h)[0] // the kth largest is now the min
}

// Quick-select — average O(n), worst O(n²)
import "math/rand"

func findKthLargestQS(nums []int, k int) int {
    target := len(nums) - k
    var qs func(l, r int) int
    qs = func(l, r int) int {
        pivotIdx := l + rand.Intn(r-l+1)
        nums[pivotIdx], nums[r] = nums[r], nums[pivotIdx]
        pivot := nums[r]
        i := l
        for j := l; j < r; j++ {
            if nums[j] <= pivot {
                nums[i], nums[j] = nums[j], nums[i]
                i++
            }
        }
        nums[i], nums[r] = nums[r], nums[i]
        if i == target { return nums[i] }
        if i < target { return qs(i+1, r) }
        return qs(l, i-1)
    }
    return qs(0, len(nums)-1)
}
```

---

## 3. Top K Frequent Elements

```go
// Min-heap approach — Time: O(n log k), Space: O(n)
import "container/heap"

type Pair struct{ val, freq int }
type PairMinHeap []Pair

func (h PairMinHeap) Len() int           { return len(h) }
func (h PairMinHeap) Less(i, j int) bool { return h[i].freq < h[j].freq }
func (h PairMinHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *PairMinHeap) Push(x interface{}) { *h = append(*h, x.(Pair)) }
func (h *PairMinHeap) Pop() interface{} {
    old := *h; n := len(old); x := old[n-1]; *h = old[:n-1]; return x
}

func topKFrequent(nums []int, k int) []int {
    freq := make(map[int]int)
    for _, n := range nums { freq[n]++ }

    h := &PairMinHeap{}
    heap.Init(h)
    for val, cnt := range freq {
        heap.Push(h, Pair{val, cnt})
        if h.Len() > k { heap.Pop(h) }
    }

    result := make([]int, k)
    for i := k - 1; i >= 0; i-- {
        result[i] = heap.Pop(h).(Pair).val
    }
    return result
}

// Bucket sort approach — O(n), O(n) space
func topKFrequentBucket(nums []int, k int) []int {
    freq := make(map[int]int)
    for _, n := range nums { freq[n]++ }

    buckets := make([][]int, len(nums)+1)
    for val, cnt := range freq {
        buckets[cnt] = append(buckets[cnt], val)
    }

    result := []int{}
    for i := len(buckets) - 1; i >= 0 && len(result) < k; i-- {
        result = append(result, buckets[i]...)
    }
    return result[:k]
}
```

---

## 4. Merge K Sorted Lists

```go
// Time: O(N log k) where N = total nodes, k = number of lists
import "container/heap"

type NodeHeap []*ListNode

func (h NodeHeap) Len() int           { return len(h) }
func (h NodeHeap) Less(i, j int) bool { return h[i].Val < h[j].Val }
func (h NodeHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *NodeHeap) Push(x interface{}) { *h = append(*h, x.(*ListNode)) }
func (h *NodeHeap) Pop() interface{} {
    old := *h; n := len(old); x := old[n-1]; *h = old[:n-1]; return x
}

func mergeKLists(lists []*ListNode) *ListNode {
    h := &NodeHeap{}
    heap.Init(h)
    for _, l := range lists {
        if l != nil { heap.Push(h, l) }
    }

    dummy := &ListNode{}
    cur := dummy
    for h.Len() > 0 {
        node := heap.Pop(h).(*ListNode)
        cur.Next = node
        cur = cur.Next
        if node.Next != nil { heap.Push(h, node.Next) }
    }
    return dummy.Next
}
```

---

## 5. Find Median from Data Stream

```go
// Two heaps — Get: O(1), Add: O(log n)
import "container/heap"

type MedianFinder struct {
    lower *MaxHeap // max-heap for lower half
    upper *MinHeap // min-heap for upper half
}

func MedianFinderConstructor() MedianFinder {
    lo := &MaxHeap{}
    hi := &MinHeap{}
    heap.Init(lo)
    heap.Init(hi)
    return MedianFinder{lo, hi}
}

func (mf *MedianFinder) AddNum(num int) {
    heap.Push(mf.lower, num)
    // Balance: largest in lower must be <= smallest in upper
    if mf.upper.Len() > 0 && (*mf.lower)[0] > (*mf.upper)[0] {
        heap.Push(mf.upper, heap.Pop(mf.lower))
    }
    // Maintain size: lower can have at most 1 extra
    if mf.lower.Len() > mf.upper.Len()+1 {
        heap.Push(mf.upper, heap.Pop(mf.lower))
    } else if mf.upper.Len() > mf.lower.Len() {
        heap.Push(mf.lower, heap.Pop(mf.upper))
    }
}

func (mf *MedianFinder) FindMedian() float64 {
    if mf.lower.Len() > mf.upper.Len() {
        return float64((*mf.lower)[0])
    }
    return float64((*mf.lower)[0]+(*mf.upper)[0]) / 2.0
}
```

---

## 6. Task Scheduler

```go
// Time: O(n), Space: O(1) — at most 26 distinct tasks
import "sort"

func leastInterval(tasks []byte, n int) int {
    freq := [26]int{}
    for _, t := range tasks { freq[t-'A']++ }

    vals := freq[:]
    sort.Ints(vals)

    maxFreq := vals[25]
    idleSlots := (maxFreq - 1) * n

    // Fill idle slots with remaining tasks (most frequent first)
    for i := 24; i >= 0 && idleSlots > 0; i-- {
        if vals[i] < maxFreq {
            idleSlots -= vals[i]
        } else {
            idleSlots -= maxFreq - 1
        }
    }
    if idleSlots < 0 { idleSlots = 0 }
    return len(tasks) + idleSlots
}
```

> **Interview Q: Two heaps for median — which half gets the new element first?**  
> Always push to the lower (max) heap first, then re-balance. This ensures the max of the lower half and min of the upper half stay ordered, guaranteeing correct median retrieval.
