# 🟠 Stored Procedures & Triggers

> **Category:** Advanced &nbsp;|&nbsp; **Tags:** `stored proc` `triggers` `materialized view`

---

## Table of Contents
1. [Stored Procedures vs Functions](#stored-procedures-vs-functions)
2. [Triggers](#triggers)
3. [Views vs Materialized Views](#views-vs-materialized-views)
4. [Interview Questions](#interview-questions)

---

## Stored Procedures vs Functions

### Stored Procedure
A **stored procedure** is a named, precompiled block of SQL (and procedural) code stored on the database server. Called explicitly via `CALL` or `EXEC`.

```sql
-- PostgreSQL: PL/pgSQL stored procedure
CREATE OR REPLACE PROCEDURE transfer_funds(
    sender_id   INT,
    receiver_id INT,
    amount      NUMERIC
)
LANGUAGE plpgsql AS $$
DECLARE
    sender_balance NUMERIC;
BEGIN
    -- Check balance
    SELECT balance INTO sender_balance FROM accounts WHERE id = sender_id FOR UPDATE;

    IF sender_balance < amount THEN
        RAISE EXCEPTION 'Insufficient funds: balance = %, requested = %',
                        sender_balance, amount;
    END IF;

    UPDATE accounts SET balance = balance - amount WHERE id = sender_id;
    UPDATE accounts SET balance = balance + amount WHERE id = receiver_id;

    INSERT INTO audit_log(from_id, to_id, amount, created_at)
    VALUES (sender_id, receiver_id, amount, NOW());

    COMMIT;
END;
$$;

-- Call it
CALL transfer_funds(1, 2, 500.00);
```

### Function (User-Defined Function — UDF)
A **function** computes and **returns a value** (scalar, table, or set). Can be used inside SELECT, WHERE, and other SQL clauses.

```sql
-- Returns a scalar value
CREATE OR REPLACE FUNCTION get_employee_rank(emp_id INT)
RETURNS INT
LANGUAGE plpgsql AS $$
DECLARE
    rank_val INT;
BEGIN
    SELECT DENSE_RANK() OVER (ORDER BY salary DESC)
    INTO rank_val
    FROM employees
    WHERE id = emp_id;

    RETURN rank_val;
END;
$$;

-- Use inside a query
SELECT name, get_employee_rank(id) AS rank FROM employees;
```

### Table-Valued Function
```sql
-- Returns a table — can be used in FROM clause
CREATE OR REPLACE FUNCTION get_dept_employees(dept_id_param INT)
RETURNS TABLE(id INT, name VARCHAR, salary NUMERIC)
LANGUAGE plpgsql AS $$
BEGIN
    RETURN QUERY
    SELECT e.id, e.name, e.salary
    FROM employees e
    WHERE e.dept_id = dept_id_param;
END;
$$;

-- Use it
SELECT * FROM get_dept_employees(3);
```

### Stored Procedure vs Function

| Feature | Stored Procedure | Function |
|---------|-----------------|---------|
| Returns value | Optional (OUT params) | Required (return type declared) |
| Used in SQL | `CALL proc()` — cannot be used inline | Can be used in `SELECT`, `WHERE`, `JOIN` |
| Transactions | Can manage transactions (COMMIT/ROLLBACK) | Generally cannot manage transactions |
| Side effects | Allowed (INSERT/UPDATE/DELETE) | Should be pure (side-effects vary by DB) |
| Use case | Business logic, multi-step operations | Computation, reusable derived values |

---

## Triggers

A **trigger** is a database object that **automatically executes a function** in response to a specific event (INSERT, UPDATE, DELETE) on a table.

### Trigger types

| Timing | Event | Description |
|--------|-------|-------------|
| BEFORE | INSERT/UPDATE/DELETE | Runs before the DML operation; can modify or cancel it |
| AFTER | INSERT/UPDATE/DELETE | Runs after the operation; can't cancel it |
| INSTEAD OF | INSERT/UPDATE/DELETE | Replaces the operation (used on views) |

### Example 1 — Audit Log Trigger
```sql
-- Audit table
CREATE TABLE employees_audit (
    id          SERIAL PRIMARY KEY,
    employee_id INT,
    action      VARCHAR(10),    -- INSERT, UPDATE, DELETE
    old_salary  NUMERIC,
    new_salary  NUMERIC,
    changed_at  TIMESTAMP DEFAULT NOW(),
    changed_by  TEXT DEFAULT current_user
);

-- Trigger function
CREATE OR REPLACE FUNCTION log_salary_change()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'UPDATE' AND OLD.salary <> NEW.salary THEN
        INSERT INTO employees_audit(employee_id, action, old_salary, new_salary)
        VALUES (NEW.id, 'UPDATE', OLD.salary, NEW.salary);
    ELSIF TG_OP = 'DELETE' THEN
        INSERT INTO employees_audit(employee_id, action, old_salary, new_salary)
        VALUES (OLD.id, 'DELETE', OLD.salary, NULL);
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Attach trigger to table
CREATE TRIGGER trg_salary_audit
AFTER UPDATE OR DELETE ON employees
FOR EACH ROW EXECUTE FUNCTION log_salary_change();
```

### Example 2 — BEFORE trigger to validate data
```sql
CREATE OR REPLACE FUNCTION validate_salary()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.salary < 0 THEN
        RAISE EXCEPTION 'Salary cannot be negative: %', NEW.salary;
    END IF;
    -- Auto-set updated_at
    NEW.updated_at = NOW();
    RETURN NEW;  -- RETURN NEW applies the (possibly modified) row
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_validate_salary
BEFORE INSERT OR UPDATE ON employees
FOR EACH ROW EXECUTE FUNCTION validate_salary();
```

### Trigger pitfalls
- **Hidden logic:** Triggers fire invisibly — makes debugging hard.
- **Cascading triggers:** A trigger fires another trigger → unexpected chains.
- **Performance impact:** Complex trigger logic on every DML slows bulk operations.
- **Order dependency:** Multiple triggers on same event — execution order can be unclear.
- **Testing difficulty:** Need full DB environment to test trigger logic.

**Best practices:**
- Keep trigger functions short and focused.
- Use triggers for cross-cutting concerns (audit, validation) — not business logic.
- Document all triggers thoroughly.

---

## Views vs Materialized Views

### Regular View
A **view** is a **saved SQL query** — a virtual table. It doesn't store data; every query against the view re-executes the underlying SQL.

```sql
-- Create view
CREATE VIEW active_employees AS
SELECT e.id, e.name, e.salary, d.name AS department
FROM employees e
JOIN departments d ON e.dept_id = d.id
WHERE e.deleted_at IS NULL;

-- Use like a table
SELECT * FROM active_employees WHERE department = 'Engineering';

-- Drop view
DROP VIEW active_employees;
```

**Pros:** Always returns fresh data; reduces query complexity; security (expose limited columns).
**Cons:** No performance gain — query executes fresh every time; complex views can be slow.

---

### Materialized View
A **materialized view** stores the **physical result set** of a query on disk. Must be **refreshed** to stay current.

```sql
-- Create materialized view
CREATE MATERIALIZED VIEW dept_salary_summary AS
SELECT
    d.name AS department,
    COUNT(e.id) AS headcount,
    AVG(e.salary) AS avg_salary,
    SUM(e.salary) AS total_payroll
FROM employees e
JOIN departments d ON e.dept_id = d.id
GROUP BY d.name;

-- Create index on materialized view
CREATE INDEX idx_mv_dept ON dept_salary_summary(department);

-- Refresh (re-run the query, update stored data)
REFRESH MATERIALIZED VIEW dept_salary_summary;

-- Refresh without locking reads (PostgreSQL)
REFRESH MATERIALIZED VIEW CONCURRENTLY dept_salary_summary;
```

**Pros:** Very fast reads (pre-computed, indexed); offload complex aggregation from hot path.
**Cons:** Stale data between refreshes; storage cost; must manage refresh schedule.

### View vs Materialized View

| Feature | View | Materialized View |
|---------|------|------------------|
| Stores data | ❌ (virtual) | ✅ (physical) |
| Always fresh | ✅ | ❌ (manual refresh needed) |
| Can be indexed | ❌ | ✅ |
| Query speed | Same as base query | Pre-computed — fast |
| Storage | None | Disk space required |
| Best for | Simplifying queries, security | Slow aggregations, reporting |

---

## Interview Questions

### Q1. What is the difference between a stored procedure and a function?

> **Answer:**
> - **Function:** Returns a value; can be used inline in SQL (in SELECT, WHERE). Should not manage transactions or have unpredictable side effects.
> - **Stored procedure:** Does not need to return a value; called with `CALL`. Can manage transactions (COMMIT/ROLLBACK). Better for multi-step business operations (e.g., order processing, fund transfers).
>
> Functions are for computation/reuse in queries. Procedures are for named, multi-step database operations.

---

### Q2. What are triggers? When would you use them?

> **Answer:**
> Triggers are automatic callbacks that fire before or after INSERT/UPDATE/DELETE on a table.
>
> **Good use cases:**
> - **Audit logging:** Record who changed what and when.
> - **Data validation:** Enforce complex business rules the DB constraint system can't express.
> - **Automatic timestamps:** Set `updated_at = NOW()` on every update.
> - **Derived/summary tables:** Keep a denormalized summary in sync.
>
> **Avoid for:** Core business logic — triggers are invisible, hard to test, and can cause performance surprises on bulk operations.

---

### Q3. What is the difference between a view and a materialized view?

> **Answer:**
> - **View:** A saved SQL query — virtual, no stored data. Every query re-runs the underlying SQL. Always fresh but has no performance advantage.
> - **Materialized view:** Stores the physical result of the query. Must be refreshed manually or on schedule. Supports indexes — very fast for repeated complex aggregations/reports.
>
> Use materialized views for expensive, infrequently-changing aggregations (daily reports, dashboards). Use regular views for simplifying access to complex joins or for row/column security.

---

### Q4. What are the pitfalls of using triggers in production?

> **Answer:**
> 1. **Hidden behavior:** Triggers fire automatically — developers not aware of them can be confused when data changes unexpectedly.
> 2. **Cascading chains:** Trigger A modifies table B, which fires trigger B, and so on — hard to trace.
> 3. **Performance impact:** Complex trigger logic multiplies with every DML, especially in bulk imports.
> 4. **Ordering issues:** Multiple triggers on the same event — execution order may not be deterministic.
> 5. **Testing difficulty:** Triggers require a full DB environment; hard to unit test.
>
> Mitigate by keeping trigger functions minimal, fully documenting them, and preferring application-layer logic where possible.

---

### Q5. How would you use a materialized view to improve dashboard performance?

> **Answer:**
> ```sql
> -- Expensive query run by dashboard 1000x/day
> SELECT region, SUM(revenue), COUNT(*) FROM orders
> WHERE created_at >= date_trunc('month', NOW())
> GROUP BY region;
>
> -- Materialize it
> CREATE MATERIALIZED VIEW monthly_revenue_by_region AS
> SELECT region, SUM(revenue) AS total_revenue, COUNT(*) AS order_count
> FROM orders
> WHERE created_at >= date_trunc('month', NOW())
> GROUP BY region;
>
> CREATE INDEX ON monthly_revenue_by_region(region);
>
> -- Refresh every hour via pg_cron or cron job
> REFRESH MATERIALIZED VIEW CONCURRENTLY monthly_revenue_by_region;
>
> -- Dashboard queries the MV instead
> SELECT * FROM monthly_revenue_by_region ORDER BY total_revenue DESC;
> ```
> This converts a slow aggregation query into a fast indexed table read.

---

<div align="center">
  <sub>← Back to <a href="Topic.md">All Topics</a></sub>
</div>
