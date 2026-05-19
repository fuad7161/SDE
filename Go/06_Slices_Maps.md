# 06 — Slices & Maps Internals

> **Topics:** slice header · append behavior · shared backing array · map internals · nil maps · sync.Map

---

## Table of Contents
1. [Slice Internals — Header, len, cap](#1-slice-internals--header-len-cap)
2. [append and Capacity Growth](#2-append-and-capacity-growth)
3. [Shared Backing Array Bugs](#3-shared-backing-array-bugs)
4. [Map Internals & Nil Maps](#4-map-internals--nil-maps)
5. [Map Concurrency & sync.Map](#5-map-concurrency--syncmap)
6. [Interview Questions](#interview-questions)

---

## 1. Slice Internals — Header, len, cap

A slice is a **three-word struct**:

```
┌─────────────┐
│  ptr        │  → pointer to underlying array
│  len  = 3   │  → number of accessible elements
│  cap  = 5   │  → total allocated capacity
└─────────────┘
      │
      ▼
  [ 1 | 2 | 3 | _ | _ ]   ← backing array (cap=5)
```

```go
package main

import "fmt"

func main() {
    // make allocates a backing array of cap=5, len=3
    s := make([]int, 3, 5)
    fmt.Println(len(s), cap(s)) // 3 5

    // Slicing does NOT copy — shares the backing array
    a := []int{1, 2, 3, 4, 5}
    b := a[1:4] // b = [2, 3, 4], len=3, cap=4 (from index 1 to end)

    fmt.Println(b)        // [2 3 4]
    fmt.Println(len(b), cap(b)) // 3 4

    // Modifying b modifies a — same backing array
    b[0] = 99
    fmt.Println(a) // [1 99 3 4 5]
}
```

---

## 2. append and Capacity Growth

`append` returns a new slice header. If `len == cap`, it allocates a **new, larger backing array** and copies the data.

```go
package main

import "fmt"

func main() {
    s := make([]int, 0, 3)

    for i := 1; i <= 6; i++ {
        prev := cap(s)
        s = append(s, i)
        if cap(s) != prev {
            fmt.Printf("grew at len=%d: cap %d → %d\n", len(s), prev, cap(s))
        }
    }
    // grew at len=4: cap 3 → 6
    // grew at len=7: cap 6 → ...
    // Growth factor is ~2x for small slices, ~1.25x for large ones (Go runtime heuristic)
}
```

### Pitfall — append doesn't modify the original slice header

```go
func addElement(s []int, v int) {
    s = append(s, v) // modifies LOCAL copy of the header
    // original slice in caller is unchanged if reallocation occurred
}

// ✅ Correct — return the new slice
func addElementOK(s []int, v int) []int {
    return append(s, v)
}
```

### Pre-allocate when size is known

```go
// ❌ Many reallocations
result := []int{}
for i := 0; i < 10000; i++ {
    result = append(result, i)
}

// ✅ Single allocation
result := make([]int, 0, 10000)
for i := 0; i < 10000; i++ {
    result = append(result, i)
}
```

---

## 3. Shared Backing Array Bugs

Two slices sharing the same backing array can **silently corrupt each other**.

```go
package main

import "fmt"

func main() {
    original := []int{1, 2, 3, 4, 5}

    // Both share the same backing array
    s1 := original[:3] // [1 2 3], cap=5
    s2 := original[2:] // [3 4 5], cap=3

    // Modifying s1 affects original and s2
    s1[2] = 99
    fmt.Println(original) // [1 2 99 4 5]
    fmt.Println(s2)       // [99 4 5] — unexpected!

    // ✅ Fix: use copy to get an independent slice
    s3 := make([]int, len(s1))
    copy(s3, s1)
    s3[0] = 777
    fmt.Println(original) // unaffected
}
```

### append on a sub-slice overwrites data

```go
a := []int{1, 2, 3, 4, 5}
b := a[:3]              // shares backing array, cap=5

b = append(b, 99)       // writes 99 at index 3 — OVERWRITES a[3]!
fmt.Println(a)          // [1 2 3 99 5] — a[3] changed!

// ✅ Use full slice expression to restrict cap
b2 := a[:3:3]           // len=3, cap=3 — append will allocate new array
b2 = append(b2, 99)
fmt.Println(a)          // [1 2 3 4 5] — unchanged
```

---

## 4. Map Internals & Nil Maps

A Go map is a **hash table** implemented as an array of buckets. Each bucket holds up to 8 key-value pairs. When load factor exceeds ~6.5, the map grows (doubles buckets) and rehashes.

```go
package main

import "fmt"

func main() {
    // Zero value of a map is nil
    var m map[string]int
    fmt.Println(m == nil)   // true
    fmt.Println(len(m))     // 0 — safe

    // ✅ Reading from nil map returns zero value — safe
    fmt.Println(m["key"])   // 0

    // ❌ Writing to nil map panics
    // m["key"] = 1  // panic: assignment to entry in nil map

    // ✅ Always initialize before writing
    m = make(map[string]int)
    m["key"] = 1
    fmt.Println(m["key"])   // 1

    // Check existence
    v, ok := m["missing"]
    fmt.Println(v, ok)      // 0 false
}
```

### Map is not ordered — iteration order is random

```go
m := map[string]int{"a": 1, "b": 2, "c": 3}
for k, v := range m {
    fmt.Printf("%s=%d\n", k, v) // order varies each run
}

// If order matters: collect keys, sort, then iterate
import "sort"
keys := make([]string, 0, len(m))
for k := range m {
    keys = append(keys, k)
}
sort.Strings(keys)
for _, k := range keys {
    fmt.Printf("%s=%d\n", k, m[k])
}
```

---

## 5. Map Concurrency & sync.Map

The built-in map is **not safe for concurrent use**. Concurrent reads are fine; concurrent read+write or write+write causes a **data race** (detected by `-race` flag).

```go
// ❌ DATA RACE — two goroutines writing to the same map
m := make(map[int]int)
go func() { m[1] = 1 }()
go func() { m[2] = 2 }()
// run with: go run -race main.go
```

### Option 1 — Mutex-protected map (most common)

```go
package main

import (
    "fmt"
    "sync"
)

type SafeMap struct {
    mu sync.RWMutex
    m  map[string]int
}

func NewSafeMap() *SafeMap {
    return &SafeMap{m: make(map[string]int)}
}

func (s *SafeMap) Set(key string, val int) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.m[key] = val
}

func (s *SafeMap) Get(key string) (int, bool) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    v, ok := s.m[key]
    return v, ok
}
```

### Option 2 — sync.Map (specialized use cases)

`sync.Map` is optimized for two specific patterns:
- A key is written **once** but read many times
- Multiple goroutines operate on **disjoint sets of keys**

```go
package main

import (
    "fmt"
    "sync"
)

func main() {
    var sm sync.Map

    sm.Store("user:1", "Alice")
    sm.Store("user:2", "Bob")

    val, ok := sm.Load("user:1")
    fmt.Println(val, ok) // Alice true

    sm.Delete("user:2")

    sm.Range(func(k, v any) bool {
        fmt.Printf("%s = %s\n", k, v)
        return true // continue iteration
    })
}
```

> **Use a mutex-protected map** for general cases. Use `sync.Map` only when profiling shows the mutex is a bottleneck and your access pattern matches the two cases above.

---

## Interview Questions

<details>
<summary><b>What is the internal structure of a slice? What are len and cap?</b></summary>

A slice is a three-field struct: a pointer to a backing array, `len` (number of accessible elements), and `cap` (total capacity of the backing array from the slice's starting position). `len ≤ cap`. Multiple slices can point into the same backing array.

</details>

<details>
<summary><b>What happens when append exceeds the capacity of a slice?</b></summary>

The runtime allocates a new, larger backing array (typically ~2x the current capacity for small slices), copies all existing elements to it, appends the new element, and returns a new slice header pointing to the new array. The original slice is unaffected. This is why `append` must always be assigned back: `s = append(s, v)`.

</details>

<details>
<summary><b>Why can two slices share the same underlying array, and when does that cause bugs?</b></summary>

Slicing (`a[i:j]`) creates a new header pointing into the same backing array — no copy. This is efficient but dangerous: modifying one slice modifies the other. The most subtle bug is appending to a sub-slice when there's remaining capacity — it overwrites data in the original. Fix: use the full slice expression `a[i:j:j]` to cap capacity at `j`, forcing append to allocate a new array.

</details>

<details>
<summary><b>What is the zero value of a map? Can you read from a nil map? Write to one?</b></summary>

The zero value is `nil`. Reading from a nil map is safe — it returns the zero value for the value type. Writing to a nil map **panics**. Always initialize with `make(map[K]V)` before writing.

</details>

<details>
<summary><b>How does Go handle map concurrency? What is sync.Map and when should you use it?</b></summary>

The built-in map is not concurrency-safe. Concurrent writes (or read+write) cause a data race. The typical solution is a `sync.RWMutex`-protected wrapper. `sync.Map` is an alternative optimized for: (1) write-once, read-many keys, or (2) goroutines operating on disjoint keys. For general-purpose concurrent maps, a mutex-protected map is simpler and often faster.

</details>
