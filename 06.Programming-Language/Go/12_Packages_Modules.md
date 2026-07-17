# 12 — Package & Module System

> **Topics:** go.mod · go.sum · build tags · init() order · internal packages · go workspace

---

## Table of Contents
1. [go.mod and go.sum](#1-gomod-and-gosum)
2. [Module Version Resolution — MVS](#2-module-version-resolution--mvs)
3. [init() Execution Order](#3-init-execution-order)
4. [Build Tags](#4-build-tags)
5. [internal Packages](#5-internal-packages)
6. [Interview Questions](#interview-questions)

---

## 1. go.mod and go.sum

`go.mod` declares the module path and its direct + indirect dependencies. `go.sum` is a lock file of cryptographic checksums — it guarantees reproducible builds.

```
myapp/
├── go.mod
├── go.sum
├── main.go
└── internal/
    └── db/
        └── db.go
```

```go
// go.mod
module github.com/yourorg/myapp

go 1.22

require (
    github.com/gin-gonic/gin v1.9.1
    github.com/jackc/pgx/v5 v5.5.5
)

require (
    // indirect dependencies — added automatically
    golang.org/x/net v0.24.0 // indirect
)
```

### Common commands

```bash
go mod init github.com/yourorg/myapp   # create go.mod
go mod tidy                            # add missing, remove unused deps
go get github.com/pkg/errors@v0.9.1   # add or upgrade dep
go get github.com/pkg/errors@none      # remove dep
go mod download                        # download to module cache
go mod verify                          # verify go.sum checksums
go list -m all                         # list all deps with versions
```

---

## 2. Module Version Resolution — MVS

Go uses **Minimum Version Selection (MVS)** — it always picks the *minimum* version that satisfies all requirements. No "latest" semantics, no automatic upgrades.

```
myapp requires:
    lib-a v1.3.0
    lib-b v1.0.0

lib-a v1.3.0 requires:
    lib-c v1.5.0

lib-b v1.0.0 requires:
    lib-c v1.2.0

Result: lib-c v1.5.0 is selected (minimum version that satisfies all)
```

### replace directive

```go
// go.mod
replace github.com/some/dep => ../local-fork   // local override for development
replace github.com/some/dep v1.2.0 => github.com/fork/dep v1.2.1 // use fork
```

### exclude directive

```go
// Exclude a known-broken version
exclude github.com/some/dep v1.1.3
```

---

## 3. init() Execution Order

`init()` functions run automatically before `main()`. A package can have multiple `init()` functions, even in the same file.

```go
// order.go
package main

import "fmt"

var x = initX() // 1. package-level var evaluated first

func initX() int {
    fmt.Println("var x initialised")
    return 42
}

func init() {
    fmt.Println("init() #1") // 2. then init() runs
}

func init() {
    fmt.Println("init() #2") // 3. multiple init() in order
}

func main() {
    fmt.Println("main()") // 4. finally main()
}

// Output:
// var x initialised
// init() #1
// init() #2
// main()
```

### Execution order across packages

```
main imports → pkg A and pkg B
pkg A imports → pkg C
pkg C imports → pkg D

Init order (depth-first, imports first):
    pkg D → pkg C → pkg A → pkg B → main
```

### Side-effect imports

```go
import _ "github.com/lib/pq" // registers postgres driver via init()
```

The blank import triggers the package's `init()` without exposing its exported symbols.

---

## 4. Build Tags

Build tags (constraints) control which files are compiled. Specified at the top of a file with `//go:build`.

```go
//go:build linux
// +build linux     // old syntax, kept for compatibility with Go <1.17

package main

import "fmt"

func platformInfo() string {
    return "running on Linux"
}
```

```go
//go:build windows

package main

func platformInfo() string {
    return "running on Windows"
}
```

### Boolean expressions in build tags

```go
//go:build linux && amd64        // Linux on AMD64
//go:build !windows              // any OS except Windows
//go:build linux || darwin       // Linux or macOS
//go:build integration           // custom tag
```

```bash
go build -tags integration ./...     # include files with integration tag
go test -tags integration ./...
GOOS=linux GOARCH=arm64 go build .   # cross-compile
```

### Test-only files

Any file ending in `_test.go` is automatically excluded from non-test builds — no build tag needed.

---

## 5. internal Packages

The `internal` directory restricts which packages can import a given package. Only code in the parent of `internal` can import it.

```
myapp/
├── main.go                     ✅ can import myapp/internal/db
├── internal/
│   └── db/
│       └── db.go               exported symbols restricted
├── api/
│   └── handler.go              ✅ can import myapp/internal/db
└── vendor/otherpkg/            ❌ CANNOT import myapp/internal/db
```

```go
// internal/db/db.go
package db

type Client struct { /* ... */ }   // exported but restricted by internal rule

func New() *Client { return &Client{} }
```

```go
// main.go — ok, within myapp
import "github.com/yourorg/myapp/internal/db"

// another module:
import "github.com/yourorg/myapp/internal/db" // ❌ compile error
// use of internal package not allowed
```

### Go workspace (go.work) — multi-module development

```bash
go work init ./myapp ./shared-lib   # create go.work

# go.work
go 1.22
use (
    ./myapp
    ./shared-lib
)
```

This lets `myapp` use the local `shared-lib` without a `replace` directive in `go.mod`.

---

## Interview Questions

<details>
<summary><b>What is the difference between go.mod and go.sum?</b></summary>

`go.mod` declares the module path, Go version, and direct/indirect dependencies with their required minimum versions. `go.sum` is a lock file containing cryptographic hashes (SHA-256) of every dependency's zip archive and `go.mod` file. `go.sum` ensures reproducible builds — the same source is used every time. You commit both files. `go mod tidy` keeps them consistent.

</details>

<details>
<summary><b>How does Go's Minimum Version Selection (MVS) work?</b></summary>

MVS selects the minimum version that satisfies all module requirements in the dependency graph. If module A needs lib-c v1.5 and module B needs lib-c v1.2, Go picks v1.5 (minimum version that satisfies everyone). There is no "latest" automatic upgrade — versions only change when `go.mod` explicitly requires a higher version. This makes builds reproducible and predictable compared to "compatible range" selectors in npm/cargo.

</details>

<details>
<summary><b>What is the execution order of init() functions?</b></summary>

Order: (1) imported packages' `init()` runs first, depth-first (deepest dependency initialises first). (2) Within a package, package-level variable declarations run in source order. (3) Then `init()` functions in the order they appear in source files (files processed in lexical filename order). (4) Multiple `init()` functions in the same file run top to bottom. (5) Finally `main()` is called. `init()` cannot be called explicitly.

</details>

<details>
<summary><b>What are build tags and how do you use them?</b></summary>

Build tags are constraints that tell the compiler which files to include. Specified with `//go:build` at the very top of a file (above `package`). Supports boolean expressions: `&&`, `||`, `!`. Common uses: OS-specific code (`//go:build linux`), architecture-specific code, custom tags for integration tests (`-tags integration`). Files with the wrong tag for the current build are silently excluded.

</details>

<details>
<summary><b>What is the purpose of the internal directory?</b></summary>

The `internal` package convention enforces access control at the module level. Only code rooted at the parent of the `internal` directory can import packages inside it. This prevents external modules from depending on implementation details you intend to keep private. It's the Go way of having package-private visibility across multiple packages within the same module without making them fully public.

</details>
