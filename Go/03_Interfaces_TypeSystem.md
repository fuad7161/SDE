# 03 — Interfaces & Type System

> **Topics:** implicit implementation · empty interface · type assertion · type switch · embedding · nil interface trap

---

## Table of Contents
1. [Implicit Interface Implementation](#1-implicit-interface-implementation)
2. [interface{} and any](#2-interface-and-any)
3. [Type Assertion vs Type Switch](#3-type-assertion-vs-type-switch)
4. [The Nil Interface Trap](#4-the-nil-interface-trap)
5. [Interface Embedding](#5-interface-embedding)
6. [Interview Questions](#interview-questions)

---

## 1. Implicit Interface Implementation

In Go, interfaces are satisfied **implicitly** — no `implements` keyword. Any type with the required methods satisfies the interface.

```go
package main

import (
    "fmt"
    "math"
)

type Shape interface {
    Area() float64
    Perimeter() float64
}

type Circle struct{ Radius float64 }
type Rectangle struct{ Width, Height float64 }

func (c Circle) Area() float64      { return math.Pi * c.Radius * c.Radius }
func (c Circle) Perimeter() float64 { return 2 * math.Pi * c.Radius }

func (r Rectangle) Area() float64      { return r.Width * r.Height }
func (r Rectangle) Perimeter() float64 { return 2 * (r.Width + r.Height) }

// Both satisfy Shape without declaring so
func printInfo(s Shape) {
    fmt.Printf("Area: %.2f  Perimeter: %.2f\n", s.Area(), s.Perimeter())
}

func main() {
    printInfo(Circle{Radius: 5})
    printInfo(Rectangle{Width: 3, Height: 4})
}
```

**vs Java/C#:** those require explicit `implements`/`:` declarations, creating tight coupling to the interface definition. Go's implicit approach lets you define interfaces *after* the concrete types exist, or in a different package entirely.

---

## 2. interface{} and any

`any` is an **alias** for `interface{}` introduced in Go 1.18. They are identical.

```go
var v1 interface{} = 42
var v2 any = 42   // exact same type

fmt.Println(v1 == v2) // true
```

An `interface{}` / `any` value holds two things internally:
- A **type descriptor** (what type is stored)
- A **pointer to the value**

```go
package main

import "fmt"

func describe(v any) {
    fmt.Printf("value: %v  type: %T\n", v, v)
}

func main() {
    describe(42)
    describe("hello")
    describe([]int{1, 2, 3})
    describe(nil)
}
// value: 42       type: int
// value: hello    type: string
// value: [1 2 3]  type: []int
// value: <nil>    type: <nil>
```

---

## 3. Type Assertion vs Type Switch

### Type assertion — extract a specific type

```go
var i any = "hello"

// Safe assertion — won't panic
s, ok := i.(string)
if ok {
    fmt.Println("string:", s) // string: hello
}

// Unsafe assertion — panics if wrong type
s2 := i.(string)  // ok if i is string
n  := i.(int)     // panic: interface conversion: interface {} is string, not int
```

### Type switch — handle multiple types

```go
package main

import "fmt"

func describe(i any) string {
    switch v := i.(type) {
    case int:
        return fmt.Sprintf("int: %d", v)
    case string:
        return fmt.Sprintf("string: %q (len=%d)", v, len(v))
    case bool:
        return fmt.Sprintf("bool: %t", v)
    case []int:
        return fmt.Sprintf("[]int with %d elements", len(v))
    case nil:
        return "nil"
    default:
        return fmt.Sprintf("unknown type: %T", v)
    }
}

func main() {
    fmt.Println(describe(42))
    fmt.Println(describe("Go"))
    fmt.Println(describe(true))
    fmt.Println(describe(nil))
}
```

---

## 4. The Nil Interface Trap

A Go interface value is nil **only when both its type and value are nil**. This is one of the most common Go gotchas.

```go
package main

import "fmt"

type MyError struct{ msg string }

func (e *MyError) Error() string { return e.msg }

// ❌ BUG — returns a non-nil interface even when there's no error
func mightFail(fail bool) error {
    var err *MyError // typed nil pointer
    if fail {
        err = &MyError{"something went wrong"}
    }
    return err // interface wraps (*MyError, nil) — NOT a nil interface!
}

func main() {
    err := mightFail(false)
    if err != nil {
        fmt.Println("BUG: this prints even though no error occurred!")
        fmt.Printf("type: %T  value: %v\n", err, err)
        // type: *MyError  value: <nil>
    }
}
```

### Why it happens

An interface holds `(type, value)`. Returning a `*MyError` nil gives `(*MyError, nil)` — type is set, so the interface is **not nil**.

```go
// ✅ FIX — return untyped nil
func mightFailFixed(fail bool) error {
    if fail {
        return &MyError{"something went wrong"}
    }
    return nil // untyped nil → (nil, nil) → truly nil interface
}
```

---

## 5. Interface Embedding

Interfaces can embed other interfaces to compose larger contracts.

```go
package main

import "fmt"

type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}

// ReadWriter embeds both — satisfied by any type implementing both methods
type ReadWriter interface {
    Reader
    Writer
}

// ---- Practical example: layered service interfaces ----

type UserStore interface {
    GetUser(id string) (*User, error)
    SaveUser(u *User) error
}

type UserCache interface {
    CachedUser(id string) (*User, bool)
    SetCache(u *User)
}

// CachingUserStore requires both capabilities
type CachingUserStore interface {
    UserStore
    UserCache
}

type User struct{ ID, Name string }

// MockStore satisfies CachingUserStore
type MockStore struct{ data map[string]*User }

func (m *MockStore) GetUser(id string) (*User, error) {
    u, ok := m.data[id]
    if !ok {
        return nil, fmt.Errorf("user %s not found", id)
    }
    return u, nil
}
func (m *MockStore) SaveUser(u *User) error     { m.data[u.ID] = u; return nil }
func (m *MockStore) CachedUser(id string) (*User, bool) { u, ok := m.data[id]; return u, ok }
func (m *MockStore) SetCache(u *User)           { m.data[u.ID] = u }

func main() {
    store := &MockStore{data: make(map[string]*User)}
    store.SaveUser(&User{"1", "Alice"})

    u, _ := store.GetUser("1")
    fmt.Println(u.Name) // Alice
}
```

---

## Interview Questions

<details>
<summary><b>How does Go implement interfaces differently from Java or C#?</b></summary>

Go uses **structural (implicit) typing** — a type satisfies an interface simply by having the required methods, without any declaration. Java/C# use **nominal typing** — you must explicitly write `implements`/`:`. Go's approach decouples the interface definition from the concrete type; you can define an interface after the fact, even in a different package, and existing types automatically satisfy it.

</details>

<details>
<summary><b>What is the difference between interface{} and any?</b></summary>

They are identical. `any` is a type alias for `interface{}` introduced in Go 1.18 for readability. Both can hold a value of any type. Internally an interface value is a `(type, value)` pair; when both are nil, the interface itself is nil.

</details>

<details>
<summary><b>What is the difference between a type assertion and a type switch?</b></summary>

A **type assertion** (`v, ok := i.(T)`) extracts one specific type from an interface. A **type switch** (`switch v := i.(type)`) branches on multiple types and is idiomatic when you need to handle several cases. Type assertion panics if the assertion fails and the two-value form is not used.

</details>

<details>
<summary><b>Can a nil interface equal a nil pointer? Explain the nil interface trap.</b></summary>

No. A nil interface has `(nil, nil)` internally. A nil pointer stored in an interface produces `(*T, nil)` — the type descriptor is set, so `err != nil` is **true** even though the pointer is nil. The fix is to return an untyped `nil` directly rather than returning a typed nil variable.

</details>

<details>
<summary><b>What is interface embedding? Give a practical example.</b></summary>

Interface embedding lets you compose multiple interfaces into one. Example: `io.ReadWriter` embeds `io.Reader` and `io.Writer`. Any type that implements both `Read` and `Write` automatically satisfies `ReadWriter`. This is the idiomatic way to build layered contracts without code duplication.

</details>
