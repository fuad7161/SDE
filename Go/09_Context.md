# 09 — Context Package

> **Topics:** cancellation · timeout · deadline · value propagation · context in HTTP/DB chains

---

## Table of Contents
1. [Main Use Cases](#1-main-use-cases)
2. [WithCancel, WithTimeout, WithDeadline](#2-withcancel-withtimeout-withdeadline)
3. [Context Values — Rules and Pitfalls](#3-context-values--rules-and-pitfalls)
4. [Propagating Cancellation — HTTP to DB](#4-propagating-cancellation--http-to-db)
5. [Ignoring Context](#5-ignoring-context)
6. [Interview Questions](#interview-questions)

---

## 1. Main Use Cases

`context.Context` carries **deadlines, cancellation signals, and request-scoped values** across API boundaries and between goroutines.

```
HTTP Request
    │
    ├─ context.WithTimeout(5s)
    │       │
    │       ├─ Service layer     (reads ctx.Done())
    │       │       │
    │       │       └─ DB query  (passes ctx to sql driver)
    │       │
    │       └─ External API call (passes ctx to http.NewRequestWithContext)
    │
    └─ If client disconnects → ctx.Done() fires → all nested work cancels
```

```go
package main

import (
    "context"
    "fmt"
    "time"
)

func doWork(ctx context.Context, name string) error {
    select {
    case <-time.After(2 * time.Second): // simulate work
        fmt.Println(name, "done")
        return nil
    case <-ctx.Done():
        fmt.Printf("%s cancelled: %v\n", name, ctx.Err())
        return ctx.Err()
    }
}

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
    defer cancel()

    if err := doWork(ctx, "task-1"); err != nil {
        fmt.Println("error:", err)
    }
    // task-1 cancelled: context deadline exceeded
}
```

---

## 2. WithCancel, WithTimeout, WithDeadline

| Constructor | Cancels when | Returns |
|---|---|---|
| `WithCancel(parent)` | `cancel()` is called manually | `ctx, cancelFunc` |
| `WithTimeout(parent, d)` | `d` duration elapses OR `cancel()` called | `ctx, cancelFunc` |
| `WithDeadline(parent, t)` | absolute `time.Time t` is reached OR `cancel()` called | `ctx, cancelFunc` |

> **Always call the cancel function** (via `defer cancel()`) to release resources even if the deadline fires first.

```go
package main

import (
    "context"
    "fmt"
    "time"
)

func main() {
    // WithCancel — manual control
    ctx1, cancel1 := context.WithCancel(context.Background())
    go func() {
        time.Sleep(500 * time.Millisecond)
        cancel1() // signal cancellation
    }()
    <-ctx1.Done()
    fmt.Println("ctx1:", ctx1.Err()) // context canceled

    // WithTimeout — fires after duration
    ctx2, cancel2 := context.WithTimeout(context.Background(), 100*time.Millisecond)
    defer cancel2()
    time.Sleep(200 * time.Millisecond)
    fmt.Println("ctx2:", ctx2.Err()) // context deadline exceeded

    // WithDeadline — fires at absolute time
    deadline := time.Now().Add(100 * time.Millisecond)
    ctx3, cancel3 := context.WithDeadline(context.Background(), deadline)
    defer cancel3()
    <-ctx3.Done()
    fmt.Println("ctx3:", ctx3.Err()) // context deadline exceeded

    // Check remaining time
    if dl, ok := ctx3.Deadline(); ok {
        fmt.Println("deadline was:", dl.Format(time.RFC3339))
    }
}
```

---

## 3. Context Values — Rules and Pitfalls

Use `context.WithValue` only for **request-scoped data** that crosses API boundaries — not for optional function parameters.

```go
package main

import (
    "context"
    "fmt"
)

// ✅ Use unexported type for context keys to avoid collisions
type contextKey string

const (
    keyRequestID contextKey = "requestID"
    keyUserID    contextKey = "userID"
)

func withRequestID(ctx context.Context, id string) context.Context {
    return context.WithValue(ctx, keyRequestID, id)
}

func getRequestID(ctx context.Context) (string, bool) {
    id, ok := ctx.Value(keyRequestID).(string)
    return id, ok
}

func handler(ctx context.Context) {
    if id, ok := getRequestID(ctx); ok {
        fmt.Println("request ID:", id)
    }
}

func main() {
    ctx := context.Background()
    ctx = withRequestID(ctx, "req-abc-123")
    ctx = context.WithValue(ctx, keyUserID, "user-42")

    handler(ctx)
}
```

### What NOT to put in context

```go
// ❌ Optional parameters — use function arguments instead
ctx = context.WithValue(ctx, "db_timeout", 30)

// ❌ Large objects — context is passed by value, object lives on heap
ctx = context.WithValue(ctx, "payload", largeStruct)

// ❌ Mutable state — context values should be read-only
ctx = context.WithValue(ctx, "counter", &myCounter)
```

---

## 4. Propagating Cancellation — HTTP to DB

```go
package main

import (
    "context"
    "database/sql"
    "fmt"
    "net/http"
    "time"
)

type UserService struct {
    db *sql.DB
}

// DB query — passes ctx so the driver cancels the query if ctx is done
func (s *UserService) GetUser(ctx context.Context, id string) (string, error) {
    var name string
    err := s.db.QueryRowContext(ctx,
        "SELECT name FROM users WHERE id = $1", id,
    ).Scan(&name)
    if err != nil {
        return "", fmt.Errorf("GetUser: %w", err)
    }
    return name, nil
}

// HTTP handler — wraps request context with a per-handler timeout
func (s *UserService) HandleGetUser(w http.ResponseWriter, r *http.Request) {
    // Inherit the request context (already cancelled if client disconnects)
    ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
    defer cancel()

    userID := r.URL.Query().Get("id")
    name, err := s.GetUser(ctx, userID)
    if err != nil {
        if ctx.Err() != nil {
            http.Error(w, "request cancelled or timed out", http.StatusGatewayTimeout)
            return
        }
        http.Error(w, "user not found", http.StatusNotFound)
        return
    }
    fmt.Fprintln(w, name)
}
```

---

## 5. Ignoring Context

```go
package main

import (
    "context"
    "fmt"
    "time"
)

// ❌ Ignores context — keeps running after caller cancels
func ignoringCtx(ctx context.Context) string {
    time.Sleep(5 * time.Second) // runs full 5s even if ctx is cancelled after 1s
    return "result"
}

// ✅ Respects context — returns early when cancelled
func respectingCtx(ctx context.Context) (string, error) {
    resultCh := make(chan string, 1)

    go func() {
        time.Sleep(5 * time.Second) // long work
        resultCh <- "result"
    }()

    select {
    case result := <-resultCh:
        return result, nil
    case <-ctx.Done():
        return "", ctx.Err() // return immediately on cancellation
    }
}

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
    defer cancel()

    result, err := respectingCtx(ctx)
    if err != nil {
        fmt.Println("cancelled:", err) // cancelled: context deadline exceeded
        return
    }
    fmt.Println(result)
}
```

---

## Interview Questions

<details>
<summary><b>What are the main use cases for context.Context?</b></summary>

Three main uses: (1) **Cancellation** — signal goroutines to stop work (client disconnect, timeout). (2) **Deadline/timeout** — bound how long an operation can run. (3) **Request-scoped values** — pass trace IDs, auth tokens across API boundaries without changing function signatures. Context is always passed as the first argument by convention.

</details>

<details>
<summary><b>What is the difference between context.WithCancel, context.WithTimeout, and context.WithDeadline?</b></summary>

`WithCancel` gives a manual cancel function — you call `cancel()` when done. `WithTimeout` sets a relative duration after which the context auto-cancels. `WithDeadline` sets an absolute `time.Time`. All three return a cancel function that must be called (via `defer`) to release timer resources, even if the deadline fires first.

</details>

<details>
<summary><b>Why should you not store large objects in a context value?</b></summary>

Context is passed by value through function calls and across goroutines — but the values stored in it live on the heap. Large objects increase memory pressure. More importantly, context values are untyped (`any`) and retrieved via type assertion, which bypasses the compiler's type safety. Large or mutable state should be passed as explicit function arguments.

</details>

<details>
<summary><b>How do you propagate cancellation from an HTTP handler down to a DB call?</b></summary>

Use `r.Context()` from the HTTP request — the framework cancels it when the client disconnects. Optionally wrap with `context.WithTimeout` for a per-handler budget. Pass this context to every downstream call: `db.QueryRowContext(ctx, ...)`, `http.NewRequestWithContext(ctx, ...)`. The DB driver, HTTP client, and gRPC stubs all honor context cancellation.

</details>

<details>
<summary><b>What happens if you ignore a context and the caller cancels?</b></summary>

The goroutine continues running — wasting CPU, memory, and potentially holding DB connections or locks. Responses sent to a cancelled request are discarded. In high-traffic services this leads to resource exhaustion. Always check `ctx.Done()` in long-running loops or select statements, and pass context to all I/O operations.

</details>
