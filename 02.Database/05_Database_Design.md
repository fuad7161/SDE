# 🟣 Database Design

> **Category:** Fundamentals &nbsp;|&nbsp; **Tags:** `ER diagrams` `many-to-many` `soft delete`

---

## Table of Contents
1. [ER Diagrams](#er-diagrams)
2. [Relationships & Cardinality](#relationships--cardinality)
3. [Implementing Relationships in SQL](#implementing-relationships-in-sql)
4. [Schema Design Trade-offs](#schema-design-trade-offs)
5. [Soft Delete Patterns](#soft-delete-patterns)
6. [Interview Questions](#interview-questions)

---

## ER Diagrams

An **Entity-Relationship (ER) diagram** is a visual blueprint of a database schema, showing entities, their attributes, and relationships.

### Notation (Crow's Foot)
```
[Users] ──|o────o<── [Orders]
   ↑          ↑          ↑
 Entity    Cardinality  Entity

||  = exactly one
|o  = zero or one
o<  = zero or many
|<  = one or many
```

### Components
- **Entity:** A table (noun) — `Users`, `Orders`, `Products`
- **Attribute:** A column — `name`, `email`, `created_at`
- **Relationship:** The connection — "User _places_ Order"
- **Cardinality:** How many of each entity participate

---

## Relationships & Cardinality

### One-to-One (1:1)
Each row in A relates to exactly one row in B.

```
[User] 1──── 1 [UserProfile]
```

```sql
CREATE TABLE user_profiles (
    user_id   INT PRIMARY KEY,
    bio       TEXT,
    avatar    VARCHAR(255),
    FOREIGN KEY (user_id) REFERENCES users(id)
);
-- user_id is both PK and FK — enforces 1:1
```

**Example:** user ↔ user_profile, person ↔ passport.

---

### One-to-Many (1:N)
One row in A relates to many rows in B. Most common relationship.

```
[Department] 1────< [Employees]
```

```sql
CREATE TABLE employees (
    id        SERIAL PRIMARY KEY,
    name      VARCHAR(100),
    dept_id   INT NOT NULL,
    FOREIGN KEY (dept_id) REFERENCES departments(id)
);
-- FK on the "many" side
```

**Example:** department ↔ employees, user ↔ orders.

---

### Many-to-Many (M:N)
One row in A relates to many rows in B, and vice versa. Requires a **junction/bridge table**.

```
[Students] ><──── [Enrollments] ────>< [Courses]
```

```sql
CREATE TABLE students  (id SERIAL PRIMARY KEY, name VARCHAR(100));
CREATE TABLE courses   (id SERIAL PRIMARY KEY, name VARCHAR(100));

-- Junction table
CREATE TABLE enrollments (
    student_id  INT NOT NULL,
    course_id   INT NOT NULL,
    enrolled_at TIMESTAMP DEFAULT NOW(),
    grade       CHAR(2),
    PRIMARY KEY (student_id, course_id),
    FOREIGN KEY (student_id) REFERENCES students(id),
    FOREIGN KEY (course_id)  REFERENCES courses(id)
);
```

**Example:** students ↔ courses, users ↔ roles, products ↔ orders.

---

## Implementing Relationships in SQL

### Referential Integrity — ON DELETE / ON UPDATE actions

```sql
FOREIGN KEY (dept_id) REFERENCES departments(id)
    ON DELETE SET NULL      -- employee's dept_id → NULL if dept deleted
    ON DELETE CASCADE       -- employee deleted when dept deleted
    ON DELETE RESTRICT      -- prevent dept deletion if employees exist
    ON DELETE NO ACTION     -- same as RESTRICT (default)
    ON UPDATE CASCADE       -- propagate PK changes to FK
```

### Self-Referential (Hierarchical)
```sql
-- Employee hierarchy — manager is also an employee
CREATE TABLE employees (
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(100),
    manager_id INT,
    FOREIGN KEY (manager_id) REFERENCES employees(id)
);
```

---

## Schema Design Trade-offs

### Normalized Schema
- Data stored in separate related tables.
- Pros: No redundancy, easier updates, less storage.
- Cons: Complex queries require many JOINs.

### Denormalized (Flat) Schema
- Frequently joined data merged into one table.
- Pros: Faster reads, simpler queries.
- Cons: Data redundancy, harder updates, more storage.

### When to denormalize
```sql
-- Normalized: 3 JOINs to get order details
SELECT o.id, u.name, p.name, oi.qty, oi.price
FROM orders o
JOIN users u ON o.user_id = u.id
JOIN order_items oi ON oi.order_id = o.id
JOIN products p ON oi.product_id = p.id;

-- Denormalized orders_flat table for read-heavy reporting
-- Stores user_name, product_name directly — no JOINs needed
SELECT * FROM orders_flat WHERE user_id = 42;
```

### Wide vs Narrow tables
- **Narrow:** Many small tables, highly normalized.
- **Wide:** Fewer tables with many columns, some redundancy.
- Wide tables common in OLAP/analytics (star schema, columnar stores).

---

## Soft Delete Patterns

**Soft delete:** Mark a record as deleted without physically removing it. Preserves history, allows recovery, maintains referential integrity.

### Pattern 1 — Boolean flag
```sql
ALTER TABLE users ADD COLUMN is_deleted BOOLEAN DEFAULT FALSE;

-- Soft delete
UPDATE users SET is_deleted = TRUE WHERE id = 42;

-- Query active users only
SELECT * FROM users WHERE is_deleted = FALSE;

-- ❌ Problem: Easy to forget the WHERE clause and expose deleted records
```

### Pattern 2 — Deleted timestamp (recommended)
```sql
ALTER TABLE users ADD COLUMN deleted_at TIMESTAMP;

-- Soft delete
UPDATE users SET deleted_at = NOW() WHERE id = 42;

-- Query active
SELECT * FROM users WHERE deleted_at IS NULL;

-- ✅ Also stores when it was deleted (useful for auditing/retention)
```

### Pattern 3 — Separate archive table
```sql
-- Move deleted rows to a separate archive
INSERT INTO users_archive SELECT * FROM users WHERE id = 42;
DELETE FROM users WHERE id = 42;

-- Active table stays clean; archive table holds history
-- ✅ Best performance on active table; ✅ no filter needed
-- ❌ More complex logic to manage
```

### Partial indexes for performance
```sql
-- Index only active (non-deleted) rows
CREATE INDEX idx_active_users ON users(email)
WHERE deleted_at IS NULL;
-- Index is small and fast — deleted rows aren't indexed
```

---

## Interview Questions

### Q1. How do you implement a many-to-many relationship in SQL?

> **Answer:**
> Create a **junction table** that holds foreign keys to both entities. This table's primary key is the composite of both foreign keys.
>
> ```sql
> -- students ↔ courses (many-to-many)
> CREATE TABLE enrollments (
>     student_id INT REFERENCES students(id),
>     course_id  INT REFERENCES courses(id),
>     enrolled_at TIMESTAMP,
>     PRIMARY KEY (student_id, course_id)
> );
> ```
> The junction table can also carry additional attributes about the relationship (e.g., `grade`, `enrolled_at`).

---

### Q2. What is the difference between ON DELETE CASCADE and ON DELETE SET NULL?

> **Answer:**
> These are referential integrity actions on a foreign key:
> - **CASCADE:** When the referenced parent row is deleted, **all child rows are also deleted** automatically.
> - **SET NULL:** When the referenced parent row is deleted, the **FK column in child rows is set to NULL** (the child row remains).
>
> Use CASCADE for strong ownership (a comment belongs to a post — delete post → delete comments). Use SET NULL when the child can exist independently (an employee can exist even if their department is removed).

---

### Q3. What is a soft delete? What are the trade-offs?

> **Answer:**
> A soft delete marks a row as deleted (via a `deleted_at` timestamp or `is_deleted` flag) without physically removing it.
>
> **Pros:** Data preserved for auditing, easy recovery, referential integrity maintained, historical reporting.
>
> **Cons:** Active queries need `WHERE deleted_at IS NULL` — easy to forget. Table grows unbounded. Indexes include "deleted" rows (use partial indexes to mitigate).
>
> **Best practice:** Use `deleted_at TIMESTAMP` (instead of boolean) to also capture when it was deleted. Use partial index on `(col) WHERE deleted_at IS NULL`.

---

### Q4. What is a self-referential (recursive) relationship? Give an example.

> **Answer:**
> A self-referential relationship is when a table has a foreign key pointing back to itself — used to model hierarchies.
>
> ```sql
> CREATE TABLE employees (
>     id         INT PRIMARY KEY,
>     name       VARCHAR(100),
>     manager_id INT REFERENCES employees(id)  -- FK to same table
> );
> ```
> An employee's `manager_id` points to another row in the same `employees` table. Query the full hierarchy using a **recursive CTE**. Other examples: category trees, org charts, threaded comments.

---

### Q5. When would you denormalize a schema?

> **Answer:**
> Denormalization trades storage/update complexity for read performance:
> - **Read-heavy workloads** where JOINs across many tables create bottlenecks.
> - **Reporting/analytics** — flat tables or pre-aggregated tables avoid complex JOINs.
> - **Caching frequently accessed joined data** — store derived values.
> - **Event sourcing / audit tables** — snapshot state at a point in time.
>
> Evaluate when profiling shows JOIN cost dominates query time, and the data changes infrequently (so redundancy doesn't cause many update anomalies).

---

<div align="center">
  <sub>← Back to <a href="Topic.md">All Topics</a></sub>
</div>
