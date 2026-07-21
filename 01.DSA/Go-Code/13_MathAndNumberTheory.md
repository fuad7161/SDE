# Math & Number Theory

> Number theory problems appear often as sub-components of larger solutions. Know your primes, GCD, and modular exponentiation.

---

## Table of Contents

1. [Sieve of Eratosthenes](#1-sieve-of-eratosthenes)
2. [GCD & LCM](#2-gcd--lcm)
3. [Fast Exponentiation (Binary Exponentiation)](#3-fast-exponentiation-binary-exponentiation)
4. [Palindrome Number](#4-palindrome-number)
5. [Roman to Integer / Integer to Roman](#5-roman-to-integer--integer-to-roman)
6. [Bonus Patterns](#6-bonus-patterns)

---

## 1. Sieve of Eratosthenes

```go
// Find all primes up to n — Time: O(n log log n), Space: O(n)
func sieve(n int) []int {
    isPrime := make([]bool, n+1)
    for i := 2; i <= n; i++ { isPrime[i] = true }

    for i := 2; i*i <= n; i++ {
        if isPrime[i] {
            for j := i * i; j <= n; j += i {
                isPrime[j] = false
            }
        }
    }

    primes := []int{}
    for i := 2; i <= n; i++ {
        if isPrime[i] {
            primes = append(primes, i)
        }
    }
    return primes
}

// Count primes less than n (LeetCode 204)
func countPrimes(n int) int {
    if n <= 2 { return 0 }
    isPrime := make([]bool, n)
    for i := 2; i < n; i++ { isPrime[i] = true }
    for i := 2; i*i < n; i++ {
        if isPrime[i] {
            for j := i * i; j < n; j += i {
                isPrime[j] = false
            }
        }
    }
    count := 0
    for _, v := range isPrime {
        if v { count++ }
    }
    return count
}

// Check if a single number is prime
func isPrime(n int) bool {
    if n < 2 { return false }
    if n == 2 { return true }
    if n%2 == 0 { return false }
    for i := 3; i*i <= n; i += 2 {
        if n%i == 0 { return false }
    }
    return true
}
```

> **Interview Q: Why does the inner loop in Sieve start at `i*i`?**  
> All composite numbers with a factor smaller than i have already been marked when we processed that smaller factor. So `i*i` is the first new composite that i can mark.

---

## 2. GCD & LCM

```go
// Euclidean GCD — Time: O(log(min(a,b)))
func gcd(a, b int) int {
    for b != 0 {
        a, b = b, a%b
    }
    return a
}

// Recursive variant
func gcdRec(a, b int) int {
    if b == 0 { return a }
    return gcdRec(b, a%b)
}

// LCM — avoid overflow with division first
func lcm(a, b int) int {
    return a / gcd(a, b) * b
}

// GCD of an array
func gcdArray(nums []int) int {
    result := nums[0]
    for i := 1; i < len(nums); i++ {
        result = gcd(result, nums[i])
    }
    return result
}

// Note: Go 1.21+ has math/big.GCD for big integers.
// For standard int, the above is the idiomatic approach.
```

---

## 3. Fast Exponentiation (Binary Exponentiation)

```go
// Compute base^exp in O(log exp)
func fastPow(base, exp int) int {
    result := 1
    for exp > 0 {
        if exp&1 == 1 {
            result *= base
        }
        base *= base
        exp >>= 1
    }
    return result
}

// Modular exponentiation — avoids overflow
func modPow(base, exp, mod int) int {
    result := 1
    base %= mod
    for exp > 0 {
        if exp&1 == 1 {
            result = result * base % mod
        }
        base = base * base % mod
        exp >>= 1
    }
    return result
}

// Pow(x, n) — LeetCode 50, handles negative exponents and float
func myPow(x float64, n int) float64 {
    if n < 0 {
        x = 1 / x
        n = -n
    }
    result := 1.0
    for n > 0 {
        if n&1 == 1 {
            result *= x
        }
        x *= x
        n >>= 1
    }
    return result
}
```

> **Interview Q: What is binary exponentiation?**  
> Instead of multiplying `base` n times (O(n)), we square at each step and multiply by base only when the current bit of n is 1. This gives O(log n) multiplications.

---

## 4. Palindrome Number

```go
// Check without converting to string — Time: O(log n)
func isPalindromeNum(x int) bool {
    if x < 0 { return false }
    if x != 0 && x%10 == 0 { return false } // e.g., 10, 100

    reversed := 0
    for x > reversed {
        reversed = reversed*10 + x%10
        x /= 10
    }
    // x == reversed (even length) OR x == reversed/10 (odd length)
    return x == reversed || x == reversed/10
}

// Check string palindrome
func isPalindromeStr(s string) bool {
    left, right := 0, len(s)-1
    for left < right {
        if s[left] != s[right] { return false }
        left++
        right--
    }
    return true
}

// Longest Palindromic Substring — Expand Around Center
func longestPalindrome(s string) string {
    start, maxLen := 0, 1

    expand := func(l, r int) {
        for l >= 0 && r < len(s) && s[l] == s[r] {
            if r-l+1 > maxLen {
                maxLen = r - l + 1
                start = l
            }
            l--
            r++
        }
    }

    for i := 0; i < len(s); i++ {
        expand(i, i)   // odd length
        expand(i, i+1) // even length
    }
    return s[start : start+maxLen]
}
```

---

## 5. Roman to Integer / Integer to Roman

```go
// Roman to Integer — Time: O(n)
func romanToInt(s string) int {
    val := map[byte]int{
        'I': 1, 'V': 5, 'X': 10, 'L': 50,
        'C': 100, 'D': 500, 'M': 1000,
    }
    result := 0
    for i := 0; i < len(s); i++ {
        if i+1 < len(s) && val[s[i]] < val[s[i+1]] {
            result -= val[s[i]] // subtractive notation: IV, IX, etc.
        } else {
            result += val[s[i]]
        }
    }
    return result
}

// Integer to Roman — Time: O(1) — at most 15 chars output
func intToRoman(num int) string {
    values := []int{1000, 900, 500, 400, 100, 90, 50, 40, 10, 9, 5, 4, 1}
    symbols := []string{"M", "CM", "D", "CD", "C", "XC", "L", "XL", "X", "IX", "V", "IV", "I"}

    result := []byte{}
    for i, v := range values {
        for num >= v {
            result = append(result, symbols[i]...)
            num -= v
        }
    }
    return string(result)
}
```

---

## 6. Bonus Patterns

### Factorial Trailing Zeros

```go
// Count factors of 5 — Time: O(log n)
func trailingZeroes(n int) int {
    count := 0
    for n >= 5 {
        n /= 5
        count += n
    }
    return count
}
```

### Integer Square Root (without math.Sqrt)

```go
// Binary search — Time: O(log n)
func mySqrt(x int) int {
    if x < 2 { return x }
    left, right := 1, x/2
    for left <= right {
        mid := left + (right-left)/2
        sq := mid * mid
        if sq == x {
            return mid
        } else if sq < x {
            left = mid + 1
        } else {
            right = mid - 1
        }
    }
    return right
}
```

### Fibonacci with Matrix Exponentiation

```go
// O(log n) Fibonacci using 2x2 matrix fast power
func fibMatrix(n int) int {
    if n <= 1 { return n }
    type mat [2][2]int
    multiply := func(a, b mat) mat {
        return mat{
            {a[0][0]*b[0][0] + a[0][1]*b[1][0], a[0][0]*b[0][1] + a[0][1]*b[1][1]},
            {a[1][0]*b[0][0] + a[1][1]*b[1][0], a[1][0]*b[0][1] + a[1][1]*b[1][1]},
        }
    }
    var power func(m mat, p int) mat
    power = func(m mat, p int) mat {
        if p == 1 { return m }
        half := power(m, p/2)
        result := multiply(half, half)
        if p%2 == 1 { result = multiply(result, m) }
        return result
    }
    base := mat{{1, 1}, {1, 0}}
    result := power(base, n)
    return result[0][1]
}
```

### Sum of Digits / Digital Root

```go
// Digital root — O(1)
func digitalRoot(n int) int {
    if n == 0 { return 0 }
    if n%9 == 0 { return 9 }
    return n % 9
}

// Sum of digits iteratively
func sumDigits(n int) int {
    if n < 0 { n = -n }
    sum := 0
    for n > 0 {
        sum += n % 10
        n /= 10
    }
    return sum
}
```

> **Go math package essentials:**
> ```go
> import "math"
> math.MaxInt    // maximum int value
> math.MinInt    // minimum int value
> math.Abs(f)    // float64 absolute value
> math.Sqrt(f)   // float64 square root
> math.Log2(f)   // float64 log base 2
> math.Pow(x, y) // float64 x^y
> ```
