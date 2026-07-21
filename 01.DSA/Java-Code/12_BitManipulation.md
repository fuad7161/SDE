# Bit Manipulation

> Bit operations are O(1) and operate directly on binary representations. Essential for space-efficient algorithms and low-level tricks.

---

## Table of Contents

1. [Bit Operations Cheat Sheet](#bit-operations-cheat-sheet)
2. [Check Power of 2](#1-check-power-of-2)
3. [Count Set Bits](#2-count-set-bits)
4. [Single Number — XOR Trick](#3-single-number--xor-trick)
5. [Subsets Using Bitmask](#4-subsets-using-bitmask)
6. [Reverse Bits](#5-reverse-bits)

---

## Bit Operations Cheat Sheet

| Operation | Expression | Result |
|---|---|---|
| AND | `a & b` | 1 if both bits are 1 |
| OR | `a \| b` | 1 if at least one bit is 1 |
| XOR | `a ^ b` | 1 if bits differ |
| NOT | `~a` | flip all bits |
| Left shift | `a << n` | multiply by 2^n |
| Right shift | `a >> n` | divide by 2^n (signed) |
| Unsigned right | `a >>> n` | divide, fills 0 from left |
| Get bit i | `(a >> i) & 1` | extract bit at position i |
| Set bit i | `a \| (1 << i)` | set bit i to 1 |
| Clear bit i | `a & ~(1 << i)` | set bit i to 0 |
| Toggle bit i | `a ^ (1 << i)` | flip bit i |
| Clear lowest set bit | `n & (n - 1)` | removes rightmost 1 |
| Isolate lowest set bit | `n & (-n)` | keeps only rightmost 1 |
| XOR identities | `x ^ x = 0`, `x ^ 0 = x` | core XOR properties |

---

## 1. Check Power of 2

A power of 2 has exactly **one bit set** in binary. `n & (n-1)` removes the lowest set bit — if result is 0, only one bit was set.

```java
// Time: O(1), Space: O(1)
public boolean isPowerOfTwo(int n) {
    return n > 0 && (n & (n - 1)) == 0;
}
// n=8:  1000 & 0111 = 0000 → true
// n=6:  0110 & 0101 = 0100 → false (non-zero, not a power of 2)

// ── Power of 4 ──
// Power of 4 must be power of 2 AND the set bit must be at an even position
public boolean isPowerOfFour(int n) {
    // 0x55555555 = 01010101...01 — bits set at even positions
    return n > 0 && (n & (n-1)) == 0 && (n & 0x55555555) != 0;
}

// ── Power of 3 (no bit trick, but O(1)) ──
public boolean isPowerOfThree(int n) {
    // 3^19 = 1162261467 is the largest power of 3 within int range
    return n > 0 && 1162261467 % n == 0;
}
```

---

## 2. Count Set Bits

### Brian Kernighan's Algorithm

```java
// Each iteration removes the lowest set bit: n & (n-1)
// Iterations = number of set bits
// Time: O(number of set bits), Space: O(1)
public int countBits(int n) {
    int count = 0;
    while (n != 0) {
        n &= (n - 1);   // clear lowest set bit
        count++;
    }
    return count;
}
// n=13 (1101): 1101→1100→1000→0000  → 3 iterations → 3 set bits

// Java built-in:
// Integer.bitCount(n) — uses hardware popcount instruction
```

### Count Bits for All Numbers 0..n (DP)

```java
// dp[i] = dp[i >> 1] + (i & 1)
// Right shift drops the LSB. Remaining is dp[i/2]. Add 1 if LSB was set.
// Time: O(n), Space: O(n)
public int[] countBitsDP(int n) {
    int[] dp = new int[n + 1];
    for (int i = 1; i <= n; i++) {
        dp[i] = dp[i >> 1] + (i & 1);
    }
    return dp;
}
// n=5 → [0, 1, 1, 2, 1, 2]
```

---

## 3. Single Number — XOR Trick

### Single Number I — All elements appear twice except one

```java
// XOR: x ^ x = 0, x ^ 0 = x
// All pairs cancel out, leaving the single element
// Time: O(n), Space: O(1)
public int singleNumber(int[] nums) {
    int result = 0;
    for (int num : nums) result ^= num;
    return result;
}
// [4,1,2,1,2]: 4^1^2^1^2 = 4^(1^1)^(2^2) = 4^0^0 = 4
```

### Single Number II — All elements appear three times except one

```java
// Use two bitmasks: 'ones' stores bits seen once, 'twos' stores bits seen twice
// When a bit is seen 3 times, it's cleared from both
// Time: O(n), Space: O(1)
public int singleNumberII(int[] nums) {
    int ones = 0, twos = 0;
    for (int num : nums) {
        ones = (ones ^ num) & ~twos;
        twos = (twos ^ num) & ~ones;
    }
    return ones;
}
```

### Single Number III — Two elements appear once, rest appear twice

```java
// Step 1: XOR all → gets xor of the two unique numbers (diff ^ diff = 0 for pairs)
// Step 2: Find any differing bit (rightmost set bit in xor)
// Step 3: Split numbers into two groups by that bit, XOR each group
// Time: O(n), Space: O(1)
public int[] singleNumberIII(int[] nums) {
    int xor = 0;
    for (int num : nums) xor ^= num;

    int diffBit = xor & (-xor);    // isolate rightmost differing bit
    int a = 0, b = 0;
    for (int num : nums) {
        if ((num & diffBit) != 0) a ^= num;
        else                      b ^= num;
    }
    return new int[]{a, b};
}
```

---

## 4. Subsets Using Bitmask

Each subset of an n-element array corresponds to a bitmask from `0` to `2^n - 1`.

```java
// Time: O(n * 2^n), Space: O(n)
public List<List<Integer>> subsets(int[] nums) {
    int n = nums.length;
    List<List<Integer>> result = new ArrayList<>();

    for (int mask = 0; mask < (1 << n); mask++) {    // iterate all 2^n subsets
        List<Integer> subset = new ArrayList<>();
        for (int i = 0; i < n; i++) {
            if ((mask & (1 << i)) != 0) {            // if bit i is set, include nums[i]
                subset.add(nums[i]);
            }
        }
        result.add(subset);
    }
    return result;
}
// nums=[1,2,3]
// mask=0 (000) → []
// mask=1 (001) → [1]
// mask=2 (010) → [2]
// mask=3 (011) → [1,2]
// mask=4 (100) → [3]
// mask=5 (101) → [1,3]
// mask=6 (110) → [2,3]
// mask=7 (111) → [1,2,3]

// ── Check if bit i is set ──
boolean isSet(int mask, int i) { return (mask >> i & 1) == 1; }
```

---

## 5. Reverse Bits

```java
// Reverse bits of a 32-bit unsigned integer
// Time: O(32) = O(1), Space: O(1)
public int reverseBits(int n) {
    int result = 0;
    for (int i = 0; i < 32; i++) {
        result = (result << 1) | (n & 1);   // shift result left, add LSB of n
        n >>>= 1;                            // unsigned right shift n
    }
    return result;
}
// n = 43261596 (00000010100101000001111010011100)
//     ↓ reversed
// result = 964176192 (00111001011110000010100101000000)

// ── Bit swap (swap bit i and j) ──
int swapBits(int n, int i, int j) {
    if (((n >> i) & 1) != ((n >> j) & 1)) {   // only swap if they differ
        n ^= (1 << i) | (1 << j);
    }
    return n;
}

// ── Add two numbers without + operator ──
public int addWithoutPlus(int a, int b) {
    while (b != 0) {
        int carry = a & b;      // bits that cause carry
        a = a ^ b;              // sum without carry
        b = carry << 1;         // carry shifted left
    }
    return a;
}
```

> **Interview Q: What does `n & (n-1)` do and why is it useful?**  
> `n & (n-1)` **clears the lowest set bit** of `n`. It's useful for: (1) checking power of 2 — if result is 0, only one bit was set; (2) counting set bits (Brian Kernighan) — iterate until `n == 0`; (3) checking if a number's binary representation is a subset of another number's bits.
