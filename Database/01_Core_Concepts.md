# 🟣 Core Concepts

> **Category:** Fundamentals &nbsp;|&nbsp; **Tags:** `ACID` `CAP theorem` `normalization` `keys` `data types` `constraints` `OLTP vs OLAP` `architecture`

---

## Table of Contents
1. [ACID Properties](#acid-properties)
2. [CAP Theorem](#cap-theorem)
3. [Normalization](#normalization)
4. [Database Keys](#database-keys)
5. [Data Types](#data-types)
6. [Constraints](#constraints)
7. [OLTP vs OLAP](#oltp-vs-olap)
8. [Database Architecture](#database-architecture)
9. [Interview Questions](#interview-questions)

---

## ACID Properties

**ACID** defines the properties that guarantee database transactions are processed reliably.

### A — Atomicity
> "All or nothing."

A transaction is treated as a single unit — either **all operations succeed** or **none are applied**.

```sql
-- Transfer $100 from Alice to Bob — must be atomic
BEGIN;
  UPDATE accounts SET balance = balance - 100 WHERE user = 'Alice';
  UPDATE accounts SET balance = balance + 100 WHERE user = 'Bob';
COMMIT;
-- If second UPDATE fails, the first is ROLLED BACK automatically
```

### C — Consistency
> "Data must always move from one valid state to another."

Every transaction must leave the database in a state that satisfies all defined rules (constraints, cascades, triggers). A transaction that would violate a constraint is rolled back.

```sql
-- Constraint: balance >= 0
UPDATE accounts SET balance = balance - 500 WHERE user = 'Alice';
-- If Alice has $300, this violates the constraint → transaction rolled back
```

### I — Isolation
> "Concurrent transactions behave as if they run serially."

Changes made by an in-progress transaction are not visible to other transactions until committed (degree of isolation depends on the **isolation level**).

```sql
-- T1: Transfer $200 from Alice to Bob
BEGIN;
  SELECT balance FROM accounts WHERE name = 'Alice';  -- 1000
  UPDATE accounts SET balance = 800 WHERE name = 'Alice';
  UPDATE accounts SET balance = 1200 WHERE name = 'Bob';
COMMIT;

-- T2 (running concurrently): Calculate total money in the system
BEGIN;
  SELECT SUM(balance) FROM accounts;
COMMIT;
```

**What isolation prevents:**
- If T2 runs its `SELECT SUM` *after* T1 debits Alice but *before* it credits Bob, it sees `800 + 1000 = 1800` instead of `2000` — money appears to have vanished. This is a **dirty read / non-repeatable read**.
- With proper isolation (e.g., `REPEATABLE READ` or `SERIALIZABLE`), T2 either sees the state *before* T1 started (`1000 + 1000 = 2000`) or *after* T1 committed (`800 + 1200 = 2000`) — never an in-between inconsistent snapshot.

### D — Durability
> "Once committed, data survives failures."

Committed transactions are written to **non-volatile storage** (WAL — Write-Ahead Log). Even if the server crashes immediately after commit, the data is not lost.

---

## CAP Theorem

**CAP Theorem** (Brewer, 2000): A distributed data store can provide at most **two of three** guarantees simultaneously:

| Property | Meaning |
|----------|---------|
| **C — Consistency** | Every read returns the most recent write (or an error) |
| **A — Availability** | Every request receives a response (no error), even if it may be stale |
| **P — Partition Tolerance** | System continues operating despite network partitions between nodes |

**In practice:** Network partitions are inevitable in distributed systems, so **P is always required**. The real trade-off is **C vs A** during a partition.

```
           CAP Triangle
               C
              / \
    CP DBs   /   \  CA DBs
            /     \
           A-------P
              AP DBs
```

### Real-world examples

| System | CAP Choice | Reasoning |
|--------|-----------|-----------|
| PostgreSQL (single node) | CA | Not distributed; partition not applicable |
| MySQL Cluster | CP | Stops responding rather than returning stale data |
| Cassandra | AP | Always accepts writes; resolves conflicts eventually |
| DynamoDB | AP (tunable) | Eventual consistency by default, strong consistency optional |
| HBase | CP | Stops serving rather than returning inconsistent data |
| Zookeeper | CP | Coordination requires strong consistency |

### PACELC Extension
CAP only describes behavior during partitions. **PACELC** extends it:
- If Partition → choose C or A
- **Else** (no partition) → choose **Latency** or **Consistency**

---

## Normalization

Normalization organizes a database to **reduce redundancy** and **improve data integrity** by decomposing tables.

### Problems without normalization
- **Insertion anomaly:** Can't add data without other unrelated data.
- **Update anomaly:** Same fact stored in multiple rows — must update all.
- **Deletion anomaly:** Deleting a row removes unrelated information.

### Normal Forms

#### 1NF — First Normal Form
**Rule:** Each column holds **atomic (indivisible) values**. No repeating groups.

```
❌ Before 1NF:
| StudentID | Name  | Courses          |
|-----------|-------|------------------|
| 1         | Alice | Math, Physics    |  ← multi-valued

✅ After 1NF:
| StudentID | Name  | Course   |
|-----------|-------|----------|
| 1         | Alice | Math     |
| 1         | Alice | Physics  |
```

#### 2NF — Second Normal Form
**Rule:** In 1NF + **no partial dependencies** (non-key columns depend on the full primary key, not just part of it).

*Only relevant when the PK is composite.*

```
❌ Before 2NF (PK = StudentID + CourseID):
| StudentID | CourseID | CourseName  | Grade |
|-----------|----------|-------------|-------|
| 1         | C01      | Mathematics | A     |
| 1         | C02      | Physics     | B     |
| 2         | C01      | Mathematics | C     |
| 3         | C02      | Physics     | A     |

Problem: "Mathematics" is stored 2× and "Physics" 2×.
CourseName depends only on CourseID, not on the full PK (StudentID + CourseID).
→ If you rename "Mathematics" you must update multiple rows (update anomaly).

✅ After 2NF: Split into two tables

Enrollments (StudentID + CourseID → Grade):
| StudentID | CourseID | Grade |
|-----------|----------|-------|
| 1         | C01      | A     |
| 1         | C02      | B     |
| 2         | C01      | C     |
| 3         | C02      | A     |

Courses (CourseID → CourseName):
| CourseID | CourseName  |
|----------|-------------|
| C01      | Mathematics |
| C02      | Physics     |

Now "Mathematics" lives in exactly one place — renaming it touches one row.
```

#### 3NF — Third Normal Form
**Rule:** In 2NF + **no transitive dependencies** (non-key columns don't depend on other non-key columns).

```
❌ Before 3NF (PK = StudentID):
| StudentID | Name    | ZipCode | City          |
|-----------|---------|---------|---------------|
| 1         | Alice   | 10001   | New York      |
| 2         | Bob     | 90001   | Los Angeles   |
| 3         | Charlie | 10001   | New York      |
| 4         | Diana   | 90001   | Los Angeles   |

Problem: "New York" and "Los Angeles" are stored multiple times.
City depends on ZipCode, and ZipCode depends on StudentID.
StudentID → ZipCode → City  (transitive dependency)
→ If zip code 10001 moves to a new city, you must update multiple rows (update anomaly).
→ If you delete Alice and Charlie, you lose the fact that 10001 = New York (deletion anomaly).

✅ After 3NF: Split into two tables

Students (StudentID → Name, ZipCode):
| StudentID | Name    | ZipCode |
|-----------|---------|---------|
| 1         | Alice   | 10001   |
| 2         | Bob     | 90001   |
| 3         | Charlie | 10001   |
| 4         | Diana   | 90001   |

ZipCodes (ZipCode → City):
| ZipCode | City          |
|---------|---------------|
| 10001   | New York      |
| 90001   | Los Angeles   |

Now each city name lives in exactly one place — changing a city touches one row.
```

#### BCNF — Boyce-Codd Normal Form
Stricter than 3NF: for every functional dependency `X → Y`, X must be a **superkey**.

```
Context: Each teacher teaches only one subject. A student can have one teacher per subject.
Candidate keys: (StudentID, TeacherID) and (StudentID, Subject) — both uniquely identify a row.

❌ In 3NF but NOT in BCNF:
| StudentID | TeacherID | Subject   |
|-----------|-----------|-----------|
| 1         | T1        | Math      |
| 1         | T2        | Physics   |
| 2         | T1        | Math      |
| 2         | T3        | Physics   |
| 3         | T2        | Physics   |

Functional dependency: TeacherID → Subject
(T1 always teaches Math, T2 always teaches Physics, T3 always teaches Physics)

Problem: TeacherID → Subject exists, but TeacherID alone is NOT a superkey.
→ "T1 teaches Math" is stored 2× — if T1 switches to Chemistry, multiple rows must change.

✅ After BCNF: Split so every determinant is a superkey

TeacherSubject (TeacherID → Subject):
| TeacherID | Subject  |
|-----------|----------|
| T1        | Math     |
| T2        | Physics  |
| T3        | Physics  |

StudentTeacher (StudentID + TeacherID → the enrollment fact):
| StudentID | TeacherID |
|-----------|-----------|
| 1         | T1        |
| 1         | T2        |
| 2         | T1        |
| 2         | T3        |
| 3         | T2        |

Now every functional dependency has a superkey on the left-hand side.
```

### When to Denormalize
- Read-heavy workloads where JOIN cost is too high.
- Reporting/analytics (flat denormalized tables are faster to scan).
- Pre-computed aggregates for dashboards.
- Trade-off: faster reads, higher storage, more complex writes.

---

## Database Keys

| Key Type | Description | Example |
|----------|-------------|---------|
| **Primary Key** | Uniquely identifies each row; NOT NULL, unique | `user_id INT PRIMARY KEY` |
| **Foreign Key** | References PK of another table; enforces referential integrity | `order.user_id → users.user_id` |
| **Candidate Key** | Any column (or set) that could be a PK | email, phone, SSN |
| **Composite Key** | PK made of multiple columns | `(student_id, course_id)` |
| **Surrogate Key** | Artificially generated PK (no business meaning) | `id SERIAL`, `UUID` |
| **Natural Key** | PK derived from real-world data | SSN, email |
| **Unique Key** | Unique but allows one NULL | `UNIQUE (email)` |

### Surrogate vs Natural Key

| | Surrogate Key | Natural Key |
|--|-------------|------------|
| Stability | Always stable | May change (e.g., email changes) |
| Business meaning | None | Has meaning |
| Performance | Compact int/UUID — fast joins | String keys — slower joins |
| Privacy | Safer (no PII in FK) | Exposes PII |
| Recommendation | ✅ Preferred for most cases | Use only when truly immutable |

---

## Data Types

Choosing the right data type affects **storage cost**, **query performance**, and **data correctness**.

### Numeric Types

| Type | Storage | Range / Use |
|------|---------|-------------|
| `SMALLINT` | 2 bytes | -32,768 to 32,767 |
| `INTEGER` | 4 bytes | ~±2 billion |
| `BIGINT` | 8 bytes | ±9.2 × 10¹⁸ |
| `SERIAL` / `BIGSERIAL` | auto-increment | Surrogate PK (PostgreSQL) |
| `NUMERIC(p,s)` | variable | Exact precision — money, scientific |
| `REAL` / `DOUBLE PRECISION` | 4 / 8 bytes | Floating point — fast, approximate |

```sql
-- NUMERIC for exact money
price NUMERIC(10,2)  -- up to 99999999.99

-- SERIAL for auto-increment PK
id SERIAL PRIMARY KEY
```

### String Types

| Type | Use |
|------|-----|
| `CHAR(n)` | Fixed-length (padded) — rarely used |
| `VARCHAR(n)` | Variable-length with max limit |
| `TEXT` | Variable-length, no limit (PostgreSQL) |
| `UUID` | 128-bit unique identifier |

```sql
-- VARCHAR vs TEXT
name VARCHAR(100)  -- enforces max length
bio  TEXT          -- no limit — use when length is unpredictable

-- UUID
id UUID PRIMARY KEY DEFAULT gen_random_uuid()
```

### Date & Time Types

| Type | Storage | Use |
|------|---------|-----|
| `DATE` | 4 bytes | Calendar date (no time) |
| `TIME` | 8 bytes | Time of day |
| `TIMESTAMP` | 8 bytes | Date + time (no timezone) |
| `TIMESTAMPTZ` | 8 bytes | Date + time + timezone (recommended) |
| `INTERVAL` | 16 bytes | Duration (e.g., '2 hours 30 minutes') |

```sql
-- Always prefer TIMESTAMPTZ
created_at TIMESTAMPTZ DEFAULT NOW()
```

### Boolean

```sql
is_active BOOLEAN DEFAULT TRUE
```

### JSON / JSONB (PostgreSQL)

| Type | Storage | Features |
|------|---------|----------|
| `JSON` | Raw text | Stores exact input; re-parse on every access |
| `JSONB` | Binary | Indexed, fast queries, decomposed — **prefer this** |

```sql
-- JSONB with indexing
CREATE TABLE events (
    id SERIAL PRIMARY KEY,
    data JSONB
);

-- GIN index for key/containment queries
CREATE INDEX idx_events_data ON events USING GIN(data);

SELECT * FROM events WHERE data @> '{"type": "click"}';
```

### Binary / Large Objects

| Type | Use |
|------|-----|
| `BYTEA` | Variable-length binary (PostgreSQL) |
| `BLOB` | Binary large object (MySQL, others) |

**Rule of thumb:** Store files in object storage (S3), store the URL/path in the database.

### Choosing the Right Type

| Principle | Example |
|-----------|---------|
| Smallest that fits | `SMALLINT` over `BIGINT` for age |
| Exact for money | `NUMERIC` not `FLOAT` |
| `TIMESTAMPTZ` over `TIMESTAMP` | Always |
| `TEXT` over `VARCHAR` | When max length is unpredictable (PostgreSQL treats them identically for performance) |
| `JSONB` over `JSON` | When you need to query inside JSON |

---

## Constraints

Constraints enforce **data integrity** at the database level — they prevent invalid data regardless of application logic.

### Types of Constraints

| Constraint | Purpose | Example |
|------------|---------|---------|
| `PRIMARY KEY` | Uniquely identifies each row; NOT NULL + UNIQUE | `id INT PRIMARY KEY` |
| `FOREIGN KEY` | Referential integrity between tables | `REFERENCES users(id)` |
| `UNIQUE` | No duplicate values (allows one NULL) | `UNIQUE (email)` |
| `NOT NULL` | Column cannot be NULL | `name TEXT NOT NULL` |
| `CHECK` | Custom validation rule | `CHECK (age >= 0 AND age <= 150)` |
| `DEFAULT` | Auto-fill when no value provided | `status TEXT DEFAULT 'active'` |
| `EXCLUSION` | No two rows overlap for specified operators (PostgreSQL) | No double-booking a room |

### CHECK Constraints

```sql
-- Validate range
ALTER TABLE products ADD CONSTRAINT positive_price
    CHECK (price > 0);

-- Validate across columns
ALTER TABLE events ADD CONSTRAINT valid_date_range
    CHECK (start_time < end_time);

-- Complex rules
ALTER TABLE employees ADD CONSTRAINT valid_email
    CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$');
```

### EXCLUSION Constraints (PostgreSQL)

```sql
-- Prevent room double-booking: no overlapping reservations
CREATE EXTENSION IF NOT EXISTS btree_gist;

CREATE TABLE reservations (
    room_id INT,
    during TSRANGE,
    EXCLUDE USING GIST (room_id WITH =, during WITH &&)
);
-- Attempting to insert overlapping reservation fails automatically
```

### NOT NULL vs DEFAULT

```sql
-- NOT NULL + DEFAULT: column always has a value
status TEXT NOT NULL DEFAULT 'pending';

-- Without DEFAULT: INSERT without status fails
-- With DEFAULT: INSERT without status → 'pending'
```

### When to Use Constraints

| Rule | Reasoning |
|------|-----------|
| Always use `NOT NULL` unless NULL is meaningful | Prevents accidental missing data |
| Always use `FOREIGN KEY` for relationships | DB enforces referential integrity |
| Use `CHECK` for business rules | Faster than app-level validation; DB-level safety net |
| Use `DEFAULT` for common values | Reduces INSERT boilerplate |
| Use `EXCLUSION` for scheduling/resource conflicts | Prevents overlapping ranges at DB level |

---

## OLTP vs OLAP

Two fundamentally different database workload types that drive schema design, indexing, and scaling decisions.

### OLTP — Online Transaction Processing

Day-to-day operations: **many small, fast reads and writes**.

| Characteristic | Detail |
|----------------|--------|
| Queries | Point lookups, short transactions |
| Row count per query | 1 – 100 rows |
| Operations | INSERT, UPDATE, DELETE, SELECT by PK |
| Latency requirement | Milliseconds |
| Data model | Normalized (3NF+) |
| Example | Banking, e-commerce, user authentication |

```
OLTP Pattern:
  SELECT * FROM orders WHERE id = 12345;           -- point lookup
  UPDATE accounts SET balance = balance - 100 WHERE id = 42;
  INSERT INTO orders (user_id, total) VALUES (7, 99.99);
```

### OLAP — Online Analytical Processing

Business intelligence: **fewer, heavier read queries** scanning large datasets.

| Characteristic | Detail |
|----------------|--------|
| Queries | Aggregations, GROUP BY, JOINs across millions of rows |
| Row count per query | 10,000 – billions |
| Operations | SELECT with aggregations, rarely UPDATE |
| Latency requirement | Seconds to minutes |
| Data model | Denormalized (star/snowflake schema) |
| Example | Reports, dashboards, analytics |

```
OLAP Pattern:
  SELECT region, SUM(revenue), COUNT(*)
  FROM orders
  WHERE created_at >= '2024-01-01'
  GROUP BY region;
```

### Comparison

| Aspect | OLTP | OLAP |
|--------|------|------|
| Purpose | Run the business | Analyze the business |
| Data volume | GBs | TBs – PBs |
| Schema | Normalized (3NF) | Denormalized (star/snowflake) |
| Indexes | B-Tree on PKs/FKs | Columnar, bitmap |
| Scaling | Vertical | Horizontal (columnar stores) |
| DB examples | PostgreSQL, MySQL | ClickHouse, Snowflake, BigQuery, Redshift |

### Hybrid (HTAP)

Some systems support both workloads: PostgreSQL with Citus, TiDB, CockroachDB.

---

## Database Architecture

Understanding how a database works internally explains why certain design decisions matter.

### Storage Model

**Page/Block-based storage:**
- Data is stored in fixed-size **pages** (typically 8 KB).
- Pages are grouped into **extents** (contiguous blocks of pages).
- A **heap** (or tablespace) holds all pages for a table.

```
Table "orders"
  Page 0: [Row 1, Row 2, Row 3, ...]
  Page 1: [Row 100, Row 101, ...]
  Page 2: [Row 200, Row 201, ...]
  ...
```

**Row-based** (OLTP): Each page contains full rows — good for point lookups.

**Column-based** (OLAP): Each page contains one column's values — good for aggregations.

### Query Processing Pipeline

```
SQL Query
    ↓
┌─────────────┐
│  Parser     │  Syntax check, build parse tree
└──────┬──────┘
       ↓
┌─────────────┐
│  Analyzer   │  Resolve table/column names, type checking
└──────┬──────┘
       ↓
┌──────────────┐
│  Optimizer   │  Generate execution plans, pick cheapest (cost-based)
└──────┬───────┘
       ↓
┌──────────────┐
│  Executor    │  Run the plan — seq scan, index scan, joins, etc.
└──────┬───────┘
       ↓
   Result Set
```

**Cost-based optimizer** estimates cost of each plan using table statistics (row count, data distribution, index availability) and picks the lowest-cost plan.

### Buffer Pool

An **in-memory cache** of disk pages that sits between the database and storage.

```
Application Request
       ↓
┌──────────────┐
│ Buffer Pool  │  (in RAM)
│ Page cache   │  ← check if page is here
└──────┬───────┘
       │
  ┌────┴────┐
  │ Hit     │ Miss → read from disk → cache in buffer pool
  └─────────┘
```

- **Buffer hit:** Page is in memory — fast (microseconds).
- **Buffer miss:** Page must be read from disk — slow (milliseconds).
- **Eviction:** When full, LRU (or clock algorithm) evicts least-used pages.

**Why this matters:** A well-tuned buffer pool size (e.g., `shared_buffers` in PostgreSQL) is one of the most impactful performance settings.

### Write-Ahead Logging (WAL)

Every change is written to a **log file** before being applied to the actual data pages.

```
1. Write change to WAL (sequential append — fast)
2. Return success to client
3. Eventually: checkpoint — flush dirty pages to disk
```

**Why WAL exists:**
- **Crash recovery:** If the DB crashes, WAL replays committed changes.
- **Durability:** Once WAL is flushed, the commit is durable even if the data page hasn't been written yet.
- **Replication:** WAL shipping is how most replication works (stream WAL to replicas).

```
WAL Timeline:
  [WAL segment 1] [WAL segment 2] [WAL segment 3] → checkpoint → [clean data pages on disk]
```

### Checkpointing

Periodically, the database flushes all **dirty pages** (modified pages in buffer pool) to disk and records a checkpoint in the WAL. After checkpoint, older WAL segments can be recycled.

```
WAL:  [seg1][seg2][seg3][seg4]  →  checkpoint  →  [seg5][seg6]...
                                       ↑
                          dirty pages flushed to data files
```

### Transaction ID (PostgreSQL)

Every transaction gets a monotonically increasing **txid**. This is used for:
- **MVCC visibility:** Determine which row version a transaction can see.
- **Vacuum:** Clean up dead tuples older than the oldest active transaction.

---

## Interview Questions

### Q1. What are ACID properties? Explain with an example.

> **Answer:**
> - **Atomicity:** A transaction is all-or-nothing. If a bank transfer debits Alice but the credit to Bob fails, the debit is rolled back.
> - **Consistency:** Transactions move the DB from one valid state to another. A constraint violation rolls back the transaction.
> - **Isolation:** Concurrent transactions don't see each other's intermediate state. Controlled by isolation levels.
> - **Durability:** Once committed, data survives crashes. Implemented via WAL (Write-Ahead Log).

---

### Q2. Explain the CAP theorem. Which two can you have?

> **Answer:**
> CAP: Consistency, Availability, Partition Tolerance. Network partitions are unavoidable in distributed systems, so P is always required. The real trade-off is **CP vs AP**:
> - **CP** (e.g., HBase, Zookeeper): Returns an error rather than stale data during a partition.
> - **AP** (e.g., Cassandra, DynamoDB): Always responds, possibly with stale data; reconciles later.
>
> Single-node databases (PostgreSQL) are effectively CA — partitions don't apply.

---

### Q3. What is normalization? Explain 1NF, 2NF, and 3NF.

> **Answer:**
> Normalization reduces data redundancy and anomalies:
> - **1NF:** Atomic values, no repeating groups, each row uniquely identifiable.
> - **2NF:** 1NF + no partial dependencies (non-key columns depend on the **whole** composite PK).
> - **3NF:** 2NF + no transitive dependencies (non-key columns depend only on the PK, not on other non-key columns).
>
> Each normal form adds a stricter constraint on how data is organized.

---

### Q4. What is the difference between a primary key and a unique key?

> **Answer:**
> - **Primary Key:** Uniquely identifies each row. Cannot be NULL. Only one per table. Automatically creates a clustered index (in most DBs).
> - **Unique Key:** Ensures uniqueness in a column but **allows one NULL** (NULLs are not equal to each other). Multiple unique keys allowed per table.
>
> A table can have only one PK but multiple unique constraints.

---

### Q5. When would you choose a surrogate key over a natural key?

> **Answer:**
> Use a **surrogate key** (auto-increment int or UUID) when:
> - The natural key can change (email, phone number).
> - The natural key is long (string) — makes foreign keys and joins expensive.
> - The natural key contains PII — you don't want it appearing in foreign keys across tables.
>
> Use a **natural key** only when it is truly immutable, compact, and unique (rare in practice).

---

### Q6. What is the difference between consistency in ACID and consistency in CAP?

> **Answer:**
> - **ACID Consistency:** A transaction moves the database from one valid state to another, respecting all constraints and rules. It's about **intra-database correctness**.
> - **CAP Consistency:** In a distributed system, every read reflects the most recent write across all nodes. It's about **cross-node agreement** (linearizability).
>
> These are completely different concepts despite sharing the same word.

---

### Q7. When would you use NUMERIC vs FLOAT for storing numbers?

> **Answer:**
> - **NUMERIC (DECIMAL):** Exact precision — stores the exact value. Use for **money**, scientific data, and anything where rounding errors are unacceptable. Slower for arithmetic.
> - **FLOAT / DOUBLE PRECISION:** Approximate — fast binary floating point. Fine for physics simulations, ML weights, or any domain where small precision loss is acceptable.
>
> Rule: Never use FLOAT for financial calculations. Use `NUMERIC(10,2)` for currency.

---

### Q8. What is the difference between OLTP and OLAP? How do they affect schema design?

> **Answer:**
> - **OLTP (Online Transaction Processing):** Fast, small transactions (point lookups, single-row inserts/updates). Normalized schemas (3NF) to reduce redundancy. Examples: banking, e-commerce.
> - **OLAP (Online Analytical Processing):** Heavy read queries scanning millions of rows (aggregations, GROUP BY). Denormalized schemas (star/snowflake) to reduce JOINs. Examples: dashboards, business intelligence.
>
> OLTP optimizes for write speed and row-level consistency. OLAP optimizes for read speed and column-level scan efficiency.

---

### Q9. What happens inside a database when you run a SQL query?

> **Answer:**
> 1. **Parser:** Validates syntax, builds a parse tree.
> 2. **Analyzer:** Resolves table/column names, checks types.
> 3. **Optimizer:** Generates multiple execution plans, estimates cost using table statistics (row count, distribution, indexes), picks the cheapest plan.
> 4. **Executor:** Runs the plan — sequential scans, index scans, nested loop/hash/merge joins.
> 5. **Result set** returned to client.
>
> You can see the optimizer's choice using `EXPLAIN` or `EXPLAIN ANALYZE`.

---

### Q10. What is a buffer pool and why does it matter?

> **Answer:**
> The buffer pool is an **in-memory cache of data pages** that sits between the database engine and disk. When a page is read, it's cached in the buffer pool so subsequent reads of the same page are served from memory (microseconds) instead of disk (milliseconds).
>
> **Why it matters:**
> - A properly sized buffer pool (e.g., PostgreSQL `shared_buffers`) is one of the most impactful performance settings.
> - If the pool is too small, pages are evicted and re-read from disk frequently (cache thrashing).
> - If too large, it may starve the OS of memory.

---

<div align="center">
  <sub>← Back to <a href="Topic.md">All Topics</a></sub>
</div>
