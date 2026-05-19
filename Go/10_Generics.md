# 10 — Generics (Go 1.18+)

> **Topics:** type parameters · constraints · comparable · generic methods limitations · interface{} vs generics

---

## Table of Contents
1. [What Problem Generics Solve](#1-what-problem-generics-solve)
2. [Type Constraints](#2-type-constraints)
3. [The comparable Constraint](#3-the-comparable-constraint)
4. [Generics with Struct Methods — Limitations](#4-generics-with-struct-methods--limitations)
5. [When to Prefer interface{} over Generics](#5-when-to-prefer-interface-over-generics)
6. [Interview Questions](#interview-questions)

---

## 1. What Problem Generics Solve

Before generics, reusable algorithms required either:
- Code duplication (one function per type)
- `interface{}` (loses type safety, requires type assertions)

```go
// ❌ Before generics — duplicated or type-unsafe
func SumInts(nums []int) int {
    var total int
    for _, n := range nums { total += n }
    return total
}
func SumFloats(nums []float64) float64 {
    var total float64
    for _, n := range nums { total += n }
    return total
}

// ✅ With generics — one function, type-safe
package main

import "fmt"

type Number interface {
    int | int8 | int16 | int32 | int64 |
    float32 | float64
}

func Sum[T Number](nums []T) T {
    var total T
    for _, n := range nums {
        total += n
    }
    return total
}

func main() {
    fmt.Println(Sum([]int{1, 2, 3}))         // 6
    fmt.Println(Sum([]float64{1.1, 2.2}))    // 3.3000...
}
```

### Generic data structures

```go
package main

import "fmt"

// Stack works for any type
type Stack[T any] struct {
    items []T
}

func (s *Stack[T]) Push(item T) {
    s.items = append(s.items, item)
}

func (s *Stack[T]) Pop() (T, bool) {
    var zero T
    if len(s.items) == 0 {
        return zero, false
    }
    top := s.items[len(s.items)-1]
    s.items = s.items[:len(s.items)-1]
    return top, true
}

func (s *Stack[T]) Len() int { return len(s.items) }

func main() {
    s := &Stack[int]{}
    s.Push(1)
    s.Push(2)
    s.Push(3)

    for s.Len() > 0 {
        v, _ := s.Pop()
        fmt.Println(v) // 3, 2, 1
    }
}
```

---

## 2. Type Constraints

Constraints define what types are permitted for a type parameter. Any interface can be used as a constraint.

```go
package main

import (
    "fmt"
    "strings"
)

// Union constraint — only these specific types
type Ordered interface {
    int | float64 | string
}

func Min[T Ordered](a, b T) T {
    if a < b {
        return a
    }
    return b
}

// Interface constraint — must have a method
type Stringer interface {
    String() string
}

func PrintAll[T Stringer](items []T) {
    for _, item := range items {
        fmt.Println(item.String())
    }
}

// ~ tilde prefix — includes types whose underlying type matches
type Integer interface {
    ~int | ~int32 | ~int64 // includes custom types like "type UserID int"
}

type UserID int
type PostID int64

func Double[T Integer](v T) T { return v * 2 }

func main() {
    fmt.Println(Min(3, 5))               // 3
    fmt.Println(Min("apple", "banana"))  // apple

    var uid UserID = 10
    fmt.Println(Double(uid)) // 20 — works because ~int includes UserID

    _ = strings.ToUpper // just importing strings
}
```

### Built-in constraints from `golang.org/x/exp/constraints`

```go
import "golang.org/x/exp/constraints"

func Max[T constraints.Ordered](a, b T) T {
    if a > b { return a }
    return b
}
```

---

## 3. The comparable Constraint

`comparable` is a built-in constraint satisfied by any type that supports `==` and `!=` (can be used as a map key).

```go
package main

import "fmt"

// Generic Set — works with any comparable type
type Set[T comparable] struct {
    items map[T]struct{}
}

func NewSet[T comparable]() *Set[T] {
    return &Set[T]{items: make(map[T]struct{})}
}

func (s *Set[T]) Add(v T)          { s.items[v] = struct{}{} }
func (s *Set[T]) Contains(v T) bool { _, ok := s.items[v]; return ok }
func (s *Set[T]) Remove(v T)       { delete(s.items, v) }
func (s *Set[T]) Len() int         { return len(s.items) }

func main() {
    intSet := NewSet[int]()
    intSet.Add(1)
    intSet.Add(2)
    intSet.Add(1)                    // duplicate ignored
    fmt.Println(intSet.Len())        // 2
    fmt.Println(intSet.Contains(2))  // true

    strSet := NewSet[string]()
    strSet.Add("go")
    strSet.Add("generics")
    fmt.Println(strSet.Contains("go")) // true
}
```

Types NOT comparable: slices, maps, functions — these cannot be used with `comparable`.

---

## 4. Generics with Struct Methods — Limitations

**You cannot add new type parameters on a method** — only the type parameters declared on the struct itself are available.

```go
package main

import "fmt"

type Pair[T any] struct {
    First, Second T
}

// ✅ Method can use T from the struct declaration
func (p Pair[T]) Swap() Pair[T] {
    return Pair[T]{First: p.Second, Second: p.First}
}

// ❌ Cannot add a new type parameter on a method
// func (p Pair[T]) Convert[U any]() Pair[U] { ... } // compile error

// ✅ Workaround — use a package-level generic function
func Convert[T, U any](p Pair[T], fn func(T) U) Pair[U] {
    return Pair[U]{First: fn(p.First), Second: fn(p.Second)}
}

func main() {
    p := Pair[int]{1, 2}
    fmt.Println(p.Swap())  // {2 1}

    s := Convert(p, func(n int) string {
        return fmt.Sprintf("%d", n)
    })
    fmt.Println(s) // {1 2}
}
```

---

## 5. When to Prefer interface{} over Generics

| Scenario | Prefer |
|---|---|
| Algorithm works the same for many types (sort, filter, map) | Generics |
| Heterogeneous collection (store different types in one list) | `interface{}` / `any` |
| Runtime polymorphism (behavior differs per type) | Interface |
| Type is known only at runtime (JSON unmarshalling) | `interface{}` |
| Performance-critical, type-safe reuse | Generics |

```go
// Generics — homogeneous, type-safe container
func Filter[T any](s []T, fn func(T) bool) []T {
    var result []T
    for _, v := range s {
        if fn(v) { result = append(result, v) }
    }
    return result
}

// interface{} — heterogeneous storage (e.g., JSON decoded object)
var decoded map[string]any
json.Unmarshal(data, &decoded) // types only known at runtime
```

---

## Interview Questions

<details>
<summary><b>What problem do generics solve in Go?</b></summary>

Generics eliminate the need to duplicate code for different types or use `interface{}` with type assertions. Before generics, writing a type-safe `Stack[int]` and `Stack[string]` required either two separate implementations or an `interface{}` implementation that lost compile-time type safety. Generics let you write the algorithm once and have the compiler specialize it per type.

</details>

<details>
<summary><b>What is a type constraint? How do you define one using an interface?</b></summary>

A constraint is an interface that restricts which types can be used as a type argument. It can specify a set of methods (behavioral constraint) or a union of specific types (`int | float64 | string`). The `~` prefix includes types whose underlying type matches (`~int` matches `type UserID int`). The built-in constraints `any` (no restriction) and `comparable` (supports ==) are the most common.

</details>

<details>
<summary><b>What is the comparable constraint used for?</b></summary>

`comparable` restricts the type parameter to types that support `==` and `!=`, which is required to use a value as a map key. It's used when building generic maps, sets, or any structure that needs to compare values for equality. Slices, maps, and functions are not comparable and cannot satisfy this constraint.

</details>

<details>
<summary><b>Can you use generics with methods on a struct? What are the limitations?</b></summary>

Yes — methods on a generic struct can use the type parameters declared on the struct. The limitation is that **methods cannot introduce their own new type parameters** — only functions at the package level can do that. The workaround is to implement the additional generic behavior as a package-level function that takes the struct as an argument.

</details>

<details>
<summary><b>When would you still prefer interface{} over generics?</b></summary>

When you need **heterogeneous** collections (storing different types in one slice), **runtime polymorphism** (behavior varies per type via method dispatch), or when types are only known at runtime (JSON decoding, reflection). Generics excel at type-safe, homogeneous algorithms (filter, map, reduce, containers). If different types need different behavior, use an interface; if they all need the same algorithm, use generics.

</details>
