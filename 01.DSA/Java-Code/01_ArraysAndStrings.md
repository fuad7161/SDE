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

## 1. Two Sum / Three Sum

**Two Sum** — find indices of two numbers that add up to a target.  
**Approach:** HashMap to store `value → index`. For each element, check if `target - num` exists.

```java
// ── Two Sum ──
// Time: O(n), Space: O(n)
public int[] twoSum(int[] nums, int target) {
    Map<Integer, Integer> map = new HashMap<>();
    for (int i = 0; i < nums.length; i++) {
        int complement = target - nums[i];
        if (map.containsKey(complement)) {
            return new int[]{map.get(complement), i};
        }
        map.put(nums[i], i);
    }
    return new int[]{};
}
```

**Three Sum** — find all unique triplets that sum to zero.  
**Approach:** Sort + Two Pointers. Fix one element, then use left/right pointers for the remaining pair.

```java
// ── Three Sum ──
// Time: O(n²), Space: O(1) excluding output
public List<List<Integer>> threeSum(int[] nums) {
    List<List<Integer>> result = new ArrayList<>();
    Arrays.sort(nums);                          // sort first

    for (int i = 0; i < nums.length - 2; i++) {
        if (i > 0 && nums[i] == nums[i - 1]) continue;   // skip duplicates

        int left = i + 1, right = nums.length - 1;
        while (left < right) {
            int sum = nums[i] + nums[left] + nums[right];
            if (sum == 0) {
                result.add(Arrays.asList(nums[i], nums[left], nums[right]));
                while (left < right && nums[left] == nums[left + 1]) left++;
                while (left < right && nums[right] == nums[right - 1]) right--;
                left++;
                right--;
            } else if (sum < 0) {
                left++;
            } else {
                right--;
            }
        }
    }
    return result;
}
```

> **Interview Q: Why do we sort for 3Sum but not 2Sum?**  
> For 2Sum we need to return indices, so sorting would disrupt them — we use a HashMap instead. For 3Sum we only need values, so sorting enables the two-pointer technique to find pairs in O(n) per fixed element, and also makes duplicate skipping simple.

---

## 2. Sliding Window

**Use when:** finding a subarray/substring satisfying some condition (max/min sum, unique chars, etc.).

### Fixed Window — Maximum Sum Subarray of Size K

```java
// Time: O(n), Space: O(1)
public int maxSumSubarray(int[] nums, int k) {
    int windowSum = 0;
    for (int i = 0; i < k; i++) windowSum += nums[i];    // first window

    int maxSum = windowSum;
    for (int i = k; i < nums.length; i++) {
        windowSum += nums[i] - nums[i - k];               // slide: add new, remove old
        maxSum = Math.max(maxSum, windowSum);
    }
    return maxSum;
}
```

### Variable Window — Longest Substring Without Repeating Characters

```java
// Time: O(n), Space: O(min(n, charset))
public int lengthOfLongestSubstring(String s) {
    Map<Character, Integer> lastSeen = new HashMap<>();
    int maxLen = 0;
    int left = 0;

    for (int right = 0; right < s.length(); right++) {
        char c = s.charAt(right);
        if (lastSeen.containsKey(c) && lastSeen.get(c) >= left) {
            left = lastSeen.get(c) + 1;   // shrink window past the duplicate
        }
        lastSeen.put(c, right);
        maxLen = Math.max(maxLen, right - left + 1);
    }
    return maxLen;
}
```

### Variable Window — Minimum Window Substring

```java
// Time: O(n), Space: O(k) where k = charset size
public String minWindow(String s, String t) {
    Map<Character, Integer> need = new HashMap<>();
    for (char c : t.toCharArray()) need.merge(c, 1, Integer::sum);

    int left = 0, matched = 0;
    int minLen = Integer.MAX_VALUE, start = 0;

    for (int right = 0; right < s.length(); right++) {
        char c = s.charAt(right);
        if (need.containsKey(c)) {
            need.put(c, need.get(c) - 1);
            if (need.get(c) == 0) matched++;
        }
        while (matched == need.size()) {
            if (right - left + 1 < minLen) {
                minLen = right - left + 1;
                start = left;
            }
            char lc = s.charAt(left++);
            if (need.containsKey(lc)) {
                if (need.get(lc) == 0) matched--;
                need.put(lc, need.get(lc) + 1);
            }
        }
    }
    return minLen == Integer.MAX_VALUE ? "" : s.substring(start, start + minLen);
}
```

> **Interview Q: How do you decide window size — fixed vs variable?**  
> Use a **fixed window** when the problem gives a specific size `k`. Use a **variable window** when you need to expand/shrink based on a condition (e.g., all unique, sum ≤ target). The template is: expand `right` always, shrink `left` when the condition is violated.

---

## 3. Kadane's Algorithm

**Problem:** Find the contiguous subarray with the maximum sum.  
**Key Idea:** At each element, decide: extend previous subarray OR start fresh.

```java
// ── Basic Kadane's ──
// Time: O(n), Space: O(1)
public int maxSubArray(int[] nums) {
    int currentSum = nums[0];
    int maxSum = nums[0];

    for (int i = 1; i < nums.length; i++) {
        // Either extend existing subarray or start new one from current element
        currentSum = Math.max(nums[i], currentSum + nums[i]);
        maxSum = Math.max(maxSum, currentSum);
    }
    return maxSum;
}

// ── With subarray indices ──
public int[] maxSubArrayWithIndices(int[] nums) {
    int currentSum = nums[0], maxSum = nums[0];
    int start = 0, end = 0, tempStart = 0;

    for (int i = 1; i < nums.length; i++) {
        if (nums[i] > currentSum + nums[i]) {
            currentSum = nums[i];
            tempStart = i;
        } else {
            currentSum += nums[i];
        }
        if (currentSum > maxSum) {
            maxSum = currentSum;
            start = tempStart;
            end = i;
        }
    }
    return new int[]{maxSum, start, end};
}

// ── Maximum Product Subarray (variant) ──
public int maxProduct(int[] nums) {
    int maxProd = nums[0], minProd = nums[0], result = nums[0];
    for (int i = 1; i < nums.length; i++) {
        if (nums[i] < 0) { int tmp = maxProd; maxProd = minProd; minProd = tmp; }
        maxProd = Math.max(nums[i], maxProd * nums[i]);
        minProd = Math.min(nums[i], minProd * nums[i]);
        result = Math.max(result, maxProd);
    }
    return result;
}
```

> **Interview Q: Why track both max and min in the product variant?**  
> A negative number can flip the minimum product to the maximum. So we always track both — when we see a negative number, we swap max and min before updating.

---

## 4. Rotate Array / Matrix

### Rotate Array by K Steps (right rotation)

```java
// Time: O(n), Space: O(1) — reverse trick
public void rotate(int[] nums, int k) {
    int n = nums.length;
    k %= n;                     // handle k > n
    reverse(nums, 0, n - 1);    // reverse entire array
    reverse(nums, 0, k - 1);    // reverse first k elements
    reverse(nums, k, n - 1);    // reverse remaining
}

private void reverse(int[] nums, int left, int right) {
    while (left < right) {
        int tmp = nums[left];
        nums[left++] = nums[right];
        nums[right--] = tmp;
    }
}
// Input:  [1,2,3,4,5,6,7], k=3
// Step 1: [7,6,5,4,3,2,1]
// Step 2: [5,6,7,4,3,2,1]
// Step 3: [5,6,7,1,2,3,4] ✓
```

### Rotate Matrix 90 Degrees (clockwise)

```java
// Time: O(n²), Space: O(1)
public void rotate(int[][] matrix) {
    int n = matrix.length;
    // Step 1: Transpose (swap matrix[i][j] with matrix[j][i])
    for (int i = 0; i < n; i++) {
        for (int j = i + 1; j < n; j++) {
            int tmp = matrix[i][j];
            matrix[i][j] = matrix[j][i];
            matrix[j][i] = tmp;
        }
    }
    // Step 2: Reverse each row
    for (int i = 0; i < n; i++) {
        int left = 0, right = n - 1;
        while (left < right) {
            int tmp = matrix[i][left];
            matrix[i][left++] = matrix[i][right];
            matrix[i][right--] = tmp;
        }
    }
}
// Counter-clockwise: reverse each row first, then transpose
```

---

## 5. Merge Intervals

**Problem:** Given a list of intervals, merge all overlapping intervals.  
**Approach:** Sort by start time, then iterate and merge if current overlaps with last merged.

```java
// Time: O(n log n), Space: O(n)
public int[][] merge(int[][] intervals) {
    Arrays.sort(intervals, (a, b) -> a[0] - b[0]);   // sort by start

    List<int[]> merged = new ArrayList<>();
    int[] current = intervals[0];

    for (int i = 1; i < intervals.length; i++) {
        if (intervals[i][0] <= current[1]) {
            // overlapping — extend the end if needed
            current[1] = Math.max(current[1], intervals[i][1]);
        } else {
            // no overlap — save current, move to next
            merged.add(current);
            current = intervals[i];
        }
    }
    merged.add(current);
    return merged.toArray(new int[0][]);
}

// ── Insert Interval (variant) ──
public int[][] insert(int[][] intervals, int[] newInterval) {
    List<int[]> result = new ArrayList<>();
    int i = 0, n = intervals.length;

    // Add all intervals ending before newInterval starts
    while (i < n && intervals[i][1] < newInterval[0]) result.add(intervals[i++]);

    // Merge overlapping intervals
    while (i < n && intervals[i][0] <= newInterval[1]) {
        newInterval[0] = Math.min(newInterval[0], intervals[i][0]);
        newInterval[1] = Math.max(newInterval[1], intervals[i][1]);
        i++;
    }
    result.add(newInterval);

    // Add remaining
    while (i < n) result.add(intervals[i++]);
    return result.toArray(new int[0][]);
}
```

---

## 6. Trapping Rain Water

**Problem:** Given heights of bars, find total water trapped after rain.  
**Approach:** Two-pointer. Water at any bar = `min(maxLeft, maxRight) - height[i]`.

```java
// Time: O(n), Space: O(1)
public int trap(int[] height) {
    int left = 0, right = height.length - 1;
    int maxLeft = 0, maxRight = 0;
    int water = 0;

    while (left < right) {
        if (height[left] < height[right]) {
            if (height[left] >= maxLeft) {
                maxLeft = height[left];
            } else {
                water += maxLeft - height[left];   // trapped water on left side
            }
            left++;
        } else {
            if (height[right] >= maxRight) {
                maxRight = height[right];
            } else {
                water += maxRight - height[right]; // trapped water on right side
            }
            right--;
        }
    }
    return water;
}

// ── Prefix/Suffix array approach (easier to understand) ──
public int trapPrefixSuffix(int[] height) {
    int n = height.length;
    int[] maxLeft = new int[n], maxRight = new int[n];

    maxLeft[0] = height[0];
    for (int i = 1; i < n; i++) maxLeft[i] = Math.max(maxLeft[i-1], height[i]);

    maxRight[n-1] = height[n-1];
    for (int i = n-2; i >= 0; i--) maxRight[i] = Math.max(maxRight[i+1], height[i]);

    int water = 0;
    for (int i = 0; i < n; i++) water += Math.min(maxLeft[i], maxRight[i]) - height[i];
    return water;
}
```

---

## 7. Product of Array Except Self

**Problem:** Return array where `output[i]` = product of all elements except `nums[i]`, **without division**.

```java
// Time: O(n), Space: O(1) output array only
public int[] productExceptSelf(int[] nums) {
    int n = nums.length;
    int[] result = new int[n];

    // Left pass: result[i] = product of all elements to the LEFT of i
    result[0] = 1;
    for (int i = 1; i < n; i++) {
        result[i] = result[i - 1] * nums[i - 1];
    }

    // Right pass: multiply by product of all elements to the RIGHT of i
    int rightProduct = 1;
    for (int i = n - 1; i >= 0; i--) {
        result[i] *= rightProduct;
        rightProduct *= nums[i];
    }

    return result;
}
// nums:     [1, 2, 3, 4]
// left:     [1, 1, 2, 6]
// right:    [24,12,4, 1]
// result:   [24,12,8, 6]
```

> **Interview Q: Why can't we use division?**  
> If the array contains a zero, division by zero is undefined. Even without zeros, division is conceptually cheating — the intended solution tests prefix/suffix product understanding.

---

## 8. Longest Common Prefix

**Problem:** Find the longest common prefix string among an array of strings.

```java
// ── Vertical scanning ──
// Time: O(S) where S = total chars, Space: O(1)
public String longestCommonPrefix(String[] strs) {
    if (strs == null || strs.length == 0) return "";

    for (int i = 0; i < strs[0].length(); i++) {
        char c = strs[0].charAt(i);
        for (int j = 1; j < strs.length; j++) {
            // stop if index out of bounds OR character doesn't match
            if (i >= strs[j].length() || strs[j].charAt(i) != c) {
                return strs[0].substring(0, i);
            }
        }
    }
    return strs[0];
}

// ── Sort-based approach ──
// After sorting, only compare first and last strings
public String longestCommonPrefixSort(String[] strs) {
    Arrays.sort(strs);
    String first = strs[0], last = strs[strs.length - 1];
    int i = 0;
    while (i < first.length() && first.charAt(i) == last.charAt(i)) i++;
    return first.substring(0, i);
}
```

---

## 9. Valid Anagram / Group Anagrams

### Valid Anagram

```java
// Two strings are anagrams if they have the same character frequencies
// Time: O(n), Space: O(1) — fixed 26-char alphabet
public boolean isAnagram(String s, String t) {
    if (s.length() != t.length()) return false;
    int[] count = new int[26];
    for (char c : s.toCharArray()) count[c - 'a']++;
    for (char c : t.toCharArray()) count[c - 'a']--;
    for (int n : count) if (n != 0) return false;
    return true;
}
```

### Group Anagrams

```java
// Group strings that are anagrams of each other
// Time: O(n * k log k) where k = max string length
public List<List<String>> groupAnagrams(String[] strs) {
    Map<String, List<String>> map = new HashMap<>();
    for (String s : strs) {
        char[] chars = s.toCharArray();
        Arrays.sort(chars);                         // canonical key
        String key = new String(chars);
        map.computeIfAbsent(key, k -> new ArrayList<>()).add(s);
    }
    return new ArrayList<>(map.values());
}

// ── O(n*k) variant using frequency count as key ──
public List<List<String>> groupAnagramsFast(String[] strs) {
    Map<String, List<String>> map = new HashMap<>();
    for (String s : strs) {
        int[] count = new int[26];
        for (char c : s.toCharArray()) count[c - 'a']++;
        String key = Arrays.toString(count);        // "[1,0,0,...,1,0,...]"
        map.computeIfAbsent(key, k -> new ArrayList<>()).add(s);
    }
    return new ArrayList<>(map.values());
}
```

> **Interview Q: What is the time complexity of Group Anagrams with sorting vs frequency counting?**  
> Sorting-based: `O(n * k log k)` where `k` is the average string length. Frequency-based: `O(n * k)` since we scan each character once and build a fixed-length key. For large inputs with long strings, frequency counting is faster.
