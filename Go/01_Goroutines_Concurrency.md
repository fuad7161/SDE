# 01 — Goroutines & Concurrency

> **Topics:** goroutine lifecycle · GMP scheduler · goroutine leaks

---

## Table of Contents
1. [Goroutine vs OS Thread](#1-goroutine-vs-os-thread)
2. [GMP Scheduler](#2-gmp-scheduler)
3. [Goroutine Leaks](#3-goroutine-leaks)
4. [Loop Variable Closure Bug](#4-loop-variable-closure-bug)
5. [How Many Goroutines?](#5-how-many-goroutines)
6. [Interview Questions](#interview-questions)

---

## 1. Goroutine vs OS Thread

| | Goroutine | OS Thread |
|---|---|---|
| Managed by | Go runtime | OS kernel |
| Stack size | ~2 KB (grows dynamically) | ~1–8 MB (fixed) |
| Creation cost | Very cheap (~1 µs) | Expensive (~1 ms) |
| Context switch | User-space (fast) | Kernel-space (slow) |
| Max count | Millions | Thousands |

A goroutine is a **lightweight, user-space thread** multiplexed onto OS threads by the Go runtime.

```go
package main

import (
    "fmt"
    "time"
)

func sayHello(name string) {
    fmt.Println("Hello,", name)
}

func main() {
    go sayHello("Alice") // launches a goroutine — non-blocking
    go sayHello("Bob")

    time.Sleep(10 * time.Millisecond) // wait for goroutines to finish
    // In real code, use sync.WaitGroup instead of Sleep
}
```

---

## 2. GMP Scheduler

The Go runtime uses a **GMP model** to schedule goroutines:

- **G** — Goroutine (the unit of work)
- **M** — Machine (OS thread)
- **P** — Processor (logical CPU, holds a run queue)

```
  G  G  G          G  G  G
  └──┴──┘          └──┴──┘
    P (run queue)    P (run queue)
    │                │
    M (OS thread)    M (OS thread)
         └───────────┘
           GOMAXPROCS
```

**Key rules:**
- `GOMAXPROCS` controls the number of Ps (default = number of CPU cores)
- A G runs on an M via a P
- If an M blocks (e.g. syscall), the runtime detaches it and assigns a new M to the P
- Work-stealing: idle Ps steal Gs from busy Ps

```go
package main

import (
    "fmt"
    "runtime"
)

func main() {
    fmt.Println("CPUs:", runtime.NumCPU())
    fmt.Println("GOMAXPROCS:", runtime.GOMAXPROCS(0)) // 0 = query without changing

    // Force single-threaded execution
    runtime.GOMAXPROCS(1)
    fmt.Println("Now GOMAXPROCS:", runtime.GOMAXPROCS(0))
}
```

---

## 3. Goroutine Leaks

A goroutine leak occurs when a goroutine is **started but never terminates** — typically because it's blocked waiting on a channel or a lock that is never released.

### Common causes

```go
// ❌ LEAK — goroutine blocks forever on a channel nobody will send to
func leak() {
    ch := make(chan int)
    go func() {
        val := <-ch // blocks forever — caller never sends
        fmt.Println(val)
    }()
    // ch goes out of scope, goroutine is stuck
}
```

### Fix — use context for cancellation

```go
package main

import (
    "context"
    "fmt"
    "time"
)

func worker(ctx context.Context, ch <-chan int) {
    for {
        select {
        case val, ok := <-ch:
            if !ok {
                return // channel closed
            }
            fmt.Println("received:", val)
        case <-ctx.Done():
            fmt.Println("worker cancelled:", ctx.Err())
            return // goroutine exits cleanly
        }
    }
}

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()

    ch := make(chan int, 5)
    go worker(ctx, ch)

    ch <- 1
    ch <- 2
    time.Sleep(3 * time.Second) // worker exits via ctx.Done after 2s
}
```

### Detection

```bash
# Check number of live goroutines at runtime
import "runtime"
fmt.Println(runtime.NumGoroutine())

# Use goleak in tests
go get go.uber.org/goleak
```

```go
func TestNoLeak(t *testing.T) {
    defer goleak.VerifyNone(t)
    // run your code
}
```

---

## 4. Loop Variable Closure Bug

Before Go 1.22, loop variables were **shared** across all goroutines in the loop.

```go
// ❌ BUG (Go < 1.22) — all goroutines print the same final value of i
for i := 0; i < 5; i++ {
    go func() {
        fmt.Println(i) // captures &i, not a copy
    }()
}
// Output: likely "5 5 5 5 5"
```

```go
// ✅ FIX 1 — shadow the variable inside the loop
for i := 0; i < 5; i++ {
    i := i // new variable scoped to this iteration
    go func() {
        fmt.Println(i)
    }()
}

// ✅ FIX 2 — pass as argument
for i := 0; i < 5; i++ {
    go func(n int) {
        fmt.Println(n)
    }(i)
}

// ✅ Go 1.22+ — loop variable is per-iteration by default (no fix needed)
```

---

## 5. How Many Goroutines?

There is **no hard limit** — the runtime can run millions. Practical limits:

- **Memory** — each goroutine starts with ~2 KB stack; 1M goroutines ≈ 2 GB minimum
- **GOMAXPROCS** — number of goroutines running *in parallel* ≤ GOMAXPROCS
- **Scheduler overhead** — too many runnable goroutines increases scheduling latency

```go
package main

import (
    "fmt"
    "runtime"
    "sync"
)

func main() {
    var wg sync.WaitGroup
    const n = 100_000

    for i := 0; i < n; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            // minimal work
        }()
    }

    wg.Wait()
    fmt.Println("goroutines still alive:", runtime.NumGoroutine()) // ~1 (main)
}
```

---

## Interview Questions

<details>
<summary><b>What is the difference between a goroutine and an OS thread?</b></summary>

Goroutines are managed by the Go runtime, not the OS. They start with a ~2 KB stack (vs ~1–8 MB for OS threads), are much cheaper to create (~1 µs vs ~1 ms), and are multiplexed onto OS threads by the GMP scheduler. You can run millions of goroutines; running thousands of OS threads is impractical.

</details>

<details>
<summary><b>How does the Go runtime schedule goroutines? Explain the GMP model.</b></summary>

Go uses the **GMP model**: Goroutines (G) run on OS threads (M) via logical processors (P). Each P has a local run queue of Gs. GOMAXPROCS controls the number of Ps. When an M blocks on a syscall, the runtime detaches it and pairs the P with a new M. Idle Ps can steal work from busy Ps (work-stealing).

</details>

<details>
<summary><b>What causes a goroutine leak? How do you detect and prevent one?</b></summary>

A leak happens when a goroutine blocks indefinitely — usually waiting on a channel or mutex that is never released. Prevention: always provide an exit path (context cancellation, closing the channel). Detection: `runtime.NumGoroutine()` during tests, or `go.uber.org/goleak`.

</details>

<details>
<summary><b>What happens when you launch a goroutine inside a loop and close over the loop variable?</b></summary>

Before Go 1.22, all goroutines capture the *same* loop variable by reference. By the time they run, the loop may have finished so they all see the final value. Fix: shadow the variable (`i := i`) or pass it as an argument. In Go 1.22+, loop variables are per-iteration.

</details>

<details>
<summary><b>How many goroutines can you run simultaneously? What limits them?</b></summary>

There's no hard limit. True parallelism is capped by `GOMAXPROCS` (defaults to CPU count). Practical limits are memory (~2 KB/goroutine initial stack) and scheduler overhead. Millions of goroutines are feasible if they spend most time blocked (e.g. waiting on I/O).

</details>
