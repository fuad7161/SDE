# 05 — Memory Management & Pointers

> **Topics:** stack vs heap · escape analysis · value vs pointer · garbage collector · allocation performance

---

## Table of Contents
1. [Stack vs Heap](#1-stack-vs-heap)
2. [Escape Analysis](#2-escape-analysis)
3. [Value vs Pointer — When to Use Each](#3-value-vs-pointer--when-to-use-each)
4. [Garbage Collector](#4-garbage-collector)
5. [Performance of Heap Allocations](#5-performance-of-heap-allocations)
6. [Interview Questions](#interview-questions)

---

## 1. Stack vs Heap

| | Stack | Heap |
|---|---|---|
| Managed by | Goroutine (automatic) | Go runtime (GC) |
| Allocation cost | Near zero | Higher (GC pressure) |
| Lifetime | Until function returns | Until GC collects |
| Access speed | Very fast | Slightly slower |
| Size | ~2 KB initial (grows) | Limited by system RAM |

```go
package main

import "fmt"

func stackAlloc() int {
    x := 42       // x stays on the stack — doesn't escape
    return x      // returned by value, copy leaves
}

func heapAlloc() *int {
    x := 42
    return &x     // x escapes to heap — its address outlives the function
}

func main() {
    a := stackAlloc()
    b := heapAlloc()
    fmt.Println(a, *b)
}
```

---

## 2. Escape Analysis

The Go compiler decides at compile time whether a variable **escapes to the heap**. You can inspect this with:

```bash
go build -gcflags="-m" ./...
# or more verbose:
go build -gcflags="-m -m" ./...
```

### Common escape triggers

```go
package main

import "fmt"

// 1. Returning a pointer — x escapes
func newInt(v int) *int {
    x := v   // "./main.go:N:2: moved to heap: x"
    return &x
}

// 2. Interface boxing — value escapes when stored in interface{}
func boxed(v any) {
    fmt.Println(v) // v may escape (fmt.Println takes interface{})
}

// 3. Closure captures a variable — it escapes
func makeCounter() func() int {
    count := 0   // count escapes — closure outlives makeCounter
    return func() int {
        count++
        return count
    }
}

// 4. Slice with unknown size — heap allocated
func dynamicSlice(n int) []int {
    return make([]int, n) // n not known at compile time → heap
}

// Does NOT escape — size known, small, not returned by pointer
func fixedSlice() [3]int {
    return [3]int{1, 2, 3}
}
```

---

## 3. Value vs Pointer — When to Use Each

```go
package main

import "fmt"

type SmallPoint struct{ X, Y float64 } // 16 bytes — prefer value

type LargeConfig struct {
    Host     string
    Port     int
    Timeout  int
    Retries  int
    Tags     []string
    Headers  map[string]string
} // large — prefer pointer

// Value receiver — good for small, immutable types
func (p SmallPoint) Distance() float64 {
    return p.X*p.X + p.Y*p.Y
}

// Pointer receiver — required when mutating, preferred for large structs
func (c *LargeConfig) SetHost(host string) {
    c.Host = host
}

func processConfig(c *LargeConfig) { // pass by pointer — avoid copying
    fmt.Println(c.Host)
}

func processPoint(p SmallPoint) { // pass by value — cheap copy, safer
    fmt.Println(p.X, p.Y)
}
```

### Decision guide

| Situation | Use |
|---|---|
| Small struct (≤ a few words) | Value |
| Struct that must be mutated | Pointer |
| Large struct | Pointer (avoid copy cost) |
| Struct used with interface | Pointer (consistent method set) |
| Primitive types (int, string) | Value |

---

## 4. Garbage Collector

Go uses a **concurrent, tri-color mark-and-sweep** GC that runs alongside the program.

**Three phases:**
1. **Mark** — traverse the object graph from roots (globals, stacks), mark live objects
2. **Sweep** — reclaim memory of unmarked (dead) objects
3. **Scavenge** — return unused memory to the OS

**Key properties:**
- Concurrent (most work happens while the program runs)
- Low-latency: typical STW (stop-the-world) pauses are < 1 ms
- GC is triggered when heap doubles since last collection (`GOGC=100` by default)

```go
package main

import (
    "fmt"
    "runtime"
)

func main() {
    var stats runtime.MemStats
    runtime.ReadMemStats(&stats)

    fmt.Printf("Alloc:      %d KB\n", stats.Alloc/1024)
    fmt.Printf("TotalAlloc: %d KB\n", stats.TotalAlloc/1024)
    fmt.Printf("NumGC:      %d\n", stats.NumGC)
    fmt.Printf("PauseTotalNs: %d ns\n", stats.PauseTotalNs)

    // Force a GC (testing/benchmarking only — never in prod)
    runtime.GC()
}
```

**Tuning:**
```bash
GOGC=200   # GC less frequently (higher memory, lower CPU)
GOGC=off   # Disable GC entirely (only for short-lived programs)
GOMEMLIMIT=500MiB  # Soft memory limit (Go 1.19+)
```

---

## 5. Performance of Heap Allocations

Each heap allocation:
1. Asks the runtime for memory
2. Is tracked by the GC (increases GC pressure)
3. May cause a GC cycle (STW pause)

### Reduce allocations — patterns

```go
package main

import (
    "fmt"
    "sync"
)

// Pattern 1: sync.Pool — reuse objects instead of allocating
var bufPool = sync.Pool{
    New: func() any { return make([]byte, 0, 1024) },
}

func processRequest(data []byte) {
    buf := bufPool.Get().([]byte)
    defer bufPool.Put(buf[:0]) // reset and return

    buf = append(buf, data...)
    fmt.Println("processed", len(buf), "bytes")
}

// Pattern 2: Pre-allocate slices when size is known
func collectIDs(n int) []int {
    ids := make([]int, 0, n) // pre-allocate capacity — no reallocs
    for i := 0; i < n; i++ {
        ids = append(ids, i)
    }
    return ids
}

// Pattern 3: Reuse structs across calls instead of creating new ones
type Processor struct {
    buf []byte // reused across calls
}

func (p *Processor) Process(data []byte) {
    p.buf = p.buf[:0]           // reset without reallocating
    p.buf = append(p.buf, data...)
}
```

```bash
# Profile allocations
go test -bench=. -benchmem ./...
# Output: BenchmarkX-8   1000   1234 ns/op   256 B/op   4 allocs/op

# Find allocation hot spots
go tool pprof -alloc_objects cpu.prof
```

---

## Interview Questions

<details>
<summary><b>When does Go allocate memory on the heap vs the stack?</b></summary>

The compiler decides via **escape analysis**. A variable stays on the stack if its lifetime is bounded by the function call. It escapes to the heap if: its address is returned or stored in a longer-lived structure; it's captured by a closure; it's assigned to an interface; or its size is unknown at compile time. Stack allocation is near-zero cost; heap allocation incurs GC overhead.

</details>

<details>
<summary><b>What is escape analysis? How do you check it?</b></summary>

Escape analysis is a compile-time analysis that determines whether a variable's lifetime can be bounded to a stack frame. Run `go build -gcflags="-m"` to see which variables "escape to heap". Common triggers: returning a pointer, storing in an interface, closures, and dynamic-size allocations.

</details>

<details>
<summary><b>What is the difference between passing a value vs a pointer to a function? When should you prefer each?</b></summary>

Passing by **value** copies the data — safe, no aliasing, good for small types (int, bool, small structs). Passing by **pointer** shares the original — required for mutation, cheaper for large structs (avoids copy). Pointer semantics introduce aliasing risks (the callee can modify the original). Rule of thumb: pointer for large structs or when mutation is needed; value otherwise.

</details>

<details>
<summary><b>How does Go's garbage collector work at a high level?</b></summary>

Go uses a **concurrent tri-color mark-and-sweep** GC. It marks live objects by traversing the object graph from roots (goroutine stacks, globals). Unmarked objects are swept (freed). Most work is concurrent with the application; stop-the-world pauses are typically sub-millisecond. GC frequency is controlled by `GOGC` (default 100 = trigger when heap doubles).

</details>

<details>
<summary><b>What are the performance implications of frequent small heap allocations?</b></summary>

Each allocation increases GC pressure and can trigger more frequent GC cycles (causing STW pauses). Mitigation strategies: `sync.Pool` to reuse objects, pre-allocating slices with `make([]T, 0, n)`, reusing struct buffers across calls, and using value types instead of pointers for small structs. Use `go test -benchmem` and `pprof -alloc_objects` to find hot spots.

</details>
