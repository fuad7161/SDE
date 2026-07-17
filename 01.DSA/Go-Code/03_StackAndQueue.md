# Stack & Queue

> Stacks handle LIFO order and are essential for parsing, monotonic patterns, and backtracking. Queues handle FIFO and underpin BFS.

---

## Table of Contents

1. [Stack & Queue Basics in Go](#1-stack--queue-basics-in-go)
2. [Valid Parentheses](#2-valid-parentheses)
3. [Min Stack](#3-min-stack)
4. [Next Greater Element](#4-next-greater-element)
5. [Implement Queue Using Stacks](#5-implement-queue-using-stacks)
6. [Largest Rectangle in Histogram](#6-largest-rectangle-in-histogram)
7. [Daily Temperatures](#7-daily-temperatures)
8. [Monotonic Stack — Template](#8-monotonic-stack--template)

---

## 1. Stack & Queue Basics in Go

```go
// Stack using a slice
stack := []int{}
stack = append(stack, 10)           // push
top := stack[len(stack)-1]          // peek
stack = stack[:len(stack)-1]        // pop
isEmpty := len(stack) == 0

// Queue using a slice
queue := []int{}
queue = append(queue, 10)           // enqueue
front := queue[0]                   // front
queue = queue[1:]                   // dequeue

// Note: for performance-critical queues, use a linked list or ring buffer.
// For interview purposes, slice-based queue is acceptable.

// Deque using two slices (or just manipulate both ends of a slice)
deque := []int{}
deque = append(deque, 10)           // push back
deque = append([]int{10}, deque...) // push front
back := deque[len(deque)-1]         // peek back
_ = back
deque = deque[:len(deque)-1]        // pop back
front2 := deque[0]                  // peek front
_ = front2
deque = deque[1:]                   // pop front
```

---

## 2. Valid Parentheses

```go
// Time: O(n), Space: O(n)
func isValid(s string) bool {
    stack := []rune{}
    pair := map[rune]rune{')': '(', '}': '{', ']': '['}

    for _, c := range s {
        if c == '(' || c == '{' || c == '[' {
            stack = append(stack, c)
        } else {
            if len(stack) == 0 || stack[len(stack)-1] != pair[c] {
                return false
            }
            stack = stack[:len(stack)-1] // pop
        }
    }
    return len(stack) == 0
}
```

> **Interview Q: What if we have `*` as a wildcard (matches `(`, `)`, or empty)?**  
> Track two variables: `lo` (minimum open count) and `hi` (maximum open count). `*` increments hi and decrements lo (clamped at 0). Invalid if hi < 0. Valid if lo == 0 at the end.

---

## 3. Min Stack

```go
// All operations O(1)
type MinStack struct {
    stack    []int
    minStack []int
}

func MinStackConstructor() MinStack {
    return MinStack{}
}

func (s *MinStack) Push(val int) {
    s.stack = append(s.stack, val)
    if len(s.minStack) == 0 || val <= s.minStack[len(s.minStack)-1] {
        s.minStack = append(s.minStack, val)
    }
}

func (s *MinStack) Pop() {
    top := s.stack[len(s.stack)-1]
    s.stack = s.stack[:len(s.stack)-1]
    if top == s.minStack[len(s.minStack)-1] {
        s.minStack = s.minStack[:len(s.minStack)-1]
    }
}

func (s *MinStack) Top() int {
    return s.stack[len(s.stack)-1]
}

func (s *MinStack) GetMin() int {
    return s.minStack[len(s.minStack)-1]
}
```

> **Interview Q: Why do we need a separate minStack instead of a single min variable?**  
> When we pop the current minimum, we need to know the *previous* minimum. A min stack tracks the minimum at each level.

---

## 4. Next Greater Element

```go
// Next Greater Element I — Time: O(n), Space: O(n)
// Given nums1 ⊆ nums2, find next greater in nums2 for each element of nums1
func nextGreaterElement(nums1 []int, nums2 []int) []int {
    nextGreater := make(map[int]int)
    stack := []int{} // monotonic decreasing stack

    for _, num := range nums2 {
        for len(stack) > 0 && stack[len(stack)-1] < num {
            top := stack[len(stack)-1]
            stack = stack[:len(stack)-1]
            nextGreater[top] = num
        }
        stack = append(stack, num)
    }
    // anything remaining in stack has no next greater → -1 (default)
    result := make([]int, len(nums1))
    for i, num := range nums1 {
        if ng, ok := nextGreater[num]; ok {
            result[i] = ng
        } else {
            result[i] = -1
        }
    }
    return result
}

// Next Greater Element II (circular array)
func nextGreaterElements(nums []int) []int {
    n := len(nums)
    result := make([]int, n)
    for i := range result { result[i] = -1 }
    stack := []int{} // stores indices

    for i := 0; i < 2*n; i++ {
        for len(stack) > 0 && nums[stack[len(stack)-1]] < nums[i%n] {
            idx := stack[len(stack)-1]
            stack = stack[:len(stack)-1]
            result[idx] = nums[i%n]
        }
        if i < n {
            stack = append(stack, i)
        }
    }
    return result
}
```

---

## 5. Implement Queue Using Stacks

```go
// Amortized O(1) Push and Pop
type MyQueue struct {
    inStack  []int
    outStack []int
}

func QueueConstructor() MyQueue {
    return MyQueue{}
}

func (q *MyQueue) Push(x int) {
    q.inStack = append(q.inStack, x)
}

func (q *MyQueue) transfer() {
    if len(q.outStack) == 0 {
        for len(q.inStack) > 0 {
            top := q.inStack[len(q.inStack)-1]
            q.inStack = q.inStack[:len(q.inStack)-1]
            q.outStack = append(q.outStack, top)
        }
    }
}

func (q *MyQueue) Pop() int {
    q.transfer()
    top := q.outStack[len(q.outStack)-1]
    q.outStack = q.outStack[:len(q.outStack)-1]
    return top
}

func (q *MyQueue) Peek() int {
    q.transfer()
    return q.outStack[len(q.outStack)-1]
}

func (q *MyQueue) Empty() bool {
    return len(q.inStack) == 0 && len(q.outStack) == 0
}
```

---

## 6. Largest Rectangle in Histogram

```go
// Time: O(n), Space: O(n)
func largestRectangleArea(heights []int) int {
    stack := []int{} // indices — monotonic increasing by height
    maxArea := 0
    // append sentinel 0 to flush the stack at the end
    heights = append(heights, 0)

    for i, h := range heights {
        for len(stack) > 0 && heights[stack[len(stack)-1]] > h {
            top := stack[len(stack)-1]
            stack = stack[:len(stack)-1]
            width := i
            if len(stack) > 0 {
                width = i - stack[len(stack)-1] - 1
            }
            area := heights[top] * width
            if area > maxArea {
                maxArea = area
            }
        }
        stack = append(stack, i)
    }
    return maxArea
}
```

> **Interview Q: Why append a 0 sentinel to heights?**  
> It guarantees the stack is fully drained at the end, processing all remaining bars whose right boundary is the end of the array.

---

## 7. Daily Temperatures

```go
// Time: O(n), Space: O(n)
func dailyTemperatures(temperatures []int) []int {
    n := len(temperatures)
    result := make([]int, n)
    stack := []int{} // indices — monotonic decreasing by temperature

    for i, temp := range temperatures {
        for len(stack) > 0 && temperatures[stack[len(stack)-1]] < temp {
            idx := stack[len(stack)-1]
            stack = stack[:len(stack)-1]
            result[idx] = i - idx
        }
        stack = append(stack, i)
    }
    return result
}
// [73,74,75,71,69,72,76,73] → [1,1,4,2,1,1,0,0]
```

---

## 8. Monotonic Stack — Template

```go
// ── Next Greater (monotonic decreasing stack) ──
func nextGreaterTemplate(nums []int) []int {
    n := len(nums)
    result := make([]int, n)
    for i := range result { result[i] = -1 }
    stack := []int{} // stores indices

    for i := 0; i < n; i++ {
        for len(stack) > 0 && nums[stack[len(stack)-1]] < nums[i] {
            idx := stack[len(stack)-1]
            stack = stack[:len(stack)-1]
            result[idx] = nums[i]
        }
        stack = append(stack, i)
    }
    return result
}

// ── Previous Smaller (monotonic increasing stack — iterate left to right) ──
func previousSmallerTemplate(nums []int) []int {
    n := len(nums)
    result := make([]int, n)
    for i := range result { result[i] = -1 }
    stack := []int{}

    for i := 0; i < n; i++ {
        for len(stack) > 0 && nums[stack[len(stack)-1]] >= nums[i] {
            stack = stack[:len(stack)-1]
        }
        if len(stack) > 0 {
            result[i] = nums[stack[len(stack)-1]]
        }
        stack = append(stack, i)
    }
    return result
}
```

> **Monotonic stack cheatsheet:**
> | Pattern | Stack type | Pop condition | Fills |
> |---|---|---|---|
> | Next Greater | Decreasing | `stack.top < current` | result[popped] |
> | Next Smaller | Increasing | `stack.top > current` | result[popped] |
> | Prev Greater | Decreasing | build left→right | result[i] from stack top |
> | Prev Smaller | Increasing | build left→right | result[i] from stack top |
