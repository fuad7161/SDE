# Go Backend Engineer (Mid/Senior) — Interview Prep

A curated set of the most commonly asked questions across Go internals, system design, databases, API practices, and behavioral rounds, each with a probable/model answer.

---

## 1. Go Language Fundamentals

### 1.1 Goroutines vs OS Threads
**Q: How do goroutines differ from OS threads, and how does Go's scheduler work?**

**A:** Goroutines are lightweight, user-space threads managed by the Go runtime, not the OS. They start with a small stack (~2KB) that grows/shrinks dynamically, versus OS threads which typically reserve 1-8MB fixed stack. Go uses an **M:N scheduler**: M goroutines are multiplexed onto N OS threads. Key components:
- **G** (goroutine), **M** (machine/OS thread), **P** (processor — holds run queues)
- `GOMAXPROCS` controls the number of Ps (effectively how many goroutines can run truly in parallel)
- The scheduler is cooperative-ish: goroutines yield at function calls, channel ops, syscalls, and GC safepoints
- Work-stealing: idle Ps steal goroutines from other Ps' local queues to balance load

This is why you can spawn 100,000 goroutines cheaply but not 100,000 OS threads.

---

### 1.2 Channels
**Q: Explain buffered vs unbuffered channels, and what happens with a nil or closed channel.**

**A:**
- **Unbuffered channel**: send blocks until a receiver is ready (synchronous handoff).
- **Buffered channel**: send only blocks when the buffer is full; receive blocks when empty.
- **Sending on a closed channel** → panic.
- **Receiving from a closed channel** → returns the zero value immediately, `ok` is `false` in `v, ok := <-ch`.
- **nil channel**: send/receive on it blocks forever (useful in `select` to disable a case dynamically).
- **Closing an already-closed channel** → panic.

**Follow-up trap:** "What happens if two goroutines both try to close the same channel?" → panic; only the sender should close, and only once — often enforced via `sync.Once` or a dedicated "closer" goroutine.

---

### 1.3 sync Package
**Q: When would you use a channel vs a Mutex for coordination?**

**A:** Rule of thumb — *"Do not communicate by sharing memory; share memory by communicating."*
- Use **channels** when passing ownership of data or signaling events/completion (pipelines, fan-out/fan-in, cancellation).
- Use **Mutex/RWMutex** when protecting shared in-memory state accessed by multiple goroutines (e.g., a cache map) where the overhead of channel-based serialization isn't worth it.
- `RWMutex` is preferable when reads vastly outnumber writes.
- `sync.Once` guarantees an initializer runs exactly once (e.g., singleton config load).
- `sync.WaitGroup` waits for a group of goroutines to finish — common pitfall: calling `Add` inside the goroutine instead of before `go func()`, causing a race.

---

### 1.4 Slices vs Arrays
**Q: Why did `append` not mutate the caller's slice the way I expected?**

**A:** A slice is a struct of `{pointer, len, cap}`. `append` only mutates the underlying array in place when there's enough capacity; otherwise it allocates a new array and the original slice variable (held by the caller) still points to the old array. This causes the classic bug where a function appends to a slice parameter but the caller doesn't see the new element — because len/cap changed only in the callee's copy of the slice header. Also: slicing (`s[1:3]`) shares the underlying array, so mutating an element through one slice can affect another aliasing slice — a common source of subtle bugs.

---

### 1.5 Interfaces
**Q: What's the difference between a nil interface and an interface holding a nil pointer?**

**A:** An interface value is internally `{type, value}`. `var i interface{} = (*MyStruct)(nil)` has a non-nil type (`*MyStruct`) but nil value — so `i == nil` is **false**, even though the underlying pointer is nil. This trips people up when returning typed nil errors: `func f() error { var e *MyError; return e }` returns a non-nil `error` interface even though `e` is nil. Fix: return `nil` explicitly, not a nil-valued typed variable.

---

### 1.6 Defer, Panic, Recover
**Q: In what order do deferred functions run, and how does recover work?**

**A:** Deferred calls execute in **LIFO** order, after the surrounding function's return value is set but before it actually returns to the caller — which is why a `defer func(){ ... }()` can modify a named return value. `recover()` only has effect when called directly inside a deferred function; it stops the panic's propagation and lets the function return normally. Common pattern:
```go
func safeDivide(a, b int) (result int, err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("recovered: %v", r)
        }
    }()
    return a / b, nil
}
```

---

### 1.7 Memory Management
**Q: How does escape analysis decide stack vs heap allocation?**

**A:** The Go compiler statically analyzes whether a variable's lifetime can be proven to stay within a function's scope. If yes, it's stack-allocated (cheap, no GC pressure). If the compiler can't prove that — e.g., a pointer to it is returned, stored in a global, sent to a channel, or passed to an interface — it "escapes" to the heap. You can check with `go build -gcflags="-m"`. Reducing heap escapes (e.g., avoiding unnecessary pointer returns for small structs) is a common performance optimization.

Go's GC is a **concurrent, tri-color mark-and-sweep** collector, designed for low latency (sub-millisecond pause targets) rather than max throughput.

---

### 1.8 Error Handling
**Q: How do errors.Is, errors.As, and error wrapping work?**

**A:** Since Go 1.13, `fmt.Errorf("...: %w", err)` wraps an error while preserving the chain. `errors.Is(err, target)` checks if `target` appears anywhere in the chain (useful for sentinel errors like `sql.ErrNoRows`). `errors.As(err, &target)` checks if any error in the chain matches a specific type and assigns it. This is preferred over manually comparing error strings, which is fragile. Custom error types should implement `Unwrap() error` to participate in the chain.

---

### 1.9 Context Package
**Q: What's context used for, and why is passing values via context discouraged?**

**A:** `context.Context` propagates **cancellation signals, deadlines, and request-scoped values** across API boundaries and goroutines. `context.WithCancel/WithTimeout/WithDeadline` create derived contexts; canceling a parent cancels all children. `ctx.Done()` returns a channel closed on cancellation — used in `select` to abort long-running work. Passing arbitrary business data via `context.WithValue` is discouraged because it's untyped (requires runtime type assertions), bypasses the compiler, and makes dependencies implicit/hidden — better to pass real parameters explicitly. Acceptable uses: request IDs, trace IDs, auth tokens tied to request lifecycle.

---

### 1.10 Concurrency Bugs
**Q: How do you detect and prevent goroutine leaks and race conditions?**

**A:**
- **Race conditions**: run tests/binaries with `go run -race` or `go test -race`; it instruments memory accesses and reports concurrent read/write conflicts.
- **Goroutine leaks**: happen when a goroutine blocks forever (e.g., writing to a channel nobody reads from, or waiting on a context that's never canceled). Prevent by always giving goroutines a cancellation path (`ctx.Done()`), using buffered channels where sender shouldn't block indefinitely, and using tools like `pprof` (`goroutine` profile) or `go.uber.org/goleak` in tests to catch leaks.

---

## 2. System Design / Architecture

### 2.1 Design a Rate Limiter
**Q: Design a rate limiter for an API gateway.**

**A:** Common algorithms:
- **Token bucket**: bucket refills at fixed rate, request consumes a token; allows bursts up to bucket size. Most widely used (also used by AWS, Stripe).
- **Sliding window log/counter**: more accurate but more memory (store timestamps or bucketed counters per window).
- **Fixed window counter**: simplest, but allows 2x burst at window boundaries.

Implementation for distributed systems: store counters in **Redis** with `INCR` + `EXPIRE`, or use a Lua script for atomicity (check-and-increment in one round trip). Key by user ID/API key + endpoint. For a single-instance Go service, `golang.org/x/time/rate` gives a token-bucket limiter out of the box.

---

### 2.2 Design a URL Shortener
**Q: How would you design a URL shortener like bit.ly?**

**A:** Core flow: `POST /shorten {long_url} → short_code`, `GET /{short_code} → 301/302 redirect`.
- **ID generation**: base62-encode an auto-incrementing ID (from a counter service or DB sequence), or use a hash of the URL + collision check. Avoid pure random with retry-on-collision at scale — prefer a dedicated ID generator (like Snowflake) for uniqueness without contention.
- **Storage**: key-value store (short_code → long_url) — Redis for hot lookups backed by a durable DB (Postgres/DynamoDB) as source of truth.
- **Caching**: heavy read-to-write ratio, so cache aggressively; use CDN/edge caching for redirects.
- **Scaling reads**: read replicas + cache; writes are much less frequent so a single primary is usually fine.
- **Analytics** (optional): async write click events to a queue (Kafka) rather than blocking the redirect path.

---

### 2.3 Handling High Throughput ("design for X req/sec")
**Q: How would you scale a service to handle 50k requests/sec?**

**A:** Structured approach:
1. **Identify the bottleneck first** — CPU-bound, I/O-bound, or DB-bound.
2. **Horizontal scaling** behind a load balancer (round robin / least-connections), stateless service instances.
3. **Caching** at multiple layers (CDN, in-memory/Redis, DB query cache) to cut DB load.
4. **Async processing** — push non-critical-path work (emails, notifications, logging) to a message queue instead of doing it synchronously.
5. **Database scaling** — read replicas, sharding by key, connection pooling to avoid exhausting DB connections.
6. **Circuit breakers/backpressure** — protect downstream services from cascading failure using something like `sony/gobreaker` in Go.
7. **Profiling** — use `pprof` to find actual hotspots rather than guessing.

---

### 2.4 Caching Strategies
**Q: Explain cache-aside vs write-through, and how you'd prevent a thundering herd.**

**A:**
- **Cache-aside (lazy loading)**: app checks cache first; on miss, reads DB and populates cache. Simple, but cache can go stale and every miss hits the DB.
- **Write-through**: writes go to cache and DB synchronously, keeping them consistent, at the cost of write latency.
- **Write-behind**: writes go to cache immediately, DB is updated asynchronously — faster writes, risk of data loss on cache failure.

**Thundering herd** (many requests missing cache simultaneously, e.g., on expiry of a hot key): mitigate with
- **Request coalescing/single-flight** (Go's `golang.org/x/sync/singleflight` — only one goroutine fetches, others wait on the result)
- **Jittered TTLs** so keys don't all expire at once
- **Stale-while-revalidate** — serve slightly stale data while refreshing in background

---

### 2.5 Message Queues
**Q: Kafka vs RabbitMQ vs Redis pub/sub — when would you use each?**

**A:**
- **Kafka**: high-throughput, durable, replayable log; consumers track their own offset. Best for event streaming, analytics pipelines, audit logs, or when multiple consumer groups need independent replay.
- **RabbitMQ**: traditional message broker with flexible routing (exchanges, topic/fanout routing), per-message ack/retry/dead-lettering. Best for task queues and complex routing needs, lower throughput than Kafka.
- **Redis pub/sub**: fire-and-forget, no persistence — messages are lost if no subscriber is listening at publish time. Good for ephemeral notifications (e.g., WebSocket fanout) but not for anything requiring durability.

---

### 2.6 CAP Theorem & Idempotency
**Q: Explain CAP theorem and why idempotency matters in distributed systems.**

**A:** CAP theorem says a distributed system can only guarantee two of **Consistency, Availability, Partition tolerance** at once during a network partition — since partitions are a fact of life in distributed systems, the real tradeoff is **C vs A** during a partition (e.g., DynamoDB favors AP, traditional RDBMS replication setups often favor CP).

**Idempotency**: an operation that produces the same result no matter how many times it's applied (e.g., "set balance to $100" vs "add $10", the latter isn't idempotent). Critical for distributed systems because retries are inevitable (timeouts, network blips) — a client retrying a non-idempotent "charge card" request could double-charge. Common solution: **idempotency keys** — client generates a unique key per logical operation; server stores the key + result and returns the cached result on retry instead of re-executing.

---

## 3. Databases

### 3.1 Indexing
**Q: How do B-tree indexes work, and when can an index hurt performance?**

**A:** A B-tree index maintains sorted keys in a balanced tree, giving O(log n) lookups, range scans, and ordered iteration — this is why they suit `WHERE`, `ORDER BY`, and `JOIN` columns. Indexes hurt when:
- **Write-heavy tables**: every insert/update/delete must also update all indexes, adding overhead.
- **Low-cardinality columns** (e.g., boolean flag): the optimizer may ignore the index anyway since a full scan is cheaper.
- **Over-indexing**: many unused indexes waste storage and slow down writes for no read benefit.

---

### 3.2 Transactions & Isolation Levels
**Q: Explain the SQL isolation levels and what anomalies each prevents.**

**A:**
| Level | Dirty Read | Non-repeatable Read | Phantom Read |
|---|---|---|---|
| Read Uncommitted | Possible | Possible | Possible |
| Read Committed | Prevented | Possible | Possible |
| Repeatable Read | Prevented | Prevented | Possible (MySQL InnoDB mostly prevents via MVCC) |
| Serializable | Prevented | Prevented | Prevented |

Most production systems default to **Read Committed** (Postgres default) or **Repeatable Read** (MySQL InnoDB default) as a balance between correctness and concurrency throughput. Serializable gives full correctness but can cause more lock contention/retries.

**Deadlocks**: two transactions each waiting on a lock the other holds. Mitigate by acquiring locks in a consistent order across the codebase, keeping transactions short, and letting the DB's deadlock detector kill one transaction (handle that retry in app code).

---

### 3.3 N+1 Query Problem
**Q: What is the N+1 query problem and how do you fix it?**

**A:** Happens when fetching a list of N parent records, then issuing a separate query per record to fetch related data (e.g., fetch 50 employees, then 50 separate queries for each employee's department). Fixes:
- **Eager loading / JOIN** — fetch parent + related data in one query.
- **Batch loading** — collect all IDs, issue one `WHERE id IN (...)` query instead of N queries.
- In Go specifically (no built-in ORM magic like Hibernate), this is usually solved by explicitly writing batched queries or using a library like `sqlc`/`gorm` Preload carefully.

---

### 3.4 SQL vs NoSQL
**Q: When would you choose NoSQL over a relational database?**

**A:** Choose **NoSQL** when: schema is highly variable/evolving, you need horizontal write scaling beyond what sharding a relational DB easily gives you, access patterns are simple key-based lookups (e.g., DynamoDB, Redis), or you're storing document-shaped data that doesn't fit tabular joins well (e.g., MongoDB for nested JSON-like records). Choose **SQL** when you need strong consistency, complex multi-table joins/aggregations, and mature transactional guarantees — which is most business/enterprise backend systems (like an HRMS).

---

### 3.5 Connection Pooling
**Q: Why is connection pooling important, and what happens if you misconfigure it?**

**A:** Opening a new DB connection is expensive (TCP handshake, auth). Pooling reuses a fixed set of live connections across requests. In Go, `database/sql` has built-in pooling — configure via `SetMaxOpenConns`, `SetMaxIdleConns`, `SetConnMaxLifetime`. Misconfiguration risks:
- **Too many max connections** → can exhaust the DB's own connection limit, especially with many service replicas each holding their own pool.
- **Too few** → requests queue and time out under load.
- **No `ConnMaxLifetime`** → stale connections behind a load balancer/proxy that silently drops idle connections can cause mysterious errors.

---

## 4. API / Backend Practicals

### 4.1 REST API Design
**Q: What are best practices for versioning, pagination, and status codes in REST APIs?**

**A:**
- **Versioning**: URL path (`/v1/users`) is simplest and most visible; header-based versioning is more "RESTful" but harder to debug/test manually.
- **Pagination**: offset-based (`?page=2&limit=20`) is simple but slow/inconsistent on large, mutating datasets; **cursor-based** (`?after=<opaque_cursor>`) scales better and avoids skipped/duplicated rows when data changes between pages.
- **Status codes**: `200` success, `201` created, `204` no content, `400` bad request (client error), `401` unauthenticated, `403` unauthorized, `404` not found, `409` conflict, `422` unprocessable entity (validation), `500` server error, `503` unavailable. Interviewers often probe whether you distinguish `401` vs `403` and `400` vs `422`.

---

### 4.2 Authentication/Authorization
**Q: JWT vs session-based auth — tradeoffs?**

**A:**
- **Sessions**: server stores session state (in-memory or Redis), client holds an opaque session ID cookie. Easy to revoke instantly (delete server-side record), but requires shared/centralized session storage across instances.
- **JWT**: self-contained, signed token holding claims; stateless — no server-side lookup needed, easy to scale horizontally. Downside: **hard to revoke before expiry** (mitigated with short expiry + refresh tokens, or a token blacklist/denylist which reintroduces state). Given your HRMS MCP JWT work — this is a great one to speak to directly with real specifics (payload design, expiry strategy, refresh flow).

**OAuth 2.0 flow** (authorization code grant): client redirects user to auth server → user authenticates & consents → auth server redirects back with a code → client exchanges code (+ client secret) for access/refresh tokens server-side. Interviewers often ask why the implicit flow is deprecated (token exposed in URL fragment, no client secret verification).

---

### 4.3 gRPC vs REST
**Q: When would you choose gRPC over REST?**

**A:** gRPC uses HTTP/2 + Protocol Buffers — smaller payloads, built-in streaming (unary, server-streaming, client-streaming, bidirectional), strongly typed contracts via `.proto` files enabling codegen across languages. Best for **internal service-to-service communication** where performance and strict contracts matter. REST/JSON is usually still preferred for **public-facing APIs** because of ubiquitous tooling, human readability, and browser-native support (gRPC-Web needs a proxy layer for browsers).

---

### 4.4 Middleware Patterns in Go
**Q: How do you implement middleware chaining in Go's net/http?**

**A:** Middleware wraps an `http.Handler` and returns a new `http.Handler`, typically via a function signature like `func(http.Handler) http.Handler`:
```go
func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next.ServeHTTP(w, r)
        log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
    })
}
```
Chain multiple with a helper:
```go
func Chain(h http.Handler, mws ...func(http.Handler) http.Handler) http.Handler {
    for i := len(mws) - 1; i >= 0; i-- {
        h = mws[i](h)
    }
    return h
}
```
Common middlewares: logging, auth/JWT validation, request ID injection, panic recovery, CORS, rate limiting.

---

### 4.5 Graceful Shutdown
**Q: How do you implement graceful shutdown in a Go HTTP server?**

**A:** Listen for OS signals (`SIGINT`/`SIGTERM`), then call `server.Shutdown(ctx)` which stops accepting new connections but lets in-flight requests finish (up to a timeout):
```go
srv := &http.Server{Addr: ":8080", Handler: router}

go func() {
    if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        log.Fatalf("listen: %s", err)
    }
}()

sigCh := make(chan os.Signal, 1)
signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
<-sigCh

ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()
if err := srv.Shutdown(ctx); err != nil {
    log.Fatalf("forced shutdown: %s", err)
}
```
Also close DB pools, flush logs/metrics, and stop background workers (via context cancellation) during this window.

---

## 5. Behavioral / Senior-Specific

### 5.1 Production Incident
**Q: Tell me about a time you debugged a difficult production incident.**

**A (structure to use — STAR):**
- **Situation**: brief context on the system and what broke.
- **Task**: what was expected of you (on-call, escalation, ownership).
- **Action**: how you narrowed down the cause — logs, metrics, `pprof`, reading recent deploys/diffs, isolating via feature flags.
- **Result**: the fix, the rollback/mitigation, and what monitoring/tests you added afterward to prevent recurrence.

Interviewers are listening for **methodical debugging process** over the specific bug — did you form hypotheses and test them, or just guess-and-check?

---

### 5.2 Code Review / Mentoring
**Q: How do you approach reviewing a junior engineer's code?**

**A:** Focus on: correctness and edge cases first, then maintainability (naming, structure), then style (usually enforced by linters/formatters like `gofmt`/`golangci-lint` so it's not a human debate). Frame feedback as questions/suggestions rather than commands ("What happens if this slice is empty?" vs "This is wrong"). Distinguish **blocking** issues (bugs, security, data races) from **nits** (style preference) so the junior isn't blocked on subjective feedback. Pair on the review live occasionally rather than only leaving async comments, especially for repeated mistakes — teaches the "why," not just the "what."

---

### 5.3 Technical Disagreement
**Q: Describe a time you disagreed with a technical decision.**

**A:** Best answers show: you stated your concern with concrete reasoning/data (not just opinion), you listened to the counter-argument genuinely, and — win or lose the debate — you committed to the team's decision afterward rather than relitigating it passively ("disagree and commit"). If you were wrong, say so honestly; if you were right, describe how you handled being vindicated without being smug, and what changed afterward (e.g., a new review process).

---

### 5.4 Tradeoffs Under Deadline
**Q: Tell me about a time you had to cut corners to hit a deadline.**

**A:** Good answers name the **specific technical debt incurred** (e.g., skipped writing tests for an edge case, hardcoded a config instead of building a settings UI), explain **why that tradeoff was acceptable at the time** (low blast radius, easy to revisit, business urgency was real), and — importantly — show you **tracked and later resolved it** (ticket filed, revisited post-launch) rather than letting it silently rot. This signals maturity: not "I never cut corners" (unrealistic) but "I cut corners deliberately and I don't lose track of them."

---
