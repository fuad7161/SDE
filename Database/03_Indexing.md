# 🔵 Indexing

> **Category:** Performance &nbsp;|&nbsp; **Tags:** `B-Tree vs Hash` `composite index` `EXPLAIN`

---

## Table of Contents
1. [What is an Index?](#what-is-an-index)
2. [B-Tree vs Hash Indexes](#b-tree-vs-hash-indexes)
3. [Composite Indexes](#composite-indexes)
4. [Covering Indexes](#covering-indexes)
5. [When Indexes Hurt](#when-indexes-hurt)
6. [EXPLAIN / EXPLAIN ANALYZE](#explain--explain-analyze)
7. [Interview Questions](#interview-questions)

---

## What is an Index?

An **index** is a separate data structure that the database maintains to speed up data retrieval. It works like a book's index — instead of scanning every page (full table scan), you jump directly to the location.

**Trade-off:**
- Reads: faster (avoids full table scan)
- Writes (INSERT/UPDATE/DELETE): slower (index must be maintained)
- Storage: extra disk space

```sql
-- Create index
CREATE INDEX idx_employees_email ON employees(email);

-- Create unique index
CREATE UNIQUE INDEX idx_users_email ON users(email);

-- Drop index
DROP INDEX idx_employees_email;
```

---

## B-Tree vs Hash Indexes

### B-Tree (Balanced Tree) Index — Default in most databases

```
                  [50]
                /      \
          [20,35]      [65,80]
          /  |  \      /  |  \
        [10][25][40] [55][70][90]
```

- Data stored in **sorted order** in leaf nodes.
- All leaf nodes linked — efficient **range scans**.
- Supports: `=`, `<`, `>`, `<=`, `>=`, `BETWEEN`, `LIKE 'prefix%'`, `ORDER BY`.

### Hash Index

```
Hash function: key → bucket
hash(email) → bucket 42 → [alice@example.com, row_ptr]
```

- O(1) lookup for exact equality.
- **Does NOT support** range queries, sorting, or LIKE.
- PostgreSQL: hash indexes on heap tables. MySQL Memory engine.

### Comparison

| Feature | B-Tree | Hash |
|---------|--------|------|
| Equality (`=`) | ✅ | ✅ (faster) |
| Range (`<`, `>`, `BETWEEN`) | ✅ | ❌ |
| `ORDER BY` | ✅ | ❌ |
| `LIKE 'prefix%'` | ✅ | ❌ |
| Default index type | ✅ | ❌ |
| Overhead | Moderate | Low |

**Use B-Tree for almost everything. Use Hash only for in-memory equality lookups.**

---

## Composite Indexes

A **composite (multi-column) index** covers multiple columns in a defined order.

```sql
CREATE INDEX idx_orders_user_status_date
ON orders(user_id, status, created_at);
```

### The Left-Prefix Rule

The index is only used if the query filters on **columns from the left side**.

```sql
-- ✅ Uses index (user_id is leftmost)
SELECT * FROM orders WHERE user_id = 1;

-- ✅ Uses index (user_id + status — left prefix)
SELECT * FROM orders WHERE user_id = 1 AND status = 'shipped';

-- ✅ Full index used
SELECT * FROM orders WHERE user_id = 1 AND status = 'shipped' AND created_at > '2024-01-01';

-- ❌ Index NOT used (skips user_id)
SELECT * FROM orders WHERE status = 'shipped';

-- ❌ Index NOT used (skips user_id and status)
SELECT * FROM orders WHERE created_at > '2024-01-01';
```

### Column Order Strategy

1. **Equality columns first** — columns used in `=` conditions come before range columns.
2. **High selectivity first** — columns that filter out the most rows.
3. **Range column last** — index stops being useful after the first range condition.

```sql
-- Query: WHERE user_id = 1 AND status = 'active' AND created_at > '2024-01-01'
-- Best index: (user_id, status, created_at)  ← equality, equality, range
-- Bad index:  (created_at, user_id, status)  ← range first kills usefulness
```

---

## Covering Indexes

A **covering index** includes all columns needed by a query — the database can satisfy the query **entirely from the index** without touching the main table (heap).

```sql
-- Query
SELECT user_id, status, created_at FROM orders WHERE user_id = 42;

-- Covering index — all 3 columns are in the index
CREATE INDEX idx_covering ON orders(user_id, status, created_at);
-- No heap access needed → very fast (index-only scan)
```

### INCLUDE columns (PostgreSQL 11+)
```sql
-- Add non-key columns to the index leaf pages without making them sort keys
CREATE INDEX idx_orders_user ON orders(user_id)
INCLUDE (status, total_amount);
-- user_id is the B-Tree key; status & total_amount are stored in leaves for covering
```

---

## When Indexes Hurt

### 1. Write-heavy tables
Every INSERT, UPDATE, DELETE must also update all indexes on the table.

```sql
-- Bulk insert into a table with 8 indexes
INSERT INTO logs SELECT * FROM staging;  -- 8× slower than unindexed
-- Solution: DROP indexes, bulk insert, REBUILD indexes
```

### 2. Low-cardinality columns
An index on a column with few distinct values (e.g., `status` with 3 values, `gender`) offers little benefit — the DB still reads most rows.

```sql
-- Low cardinality: ~50% of rows per value
CREATE INDEX idx_users_gender ON users(gender);
-- Optimizer may choose full scan anyway — less overhead than index + heap reads
```

### 3. Small tables
Full table scan on a 100-row table is faster than index lookup + heap fetch overhead.

### 4. Redundant / duplicate indexes
```sql
-- Redundant: idx(a,b) already covers queries on (a) alone
CREATE INDEX idx_a ON t(a);        -- redundant if idx_ab exists
CREATE INDEX idx_ab ON t(a, b);
```

### 5. Unused indexes
Check `pg_stat_user_indexes` (PostgreSQL) for indexes with zero scans — they waste storage and slow writes.

---

## EXPLAIN / EXPLAIN ANALYZE

`EXPLAIN` shows the **query execution plan** — how the database intends to execute a query.

```sql
EXPLAIN SELECT * FROM employees WHERE dept_id = 3;

-- Output:
Seq Scan on employees  (cost=0.00..450.00 rows=50 width=32)
  Filter: (dept_id = 3)
```

```sql
EXPLAIN ANALYZE SELECT * FROM employees WHERE dept_id = 3;

-- Output:
Index Scan using idx_dept_id on employees
  (cost=0.29..8.31 rows=50 width=32) (actual time=0.025..0.143 rows=48 loops=1)
  Index Cond: (dept_id = 3)
Planning Time: 0.3 ms
Execution Time: 0.2 ms
```

### Key terms to understand

| Term | Meaning |
|------|---------|
| **Seq Scan** | Full table scan — reads every row |
| **Index Scan** | Uses index → heap fetch for each row |
| **Index Only Scan** | Covering index — no heap access |
| **Bitmap Heap Scan** | Collects matching pages from index, then reads heap in bulk |
| **Nested Loop** | For each row in outer, scan inner — good for small sets |
| **Hash Join** | Build hash table from inner, probe with outer — good for large equal sets |
| **Merge Join** | Both sides sorted — efficient for large sorted datasets |
| **cost=X..Y** | Startup cost .. total cost (in arbitrary units) |
| **rows=N** | Estimated rows |
| **actual time** | Real time (only in ANALYZE) |

### Red flags in EXPLAIN
- **Seq Scan on large table** — missing index
- **High row estimate vs actual** — stale statistics, run `ANALYZE`
- **Nested Loop with large outer set** — may need index on inner table
- **Sort** node — missing index for ORDER BY

---

## Interview Questions

### Q1. What is the difference between a B-Tree and a Hash index?

> **Answer:**
> - **B-Tree:** Data stored in sorted order in a balanced tree. Supports equality, range queries (`<`, `>`, `BETWEEN`), prefix LIKE, and ORDER BY. Default index type.
> - **Hash:** Maps key to a bucket via hash function. O(1) equality lookup. Cannot support range queries, sorting, or prefix matching.
>
> Use B-Tree for general-purpose indexing. Hash indexes are rarely used in practice outside in-memory tables.

---

### Q2. What is the left-prefix rule for composite indexes?

> **Answer:**
> A composite index `(A, B, C)` can only be used if the query filters on A, or A+B, or A+B+C — always starting from the leftmost column. A query filtering only on B or C cannot use this index.
>
> The index stops being effective after the first **range condition** (`<`, `>`, `BETWEEN`, `LIKE`). So put equality conditions first and the range column last.

---

### Q3. What is a covering index? Why is it faster?

> **Answer:**
> A covering index includes all columns referenced by a query (WHERE, SELECT, ORDER BY). When the database can answer the query entirely from the index without accessing the main table (heap), it performs an **index-only scan** — much faster because it avoids random I/O to fetch heap pages.

---

### Q4. Why would adding an index make performance worse?

> **Answer:**
> - **Write overhead:** Every INSERT/UPDATE/DELETE must update all indexes on the table.
> - **Low cardinality:** If the column has few distinct values, the index doesn't filter much and the optimizer may prefer a full scan.
> - **Small table:** Index lookup + heap fetch has overhead; full scan may be faster.
> - **Redundant indexes** waste write performance and storage.
> - The optimizer may choose a bad index due to stale statistics — run `ANALYZE`.

---

### Q5. How do you read an EXPLAIN plan? What does "Seq Scan" mean?

> **Answer:**
> `EXPLAIN` shows the execution plan the query optimizer chose. Key nodes:
> - **Seq Scan** — full table scan. Bad for large tables; indicates a missing or unused index.
> - **Index Scan** — uses an index to find rows, then fetches from heap.
> - **Index Only Scan** — uses a covering index; no heap access — fastest.
> - **cost=X..Y** — estimated cost (startup..total) in planner units.
> - `EXPLAIN ANALYZE` adds actual timing and row counts — use this to compare estimated vs actual rows (large discrepancies indicate stale statistics → run `ANALYZE`).

---

### Q6. How would you optimize a slow query on a large table?

> **Answer:**
> 1. Run `EXPLAIN ANALYZE` to identify the bottleneck.
> 2. Check if it's doing a **Seq Scan** on a large table — add an appropriate index.
> 3. Make predicates **sargable** — avoid functions on indexed columns; avoid leading wildcards.
> 4. Consider a **composite index** if filtering on multiple columns — put equality columns first.
> 5. Consider a **covering index** if the same columns are frequently selected and filtered.
> 6. Check for **stale statistics** — run `ANALYZE`.
> 7. Rewrite correlated subqueries as JOINs or CTEs.
> 8. Avoid `SELECT *` — fetch only needed columns.

---

<div align="center">
  <sub>← Back to <a href="Topic.md">All Topics</a></sub>
</div>
