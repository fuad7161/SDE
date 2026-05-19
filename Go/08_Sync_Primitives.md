# 08 — Sync Primitives

> **Topics:** sync.Mutex · sync.RWMutex · sync.WaitGroup · sync.Once · sync.Pool · atomic

---

## Table of Contents
1. [sync.Mutex vs Channel](#1-syncmutex-vs-channel)
2. [sync.RWMutex](#2-syncrwmutex)
3. [sync.WaitGroup](#3-syncwaitgroup)
4. [sync.Once](#4-synconce)
5. [sync.Pool](#5-syncpool)
6. [sync/atomic](#6-syncatomic)
7. [Interview Questions](#interview-questions)

---

## 1. sync.Mutex vs Channel

| | `sync.Mutex` | Channel |
|---|---|---|
| Best for | Protecting shared state | Communicating data / signaling |
| Mental model | Lock a critical section | Pass ownership of data |
| Ownership | Shared (any goroutine can lock/unlock) | Explicit (sender → receiver) |
| Composability | Lower | Higher (select, fan-in, pipelines) |

```go
package main

import (
    "fmt"
    "sync"
)

// Mutex — protecting shared state (a counter)
type Counter struct {
    mu    sync.Mutex
    value int
}

func (c *Counter) Inc() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.value++
}

func (c *Counter) Value() int {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.value
}

func main() {
    c := &Counter{}
    var wg sync.WaitGroup

    for i := 0; i < 1000; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            c.Inc()
        }()
    }
    wg.Wait()
    fmt.Println(c.Value()) // always 1000
}
```

---

## 2. sync.RWMutex

Allows **multiple concurrent readers** but only **one writer** at a time. Use it when reads vastly outnumber writes.

```go
package main

import (
    "fmt"
    "sync"
)

type Cache struct {
    mu    sync.RWMutex
    items map[string]string
}

func NewCache() *Cache {
    return &Cache{items: make(map[string]string)}
}

func (c *Cache) Set(key, val string) {
    c.mu.Lock()         // exclusive write lock
    defer c.mu.Unlock()
    c.items[key] = val
}

func (c *Cache) Get(key string) (string, bool) {
    c.mu.RLock()        // shared read lock — multiple goroutines can hold this simultaneously
    defer c.mu.RUnlock()
    v, ok := c.items[key]
    return v, ok
}

func main() {
    cache := NewCache()
    cache.Set("user:1", "Alice")

    var wg sync.WaitGroup
    for i := 0; i < 5; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            v, _ := cache.Get("user:1") // all 5 goroutines read concurrently
            fmt.Println(v)
        }()
    }
    wg.Wait()
}
```

> **Note:** `RWMutex` can be slower than `Mutex` under high write contention due to reader-writer coordination overhead. Benchmark before assuming it's faster.

---

## 3. sync.WaitGroup

Waits for a collection of goroutines to finish.

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

func worker(id int, wg *sync.WaitGroup) {
    defer wg.Done() // called when this goroutine finishes
    fmt.Printf("worker %d starting\n", id)
    time.Sleep(time.Duration(id) * 100 * time.Millisecond)
    fmt.Printf("worker %d done\n", id)
}

func main() {
    var wg sync.WaitGroup

    for i := 1; i <= 3; i++ {
        wg.Add(1)       // increment BEFORE launching goroutine
        go worker(i, &wg)
    }

    wg.Wait()           // blocks until all Done() calls
    fmt.Println("all workers finished")
}
```

### Collecting results from goroutines

```go
func processAll(items []int) []int {
    results := make([]int, len(items))
    var wg sync.WaitGroup

    for i, item := range items {
        wg.Add(1)
        go func(idx, val int) {
            defer wg.Done()
            results[idx] = val * 2 // safe — each goroutine writes to a unique index
        }(i, item)
    }

    wg.Wait()
    return results
}
```

---

## 4. sync.Once

Executes a function **exactly once**, regardless of how many goroutines call it. Used for lazy singleton initialization.

```go
package main

import (
    "fmt"
    "sync"
)

type Database struct {
    url string
}

var (
    db   *Database
    once sync.Once
)

func GetDB() *Database {
    once.Do(func() {
        fmt.Println("initializing database connection") // runs only once
        db = &Database{url: "postgres://localhost/mydb"}
    })
    return db
}

func main() {
    var wg sync.WaitGroup
    for i := 0; i < 5; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            d := GetDB()
            fmt.Println("using db:", d.url)
        }()
    }
    wg.Wait()
    // "initializing database connection" prints exactly once
}
```

> **sync.Once vs init():** `Once` is lazy (runs on first call) and can return values via closure. `init()` always runs at startup and can't be deferred.

---

## 5. sync.Pool

Reuses temporary objects to reduce GC pressure. Objects may be **collected at any GC cycle** — don't use Pool for permanent storage.

```go
package main

import (
    "bytes"
    "fmt"
    "sync"
)

var bufPool = sync.Pool{
    New: func() any {
        return new(bytes.Buffer) // called when pool is empty
    },
}

func processRequest(data string) string {
    buf := bufPool.Get().(*bytes.Buffer) // get from pool
    defer func() {
        buf.Reset()
        bufPool.Put(buf) // return to pool — don't use buf after this
    }()

    buf.WriteString("processed: ")
    buf.WriteString(data)
    return buf.String()
}

func main() {
    fmt.Println(processRequest("hello"))
    fmt.Println(processRequest("world"))
}
```

**When to use sync.Pool:**
- Short-lived objects created frequently (per-request buffers, encoder instances)
- Profiling shows allocation is a bottleneck
- `encoding/json`, `fmt`, `net/http` use it internally

---

## 6. sync/atomic

Atomic operations are **lock-free** — the CPU guarantees their atomicity without a mutex. Use for simple counters and flags where mutex overhead is measurable.

```go
package main

import (
    "fmt"
    "sync"
    "sync/atomic"
)

func main() {
    var counter int64

    var wg sync.WaitGroup
    for i := 0; i < 1000; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            atomic.AddInt64(&counter, 1) // atomic increment — no mutex needed
        }()
    }
    wg.Wait()
    fmt.Println(atomic.LoadInt64(&counter)) // always 1000
}
```

### atomic.Value — store/load any value atomically

```go
package main

import (
    "fmt"
    "sync/atomic"
)

type Config struct {
    MaxConnections int
    Timeout        int
}

var currentConfig atomic.Value

func updateConfig(cfg Config) {
    currentConfig.Store(cfg) // atomic write
}

func getConfig() Config {
    return currentConfig.Load().(Config) // atomic read
}

func main() {
    currentConfig.Store(Config{MaxConnections: 100, Timeout: 30})

    cfg := getConfig()
    fmt.Println(cfg.MaxConnections) // 100

    updateConfig(Config{MaxConnections: 200, Timeout: 60})
    fmt.Println(getConfig().MaxConnections) // 200
}
```

| | Mutex | atomic |
|---|---|---|
| Use for | Complex state, multiple fields | Single counters, flags, pointers |
| Overhead | Higher (kernel involvement possible) | Lower (CPU instruction) |
| Composability | Can protect any code block | Only predefined operations |

---

## Interview Questions

<details>
<summary><b>When would you use sync.Mutex vs a channel?</b></summary>

Use a **mutex** when protecting shared state — a cache, counter, or map that multiple goroutines read/write. Use a **channel** when transferring ownership of data or signaling events between goroutines. The Go proverb: "Share memory by communicating" — but in practice, a mutex is often simpler and more performant for pure state protection.

</details>

<details>
<summary><b>What is sync.RWMutex and when does it give a performance advantage?</b></summary>

`RWMutex` allows many goroutines to hold a read lock simultaneously, but only one goroutine can hold a write lock. It's faster than `Mutex` when reads vastly outnumber writes (e.g., a read-heavy cache). It can be slower under heavy write contention due to coordination overhead — always benchmark.

</details>

<details>
<summary><b>What is sync.Once used for? Give a real-world use case.</b></summary>

`sync.Once` executes a function exactly once, thread-safely. Use cases: lazy singleton initialization (database connection, config loading), one-time setup work triggered at first use rather than at startup. It's safer than `init()` because it's lazy and can close over values.

</details>

<details>
<summary><b>What is sync.Pool? When would you use it?</b></summary>

`sync.Pool` is a cache of reusable objects that reduces GC pressure by avoiding repeated allocations. You `Get()` an object (or a new one is created via `New`), use it, then `Put()` it back. Use it for frequently allocated temporary objects like byte buffers or encoder instances. Don't use it for permanent storage — objects can be evicted at any GC cycle.

</details>

<details>
<summary><b>What is the difference between sync/atomic operations and a mutex?</b></summary>

Atomic operations (`AddInt64`, `LoadInt64`, `StoreInt64`, `CompareAndSwap`) are single CPU instructions guaranteed to be atomic by hardware — no locking, no goroutine blocking. They work only on single primitive values. A mutex protects arbitrary code blocks and multiple fields but involves goroutine scheduling overhead. Use atomics for simple counters/flags where profiling shows mutex overhead matters; use a mutex for everything else.

</details>
