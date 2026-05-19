# Hashing

> Hash maps provide O(1) average-case lookup, insertion, and deletion. In Go, the built-in `map` type is your hash table.

---

## Table of Contents

1. [Map Basics in Go](#1-map-basics-in-go)
2. [Two Sum](#2-two-sum)
3. [Subarray Sum Equals K](#3-subarray-sum-equals-k)
4. [Longest Consecutive Sequence](#4-longest-consecutive-sequence)
5. [4Sum II](#5-4sum-ii)

---

## 1. Map Basics in Go

```go
// Create
m := make(map[string]int)
m["apple"] = 5

// Read with existence check
val, ok := m["banana"]
if !ok {
    fmt.Println("key not found")
}

// Delete
delete(m, "apple")

// Iterate
for key, value := range m {
    fmt.Println(key, value)
}

// Zero value: accessing a missing key returns 0 (for int), "" (for string), etc.
m["newKey"]++  // safe even if "newKey" doesn't exist yet

// Map as a set
set := make(map[int]bool)
set[42] = true
if set[42] {
    fmt.Println("exists")
}

// Struct as a map key (must be comparable — no slices/maps as fields)
type Point struct{ X, Y int }
grid := make(map[Point]bool)
grid[Point{1, 2}] = true
```

---

## 2. Two Sum

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

// Two Sum in sorted array — O(n) time, O(1) space (two pointers)
func twoSumSorted(numbers []int, target int) []int {
    left, right := 0, len(numbers)-1
    for left < right {
        sum := numbers[left] + numbers[right]
        switch {
        case sum == target:
            return []int{left + 1, right + 1} // 1-indexed
        case sum < target:
            left++
        default:
            right--
        }
    }
    return nil
}
```

---

## 3. Subarray Sum Equals K

```go
// Time: O(n), Space: O(n)
// Count subarrays whose sum equals k
func subarraySum(nums []int, k int) int {
    prefixCount := map[int]int{0: 1} // prefix sum → count
    count, currentSum := 0, 0

    for _, num := range nums {
        currentSum += num
        // If (currentSum - k) exists, those subarrays end here
        count += prefixCount[currentSum-k]
        prefixCount[currentSum]++
    }
    return count
}
```

> **Interview Q: Why initialize `prefixCount[0] = 1`?**  
> This handles subarrays that start from index 0 (where the prefix sum itself equals k). Without it, we'd miss subarrays like `[k]` at the beginning.

---

## 4. Longest Consecutive Sequence

```go
// Time: O(n), Space: O(n)
func longestConsecutive(nums []int) int {
    numSet := make(map[int]bool)
    for _, n := range nums { numSet[n] = true }

    longest := 0
    for num := range numSet {
        // Only start counting from the beginning of a sequence
        if !numSet[num-1] {
            cur := num
            streak := 1
            for numSet[cur+1] {
                cur++
                streak++
            }
            if streak > longest { longest = streak }
        }
    }
    return longest
}
```

> **Interview Q: Why check `!numSet[num-1]` before starting a streak?**  
> It ensures we only begin counting from the smallest element of a consecutive run. Without it, every element would start a count, making the algorithm O(n²) instead of O(n).

---

## 5. 4Sum II

```go
// Time: O(n²), Space: O(n²)
// Count tuples (i,j,k,l) such that A[i]+B[j]+C[k]+D[l] == 0
func fourSumCount(nums1 []int, nums2 []int, nums3 []int, nums4 []int) int {
    pairSums := make(map[int]int)
    for _, a := range nums1 {
        for _, b := range nums2 {
            pairSums[a+b]++
        }
    }

    count := 0
    for _, c := range nums3 {
        for _, d := range nums4 {
            count += pairSums[-(c + d)]
        }
    }
    return count
}

// ── 4Sum — find all unique quadruplets summing to target ──
// Time: O(n³), Space: O(1) excluding output
import "sort"

func fourSum(nums []int, target int) [][]int {
    sort.Ints(nums)
    result := [][]int{}
    n := len(nums)

    for i := 0; i < n-3; i++ {
        if i > 0 && nums[i] == nums[i-1] { continue }
        for j := i + 1; j < n-2; j++ {
            if j > i+1 && nums[j] == nums[j-1] { continue }
            left, right := j+1, n-1
            for left < right {
                sum := nums[i] + nums[j] + nums[left] + nums[right]
                switch {
                case sum == target:
                    result = append(result, []int{nums[i], nums[j], nums[left], nums[right]})
                    for left < right && nums[left] == nums[left+1] { left++ }
                    for left < right && nums[right] == nums[right-1] { right-- }
                    left++
                    right--
                case sum < target:
                    left++
                default:
                    right--
                }
            }
        }
    }
    return result
}
```

---

## Hashing Patterns Summary

| Pattern | Key Idea | Example |
|---|---|---|
| Complement lookup | Store seen values; check `target - current` | Two Sum |
| Prefix sum + map | `count += map[prefixSum - k]` | Subarray Sum = K |
| Frequency count | `map[item]++` then query | Top K Frequent |
| Seen-set dedup | `map[val]bool` for O(1) membership | Longest Consecutive |
| Pair sums | Build AB sums in map, query with -(C+D) | 4Sum II |
