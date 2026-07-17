# Bit Manipulation

> Bit manipulation is O(1) and avoids division/modulo. In Go, all integer types support bitwise operators.

---

## Table of Contents

1. [Bit Operations Cheat Sheet](#1-bit-operations-cheat-sheet)
2. [Check / Count / Set Bits](#2-check--count--set-bits)
3. [Power of Two](#3-power-of-two)
4. [Single Number I, II, III](#4-single-number-i-ii-iii)
5. [Subsets Using Bitmask](#5-subsets-using-bitmask)
6. [Reverse Bits](#6-reverse-bits)

---

## 1. Bit Operations Cheat Sheet

```go
// Basic operators (same as C/Java)
a & b   // AND  — both bits 1
a | b   // OR   — at least one bit 1
a ^ b   // XOR  — exactly one bit 1
^a      // NOT  — flip all bits (Go uses ^ as unary NOT, not ~)
a << n  // left shift  — multiply by 2^n
a >> n  // right shift — divide by 2^n (arithmetic for signed)

// Common tricks
x & 1          // check if x is odd (1 = odd, 0 = even)
x & (x - 1)    // clear lowest set bit
x & (-x)       // isolate lowest set bit
x | (1 << k)   // set bit k
x & ^(1 << k)  // clear bit k   (^(1<<k) = NOT mask in Go)
x ^ (1 << k)   // toggle bit k
(x >> k) & 1   // check bit k

// Integer sizes in Go
// int  — platform-dependent (64-bit on 64-bit systems)
// int32, int64, uint32, uint64 — explicit sizes
// Use uint32 for problems that treat the number as 32 bits (e.g., Reverse Bits)
```

---

## 2. Check / Count / Set Bits

```go
// Check if bit k is set
func isBitSet(x, k int) bool {
    return (x>>k)&1 == 1
}

// Count set bits — Kernighan's algorithm
// Time: O(number of set bits)
func countBits(x int) int {
    count := 0
    for x != 0 {
        x &= x - 1 // clear lowest set bit
        count++
    }
    return count
}

// Count set bits — popcount via DP (for range [0, n])
// "Counting Bits" problem
func countBitsRange(n int) []int {
    dp := make([]int, n+1)
    for i := 1; i <= n; i++ {
        dp[i] = dp[i>>1] + (i & 1) // dp[i/2] + last bit
    }
    return dp
}

// Hamming distance between x and y
func hammingDistance(x int, y int) int {
    return countBits(x ^ y)
}
```

---

## 3. Power of Two

```go
// Power of 2 — exactly one bit set
func isPowerOfTwo(n int) bool {
    return n > 0 && n&(n-1) == 0
}

// Power of 4 — power of 2 AND the set bit is at an even position
// 0x55555555 = 01010101...01 in binary
func isPowerOfFour(n int) bool {
    return n > 0 && n&(n-1) == 0 && n&0x55555555 != 0
}
```

---

## 4. Single Number I, II, III

### Single Number I — one element appears once, rest twice

```go
// XOR all elements: duplicates cancel out
// Time: O(n), Space: O(1)
func singleNumber(nums []int) int {
    result := 0
    for _, n := range nums {
        result ^= n
    }
    return result
}
```

### Single Number II — one element appears once, rest three times

```go
// Count set bits for each position mod 3
// Time: O(32n) = O(n), Space: O(1)
func singleNumberII(nums []int) int {
    result := 0
    for bit := 0; bit < 32; bit++ {
        sum := 0
        for _, n := range nums {
            sum += (n >> bit) & 1
        }
        if sum%3 != 0 {
            result |= 1 << bit
        }
    }
    // Handle sign for 32-bit interpretation
    if result >= (1 << 31) {
        result -= (1 << 32)
    }
    return result
}

// Elegant two-variable solution
func singleNumberIIAlt(nums []int) int {
    ones, twos := 0, 0
    for _, n := range nums {
        ones = (ones ^ n) & ^twos
        twos = (twos ^ n) & ^ones
    }
    return ones
}
```

### Single Number III — two elements appear once, rest twice

```go
// Time: O(n), Space: O(1)
func singleNumberIII(nums []int) []int {
    // XOR of all → XOR of the two unique numbers
    xor := 0
    for _, n := range nums { xor ^= n }

    // Find any differing bit (rightmost set bit)
    diffBit := xor & (-xor)

    // Split into two groups and XOR each
    a, b := 0, 0
    for _, n := range nums {
        if n&diffBit != 0 {
            a ^= n
        } else {
            b ^= n
        }
    }
    return []int{a, b}
}
```

> **Interview Q: Why use `xor & (-xor)` to find the differing bit?**  
> `-xor` in two's complement flips all bits and adds 1. `xor & (-xor)` isolates the lowest set bit — this bit is 1 in one unique number and 0 in the other, allowing us to separate them.

---

## 5. Subsets Using Bitmask

```go
// Enumerate all 2^n subsets via bitmask
// Time: O(n * 2^n), Space: O(n)
func subsetsWithBitmask(nums []int) [][]int {
    n := len(nums)
    total := 1 << n // 2^n
    result := make([][]int, 0, total)

    for mask := 0; mask < total; mask++ {
        subset := []int{}
        for i := 0; i < n; i++ {
            if mask&(1<<i) != 0 {
                subset = append(subset, nums[i])
            }
        }
        result = append(result, subset)
    }
    return result
}

// Iterate over all subsets of a given set (bit trick)
// Enumerate all subsets of the bitmask `s`
func enumerateSubsets(s int) {
    for sub := s; sub > 0; sub = (sub - 1) & s {
        // process `sub`
        _ = sub
    }
    // process empty set (sub == 0) if needed
}
```

---

## 6. Reverse Bits

```go
// Reverse the 32-bit binary representation
// Time: O(32) = O(1)
func reverseBits(num uint32) uint32 {
    result := uint32(0)
    for i := 0; i < 32; i++ {
        result = (result << 1) | (num & 1)
        num >>= 1
    }
    return result
}

// With memoization (when called many times)
var cache = map[byte]uint32{}

func reverseByteCached(b byte) uint32 {
    if v, ok := cache[b]; ok { return v }
    result := uint32(0)
    for i := 0; i < 8; i++ {
        result = (result << 1) | uint32(b&1)
        b >>= 1
    }
    cache[b] = result
    return result
}

func reverseBitsCached(num uint32) uint32 {
    return reverseByteCached(byte(num)) << 24 |
           reverseByteCached(byte(num>>8)) << 16 |
           reverseByteCached(byte(num>>16)) << 8 |
           reverseByteCached(byte(num>>24))
}
```

---

## Bit Manipulation Quick Reference

| Goal | Expression | Example (x=6=110) |
|---|---|---|
| Check if odd | `x & 1` | `6 & 1 = 0` (even) |
| Lowest set bit | `x & (-x)` | `6 & -6 = 2` (010) |
| Clear lowest set bit | `x & (x-1)` | `6 & 5 = 4` (100) |
| Set bit k | `x \| (1 << k)` | `6 \| (1<<0) = 7` |
| Clear bit k | `x & ^(1 << k)` | `6 & ^(1<<1) = 4` |
| Toggle bit k | `x ^ (1 << k)` | `6 ^ (1<<0) = 7` |
| Check bit k | `(x >> k) & 1` | `(6>>1) & 1 = 1` |
| XOR cancel duplicates | `a ^ a = 0` | `5 ^ 5 = 0` |
