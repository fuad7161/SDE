# Hashing

> HashMaps and HashSets give O(1) average lookup, insert, and delete. Essential for frequency counting, prefix sums, and deduplication.

---

## Table of Contents

1. [Two Sum — Classic](#1-two-sum--classic)
2. [Subarray Sum Equals K](#2-subarray-sum-equals-k)
3. [Longest Consecutive Sequence](#3-longest-consecutive-sequence)
4. [4Sum / Pair with Target](#4-4sum--pair-with-target)

---

## HashMap Quick Reference

```java
Map<Integer, Integer> map = new HashMap<>();
map.put(key, value);
map.get(key);                          // null if not found
map.getOrDefault(key, 0);             // safe default
map.containsKey(key);
map.merge(key, 1, Integer::sum);      // increment count
map.computeIfAbsent(key, k -> new ArrayList<>()).add(val);

// Iterate
for (Map.Entry<Integer, Integer> entry : map.entrySet()) {
    entry.getKey(); entry.getValue();
}

Set<Integer> set = new HashSet<>();
set.add(x); set.contains(x); set.remove(x);
```

---

## 1. Two Sum — Classic

```java
// Find two indices such that nums[i] + nums[j] == target
// Time: O(n), Space: O(n)
public int[] twoSum(int[] nums, int target) {
    Map<Integer, Integer> map = new HashMap<>();  // value → index

    for (int i = 0; i < nums.length; i++) {
        int complement = target - nums[i];
        if (map.containsKey(complement)) {
            return new int[]{map.get(complement), i};
        }
        map.put(nums[i], i);
    }
    return new int[]{};
}
// nums=[2,7,11,15], target=9 → [0,1] (2+7=9)

// ── Two Sum II — Sorted array (use two pointers, O(1) space) ──
public int[] twoSumSorted(int[] nums, int target) {
    int lo = 0, hi = nums.length - 1;
    while (lo < hi) {
        int sum = nums[lo] + nums[hi];
        if (sum == target)      return new int[]{lo + 1, hi + 1};  // 1-indexed
        else if (sum < target)  lo++;
        else                    hi--;
    }
    return new int[]{};
}
```

---

## 2. Subarray Sum Equals K

**Problem:** Count subarrays whose sum equals `k`.  
**Key insight:** `sum[i..j] = prefixSum[j] - prefixSum[i-1]`. So if `prefixSum[j] - k` was seen before, we found a valid subarray.

```java
// Time: O(n), Space: O(n)
public int subarraySum(int[] nums, int k) {
    Map<Integer, Integer> prefixCount = new HashMap<>();
    prefixCount.put(0, 1);   // empty prefix (sum 0 seen once)

    int count = 0, sum = 0;
    for (int num : nums) {
        sum += num;
        // If (sum - k) was seen before, there's a subarray ending here with sum = k
        count += prefixCount.getOrDefault(sum - k, 0);
        prefixCount.merge(sum, 1, Integer::sum);
    }
    return count;
}
// nums=[1,1,1], k=2 → count=2 (indices [0,1] and [1,2])

// ── Maximum length subarray with sum = k ──
public int maxSubarrayLen(int[] nums, int k) {
    Map<Integer, Integer> firstSeen = new HashMap<>();
    firstSeen.put(0, -1);   // prefix sum 0 first seen at index -1
    int maxLen = 0, sum = 0;

    for (int i = 0; i < nums.length; i++) {
        sum += nums[i];
        if (firstSeen.containsKey(sum - k)) {
            maxLen = Math.max(maxLen, i - firstSeen.get(sum - k));
        }
        firstSeen.putIfAbsent(sum, i);   // only store first occurrence
    }
    return maxLen;
}
```

> **Interview Q: Why `prefixCount.put(0, 1)` at the start?**  
> This handles the case where a subarray starting from index 0 sums to `k`. Without it, `sum - k = 0` would not be found in the map, missing those subarrays.

---

## 3. Longest Consecutive Sequence

**Problem:** Given unsorted array, find the length of the longest consecutive sequence. Requires O(n) time.

```java
// Time: O(n), Space: O(n)
// Key: only start counting from the beginning of a sequence
public int longestConsecutive(int[] nums) {
    Set<Integer> set = new HashSet<>();
    for (int num : nums) set.add(num);

    int maxLen = 0;
    for (int num : set) {
        // Only start counting if num is the beginning of a sequence
        if (!set.contains(num - 1)) {
            int current = num;
            int length = 1;
            while (set.contains(current + 1)) {
                current++;
                length++;
            }
            maxLen = Math.max(maxLen, length);
        }
    }
    return maxLen;
}
// nums=[100,4,200,1,3,2] → longest=[1,2,3,4], length=4

// Why O(n)? Each number is visited at most twice —
// once in the outer loop, once in the inner while loop.
// The !set.contains(num-1) check ensures the while loop
// runs only for sequence starts.
```

---

## 4. 4Sum / Pair with Target

### 4Sum — Find all unique quadruplets summing to target

```java
// Sort + Two Pointer + Two outer loops
// Time: O(n³), Space: O(1) excluding output
public List<List<Integer>> fourSum(int[] nums, int target) {
    Arrays.sort(nums);
    List<List<Integer>> result = new ArrayList<>();
    int n = nums.length;

    for (int i = 0; i < n - 3; i++) {
        if (i > 0 && nums[i] == nums[i-1]) continue;   // skip outer duplicate

        for (int j = i + 1; j < n - 2; j++) {
            if (j > i + 1 && nums[j] == nums[j-1]) continue;  // skip inner duplicate

            int left = j + 1, right = n - 1;
            while (left < right) {
                long sum = (long) nums[i] + nums[j] + nums[left] + nums[right];
                if (sum == target) {
                    result.add(Arrays.asList(nums[i], nums[j], nums[left], nums[right]));
                    while (left < right && nums[left]  == nums[left+1])  left++;
                    while (left < right && nums[right] == nums[right-1]) right--;
                    left++; right--;
                } else if (sum < target) {
                    left++;
                } else {
                    right--;
                }
            }
        }
    }
    return result;
}
```

### Count Pairs with Given Difference

```java
// Count pairs (i, j) where nums[j] - nums[i] == k (i < j), k > 0
// Time: O(n), Space: O(n)
public int countPairsWithDiff(int[] nums, int k) {
    Map<Integer, Integer> freq = new HashMap<>();
    int count = 0;
    for (int num : nums) {
        count += freq.getOrDefault(num - k, 0);  // look for num - k (already seen)
        freq.merge(num, 1, Integer::sum);
    }
    return count;
}

// ── Two Sum with HashMap — all pairs summing to target ──
public List<int[]> allPairs(int[] nums, int target) {
    List<int[]> result = new ArrayList<>();
    Set<Integer> seen = new HashSet<>();
    Set<Integer> used = new HashSet<>();    // avoid duplicate pairs

    for (int num : nums) {
        int complement = target - num;
        if (seen.contains(complement) && !used.contains(num)) {
            result.add(new int[]{complement, num});
            used.add(num);
            used.add(complement);
        }
        seen.add(num);
    }
    return result;
}
```
