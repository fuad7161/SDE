# 🟣 Core Concepts

> **Category:** Fundamentals &nbsp;|&nbsp; **Tags:** `ACID` `CAP theorem` `normalization` `keys`

---

## Table of Contents
1. [ACID Properties](#acid-properties)
2. [CAP Theorem](#cap-theorem)
3. [Normalization](#normalization)
4. [Database Keys](#database-keys)
5. [Interview Questions](#interview-questions)

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
❌ Before 3NF:
| StudentID | ZipCode | City   |
→ City depends on ZipCode, not on StudentID (transitive)

✅ After 3NF:
Students(StudentID, ZipCode)
ZipCodes(ZipCode, City)
```

#### BCNF — Boyce-Codd Normal Form
Stricter than 3NF: for every functional dependency `X → Y`, X must be a **superkey**.

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

<div align="center">
  <sub>← Back to <a href="Topic.md">All Topics</a></sub>
</div>
