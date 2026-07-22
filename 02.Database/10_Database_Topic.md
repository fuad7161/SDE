# Database Interview Guide — Explanations & Answers

---

## 1. Database Fundamentals

### Q1. What is ACID?
**Explanation:** ACID is the set of guarantees a database transaction provides to keep data reliable even during crashes, concurrent access, or errors.

**Answer:**
- **Atomicity** — a transaction is all-or-nothing. If one part fails, the whole transaction rolls back.
- **Consistency** — a transaction moves the database from one valid state to another, respecting constraints, triggers, and rules.
- **Isolation** — concurrent transactions don't see each other's uncommitted changes.
- **Durability** — once committed, data survives crashes (usually via WAL/redo logs).

**Follow-ups:**
- *Why do we need ACID?* Without it, concurrent writes could corrupt data, partial failures could leave inconsistent state, and crashes could lose committed work.
- *Which databases fully support ACID?* PostgreSQL, MySQL (InnoDB), Oracle, SQL Server, SQLite.
- *Does MongoDB support ACID?* Yes, since v4.0, MongoDB supports multi-document ACID transactions, though single-document operations were already atomic before that.

---

### Q2. Explain database normalization.
**Explanation:** Normalization organizes tables to reduce redundancy and avoid update anomalies by progressively applying stricter rules.

**Answer:**
- **1NF** — atomic columns, no repeating groups.
- **2NF** — 1NF + no partial dependency (non-key columns depend on the whole composite key).
- **3NF** — 2NF + no transitive dependency (non-key columns depend only on the key, not on other non-key columns).
- **BCNF** — stricter 3NF: every determinant must be a candidate key.

**Follow-ups:**
- *Why normalize?* Reduces duplicate data, prevents update/insert/delete anomalies, keeps data consistent.
- *When should you denormalize?* When read performance matters more than write simplicity — reporting, analytics, read-heavy APIs.

---

### Q3. Normalization vs Denormalization
**Explanation:** A trade-off between data integrity/storage efficiency and query speed.

**Answer:**
- **Normalization advantages:** less redundancy, easier and safer updates (single source of truth).
- **Normalization disadvantage:** more joins needed to reconstruct data, which can hurt read performance.
- **When denormalization is better:** reporting, analytics dashboards, read-heavy systems where join cost outweighs the risk of redundant data.

---

### Q4. Primary Key vs Unique Key
**Explanation:** Both enforce uniqueness but serve different roles.

**Answer:**
| Aspect | Primary Key | Unique Key |
|---|---|---|
| NULLs | Not allowed | One (or more, depending on DB) NULLs allowed |
| Count per table | Only one | Multiple allowed |
| Index | Automatically creates a clustered index (in most DBs) | Creates a non-clustered/unique index |

---

### Q5. Clustered vs Non-clustered Index
**Explanation:** Determines how data is physically stored vs. how it's looked up.

**Answer:**
- **Clustered index:** the table's rows are physically stored in the order of the index key. There can be only one per table. Typically built on the primary key.
- **Non-clustered index:** a separate B+ Tree structure whose leaf nodes store pointers back to the actual row, not the row itself. A table can have many.
- **Performance:** clustered index range scans avoid an extra pointer hop; non-clustered indexes need an extra lookup ("bookmark lookup") to fetch the full row unless it's a covering index.

---

## 2. Indexing (Very Frequently Asked)

### Q6. What is an Index?
**Explanation:** A data structure (usually a B+ Tree) that speeds up lookups at the cost of extra storage and slower writes.

**Answer:** An index maps column values to row locations so the database can find rows without scanning the whole table (O(log n) instead of O(n)). Trade-offs: extra disk space, and every INSERT/UPDATE/DELETE must also update the index, slowing writes.

---

### Q7. How does a B+ Tree work?
**Explanation:** Senior interviews almost always ask this.

**Answer:**
- **Internal nodes** store only keys used to navigate to children (no actual row data).
- **Leaf nodes** store the actual keys plus pointers to rows, and are linked together for fast range scans.
- The tree is **balanced** — all leaves are at the same depth.
- Lookup, insert, and delete are all **O(log n)** because high fan-out keeps the tree height low even as data grows.

---

### Q8. Why don't databases use Binary Search Trees?
**Explanation:** BSTs are great in memory but poor for disk-based storage.

**Answer:**
- A plain BST can become **unbalanced**, degrading to O(n).
- Each BST node holds only one key, so the tree gets tall, requiring many **disk I/O** operations (each node = a potential disk seek).
- B+ Trees have very high **fan-out** (hundreds of keys per node), keeping **tree height** low (often 3–4 levels for millions of rows), minimizing disk I/O.

---

### Q9. Hash Index vs B-tree Index
**Explanation:** Need comparison — equality search vs range search.

**Answer:**
- **Hash Index:** O(1) average lookup, but only supports equality search (`=`). No range queries or sorting.
- **B-tree Index:** O(log n) lookup, but supports equality *and* range search (`<`, `>`, `BETWEEN`), ordering, and prefix matching — why it's the default in most databases.

---

### Q10. Composite Index
**Explanation:** An index on multiple columns, ordered left to right.

**Example:**
```sql
CREATE INDEX idx_user_age ON users(country, age);
```

**Answer:**
- `WHERE country='BD'` → **YES**, uses the index (leftmost column).
- `WHERE age=25` → **NO**, can't use the index because `age` isn't the leftmost column.

**Leftmost Prefix Rule:** A composite index on `(a, b, c)` can be used for queries filtering on `a`, `(a,b)`, or `(a,b,c)` — but not on `b`, `c`, or `(b,c)` alone.

---

### Q11. Covering Index
**Explanation:** What is it, and why does it improve performance?

**Answer:** A covering index contains all the columns a query needs, so the database can answer the query entirely from the index without touching the actual table (no "bookmark lookup"). This avoids extra disk reads to fetch full rows.

---

### Q12. Why are too many indexes bad?
**Explanation:** Indexes aren't free — they're extra structures the DB must maintain.

**Answer:**
- **Insert slower** — every insert must also insert into each index.
- **Update slower** — updating an indexed column means updating the index too.
- **Delete slower** — index entries must also be removed.
- **More storage** — each index consumes additional disk space.

---

## 3. SQL Questions

### Q13. Difference between WHERE and HAVING
**Answer:** `WHERE` filters rows before grouping/aggregation; `HAVING` filters groups after `GROUP BY`. Aggregate functions can't be used in `WHERE`, but can in `HAVING`.

---

### Q14. INNER JOIN vs LEFT JOIN vs RIGHT JOIN vs FULL JOIN
**Answer:**
- **INNER JOIN** — only matching rows in both tables.
- **LEFT JOIN** — all rows from the left table, matched rows from the right (NULLs if no match).
- **RIGHT JOIN** — all rows from the right table, matched rows from the left.
- **FULL JOIN** — all rows from both tables, NULLs where there's no match on either side.

---

### Q15. Explain different JOIN algorithms.
**Explanation:** Senior interviews.

**Answer:**
- **Nested Loop Join:** for each row in the outer table, scan the inner table for matches. Good for small tables or an indexed inner side.
- **Hash Join:** builds a hash table on the smaller input's join key, then probes with the other input. Good for large, unsorted equality joins.
- **Merge Join:** both inputs are sorted on the join key, then merged like a zipper. Efficient when inputs are already sorted (e.g., from an index).

---

### Q16. GROUP BY
**Explanation:** How does it work internally?

**Answer:** The database groups rows sharing the same value(s) via **hash-based** grouping (build a hash table keyed by group values) or **sort-based** grouping (sort rows by group key, then aggregate consecutive rows). Aggregate functions are then computed per group.

---

### Q17. DISTINCT vs GROUP BY
**Answer:** `DISTINCT` removes duplicate rows from the result set. `GROUP BY` groups rows to apply aggregate functions per group. They can produce similar results for simple cases, but `GROUP BY` is meant for aggregation.

---

### Q18. UNION vs UNION ALL
**Answer:** `UNION` combines result sets and removes duplicates (requires a sort/hash dedup step, slower). `UNION ALL` keeps duplicates and is faster since no dedup is needed.

---

### Q19. EXISTS vs IN
**Explanation:** Performance discussion.

**Answer:** `EXISTS` stops as soon as it finds one matching row (short-circuits), often faster for correlated subqueries with large inner result sets. `IN` evaluates the full list/subquery first. `EXISTS` is also safer with NULLs, since `NOT IN` behaves unexpectedly when the list contains NULL.

---

### Q20. DELETE vs TRUNCATE vs DROP
**Answer:**
- **DELETE** — removes rows (optionally with `WHERE`), logged row-by-row, can be rolled back, triggers fire, keeps table structure.
- **TRUNCATE** — removes all rows quickly, minimally logged, resets auto-increment, keeps table structure.
- **DROP** — removes the entire table (structure + data).

---

## 4. Transactions

### Q21. What is a transaction?
**Answer:** A sequence of one or more SQL operations executed as a single logical unit of work — either all changes are applied, or none are, guaranteeing ACID properties.

---

### Q22. Transaction Lifecycle
**Answer:**
1. **BEGIN** — start the transaction.
2. **SQL** — execute statements.
3. **COMMIT** — persist all changes permanently.
4. **ROLLBACK** — undo all changes since BEGIN (used on error).

---

### Q23. What is Isolation Level?
**Explanation:** Need all four.

**Answer:**
- **Read Uncommitted** — can see uncommitted changes from other transactions (dirty reads possible).
- **Read Committed** — only sees committed data; each query sees a fresh snapshot (Postgres default).
- **Repeatable Read** — the whole transaction sees a consistent snapshot from its start.
- **Serializable** — strictest; transactions behave as if executed one after another.

---

### Q24. Explain Dirty Read
**Answer:** A transaction reads data another transaction wrote but hasn't committed. If that transaction rolls back, the first transaction acted on data that never really existed.

---

### Q25. Explain Non-repeatable Read
**Answer:** A transaction reads the same row twice and gets different values because another transaction updated and committed a change in between.

---

### Q26. Explain Phantom Read
**Answer:** A transaction re-runs the same query twice and gets a different set of rows because another transaction inserted or deleted matching rows in between.

---

### Q27. Which isolation level does PostgreSQL use?
**Answer:** By default, **Read Committed**. PostgreSQL also supports Repeatable Read and Serializable, implemented via MVCC snapshots.

---

## 5. Locks

### Q28. Optimistic Locking
**Answer:** Assumes conflicts are rare. No locks held while reading; a version number is checked at update time. If it changed since it was read, the update is rejected and retried. Good for low-contention workloads.

---

### Q29. Pessimistic Locking
**Answer:** Assumes conflicts are likely. Locks the row (`SELECT ... FOR UPDATE`) as soon as read, blocking others until released. Good for high-contention workloads, reduces concurrency.

---

### Q30. Row Lock vs Table Lock
**Answer:** A **row lock** locks only the specific row(s), allowing high concurrency elsewhere. A **table lock** locks the entire table, blocking all other transactions — more restrictive but sometimes necessary for schema changes or bulk operations.

---

### Q31. What causes deadlock?
**Explanation:** Very common. Need example.

**Answer:** Two (or more) transactions each hold a lock the other needs, forming a circular wait.

**Example:** Transaction A locks row 1, then wants row 2. Transaction B locks row 2, then wants row 1. Neither can proceed — the database detects the cycle and kills one transaction to break it.

---

### Q32. How do you avoid deadlocks?
**Answer:**
- Acquire locks in a **consistent order** across transactions.
- Keep transactions **short**.
- Use lower isolation levels when possible.
- Add appropriate **indexes** to reduce lock scope.
- Use **timeouts** and retry logic for deadlock victims.

---

## 6. Query Optimization

### Q33. Why is my query slow?
**Answer:** Missing index on filtered/joined columns, too many joins, a sequential scan on a large table, an unnecessarily large result set, or a bad execution plan from stale statistics.

---

### Q34. How do you analyze a slow query?
**Answer:** In PostgreSQL, use `EXPLAIN` to see the planned strategy without running the query, or `EXPLAIN ANALYZE` to run it and show real timing/row counts — helping spot mismatches between estimated and actual rows.

---

### Q35. What is Query Execution Plan?
**Answer:** The step-by-step strategy the query planner chooses — which scan types, join algorithms, and order of operations — based on table statistics and available indexes.

---

### Q36. What is Sequential Scan?
**Answer:** The database reads every row in a table, checking each against the query condition. Efficient for small tables or when many rows match; slow for large tables with selective filters.

---

### Q37. Index Scan vs Bitmap Scan
**Explanation:** Senior-level topic.

**Answer:**
- **Index Scan:** walks the index and fetches each matching row one at a time — good when few rows match.
- **Bitmap Scan:** builds an in-memory bitmap of matching row locations from the index, then fetches rows in physical order in one pass — good when a moderate number of rows match, reducing random I/O.

---

## 7. PostgreSQL (Very Common for Go)

### Q38. MVCC
**Explanation:** Must know.

**Answer:** Instead of locking rows for reads, PostgreSQL keeps multiple versions of a row. Each transaction sees a consistent **snapshot** as of when it started, so readers never block writers and writers never block readers. Old row versions become dead tuples, cleaned up by VACUUM.

---

### Q39. VACUUM
**Answer:** Reclaims storage from dead tuples (old row versions left by MVCC after UPDATE/DELETE) and updates planner statistics. Without it, tables bloat and performance degrades.

---

### Q40. Autovacuum
**Answer:** A background process that automatically runs VACUUM (and ANALYZE) when tables cross configured thresholds of dead tuples, so DBAs don't need to run it manually.

---

### Q41. WAL (Write Ahead Log)
**Explanation:** Very common.

**Answer:** Before any data change is applied to data files, it's first written to the Write-Ahead Log. This guarantees durability — if the database crashes, it can replay the WAL to recover committed transactions not yet flushed to disk. It's also the foundation for replication.

---

### Q42. Checkpoint
**Answer:** A point where PostgreSQL flushes all dirty pages from memory to disk and marks that point in the WAL, limiting how much WAL needs replaying during crash recovery.

---

### Q43. HOT Update
**Answer:** Heap-Only Tuple update — if an updated row's indexed columns didn't change and there's space on the same page, PostgreSQL updates the row without touching any indexes, reducing write overhead and index bloat.

---

### Q44. TOAST
**Answer:** The Oversized-Attribute Storage Technique — PostgreSQL's mechanism for storing very large column values outside the main table row, in a separate TOAST table, compressed and/or chunked.

---

## 8. Replication

### Q45. Primary-Replica Replication
**Answer:** The primary handles all writes and streams changes (via WAL) to one or more replicas, which apply those changes to stay in sync and can serve read queries.

---

### Q46. Synchronous vs Asynchronous Replication
**Answer:**
- **Synchronous:** the primary waits for at least one replica to confirm before acknowledging commit. Safer, higher write latency.
- **Asynchronous:** the primary commits and acknowledges immediately without waiting. Faster, but a crash right after commit could lose recent transactions on failover.

---

### Q47. Replication Lag
**Answer:** The delay between a write on the primary and that change becoming visible on a replica. Caused by network latency, replica load, or long-running queries blocking WAL replay. High lag means stale reads from replicas.

---

### Q48. Read Replica
**Explanation:** When should you use it?

**Answer:** Use read replicas to scale read-heavy workloads, for geographically distributed reads, and for analytics/reporting without impacting the primary. Not suitable when the app needs strongly consistent, up-to-the-millisecond reads.

---

## 9. Scaling

### Q49. Vertical vs Horizontal Scaling
**Answer:**
- **Vertical:** add more resources to a single server. Simple, but has a hardware ceiling and causes downtime during upgrades.
- **Horizontal:** add more servers/nodes and distribute the load. More complex, but scales further with better fault tolerance.

---

### Q50. Sharding
**Explanation:** Splitting data across multiple independent database instances.

**Answer:**
- **Pros:** near-linear write scalability, smaller datasets per node, fault isolation.
- **Cons:** cross-shard joins/transactions are hard, added operational complexity, difficult rebalancing.
- **Shard key:** the column deciding which shard a row belongs to (e.g., `user_id`). A poor shard key causes hot shards.

---

### Q51. Partitioning
**Answer:** **Partitioning** splits a large table into smaller pieces within the same database instance, usually by range or list, for performance and manageability. **Sharding** splits data across multiple separate database instances/servers for horizontal write scalability. Partitioning is a single-node optimization; sharding is a distributed strategy (often combined).

---

## 10. Database Design

### Q52. Design a URL Shortener Database
**Answer:** Core table: `urls(id, short_code UNIQUE, long_url, user_id, created_at, expires_at, click_count)`. Generate `short_code` via base62 encoding of an auto-increment ID or a hash, with a unique index on `short_code`. Optionally a `clicks` table or counter cache for analytics, with a cache (Redis) in front for hot redirects.

---

### Q53. Design an E-commerce Database
**Answer:** Core tables:
- **Users**(id, name, email, password_hash, address, created_at)
- **Products**(id, name, description, price, category_id)
- **Inventory**(product_id, warehouse_id, quantity)
- **Orders**(id, user_id, status, total, created_at)
- **Order_Items**(order_id, product_id, quantity, price_at_purchase)

Orders → Users (many-to-one); Order_Items → Orders/Products (many-to-many via join table); Inventory tracked separately for multi-warehouse support.

---

### Q54. Design Chat Application Database
**Answer:** Core tables: `users`, `conversations`, `conversation_participants` (many-to-many), `messages(id, conversation_id, sender_id, content, created_at)`. Index `messages(conversation_id, created_at)` for fast pagination; consider partitioning/archiving old messages. Read receipts via a separate `message_reads(message_id, user_id, read_at)` table.

---

### Q55. How would you model many-to-many relationships?
**Answer:** Use a **join table** with foreign keys to both related tables, e.g., `students`, `courses`, `student_courses(student_id, course_id)`. The join table can also hold extra attributes about the relationship (e.g., `enrolled_at`, `grade`).

---

### Q56. Soft Delete vs Hard Delete
**Answer:**
- **Soft delete:** mark a row as deleted (`deleted_at`/`is_deleted`) without removing it. Preserves history, allows recovery, but requires filtering in every query and can bloat tables.
- **Hard delete:** physically removes the row. Frees storage immediately, simpler queries, but unrecoverable and breaks audit trails.

---

## 11. Concurrency

### Q57. Two users update same row. How do you prevent data loss?
**Explanation:** Classic "lost update" problem.

**Answer:** Use **optimistic locking** with a **version column**: `WHERE id=? AND version=?`, incrementing version on update. If zero rows are affected, someone else updated first — retry or notify the user. Alternatively use pessimistic locking (`SELECT ... FOR UPDATE`) for frequent conflicts.

---

### Q58. Bank Transfer Problem
**Explanation:** Transfer $100 from Account A to Account B — how to ensure consistency?

**Answer:** Wrap both operations in a single transaction:
```sql
BEGIN;
UPDATE accounts SET balance = balance - 100 WHERE id = 'A';
UPDATE accounts SET balance = balance + 100 WHERE id = 'B';
COMMIT;
```
This guarantees atomicity (both updates happen or neither does). Locking the rows prevents concurrent transactions from reading a stale balance. Always update accounts in a **consistent order** across all transfers to avoid deadlocks.

---

## 12. Golang + Database

### Q59. Why use connection pooling?
**Answer:** Opening a new connection is expensive (TCP handshake, auth, resource allocation). A pool keeps ready-to-use connections that the app reuses, reducing latency and resource overhead, and prevents overwhelming the database with too many connections.

---

### Q60. How does database/sql manage connections?
**Answer:** Go's `database/sql` maintains an internal pool. It lazily opens connections as needed (up to `SetMaxOpenConns`), reuses idle ones (up to `SetMaxIdleConns`), and recycles connections after `SetConnMaxLifetime`. Each query borrows a connection and returns it when done.

---

### Q61. What happens if rows aren't closed?
**Explanation:** Need `defer rows.Close()`.

**Answer:** The underlying connection is never returned to the pool — it stays checked out until closed or garbage collected. Under load this exhausts the pool, causing new queries to block or time out. Always `defer rows.Close()` right after checking the query error.

---

### Q62. Context timeout in SQL queries
**Explanation:** `ctx, cancel := context.WithTimeout(...)`.

**Answer:** Passing a context with a timeout (e.g., `db.QueryContext(ctx, ...)`) ensures a query is automatically cancelled if it takes too long, preventing a slow query from hanging a request indefinitely, holding a connection, and cascading into resource exhaustion.

---

### Q63. SQL Injection
**Explanation:** How to prevent it.

**Answer:** Occurs when user input is concatenated directly into a SQL string, letting an attacker inject arbitrary SQL. Prevention: always use **parameterized queries** (placeholders like `$1` or `?`) so the driver sends query and data separately — never build SQL via string concatenation with user input.

---

### Q64. Prepared Statement
**Explanation:** Benefits.

**Answer:** Parsed and planned by the database once, then executed multiple times with different parameters. Benefits: prevents SQL injection, reduces re-parsing/re-planning overhead, slightly improves performance for repeated queries.

---

### Q65. ORM vs Raw SQL
**Explanation:** Trade-offs.

**Answer:**
- **ORM:** faster development, less boilerplate, database-agnostic, built-in injection protection — but can generate inefficient queries (N+1 problem) and abstracts away performance details.
- **Raw SQL:** full control over performance and exact behavior — but more boilerplate, manual result mapping, and higher injection risk if not parameterized carefully.

---

## 13. Redis + Database

### Q66. Why cache database queries?
**Answer:** Databases are slower (disk I/O, query planning) than an in-memory cache. Caching frequently-read, rarely-changing data reduces database load, lowers latency, and increases read throughput.

---

### Q67. Cache Aside Pattern
**Answer:** The app checks the cache first. On a miss, it reads from the database, then writes the result into the cache. On writes, the app typically invalidates (deletes) the cache entry so the next read repopulates it fresh. The most common caching pattern.

---

### Q68. Write Through Cache
**Answer:** Every write goes to the cache and database simultaneously, keeping them in sync. Reads are always fast since data is guaranteed to be in cache after a write, but writes are slightly slower.

---

### Q69. Cache Invalidation
**Explanation:** Classic interview topic.

**Answer:** Strategies: **TTL-based expiration** (simple, can serve stale data), **explicit invalidation** (delete/update the key when underlying data changes), and **write-through** (always current). The right choice depends on staleness tolerance vs. extra invalidation logic.

---

### Q70. Prevent Cache Stampede
**Explanation:** Also known as "thundering herd."

**Answer:** When a popular key expires, many concurrent requests can simultaneously hit the database. Prevent with a **lock/mutex** so only one request repopulates the cache, **request coalescing**, **early/proactive refresh** before actual expiry, or **randomized TTL jitter** so many keys don't expire at once.

---

## 14. Scenario-Based Questions

### Q71. Your API suddenly became slow. How do you investigate?
**Answer:** Check dashboards (CPU, memory, DB connections, latency percentiles) to isolate app vs. network vs. database. Check for slow queries (`EXPLAIN ANALYZE`, slow-query log), connection pool exhaustion, recent deploys/config changes, and blocking/locking queries.

---

### Q72. Database CPU reaches 100%. What do you do?
**Answer:** Identify the top CPU-consuming queries (`pg_stat_activity`/`pg_stat_statements` in Postgres). Look for missing indexes causing sequential scans, runaway queries, or a traffic spike. Kill/optimize offending queries, add indexes, consider read replicas or caching.

---

### Q73. A query takes 10 seconds. How do you optimize it?
**Answer:** Run `EXPLAIN ANALYZE`. Look for sequential scans on large tables (add index), poor join order/algorithm, missing statistics (run `ANALYZE`), or an overly broad `SELECT *`. Rewrite, add a covering index, or denormalize for a recurring hot path.

---

### Q74. Millions of inserts per second. Database design?
**Answer:** Batch inserts instead of single-row inserts, partition/shard the table to spread write load, minimize indexes on the hot table, consider an append-only/log-structured engine, use asynchronous replication, and offload to a queue (Kafka) with consumers batching writes.

---

### Q75. Read-heavy system. How would you scale?
**Answer:** Add read replicas, introduce a caching layer (Redis), add appropriate indexes, and consider denormalization for expensive joins. A CDN can help for cacheable content.

---

### Q76. Write-heavy system. How would you scale?
**Answer:** Use sharding to parallelize writes, batch writes where possible, use a message queue to buffer bursts, minimize synchronous secondary work (indexes, triggers) on the write path, and consider write-optimized storage engines.

---

### Q77. User reports missing data after deployment. How do you debug?
**Answer:** Check recent migration scripts for destructive changes (dropped columns, bad WHERE clauses in UPDATE/DELETE), check application logs around the deployment time, check if it's a caching issue rather than real data loss, and check backups/WAL to determine if and when the data was removed.

---

### Q78. Duplicate orders are being created. How would you fix it?
**Answer:** Add a **unique constraint** on an idempotency key (e.g., `idempotency_key` or `(user_id, cart_id)`) so the database rejects duplicates. Use a client-provided idempotency key for retried requests, check-before-insert within a transaction, or `INSERT ... ON CONFLICT DO NOTHING`.

---

### Q79. How would you migrate a database with zero downtime?
**Answer:** Use the **expand-contract pattern**: add new columns/tables without removing old ones (expand); deploy code writing to both schemas; backfill historical data; switch reads to the new schema; once verified, remove the old schema (contract). Use online schema-change tools (e.g., `pg_repack`, `gh-ost`) for large tables.

---

### Q80. Your production database is full. What do you do?
**Answer:** Check disk usage and free emergency space (rotate/delete old logs, temp files). Identify what's consuming space — bloated tables/indexes (run VACUUM), unbounded log tables, or missing retention policies. Longer-term: monitoring/alerts on disk usage, data archiving, and scaling storage or partitioning old data out.

---

## 15. Advanced Topics (Senior)

### Q81. CAP Theorem
**Answer:** In a distributed system, you can only guarantee two of three during a network partition:
- **Consistency** — every read gets the most recent write.
- **Availability** — every request gets a response.
- **Partition Tolerance** — the system keeps working despite network partitions.

Since partitions are unavoidable, the real trade-off is between Consistency (CP) and Availability (AP).

---

### Q82. Eventual Consistency
**Answer:** If no new writes occur, all replicas will eventually converge to the same value — but reads immediately after a write might return stale data. Common in AP systems (DynamoDB, Cassandra) that prioritize availability.

---

### Q83. CQRS
**Answer:** Command Query Responsibility Segregation — separates the write model (commands) from the read model (queries), often using different data stores optimized for each. Useful for complex domains and high read/write asymmetry, but adds complexity (syncing the two models).

---

### Q84. Database Failover
**Answer:** Automatically (or manually) promoting a replica to primary when the current primary fails, minimizing downtime. Involves failure detection (health checks), electing/promoting a replica, and redirecting traffic — often via tools like Patroni, Consul, or a cloud provider's managed failover.

---

### Q85. Leader Election
**Answer:** How distributed nodes agree on which single node acts as leader, especially after a failure. Typically implemented via consensus algorithms like **Raft** or **Paxos**, or coordination services like **etcd**/**ZooKeeper**.

---

### Q86. Idempotency in Database Operations
**Answer:** An operation is idempotent if performing it multiple times has the same effect as once. Important for safe retries after network timeouts. Achieved via idempotency keys, `UPSERT`/`INSERT ... ON CONFLICT`, or checking existing state before writing.

---

### Q87. Outbox Pattern
**Explanation:** Very common in microservices.

**Answer:** To reliably publish events after a DB write (avoiding the "dual write" problem), the app writes business data and an event record into an `outbox` table in the **same transaction**. A background process (or CDC tool like Debezium) reads the outbox and publishes events to the message broker, guaranteeing at-least-once delivery consistent with the DB state.

---

### Q88. Saga Pattern
**Answer:** Manages distributed transactions across services without a global two-phase commit. A saga is a sequence of local transactions; if a step fails, **compensating transactions** undo previous steps. Implemented as **choreography** (services react to each other's events) or **orchestration** (a central coordinator directs the steps).

---

### Q89. Why PostgreSQL over MySQL?
**Answer:** Richer feature set (advanced indexing like GiST/GIN, native JSON/JSONB with indexing, full-text search, arrays, window functions, CTEs), stricter standards compliance, robust MVCC, extensibility (custom types, PostGIS). MySQL can still be preferred for simpler read-heavy workloads and wide hosting support — the right choice depends on the use case.

---

### Q90. Explain MVCC Internally
**Answer:** Every row version has hidden `xmin`/`xmax` system columns marking the transaction IDs that created and (if applicable) deleted/updated it. When a transaction starts, it gets a **snapshot** of running/committed transaction IDs. Reading a row checks its `xmin`/`xmax` against the snapshot to decide visibility. Updates create a new row version rather than overwriting in place (old versions become dead tuples, cleaned up by VACUUM). This lets readers and writers proceed concurrently without blocking each other.

---

## Top 20 "Must Master" Questions
*(For mid/senior Go backend interviews)*

1. ACID properties
2. Isolation levels
3. Deadlocks and prevention
4. MVCC
5. WAL (Write-Ahead Logging)
6. B+ Tree indexes
7. Composite indexes and the leftmost prefix rule
8. EXPLAIN ANALYZE
9. Query optimization techniques
10. Transactions in PostgreSQL
11. Optimistic vs pessimistic locking
12. Database connection pooling in Go
13. Prepared statements and SQL injection prevention
14. Read replicas and replication lag
15. Sharding vs partitioning
16. Redis caching patterns (Cache Aside, Write Through)
17. Soft delete vs hard delete
18. Zero-downtime database migrations
19. Idempotency for APIs and database writes
20. Outbox pattern for reliable event publishing
