# Math & Number Theory

> Clean O(1) or O(√n) solutions expected. Know Euclidean GCD, Sieve, and modular arithmetic.

---

## Table of Contents

1. [Sieve of Eratosthenes](#1-sieve-of-eratosthenes)
2. [GCD / LCM — Euclidean Algorithm](#2-gcd--lcm--euclidean-algorithm)
3. [Fast Exponentiation](#3-fast-exponentiation)
4. [Palindrome Number](#4-palindrome-number)
5. [Roman to Integer / Integer to Roman](#5-roman-to-integer--integer-to-roman)

---

## 1. Sieve of Eratosthenes

**Find all primes up to n.**  
**Algorithm:** Start with all numbers marked prime. For each prime `p`, mark all multiples `p², p²+p, p²+2p, ...` as composite.

```java
// Time: O(n log log n), Space: O(n)
public boolean[] sieve(int n) {
    boolean[] isPrime = new boolean[n + 1];
    Arrays.fill(isPrime, true);
    isPrime[0] = isPrime[1] = false;

    for (int p = 2; (long) p * p <= n; p++) {
        if (isPrime[p]) {
            for (int multiple = p * p; multiple <= n; multiple += p) {
                isPrime[multiple] = false;  // mark multiples as composite
            }
        }
    }
    return isPrime;
}

// ── Count primes less than n ──
public int countPrimes(int n) {
    if (n < 2) return 0;
    boolean[] isPrime = new boolean[n];
    Arrays.fill(isPrime, true);
    isPrime[0] = isPrime[1] = false;

    for (int p = 2; (long) p * p < n; p++) {
        if (isPrime[p]) {
            for (int m = p * p; m < n; m += p) isPrime[m] = false;
        }
    }
    int count = 0;
    for (boolean b : isPrime) if (b) count++;
    return count;
}

// ── Simple primality check for a single number ──
// Time: O(√n)
public boolean isPrime(int n) {
    if (n < 2) return false;
    if (n == 2) return true;
    if (n % 2 == 0) return false;
    for (int i = 3; (long) i * i <= n; i += 2) {
        if (n % i == 0) return false;
    }
    return true;
}
```

> **Interview Q: Why do we start marking from p² in the Sieve?**  
> Any composite multiple of `p` smaller than `p²` has a prime factor less than `p` and was already marked by a previous prime. Starting at `p²` avoids redundant work. This is why the outer loop only needs to go up to `√n`.

---

## 2. GCD / LCM — Euclidean Algorithm

```java
// GCD (Greatest Common Divisor) — Euclidean Algorithm
// gcd(a, b) = gcd(b, a % b), base: gcd(a, 0) = a
// Time: O(log(min(a, b))), Space: O(log n) recursive / O(1) iterative

// Recursive
public int gcd(int a, int b) {
    return b == 0 ? a : gcd(b, a % b);
}

// Iterative
public int gcdIterative(int a, int b) {
    while (b != 0) {
        int temp = b;
        b = a % b;
        a = temp;
    }
    return a;
}

// LCM (Least Common Multiple)
// lcm(a, b) = (a * b) / gcd(a, b)
// Use long to avoid overflow
public long lcm(int a, int b) {
    return (long) a / gcd(a, b) * b;  // divide first to reduce overflow risk
}

// Example:
// gcd(48, 18): 48%18=12 → gcd(18,12): 18%12=6 → gcd(12,6): 12%6=0 → gcd(6,0)=6
// lcm(48,18) = 48*18/6 = 144

// ── GCD of array ──
public int gcdArray(int[] nums) {
    int result = nums[0];
    for (int num : nums) result = gcd(result, num);
    return result;
}

// ── Check if two numbers are coprime ──
boolean isCoprime(int a, int b) { return gcd(a, b) == 1; }
```

---

## 3. Fast Exponentiation

**Compute `base^exp % mod` in O(log exp) instead of O(exp).**  
**Idea:** Square the base at each step, halving the exponent.

```java
// Modular fast exponentiation — Time: O(log exp), Space: O(1)
public long power(long base, long exp, long mod) {
    long result = 1;
    base %= mod;

    while (exp > 0) {
        if ((exp & 1) == 1) {          // if current bit is set (odd exponent)
            result = result * base % mod;
        }
        base = base * base % mod;      // square the base
        exp >>= 1;                     // right shift (divide exponent by 2)
    }
    return result;
}
// power(2, 10, 1000) = 1024 % 1000 = 24
// 2^10: exp=10(1010) → base=2,4,16,256 → picks base when bit=1 → 4*256=1024

// ── Without modulo (returns double for large numbers) ──
public double myPow(double base, int exp) {
    if (exp < 0) { base = 1.0 / base; exp = -exp; }
    double result = 1.0;
    while (exp > 0) {
        if ((exp & 1) == 1) result *= base;
        base *= base;
        exp >>= 1;
    }
    return result;
}
```

> **Interview Q: How does binary exponentiation achieve O(log n)?**  
> Instead of multiplying `base` n times, we express n in binary and only square `base` `log(n)` times. At each step, we either use or skip the current power based on the bit. This reduces multiplications from n to log(n).

---

## 4. Palindrome Number

```java
// ── Without string conversion ──
// Time: O(log n), Space: O(1)
public boolean isPalindrome(int x) {
    if (x < 0 || (x % 10 == 0 && x != 0)) return false;   // negatives and trailing zeros

    int reversed = 0;
    // Only reverse half the number (stop when reversed >= x)
    while (x > reversed) {
        reversed = reversed * 10 + x % 10;
        x /= 10;
    }
    // For even digits: x == reversed
    // For odd digits:  x == reversed / 10 (middle digit doesn't matter)
    return x == reversed || x == reversed / 10;
}
// 121 → reversed half: 12 → x=1, reversed=12 → 1 == 12/10=1 ✓
// 1221 → reversed half: 12 → x=12, reversed=12 → 12==12 ✓
// 123 → reversed half: 32 → x=1, reversed=32 → 1!=32 ✗

// ── With string conversion (simple but uses extra space) ──
public boolean isPalindromeString(int x) {
    if (x < 0) return false;
    String s = Integer.toString(x);
    int lo = 0, hi = s.length() - 1;
    while (lo < hi) {
        if (s.charAt(lo++) != s.charAt(hi--)) return false;
    }
    return true;
}
```

---

## 5. Roman to Integer / Integer to Roman

### Roman to Integer

```java
// Time: O(n), Space: O(1)
public int romanToInt(String s) {
    Map<Character, Integer> val = new HashMap<>();
    val.put('I', 1);   val.put('V', 5);   val.put('X', 10);
    val.put('L', 50);  val.put('C', 100); val.put('D', 500);
    val.put('M', 1000);

    int result = 0;
    for (int i = 0; i < s.length(); i++) {
        int curr = val.get(s.charAt(i));
        int next = (i + 1 < s.length()) ? val.get(s.charAt(i + 1)) : 0;

        if (curr < next) result -= curr;   // subtractive notation (IV, IX, XL...)
        else             result += curr;
    }
    return result;
}
// "MCMXCIV" → M=1000, CM=900, XC=90, IV=4 → 1994
```

### Integer to Roman

```java
// Time: O(1) — bounded by the 13 Roman numeral values
public String intToRoman(int num) {
    int[]    values = {1000,900,500,400,100,90,50,40,10,9,5,4,1};
    String[] symbols = {"M","CM","D","CD","C","XC","L","XL","X","IX","V","IV","I"};

    StringBuilder sb = new StringBuilder();
    for (int i = 0; i < values.length; i++) {
        while (num >= values[i]) {
            sb.append(symbols[i]);
            num -= values[i];
        }
    }
    return sb.toString();
}
// 1994 → 1994-1000=994→M, 994-900=94→CM, 94-90=4→XC, 4-4=0→IV → "MCMXCIV"
```

---

## Bonus: Common Math Patterns

```java
// ── Count digits ──
int digitCount(int n) { return (int) Math.log10(n) + 1; }  // O(1)

// ── Integer square root (floor) ──
int isqrt(int n) {
    int lo = 0, hi = n;
    while (lo < hi) {
        int mid = lo + (hi - lo + 1) / 2;
        if ((long) mid * mid <= n) lo = mid;
        else hi = mid - 1;
    }
    return lo;
}

// ── Modular arithmetic identities ──
// (a + b) % m = ((a % m) + (b % m)) % m
// (a * b) % m = ((a % m) * (b % m)) % m
// (a - b + m) % m   ← ensures non-negative result

// ── Number of trailing zeros in n! ──
// Each trailing zero = one factor of 10 = one 5 (2s are always more plentiful)
int trailingZeros(int n) {
    int count = 0;
    while (n >= 5) { n /= 5; count += n; }
    return count;
}

// ── Check if n is a perfect square without sqrt() ──
boolean isPerfectSquare(int n) {
    if (n < 1) return false;
    long lo = 1, hi = n;
    while (lo < hi) {
        long mid = lo + (hi - lo + 1) / 2;
        if (mid * mid <= n) lo = mid;
        else hi = mid - 1;
    }
    return lo * lo == n;
}
```
