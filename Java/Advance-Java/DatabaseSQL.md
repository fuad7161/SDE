# Database & SQL — In-Depth Notes

---

## Table of Contents

1. [Indexes (B-Tree, Composite, Covering Index)](#1-indexes-b-tree-composite-covering-index)
2. [Transaction Isolation Levels](#2-transaction-isolation-levels)
3. [EXPLAIN / EXPLAIN ANALYZE](#3-explain--explain-analyze)
4. [Connection Pooling (HikariCP)](#4-connection-pooling-hikaricp)
5. [Deadlocks in DB](#5-deadlocks-in-db)

---

## 1. Indexes (B-Tree, Composite, Covering Index)

### What Is an Index?

A data structure that speeds up row lookups at the cost of extra storage and slower writes (INSERT/UPDATE/DELETE must maintain the index).

```
Without index:                 With index on email:
Full table scan               B-Tree lookup
  row 1 → check email          root → branch → leaf → row pointer
  row 2 → check email          O(log n) instead of O(n)
  ...
  row 1,000,000 → found!
```

---

### B-Tree Index (Default)

Balanced tree structure — suited for **equality**, **range**, **ORDER BY**, and **inequality** queries.

```sql
-- Create index
CREATE INDEX idx_users_email ON users(email);

-- Used by:
SELECT * FROM users WHERE email = 'alice@example.com';          -- equality ✅
SELECT * FROM users WHERE email > 'b';                          -- range ✅
SELECT * FROM users WHERE email LIKE 'alice%';                  -- prefix ✅
SELECT * FROM users WHERE email LIKE '%alice';                  -- suffix ❌ (full scan)
SELECT * FROM users ORDER BY email;                             -- sort ✅
```

**B-Tree structure:**
```
                    [M]
                 /       \
           [D, H]         [R, V]
          /  |  \         /  |  \
        [A] [E] [I,L]  [N,P] [S] [W,Z]
         ↑   ↑    ↑      ↑    ↑    ↑
       heap heap heap   heap heap heap  ← row pointers
```

---

### Composite (Multi-Column) Index

An index on **two or more columns** — key rule: **leftmost prefix** must be used.

```sql
CREATE INDEX idx_orders_user_date ON orders(user_id, created_at);

-- ✅ Uses index (user_id is leftmost)
SELECT * FROM orders WHERE user_id = 5;
SELECT * FROM orders WHERE user_id = 5 AND created_at > '2024-01-01';
SELECT * FROM orders WHERE user_id = 5 ORDER BY created_at;

-- ❌ Does NOT use index (skips leftmost column)
SELECT * FROM orders WHERE created_at > '2024-01-01';
```

**Column order matters** — put high-cardinality (many distinct values) and most-filtered columns first:

```sql
-- Good: filter by user_id first (high cardinality), then status
CREATE INDEX idx ON orders(user_id, status, created_at);

-- Query pattern should match index prefix
WHERE user_id = ?                          -- uses index
WHERE user_id = ? AND status = ?           -- uses index
WHERE user_id = ? AND status = ? AND ...   -- uses index
WHERE status = ?                           -- does NOT use index (skips user_id)
```

---

### Covering Index

An index that contains **all columns needed by a query** — the DB can answer the query entirely from the index without accessing the table (heap) at all.

```sql
-- Query needs: user_id (filter), total_amount (select), status (select)
CREATE INDEX idx_covering ON orders(user_id, total_amount, status);

-- This query is answered entirely from the index — no heap access
SELECT total_amount, status FROM orders WHERE user_id = 5;
-- PostgreSQL EXPLAIN will show "Index Only Scan" ← the fastest possible
```

---

### Other Index Types

| Type | Use case | DB Support |
|---|---|---|
| `HASH` | Equality only (`=`), faster than B-Tree | PostgreSQL, MySQL MEMORY |
| `GIN` | Full-text search, arrays, JSONB | PostgreSQL |
| `GiST` | Geometric data, range types | PostgreSQL |
| `BRIN` | Very large tables with naturally ordered data (timestamps) | PostgreSQL |
| `FULLTEXT` | Text search with relevance scoring | MySQL |

---

### Index Best Practices

```sql
-- ✅ Index foreign keys (always)
CREATE INDEX idx_orders_customer_id ON orders(customer_id);

-- ✅ Index columns used in WHERE, JOIN, ORDER BY frequently
CREATE INDEX idx_products_category_price ON products(category_id, price);

-- ✅ Partial index — index only a subset of rows
CREATE INDEX idx_active_users ON users(email) WHERE active = true;
-- Smaller index, only used for active user queries

-- ✅ Unique index — enforces uniqueness AND speeds up lookups
CREATE UNIQUE INDEX idx_users_email ON users(email);

-- ❌ Don't index low-cardinality columns alone (boolean, status with 2-3 values)
-- e.g., CREATE INDEX ON orders(is_deleted) → NOT useful (50% of rows → full scan anyway)
```

---

## 2. Transaction Isolation Levels

### Concurrency Problems

| Problem | Description |
|---|---|
| **Dirty Read** | Thread A reads data that Thread B wrote but **not yet committed** (could be rolled back) |
| **Non-Repeatable Read** | Thread A reads a row, Thread B **updates and commits** it, Thread A reads it again → different value |
| **Phantom Read** | Thread A queries a set of rows, Thread B **inserts/deletes** rows matching the query, Thread A queries again → different result set |
| **Lost Update** | Two transactions read and then update the same row — one update overwrites the other |

---

### Isolation Levels

#### READ UNCOMMITTED
```sql
SET TRANSACTION ISOLATION LEVEL READ UNCOMMITTED;

-- Thread A (not yet committed):
UPDATE accounts SET balance = 0 WHERE id = 1;

-- Thread B can read this dirty data:
SELECT balance FROM accounts WHERE id = 1;  -- returns 0 (uncommitted!)

-- Thread A rolls back → Thread B used invalid data
```

#### READ COMMITTED (PostgreSQL default)
```sql
-- Each statement sees only committed data at that moment
-- Dirty reads impossible, but non-repeatable reads possible

-- Thread A:
SELECT balance FROM accounts WHERE id = 1;  -- returns 1000

-- Thread B commits: UPDATE accounts SET balance = 500 WHERE id = 1;

-- Thread A (same transaction):
SELECT balance FROM accounts WHERE id = 1;  -- returns 500 (changed!)
```

#### REPEATABLE READ (MySQL InnoDB default)
```sql
-- All reads in a transaction see the snapshot from the first read
-- Non-repeatable reads impossible; phantom reads possible in some DBs

-- Thread A:
SELECT balance FROM accounts WHERE id = 1;  -- returns 1000 (snapshot taken)

-- Thread B commits: UPDATE accounts SET balance = 500 WHERE id = 1;

-- Thread A (same transaction):
SELECT balance FROM accounts WHERE id = 1;  -- still returns 1000 (snapshot)
```

#### SERIALIZABLE
```sql
-- Transactions execute as if serial (one after another)
-- All anomalies prevented; lowest concurrency, highest overhead

SET TRANSACTION ISOLATION LEVEL SERIALIZABLE;

-- Thread A reads range:
SELECT SUM(balance) FROM accounts WHERE type = 'savings';

-- Thread B inserts new savings account → blocked or causes serialization failure
```

### Summary Table

| Isolation Level | Dirty Read | Non-Repeatable Read | Phantom Read |
|---|---|---|---|
| READ UNCOMMITTED | ✅ Possible | ✅ Possible | ✅ Possible |
| READ COMMITTED | ❌ Prevented | ✅ Possible | ✅ Possible |
| REPEATABLE READ | ❌ Prevented | ❌ Prevented | ✅ Possible* |
| SERIALIZABLE | ❌ Prevented | ❌ Prevented | ❌ Prevented |

*PostgreSQL REPEATABLE READ also prevents phantom reads.

### In Spring

```java
@Transactional(isolation = Isolation.READ_COMMITTED)
public BigDecimal getBalance(Long accountId) { ... }

@Transactional(isolation = Isolation.SERIALIZABLE)
public void transfer(Long from, Long to, BigDecimal amount) { ... }
```

---

## 3. EXPLAIN / EXPLAIN ANALYZE

### `EXPLAIN` — Query Plan (No Execution)

Shows the **plan** the query planner chose — cost estimates, scan types, join algorithms.

```sql
EXPLAIN SELECT * FROM orders WHERE user_id = 5 AND status = 'PENDING';

-- Output:
Index Scan using idx_orders_user_status on orders  (cost=0.43..8.45 rows=3 width=120)
  Index Cond: ((user_id = 5) AND (status = 'PENDING'))
```

### `EXPLAIN ANALYZE` — Execute + Measure (PostgreSQL)

Actually **runs** the query and shows real timing alongside the plan:

```sql
EXPLAIN ANALYZE
SELECT o.id, c.name, SUM(oi.price) AS total
FROM orders o
JOIN customers c ON c.id = o.customer_id
JOIN order_items oi ON oi.order_id = o.id
WHERE o.created_at > NOW() - INTERVAL '30 days'
GROUP BY o.id, c.name;

-- Output (PostgreSQL):
HashAggregate  (cost=1245.30..1267.80 rows=1500 width=48)
               (actual time=234.1..235.2 rows=1482 loops=1)
  Group Key: o.id, c.name
  ->  Hash Join  (cost=320.00..1182.80 rows=6250 width=40)
                 (actual time=12.3..198.4 rows=7410 loops=1)
        Hash Cond: (oi.order_id = o.id)
        ->  Seq Scan on order_items oi  (cost=0.00..680.00 rows=25000 width=16)
                                        (actual time=0.1..45.2 rows=25000 loops=1)
        ->  Hash  (cost=287.50..287.50 rows=1500 width=32)
                  (actual time=11.8..11.8 rows=1482 loops=1)
              ->  Hash Join  (cost=...)
                    ->  Index Scan on orders o (actual time=0.02..8.4 rows=1482 loops=1)
                          Index Cond: (created_at > ...)
Planning time: 1.8 ms
Execution time: 235.9 ms
```

### Key Terms to Look For

| Term | Meaning | Action |
|---|---|---|
| `Seq Scan` | Full table scan | Consider adding an index |
| `Index Scan` | Using index, fetching heap rows | Good — index is being used |
| `Index Only Scan` | Covered by index — no heap access | Best case |
| `Bitmap Index Scan` | Multiple index ranges combined | Good for multiple conditions |
| `Hash Join` | Hashing smaller table for join | Good for large tables |
| `Nested Loop` | O(n×m) join — fine for small tables | Bad for large × large |
| `rows=` estimate vs actual | Big difference = stale statistics | Run `ANALYZE tablename` |
| High `cost=` | Expensive operation | Investigate with EXPLAIN |

```sql
-- Update statistics so planner makes better decisions
ANALYZE orders;

-- Verbose output with buffers (cache hits vs disk reads)
EXPLAIN (ANALYZE, BUFFERS, FORMAT TEXT)
SELECT * FROM orders WHERE user_id = 5;
```

---

## 4. Connection Pooling (HikariCP)

### Why Connection Pooling?

Opening a DB connection is expensive (~50–100ms): TCP handshake, authentication, session setup.  
A **pool** maintains a set of pre-opened connections and reuses them.

```
Application thread → borrow connection from pool
                          ↓
                   execute SQL
                          ↓
                   return connection to pool (not closed)

Next thread → immediately borrows the already-open connection
```

### HikariCP (Default in Spring Boot)

The fastest JDBC connection pool — minimal overhead, smart connection health checking.

```yaml
# application.yml
spring:
  datasource:
    url: jdbc:postgresql://localhost:5432/mydb
    username: user
    password: secret
    hikari:
      pool-name: MyPool
      minimum-idle: 5              # min idle connections kept alive
      maximum-pool-size: 20        # max total connections in pool
      idle-timeout: 600000         # remove idle conn after 10 min (ms)
      connection-timeout: 30000    # throw exception if no conn available in 30s
      max-lifetime: 1800000        # retire connection after 30 min (prevents stale conns)
      keepalive-time: 60000        # send keepalive query every 60s
      connection-test-query: SELECT 1  # validate connection before borrow (if needed)
```

### Pool Sizing

```
Rule of thumb (OLTP web apps):
  pool_size = (core_count * 2) + effective_spindle_count

For PostgreSQL specifically (from PgBouncer docs):
  pool_size = Connections that can be ACTIVELY used simultaneously
            = number of threads doing DB work at once
            ≠ total application threads

  Too large → context switching overhead, DB connection overhead
  Too small → threads queue waiting for connections
```

```java
// Programmatic HikariCP setup
HikariConfig config = new HikariConfig();
config.setJdbcUrl("jdbc:postgresql://localhost/mydb");
config.setUsername("user");
config.setPassword("secret");
config.setMaximumPoolSize(20);
config.setMinimumIdle(5);
config.setConnectionTimeout(30_000);
config.setMaxLifetime(1_800_000);

HikariDataSource dataSource = new HikariDataSource(config);
```

### Monitoring

```java
// HikariCP exposes metrics via MeterRegistry (Micrometer)
// With Spring Boot Actuator:
management.endpoints.web.exposure.include=health,metrics
# GET /actuator/metrics/hikaricp.connections.active
# GET /actuator/metrics/hikaricp.connections.pending
# GET /actuator/metrics/hikaricp.connections.timeout
```

### Common Issues

| Problem | Symptom | Fix |
|---|---|---|
| Pool exhaustion | `Connection is not available, request timed out` | Increase pool size or fix slow queries |
| Connection leak | Pool slowly drains | Ensure connections returned (use try-with-resources) |
| Stale connections | `Connection reset` errors | Set `max-lifetime` < DB timeout; enable keepalive |
| Too large pool | Slow queries, DB CPU spike | Reduce pool size; DB can't handle too many concurrent connections |

---

## 5. Deadlocks in DB

### What Is a DB Deadlock?

Two or more transactions each hold a lock the other needs — circular wait — neither can proceed.

```
Transaction 1:                    Transaction 2:
LOCK row A (id=1)                 LOCK row B (id=2)
  ... waiting for row B ...         ... waiting for row A ...
  ← DEADLOCK DETECTED →
```

```sql
-- Transaction 1
BEGIN;
UPDATE accounts SET balance = balance - 100 WHERE id = 1;  -- locks row 1
-- (paused)
UPDATE accounts SET balance = balance + 100 WHERE id = 2;  -- waits for row 2

-- Transaction 2 (concurrent)
BEGIN;
UPDATE accounts SET balance = balance - 50 WHERE id = 2;   -- locks row 2
UPDATE accounts SET balance = balance + 50 WHERE id = 1;   -- waits for row 1 → DEADLOCK
```

The DB detects the cycle and **kills one transaction** (the "deadlock victim"), rolling it back and returning an error to the application.

---

### Prevention Strategies

#### 1. Consistent Lock Ordering

Always acquire locks in the **same order** across all transactions:

```java
// BAD — T1 locks account 1 then 2; T2 locks 2 then 1
void transfer(Long fromId, Long toId) {
    Account from = repository.findById(fromId);
    Account to   = repository.findById(toId);
    ...
}

// GOOD — always lock lower ID first
void transfer(Long fromId, Long toId) {
    Long firstId  = Math.min(fromId, toId);
    Long secondId = Math.max(fromId, toId);
    Account first  = repository.findByIdForUpdate(firstId);
    Account second = repository.findByIdForUpdate(secondId);
    ...
}
```

#### 2. Short Transactions

Keep transactions as **short as possible** — hold locks for the minimum time:

```java
// BAD — long transaction with external call holding DB lock
@Transactional
public void processOrder(Long orderId) {
    Order order = orderRepo.findById(orderId);   // lock held
    callExternalPaymentApi();                    // 2-3 seconds — lock still held!
    orderRepo.save(order);
}

// GOOD — external call outside transaction
public void processOrder(Long orderId) {
    String paymentId = callExternalPaymentApi();  // outside transaction
    saveOrderWithPayment(orderId, paymentId);      // short transaction
}

@Transactional
void saveOrderWithPayment(Long orderId, String paymentId) { ... }
```

#### 3. Use `SELECT FOR UPDATE SKIP LOCKED` (Queue Pattern)

```sql
-- Claim a pending job without deadlocking with other workers
SELECT * FROM jobs
WHERE status = 'PENDING'
ORDER BY created_at
LIMIT 1
FOR UPDATE SKIP LOCKED;   -- skip rows locked by other transactions
```

```java
// Spring Data
@Lock(LockModeType.PESSIMISTIC_WRITE)
@QueryHints(@QueryHint(name = "javax.persistence.lock.timeout", value = "-2"))  // SKIP_LOCKED
@Query("SELECT j FROM Job j WHERE j.status = 'PENDING' ORDER BY j.createdAt")
Optional<Job> claimNextJob();
```

#### 4. Retry on Deadlock

```java
@Retryable(
    value = {CannotAcquireLockException.class, DeadlockLoserDataAccessException.class},
    maxAttempts = 3,
    backoff = @Backoff(delay = 100, multiplier = 2)
)
@Transactional
public void transfer(Long from, Long to, BigDecimal amount) {
    // will retry up to 3 times with 100ms, 200ms backoff
}
```

### Detecting Deadlocks

```sql
-- PostgreSQL: view current locks
SELECT pid, relation::regclass, mode, granted
FROM pg_locks l
JOIN pg_stat_activity a ON a.pid = l.pid
WHERE NOT granted;

-- PostgreSQL: deadlock info in log
-- log_min_messages = ERROR will show deadlock errors

-- MySQL: last deadlock info
SHOW ENGINE INNODB STATUS;
-- Look for "LATEST DETECTED DEADLOCK" section
```

### Deadlock vs DB-Level Lock Wait Timeout

```sql
-- PostgreSQL — set lock wait timeout to avoid indefinite blocking
SET lock_timeout = '5s';

-- Spring / JPA
@QueryHints(@QueryHint(name = "javax.persistence.lock.timeout", value = "5000"))
```
