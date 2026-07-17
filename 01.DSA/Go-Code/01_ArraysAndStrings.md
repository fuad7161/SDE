# Arrays & Strings

> Most asked category in coding interviews. Focus on Two Pointer, Sliding Window, and prefix sums.

---

## Table of Contents

1. [Two Sum / Three Sum](#1-two-sum--three-sum)
2. [Sliding Window](#2-sliding-window)
3. [Kadane's Algorithm](#3-kadanes-algorithm)
4. [Rotate Array / Matrix](#4-rotate-array--matrix)
5. [Merge Intervals](#5-merge-intervals)
6. [Trapping Rain Water](#6-trapping-rain-water)
7. [Product of Array Except Self](#7-product-of-array-except-self)
8. [Longest Common Prefix](#8-longest-common-prefix)
9. [Valid Anagram / Group Anagrams](#9-valid-anagram--group-anagrams)

---

## Go Array/Slice Basics

```go
// Slice (dynamic array — used instead of fixed arrays)
nums := []int{1, 2, 3, 4, 5}
nums = append(nums, 6)
length := len(nums)

// 2D slice
matrix := [][]int{{1, 2}, {3, 4}}

// String iteration
s := "hello"
for i, ch := range s {   // ch is rune (Unicode code point)
    _ = i
    _ = ch
}
```

---

## 1. Two Sum / Three Sum

**Two Sum** — find indices of two numbers that add up to a target.

```go
// Time: O(n), Space: O(n)
func twoSum(nums []int, target int) []int {
    seen := make(map[int]int) // value → index
    for i, num := range nums {
        complement := target - num
        if j, ok := seen[complement]; ok {
            return []int{j, i}
        }
        seen[num] = i
    }
    return nil
}
```

**Three Sum** — find all unique triplets that sum to zero.

```go
// Time: O(n²), Space: O(1) excluding output
import "sort"

func threeSum(nums []int) [][]int {
    sort.Ints(nums)
    result := [][]int{}

    for i := 0; i < len(nums)-2; i++ {
        if i > 0 && nums[i] == nums[i-1] {
            continue // skip duplicates
        }
        left, right := i+1, len(nums)-1
        for left < right {
            sum := nums[i] + nums[left] + nums[right]
            switch {
            case sum == 0:
                result = append(result, []int{nums[i], nums[left], nums[right]})
                for left < right && nums[left] == nums[left+1] { left++ }
                for left < right && nums[right] == nums[right-1] { right-- }
                left++
                right--
            case sum < 0:
                left++
            default:
                right--
            }
        }
    }
    return result
}
```

> **Interview Q: Why sort for 3Sum but not 2Sum?**  
> For 2Sum we need indices — sorting destroys them, so we use a map. For 3Sum we only need values, so sorting enables the two-pointer technique and makes duplicate skipping trivial.

---

## 2. Sliding Window

### Fixed Window — Maximum Sum Subarray of Size K

```go
// Time: O(n), Space: O(1)
func maxSumSubarray(nums []int, k int) int {
    windowSum := 0
    for i := 0; i < k; i++ {
        windowSum += nums[i]
    }
    maxSum := windowSum
    for i := k; i < len(nums); i++ {
        windowSum += nums[i] - nums[i-k] // slide: add new, remove old
        if windowSum > maxSum {
            maxSum = windowSum
        }
    }
    return maxSum
}
```

### Variable Window — Longest Substring Without Repeating Characters

```go
// Time: O(n), Space: O(min(n, charset))
func lengthOfLongestSubstring(s string) int {
    lastSeen := make(map[byte]int)
    maxLen, left := 0, 0

    for right := 0; right < len(s); right++ {
        c := s[right]
        if idx, ok := lastSeen[c]; ok && idx >= left {
            left = idx + 1 // shrink window past the duplicate
        }
        lastSeen[c] = right
        if right-left+1 > maxLen {
            maxLen = right - left + 1
        }
    }
    return maxLen
}
```

### Variable Window — Minimum Window Substring

```go
// Time: O(n), Space: O(k)
func minWindow(s string, t string) string {
    need := make(map[byte]int)
    for i := 0; i < len(t); i++ {
        need[t[i]]++
    }

    left, matched := 0, 0
    minLen, start := len(s)+1, 0

    for right := 0; right < len(s); right++ {
        c := s[right]
        if _, ok := need[c]; ok {
            need[c]--
            if need[c] == 0 {
                matched++
            }
        }
        for matched == len(need) {
            if right-left+1 < minLen {
                minLen = right - left + 1
                start = left
            }
            lc := s[left]
            left++
            if _, ok := need[lc]; ok {
                if need[lc] == 0 {
                    matched--
                }
                need[lc]++
            }
        }
    }
    if minLen == len(s)+1 {
        return ""
    }
    return s[start : start+minLen]
}
```

> **Interview Q: Fixed vs variable window — how to decide?**  
> **Fixed window** when size `k` is given. **Variable window** when you expand `right` always and shrink `left` when a condition is violated (e.g., duplicate found, sum exceeded).

---

## 3. Kadane's Algorithm

```go
// Time: O(n), Space: O(1)
func maxSubArray(nums []int) int {
    currentSum := nums[0]
    maxSum := nums[0]

    for i := 1; i < len(nums); i++ {
        if currentSum+nums[i] > nums[i] {
            currentSum += nums[i]
        } else {
            currentSum = nums[i] // start fresh
        }
        if currentSum > maxSum {
            maxSum = currentSum
        }
    }
    return maxSum
}

// ── Maximum Product Subarray (variant) ──
func maxProduct(nums []int) int {
    maxP, minP, result := nums[0], nums[0], nums[0]
    for i := 1; i < len(nums); i++ {
        if nums[i] < 0 {
            maxP, minP = minP, maxP // swap on negative
        }
        if nums[i] > maxP*nums[i] { maxP = nums[i] } else { maxP *= nums[i] }
        if nums[i] < minP*nums[i] { minP = nums[i] } else { minP *= nums[i] }
        if maxP > result { result = maxP }
    }
    return result
}
```

> **Interview Q: Why track both max and min in the product variant?**  
> A negative number flips the minimum to maximum. Always track both — swap them when encountering a negative number.

---

## 4. Rotate Array / Matrix

### Rotate Array by K Steps (right rotation)

```go
// Time: O(n), Space: O(1) — reverse trick
func rotate(nums []int, k int) {
    n := len(nums)
    k %= n
    reverse(nums, 0, n-1)
    reverse(nums, 0, k-1)
    reverse(nums, k, n-1)
}

func reverse(nums []int, left, right int) {
    for left < right {
        nums[left], nums[right] = nums[right], nums[left]
        left++
        right--
    }
}
// [1,2,3,4,5,6,7], k=3 → [5,6,7,1,2,3,4]
```

### Rotate Matrix 90 Degrees Clockwise

```go
// Time: O(n²), Space: O(1)
func rotateMatrix(matrix [][]int) {
    n := len(matrix)
    // Step 1: Transpose
    for i := 0; i < n; i++ {
        for j := i + 1; j < n; j++ {
            matrix[i][j], matrix[j][i] = matrix[j][i], matrix[i][j]
        }
    }
    // Step 2: Reverse each row
    for i := 0; i < n; i++ {
        left, right := 0, n-1
        for left < right {
            matrix[i][left], matrix[i][right] = matrix[i][right], matrix[i][left]
            left++
            right--
        }
    }
}
```

---

## 5. Merge Intervals

```go
// Time: O(n log n), Space: O(n)
import "sort"

func merge(intervals [][]int) [][]int {
    sort.Slice(intervals, func(i, j int) bool {
        return intervals[i][0] < intervals[j][0]
    })

    merged := [][]int{intervals[0]}

    for i := 1; i < len(intervals); i++ {
        last := merged[len(merged)-1]
        if intervals[i][0] <= last[1] {
            // overlapping — extend end
            if intervals[i][1] > last[1] {
                last[1] = intervals[i][1]
            }
        } else {
            merged = append(merged, intervals[i])
        }
    }
    return merged
}
```

---

## 6. Trapping Rain Water

```go
// Two-pointer — Time: O(n), Space: O(1)
func trap(height []int) int {
    left, right := 0, len(height)-1
    maxLeft, maxRight := 0, 0
    water := 0

    for left < right {
        if height[left] < height[right] {
            if height[left] >= maxLeft {
                maxLeft = height[left]
            } else {
                water += maxLeft - height[left]
            }
            left++
        } else {
            if height[right] >= maxRight {
                maxRight = height[right]
            } else {
                water += maxRight - height[right]
            }
            right--
        }
    }
    return water
}
```

---

## 7. Product of Array Except Self

```go
// Time: O(n), Space: O(1) — output array only
func productExceptSelf(nums []int) []int {
    n := len(nums)
    result := make([]int, n)

    // Left pass: result[i] = product of all elements left of i
    result[0] = 1
    for i := 1; i < n; i++ {
        result[i] = result[i-1] * nums[i-1]
    }

    // Right pass: multiply by product of all elements right of i
    rightProduct := 1
    for i := n - 1; i >= 0; i-- {
        result[i] *= rightProduct
        rightProduct *= nums[i]
    }
    return result
}
// [1,2,3,4] → [24,12,8,6]
```

---

## 8. Longest Common Prefix

```go
// Vertical scanning — Time: O(S), Space: O(1)
func longestCommonPrefix(strs []string) string {
    if len(strs) == 0 {
        return ""
    }
    for i := 0; i < len(strs[0]); i++ {
        c := strs[0][i]
        for j := 1; j < len(strs); j++ {
            if i >= len(strs[j]) || strs[j][i] != c {
                return strs[0][:i]
            }
        }
    }
    return strs[0]
}
```

---

## 9. Valid Anagram / Group Anagrams

### Valid Anagram

```go
// Time: O(n), Space: O(1) — 26-char alphabet
func isAnagram(s string, t string) bool {
    if len(s) != len(t) {
        return false
    }
    count := [26]int{}
    for i := 0; i < len(s); i++ {
        count[s[i]-'a']++
        count[t[i]-'a']--
    }
    for _, v := range count {
        if v != 0 {
            return false
        }
    }
    return true
}
```

### Group Anagrams

```go
// Time: O(n * k log k), Space: O(n*k)
import "sort"

func groupAnagrams(strs []string) [][]string {
    groups := make(map[string][]string)
    for _, s := range strs {
        chars := []byte(s)
        sort.Slice(chars, func(i, j int) bool { return chars[i] < chars[j] })
        key := string(chars)
        groups[key] = append(groups[key], s)
    }
    result := make([][]string, 0, len(groups))
    for _, v := range groups {
        result = append(result, v)
    }
    return result
}

// ── O(n*k) — frequency count as key ──
import "fmt"

func groupAnagramsFast(strs []string) [][]string {
    groups := make(map[[26]int][]string)
    for _, s := range strs {
        var key [26]int
        for _, c := range s {
            key[c-'a']++
        }
        groups[key] = append(groups[key], s)
    }
    result := make([][]string, 0, len(groups))
    for _, v := range groups {
        result = append(result, v)
    }
    return result
}
```

> **Interview Q: What is the time complexity of Group Anagrams with sorting vs frequency counting?**  
> Sorting-based: `O(n * k log k)`. Frequency-based: `O(n * k)`. In Go, a `[26]int` array is a valid map key — it's comparable — making the frequency approach clean and fast.
