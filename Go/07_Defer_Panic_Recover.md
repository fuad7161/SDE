# 07 — defer, panic, and recover

> **Topics:** defer execution order · defer with named returns · panic/recover pattern · performance cost

---

## Table of Contents
1. [defer — Execution Order](#1-defer--execution-order)
2. [defer and Named Return Values](#2-defer-and-named-return-values)
3. [panic and recover](#3-panic-and-recover)
4. [Typical recover Pattern](#4-typical-recover-pattern)
5. [Performance Cost of defer](#5-performance-cost-of-defer)
6. [Interview Questions](#interview-questions)

---

## 1. defer — Execution Order

`defer` pushes a function call onto a **LIFO stack**. All deferred calls execute when the surrounding function returns (normally or via panic), in **reverse order** of declaration.

```go
package main

import "fmt"

func main() {
    defer fmt.Println("first deferred")  // runs last
    defer fmt.Println("second deferred") // runs second
    defer fmt.Println("third deferred")  // runs first
    fmt.Println("main body")
}
// Output:
// main body
// third deferred
// second deferred
// first deferred
```

### Arguments are evaluated immediately, not at call time

```go
func main() {
    x := 10
    defer fmt.Println("deferred x =", x) // captures x = 10 right now

    x = 99
    fmt.Println("current x =", x)
}
// Output:
// current x = 99
// deferred x = 10   ← not 99
```

### Classic defer use — cleanup

```go
package main

import (
    "fmt"
    "os"
)

func readFile(name string) (string, error) {
    f, err := os.Open(name)
    if err != nil {
        return "", err
    }
    defer f.Close() // guaranteed to run, even if an error occurs below

    buf := make([]byte, 1024)
    n, err := f.Read(buf)
    if err != nil {
        return "", err
    }
    return string(buf[:n]), nil
}
```

---

## 2. defer and Named Return Values

Deferred functions can **read and modify named return values** — the named return variable is in scope.

```go
package main

import "fmt"

// Named return — defer CAN modify it
func divide(a, b float64) (result float64, err error) {
    defer func() {
        if err != nil {
            result = 0 // normalize on error
        }
    }()

    if b == 0 {
        err = fmt.Errorf("division by zero")
        return // named return: result=0, err=set
    }
    result = a / b
    return
}

// Practical: wrap close errors without losing the original error
func writeToFile(name, data string) (err error) {
    f, err := os.Create(name)
    if err != nil {
        return
    }
    defer func() {
        cerr := f.Close()
        if cerr != nil && err == nil {
            err = cerr // capture close error only if no prior error
        }
    }()

    _, err = f.WriteString(data)
    return
}
```

---

## 3. panic and recover

- `panic` stops the current goroutine's normal execution, unwinds the stack running deferred functions, and crashes the program if uncaught.
- `recover` is only useful **inside a deferred function** — it stops the panic and returns the panic value.

```go
package main

import "fmt"

func safeDivide(a, b int) (result int, err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("recovered from panic: %v", r)
        }
    }()

    result = a / b // panics if b == 0 (integer division by zero)
    return
}

func main() {
    result, err := safeDivide(10, 2)
    fmt.Println(result, err) // 5 <nil>

    result, err = safeDivide(10, 0)
    fmt.Println(result, err) // 0 recovered from panic: runtime error: integer divide by zero
}
```

**recover outside a deferred function returns nil:**

```go
func bad() {
    r := recover() // does nothing — not in a defer
    fmt.Println(r) // <nil>
}
```

---

## 4. Typical recover Pattern

Used to convert panics into errors — common in HTTP servers, worker pools, and library code.

```go
package main

import (
    "fmt"
    "log"
    "net/http"
)

// Middleware — catches panics in HTTP handlers, returns 500
func recoveryMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                log.Printf("panic recovered: %v", err)
                http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            }
        }()
        next.ServeHTTP(w, r)
    })
}

// Worker pool — one goroutine panic must not kill the whole pool
func safeWorker(id int, jobs <-chan func()) {
    defer func() {
        if r := recover(); r != nil {
            log.Printf("worker %d recovered: %v", id, r)
            go safeWorker(id, jobs) // restart the worker
        }
    }()
    for job := range jobs {
        job()
    }
}

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        panic("simulated panic")
    })
    fmt.Println("server starting")
    http.ListenAndServe(":8080", recoveryMiddleware(mux))
}
```

---

## 5. Performance Cost of defer

In Go < 1.14, `defer` had significant overhead (~50–100 ns per call). Since **Go 1.14**, `defer` in most cases is **open-coded** (inlined by the compiler), making it near-free for simple, non-looping cases.

```go
package main

import "testing"

func withDefer(t *testing.B) {
    for i := 0; i < t.N; i++ {
        func() {
            defer func() {}() // negligible cost in Go 1.14+
        }()
    }
}

func withoutDefer(t *testing.B) {
    for i := 0; i < t.N; i++ {
        func() {}()
    }
}
```

### When defer IS still costly

```go
// ❌ Defer inside a loop — stack grows with each iteration
// Each iteration pushes a new defer frame
func processItems(items []string) {
    for _, item := range items {
        defer fmt.Println(item) // all run at end of function, not end of iteration
    }
}

// ✅ Use a nested function or explicit cleanup
func processItemsOK(items []string) {
    for _, item := range items {
        func() {
            // defer is scoped to this anonymous function
            defer fmt.Println(item)
            // ... process item
        }()
    }
}
```

---

## Interview Questions

<details>
<summary><b>In what order do multiple defer statements execute?</b></summary>

LIFO — last in, first out. The last `defer` statement reached executes first when the function returns. This mirrors a stack: each `defer` pushes a call; when the function exits, calls are popped and executed in reverse order.

</details>

<details>
<summary><b>Do deferred functions run if a panic occurs?</b></summary>

Yes. When a panic occurs, the runtime unwinds the current goroutine's call stack executing all deferred functions in LIFO order. This is what makes `defer f.Close()` and recovery middleware reliable — they run regardless of how the function exits.

</details>

<details>
<summary><b>How do you recover from a panic? What is the typical pattern?</b></summary>

Call `recover()` **inside a deferred function**. It returns the value passed to `panic` and stops the panic unwinding. Outside a deferred function, `recover()` returns nil and has no effect. Typical pattern: `defer func() { if r := recover(); r != nil { /* handle */ } }()`. Used in HTTP middleware, worker pools, and library boundary code.

</details>

<details>
<summary><b>Can a deferred function modify a named return value?</b></summary>

Yes. Named return variables are in scope for deferred functions. This is used to: (1) add context to an error before returning, (2) capture a `Close()` error without overwriting an earlier error, (3) enforce invariants on the return value.

</details>

<details>
<summary><b>What are the performance costs of using defer in a tight loop?</b></summary>

Since Go 1.14, simple defers are open-coded (nearly free). However, `defer` inside a loop accumulates entries on a heap-allocated defer chain, which has real overhead. More importantly, deferred calls inside a loop don't execute until the enclosing function returns — which is usually wrong semantically. The fix is to wrap loop body logic in an anonymous function so defer fires at end of each iteration.

</details>
