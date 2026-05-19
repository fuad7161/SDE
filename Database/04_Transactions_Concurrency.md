# 🟠 Transactions & Concurrency

> **Category:** Advanced &nbsp;|&nbsp; **Tags:** `isolation levels` `phantom read` `deadlock`

---

## Table of Contents
1. [What is a Transaction?](#what-is-a-transaction)
2. [Isolation Levels](#isolation-levels)
3. [Read Anomalies](#read-anomalies)
4. [Optimistic vs Pessimistic Locking](#optimistic-vs-pessimistic-locking)
5. [Deadlocks](#deadlocks)
6. [Interview Questions](#interview-questions)

---

## What is a Transaction?

A **transaction** is a sequence of one or more SQL operations executed as a single logical unit that satisfies ACID properties.

```sql
BEGIN;                                          -- start transaction

UPDATE accounts SET balance = balance - 500
WHERE account_id = 1;                          -- debit

UPDATE accounts SET balance = balance + 500
WHERE account_id = 2;                          -- credit

COMMIT;                                         -- make permanent
-- or ROLLBACK; to undo everything
```

### Savepoints — Partial rollback
```sql
BEGIN;
  UPDATE orders SET status = 'processing' WHERE id = 101;
  SAVEPOINT sp1;

  UPDATE inventory SET qty = qty - 1 WHERE product_id = 42;
  -- Something went wrong with inventory update
  ROLLBACK TO SAVEPOINT sp1;   -- undo inventory update, keep order update

COMMIT;
```

---

## Isolation Levels

SQL standard defines 4 isolation levels, trading **data accuracy** for **concurrency/performance**:

| Level | Dirty Read | Non-Repeatable Read | Phantom Read |
|-------|-----------|---------------------|-------------|
| **Read Uncommitted** | ✅ Possible | ✅ Possible | ✅ Possible |
| **Read Committed** | ❌ Prevented | ✅ Possible | ✅ Possible |
| **Repeatable Read** | ❌ Prevented | ❌ Prevented | ✅ Possible |
| **Serializable** | ❌ Prevented | ❌ Prevented | ❌ Prevented |

Higher isolation = more locks = less concurrency = lower throughput.

### Setting isolation level
```sql
-- PostgreSQL
BEGIN TRANSACTION ISOLATION LEVEL REPEATABLE READ;

-- MySQL
SET SESSION TRANSACTION ISOLATION LEVEL SERIALIZABLE;
```

### Defaults
- **PostgreSQL:** Read Committed
- **MySQL InnoDB:** Repeatable Read
- **SQL Server:** Read Committed
- **Oracle:** Read Committed

---

## Read Anomalies

### 1. Dirty Read
Transaction A reads **uncommitted data** written by transaction B. If B rolls back, A has read data that never existed.

```
T1: UPDATE accounts SET balance = 1000 WHERE id = 1;   -- not yet committed
T2: SELECT balance FROM accounts WHERE id = 1;  → reads 1000  (dirty!)
T1: ROLLBACK;  -- balance is back to 500, but T2 already acted on 1000
```
**Prevented by:** Read Committed and higher.

---

### 2. Non-Repeatable Read
Transaction A reads the same row **twice** and gets **different values** because transaction B updated and committed in between.

```
T1: SELECT balance FROM accounts WHERE id = 1;  → 500
T2: UPDATE accounts SET balance = 800 WHERE id = 1; COMMIT;
T1: SELECT balance FROM accounts WHERE id = 1;  → 800  (changed!)
```
**Prevented by:** Repeatable Read and higher.

---

### 3. Phantom Read
Transaction A executes the **same query twice** and gets **different sets of rows** because transaction B inserted or deleted rows in between.

```
T1: SELECT COUNT(*) FROM orders WHERE user_id = 5;  → 3
T2: INSERT INTO orders (user_id, ...) VALUES (5, ...); COMMIT;
T1: SELECT COUNT(*) FROM orders WHERE user_id = 5;  → 4  (phantom row!)
```
**Prevented by:** Serializable only (Repeatable Read prevents data change but not new row insertion).

> **Note:** In PostgreSQL, Repeatable Read also prevents phantom reads (uses MVCC snapshot).

---

## Optimistic vs Pessimistic Locking

### Pessimistic Locking
Assumes conflicts **will happen** — acquires a lock before reading, blocks other transactions.

```sql
-- SELECT FOR UPDATE — locks the row until transaction ends
BEGIN;
SELECT * FROM seats WHERE id = 42 FOR UPDATE;
-- Other transactions trying to read this row will WAIT
UPDATE seats SET status = 'booked' WHERE id = 42;
COMMIT;
```

- **Good for:** High-contention scenarios (seat booking, bank transfers).
- **Bad for:** Long transactions — locks block other users.

### Optimistic Locking
Assumes conflicts are **rare** — reads without locking, checks for conflicts at write time using a **version column**.

```sql
-- Table has a `version` column
-- Read
SELECT id, seat_no, status, version FROM seats WHERE id = 42;
-- Returns: id=42, status='available', version=5

-- Update — only succeeds if version hasn't changed
UPDATE seats
SET status = 'booked', version = version + 1
WHERE id = 42 AND version = 5;

-- If 0 rows affected → conflict detected → retry or fail
```

- **Good for:** Low-contention, read-heavy scenarios.
- **Bad for:** High-contention (many retries degrade performance).

### Comparison

| | Pessimistic | Optimistic |
|--|------------|-----------|
| Locking | Lock at read time | No lock; check at write time |
| Conflict handling | Prevented upfront | Detected at commit |
| Throughput | Lower (waiting) | Higher (no waiting) |
| Best for | High contention | Low contention |
| Risk | Deadlocks | Retry storms |

---

## Deadlocks

A **deadlock** occurs when two (or more) transactions are **each waiting for a lock held by the other**, creating a cycle with no resolution.

```
T1 holds lock on Row A, wants lock on Row B
T2 holds lock on Row B, wants lock on Row A
→ Neither can proceed → Deadlock
```

### Example
```sql
-- Transaction 1:
BEGIN;
UPDATE accounts SET balance = balance - 100 WHERE id = 1;  -- locks row 1
-- ... (T2 runs here) ...
UPDATE accounts SET balance = balance + 100 WHERE id = 2;  -- WAITS for T2's lock on row 2

-- Transaction 2:
BEGIN;
UPDATE accounts SET balance = balance - 50 WHERE id = 2;   -- locks row 2
UPDATE accounts SET balance = balance + 50 WHERE id = 1;   -- WAITS for T1's lock on row 1

-- DEADLOCK → DB kills one transaction, returns error to application
```

### Detection
Most databases detect deadlocks automatically using a **wait-for graph**. When a cycle is detected, one transaction is chosen as the **victim** and rolled back (usually the one that has done the least work).

### Prevention Strategies

| Strategy | How |
|----------|-----|
| **Lock ordering** | Always acquire locks in the **same order** (e.g., always lock lower ID first) |
| **Short transactions** | Minimize time locks are held |
| **Avoid user interaction inside transactions** | Never wait for user input with locks held |
| **Use SELECT FOR UPDATE NOWAIT** | Fail immediately if lock unavailable instead of waiting |
| **Retry logic** | Catch deadlock errors and retry the transaction |

```sql
-- Lock with timeout instead of waiting forever
SELECT * FROM seats WHERE id = 42 FOR UPDATE NOWAIT;
-- Raises error immediately if locked

SELECT * FROM seats WHERE id = 42 FOR UPDATE SKIP LOCKED;
-- Skips locked rows — useful for job queue processing
```

---

## Interview Questions

### Q1. What are database isolation levels? Which anomalies does each prevent?

> **Answer:**
> | Level | Prevents |
> |-------|---------|
> | Read Uncommitted | Nothing — allows dirty reads |
> | Read Committed | Dirty reads |
> | Repeatable Read | Dirty reads + non-repeatable reads |
> | Serializable | All anomalies including phantom reads |
>
> Higher isolation level = stronger guarantees but less concurrency. Most applications use **Read Committed** (default in PostgreSQL). Use **Repeatable Read** when you need stable snapshots (e.g., reports). Use **Serializable** for financial transactions.

---

### Q2. What is the difference between a dirty read, non-repeatable read, and phantom read?

> **Answer:**
> - **Dirty read:** Reading uncommitted changes from another transaction. If that transaction rolls back, you've read "ghost" data.
> - **Non-repeatable read:** Same row is read twice in a transaction; returns different values because another transaction modified and committed the row between reads.
> - **Phantom read:** Same query is run twice; returns different rows because another transaction inserted or deleted rows matching the query's condition.

---

### Q3. What is the difference between optimistic and pessimistic locking?

> **Answer:**
> - **Pessimistic locking:** Lock the row at read time (`SELECT FOR UPDATE`). No one else can modify it until you release the lock. Prevents conflicts but reduces concurrency.
> - **Optimistic locking:** Read without locking; include a version number. At write time, check that the version hasn't changed. If it has, someone else modified the record — retry.
>
> Use pessimistic when conflicts are frequent (bank balance updates). Use optimistic when reads far outweigh writes and conflicts are rare.

---

### Q4. What is a deadlock? How do you prevent it?

> **Answer:**
> A deadlock is when two transactions each hold a lock the other needs, causing both to wait forever. Databases detect this via a wait-for graph and roll back one transaction (the victim).
>
> **Prevention:**
> 1. **Always acquire locks in the same order** across all transactions (e.g., always lock lower ID first).
> 2. Keep transactions **short** — minimize lock hold time.
> 3. Use `SELECT FOR UPDATE NOWAIT` or `SKIP LOCKED` to fail fast instead of waiting.
> 4. Implement **retry logic** in application code for deadlock errors.

---

### Q5. What is MVCC (Multi-Version Concurrency Control)?

> **Answer:**
> MVCC allows **readers and writers to not block each other** by maintaining multiple versions of rows.
>
> When a row is updated:
> - The old version is kept with a timestamp/transaction ID.
> - Readers see a **snapshot** of the data as of their transaction start — they never block on writers.
> - Writers create new versions — they never block readers.
>
> PostgreSQL uses MVCC extensively, which is why it avoids many locking issues even at higher isolation levels. Old row versions are cleaned up by the `VACUUM` process.

---

### Q6. What does `SELECT FOR UPDATE SKIP LOCKED` do? Give a use case.

> **Answer:**
> `SKIP LOCKED` means: "find rows matching the WHERE clause that are **not already locked**, lock and return them — skip any that are currently locked by another transaction."
>
> **Use case — Job queue:**
> ```sql
> BEGIN;
> SELECT id, payload FROM jobs
> WHERE status = 'pending'
> ORDER BY created_at
> LIMIT 1
> FOR UPDATE SKIP LOCKED;
> -- Worker processes the job
> UPDATE jobs SET status = 'done' WHERE id = :id;
> COMMIT;
> ```
> Multiple workers can run this concurrently — each picks a different job with no blocking or deadlocking.

---

<div align="center">
  <sub>← Back to <a href="Topic.md">All Topics</a></sub>
</div>
