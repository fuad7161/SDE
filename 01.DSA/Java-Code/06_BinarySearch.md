# Binary Search

> Binary search is not just for sorted arrays — apply it whenever the search space can be halved. Master the three templates.

---

## Table of Contents

1. [Classic Binary Search](#1-classic-binary-search)
2. [Search in Rotated Sorted Array](#2-search-in-rotated-sorted-array)
3. [Find Peak Element](#3-find-peak-element)
4. [Kth Smallest in Matrix](#4-kth-smallest-in-matrix)
5. [Binary Search on Answer](#5-binary-search-on-answer)
6. [Lower / Upper Bound](#6-lower--upper-bound)

---

## Binary Search Templates

```
Template 1 — Exact match:        while (lo <= hi)
Template 2 — Left boundary:      while (lo < hi),  lo = mid + 1 or hi = mid
Template 3 — Right boundary:     while (lo < hi),  lo = mid + 1 or hi = mid - 1
```

**Avoid overflow:** Use `mid = lo + (hi - lo) / 2` instead of `(lo + hi) / 2`.

---

## 1. Classic Binary Search

```java
// Find target in sorted array. Return index or -1.
// Time: O(log n), Space: O(1)
public int search(int[] nums, int target) {
    int lo = 0, hi = nums.length - 1;

    while (lo <= hi) {
        int mid = lo + (hi - lo) / 2;  // prevents overflow
        if (nums[mid] == target) return mid;
        else if (nums[mid] < target)   lo = mid + 1;
        else                           hi = mid - 1;
    }
    return -1;
}

// ── Recursive version ──
public int searchRecursive(int[] nums, int target, int lo, int hi) {
    if (lo > hi) return -1;
    int mid = lo + (hi - lo) / 2;
    if (nums[mid] == target) return mid;
    if (nums[mid] < target)  return searchRecursive(nums, target, mid + 1, hi);
    return searchRecursive(nums, target, lo, mid - 1);
}
```

---

## 2. Search in Rotated Sorted Array

**Key insight:** Even in a rotated array, **one half is always sorted**. Use that to eliminate halves.

```java
// Time: O(log n), Space: O(1)
public int searchRotated(int[] nums, int target) {
    int lo = 0, hi = nums.length - 1;

    while (lo <= hi) {
        int mid = lo + (hi - lo) / 2;
        if (nums[mid] == target) return mid;

        // Left half is sorted
        if (nums[lo] <= nums[mid]) {
            if (nums[lo] <= target && target < nums[mid]) hi = mid - 1;
            else lo = mid + 1;
        }
        // Right half is sorted
        else {
            if (nums[mid] < target && target <= nums[hi]) lo = mid + 1;
            else hi = mid - 1;
        }
    }
    return -1;
}

// ── Find minimum in rotated sorted array ──
public int findMin(int[] nums) {
    int lo = 0, hi = nums.length - 1;
    while (lo < hi) {
        int mid = lo + (hi - lo) / 2;
        if (nums[mid] > nums[hi]) lo = mid + 1;  // min is in right half
        else                      hi = mid;        // min is in left half (mid could be min)
    }
    return nums[lo];
}
```

---

## 3. Find Peak Element

**Peak:** An element greater than its neighbors. Return any peak index.  
**Key insight:** If `nums[mid] < nums[mid+1]`, a peak exists to the right.

```java
// Time: O(log n), Space: O(1)
public int findPeakElement(int[] nums) {
    int lo = 0, hi = nums.length - 1;

    while (lo < hi) {
        int mid = lo + (hi - lo) / 2;
        if (nums[mid] < nums[mid + 1]) {
            lo = mid + 1;   // ascending slope → peak is to the right
        } else {
            hi = mid;       // descending slope → peak is at mid or to the left
        }
    }
    return lo;   // lo == hi, this is a peak
}
// nums = [1, 2, 3, 1] → peak at index 2 (value 3)
// nums = [1, 2, 1, 3, 5, 6, 4] → peak at index 1 or 5 (either valid)
```

---

## 4. Kth Smallest in Matrix

**Matrix:** n×n sorted row-wise and column-wise. Find the kth smallest element.

```java
// Binary search on value space (not index)
// Time: O(n log(max-min)), Space: O(1)
public int kthSmallest(int[][] matrix, int k) {
    int n = matrix.length;
    int lo = matrix[0][0], hi = matrix[n-1][n-1];

    while (lo < hi) {
        int mid = lo + (hi - lo) / 2;
        int count = countLessOrEqual(matrix, mid, n);

        if (count < k) lo = mid + 1;
        else           hi = mid;
    }
    return lo;
}

// Count elements <= mid using staircase traversal from bottom-left
private int countLessOrEqual(int[][] matrix, int target, int n) {
    int count = 0, row = n - 1, col = 0;
    while (row >= 0 && col < n) {
        if (matrix[row][col] <= target) {
            count += row + 1;   // all elements above in this column are also <= target
            col++;
        } else {
            row--;
        }
    }
    return count;
}
```

---

## 5. Binary Search on Answer

**Pattern:** When asking "find the minimum/maximum X such that condition holds", binary search on the answer space, not an array.

### Capacity to Ship Packages Within D Days

```java
// Find minimum weight capacity so all packages ship within D days
// Time: O(n log(sum)), Space: O(1)
public int shipWithinDays(int[] weights, int days) {
    int lo = Arrays.stream(weights).max().getAsInt();    // min capacity = heaviest package
    int hi = Arrays.stream(weights).sum();               // max capacity = ship all in 1 day

    while (lo < hi) {
        int mid = lo + (hi - lo) / 2;
        if (canShip(weights, days, mid)) hi = mid;   // try smaller
        else                            lo = mid + 1; // need more capacity
    }
    return lo;
}

private boolean canShip(int[] weights, int days, int capacity) {
    int daysNeeded = 1, currentLoad = 0;
    for (int w : weights) {
        if (currentLoad + w > capacity) {
            daysNeeded++;
            currentLoad = 0;
        }
        currentLoad += w;
    }
    return daysNeeded <= days;
}
```

### Koko Eating Bananas

```java
// Find minimum speed k so Koko eats all bananas within h hours
public int minEatingSpeed(int[] piles, int h) {
    int lo = 1, hi = Arrays.stream(piles).max().getAsInt();

    while (lo < hi) {
        int mid = lo + (hi - lo) / 2;
        long hours = 0;
        for (int pile : piles) hours += (pile + mid - 1) / mid;  // ceiling division

        if (hours <= h) hi = mid;
        else            lo = mid + 1;
    }
    return lo;
}
```

> **Interview Q: How do you identify a "binary search on answer" problem?**  
> Look for: (1) asking for minimum or maximum of some variable, (2) a monotonic feasibility condition — if X works, then X+1 (or X-1) also works, (3) the answer lies in a clearly defined range. Binary search that condition over the answer space.

---

## 6. Lower / Upper Bound

```java
// Lower bound — first index where nums[i] >= target
public int lowerBound(int[] nums, int target) {
    int lo = 0, hi = nums.length;
    while (lo < hi) {
        int mid = lo + (hi - lo) / 2;
        if (nums[mid] < target) lo = mid + 1;
        else                    hi = mid;
    }
    return lo;  // first position where element >= target
}

// Upper bound — first index where nums[i] > target
public int upperBound(int[] nums, int target) {
    int lo = 0, hi = nums.length;
    while (lo < hi) {
        int mid = lo + (hi - lo) / 2;
        if (nums[mid] <= target) lo = mid + 1;
        else                     hi = mid;
    }
    return lo;  // first position where element > target
}

// Count of target occurrences = upperBound - lowerBound
public int countOccurrences(int[] nums, int target) {
    return upperBound(nums, target) - lowerBound(nums, target);
}

// Java built-in:
// Arrays.binarySearch(nums, target) returns index (may be negative if not found)
// Negative return means insertion point = -(result) - 1
```

> **Interview Q: What is the difference between lower bound and upper bound?**  
> **Lower bound** finds the first position where the element is `>= target` (leftmost occurrence). **Upper bound** finds the first position where the element is `> target` (one past the rightmost occurrence). Their difference gives the count of occurrences of `target`. This is the foundation of `std::lower_bound` in C++ and is commonly implemented manually in Java interviews.
