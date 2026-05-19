# 11 — Testing

> **Topics:** table-driven tests · benchmarks · mocking · t.Parallel · testing unexported functions

---

## Table of Contents
1. [Table-Driven Tests](#1-table-driven-tests)
2. [Benchmarks](#2-benchmarks)
3. [Mocking Without a Framework](#3-mocking-without-a-framework)
4. [t.Parallel](#4-tparallel)
5. [Testing Unexported Functions](#5-testing-unexported-functions)
6. [Interview Questions](#interview-questions)

---

## 1. Table-Driven Tests

Table-driven tests are idiomatic Go — one test function, a slice of cases, a loop. They eliminate duplicated setup code and make adding new cases trivial.

```go
// math.go
package math

func Divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, fmt.Errorf("division by zero")
    }
    return a / b, nil
}
```

```go
// math_test.go
package math

import (
    "errors"
    "testing"
)

func TestDivide(t *testing.T) {
    tests := []struct {
        name    string
        a, b    float64
        want    float64
        wantErr bool
    }{
        {name: "basic division",    a: 10, b: 2,  want: 5,   wantErr: false},
        {name: "divide by zero",   a: 10, b: 0,  want: 0,   wantErr: true},
        {name: "negative numbers", a: -6, b: 3,  want: -2,  wantErr: false},
        {name: "fractional result", a: 1, b: 3, want: 0.333, wantErr: false},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := Divide(tt.a, tt.b)

            if (err != nil) != tt.wantErr {
                t.Fatalf("Divide(%v, %v) error = %v, wantErr %v", tt.a, tt.b, err, tt.wantErr)
            }
            if !tt.wantErr && got != tt.want {
                t.Errorf("Divide(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
            }
        })
    }
}
```

```bash
go test ./...
go test -v -run TestDivide
go test -run TestDivide/basic_division  # run single sub-test
```

---

## 2. Benchmarks

Benchmark functions must start with `Benchmark` and take `*testing.B`. The `b.N` loop is run until the timing stabilizes.

```go
// sort_test.go
package sort

import (
    "math/rand"
    "sort"
    "testing"
)

func BenchmarkSort(b *testing.B) {
    data := rand.Perm(1000) // 1000 random ints

    b.ResetTimer() // exclude setup from measurement
    for i := 0; i < b.N; i++ {
        sort.Ints(data)
    }
}

// Benchmark with sub-benchmarks for different sizes
func BenchmarkSortSizes(b *testing.B) {
    sizes := []int{100, 1000, 10000}
    for _, n := range sizes {
        b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
            data := rand.Perm(n)
            b.ResetTimer()
            for i := 0; i < b.N; i++ {
                input := make([]int, len(data))
                copy(input, data)
                sort.Ints(input)
            }
        })
    }
}
```

```bash
go test -bench=.               # run all benchmarks
go test -bench=BenchmarkSort   # specific benchmark
go test -bench=. -benchmem     # include memory allocations
go test -bench=. -count=5      # run 5 times for stability

# Output:
# BenchmarkSort-8   50000   28451 ns/op   0 B/op   0 allocs/op
```

---

## 3. Mocking Without a Framework

Go's interface-based design makes mocking simple — define an interface, write a mock struct that implements it.

```go
// user_service.go
package service

type UserRepository interface {
    GetByID(id string) (*User, error)
    Save(u *User) error
}

type UserService struct {
    repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
    return &UserService{repo: repo}
}

func (s *UserService) GetUser(id string) (*User, error) {
    u, err := s.repo.GetByID(id)
    if err != nil {
        return nil, fmt.Errorf("GetUser: %w", err)
    }
    return u, nil
}
```

```go
// user_service_test.go
package service

import (
    "errors"
    "testing"
)

// Mock — implements UserRepository without a real DB
type mockUserRepo struct {
    users   map[string]*User
    saveErr error
}

func (m *mockUserRepo) GetByID(id string) (*User, error) {
    if u, ok := m.users[id]; ok {
        return u, nil
    }
    return nil, errors.New("not found")
}

func (m *mockUserRepo) Save(u *User) error { return m.saveErr }

func TestGetUser(t *testing.T) {
    t.Run("found", func(t *testing.T) {
        repo := &mockUserRepo{
            users: map[string]*User{"1": {ID: "1", Name: "Alice"}},
        }
        svc := NewUserService(repo)
        u, err := svc.GetUser("1")
        if err != nil {
            t.Fatal(err)
        }
        if u.Name != "Alice" {
            t.Errorf("got %s, want Alice", u.Name)
        }
    })

    t.Run("not found", func(t *testing.T) {
        repo := &mockUserRepo{users: map[string]*User{}}
        svc := NewUserService(repo)
        _, err := svc.GetUser("999")
        if err == nil {
            t.Fatal("expected error, got nil")
        }
    })
}
```

---

## 4. t.Parallel

`t.Parallel()` allows a test to run concurrently with other parallel tests, speeding up the suite.

```go
func TestFeatureA(t *testing.T) {
    t.Parallel() // this test can run concurrently with other Parallel tests

    // test body
}

func TestFeatureB(t *testing.T) {
    t.Parallel()
    // test body
}
```

### Parallel sub-tests — the closure bug

```go
func TestParallelSubtests(t *testing.T) {
    tests := []struct{ name, input string }{ /* ... */ }

    for _, tt := range tests {
        tt := tt // ✅ shadow loop variable BEFORE launching subtest (required pre-Go 1.22)
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel()
            // use tt safely
        })
    }
}
```

### When NOT to use t.Parallel

- Tests that share global state (package-level vars, databases, files)
- Tests that start servers on fixed ports
- Tests using `t.Setenv` (the testing package resets env vars — not safe in parallel)

---

## 5. Testing Unexported Functions

Two idiomatic approaches:

### Approach 1 — Test in the same package (white-box testing)

```go
// File: parser.go
package parser

func parseToken(s string) (string, error) { /* unexported */ }

// File: parser_test.go — note: same package name
package parser  // NOT "parser_test"

func TestParseToken(t *testing.T) {
    got, err := parseToken("Bearer abc123") // can access unexported function
    // ...
}
```

### Approach 2 — Export for test only via export_test.go

```go
// export_test.go — only compiled during testing (file ends in _test.go)
package parser

// Re-export the unexported function for the test package
var ParseToken = parseToken

// parser_test.go — uses the black-box test package
package parser_test

import "yourmodule/parser"

func TestParseToken(t *testing.T) {
    got, err := parser.ParseToken("Bearer abc123")
    // ...
}
```

---

## Interview Questions

<details>
<summary><b>What is a table-driven test and why is it idiomatic in Go?</b></summary>

A table-driven test defines test cases as a slice of structs and iterates over them in a single test function, using `t.Run` for sub-tests. It's idiomatic because: it eliminates duplicated setup/teardown code, adding new cases requires only a new struct literal, each case has a name that appears in test output, and subtests can be run individually with `-run`. The Go standard library itself uses this pattern extensively.

</details>

<details>
<summary><b>How do you benchmark a function using the testing package?</b></summary>

Write a function named `BenchmarkXxx(b *testing.B)`. Inside, call `b.ResetTimer()` after any setup, then loop `b.N` times running the code under test. Run with `go test -bench=. -benchmem` to see ns/op and allocations. Use `b.Run` for sub-benchmarks across different input sizes.

</details>

<details>
<summary><b>How do you mock dependencies in Go without a framework?</b></summary>

Define the dependency as an interface. In production code, inject the real implementation. In tests, create a struct that implements the same interface with test-controlled behavior. No framework needed — Go's implicit interfaces mean any struct with the right methods satisfies the interface. For complex scenarios, `gomock` or `testify/mock` can generate mocks from interfaces automatically.

</details>

<details>
<summary><b>What is t.Parallel() and when should you use it?</b></summary>

`t.Parallel()` marks a test to run concurrently with other parallel tests, reducing total test suite time. Use it for tests with no shared mutable state. In table-driven subtests, shadow the loop variable before calling `t.Parallel()` to avoid the closure bug (required pre-Go 1.22). Avoid with shared global state, fixed ports, or `t.Setenv`.

</details>

<details>
<summary><b>How do you test unexported functions?</b></summary>

Two ways: (1) **White-box** — put `_test.go` files in the same package (`package foo`, not `package foo_test`); unexported symbols are directly accessible. (2) **export_test.go** — create a file with `_test.go` suffix in the same package that re-exports unexported symbols as exported variables/functions; only compiled during testing, invisible to production code.

</details>
