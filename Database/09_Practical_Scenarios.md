# 🟢 Practical / Scenario-based

> **Category:** Queries &nbsp;|&nbsp; **Tags:** `URL shortener` `social feed` `Nth salary` `pagination`

---

## Table of Contents
1. [Design a URL Shortener Schema](#1-design-a-url-shortener-schema)
2. [Design a Social Media Feed Schema](#2-design-a-social-media-feed-schema)
3. [Find the Nth Highest Salary](#3-find-the-nth-highest-salary)
4. [Detect Duplicate Records](#4-detect-duplicate-records)
5. [Pagination Strategies](#5-pagination-strategies)
6. [Bonus: More Classic SQL Puzzles](#6-bonus-more-classic-sql-puzzles)

---

## 1. Design a URL Shortener Schema

### Requirements
- Store original URLs and their short codes
- Track click counts per link
- Track per-click analytics (IP, user agent, timestamp)
- Support optional expiry
- Support user ownership

### Schema

```sql
-- Users (optional — for user-owned links)
CREATE TABLE users (
    id         BIGSERIAL PRIMARY KEY,
    email      VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Short links
CREATE TABLE short_links (
    id           BIGSERIAL PRIMARY KEY,
    short_code   VARCHAR(12)  UNIQUE NOT NULL,    -- e.g., "abc123"
    original_url TEXT         NOT NULL,
    user_id      BIGINT       REFERENCES users(id) ON DELETE SET NULL,
    click_count  BIGINT       DEFAULT 0,
    expires_at   TIMESTAMP,                        -- NULL = never expires
    created_at   TIMESTAMP    DEFAULT NOW(),
    is_active    BOOLEAN      DEFAULT TRUE
);

CREATE INDEX idx_short_links_code   ON short_links(short_code);   -- primary lookup
CREATE INDEX idx_short_links_user   ON short_links(user_id);
CREATE INDEX idx_short_links_expiry ON short_links(expires_at)
    WHERE expires_at IS NOT NULL;                                  -- partial index

-- Click analytics (high-volume inserts)
CREATE TABLE link_clicks (
    id           BIGSERIAL PRIMARY KEY,
    short_link_id BIGINT     NOT NULL REFERENCES short_links(id) ON DELETE CASCADE,
    clicked_at   TIMESTAMP   DEFAULT NOW(),
    ip_address   INET,
    user_agent   TEXT,
    referrer     TEXT,
    country_code CHAR(2)
);

-- Partition by month for performance at scale
-- In PostgreSQL: PARTITION BY RANGE (clicked_at)

CREATE INDEX idx_clicks_link_id ON link_clicks(short_link_id, clicked_at DESC);
```

### Query: Redirect flow
```sql
-- Lookup + validate (atomic with update)
UPDATE short_links
SET click_count = click_count + 1
WHERE short_code = 'abc123'
  AND is_active = TRUE
  AND (expires_at IS NULL OR expires_at > NOW())
RETURNING original_url;
```

### Query: Top links for a user
```sql
SELECT short_code, original_url, click_count
FROM short_links
WHERE user_id = 42
ORDER BY click_count DESC
LIMIT 10;
```

### Key design decisions
- `short_code` indexed for fast O(log n) lookup on every redirect.
- `click_count` is a counter on the main table for fast display — eventual consistency acceptable.
- `link_clicks` is append-only (no updates) — good for partitioning by time.
- Use a base62 encoder (`[a-zA-Z0-9]`) on the auto-increment ID to generate short codes.

---

## 2. Design a Social Media Feed Schema

### Requirements
- Users can post (text, images)
- Users can follow other users
- Users can like and comment on posts
- Feed shows posts from followed users, sorted by recency

### Schema

```sql
CREATE TABLE users (
    id         BIGSERIAL PRIMARY KEY,
    username   VARCHAR(50)  UNIQUE NOT NULL,
    bio        TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE posts (
    id           BIGSERIAL PRIMARY KEY,
    user_id      BIGINT    NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content      TEXT,
    media_url    TEXT,                     -- image/video link
    like_count   INT DEFAULT 0,           -- denormalized counter
    comment_count INT DEFAULT 0,
    created_at   TIMESTAMP DEFAULT NOW(),
    deleted_at   TIMESTAMP                -- soft delete
);

CREATE INDEX idx_posts_user_time ON posts(user_id, created_at DESC)
    WHERE deleted_at IS NULL;

-- Follow relationships
CREATE TABLE follows (
    follower_id  BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    followee_id  BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at   TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (follower_id, followee_id)
);

CREATE INDEX idx_follows_followee ON follows(followee_id);  -- "who follows me"
CREATE INDEX idx_follows_follower ON follows(follower_id);  -- "who do I follow"

-- Likes
CREATE TABLE likes (
    user_id    BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    post_id    BIGINT NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (user_id, post_id)
);

-- Comments
CREATE TABLE comments (
    id         BIGSERIAL PRIMARY KEY,
    post_id    BIGINT NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    user_id    BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content    TEXT   NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX idx_comments_post ON comments(post_id, created_at DESC)
    WHERE deleted_at IS NULL;
```

### Feed Query (Pull model)
```sql
-- Get feed for user_id = 42 (posts from followed users)
SELECT p.id, p.content, p.like_count, p.created_at,
       u.username, u.avatar_url
FROM posts p
JOIN follows f ON f.followee_id = p.user_id
JOIN users u   ON u.id = p.user_id
WHERE f.follower_id = 42
  AND p.deleted_at IS NULL
ORDER BY p.created_at DESC
LIMIT 20
OFFSET 0;
```

### Key design decisions

| Decision | Reasoning |
|----------|----------|
| `like_count` on posts | Denormalized counter — avoid `COUNT(*)` from likes table on every feed render |
| Soft delete on posts | Preserve data; consistent referential integrity |
| Composite index on `(user_id, created_at DESC)` | Feed queries always filter by user and sort by time |
| Pull model vs push model | Pull: query at read time (simpler, slower at scale). Push: fan-out writes to each follower's feed table (faster reads, complex writes for celebrities with 10M followers) |

---

## 3. Find the Nth Highest Salary

Classic SQL puzzle — multiple solutions.

### Setup
```sql
CREATE TABLE employees (
    id     INT,
    name   VARCHAR(100),
    salary DECIMAL(10,2)
);

INSERT INTO employees VALUES
(1, 'Alice', 90000),
(2, 'Bob',   75000),
(3, 'Carol', 90000),  -- same as Alice
(4, 'Dave',  60000),
(5, 'Eve',   75000);
```

### Solution 1 — DENSE_RANK (recommended)
```sql
-- Find the 2nd highest unique salary
SELECT salary FROM (
    SELECT salary,
           DENSE_RANK() OVER (ORDER BY salary DESC) AS rnk
    FROM employees
) ranked
WHERE rnk = 2;
-- Result: 75000
```

### Solution 2 — Subquery with DISTINCT
```sql
-- Nth highest: N = 2
SELECT MIN(salary) FROM (
    SELECT DISTINCT salary FROM employees
    ORDER BY salary DESC
    LIMIT 2    -- top 2 distinct salaries
) top_n;
-- Result: 75000
```

### Solution 3 — Correlated subquery
```sql
-- Find salary that has exactly N-1 distinct salaries above it
SELECT DISTINCT salary FROM employees e1
WHERE 1 = (          -- N-1 = 1
    SELECT COUNT(DISTINCT salary) FROM employees e2
    WHERE e2.salary > e1.salary
);
-- Result: 75000 (1 distinct salary above it: 90000)
```

### Handling NULLs / no result
```sql
-- Return NULL if fewer than N distinct salaries exist
SELECT COALESCE(
    (SELECT DISTINCT salary FROM employees
     ORDER BY salary DESC
     LIMIT 1 OFFSET 1),   -- OFFSET N-1
    NULL
) AS nth_salary;
```

---

## 4. Detect Duplicate Records

### Find duplicate emails
```sql
-- Show all emails that appear more than once
SELECT email, COUNT(*) AS occurrences
FROM users
GROUP BY email
HAVING COUNT(*) > 1;
```

### Find the duplicate rows with all columns
```sql
-- Show all rows that have duplicate (first_name, last_name, email)
SELECT * FROM users
WHERE (first_name, last_name, email) IN (
    SELECT first_name, last_name, email
    FROM users
    GROUP BY first_name, last_name, email
    HAVING COUNT(*) > 1
)
ORDER BY email, id;
```

### Delete duplicates — keep the row with the lowest ID
```sql
-- Using CTE + ROW_NUMBER
WITH duplicates AS (
    SELECT id,
           ROW_NUMBER() OVER (PARTITION BY email ORDER BY id ASC) AS rn
    FROM users
)
DELETE FROM users
WHERE id IN (
    SELECT id FROM duplicates WHERE rn > 1
);
-- Keeps the first (lowest id) occurrence; deletes the rest
```

---

## 5. Pagination Strategies

### OFFSET-based Pagination (Simple, but flawed)
```sql
-- Page 1 (items 1-20)
SELECT * FROM posts ORDER BY created_at DESC LIMIT 20 OFFSET 0;

-- Page 2 (items 21-40)
SELECT * FROM posts ORDER BY created_at DESC LIMIT 20 OFFSET 20;

-- Page N
SELECT * FROM posts ORDER BY created_at DESC LIMIT 20 OFFSET (N-1)*20;
```

**Problems with OFFSET:**
- **Performance degrades** with large offsets — DB must scan and discard `OFFSET` rows.
- **Drift:** If rows are inserted/deleted between page requests, items can be skipped or duplicated.
- `OFFSET 1000000` on a 10M-row table = scan 1M rows and discard.

---

### Cursor-based Pagination (Keyset Pagination) — Recommended
```sql
-- First page
SELECT id, title, created_at FROM posts
WHERE deleted_at IS NULL
ORDER BY created_at DESC, id DESC
LIMIT 20;
-- Returns last item: created_at = '2024-05-01 12:00:00', id = 5001

-- Next page: pass the last item's cursor values
SELECT id, title, created_at FROM posts
WHERE deleted_at IS NULL
  AND (created_at, id) < ('2024-05-01 12:00:00', 5001)   -- keyset condition
ORDER BY created_at DESC, id DESC
LIMIT 20;
```

**Advantages:**
- O(log n) — uses the index directly, no skipping rows.
- Stable — inserting new rows doesn't shift results.
- Works well for infinite scroll.

**Disadvantages:**
- Can't jump to arbitrary page numbers.
- Cursor must be included in client requests.
- Sorting must use a unique or tie-broken column.

### Comparison

| Feature | OFFSET | Cursor-based |
|---------|--------|-------------|
| Performance at depth | ❌ O(n) degrades | ✅ O(log n) always |
| Stable (no drift) | ❌ | ✅ |
| Random page access | ✅ | ❌ |
| Total page count | ✅ | ❌ (expensive) |
| Implementation | Simple | Slightly complex |
| Best for | Admin UIs, small tables | APIs, infinite scroll, large tables |

---

## 6. Bonus: More Classic SQL Puzzles

### Employees earning more than their manager
```sql
SELECT e.name AS employee, e.salary,
       m.name AS manager,  m.salary AS manager_salary
FROM employees e
JOIN employees m ON e.manager_id = m.id
WHERE e.salary > m.salary;
```

### Departments with no employees
```sql
SELECT d.name FROM departments d
LEFT JOIN employees e ON e.dept_id = d.id
WHERE e.id IS NULL;
```

### Running total
```sql
SELECT id, amount, created_at,
       SUM(amount) OVER (ORDER BY created_at) AS running_total
FROM transactions;
```

### Month-over-month revenue growth
```sql
WITH monthly AS (
    SELECT DATE_TRUNC('month', created_at) AS month,
           SUM(total) AS revenue
    FROM orders
    GROUP BY 1
)
SELECT month, revenue,
       LAG(revenue) OVER (ORDER BY month) AS prev_month_revenue,
       ROUND((revenue - LAG(revenue) OVER (ORDER BY month))
             / LAG(revenue) OVER (ORDER BY month) * 100, 2) AS pct_growth
FROM monthly
ORDER BY month;
```

---

<div align="center">
  <sub>← Back to <a href="Topic.md">All Topics</a></sub>
</div>
