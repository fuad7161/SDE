# 🟢 SQL & Queries

> **Category:** Queries &nbsp;|&nbsp; **Tags:** `JOINs` `window functions` `CTEs` `N+1`

---

## Table of Contents
1. [JOINs](#joins)
2. [Aggregations & GROUP BY](#aggregations--group-by)
3. [Window Functions](#window-functions)
4. [Subqueries vs CTEs](#subqueries-vs-ctes)
5. [Query Optimization & N+1](#query-optimization--n1)
6. [UNION vs UNION ALL / EXISTS vs IN](#union-vs-union-all--exists-vs-in)
7. [Interview Questions](#interview-questions)

---

## JOINs

```sql
-- Sample tables
-- employees(id, name, dept_id, salary, manager_id)
-- departments(id, name)
```

### INNER JOIN — Only matching rows from both tables
```sql
SELECT e.name, d.name AS department
FROM employees e
INNER JOIN departments d ON e.dept_id = d.id;
-- Returns only employees who have a valid dept_id
```

### LEFT JOIN — All rows from left + matched rows from right (NULLs for no match)
```sql
SELECT e.name, d.name AS department
FROM employees e
LEFT JOIN departments d ON e.dept_id = d.id;
-- Returns ALL employees; dept is NULL if no department assigned
```

### RIGHT JOIN — All rows from right + matched rows from left
```sql
SELECT e.name, d.name AS department
FROM employees e
RIGHT JOIN departments d ON e.dept_id = d.id;
-- Returns ALL departments; employee is NULL if dept has no employees
```

### FULL OUTER JOIN — All rows from both, NULLs where no match
```sql
SELECT e.name, d.name
FROM employees e
FULL OUTER JOIN departments d ON e.dept_id = d.id;
-- All employees + all departments, matched where possible
```

### CROSS JOIN — Cartesian product (every row × every row)
```sql
SELECT e.name, d.name
FROM employees e
CROSS JOIN departments d;
-- 10 employees × 5 departments = 50 rows
```

### SELF JOIN — Join a table with itself
```sql
-- Find each employee and their manager's name
SELECT e.name AS employee, m.name AS manager
FROM employees e
LEFT JOIN employees m ON e.manager_id = m.id;
```

### Finding rows with no match (Anti-Join)
```sql
-- Employees with no department (using LEFT JOIN)
SELECT e.name FROM employees e
LEFT JOIN departments d ON e.dept_id = d.id
WHERE d.id IS NULL;

-- Same using NOT EXISTS
SELECT name FROM employees e
WHERE NOT EXISTS (
    SELECT 1 FROM departments d WHERE d.id = e.dept_id
);
```

---

## Aggregations & GROUP BY

```sql
-- Basic aggregation
SELECT dept_id,
       COUNT(*)           AS headcount,
       AVG(salary)        AS avg_salary,
       MAX(salary)        AS max_salary,
       MIN(salary)        AS min_salary,
       SUM(salary)        AS total_payroll
FROM employees
GROUP BY dept_id;
```

### HAVING — Filter on aggregated results
```sql
-- Departments with more than 5 employees
SELECT dept_id, COUNT(*) AS headcount
FROM employees
GROUP BY dept_id
HAVING COUNT(*) > 5;
-- WHERE filters rows BEFORE aggregation
-- HAVING filters groups AFTER aggregation
```

### Execution order of a SELECT statement
```
FROM → JOIN → WHERE → GROUP BY → HAVING → SELECT → DISTINCT → ORDER BY → LIMIT
```

---

## Window Functions

Window functions perform calculations **across a set of rows related to the current row**, without collapsing them into a single output row (unlike GROUP BY).

```sql
SELECT
    name,
    dept_id,
    salary,
    -- Rank within each department
    ROW_NUMBER()  OVER (PARTITION BY dept_id ORDER BY salary DESC) AS row_num,
    RANK()        OVER (PARTITION BY dept_id ORDER BY salary DESC) AS rank,
    DENSE_RANK()  OVER (PARTITION BY dept_id ORDER BY salary DESC) AS dense_rank,

    -- Running total
    SUM(salary)   OVER (PARTITION BY dept_id ORDER BY salary DESC) AS running_total,

    -- Previous and next salary in same dept
    LAG(salary, 1)  OVER (PARTITION BY dept_id ORDER BY salary) AS prev_salary,
    LEAD(salary, 1) OVER (PARTITION BY dept_id ORDER BY salary) AS next_salary,

    -- Percentile rank
    PERCENT_RANK() OVER (PARTITION BY dept_id ORDER BY salary) AS pct_rank

FROM employees;
```

### ROW_NUMBER vs RANK vs DENSE_RANK

Given salaries: 100, 90, 90, 80

| salary | ROW_NUMBER | RANK | DENSE_RANK |
|--------|-----------|------|-----------|
| 100 | 1 | 1 | 1 |
| 90 | 2 | 2 | 2 |
| 90 | 3 | 2 | 2 |
| 80 | 4 | 4 | 3 |

- `ROW_NUMBER` — always unique, arbitrary tie-breaking
- `RANK` — ties get same rank, **skips** next rank (1,2,2,4)
- `DENSE_RANK` — ties get same rank, **no skip** (1,2,2,3)

### Find Nth highest salary using window function
```sql
-- 3rd highest salary
SELECT salary FROM (
    SELECT salary,
           DENSE_RANK() OVER (ORDER BY salary DESC) AS rnk
    FROM employees
) ranked
WHERE rnk = 3;
```

---

## Subqueries vs CTEs

### Subquery
```sql
-- Find employees earning more than the department average
SELECT name, salary, dept_id
FROM employees e
WHERE salary > (
    SELECT AVG(salary)
    FROM employees
    WHERE dept_id = e.dept_id   -- correlated subquery (executes once per row)
);
```

### CTE (Common Table Expression) — WITH clause
```sql
-- Same query using CTE — more readable, can be referenced multiple times
WITH dept_avg AS (
    SELECT dept_id, AVG(salary) AS avg_sal
    FROM employees
    GROUP BY dept_id
)
SELECT e.name, e.salary, e.dept_id
FROM employees e
JOIN dept_avg d ON e.dept_id = d.dept_id
WHERE e.salary > d.avg_sal;
```

### Recursive CTE — Traverse hierarchies
```sql
-- Get full management chain for employee ID 5
WITH RECURSIVE org_chart AS (
    -- Base case
    SELECT id, name, manager_id, 0 AS depth
    FROM employees
    WHERE id = 5

    UNION ALL

    -- Recursive case
    SELECT e.id, e.name, e.manager_id, oc.depth + 1
    FROM employees e
    JOIN org_chart oc ON e.id = oc.manager_id
)
SELECT * FROM org_chart;
```

### When to use what

| | Subquery | CTE |
|--|---------|-----|
| Readability | Hard to read when nested | Easier — named, modular |
| Reuse | Cannot reuse in same query | Can reference multiple times |
| Recursion | ❌ | ✅ |
| Performance | Correlated = slow (row-by-row) | Non-correlated = same as subquery |
| Debugging | Difficult | Easy — test CTE independently |

---

## Query Optimization & N+1

### The N+1 Problem
```
// Pseudo-code — fetching posts and their authors
posts = SELECT * FROM posts LIMIT 100;   -- 1 query

for each post:
    author = SELECT * FROM users WHERE id = post.author_id;  -- N queries!

Total: 1 + N queries → N+1 problem
```

**Fix:** Use a JOIN to fetch everything in one query:
```sql
SELECT p.*, u.name AS author_name
FROM posts p
JOIN users u ON p.author_id = u.id
LIMIT 100;
-- 1 query
```

Or use `IN` / batch fetch in application code.

### General optimization tips
```sql
-- ❌ Full table scan — avoid SELECT *
SELECT * FROM orders WHERE YEAR(created_at) = 2024;

-- ✅ Use sargable predicates (index-friendly)
SELECT id, total FROM orders
WHERE created_at >= '2024-01-01' AND created_at < '2025-01-01';

-- ❌ Function on indexed column = index not used
WHERE LOWER(email) = 'alice@example.com'

-- ✅ Store normalized form in DB, or use functional index
WHERE email = 'alice@example.com'

-- ❌ Leading wildcard = full scan
WHERE name LIKE '%alice%'

-- ✅ Prefix search uses index
WHERE name LIKE 'alice%'
```

---

## UNION vs UNION ALL / EXISTS vs IN

### UNION vs UNION ALL
```sql
-- UNION: combines results, removes duplicates (sorts internally — slower)
SELECT name FROM employees
UNION
SELECT name FROM contractors;

-- UNION ALL: combines results, keeps duplicates (faster — no dedup step)
SELECT name FROM employees
UNION ALL
SELECT name FROM contractors;
```
**Use UNION ALL when duplicates don't matter or are impossible** — it's always faster.

### EXISTS vs IN
```sql
-- IN: evaluates the full subquery, builds a list
SELECT * FROM employees
WHERE dept_id IN (SELECT id FROM departments WHERE name = 'Engineering');

-- EXISTS: stops as soon as a match is found (short-circuit)
SELECT * FROM employees e
WHERE EXISTS (
    SELECT 1 FROM departments d
    WHERE d.id = e.dept_id AND d.name = 'Engineering'
);
```

| | IN | EXISTS |
|--|-----|--------|
| Handles NULLs | Poorly — `IN (NULL, 1)` never matches | Correctly |
| Large subquery | Slower (builds full list) | Faster (short-circuits) |
| Correlated | Not naturally | Natural use case |
| Best for | Small, static list | Large subquery / existence check |

---

## Interview Questions

### Q1. What is the difference between INNER JOIN and LEFT JOIN?

> **Answer:**
> - **INNER JOIN:** Returns only rows where there is a **match in both tables**.
> - **LEFT JOIN:** Returns **all rows from the left table**, plus matched rows from the right. If no match, right-side columns are NULL.
>
> Use LEFT JOIN when you want to keep all records from one table regardless of whether a related record exists in the other (e.g., all employees even if they're not in a department).

---

### Q2. What is the difference between WHERE and HAVING?

> **Answer:**
> - **WHERE** filters **individual rows before** aggregation/grouping.
> - **HAVING** filters **groups after** GROUP BY aggregation.
>
> You cannot use aggregate functions (COUNT, SUM, AVG) in WHERE — use HAVING for that.
> ```sql
> SELECT dept_id, COUNT(*) FROM employees
> WHERE salary > 50000           -- filters rows first
> GROUP BY dept_id
> HAVING COUNT(*) > 3;           -- filters groups
> ```

---

### Q3. What is the difference between ROW_NUMBER, RANK, and DENSE_RANK?

> **Answer:**
> All rank rows within a partition. Difference is in how they handle **ties**:
> - `ROW_NUMBER()` — always unique (1,2,3,4) — arbitrary tie-breaking.
> - `RANK()` — same rank for ties, **skips** next ranks (1,2,2,4).
> - `DENSE_RANK()` — same rank for ties, **no skip** (1,2,2,3).
>
> Use `DENSE_RANK()` for "top-N" queries where you want no gaps in ranking.

---

### Q4. What is a correlated subquery? How does it differ from a regular subquery?

> **Answer:**
> - **Regular subquery:** Executes once, independently of the outer query. Result is reused.
> - **Correlated subquery:** References columns from the outer query, so it **re-executes for every row** of the outer query.
>
> Correlated subqueries are often slow on large datasets. Replace with a JOIN or CTE + JOIN for better performance.

---

### Q5. What is the N+1 query problem? How do you fix it?

> **Answer:**
> N+1 occurs when fetching a list of N records and then issuing **one additional query per record** to fetch related data. Total = N+1 queries.
>
> **Fix:** Replace with a single JOIN query, or use batch/eager loading:
> ```sql
> -- ❌ N+1
> SELECT * FROM posts;  -- then for each post: SELECT * FROM users WHERE id = ?
>
> -- ✅ Single query
> SELECT p.*, u.name FROM posts p JOIN users u ON p.author_id = u.id;
> ```

---

### Q6. When would you use EXISTS instead of IN?

> **Answer:**
> - **EXISTS** is better when the subquery can be **large** — it short-circuits as soon as a match is found.
> - **IN** can behave unexpectedly with **NULLs** — `WHERE id NOT IN (1, NULL)` returns no rows because NULL comparisons are unknown.
> - For **correlated existence checks**, EXISTS is the natural choice and often more readable.
>
> Use `IN` for small, known, static lists. Use `EXISTS` for subquery-based existence tests.

---

<div align="center">
  <sub>← Back to <a href="Topic.md">All Topics</a></sub>
</div>
