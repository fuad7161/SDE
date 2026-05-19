# Stack & Queue

> Stack is LIFO; Queue is FIFO. Monotonic stack solves "next greater/smaller" problems in O(n).

---

## Table of Contents

1. [Valid Parentheses](#1-valid-parentheses)
2. [Min Stack](#2-min-stack)
3. [Next Greater Element](#3-next-greater-element)
4. [Implement Queue using Stack](#4-implement-queue-using-stack)
5. [Largest Rectangle in Histogram](#5-largest-rectangle-in-histogram)
6. [Daily Temperatures](#6-daily-temperatures)
7. [Monotonic Stack — Pattern Summary](#7-monotonic-stack--pattern-summary)

---

## 1. Valid Parentheses

**Problem:** Given a string of brackets, determine if it is valid (properly opened and closed).

```java
// Time: O(n), Space: O(n)
public boolean isValid(String s) {
    Deque<Character> stack = new ArrayDeque<>();
    for (char c : s.toCharArray()) {
        if (c == '(' || c == '[' || c == '{') {
            stack.push(c);
        } else {
            if (stack.isEmpty()) return false;
            char top = stack.pop();
            if (c == ')' && top != '(') return false;
            if (c == ']' && top != '[') return false;
            if (c == '}' && top != '{') return false;
        }
    }
    return stack.isEmpty();   // all opened must be closed
}
```

> **Tip:** Use `Deque<Character>` with `push`/`pop` for stack behavior in Java. `Stack<>` class is legacy and synchronized (slower).

---

## 2. Min Stack

**Design:** Stack that supports `push`, `pop`, `top`, and `getMin` — all in O(1).  
**Trick:** Store the minimum alongside each element using a second stack.

```java
class MinStack {
    private Deque<Integer> stack = new ArrayDeque<>();
    private Deque<Integer> minStack = new ArrayDeque<>();

    public void push(int val) {
        stack.push(val);
        // Push new min (keep current min if val is larger)
        int currentMin = minStack.isEmpty() ? val : Math.min(val, minStack.peek());
        minStack.push(currentMin);
    }

    public void pop() {
        stack.pop();
        minStack.pop();
    }

    public int top() {
        return stack.peek();
    }

    public int getMin() {
        return minStack.peek();
    }
}

// ── Alternative: store (val, minAtTime) pairs in one stack ──
class MinStackV2 {
    private Deque<int[]> stack = new ArrayDeque<>();

    public void push(int val) {
        int min = stack.isEmpty() ? val : Math.min(val, stack.peek()[1]);
        stack.push(new int[]{val, min});
    }
    public void pop()       { stack.pop(); }
    public int top()        { return stack.peek()[0]; }
    public int getMin()     { return stack.peek()[1]; }
}
```

> **Interview Q: Why use a parallel min-stack instead of recalculating?**  
> When a new minimum is pushed, the old minimum is "remembered" in the minStack. When that minimum is popped, the previous minimum is automatically restored. Recalculating after every pop would be O(n).

---

## 3. Next Greater Element

**Problem:** For each element, find the first greater element to its right. Return -1 if none.

```java
// ── Next Greater Element I ──
// Time: O(n), Space: O(n) — monotonic decreasing stack
public int[] nextGreaterElement(int[] nums) {
    int n = nums.length;
    int[] result = new int[n];
    Arrays.fill(result, -1);
    Deque<Integer> stack = new ArrayDeque<>();  // stores indices

    for (int i = 0; i < n; i++) {
        // Pop all elements smaller than current — current is their next greater
        while (!stack.isEmpty() && nums[stack.peek()] < nums[i]) {
            result[stack.pop()] = nums[i];
        }
        stack.push(i);
    }
    return result;
}
// [2, 1, 2, 4, 3]
// Result: [4, 2, 4, -1, -1]

// ── Next Greater Element II (circular array) ──
public int[] nextGreaterElements(int[] nums) {
    int n = nums.length;
    int[] result = new int[n];
    Arrays.fill(result, -1);
    Deque<Integer> stack = new ArrayDeque<>();

    // Traverse twice to simulate circular behavior
    for (int i = 0; i < 2 * n; i++) {
        int idx = i % n;
        while (!stack.isEmpty() && nums[stack.peek()] < nums[idx]) {
            result[stack.pop()] = nums[idx];
        }
        if (i < n) stack.push(idx);
    }
    return result;
}
```

---

## 4. Implement Queue using Stack

**Two-stack trick:** `stackIn` for push, `stackOut` for pop. Transfer from `stackIn` to `stackOut` lazily.

```java
class MyQueue {
    private Deque<Integer> stackIn  = new ArrayDeque<>();
    private Deque<Integer> stackOut = new ArrayDeque<>();

    public void push(int x) {
        stackIn.push(x);
    }

    public int pop() {
        transfer();
        return stackOut.pop();
    }

    public int peek() {
        transfer();
        return stackOut.peek();
    }

    public boolean empty() {
        return stackIn.isEmpty() && stackOut.isEmpty();
    }

    private void transfer() {
        // Only transfer when stackOut is empty (lazy approach)
        if (stackOut.isEmpty()) {
            while (!stackIn.isEmpty()) {
                stackOut.push(stackIn.pop());
            }
        }
    }
}
// Amortized O(1) per operation — each element transferred at most once
```

> **Interview Q: What is the amortized time complexity?**  
> Each element is pushed to `stackIn` once and popped from `stackOut` once — O(2) per element total = **amortized O(1)** per operation, though a single pop can be O(n) in the worst case.

---

## 5. Largest Rectangle in Histogram

**Problem:** Given bar heights, find the largest rectangle area in the histogram.  
**Approach:** Monotonic increasing stack — for each bar, compute the max width rectangle using it as the shortest bar.

```java
// Time: O(n), Space: O(n)
public int largestRectangleArea(int[] heights) {
    Deque<Integer> stack = new ArrayDeque<>();  // indices, monotonic increasing height
    int maxArea = 0;
    int n = heights.length;

    for (int i = 0; i <= n; i++) {
        int currHeight = (i == n) ? 0 : heights[i];  // sentinel 0 at end to flush stack

        while (!stack.isEmpty() && heights[stack.peek()] > currHeight) {
            int height = heights[stack.pop()];
            int width = stack.isEmpty() ? i : i - stack.peek() - 1;
            maxArea = Math.max(maxArea, height * width);
        }
        stack.push(i);
    }
    return maxArea;
}
// heights = [2,1,5,6,2,3]
// Max area = 10 (bars 3,4 with height 5: width=2, 5*2=10)
```

---

## 6. Daily Temperatures

**Problem:** For each day, find how many days until a warmer temperature. Return 0 if no warmer day exists.

```java
// Time: O(n), Space: O(n) — monotonic decreasing stack of indices
public int[] dailyTemperatures(int[] temperatures) {
    int n = temperatures.length;
    int[] result = new int[n];
    Deque<Integer> stack = new ArrayDeque<>();  // indices

    for (int i = 0; i < n; i++) {
        while (!stack.isEmpty() && temperatures[stack.peek()] < temperatures[i]) {
            int idx = stack.pop();
            result[idx] = i - idx;   // days to wait = distance
        }
        stack.push(i);
    }
    // Remaining indices in stack have result 0 (default)
    return result;
}
// temperatures = [73,74,75,71,69,72,76,73]
// result       = [ 1, 1, 4, 2, 1, 1, 0, 0]
```

---

## 7. Monotonic Stack — Pattern Summary

A **monotonic stack** maintains elements in increasing or decreasing order.

| Problem Type | Stack Order | What to push | When to pop |
|---|---|---|---|
| Next Greater Element | Decreasing | index | when `nums[i] > nums[top]` |
| Next Smaller Element | Increasing | index | when `nums[i] < nums[top]` |
| Previous Greater | Decreasing | index | keep stack state as you go |
| Daily Temperatures | Decreasing | index | when temp[i] > temp[top] |
| Histogram Rectangle | Increasing | index | when height[i] < height[top] |

**Template:**

```java
// Next Greater — monotonic decreasing stack
Deque<Integer> stack = new ArrayDeque<>();
for (int i = 0; i < n; i++) {
    while (!stack.isEmpty() && nums[stack.peek()] < nums[i]) {
        // nums[i] is the answer for the element at stack.peek()
        result[stack.pop()] = nums[i];
    }
    stack.push(i);
}
// Unprocessed elements in stack → no next greater, answer = -1
```

> **Interview Q: How do you recognize a monotonic stack problem?**  
> Look for patterns like "next/previous greater/smaller element", "span of consecutive elements", or "areas/distances based on surrounding heights". If a brute force O(n²) solution naturally involves comparing each element with all others, a monotonic stack likely gives O(n).
