# Binary Search

> Binary search is not just for sorted arrays. Any time you can **define a monotonic condition**, you can binary search on the answer.

---

## Table of Contents

1. [Classic Binary Search](#1-classic-binary-search)
2. [Search in Rotated Sorted Array](#2-search-in-rotated-sorted-array)
3. [Find Peak Element](#3-find-peak-element)
4. [Kth Smallest in a Sorted Matrix](#4-kth-smallest-in-a-sorted-matrix)
5. [Binary Search on the Answer](#5-binary-search-on-the-answer)
6. [Lower Bound & Upper Bound](#6-lower-bound--upper-bound)

---

## 1. Classic Binary Search

```go
// Time: O(log n), Space: O(1)
func search(nums []int, target int) int {
    left, right := 0, len(nums)-1
    for left <= right {
        mid := left + (right-left)/2 // avoids overflow
        switch {
        case nums[mid] == target:
            return mid
        case nums[mid] < target:
            left = mid + 1
        default:
            right = mid - 1
        }
    }
    return -1
}
```

> **Interview Q: Why `left + (right-left)/2` instead of `(left+right)/2`?**  
> `left + right` can overflow `int` for large indices. The subtraction form is always safe.

---

## 2. Search in Rotated Sorted Array

```go
// Time: O(log n) — one sorted half always exists
func searchRotated(nums []int, target int) int {
    left, right := 0, len(nums)-1
    for left <= right {
        mid := left + (right-left)/2
        if nums[mid] == target { return mid }

        // Determine which half is sorted
        if nums[left] <= nums[mid] { // left half is sorted
            if nums[left] <= target && target < nums[mid] {
                right = mid - 1
            } else {
                left = mid + 1
            }
        } else { // right half is sorted
            if nums[mid] < target && target <= nums[right] {
                left = mid + 1
            } else {
                right = mid - 1
            }
        }
    }
    return -1
}

// Find minimum in rotated sorted array (no duplicates)
func findMin(nums []int) int {
    left, right := 0, len(nums)-1
    for left < right {
        mid := left + (right-left)/2
        if nums[mid] > nums[right] {
            left = mid + 1 // min is in right half
        } else {
            right = mid // mid might be the min
        }
    }
    return nums[left]
}
```

---

## 3. Find Peak Element

```go
// A peak is any element greater than its neighbors
// Time: O(log n)
func findPeakElement(nums []int) int {
    left, right := 0, len(nums)-1
    for left < right {
        mid := left + (right-left)/2
        if nums[mid] > nums[mid+1] {
            right = mid // peak is on the left (including mid)
        } else {
            left = mid + 1 // peak is on the right
        }
    }
    return left
}
```

---

## 4. Kth Smallest in a Sorted Matrix

```go
// Binary search on value — Time: O(n log(max-min))
func kthSmallest(matrix [][]int, k int) int {
    n := len(matrix)
    left, right := matrix[0][0], matrix[n-1][n-1]

    // count elements <= mid using staircase traversal — O(n)
    countLE := func(mid int) int {
        count := 0
        row, col := n-1, 0
        for row >= 0 && col < n {
            if matrix[row][col] <= mid {
                count += row + 1
                col++
            } else {
                row--
            }
        }
        return count
    }

    for left < right {
        mid := left + (right-left)/2
        if countLE(mid) < k {
            left = mid + 1
        } else {
            right = mid
        }
    }
    return left
}
```

---

## 5. Binary Search on the Answer

### Split Array Largest Sum (Minimize the maximum sum of sub-arrays)

```go
// "Is it possible to split into m parts, each with sum ≤ limit?"
// Time: O(n log(sum))
func splitArray(nums []int, m int) int {
    left, right := 0, 0
    for _, n := range nums {
        if n > left { left = n }
        right += n
    }

    canSplit := func(limit int) bool {
        parts, cur := 1, 0
        for _, n := range nums {
            if cur+n > limit {
                parts++
                cur = 0
            }
            cur += n
        }
        return parts <= m
    }

    for left < right {
        mid := left + (right-left)/2
        if canSplit(mid) {
            right = mid
        } else {
            left = mid + 1
        }
    }
    return left
}
```

### Koko Eating Bananas

```go
// Time: O(n log(max))
import "math"

func minEatingSpeed(piles []int, h int) int {
    left, right := 1, 0
    for _, p := range piles {
        if p > right { right = p }
    }

    canFinish := func(speed int) bool {
        hours := 0
        for _, p := range piles {
            hours += (p + speed - 1) / speed // ceiling division
        }
        return hours <= h
    }

    for left < right {
        mid := left + (right-left)/2
        if canFinish(mid) {
            right = mid
        } else {
            left = mid + 1
        }
    }
    return left
}
```

> **Binary Search on Answer Template:**
> ```go
> left, right := minPossible, maxPossible
> for left < right {
>     mid := left + (right-left)/2
>     if feasible(mid) {
>         right = mid     // minimize: tighten from above
>         // left = mid+1 // maximize: tighten from below
>     } else {
>         left = mid + 1  // minimize: push up
>         // right = mid  // maximize: push down
>     }
> }
> return left
> ```

---

## 6. Lower Bound & Upper Bound

```go
// Lower bound — first index where nums[i] >= target
func lowerBound(nums []int, target int) int {
    left, right := 0, len(nums)
    for left < right {
        mid := left + (right-left)/2
        if nums[mid] < target {
            left = mid + 1
        } else {
            right = mid
        }
    }
    return left
}

// Upper bound — first index where nums[i] > target
func upperBound(nums []int, target int) int {
    left, right := 0, len(nums)
    for left < right {
        mid := left + (right-left)/2
        if nums[mid] <= target {
            left = mid + 1
        } else {
            right = mid
        }
    }
    return left
}

// Search Range [first, last] occurrence of target
func searchRange(nums []int, target int) [2]int {
    first := lowerBound(nums, target)
    if first == len(nums) || nums[first] != target {
        return [2]int{-1, -1}
    }
    last := upperBound(nums, target) - 1
    return [2]int{first, last}
}

// Go standard library equivalents
import "sort"
// sort.SearchInts(nums, target) → lower bound (first index >= target)
```

> **Interview Q: When does `left < right` vs `left <= right` apply?**  
> Use `left <= right` when you return inside the loop (classic search for exact match). Use `left < right` when you converge the window (lower/upper bound and BS on answer) — the loop exits with `left == right` being the answer.
