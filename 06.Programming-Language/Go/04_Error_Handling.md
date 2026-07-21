# 04 — Error Handling

> **Topics:** error interface · errors.Is/As · wrapping with %w · custom errors · sentinel errors

---

## Table of Contents
1. [Custom Error Types](#1-custom-error-types)
2. [Error Wrapping — %w vs %v](#2-error-wrapping--w-vs-v)
3. [errors.Is vs errors.As](#3-errorsis-vs-errorsas)
4. [Sentinel Errors](#4-sentinel-errors)
5. [Reducing if err != nil Repetition](#5-reducing-if-err--nil-repetition)
6. [Interview Questions](#interview-questions)

---

## 1. Custom Error Types

The `error` interface has one method: `Error() string`. Any type that implements it is an error.

```go
package main

import (
    "errors"
    "fmt"
)

// Simple custom error type
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation failed on field %q: %s", e.Field, e.Message)
}

// HTTP-style error with a status code
type HTTPError struct {
    Code    int
    Message string
}

func (e *HTTPError) Error() string {
    return fmt.Sprintf("HTTP %d: %s", e.Code, e.Message)
}

func validateAge(age int) error {
    if age < 0 {
        return &ValidationError{Field: "age", Message: "must be non-negative"}
    }
    if age > 150 {
        return &ValidationError{Field: "age", Message: "unrealistically large"}
    }
    return nil
}

func main() {
    if err := validateAge(-1); err != nil {
        fmt.Println(err)
        // validation failed on field "age": must be non-negative

        var ve *ValidationError
        if errors.As(err, &ve) {
            fmt.Println("field:", ve.Field)
        }
    }
}
```

---

## 2. Error Wrapping — %w vs %v

| Verb | Effect | Unwrappable? |
|---|---|---|
| `fmt.Errorf("... %w", err)` | Wraps err — preserves the chain | ✅ Yes |
| `fmt.Errorf("... %v", err)` | Formats as string — loses the original | ❌ No |

```go
package main

import (
    "errors"
    "fmt"
)

var ErrNotFound = errors.New("not found")

func getUser(id int) error {
    if id == 0 {
        // Wrap with %w — caller can unwrap and inspect
        return fmt.Errorf("getUser(id=%d): %w", id, ErrNotFound)
    }
    return nil
}

func main() {
    err := getUser(0)
    fmt.Println(err)
    // getUser(id=0): not found

    // errors.Is works through the chain because we used %w
    fmt.Println(errors.Is(err, ErrNotFound)) // true

    // With %v the chain is broken
    err2 := fmt.Errorf("getUser: %v", ErrNotFound)
    fmt.Println(errors.Is(err2, ErrNotFound)) // false
}
```

### Unwrap chain

```go
// errors.Unwrap returns the next error in the chain
wrapped := fmt.Errorf("layer2: %w", fmt.Errorf("layer1: %w", ErrNotFound))

fmt.Println(errors.Unwrap(wrapped))          // layer1: not found
fmt.Println(errors.Is(wrapped, ErrNotFound)) // true — walks the full chain
```

---

## 3. errors.Is vs errors.As

| Function | Purpose | Matches by |
|---|---|---|
| `errors.Is(err, target)` | Check if a specific **value** is in the chain | Value equality |
| `errors.As(err, &target)` | Extract a specific **type** from the chain | Type match |

```go
package main

import (
    "errors"
    "fmt"
)

type DBError struct {
    Code    int
    Message string
}

func (e *DBError) Error() string {
    return fmt.Sprintf("db error %d: %s", e.Code, e.Message)
}

var ErrTimeout = errors.New("timeout")

func query() error {
    dbErr := &DBError{Code: 503, Message: "connection refused"}
    return fmt.Errorf("query failed: %w", dbErr)
}

func ping() error {
    return fmt.Errorf("ping: %w", ErrTimeout)
}

func main() {
    // errors.Is — match a sentinel value
    err := ping()
    if errors.Is(err, ErrTimeout) {
        fmt.Println("operation timed out") // prints this
    }

    // errors.As — extract a typed error from the chain
    err2 := query()
    var dbErr *DBError
    if errors.As(err2, &dbErr) {
        fmt.Printf("DB error code: %d\n", dbErr.Code) // DB error code: 503
    }
}
```

---

## 4. Sentinel Errors

Sentinel errors are **package-level error variables** used as known, comparable values.

```go
package main

import (
    "errors"
    "fmt"
)

// Sentinel errors — exported so callers can check against them
var (
    ErrNotFound   = errors.New("not found")
    ErrPermission = errors.New("permission denied")
    ErrConflict   = errors.New("conflict")
)

func findItem(id int) error {
    if id <= 0 {
        return fmt.Errorf("findItem: %w", ErrNotFound)
    }
    return nil
}

func main() {
    err := findItem(0)

    switch {
    case errors.Is(err, ErrNotFound):
        fmt.Println("item does not exist")
    case errors.Is(err, ErrPermission):
        fmt.Println("access denied")
    default:
        fmt.Println("unexpected error:", err)
    }
}
```

**When to avoid sentinels:** They create package-level coupling — callers import your package just to compare errors. For richer context (status codes, fields) prefer custom error types with `errors.As`.

---

## 5. Reducing if err != nil Repetition

### Pattern 1 — errWriter (accumulate first error)

```go
package main

import (
    "fmt"
    "io"
)

// errWriter stops writing after the first error
type errWriter struct {
    w   io.Writer
    err error
}

func (ew *errWriter) write(p []byte) {
    if ew.err != nil {
        return // skip all subsequent writes
    }
    _, ew.err = ew.w.Write(p)
}

func writeAll(w io.Writer, chunks ...[]byte) error {
    ew := &errWriter{w: w}
    for _, chunk := range chunks {
        ew.write(chunk)
    }
    return ew.err
}
```

### Pattern 2 — functional pipeline with early exit

```go
type Result struct {
    Value int
    Err   error
}

func step1(v int) Result {
    if v < 0 {
        return Result{Err: fmt.Errorf("step1: negative value %d", v)}
    }
    return Result{Value: v * 2}
}

func step2(r Result) Result {
    if r.Err != nil {
        return r // propagate error without doing work
    }
    return Result{Value: r.Value + 10}
}

func process(v int) (int, error) {
    r := step2(step1(v))
    return r.Value, r.Err
}
```

### Pattern 3 — named return + defer for cleanup chains

```go
func openAndProcess(path string) (err error) {
    f, err := os.Open(path)
    if err != nil {
        return fmt.Errorf("open: %w", err)
    }
    defer func() {
        if cerr := f.Close(); cerr != nil && err == nil {
            err = fmt.Errorf("close: %w", cerr) // capture close error only if no prior error
        }
    }()

    // ... process f
    return nil
}
```

---

## Interview Questions

<details>
<summary><b>How do you create a custom error type in Go?</b></summary>

Implement the `error` interface by adding an `Error() string` method to any type. Use a struct to carry extra context (field name, HTTP status code, etc.). Return a pointer (`*MyError`) so `errors.As` can match it correctly.

</details>

<details>
<summary><b>What is the difference between fmt.Errorf("... %w", err) and %v?</b></summary>

`%w` **wraps** the error, preserving the chain so `errors.Is` and `errors.As` can inspect it. `%v` formats the error as a plain string, destroying the chain — callers can read the message but cannot programmatically match the original error type or value.

</details>

<details>
<summary><b>What is errors.Is() vs errors.As()? When do you use each?</b></summary>

`errors.Is(err, target)` checks if a specific **value** appears anywhere in the error chain — ideal for sentinel errors (`ErrNotFound`, `io.EOF`). `errors.As(err, &target)` finds the first error in the chain that matches a **type** and sets `target` to it — use it when you need to access fields on a custom error struct.

</details>

<details>
<summary><b>What are sentinel errors and when should you avoid them?</b></summary>

Sentinel errors are package-level `var` values (`var ErrNotFound = errors.New("not found")`). They're good for well-known, stable conditions (`io.EOF`, `sql.ErrNoRows`). Avoid them when you need to carry context (status codes, field names) — use custom types instead. They also create import coupling between packages.

</details>

<details>
<summary><b>How would you handle errors in a chain of function calls without repeating if err != nil?</b></summary>

Three common patterns: (1) **errWriter** — an accumulator that stores the first error and skips subsequent operations, used in `encoding/binary`. (2) **Pipeline** — each step returns a result struct; downstream steps check and forward the error without doing work. (3) **Named return + defer** — especially for cleanup (e.g. closing resources and capturing close errors without overwriting an earlier error).

</details>
